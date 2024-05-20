package values

import (
	"os"

	"sigs.k8s.io/yaml"
)

type baseValues struct {
	Global global `json:"global",yaml:"global"`
}

type global struct {
	Metadata metadata `json:"metadata",yaml:"metadata"`
}

type metadata struct {
	Labels map[string]string `json:"labels",yaml:"labels"`
}

// BuildBaseValues generates a base set of cluster values that are relevant for all test clusters
func BuildBaseValues() string {
	v := baseValues{
		Global: global{
			Metadata: metadata{
				Labels: map[string]string{
					"monitoring.giantswarm.io/prometheus-volume-size": "small",
				},
			},
		},
	}

	// If found, populate details about Tekton run as labels
	if os.Getenv("TEKTON_PIPELINE_RUN") != "" {
		v.Global.Metadata.Labels["cicd.giantswarm.io/pipelinerun"] = os.Getenv("TEKTON_PIPELINE_RUN")
	}
	if os.Getenv("TEKTON_TASK_RUN") != "" {
		v.Global.Metadata.Labels["cicd.giantswarm.io/taskrun"] = os.Getenv("TEKTON_TASK_RUN")
	}

	bytes, _ := yaml.Marshal(v)
	return string(bytes)
}
