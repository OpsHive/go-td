package models

type Service struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name      string            `yaml:"name"`
	Labels    map[string]string `yaml:"labels"`
	Namespace string            `yaml:"namespace"`
}

type Spec struct {
	Template Template `yaml:"template"`
}

type Template struct {
	Spec SpecTemplate `yaml:"spec"`
}

type SpecTemplate struct {
	ImagePullSecrets []ImagePullSecret `yaml:"imagePullSecrets"`
	Containers       []Container       `yaml:"containers"`
}

type ImagePullSecret struct {
	Name string `yaml:"name"`
}

type Container struct {
	Image string   `yaml:"image"`
	Ports []Port   `yaml:"ports"`
	Env   []EnvVar `yaml:"env"`
}

type Port struct {
	ContainerPort int `yaml:"containerPort"`
}

type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
