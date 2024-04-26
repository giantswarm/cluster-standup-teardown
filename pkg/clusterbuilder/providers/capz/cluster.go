package capz

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/pkg/values"

	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/utils"
)

var (
	//go:embed values/cluster_values.yaml
	baseClusterValues string
	//go:embed values/default-apps_values.yaml
	baseDefaultAppsValues string
)

// ClusterBuilder is the CAPZ ClusterBuilder
type ClusterBuilder struct{}

// NewClusterApp builds a new CAPZ cluster App
func (c *ClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesFile string, defaultAppsValuesFile string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderAzure).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustOverlayValues(baseClusterValues, clusterValuesFile),
			values.MustOverlayValues(baseDefaultAppsValues, defaultAppsValuesFile),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		)
}