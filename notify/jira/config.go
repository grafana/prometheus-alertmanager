// Copyright 2015 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jira

import (
	"errors"

	commoncfg "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"

	amcommoncfg "github.com/prometheus/alertmanager/config/common"
)

// DefaultJiraConfig defines default values for Jira configurations.
var DefaultJiraConfig = JiraConfig{
	NotifierConfig: amcommoncfg.NotifierConfig{
		VSendResolved: true,
	},
	Summary:     `{{ template "jira.default.summary" . }}`,
	Description: `{{ template "jira.default.description" . }}`,
	Priority:    `{{ template "jira.default.priority" . }}`,
}

type JiraConfig struct {
	amcommoncfg.NotifierConfig `yaml:",inline" json:",inline"`
	HTTPConfig                 *commoncfg.HTTPClientConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	APIURL *amcommoncfg.URL `yaml:"api_url,omitempty" json:"api_url,omitempty"`

	Project     string   `yaml:"project,omitempty" json:"project,omitempty"`
	Summary     string   `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Labels      []string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Priority    string   `yaml:"priority,omitempty" json:"priority,omitempty"`
	IssueType   string   `yaml:"issue_type,omitempty" json:"issue_type,omitempty"`

	ReopenTransition  string         `yaml:"reopen_transition,omitempty" json:"reopen_transition,omitempty"`
	ResolveTransition string         `yaml:"resolve_transition,omitempty" json:"resolve_transition,omitempty"`
	WontFixResolution string         `yaml:"wont_fix_resolution,omitempty" json:"wont_fix_resolution,omitempty"`
	ReopenDuration    model.Duration `yaml:"reopen_duration,omitempty" json:"reopen_duration,omitempty"`

	Fields map[string]any `yaml:"fields,omitempty" json:"custom_fields,omitempty"`
}

func (c *JiraConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultJiraConfig
	type plain JiraConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	if c.Project == "" {
		return errors.New("missing project in jira_config")
	}
	if c.IssueType == "" {
		return errors.New("missing issue_type in jira_config")
	}

	return nil
}
