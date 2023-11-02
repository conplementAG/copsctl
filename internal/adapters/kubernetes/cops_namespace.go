package kubernetes

type CopsNamespaceResponse struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		NamespaceAdminUsers           []string             `yaml:"namespaceAdminUsers"`
		NamespaceAdminServiceAccounts []CopsServiceAccount `yaml:"namespaceAdminServiceAccounts"`
		Project                       struct {
			Name       string `yaml:"name"`
			CostCenter string `yaml:"costCenter"`
		} `yaml:"project"`
	} `yaml:"spec"`
}

type CopsServiceAccount struct {
	ServiceAccount string `yaml:"serviceAccount"`
	Namespace      string `yaml:"namespace"`
}
