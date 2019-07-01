package kubernetes

type ConfigResponse struct {
	Kind        string `json:"kind"`
	APIVersion  string `json:"apiVersion"`
	Preferences struct {
	} `json:"preferences"`
	Clusters []Cluster `json:"clusters"`
	Users    []struct {
		Name string `json:"name"`
		User struct {
			ClientCertificateData string `json:"client-certificate-data"`
			ClientKeyData         string `json:"client-key-data"`
			Token                 string `json:"token"`
		} `json:"user"`
	} `json:"users"`
	Contexts       []Context `json:"contexts"`
	CurrentContext string    `json:"current-context"`
}

type Context struct {
	Name    string `json:"name"`
	Context struct {
		Cluster string `json:"cluster"`
		User    string `json:"user"`
	} `json:"context"`
}

type Cluster struct {
	Name    string `json:"name"`
	Cluster struct {
		Server                string `json:"server"`
		InsecureSkipTLSVerify bool   `json:"insecure-skip-tls-verify"`
	} `json:"cluster"`
}
