package scope

import (
	"math/rand"
	"sort"
)

type ScopeNamespace struct {
	ClusterID     string
	NamespaceName string
}

func sortNamespaces(namespacesByCluster map[string][]string) ([]string, map[string][]string) {
	orderedClusters := make([]string, 0, len(namespacesByCluster))
	orderedNamespacesByCluster := make(map[string][]string, len(namespacesByCluster))
	for clusterID, namespaces := range namespacesByCluster {
		orderedClusters = append(orderedClusters, clusterID)
		orderedNamespaces := make([]string, 0, len(namespaces))
		for _, ns := range namespaces {
			orderedNamespaces = append(orderedNamespaces, ns)
		}
		sort.Strings(orderedNamespaces)
		orderedNamespacesByCluster[clusterID] = orderedNamespaces
	}
	sort.Strings(orderedClusters)
	return orderedClusters, orderedNamespacesByCluster
}

func shuffleNamespaces(namespacesByCluster map[string][]string) ([]string, map[string][]string) {
	numClusters := len(namespacesByCluster)
	numNamespaces := 0
	for _, namespaces := range namespacesByCluster {
		numNamespaces += len(namespaces)
	}
	randSource := rand.NewSource(int64(10000*numClusters + numNamespaces))
	randGen := rand.New(randSource)
	clusterIDs := make([]string, 0, len(namespacesByCluster))
	for clusterID, _ := range namespacesByCluster {
		clusterIDs = append(clusterIDs, clusterID)
	}
	sort.Strings(clusterIDs)
	shuffledClusterIDs := make([]string, 0, len(clusterIDs))
	for _, clusterID := range clusterIDs {
		shuffledClusterIDs = append(shuffledClusterIDs, clusterID)
	}
	for i := 0; i < numClusters; i++ {
		swapIx := int(randGen.Int31n(int32(numClusters - i)))
		if swapIx == 0 || i+swapIx >= numClusters {
			continue
		}
		shuffledClusterIDs[i], shuffledClusterIDs[i+swapIx] = shuffledClusterIDs[i+swapIx], shuffledClusterIDs[i]
	}
	shuffledNamespacesByClusterID := make(map[string][]string, len(shuffledClusterIDs))
	for _, clusterID := range shuffledClusterIDs {
		shuffledNamespaces := make([]string, 0, len(namespacesByCluster[clusterID]))
		for _, ns := range namespacesByCluster[clusterID] {
			shuffledNamespaces = append(shuffledNamespaces, ns)
		}
		numNamespaces := len(shuffledNamespaces)
		for i := 0; i < numNamespaces; i++ {
			swapIx := int(randGen.Int31n(int32(numNamespaces - i)))
			if swapIx == 0 || i+swapIx >= numNamespaces {
				continue
			}
			shuffledNamespaces[i], shuffledNamespaces[i+swapIx] = shuffledNamespaces[i+swapIx], shuffledNamespaces[i]
		}
		shuffledNamespacesByClusterID[clusterID] = shuffledNamespaces
	}
	return shuffledClusterIDs, shuffledNamespacesByClusterID
}

func selectNamespaces(clusterIDs []string, namespacesByCluster map[string][]string, sizes []int) [][]ScopeNamespace {
	output := make([][]ScopeNamespace, 0, len(sizes))
	for _, size := range sizes {
		selectedNamespaces := make([]ScopeNamespace, 0, size)
		selectedNamespaceCount := 0
		consumedClusters := 0
		for consumedClusters < len(clusterIDs) && selectedNamespaceCount < size {
			selectedNamespacesInCluster := 0
			clusterID := clusterIDs[consumedClusters]
			clusterNamespaces := namespacesByCluster[clusterID]
			for selectedNamespaceCount < size && selectedNamespacesInCluster < len(clusterNamespaces) {
				namespace := clusterNamespaces[selectedNamespacesInCluster]
				selectedNamespaces = append(selectedNamespaces, ScopeNamespace{ClusterID: clusterID, NamespaceName: namespace})
				selectedNamespacesInCluster++
				selectedNamespaceCount++
			}
			consumedClusters++
		}
		output = append(output, selectedNamespaces)
	}
	return output
}

func SelectNamespacesOrdered(namespacesByCluster map[string][]string, sizes []int) [][]ScopeNamespace {
	sortedClusterIDs, sortedNamespacesByCluster := sortNamespaces(namespacesByCluster)
	return selectNamespaces(sortedClusterIDs, sortedNamespacesByCluster, sizes)
}

func SelectNamespacesRandom(namespacesByCluster map[string][]string, sizes []int) [][]ScopeNamespace {
	_, sortedNamespacesByCluster := sortNamespaces(namespacesByCluster)
	shuffledClusterIDs, shuffledNamespacesByCluster := shuffleNamespaces(sortedNamespacesByCluster)
	return selectNamespaces(shuffledClusterIDs, shuffledNamespacesByCluster, sizes)
}
