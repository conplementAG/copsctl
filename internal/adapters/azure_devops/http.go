package azure_devops

import (
	"encoding/base64"
	"fmt"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"log"

	"github.com/imroc/req"
	"github.com/spf13/viper"
)

func prepareHeaders(username string, accessToken string) req.Header {
	if viper.GetBool(flags.Verbose) {
		log.Println("Using username: " + username)
		log.Println("Using token: " + accessToken)
	}

	headers := make(map[string]string)

	headers["Accept"] = "application/json"
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = createBasicAuthHeader(username, accessToken)

	return headers
}

func createBasicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func buildBaseUrl(organization string, project string) string {
	return "https://dev.azure.com/" + organization + "/" + project + "/_apis/serviceendpoint/endpoints"
}

func searchForEndpoint(username string, accesstoken string, name string, organization string, project string) endpointsByNameQueryResponse {
	headers := prepareHeaders(username, accesstoken)

	param := req.QueryParam{
		"api-version":   "5.0-preview.2",
		"includeFailed": "true",
		"endpointNames": name,
	}

	reqUrl := buildBaseUrl(organization, project)

	r, err := req.Get(reqUrl, headers, param)

	if viper.GetBool(flags.Verbose) {
		log.Printf("%+v", r)
	}

	if err != nil {
		panic(err)
	}

	statusCode := r.Response().StatusCode

	if statusCode != 200 {
		panic(fmt.Sprintf("[searchForEndpoint] HTTP request failed with status code %d", statusCode))
	}

	var response endpointsByNameQueryResponse

	err = r.ToJSON(&response)

	if err != nil {
		panic(err)
	}

	return response
}

func deleteEndpoint(username string, accessToken string, endpointId string, organization string, project string) {
	header := prepareHeaders(username, accessToken)

	param := req.QueryParam{
		"api-version": "5.0-preview.2",
	}

	reqUrl := buildBaseUrl(organization, project) + "/" + endpointId

	r, err := req.Delete(reqUrl, header, param)

	if viper.GetBool(flags.Verbose) {
		log.Printf("%+v", r)
	}

	if err != nil {
		panic(err)
	}

	statusCode := r.Response().StatusCode

	if statusCode != 204 {
		panic(fmt.Sprintf("[Endpoint Delete] HTTP request failed with status code %d", statusCode))
	}
}

func createEndpoint(username string, accessToken string, name string, organization string, project string, kubernetesMasterUrl string, token string, certificate string) {
	headers := prepareHeaders(username, accessToken)

	param := req.QueryParam{
		"api-version": "5.0-preview.2",
	}

	createEndpointRequest := &createEndpointRequest{
		Name: name,
		Type: "kubernetes",
		URL:  kubernetesMasterUrl,
	}

	createEndpointRequest.Authorization.Parameters.APIToken = token
	createEndpointRequest.Authorization.Parameters.ServiceAccountCertificate = certificate
	createEndpointRequest.Authorization.Scheme = "Token"
	createEndpointRequest.Data.AuthorizationType = "ServiceAccount"
	createEndpointRequest.Data.AcceptUntrustedCerts = "false"

	body := req.BodyJSON(createEndpointRequest)

	reqUrl := buildBaseUrl(organization, project)

	r, err := req.Post(reqUrl, headers, param, body)

	if viper.GetBool(flags.Verbose) {
		log.Printf("%+v", r)
	}

	if err != nil {
		panic(err)
	}

	statusCode := r.Response().StatusCode

	if statusCode != 200 {
		panic(fmt.Sprintf("[createEndpoint] HTTP request failed with status code %d", statusCode))
	}
}

func updateEndpoint(username string, accessToken string, endpointId string, name string, organization string, project string, kubernetesMasterUrl string, token string, certificate string) {
	headers := prepareHeaders(username, accessToken)

	param := req.QueryParam{
		"api-version": "5.0-preview.2",
	}

	updateEndpointRequest := &updateEndpointRequest{
		Name: name,
		Type: "kubernetes",
		URL:  kubernetesMasterUrl,
		ID:   endpointId,
	}

	updateEndpointRequest.Authorization.Parameters.APIToken = token
	updateEndpointRequest.Authorization.Parameters.ServiceAccountCertificate = certificate
	updateEndpointRequest.Authorization.Scheme = "Token"
	updateEndpointRequest.Data.AuthorizationType = "ServiceAccount"
	updateEndpointRequest.Data.AcceptUntrustedCerts = "false"

	body := req.BodyJSON(updateEndpointRequest)

	reqUrl := buildBaseUrl(organization, project) + "/" + endpointId

	r, err := req.Put(reqUrl, headers, param, body)

	if viper.GetBool(flags.Verbose) {
		log.Printf("%+v", r)
	}

	if err != nil {
		panic(err)
	}

	statusCode := r.Response().StatusCode

	if statusCode != 200 {
		panic(fmt.Sprintf("[updateEndpoint] HTTP request failed with status code %d", statusCode))
	}
}

type endpointsByNameQueryResponse struct {
	Count int `json:"count"`
	Value []struct {
		Data struct {
			AuthorizationType    string `json:"authorizationType"`
			AcceptUntrustedCerts string `json:"acceptUntrustedCerts"`
		} `json:"data"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		URL       string `json:"url"`
		CreatedBy struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			UniqueName  string `json:"uniqueName"`
		} `json:"createdBy"`
		Authorization struct {
			Scheme string `json:"scheme"`
		} `json:"authorization"`
		IsShared bool `json:"isShared"`
		IsReady  bool `json:"isReady"`
	} `json:"value"`
}

type createEndpointRequest struct {
	Description         string      `json:"description"`
	AdministratorsGroup interface{} `json:"administratorsGroup"`
	Authorization       struct {
		Parameters struct {
			APIToken                  string `json:"apiToken"`
			ServiceAccountCertificate string `json:"serviceAccountCertificate"`
		} `json:"parameters"`
		Scheme string `json:"scheme"`
	} `json:"authorization"`
	CreatedBy interface{} `json:"createdBy"`
	Data      struct {
		AuthorizationType    string `json:"authorizationType"`
		AcceptUntrustedCerts string `json:"acceptUntrustedCerts"`
	} `json:"data"`
	Name            string      `json:"name"`
	Type            string      `json:"type"`
	URL             string      `json:"url"`
	ReadersGroup    interface{} `json:"readersGroup"`
	GroupScopeID    interface{} `json:"groupScopeId"`
	IsReady         bool        `json:"isReady"`
	OperationStatus interface{} `json:"operationStatus"`
}

type updateEndpointRequest struct {
	Description         string      `json:"description"`
	AdministratorsGroup interface{} `json:"administratorsGroup"`
	Authorization       struct {
		Parameters struct {
			APIToken                  string `json:"apiToken"`
			ServiceAccountCertificate string `json:"serviceAccountCertificate"`
		} `json:"parameters"`
		Scheme string `json:"scheme"`
	} `json:"authorization"`
	CreatedBy interface{} `json:"createdBy"`
	Data      struct {
		AuthorizationType    string `json:"authorizationType"`
		AcceptUntrustedCerts string `json:"acceptUntrustedCerts"`
	} `json:"data"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	IsReady  bool   `json:"isReady"`
	IsShared bool   `json:"isShared"`
	ID       string `json:"id"`
	Owner    string `json:"owner"`
}
