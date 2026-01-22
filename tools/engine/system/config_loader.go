package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Josepavese/asdp/engine/domain"
	"gopkg.in/yaml.v3"
)

type ValidationReport struct {
	Errors   []string
	Warnings []string
}

// LoadConfig loads the ASDP configuration in a hierarchical manner:
// 1. Embedded Defaults (Hardcoded)
// 2. Global Config (~/.asdp/config.yaml) - Overrides defaults
// 3. Project Config (<root>/.asdp.yaml) - Overrides global
func LoadConfig(projectRoot string) (*domain.Config, error) {
	// 1. Start with Defaults
	cfg := domain.DefaultConfig()

	// 2. Load Global Config
	homeDir, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Join(homeDir, ".asdp", "config.yaml")
		if _, err := os.Stat(globalPath); err == nil {
			if err := loadAndMerge(globalPath, cfg); err != nil {
				return nil, fmt.Errorf("failed to load global config: %w", err)
			}
		}
	}

	// 3. Load Project Config
	if projectRoot != "" {
		projectPath := filepath.Join(projectRoot, ".asdp.yaml")
		if _, err := os.Stat(projectPath); err == nil {
			if err := loadAndMerge(projectPath, cfg); err != nil {
				return nil, fmt.Errorf("failed to load project config: %w", err)
			}
		}
	}

	return cfg, nil
}

// ConfigurationLoaderImpl implements domain.ConfigurationLoader
type ConfigurationLoaderImpl struct{}

func NewConfigurationLoader() *ConfigurationLoaderImpl {
	return &ConfigurationLoaderImpl{}
}

func (l *ConfigurationLoaderImpl) LoadForProject(baseConfig *domain.Config, projectRoot string) (*domain.Config, error) {
	// Deep copy base config to avoid mutating shared state
	// YAML marshaling is a cheap way to deep copy for now
	data, _ := yaml.Marshal(baseConfig)
	newCfg := &domain.Config{}
	yaml.Unmarshal(data, newCfg)

	// Apply project override
	if projectRoot != "" {
		projectPath := filepath.Join(projectRoot, ".asdp.yaml")
		if _, err := os.Stat(projectPath); err == nil {
			if err := loadAndMerge(projectPath, newCfg); err != nil {
				return nil, fmt.Errorf("failed to merge project config: %w", err)
			}
		}
	}
	return newCfg, nil
}

func loadAndMerge(path string, cfg *domain.Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal merges into existing struct
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}
