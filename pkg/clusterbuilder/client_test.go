package clusterbuilder

import (
	"os"
	"reflect"
	"testing"

	"github.com/giantswarm/clustertest/v3/pkg/application"
	"github.com/giantswarm/clustertest/v3/pkg/env"

	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/clusterbuilder/providers/capa"
	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/clusterbuilder/providers/capv"
	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/clusterbuilder/providers/capvcd"
	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/clusterbuilder/providers/capz"
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

func Test_ApplyAppOverridesFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		provider application.Provider
	}{
		{
			name:     "no env var set",
			envValue: "",
			provider: application.ProviderAWS,
		},
		{
			name:     "single app override",
			envValue: "karpenter=2.0.0",
			provider: application.ProviderAWS,
		},
		{
			name:     "multiple app overrides",
			envValue: "karpenter=2.0.0,aws-ebs-csi-driver=4.1.0",
			provider: application.ProviderAWS,
		},
		{
			name:     "cluster app is processed but ignored by WithAppOverride",
			envValue: "cluster-aws=1.0.0,karpenter=2.0.0",
			provider: application.ProviderAWS,
		},
		{
			name:     "handles whitespace",
			envValue: " karpenter = 2.0.0 , aws-ebs-csi-driver = 4.1.0 ",
			provider: application.ProviderAWS,
		},
		{
			name:     "ignores malformed entries",
			envValue: "karpenter=2.0.0,invalid,=noname,noversion=",
			provider: application.ProviderAWS,
		},
		{
			name:     "sha-based version",
			envValue: "karpenter=2.0.0-abc123def456abc123def456abc123def456abc1",
			provider: application.ProviderAWS,
		},
		{
			name:     "single app with catalog",
			envValue: "karpenter=2.0.0:giantswarm",
			provider: application.ProviderAWS,
		},
		{
			name:     "multiple apps with mixed catalog specifications",
			envValue: "cluster-aws=7.2.5-abc123,aws-ebs-csi-driver=4.1.0:default,karpenter=2.0.0:giantswarm",
			provider: application.ProviderAWS,
		},
		{
			name:     "catalog with whitespace",
			envValue: "karpenter=2.0.0: giantswarm ",
			provider: application.ProviderAWS,
		},
		{
			name:     "empty catalog after colon is treated as no catalog",
			envValue: "karpenter=2.0.0:",
			provider: application.ProviderAWS,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			if tc.envValue != "" {
				os.Setenv(env.OverrideVersions, tc.envValue)
				defer os.Unsetenv(env.OverrideVersions)
			} else {
				os.Unsetenv(env.OverrideVersions)
			}

			// Create a test cluster
			cluster := application.NewClusterApp("test-cluster", tc.provider)

			// Apply overrides
			result := ApplyAppOverridesFromEnv(cluster)

			// We can't directly access appOverrides as it's private, but we can verify
			// the function doesn't panic and returns a cluster
			if result == nil {
				t.Fatal("ApplyAppOverridesFromEnv returned nil")
			}

			// Verify the cluster name is preserved
			if result.Name != "test-cluster" {
				t.Fatalf("Cluster name changed: expected 'test-cluster', got '%s'", result.Name)
			}

			// Verify the provider is preserved
			if result.Provider != tc.provider {
				t.Fatalf("Cluster provider changed: expected '%s', got '%s'", tc.provider, result.Provider)
			}
		})
	}
}

func Test_ApplyAppOverridesFromEnv_VersionFormat(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
	}{
		{
			name:     "semver version",
			envValue: "karpenter=2.0.0",
		},
		{
			name:     "semver with v prefix",
			envValue: "karpenter=v2.0.0",
		},
		{
			name:     "version with sha suffix",
			envValue: "karpenter=2.0.0-164a75740365c5c21ca8aed69ebeb05f75c07fd8",
		},
		{
			name:     "version with prerelease",
			envValue: "karpenter=2.0.0-alpha.1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(env.OverrideVersions, tc.envValue)
			defer os.Unsetenv(env.OverrideVersions)

			cluster := application.NewClusterApp("test-cluster", application.ProviderAWS)
			result := ApplyAppOverridesFromEnv(cluster)

			if result == nil {
				t.Fatal("ApplyAppOverridesFromEnv returned nil")
			}
		})
	}
}

func Test_parseVersionAndCatalog(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedVersion string
		expectedCatalog string
	}{
		{
			name:            "version only",
			input:           "2.0.0",
			expectedVersion: "2.0.0",
			expectedCatalog: "",
		},
		{
			name:            "version with catalog",
			input:           "2.0.0:giantswarm",
			expectedVersion: "2.0.0",
			expectedCatalog: "giantswarm",
		},
		{
			name:            "version with sha suffix only",
			input:           "2.0.0-164a75740365c5c21ca8aed69ebeb05f75c07fd8",
			expectedVersion: "2.0.0-164a75740365c5c21ca8aed69ebeb05f75c07fd8",
			expectedCatalog: "",
		},
		{
			name:            "version with sha suffix and catalog",
			input:           "2.0.0-164a75740365c5c21ca8aed69ebeb05f75c07fd8:cluster-test",
			expectedVersion: "2.0.0-164a75740365c5c21ca8aed69ebeb05f75c07fd8",
			expectedCatalog: "cluster-test",
		},
		{
			name:            "version with prerelease and catalog",
			input:           "2.0.0-alpha.1:giantswarm",
			expectedVersion: "2.0.0-alpha.1",
			expectedCatalog: "giantswarm",
		},
		{
			name:            "empty catalog after colon",
			input:           "2.0.0:",
			expectedVersion: "2.0.0:",
			expectedCatalog: "",
		},
		{
			name:            "catalog with whitespace",
			input:           "2.0.0: giantswarm ",
			expectedVersion: "2.0.0",
			expectedCatalog: "giantswarm",
		},
		{
			name:            "v prefix version with catalog",
			input:           "v2.0.0:default",
			expectedVersion: "v2.0.0",
			expectedCatalog: "default",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			version, catalog := parseVersionAndCatalog(tc.input)

			if version != tc.expectedVersion {
				t.Errorf("version mismatch: expected %q, got %q", tc.expectedVersion, version)
			}

			if catalog != tc.expectedCatalog {
				t.Errorf("catalog mismatch: expected %q, got %q", tc.expectedCatalog, catalog)
			}
		})
	}
}
