package capa

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/v3/pkg/values"

	"github.com/giantswarm/clustertest/v3/pkg/application"
	"github.com/giantswarm/clustertest/v3/pkg/organization"
	"github.com/giantswarm/clustertest/v3/pkg/utils"
)

var (
	//go:embed values/managed-cluster_values.yaml
	baseManagedClusterValues string
	//go:embed values/managed-default-apps_values.yaml
	baseManagedDefaultAppsValues string
)

// ManagedClusterBuilder is the CAPA EKS ClusterBuilder
type ManagedClusterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new CAPA EKS cluster App
func (c *ManagedClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string, defaultAppsValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderEKS).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustMergeValues(append([]string{baseManagedClusterValues}, clusterValuesOverrides...)...),
			values.MustMergeValues(append([]string{baseManagedDefaultAppsValues}, defaultAppsValuesOverrides...)...),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		)
}

// KubeContext returns the known KubeConfig context that this builder expects
func (c *ManagedClusterBuilder) KubeContext() string {
	if c.CustomKubeContext != "" {
		return c.CustomKubeContext
	}
	return "eks"
}
