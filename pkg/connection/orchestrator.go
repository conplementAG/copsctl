package connection

import (
	"log"
	"strings"

	"github.com/conplementAG/copsctl/pkg/adapters/kubernetes"
	"github.com/conplementAG/copsctl/pkg/common/logging"
	"github.com/spf13/viper"
)

// Connect tries to set the current context in kubectl to the determined cluster context
func Connect() {
	environmentTag := viper.GetString("environment-tag")

	config := kubernetes.GetCurrentConfig()
	selectedContext := ""
	if len(config.Contexts) > 0 {
		for _, context := range config.Contexts {
			if strings.HasPrefix(context.Context.User, "clusterUser_"+environmentTag+"-") {
				selectedContext = context.Name
				break
			}
		}
	}

	if selectedContext != "" {
		log.Printf("Connecting to cluster context: %s\n", selectedContext)
		kubernetes.UseContext(selectedContext)
		logging.LogSuccess("kubectl successfully setup")
	} else {
		panic("Could not find a proper context in your .kube/config")
	}
}
