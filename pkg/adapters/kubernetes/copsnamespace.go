package kubernetes

type CopsNamespaceResponse struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		NamespaceAdminUsers []string `yaml:"namespaceAdminUsers"`
	} `yaml:"spec"`
}
