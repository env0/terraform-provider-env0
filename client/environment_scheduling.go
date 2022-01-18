package client

import "errors"

func (self *ApiClient) EnvironmentScheduling(environmentId string) (EnvironmentScheduling, error) {
	var result EnvironmentScheduling

	err := self.http.Get("/scheduling/environments/"+environmentId, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self *ApiClient) EnvironmentSchedulingUpdate(environmentId string, payload EnvironmentScheduling) (EnvironmentScheduling, error) {
	var result EnvironmentScheduling

	if payload.Deploy.Cron == payload.Destroy.Cron {
		return EnvironmentScheduling{}, errors.New("deploy and destroy cron expressions must not be the same")
	}

	err := self.http.Put("/scheduling/environments/"+environmentId, payload, &result)
	if err != nil {
		return EnvironmentScheduling{}, err
	}

	return result, nil
}

func (self *ApiClient) EnvironmentSchedulingDelete(environmentId string) error {
	err := self.http.Put("/scheduling/environments/"+environmentId, EnvironmentScheduling{}, nil)
	if err != nil {
		return err
	}

	return nil
}
