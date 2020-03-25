package kubernetes

type ConfigResponse struct {
	Kind        string `json:"kind"`
	APIVersion  string `json:"apiVersion"`
	Preferences struct {
	} `json:"preferences"`
	Clusters       []Cluster `json:"clusters"`
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
		Server string `json:"server"`
	} `json:"cluster"`
}
