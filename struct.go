package main

import "net/url"

type App struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	Type        string `json:"type"`
	Scalability struct {
		MinReplicas                    int `json:"minReplicas"`
		MaxReplicas                    int `json:"maxReplicas"`
		TargetCPUUtilizationPercentage int `json:"targetCPUUtilizationPercentage"`
	} `json:"scalability"`
	ImageName      string `json:"imageName"`
	DockerContext  string `json:"dockerContext,omitempty"`
	DockerFilePath string `json:"dockerFilePath,omitempty"`
	Ports          struct {
		External struct {
			HTTP      int  `json:"http"`
			WEBSOCKET bool `json:"websocket,omitempty"`
		} `json:"external"`
		Internal struct {
			HTTP int `json:"http"`
		} `json:"internal"`
	} `json:"ports"`
	Readiness struct {
		Path string `json:"path,omitempty"`
		Port int    `json:"port,omitempty"`
		Exec struct {
			Command []string `json:"command,omitempty"`
		} `json:"exec,omitempty"`
		InitialDelaySeconds int `json:"initialDelaySeconds"`
		PeriodSeconds       int `json:"periodSeconds"`
		FailureThreshold    int `json:"failureThreshold"`
		TimeoutSeconds      int `json:"timeoutSeconds"`
		SuccessThreshold    int `json:"successThreshold"`
	} `json:"readiness"`
	Health struct {
		Path string `json:"path,omitempty"`
		Port int    `json:"port,omitempty"`
		Exec struct {
			Command []string `json:"command,omitempty"`
		} `json:"exec,omitempty"`
		InitialDelaySeconds int `json:"initialDelaySeconds"`
		PeriodSeconds       int `json:"periodSeconds"`
		FailureThreshold    int `json:"failureThreshold"`
		TimeoutSeconds      int `json:"timeoutSeconds"`
		SuccessThreshold    int `json:"successThreshold"`
	} `json:"health"`
	PreStop struct {
		HTTPGet struct {
			Path string `json:"path,omitempty"`
			Port int    `json:"port,omitempty"`
		} `json:"httpGet,omitempty"`
		Exec struct {
			Command []string `json:"command,omitempty"`
		} `json:"exec,omitempty"`
	} `json:"preStop"`
	Resources struct {
		Requests struct {
			Cpu    string `json:"cpu,omitempty"`
			Memory string `json:"memory,omitempty"`
		} `json:"requests,omitempty"`
		Limits struct {
			Cpu    string `json:"cpu,omitempty"`
			Memory string `json:"memory,omitempty"`
		} `json:"limits,omitempty"`
	} `json:"resources,omitempty"`
	Secretmounts []struct {
		Mountpath  string `json:"mountpath"`
		Secretname string `json:"secretname"`
	} `json:"secretmounts"`
}

type Ingresstmpl struct {
	Application_name string
	Build_id         string
	Deploy           App
	Hostnames        []string
	Namespace        string
}

type Autoscalertmpl struct {
	Deploy      App
	Deploy_name string
	Namespace   string
}

type Servicetmpl struct {
	Application_name string
	Build_id         string
	Cluster_ip       string
	Deploy           App
	Namespace        string
}

type Clinictmpl struct {
	Clinic_name     string
	Clinic_hostname string
	Namespace       string
}

type Deploytmpl struct {
	Application_name           string
	Bamboo_AWS_HOSTNAME        string
	Bamboo_deploy_release      string
	Build_id                   string
	Build_nr                   string
	CONSUL_APPLICATION         string
	CONSUL_ENVIRONMENT         string
	CONSUL_FULL_URL            *url.URL
	CONSUL_PASSWORD            string
	CONSUL_URL                 *url.URL
	CONSUL_USERNAME            string
	Deploy                     App
	Deploy_name                string
	Git_repo                   *url.URL
	Namespace                  string
	NEW_RELIC_API_KEY_PASSWORD string
	NEW_RELIC_API_URL          *url.URL
	NEW_RELIC_LICENSE_KEY      string
	Ssh_key                    string
}
