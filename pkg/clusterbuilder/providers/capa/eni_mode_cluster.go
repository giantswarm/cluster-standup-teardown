package capa

import (
	_ "embed"

	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/utils"

	"github.com/giantswarm/cluster-standup-teardown/pkg/values"
)

var (
	//go:embed values/cilium-eni-mode-cluster_values.yaml
	baseCiliumEniModeClusterValues string
	//go:embed values/cilium-eni-mode-default-apps_values.yaml
	baseCiliumEniModeDefaultAppsValues string
)

// CiliumEniModeClusterBuilder is the CAPA ClusterBuilder for Cilium ENI mode
type CiliumEniModeClusterBuilder struct{}

// NewClusterApp builds a new CAPA cluster App for Cilium ENI mode
func (c *CiliumEniModeClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string, defaultAppsValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderAWS).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustMergeValues(append([]string{baseCiliumEniModeClusterValues}, clusterValuesOverrides...)...),
			values.MustMergeValues(append([]string{baseCiliumEniModeDefaultAppsValues}, defaultAppsValuesOverrides...)...),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		)
}

// KubeContext returns the known KubeConfig context that this builder expects
func (c *CiliumEniModeClusterBuilder) KubeContext() string {
	return "capa"
}
