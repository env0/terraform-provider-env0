package env0

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
)

func getCredentialsByName(name string, prefix string, meta interface{}) (client.Credentials, error) {
	apiClient := meta.(client.ApiClientInterface)

	credentialsList, err := apiClient.CloudCredentialsList()
	if err != nil {
		return client.Credentials{}, err
	}

	var foundCredentials []client.Credentials
	for _, credentials := range credentialsList {
		if credentials.Name == name && strings.HasPrefix(credentials.Type, prefix) {
			foundCredentials = append(foundCredentials, credentials)
		}
	}

	if len(foundCredentials) == 0 {
		return client.Credentials{}, fmt.Errorf("credentials with name %v not found", name)
	}

	if len(foundCredentials) > 1 {
		return client.Credentials{}, fmt.Errorf("found multiple credentials with name: %s. Use id instead or make sure credential names are unique %v", name, foundCredentials)
	}

	return foundCredentials[0], nil
}

func getCredentialsById(id string, prefix string, meta interface{}) (client.Credentials, error) {
	apiClient := meta.(client.ApiClientInterface)
	credentials, err := apiClient.CloudCredentials(id)
	if err != nil {
		if _, ok := err.(*client.NotFoundError); ok {
			return client.Credentials{}, errors.New("credentials not found")
		}
		return client.Credentials{}, err
	}

	if !strings.HasPrefix(credentials.Type, prefix) {
		return client.Credentials{}, fmt.Errorf("credentials type mistmatch %s", credentials.Type)
	}

	return credentials, nil
}

func getCredentials(id string, prefix string, meta interface{}) (client.Credentials, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		log.Println("[INFO] Resolving credentials by id: ", id)
		return getCredentialsById(id, prefix, meta)
	} else {
		log.Println("[INFO] Resolving credentials by name: ", id)
		return getCredentialsByName(id, prefix, meta)
	}
}