package config

import (
	"github.com/NubeIO/git/pkg/git"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sync"
)

type Config struct {
	Path         string `yaml:"path" json:"path"`
	GitToken     string `yaml:"git_token" json:"gitToken"`
	GitOwner     string `yaml:"git_owner" json:"gitOwner"`
	ReleasesRepo string `yaml:"releases_repo" json:"releasesRepo"`
	mutex        sync.Mutex
	configPath   string
}

func New(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config *Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (inst *Config) SaveConfig(config *Config) error {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(inst.configPath, data, 0644)
	return err
}

func (inst *Config) GetConfig() *Config {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	return inst
}

func (inst *Config) GetGitToken() string {
	config := inst.GetConfig()
	decoded := git.DecodeToken(config.GitToken)
	return decoded
}
