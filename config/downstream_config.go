package config

type DownstreamConfig struct {
	TaskServiceURL     string `yaml:"task_service_url"`
	ResourceServiceURL string `yaml:"resource_service_url"`
}

func (cfg *DownstreamConfig) merge(next DownstreamConfig, section sectionValues) {
	if section.has("task_service_url") {
		cfg.TaskServiceURL = next.TaskServiceURL
	}
	if section.has("resource_service_url") {
		cfg.ResourceServiceURL = next.ResourceServiceURL
	}
}
