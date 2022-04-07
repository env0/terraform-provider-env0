package client

import "errors"

type EnvironmentSchedulingExpression struct {
	Cron    string `json:"cron,omitempty"`
	Enabled bool   `json:"enabled"`
}

type EnvironmentScheduling struct {
	Deploy  *EnvironmentSchedulingExpression `json:"deploy,omitempty"`
	Destroy *EnvironmentSchedulingExpression `json:"destroy,omitempty"`
}

func (client *ApiClient) EnvironmentScheduling(environmentId string) (EnvironmentScheduling, error) {
	var result EnvironmentScheduling

	err := client.http.Get("/scheduling/environments/"+environmentId, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (client *ApiClient) EnvironmentSchedulingUpdate(environmentId string, payload EnvironmentScheduling) (EnvironmentScheduling, error) {
	var result EnvironmentScheduling

	if payload.Deploy != nil && payload.Destroy != nil {
		if payload.Deploy.Cron == payload.Destroy.Cron {
			return EnvironmentScheduling{}, errors.New("deploy and destroy cron expressions must not be the same")
		}
	}

	err := client.http.Put("/scheduling/environments/"+environmentId, payload, &result)
	if err != nil {
		return EnvironmentScheduling{}, err
	}

	return result, nil
}

func (client *ApiClient) EnvironmentSchedulingDelete(environmentId string) error {
	err := client.http.Put("/scheduling/environments/"+environmentId, EnvironmentScheduling{}, &EnvironmentScheduling{})
	if err != nil {
		return err
	}

	return nil
}
