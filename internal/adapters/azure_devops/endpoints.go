package azure_devops

import "log"

func CreateServiceEndpoint(username string, accesstoken string, name string, organization string, project string, masterPlaneFqdn string, token string, certificate string) {
	existingEndpointId := getEndpointId(username, accesstoken, name, organization, project)
	if existingEndpointId != "" {
		log.Println("Endpoint with name " + name + " in " + organization + ", project " + project + " already exists. Connection will be updated")
		updateEndpoint(username, accesstoken, existingEndpointId, name, organization, project, masterPlaneFqdn, token, certificate)
	} else {
		createEndpoint(username, accesstoken, name, organization, project, masterPlaneFqdn, token, certificate)
	}
}

func RemoveServiceEndpoint(username string, accesstoken string, name string, organization string, project string) {
	endpoints := searchForEndpoint(username, accesstoken, name, organization, project)
	if len(endpoints.Value) > 1 {
		panic("Endpoint delete cannot be performed since multiple endpoints with the name have been found.")
	}

	if len(endpoints.Value) == 0 {
		// nothing to delete, all good
		return
	}

	endpointId := endpoints.Value[0].ID

	deleteEndpoint(username, accesstoken, endpointId, organization, project)

	log.Println("Azure DevOps endpoint " + name + " deleted, organization " + organization + ", project " + project)
}

func getEndpointId(username string, accesstoken string, name string, organization string, project string) string {
	response := searchForEndpoint(username, accesstoken, name, organization, project)
	endpointId := ""
	if len(response.Value) > 0 {
		endpointId = response.Value[0].ID
	}

	return endpointId
}
