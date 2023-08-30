package models

type ImageConfig struct {
	Repository string `yaml:"repository"`
	PullPolicy string `yaml:"pullPolicy"`
	Tag        string `yaml:"tag"`
}

type ServiceAccountConfig struct {
	Create      bool              `yaml:"create"`
	Annotations map[string]string `yaml:"annotations"`
	Name        string            `yaml:"name"`
}

type PathConfig struct {
	Path     string `yaml:"path"`
	PathType string `yaml:"pathType"`
}

type HostConfig struct {
	Host  string       `yaml:"host"`
	Paths []PathConfig `yaml:"paths"`
}

type IngressConfig struct {
	Enabled     bool              `yaml:"enabled"`
	ClassName   string            `yaml:"className"`
	Annotations map[string]string `yaml:"annotations"`
	Hosts       []HostConfig      `yaml:"hosts"`
	TLS         []interface{}     `yaml:"tls"`
}

type NginxConfig struct {
	ReplicaCount       int                    `yaml:"replicaCount"`
	Image              ImageConfig            `yaml:"image"`
	ImagePullSecrets   []interface{}          `yaml:"imagePullSecrets"`
	NameOverride       string                 `yaml:"nameOverride"`
	FullnameOverride   string                 `yaml:"fullnameOverride"`
	ServiceAccount     ServiceAccountConfig   `yaml:"serviceAccount"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations"`
	PodSecurityContext map[string]interface{} `yaml:"podSecurityContext"`
	SecurityContext    map[string]interface{} `yaml:"securityContext"`
	Service            struct {
		Type string `yaml:"type"`
		Port int    `yaml:"port"`
	} `yaml:"service"`
	Ingress   IngressConfig `yaml:"ingress"`
	Resources struct {
		Limits   map[string]string `yaml:"limits"`
		Requests map[string]string `yaml:"requests"`
	} `yaml:"resources"`
	Autoscaling struct {
		Enabled                        bool `yaml:"enabled"`
		MinReplicas                    int  `yaml:"minReplicas"`
		MaxReplicas                    int  `yaml:"maxReplicas"`
		TargetCPUUtilizationPercentage int  `yaml:"targetCPUUtilizationPercentage"`
	} `yaml:"autoscaling"`
	NodeSelector map[string]string      `yaml:"nodeSelector"`
	Tolerations  []interface{}          `yaml:"tolerations"`
	Affinity     map[string]interface{} `yaml:"affinity"`
}
