package azure_devops

import (
	"github.com/sirupsen/logrus"
)

func CreateServiceEndpoint(username string, accessToken string, name string, organization string, project string, masterPlaneFqdn string, token string, certificate string) {
	existingEndpointId := getEndpointId(username, accessToken, name, organization, project)
	if existingEndpointId != "" {
		logrus.Println("Endpoint with name " + name + " in " + organization + ", project " + project + " already exists. Connection will be updated")
		updateEndpoint(username, accessToken, existingEndpointId, name, organization, project, masterPlaneFqdn, token, certificate)
	} else {
		createEndpoint(username, accessToken, name, organization, project, masterPlaneFqdn, token, certificate)
	}
}

func getEndpointId(username string, accessToken string, name string, organization string, project string) string {
	response := searchForEndpoint(username, accessToken, name, organization, project)
	endpointId := ""
	if len(response.Value) > 0 {
		endpointId = response.Value[0].ID
	}

	return endpointId
}
