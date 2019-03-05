package goOpenstackAuth

import (
	"net/http/httptest"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

//
// Mock authentication interface
//

type MockV3 struct {
	Options     AuthOptions
	tokenResult *tokens.CreateResult
	TestServer  *httptest.Server
}

func NewMockAuthenticationV3(authOpts AuthOptions) Authentication {
	return &MockV3{Options: authOpts}
}

func (a *MockV3) GetOptions() *AuthOptions {
	return &a.Options
}

func (a *MockV3) GetToken() (*tokens.Token, error) {
	token := tokens.Token{ID: "test_token_id"}
	return &token, nil
}

func (a *MockV3) GetServiceEndpoint(serviceType, region, serviceInterface string) (string, error) {
	if a.TestServer != nil {
		return a.TestServer.URL, nil
	} else {
		// get entry from catalog
		serviceEntry, err := getServiceEntry(serviceType, &Catalog1)
		if err != nil {
			return "", err
		}

		// get endpoint
		endpoint, err := getServiceEndpoint(region, serviceInterface, serviceEntry)
		if err != nil {
			return "", err
		}

		return endpoint, nil
	}
	return "", nil
}

func (a *MockV3) GetProject() (*Project, error) {
	return extractProject(CommonResult1)
}

var Catalog1 = tokens.ServiceCatalog{
	Entries: []tokens.CatalogEntry{
		{ID: "s-8be070817", Name: "Arc", Type: "arc", Endpoints: []tokens.Endpoint{
			{ID: "e-904f431c9", Region: "staging", Interface: "internal", URL: "https://arc-app-staging/internal"},
			{ID: "e-904f431c9", Region: "staging", Interface: "admin", URL: "https://arc-app-staging/admin"},
			{ID: "e-884f431c9", Region: "staging", Interface: "public", URL: "https://arc-app-staging/public"},
			{ID: "e-904f431c9", Region: "production", Interface: "internal", URL: "https://arc-app-prod/internal"},
			{ID: "e-904f431c9", Region: "production", Interface: "admin", URL: "https://arc-app-prod/admin"},
			{ID: "e-884f431c9", Region: "production", Interface: "public", URL: "https://arc-app-prod/public"},
		}},
		{ID: "s-d5e793744", Name: "Lyra", Type: "automation", Endpoints: []tokens.Endpoint{
			{ID: "e-54b8d28fc", Region: "staging", Interface: "public", URL: "https://lyra-app-staging"},
		}},
	},
}

var CommonResult1 = map[string]interface{}{"token": map[string]interface{}{"project": map[string]string{"id": "p-9597d2775", "domain_id": "o-monsoon2", "name": ""}}}
