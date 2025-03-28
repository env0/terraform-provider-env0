package client

type KubernetesCredentialsType string

const (
	KubeconfigCredentialsType KubernetesCredentialsType = "K8S_KUBECONFIG_FILE"
	AwsEksCredentialsType     KubernetesCredentialsType = "K8S_EKS_AUTH"
	AzureAksCredentialsType   KubernetesCredentialsType = "K8S_AZURE_AKS_AUTH"
	GcpGkeCredentialsType     KubernetesCredentialsType = "K8S_GCP_GKE_AUTH"
)

type KubernetesCredentialsCreatePayload struct {
	Name  string                    `json:"name"`
	Type  KubernetesCredentialsType `json:"type"`
	Value any                       `json:"value"`
}

type KubernetesCredentialsUpdatePayload struct {
	Type  KubernetesCredentialsType `json:"type"`
	Value any                       `json:"value"`
}

// K8S_KUBECONFIG_FILE
type KubeconfigFileValue struct {
	KubeConfig string `json:"kubeConfig"`
}

// K8S_EKS_AUTH
type AwsEksValue struct {
	ClusterName   string `json:"clusterName"`
	ClusterRegion string `json:"clusterRegion"`
}

// K8S_AZURE_AKS_AUTH
type AzureAksValue struct {
	ClusterName   string `json:"clusterName"`
	ResourceGroup string `json:"resourceGroup"`
}

// K8S_GCP_GKE_AUTH
type GcpGkeValue struct {
	ClusterName   string `json:"clusterName"`
	ComputeRegion string `json:"computeRegion"`
}

func (client *ApiClient) KubernetesCredentialsCreate(payload *KubernetesCredentialsCreatePayload) (*Credentials, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizatioId := &struct {
		OrganizationId string `json:"organizationId"`
		Name           string `json:"name"`
		Type           string `json:"type"`
		Value          any    `json:"value"`
	}{
		OrganizationId: organizationId,
		Name:           payload.Name,
		Type:           string(payload.Type),
		Value:          payload.Value,
	}

	var result Credentials
	if err := client.http.Post("/credentials", payloadWithOrganizatioId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) KubernetesCredentialsUpdate(id string, payload *KubernetesCredentialsUpdatePayload) (*Credentials, error) {
	var result Credentials
	if err := client.http.Patch("/credentials/"+id, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
