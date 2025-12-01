package capv

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/v2/pkg/values"

	applicationv1alpha1 "github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/clustertest/v2/pkg/application"
	"github.com/giantswarm/clustertest/v2/pkg/organization"
	"github.com/giantswarm/clustertest/v2/pkg/utils"
)

const (
	VSphereCredSecretName      = "vsphere-credentials" //nolint:gosec
	VSphereCredSecretNamespace = "org-giantswarm"
)

var (
	//go:embed values/cluster_values.yaml
	baseClusterValues string
	//go:embed values/default-apps_values.yaml
	baseDefaultAppsValues string
)

// ClusterBuilder is the CAPV ClusterBuilder
type ClusterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new CAPV cluster App
func (c *ClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string, defaultAppsValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderVSphere).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustMergeValues(append([]string{baseClusterValues}, clusterValuesOverrides...)...),
			values.MustMergeValues(append([]string{baseDefaultAppsValues}, defaultAppsValuesOverrides...)...),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		).
		WithExtraConfigs([]applicationv1alpha1.AppExtraConfig{
			{
				Kind:      "secret",
				Name:      VSphereCredSecretName,
				Namespace: VSphereCredSecretNamespace,
				Priority:  25,
			},
		})
}

// KubeContext returns the known KubeConfig context that this builder expects
func (c *ClusterBuilder) KubeContext() string {
	if c.CustomKubeContext != "" {
		return c.CustomKubeContext
	}
	return "capv"
}
