package capa

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/pkg/values"

	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/utils"
)

var (
	//go:embed values/managed-cluster_values.yaml
	baseManagedClusterValues string
	//go:embed values/managed-default-apps_values.yaml
	baseManagedDefaultAppsValues string
)

// ClusterBuilder is the CAPA EKS ClusterBuilder
type ManagedClusterBuilder struct{}

// NewClusterApp builds a new CAPA EKS cluster App
func (c *ManagedClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesFile string, defaultAppsValuesFile string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderEKS).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustOverlayValues(baseManagedClusterValues, clusterValuesFile),
			values.MustOverlayValues(baseManagedDefaultAppsValues, defaultAppsValuesFile),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		)
}
