package client

func (client *ApiClient) EnvironmentDriftDetection(environmentId string) (EnvironmentSchedulingExpression, error) {
	var result EnvironmentSchedulingExpression

	err := client.http.Get("/scheduling/drift-detection/environments/"+environmentId, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (client *ApiClient) EnvironmentUpdateDriftDetection(environmentId string, payload EnvironmentSchedulingExpression) (EnvironmentSchedulingExpression, error) {
	var result EnvironmentSchedulingExpression

	err := client.http.Patch("/scheduling/drift-detection/environments/"+environmentId, payload, &result)
	if err != nil {
		return EnvironmentSchedulingExpression{}, err
	}

	return result, nil
}

func (client *ApiClient) EnvironmentStopDriftDetection(environmentId string) error {
	err := client.http.Patch("/scheduling/drift-detection/environments/"+environmentId, EnvironmentSchedulingExpression{Enabled: false}, &EnvironmentScheduling{})
	if err != nil {
		return err
	}

	return nil
}
