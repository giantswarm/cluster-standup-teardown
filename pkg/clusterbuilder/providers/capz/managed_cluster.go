package capz

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/v6/pkg/values"

	"github.com/giantswarm/clustertest/v5/pkg/application"
	"github.com/giantswarm/clustertest/v5/pkg/organization"
	"github.com/giantswarm/clustertest/v5/pkg/utils"
)

// ProviderAKS is the provider name used for the CAPZ managed (AKS) cluster app.
//
// The clustertest `application` package does not yet expose a constant for AKS,
// so we define it here. The provider name determines the cluster app to use
// (i.e. `cluster-aks`).
const ProviderAKS application.Provider = "aks"

var (
	//go:embed values/managed-cluster_values.yaml
	baseManagedClusterValues string
)

// ManagedClusterBuilder is the CAPZ AKS ClusterBuilder
type ManagedClusterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new CAPZ AKS cluster App
func (c *ManagedClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, ProviderAKS).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustMergeValues(append([]string{baseManagedClusterValues}, clusterValuesOverrides...)...),
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
	return "aks"
}
