package main

type sk8config struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Site      string            `json:"site,omitempty"`
	Image     string            `json:"image,omitempty"`
	Version   string            `json:"version,omitempty"`
	Registry  registry          `json:"registry,omitempty"`
	Port      int               `json:"port,omitempty"`
	Parents   []string          `json:"parents,omitempty"`
	Custom    interface{}       `json:"custom,omitempty"`
	Notes     map[string]string `json:"notes,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Extra     extra             `json:"extra,omitempty"`
	Env       envmap            `json:"env,omitempty"`
	Volume    []volumeType      `json:"volume,omitempty"`
	Templates map[string]string `json:"templates,omitempty"`
	Features  []string          `json:"features"`
	//	URL       string            `json:"url,omitempty"`
}

func (sk8 *sk8config) HasFeature(feat string) bool {
	if sk8.Features == nil {
		return false
	}
	for _, f := range sk8.Features {
		if feat == f {
			return true
		}
	}
	return false
}

type registry struct {
	Host string `json:"host,omitempty"`
	Path string `json:"path,omitempty"`
}

type envmap struct {
	Values map[string]string    `json:"values,omitempty"`
	Config map[string]envconfig `json:"config,omitempty"`
	Secret map[string]envconfig `json:"secret,omitempty"`
	Fields map[string]string    `json:"fields,omitempty"`
}

type envconfig struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

type extra struct {
	Replicas  *int   `json:"replicas,omitempty"`
	History   *int   `json:"history,omitempty"`
	Liveness  *probe `json:"liveness,omitempty"`
	Readyness *probe `json:"readyness,omitempty"`
}

type probe struct {
	Path                string `json:"path,omitempty"`
	Port                *int   `json:"port,omitempty"`
	InitialDelaySeconds *int   `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      *int   `json:"timeoutSeconds,omitempty"`
}

type volumeType struct {
	Name     string        `json:"name,omitempty"`
	Path     string        `json:"path,omitempty"`
	ReadOnly bool          `json:"readonly,omitempty"`
	EmptyDir bool          `json:"empty,omitempty"`
	Config   *volumeSource `json:"config,omitempty"`
	Secret   *volumeSource `json:"secret,omitempty"`
}

type volumeSource struct {
	Name  string            `json:"name,omitempty"`
	Items map[string]string `json:"items,omitempty"`
}

type volumeConfigItem struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path,omitempty"`
}
