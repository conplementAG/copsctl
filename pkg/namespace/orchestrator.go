package namespace

import (
	"strings"

	"github.com/conplementAG/copsctl/pkg/adapters/kubernetes"
	"github.com/conplementAG/copsctl/pkg/common/fileprocessing"
	"github.com/conplementAG/copsctl/pkg/common/logging"
	"github.com/spf13/viper"
)

// Create creates a CopsNamespace Custom-Resource-Definition with the given name and user
func Create() {
	namespaceName := viper.GetString("name")
	adminUsername := viper.GetString("user")

	copsnamespace := renderTemplate(namespaceName, adminUsername)

	temporaryDirectory, temporaryFile := fileprocessing.WriteStringToTemporaryFile(copsnamespace, "copsnamespace.yaml")
	kubernetes.Apply(temporaryFile)
	fileprocessing.DeletePath(temporaryDirectory)

	logging.LogSuccessf("Cops namespace %s successfully created\n", namespaceName)
}

func renderTemplate(namespaceName string, adminUsername string) string {
	copsnamespace := strings.Replace(copsNamespaceTemplate, "{{ namespaceName }}", namespaceName, -1)
	copsnamespace = strings.Replace(copsnamespace, "{{ adminUsername }}", adminUsername, -1)
	return copsnamespace
}
