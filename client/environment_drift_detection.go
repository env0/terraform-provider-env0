package client

func (self *ApiClient) EnvironmentDriftDetection(environmentId string) (EnvironmentSchedulingExpression, error) {
	var result EnvironmentSchedulingExpression

	err := self.http.Get("/scheduling/drift-detection/environments/"+environmentId, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self *ApiClient) EnvironmentUpdateDriftDetection(environmentId string, payload EnvironmentSchedulingExpression) (EnvironmentSchedulingExpression, error) {
	var result EnvironmentSchedulingExpression

	err := self.http.Patch("/scheduling/drift-detection/environments/"+environmentId, payload, &result)
	if err != nil {
		return EnvironmentSchedulingExpression{}, err
	}

	return result, nil
}

func (self *ApiClient) EnvironmentStopDriftDetection(environmentId string) error {
	err := self.http.Patch("/scheduling/drift-detection/environments/"+environmentId, EnvironmentSchedulingExpression{Enabled: false}, &EnvironmentScheduling{})
	if err != nil {
		return err
	}

	return nil
}
