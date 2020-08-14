package config

import (
	"io/ioutil"
	"regexp"

	"github.com/depscloud/api/v1alpha/schema"

	"github.com/ghodss/yaml"
)

// Config represents the structure of a config file
type Config struct {
	CompanyPatterns      []string `json:"company_patterns"`
	CompanyPatternsRegex []*regexp.Regexp
}

// IsCompanyModule determines if the provided module matches any provided regular expressions.
func (cfg *Config) IsCompanyModule(module *schema.Module) bool {
	for _, excludePatterns := range cfg.CompanyPatternsRegex {
		if excludePatterns.MatchString(module.GetOrganization()) {
			return true
		}

		if excludePatterns.MatchString(module.GetModule()) {
			return true
		}
	}

	return false
}

// Load reads and parses the provided config file into the appropriate structure
func Load(configFile string) (*Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	cfg.CompanyPatternsRegex = make([]*regexp.Regexp, len(cfg.CompanyPatterns))
	for i, pattern := range cfg.CompanyPatterns {
		cfg.CompanyPatternsRegex[i], err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
