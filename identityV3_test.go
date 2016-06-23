package goOpenstackAuth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
)

func resetAuthentication() {
	AuthenticationV3 = newMockAuthenticationV3
}

func TestAuthenticationTokenSuccess(t *testing.T) {
	resetAuthentication()

	options := AuthV3Options{
		IdentityEndpoint: "http://some_test_url",
		UserId:           "miau",
		Password:         "123456789",
		ProjectId:        "bup",
	}

	a := AuthenticationV3(options)
	token, err := a.GetToken()
	if err != nil {
		t.Error(fmt.Sprint(`Expected to not get an error. `, err.Error()))
		return
	}

	if !strings.Contains(token.ID, "test_token_id") {
		diffString := StringDiff(token.ID, "test_token_id")
		t.Error(fmt.Sprintf("Token does not match. \n \n %s", diffString))
	}
}

func TestAuthenticationEndpointSuccess(t *testing.T) {
	resetAuthentication()

	options := AuthV3Options{
		IdentityEndpoint: "http://some_test_url",
		UserId:           "miau",
		Password:         "123456789",
		ProjectId:        "bup",
	}

	a := AuthenticationV3(options)
	endpoint, err := a.GetServiceEndpoint("arc", "staging", "public")
	if err != nil {
		t.Error(fmt.Sprint(`Expected to not get an error. `, err.Error()))
		return
	}

	if !strings.Contains(endpoint, "https://arc-app-staging/public") {
		diffString := StringDiff(endpoint, "https://arc-app-staging/public")
		t.Error(fmt.Sprintf("Endpoint does not match. \n \n %s", diffString))
	}
}

func TestAuthenticationEndpointNotGivenARegion(t *testing.T) {
	resetAuthentication()

	options := AuthV3Options{
		IdentityEndpoint: "http://some_test_url",
		UserId:           "miau",
		Password:         "123456789",
		ProjectId:        "bup",
	}

	a := AuthenticationV3(options)
	endpoint, err := a.GetServiceEndpoint("arc", "", "public")
	if err != nil {
		t.Error(fmt.Sprint(`Expected to not get an error. `, err.Error()))
		return
	}

	if !strings.Contains(endpoint, "https://arc-app-staging/public") {
		diffString := StringDiff(endpoint, "https://arc-app-staging/public")
		t.Error(fmt.Sprintf("Endpoint does not match. \n \n %s", diffString))
	}
}

func TestAuthenticationGivenARegion(t *testing.T) {
	resetAuthentication()

	options := AuthV3Options{
		IdentityEndpoint: "http://some_test_url",
		UserId:           "miau",
		Password:         "123456789",
		ProjectId:        "bup",
	}

	a := AuthenticationV3(options)
	endpoint, err := a.GetServiceEndpoint("arc", "production", "internal")
	if err != nil {
		t.Error(fmt.Sprint(`Expected to not get an error. `, err.Error()))
		return
	}

	if !strings.Contains(endpoint, "https://arc-app-prod/internal") {
		diffString := StringDiff(endpoint, "https://arc-app-prod/internal")
		t.Error(fmt.Sprintf("Endpoint does not match. \n \n %s", diffString))
	}
}

func TestAuthenticationGivenAWrongRegion(t *testing.T) {
	resetAuthentication()

	options := AuthV3Options{
		IdentityEndpoint: "http://some_test_url",
		UserId:           "miau",
		Password:         "123456789",
		ProjectId:        "bup",
	}

	a := AuthenticationV3(options)
	endpoint, err := a.GetServiceEndpoint("arc", "non_exisitng_region", "internal")
	if err != nil {
		t.Error(fmt.Sprint(`Expected to not get an error. `, err.Error()))
		return
	}

	if endpoint != "" {
		t.Error("Endpoint should be empty")
	}
}

//
// Mock authentication interface
//

type MockV3 struct {
	Options     AuthV3Options
	tokenResult *tokens.CreateResult
}

func newMockAuthenticationV3(authOpts AuthV3Options) Authentication {
	return &MockV3{Options: authOpts}
}

func (a *MockV3) GetToken() (*tokens.Token, error) {
	token := tokens.Token{ID: "test_token_id"}
	return &token, nil
}

func (a *MockV3) GetServiceEndpoint(serviceType, region, serviceInterface string) (string, error) {
	// get entry from catalog
	serviceEntry, err := getServiceEntry(serviceType, &catalog1)
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

var catalog1 = tokens.ServiceCatalog{
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
			{ID: "e-54b8d28fc", Region: "staging", Interface: "public", URL: "https://lyra-app"},
		}},
	},
}
