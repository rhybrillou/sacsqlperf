package scope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortNamespaces(t *testing.T) {
	testCases := []struct {
		name                              string
		namespacesByCluster               map[string][]string
		expectedSortedClusterIDs          []string
		expectedSortedNamespacesByCluster map[string][]string
	}{
		{
			name: "1 cluster, 1 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA"},
			},
			expectedSortedClusterIDs: []string{"Cluster1"},
			expectedSortedNamespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA"},
			},
		},
		{
			name: "2 cluster, 1 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster2": []string{},
				"Cluster1": []string{"namespaceA"},
			},
			expectedSortedClusterIDs: []string{"Cluster1", "Cluster2"},
			expectedSortedNamespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA"},
				"Cluster2": []string{},
			},
		},
		{
			name: "2 cluster, 10 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster2": []string{"namespaceJ", "namespaceG", "namespaceA", "namespaceE", "namespaceH", "namespaceI", "namespaceC"},
				"Cluster1": []string{"namespaceD", "namespaceF", "namespaceB"},
			},
			expectedSortedClusterIDs: []string{"Cluster1", "Cluster2"},
			expectedSortedNamespacesByCluster: map[string][]string{
				"Cluster2": []string{"namespaceA", "namespaceC", "namespaceE", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster1": []string{"namespaceB", "namespaceD", "namespaceF"},
			},
		},
		{
			name: "1 cluster, 10 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceB", "namespaceH", "namespaceF", "namespaceI", "namespaceA", "namespaceD", "namespaceG", "namespaceJ", "namespaceC", "namespaceE"},
			},
			expectedSortedClusterIDs: []string{"Cluster1"},
			expectedSortedNamespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
			},
		},
		{
			name: "2 cluster, 4-16 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceD", "namespaceA", "namespaceB", "namespaceC"},
				"Cluster2": []string{"namespaceO", "namespaceG", "namespaceS", "namespaceT", "namespaceE", "namespaceH", "namespaceL", "namespaceP", "namespaceF", "namespaceQ", "namespaceN", "namespaceK", "namespaceM", "namespaceJ", "namespaceI", "namespaceR"},
			},
			expectedSortedClusterIDs: []string{"Cluster1", "Cluster2"},
			expectedSortedNamespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD"},
				"Cluster2": []string{"namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ", "namespaceK", "namespaceL", "namespaceM", "namespaceN", "namespaceO", "namespaceP", "namespaceQ", "namespaceR", "namespaceS", "namespaceT"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(it *testing.T) {
			sortedClusterIDs, sortedNamespacesByClusterID := sortNamespaces(tc.namespacesByCluster)
			assert.Equal(t, tc.expectedSortedClusterIDs, sortedClusterIDs)
			assert.Equal(t, tc.expectedSortedNamespacesByCluster, sortedNamespacesByClusterID)
		})
	}
}

func TestShuffleNamespaces(t *testing.T) {
	testCases := []struct {
		name                                string
		namespacesByCluster                 map[string][]string
		expectedShuffledClusterIDs          []string
		expectedShuffledNamespacesByCluster map[string][]string
	}{
		{
			name: "1 cluster, 1 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA"},
			},
			expectedShuffledClusterIDs: []string{"Cluster1"},
			expectedShuffledNamespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA"},
			},
		},
		{
			name: "2 cluster, 1 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA"},
				"Cluster2": []string{},
			},
			expectedShuffledClusterIDs: []string{"Cluster1", "Cluster2"},
			expectedShuffledNamespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA"},
				"Cluster2": []string{},
			},
		},
		{
			name: "2 cluster, 10 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceB", "namespaceD", "namespaceF"},
				"Cluster2": []string{"namespaceA", "namespaceC", "namespaceE", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
			},
			expectedShuffledClusterIDs: []string{"Cluster2", "Cluster1"},
			expectedShuffledNamespacesByCluster: map[string][]string{
				"Cluster2": []string{"namespaceJ", "namespaceE", "namespaceG", "namespaceA", "namespaceH", "namespaceC", "namespaceI"},
				"Cluster1": []string{"namespaceF", "namespaceB", "namespaceD"},
			},
		},
		{
			name: "1 cluster, 10 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
			},
			expectedShuffledClusterIDs: []string{"Cluster1"},
			expectedShuffledNamespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceE", "namespaceA", "namespaceD", "namespaceI", "namespaceC", "namespaceG", "namespaceJ", "namespaceB", "namespaceH", "namespaceF"},
			},
		},
		{
			name: "2 cluster, 4-16 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD"},
				"Cluster2": []string{"namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ", "namespaceK", "namespaceL", "namespaceM", "namespaceN", "namespaceO", "namespaceP", "namespaceQ", "namespaceR", "namespaceS", "namespaceT"},
			},
			expectedShuffledClusterIDs: []string{"Cluster1", "Cluster2"},
			expectedShuffledNamespacesByCluster: map[string][]string{
				"Cluster1": []string{"namespaceD", "namespaceA", "namespaceB", "namespaceC"},
				"Cluster2": []string{"namespaceO", "namespaceG", "namespaceS", "namespaceT", "namespaceE", "namespaceH", "namespaceL", "namespaceP", "namespaceF", "namespaceQ", "namespaceN", "namespaceK", "namespaceM", "namespaceJ", "namespaceI", "namespaceR"},
			},
		},
		{
			name: "10 cluster, 100 namespace",
			namespacesByCluster: map[string][]string{
				"Cluster1":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster2":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster3":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster4":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster5":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster6":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster7":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster8":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster9":  []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
				"Cluster10": []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD", "namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ"},
			},
			expectedShuffledClusterIDs: []string{"Cluster10", "Cluster3", "Cluster8", "Cluster6", "Cluster7", "Cluster4", "Cluster5", "Cluster2", "Cluster9", "Cluster1"},
			expectedShuffledNamespacesByCluster: map[string][]string{
				"Cluster1":  []string{"namespaceF", "namespaceA", "namespaceI", "namespaceG", "namespaceC", "namespaceD", "namespaceB", "namespaceJ", "namespaceE", "namespaceH"},
				"Cluster2":  []string{"namespaceB", "namespaceF", "namespaceJ", "namespaceH", "namespaceE", "namespaceG", "namespaceC", "namespaceA", "namespaceD", "namespaceI"},
				"Cluster3":  []string{"namespaceJ", "namespaceE", "namespaceA", "namespaceC", "namespaceI", "namespaceF", "namespaceG", "namespaceH", "namespaceD", "namespaceB"},
				"Cluster4":  []string{"namespaceE", "namespaceG", "namespaceJ", "namespaceH", "namespaceC", "namespaceB", "namespaceF", "namespaceD", "namespaceI", "namespaceA"},
				"Cluster5":  []string{"namespaceF", "namespaceH", "namespaceE", "namespaceA", "namespaceJ", "namespaceC", "namespaceG", "namespaceD", "namespaceI", "namespaceB"},
				"Cluster6":  []string{"namespaceA", "namespaceE", "namespaceD", "namespaceJ", "namespaceB", "namespaceI", "namespaceH", "namespaceG", "namespaceF", "namespaceC"},
				"Cluster7":  []string{"namespaceH", "namespaceF", "namespaceE", "namespaceG", "namespaceJ", "namespaceC", "namespaceB", "namespaceA", "namespaceI", "namespaceD"},
				"Cluster8":  []string{"namespaceF", "namespaceJ", "namespaceE", "namespaceD", "namespaceC", "namespaceG", "namespaceA", "namespaceI", "namespaceH", "namespaceB"},
				"Cluster9":  []string{"namespaceB", "namespaceD", "namespaceG", "namespaceI", "namespaceJ", "namespaceA", "namespaceE", "namespaceH", "namespaceC", "namespaceF"},
				"Cluster10": []string{"namespaceB", "namespaceD", "namespaceC", "namespaceI", "namespaceF", "namespaceE", "namespaceJ", "namespaceG", "namespaceA", "namespaceH"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(it *testing.T) {
			shuffledClusterIDs, shuffledNamespacesByClusterID := shuffleNamespaces(tc.namespacesByCluster)
			assert.Equal(it, tc.expectedShuffledClusterIDs, shuffledClusterIDs)
			assert.Equal(it, tc.expectedShuffledNamespacesByCluster, shuffledNamespacesByClusterID)
		})
	}
}

func TestSelectNamespacesOrdered(t *testing.T) {
	randomNamespacesByCluster := map[string][]string{
		"Cluster1": []string{"namespaceD", "namespaceA", "namespaceB", "namespaceC"},
		"Cluster2": []string{"namespaceO", "namespaceG", "namespaceS", "namespaceT", "namespaceE", "namespaceH", "namespaceL", "namespaceP", "namespaceF", "namespaceQ", "namespaceN", "namespaceK", "namespaceM", "namespaceJ", "namespaceI", "namespaceR"},
	}
	scopeNamespaceA := ScopeNamespace{ClusterID: "Cluster1", NamespaceName: "namespaceA"}
	scopeNamespaceB := ScopeNamespace{ClusterID: "Cluster1", NamespaceName: "namespaceB"}
	scopeNamespaceC := ScopeNamespace{ClusterID: "Cluster1", NamespaceName: "namespaceC"}
	scopeNamespaceD := ScopeNamespace{ClusterID: "Cluster1", NamespaceName: "namespaceD"}
	scopeNamespaceE := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceE"}
	scopeNamespaceF := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceF"}
	scopeNamespaceG := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceG"}
	scopeNamespaceH := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceH"}
	scopeNamespaceI := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceI"}
	scopeNamespaceJ := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceJ"}
	scopeNamespaceK := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceK"}
	scopeNamespaceL := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceL"}
	scopeNamespaceM := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceM"}
	scopeNamespaceN := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceN"}
	scopeNamespaceO := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceO"}
	scopeNamespaceP := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceP"}
	scopeNamespaceQ := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceQ"}
	scopeNamespaceR := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceR"}
	scopeNamespaceS := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceS"}
	scopeNamespaceT := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceT"}
	expectedSelectedNamespaces := [][]ScopeNamespace{
		{scopeNamespaceA},
		{scopeNamespaceA, scopeNamespaceB, scopeNamespaceC, scopeNamespaceD, scopeNamespaceE},
		{
			scopeNamespaceA, scopeNamespaceB, scopeNamespaceC, scopeNamespaceD, scopeNamespaceE,
			scopeNamespaceF, scopeNamespaceG, scopeNamespaceH, scopeNamespaceI, scopeNamespaceJ,
			scopeNamespaceK, scopeNamespaceL, scopeNamespaceM, scopeNamespaceN, scopeNamespaceO,
			scopeNamespaceP, scopeNamespaceQ, scopeNamespaceR, scopeNamespaceS, scopeNamespaceT,
		},
	}
	scopeSizes := []int{1, 5, 100}
	selectedNamespaces := SelectNamespacesOrdered(randomNamespacesByCluster, scopeSizes)
	assert.Equal(t, expectedSelectedNamespaces, selectedNamespaces)
}

func TestSelectNamespacesRandom(t *testing.T) {
	randomNamespacesByCluster := map[string][]string{
		"Cluster1": []string{"namespaceA", "namespaceB", "namespaceC", "namespaceD"},
		"Cluster2": []string{"namespaceE", "namespaceF", "namespaceG", "namespaceH", "namespaceI", "namespaceJ", "namespaceK", "namespaceL", "namespaceM", "namespaceN", "namespaceO", "namespaceP", "namespaceQ", "namespaceR", "namespaceS", "namespaceT"},
	}
	scopeNamespaceA := ScopeNamespace{ClusterID: "Cluster1", NamespaceName: "namespaceA"}
	scopeNamespaceB := ScopeNamespace{ClusterID: "Cluster1", NamespaceName: "namespaceB"}
	scopeNamespaceC := ScopeNamespace{ClusterID: "Cluster1", NamespaceName: "namespaceC"}
	scopeNamespaceD := ScopeNamespace{ClusterID: "Cluster1", NamespaceName: "namespaceD"}
	scopeNamespaceE := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceE"}
	scopeNamespaceF := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceF"}
	scopeNamespaceG := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceG"}
	scopeNamespaceH := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceH"}
	scopeNamespaceI := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceI"}
	scopeNamespaceJ := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceJ"}
	scopeNamespaceK := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceK"}
	scopeNamespaceL := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceL"}
	scopeNamespaceM := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceM"}
	scopeNamespaceN := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceN"}
	scopeNamespaceO := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceO"}
	scopeNamespaceP := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceP"}
	scopeNamespaceQ := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceQ"}
	scopeNamespaceR := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceR"}
	scopeNamespaceS := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceS"}
	scopeNamespaceT := ScopeNamespace{ClusterID: "Cluster2", NamespaceName: "namespaceT"}
	expectedSelectedNamespaces := [][]ScopeNamespace{
		{scopeNamespaceD},
		{scopeNamespaceD, scopeNamespaceA, scopeNamespaceB, scopeNamespaceC, scopeNamespaceO},
		{
			scopeNamespaceD, scopeNamespaceA, scopeNamespaceB, scopeNamespaceC, scopeNamespaceO,
			scopeNamespaceG, scopeNamespaceS, scopeNamespaceT, scopeNamespaceE, scopeNamespaceH,
			scopeNamespaceL, scopeNamespaceP, scopeNamespaceF, scopeNamespaceQ, scopeNamespaceN,
			scopeNamespaceK, scopeNamespaceM, scopeNamespaceJ, scopeNamespaceI, scopeNamespaceR,
		},
	}
	scopeSizes := []int{1, 5, 100}
	selectedNamespaces := SelectNamespacesRandom(randomNamespacesByCluster, scopeSizes)
	assert.Equal(t, expectedSelectedNamespaces, selectedNamespaces)
}
