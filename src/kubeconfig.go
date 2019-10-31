package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type kubeconfig struct {
	Kind           string        `json:"kind"`
	Clusters       []kubeCluster `json:"clusters"`
	Contexts       []kubeContext `json:"contexts"`
	CurrentContext string        `json:"current-context"`
}

type kubeCluster struct {
	Name    string          `json:"name"`
	Context kubeClusterData `json:"cluster"`
}

type kubeClusterData struct {
	Server string `yaml:"server"`
}

type kubeContext struct {
	Name    string          `json:"name"`
	Context kubeContextData `json:"context"`
}

type kubeContextData struct {
	Cluster string `json:"cluster"`
	User    string `json:"user"`
}

func validateKubeconfig(validate bool) bool {
	fn := os.Getenv("KUBECONFIG")
	if fn == "" {
		log.Warning("Env 'KUBECONFIG' is missing")
		return false
	}

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		if validate {
			panic(err)
		}
		return false
	}

	var kc kubeconfig
	var o interface{}
	//yaml_mapstr.Unmarshal(buf, &o)
	err = yamlUnmarshal(buf, &o)
	if err != nil {
		log.Error("Unable to read kubeconfig: %s", err.Error())
		return false
	}
	buf, err = json.Marshal(o)
	json.Unmarshal(buf, &kc)

	var ctx *kubeContext
	for _, c := range kc.Contexts {
		if c.Name == kc.CurrentContext {
			ctx = &c
		}
	}

	if ctx == nil {
		log.Error("Current-context doesn't exist in the KUBECONFIG")
		return false
	}

	var cluster *kubeCluster
	for _, cl := range kc.Clusters {
		if cl.Name == ctx.Context.Cluster {
			cluster = &cl
		}
	}

	if cluster == nil {
		log.Error("Can't find the clusterdata for the current context")
		return false
	}

	log.Noticef("Current Kubeconfig: %s", cluster.Context.Server)
	return true
}
