package domain

type ConfigurationLoader interface {
	LoadForProject(baseConfig *Config, projectRoot string) (*Config, error)
}
