package azuredevops

import "log"

func CreateServiceEndpoint(name string, organization string, project string, masterPlaneFqdn string, token string, certificate string) {
	existingEndpointId := getEndpointId(name, organization, project)
	if existingEndpointId != "" {
		log.Println("Endpoint with name " + name + " in " + organization + ", project " + project + " already exists. Connection will be updated")
		updateEndpoint(existingEndpointId, name, organization, project, masterPlaneFqdn, token, certificate)
	} else {
		createEndpoint(name, organization, project, masterPlaneFqdn, token, certificate)
	}
}

func RemoveServiceEndpoint(name string, organization string, project string) {
	endpoints := searchForEndpoint(name, organization, project)
	if len(endpoints.Value) > 1 {
		panic("Endpoint delete cannot be performed since multiple endpoints with the name have been found.")
	}

	if len(endpoints.Value) == 0 {
		// nothing to delete, all good
		return
	}

	endpointId := endpoints.Value[0].ID

	deleteEndpoint(endpointId, organization, project)

	log.Println("Azure DevOps endpoint " + name + " deleted, organization " + organization + ", project " + project)
}

func getEndpointId(name string, organization string, project string) string {
	response := searchForEndpoint(name, organization, project)
	endpointId := ""
	if len(response.Value) > 0{
		endpointId = response.Value[0].ID
	}

	return endpointId
}
