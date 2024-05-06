package client

type KubernetesCrednetialsType string

const (
	KubeconfigCredentialsType KubernetesCrednetialsType = "K8S_KUBECONFIG_FILE"
	AwsEksCredentialsType     KubernetesCrednetialsType = "K8S_EKS_AUTH"
	AzureAksCredentialsType   KubernetesCrednetialsType = "K8S_AZURE_AKS_AUTH"
	GcpGkeCredentialsType     KubernetesCrednetialsType = "K8S_GCP_GKE_AUTH"
)

type KubernetesCredentialsCreatePayload struct {
	Name  string                    `json:"name"`
	Type  KubernetesCrednetialsType `json:"type"`
	Value interface{}               `json:"value"`
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
		OrganizationId string      `json:"organizationId"`
		Name           string      `json:"name"`
		Type           string      `json:"type"`
		Value          interface{} `json:"value"`
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
