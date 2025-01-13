package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rhybrillou/sacsqlperf/src/pkg/db"
	"github.com/rhybrillou/sacsqlperf/src/pkg/query"
	"github.com/rhybrillou/sacsqlperf/src/pkg/scope"
)

var (
	scopeSizes = []int{10, 20, 50, 100, 200, 500, 1000, 2000}

	testedQueries = []*query.Query{
		{
			Statement: "select",
			StatementTargets: []string{
				"distinct(images.Id) as Image_Sha",
				"images.RiskScore as image_risk_score",
			},
			TargetTables: []string{"images"},
			InnerJoins: []query.InnerJoin{
				{
					Left:  query.QualifiedColumn{TableName: "images", ColumnName: "Id"},
					Right: query.QualifiedColumn{TableName: "deployments_containers", ColumnName: "Image_Id"},
				},
				{
					Left:  query.QualifiedColumn{TableName: "deployments_containers", ColumnName: "deployments_Id"},
					Right: query.QualifiedColumn{TableName: "deployments", ColumnName: "Id"},
				},
			},
			OrderBy: []query.OrderColumn{
				{Reversed: true, Column: query.QualifiedColumn{TableName: "images", ColumnName: "RiskScore"}},
			},
			QueryPagination: &query.Pagination{
				Limit: 6,
			},
			ScopeLevel:           "namespace",
			ScopeTable:           "deployments",
			ScopeClusterColumn:   "ClusterId",
			ScopeNamespaceColumn: "Namespace",
		}, /*
			{
				Statement: "select",
				StatementTargets: []string{
					"count(*)",
				},
				TargetTables: []string{"alerts"},
				InnerJoins:   []query.InnerJoin{},
				WhereClause: &query.WcAnd{
					Operands: []query.WhereClausePart{
						&query.QualifiedColumn{TableName: "alerts", ColumnName: "Policy_Severity", Value: 3},
						&query.WcOr{
							Operands: []query.WhereClausePart{
								&query.QualifiedColumn{TableName: "alerts", ColumnName: "State", Value: 0},
								&query.QualifiedColumn{TableName: "alerts", ColumnName: "State", Value: 3},
							},
						},
					},
				},
				ScopeLevel:           "namespace",
				ScopeTable:           "alerts",
				ScopeClusterColumn:   "ClusterId",
				ScopeNamespaceColumn: "Namespace",
			}, */
		{
			Statement: "select",
			StatementTargets: []string{
				"policy_severity",
				"count(*)",
			},
			TargetTables: []string{"alerts"},
			InnerJoins:   []query.InnerJoin{},
			WhereClause: &query.WcAnd{
				Operands: []query.WhereClausePart{
					&query.WcOr{
						Operands: []query.WhereClausePart{
							&query.QualifiedColumn{TableName: "alerts", ColumnName: "State", Value: 0},
							&query.QualifiedColumn{TableName: "alerts", ColumnName: "State", Value: 3},
						},
					},
				},
			},
			GroupBy: []query.QualifiedColumn{
				{TableName: "alerts", ColumnName: "Policy_Severity"},
			},
			ScopeLevel:           "namespace",
			ScopeTable:           "alerts",
			ScopeClusterColumn:   "ClusterId",
			ScopeNamespaceColumn: "Namespace",
		},
	}
)

func main() {
	// Ensure the logs are available for a while after the execution completed.
	defer done()
	ctx := context.Background()
	_ = ctx
	fmt.Println("Starting SQL performance tests")

	db, err := db.GetDBConn(ctx)
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
	orderedScopeNamespaces := scope.SelectNamespacesOrdered(namespacesByCluster, scopeSizes)
	pseudoRandomScopeNamespaces := scope.SelectNamespacesRandom(namespacesByCluster, scopeSizes)
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

func explain(ctx context.Context, db *pgxpool.Pool, request *query.Query) error {
	stmt, bindValues := request.ForExecution()
	explainRows, err := db.Query(ctx, fmt.Sprintf("explain (verbose, analyze, buffers, settings) %s", stmt), bindValues...)
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

func injectSACFilter(request *query.Query, scope []scope.ScopeNamespace) *query.Query {
	if request == nil {
		return nil
	}
	if len(scope) <= 0 {
		return request
	}
	if request.ScopeLevel != "cluster" && request.ScopeLevel != "namespace" {
		return request
	}
	namespacesByCluster := make(map[string][]string, 0)
	for _, ns := range scope {
		namespacesByCluster[ns.ClusterID] = append(namespacesByCluster[ns.ClusterID], ns.NamespaceName)
	}
	whereClusters := make([]query.WhereClausePart, 0, len(namespacesByCluster))
	for clusterID, namespaces := range namespacesByCluster {
		clusterColumnPart := &query.QualifiedColumn{
			TableName:  request.ScopeTable,
			ColumnName: request.ScopeClusterColumn,
			Value:      clusterID,
		}
		switch request.ScopeLevel {
		case "cluster":
			whereClusters = append(whereClusters, clusterColumnPart)
		case "namespace":
			whereClusterNamespaces := make([]query.WhereClausePart, 0, len(namespaces))
			for _, ns := range namespaces {
				namespaceColumnPart := &query.QualifiedColumn{
					TableName:  request.ScopeTable,
					ColumnName: request.ScopeNamespaceColumn,
					Value:      ns,
				}
				whereClusterNamespaces = append(whereClusterNamespaces, namespaceColumnPart)
			}
			whereClusters = append(whereClusters, &query.WcAnd{
				Operands: []query.WhereClausePart{
					clusterColumnPart,
					&query.WcOr{
						Operands: whereClusterNamespaces,
					},
				},
			})
		default:
			continue
		}
	}
	result := &query.Query{
		Statement:            request.Statement,
		StatementTargets:     request.StatementTargets,
		TargetTables:         request.TargetTables,
		InnerJoins:           request.InnerJoins,
		OrderBy:              request.OrderBy,
		GroupBy:              request.GroupBy,
		QueryPagination:      request.QueryPagination,
		ScopeLevel:           request.ScopeLevel,
		ScopeTable:           request.ScopeTable,
		ScopeClusterColumn:   request.ScopeClusterColumn,
		ScopeNamespaceColumn: request.ScopeNamespaceColumn,
	}
	if request.WhereClause != nil {
		result.WhereClause = &query.WcAnd{
			Operands: []query.WhereClausePart{
				&query.WcOr{
					Operands: whereClusters,
				},
				request.WhereClause,
			},
		}
	} else {
		result.WhereClause = &query.WcOr{
			Operands: whereClusters,
		}
	}
	return result
}
