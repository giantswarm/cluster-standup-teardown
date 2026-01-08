package values

import (
	"github.com/giantswarm/clustertest/v3/pkg/utils"
	"sigs.k8s.io/yaml"
)

type baseValues struct {
	Global global `json:"global" yaml:"global"`
}

type global struct {
	Metadata metadata `json:"metadata" yaml:"metadata"`
}

type metadata struct {
	Labels map[string]string `json:"labels" yaml:"labels"`
}

// BuildBaseValues generates a base set of cluster values that are relevant for all test clusters
func BuildBaseValues() string {
	v := baseValues{
		Global: global{
			Metadata: metadata{
				Labels: utils.GetBaseLabels(),
			},
		},
	}

	v.Global.Metadata.Labels["monitoring.giantswarm.io/prometheus-volume-size"] = "small"

	bytes, _ := yaml.Marshal(v)
	return string(bytes)
}
