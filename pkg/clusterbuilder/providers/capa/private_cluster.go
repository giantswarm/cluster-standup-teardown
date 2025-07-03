package capa

import (
	_ "embed"
	"fmt"
	"math/rand"

	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/utils"

	"github.com/giantswarm/cluster-standup-teardown/pkg/values"
)

var (
	//go:embed values/private-cluster_values.yaml
	basePrivateClusterValues string
	//go:embed values/private-default-apps_values.yaml
	basePrivateDefaultAppsValues string
)

// PrivateClusterBuilder is the private CAPA ClusterBuilder
type PrivateClusterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new private CAPA cluster App
func (c *PrivateClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string, defaultAppsValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	// WC CIDRs have to not overlap and be in the 10.129. - 10.159. range, so
	// we select a random number in that range and set it as the second octet.
	//nolint:gosec
	randomOctet := rand.Intn(31) + 129 // Generates a number between 129 and 159 inclusive.
	cidrOctet := fmt.Sprintf("%d", randomOctet)
	templateValues := &application.TemplateValues{
		ClusterName:  clusterName,
		Organization: orgName,
		ExtraValues: map[string]string{
			"CIDRSecondOctet": cidrOctet,
		},
	}

	return application.NewClusterApp(clusterName, application.ProviderAWS).
		WithOrg(organization.New(orgName)).
		WithAppValues(

			values.MustMergeValues(append([]string{basePrivateClusterValues}, clusterValuesOverrides...)...),
			values.MustMergeValues(append([]string{basePrivateDefaultAppsValues}, defaultAppsValuesOverrides...)...),
			templateValues,
		)
}

// KubeContext returns the known KubeConfig context that this builder expects
func (c *PrivateClusterBuilder) KubeContext() string {
	if c.CustomKubeContext != "" {
		return c.CustomKubeContext
	}
	return "capa-private-proxy"
}
