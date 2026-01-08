package clusterbuilder

import (
	"reflect"
	"testing"

	"github.com/giantswarm/cluster-standup-teardown/v3/pkg/clusterbuilder/providers/capa"
	"github.com/giantswarm/cluster-standup-teardown/v3/pkg/clusterbuilder/providers/capv"
	"github.com/giantswarm/cluster-standup-teardown/v3/pkg/clusterbuilder/providers/capvcd"
	"github.com/giantswarm/cluster-standup-teardown/v3/pkg/clusterbuilder/providers/capz"
)

func Test_GetClusterBuilderForContext(t *testing.T) {
	tests := []struct {
		inputValues   string
		expected      ClusterBuilder
		expectedError bool
	}{
		{
			inputValues:   "capa",
			expected:      &capa.ClusterBuilder{},
			expectedError: false,
		},
		{
			inputValues:   "capa-private-proxy",
			expected:      &capa.PrivateClusterBuilder{},
			expectedError: false,
		},
		{
			inputValues:   "eks",
			expected:      &capa.ManagedClusterBuilder{},
			expectedError: false,
		},
		{
			inputValues:   "capv",
			expected:      &capv.ClusterBuilder{},
			expectedError: false,
		},
		{
			inputValues:   "capvcd",
			expected:      &capvcd.ClusterBuilder{},
			expectedError: false,
		},
		{
			inputValues:   "capz",
			expected:      &capz.ClusterBuilder{},
			expectedError: false,
		},
		{
			inputValues:   "CAPA",
			expected:      &capa.ClusterBuilder{},
			expectedError: false,
		},
		{
			inputValues:   "unknown",
			expected:      nil,
			expectedError: true,
		},
		{
			inputValues:   "",
			expected:      nil,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.inputValues, func(t *testing.T) {
			cb, err := GetClusterBuilderForContext(tc.inputValues)

			if err != nil && !tc.expectedError {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(cb, tc.expected) {
				t.Fatalf("Actual value didn't match expected value\n\nexpected: %q\n\nactual: %q", tc.expected, cb)
			}
		})
	}
}
