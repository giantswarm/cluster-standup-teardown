package teardown

import (
	"context"

	"github.com/giantswarm/clustertest/v2"
	"github.com/giantswarm/clustertest/v2/pkg/application"
	"github.com/giantswarm/clustertest/v2/pkg/logger"
)

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
	return c.Framework.DeleteCluster(context.Background(), cluster)
}
