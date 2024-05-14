package teardown

import (
	"context"
	"os"
	"strings"

	"github.com/giantswarm/clustertest"
	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/logger"
)

const keepWCEnvVar = "E2E_WC_KEEP" //nolint:gosec

// Client is the client responsible for handling cluster teardown
type Client struct {
	Framework *clustertest.Framework
}

// New returns a new teardown client
func New(framework *clustertest.Framework) *Client {
	return &Client{
		Framework: framework,
	}
}

// Teardown handles removing the given workload cluster from the MC
func (c *Client) Teardown(cluster *application.Cluster) error {
	logger.Log("Deleting cluster: %s", cluster.Name)

	keep := strings.ToLower(os.Getenv(keepWCEnvVar))
	if keep != "" && keep != "false" {
		logger.Log("⚠️ The %s env var is set, skipping deletion of workload cluster", keepWCEnvVar)
		logger.Log("⚠️ This means the Cluster '%s' will remain on the management cluster only until the cluster-cleaner decides to remove it later. To disable the cluster-cleaner behavior please manually add the 'alpha.giantswarm.io/ignore-cluster-deletion' annotation to your test cluster.", cluster.Name)
		logger.Log("⚠️ Please be sure to manually delete the '%s' Organisation when you are finished.", cluster.Organization.Name)

		return nil
	}

	return c.Framework.DeleteCluster(context.Background(), cluster)
}
