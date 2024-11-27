package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var (
	scopeSizes = []int{10, 20, 50, 100, 200, 500, 1000, 2000}

	testedQueries = []*query{
		{
			statement: "select",
			statementTargets: []string{
				"distinct(images.Id) as Image_Sha",
				"images.RiskScore as image_risk_score",
			},
			targetTables: []string{"images"},
			innerJoins: []innerJoin{
				{
					left:  qualifiedColumn{tableName: "images", columnName: "Id"},
					right: qualifiedColumn{tableName: "deployments_containers", columnName: "Image_Id"},
				},
				{
					left:  qualifiedColumn{tableName: "deployments_containers", columnName: "deployments_Id"},
					right: qualifiedColumn{tableName: "deployments", columnName: "Id"},
				},
			},
			orderBy: []orderColumn{
				{reversed: true, column: qualifiedColumn{tableName: "images", columnName: "RiskScore"}},
			},
			queryPagination: &pagination{
				limit: 6,
			},
			scopeLevel:           "namespace",
			scopeTable:           "deployments",
			scopeClusterColumn:   "ClusterId",
			scopeNamespaceColumn: "Namespace",
		},
	}
)

func main() {
	// Ensure the logs are available for a while after the execution completed.
	defer done()
	ctx := context.Background()
	_ = ctx
	fmt.Println("Starting SQL performance tests")

	db, err := getDBConn(ctx)
	if err != nil {
		fmt.Printf("Error getting DB connection: %v\n", err)
		return
	}
	defer db.Close()
	dbRow := db.QueryRow(ctx, "select current_database()")
	var dbName string
	err = dbRow.Scan(&dbName)
	fmt.Println("connected to", dbName)
	fmt.Println("Querying namespaces")
	const namespaceCountStatement = "select count(*) from namespaces"
	fmt.Println("Running", namespaceCountStatement)
	row := db.QueryRow(ctx, namespaceCountStatement)
	var namespaceCount int
	err = row.Scan(&namespaceCount)
	fmt.Printf("Found %d namespaces\n", namespaceCount)
	rows, err := db.Query(ctx, "select clusterid, name from namespaces")
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		fmt.Printf("Error querying namespaces: %v\n", err)
		return
	}
	fmt.Println("Namespace query complete")
	defer rows.Close()
	c := 0
	namespacesByCluster := make(map[string][]string, 0)
	for rows.Next() {
		c++
		// fmt.Print(".")
		var clusterID string
		var namespaceName string
		err = rows.Scan(&clusterID, &namespaceName)
		if err != nil {
			continue
		}
		if _, found := namespacesByCluster[clusterID]; !found {
			namespacesByCluster[clusterID] = make([]string, 0, namespaceCount)
		}
		namespacesByCluster[clusterID] = append(namespacesByCluster[clusterID], namespaceName)
	}
	fmt.Printf("Selected %d namespaces\n", c)
	for clusterID, namespaces := range namespacesByCluster {
		fmt.Printf("Found %d namespaces for cluster %q\n", len(namespaces), clusterID)
	}
	orderedScopeNamespaces := selectNamespacesOrdered(namespacesByCluster, scopeSizes)
	pseudoRandomScopeNamespaces := selectNamespacesRandom(namespacesByCluster, scopeSizes)
	_ = orderedScopeNamespaces
	_ = pseudoRandomScopeNamespaces

	for ix, q := range testedQueries {
		_ = ix
		fmt.Println("index", ix)
		stmt, _ := q.ForExecution()
		fmt.Println(stmt)
		err = explain(ctx, db, q)
		if err != nil {
			fmt.Printf("Error querying for execution plan: %v\n", err)
		}
		for _, scope := range orderedScopeNamespaces {
			fmt.Printf("Getting plan for %d ordered namespaces\n", len(scope))
			sq := injectSACFilter(q, scope)
			err = explain(ctx, db, sq)
			if err != nil {
				fmt.Printf("Error querying for execution plan: %v\n", err)
			}
		}
		for _, scope := range pseudoRandomScopeNamespaces {
			fmt.Printf("Getting plan for %d random namespaces\n", len(scope))
			sq := injectSACFilter(q, scope)
			err = explain(ctx, db, sq)
			if err != nil {
				fmt.Printf("Error querying for execution plan: %v\n", err)
			}
		}
	}
}

func explain(ctx context.Context, db *pgxpool.Pool, request *query) error {
	stmt, bindValues := request.ForExecution()
	explainRows, err := db.Query(ctx, fmt.Sprintf("explain %s", stmt), bindValues...)
	if err != nil {
		return err
	}
	defer explainRows.Close()
	for explainRows.Next() {
		var explainResult string
		err = explainRows.Scan(&explainResult)
		if err != nil {
			return err
		}
		fmt.Println(explainResult)
	}
	return nil
}

func done() {
	fmt.Println("Sleeping an hour")
	time.Sleep(time.Hour)
}

func injectSACFilter(request *query, scope []scopeNamespace) *query {
	if request == nil {
		return nil
	}
	if len(scope) <= 0 {
		return request
	}
	if request.scopeLevel != "cluster" && request.scopeLevel != "namespace" {
		return request
	}
	namespacesByCluster := make(map[string][]string, 0)
	for _, ns := range scope {
		namespacesByCluster[ns.ClusterID] = append(namespacesByCluster[ns.ClusterID], ns.NamespaceName)
	}
	whereClusters := make([]whereClausePart, 0, len(namespacesByCluster))
	for clusterID, namespaces := range namespacesByCluster {
		clusterColumnPart := &qualifiedColumn{
			tableName:  request.scopeTable,
			columnName: request.scopeClusterColumn,
			value:      clusterID,
		}
		switch request.scopeLevel {
		case "cluster":
			whereClusters = append(whereClusters, clusterColumnPart)
		case "namespace":
			whereClusterNamespaces := make([]whereClausePart, 0, len(namespaces))
			for _, ns := range namespaces {
				namespaceColumnPart := &qualifiedColumn{
					tableName:  request.scopeTable,
					columnName: request.scopeNamespaceColumn,
					value:      ns,
				}
				whereClusterNamespaces = append(whereClusterNamespaces, namespaceColumnPart)
			}
			whereClusters = append(whereClusters, &wcAnd{
				operands: []whereClausePart{
					clusterColumnPart,
					&wcOr{
						operands: whereClusterNamespaces,
					},
				},
			})
		default:
			continue
		}
	}
	result := &query{
		statement:            request.statement,
		statementTargets:     request.statementTargets,
		targetTables:         request.targetTables,
		innerJoins:           request.innerJoins,
		orderBy:              request.orderBy,
		queryPagination:      request.queryPagination,
		scopeLevel:           request.scopeLevel,
		scopeTable:           request.scopeTable,
		scopeClusterColumn:   request.scopeClusterColumn,
		scopeNamespaceColumn: request.scopeNamespaceColumn,
	}
	if request.whereClause != nil {
		result.whereClause = &wcAnd{
			operands: []whereClausePart{
				&wcOr{
					operands: whereClusters,
				},
				request.whereClause,
			},
		}
	} else {
		result.whereClause = &wcOr{
			operands: whereClusters,
		}
	}
	return result
}
