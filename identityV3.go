package goOpenstackAuth

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
)

type AuthV3Options struct {
	IdentityEndpoint  string
	Region            string
	Username          string
	UserId            string
	Password          string
	ProjectName       string
	ProjectId         string
	UserDomainName    string
	UserDomainId      string
	ProjectDomainName string
	ProjectDomainId   string
}

var (
	AuthenticationV3 = NewAuthV3
)

type Authentication interface {
	GetToken() (*tokens.Token, error)
	GetServiceEndpoint(serviceType, region, serviceInterface string) (string, error)
	GetProject() (*Project, error)
}

type AuthV3 struct {
	Options     AuthV3Options
	tokenResult *tokens.CreateResult
	client      *gophercloud.ServiceClient
}

type Project struct {
	Name     string `mapstructure:"name"`
	ID       string `mapstructure:"id"`
	DomainID string `mapstructure:"domain_id"`
}

func NewAuthV3(authOpts AuthV3Options) Authentication {
	return &AuthV3{Options: authOpts}
}

func (a *AuthV3) GetToken() (*tokens.Token, error) {
	var err error
	err = a.createTokenCommonResult()
	if err != nil {
		return nil, err
	}

	token, err := a.tokenResult.Extract()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *AuthV3) GetServiceEndpoint(serviceType, region, serviceInterface string) (string, error) {
	if a.tokenResult == nil {
		a.GetToken()
	}

	// get catalog
	catalog, err := a.tokenResult.ExtractServiceCatalog()
	if err != nil {
		return "", err
	}

	// get entry from catalog
	serviceEntry, err := getServiceEntry(serviceType, catalog)
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

func (a *AuthV3) GetProject() (*Project, error) {
	if a.tokenResult == nil {
		_, err := a.GetToken()
		if err != nil {
			return nil, err
		}
	}

	return extractProject(a.tokenResult.Body)
}

func (a *AuthV3) getAuthOptions() gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: a.Options.IdentityEndpoint,
		Username:         a.Options.Username,
		UserID:           a.Options.UserId,
		Password:         a.Options.Password,
		DomainName:       a.Options.UserDomainName,
		DomainID:         a.Options.UserDomainId,
	}
}

func (a *AuthV3) getClient() (*gophercloud.ServiceClient, error) {
	// get provider client struct
	provider, err := openstack.AuthenticatedClient(a.getAuthOptions())
	if err != nil {
		return nil, err
	}
	return openstack.NewIdentityV3(provider), nil
}

func (a *AuthV3) createTokenCommonResult() error {
	scope := tokens.Scope{
		ProjectName: a.Options.ProjectName,
		ProjectID:   a.Options.ProjectId,
		DomainName:  a.Options.ProjectDomainName,
		DomainID:    a.Options.ProjectDomainId,
	}

	// init the v3 client
	var err error
	if a.client == nil {
		a.client, err = a.getClient()
		if err != nil {
			return err
		}
	}

	// create common result
	result := tokens.Create(a.client, a.getAuthOptions(), &scope)

	// save common result
	a.tokenResult = &result

	return nil
}

// private

func extractProject(body interface{}) (*Project, error) {
	var response struct {
		Token struct {
			Project `mapstructure:"project"`
		} `mapstructure:"token"`
	}

	err := mapstructure.Decode(body, &response)
	if err != nil {
		return nil, err
	}

	return &Project{
		ID:       response.Token.ID,
		Name:     response.Token.Name,
		DomainID: response.Token.DomainID,
	}, nil
}

func getServiceEndpoint(region string, serviceInterface string, entry *tokens.CatalogEntry) (string, error) {
	if entry != nil && len(entry.Endpoints) > 0 {
		var endpoint string
		for _, ep := range entry.Endpoints {
			if region != "" {
				if ep.Interface == serviceInterface && ep.Region == region {
					endpoint = ep.URL
					break
				}
			} else {
				if ep.Interface == serviceInterface {
					endpoint = ep.URL
					break
				}
			}
		}
		return endpoint, nil
	} else {
		return "", fmt.Errorf("Authenticate: getServicePublicEndpoint: entry nil or no endpoints found for %+v.", entry)
	}
	return "", nil
}

func getServiceEntry(serviceType string, catalog *tokens.ServiceCatalog) (*tokens.CatalogEntry, error) {
	if catalog != nil && len(catalog.Entries) > 0 {
		serviceEntry := tokens.CatalogEntry{}
		for _, service := range catalog.Entries {
			if service.Type == serviceType {
				serviceEntry = service
				break
			}
		}
		return &serviceEntry, nil
	} else {
		return nil, fmt.Errorf("Authenticate: GetServicePublicEndpoint: catalog nil or emtpy.")
	}

	return nil, nil
}
