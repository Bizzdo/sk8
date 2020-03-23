package main

import (
	"strings"
)

type configType int

const (
	typeSK8     configType = 1
	typeKubectl configType = 2
)

// SK8config s the root of the YAML config for a service
type SK8config struct {
	Name         string            `json:"name,omitempty"`
	Namespace    string            `json:"namespace,omitempty"`
	Site         string            `json:"site,omitempty"`
	Image        string            `json:"image,omitempty"`
	Version      string            `json:"version,omitempty"`
	ImageVersion *string           `json:"imageversion,omitempty"`
	Registry     Registry          `json:"registry,omitempty"`
	Port         int               `json:"port,omitempty"`
	Parents      []string          `json:"parents,omitempty"`
	Custom       interface{}       `json:"custom,omitempty"`
	Notes        map[string]string `json:"notes,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
	Extra        Extra             `json:"extra,omitempty"`
	Env          EnvMap            `json:"env,omitempty"`
	Volume       []VolumeType      `json:"volume,omitempty"`
	Templates    map[string]string `json:"templates,omitempty"`
	Features     []string          `json:"features,omitempty"`
	Override     *SK8config        `json:"override,omitempty"`
	cfgType      configType
	Kind         string       `json:"kind,omitempty"`
	RawMetadata  *rawMetadata `json:"metadata,omitempty"`
	RawYAML      []byte       `json:"rawYaml,omitempty"`
	Containers   []Container  `json:"containers,omitempty"`
}

type rawMetadata struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type Container struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	Env       EnvMap `json:"env,omitempty"`
	Liveness  *Probe `json:"liveness,omitempty"`
	Readyness *Probe `json:"readyness,omitempty"`
}

// HasFeature is a helper for the templating to test if a feature is used or not
func (sk8 *SK8config) HasFeature(feat string) bool {
	if sk8.Features == nil {
		return false
	}
	for _, f := range sk8.Features {
		if strings.EqualFold(feat, f) {
			return true
		}
	}
	return false
}

// Registry is used to specified the Docker-image source
type Registry struct {
	Host string `json:"host,omitempty"`
	Path string `json:"path,omitempty"`
}

// EnvMap contains all the Enviroment-settings for the service, for all sources
type EnvMap struct {
	Values map[string]string    `json:"values,omitempty"`
	Config map[string]EnvConfig `json:"config,omitempty"`
	Secret map[string]EnvConfig `json:"secret,omitempty"`
	Fields map[string]string    `json:"fields,omitempty"`
}

// EnvConfig contains name and key for mapping env-vars from ConfigMaps and Secrets
type EnvConfig struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Extra has properties for replica-count, history-length and probes
type Extra struct {
	Replicas  *int   `json:"replicas,omitempty"`
	History   *int   `json:"history,omitempty"`
	Liveness  *Probe `json:"liveness,omitempty"`
	Readyness *Probe `json:"readyness,omitempty"`
}

// Probe describes the liveness/readyness-probe definition
type Probe struct {
	Path                string `json:"path,omitempty"`
	Port                *int   `json:"port,omitempty"`
	InitialDelaySeconds *int   `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      *int   `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       *int   `json:"periodSeconds,omitempty"`
}

// VolumeType is how to get a "file" in the pod from various sources
type VolumeType struct {
	Name     string        `json:"name,omitempty"`
	Path     string        `json:"path,omitempty"`
	HostDir  string        `json:"hostdir,omitempty"`
	HostFile string        `json:"hostfile,omitempty"`
	ReadOnly bool          `json:"readonly,omitempty"`
	EmptyDir bool          `json:"empty,omitempty"`
	Config   *VolumeSource `json:"config,omitempty"`
	Secret   *VolumeSource `json:"secret,omitempty"`
}

// VolumeSource is what files to actually map from a source-definition
type VolumeSource struct {
	Name  string            `json:"name,omitempty"`
	Items map[string]string `json:"items,omitempty"`
}

// type VolumeConfigItem struct {
// 	Key  string `json:"key,omitempty"`
// 	Path string `json:"path,omitempty"`
// }
