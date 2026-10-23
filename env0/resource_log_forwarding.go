package env0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

func resourceLogForwarding() *schema.Resource {
	logForwardingTypes := make([]string, len(client.LogForwardingTypes))
	for i, logForwardingType := range client.LogForwardingTypes {
		logForwardingTypes[i] = string(logForwardingType)
	}

	return &schema.Resource{
		Description: "Forwards deployment logs, audit logs or both to an external logging service. " +
			"An organization has at most one log forwarding configuration per type. " +
			"The env0 API never returns the configuration value, so Terraform cannot detect changes made to it outside Terraform. " +
			"An import cannot read the value either: the first apply after an import writes the configured value.",

		CreateContext: resourceLogForwardingCreate,
		ReadContext:   resourceLogForwardingRead,
		UpdateContext: resourceLogForwardingUpdate,
		DeleteContext: resourceLogForwardingDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceLogForwardingImport},

		CustomizeDiff: resourceLogForwardingCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:             schema.TypeString,
				Description:      "the logging service: " + strings.Join(logForwardingTypes, ", "),
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: NewStringInValidator(logForwardingTypes),
			},
			"value": {
				Type: schema.TypeString,
				Description: "the logging service configuration, as a JSON object (use `jsonencode`). Its fields depend on the type, see https://docs.envzero.com/docs/logs-forwarding. " +
					"With deployment log forwarding, secret fields must be secret references such as `$${ssm:<name>}` unless the organization allows saving secrets on env0. " +
					"Not returned by the env0 API, so changes made outside Terraform are not detected",
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: validateLogForwardingValue,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
			"deployment_log_forwarding": {
				Type:        schema.TypeBool,
				Description: "forward deployment logs (at least one of deployment_log_forwarding and audit_log_forwarding must be true)",
				Optional:    true,
				Default:     false,
			},
			"audit_log_forwarding": {
				Type:        schema.TypeBool,
				Description: "forward audit logs (at least one of deployment_log_forwarding and audit_log_forwarding must be true)",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func validateLogForwardingValue(i any, path cty.Path) diag.Diagnostics {
	var value map[string]any
	if err := json.Unmarshal([]byte(i.(string)), &value); err != nil || value == nil {
		return diag.Errorf("value must be a JSON object")
	}

	return nil
}

func resourceLogForwardingCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta any) error {
	if !d.NewValueKnown("deployment_log_forwarding") || !d.NewValueKnown("audit_log_forwarding") {
		return nil
	}

	if !d.Get("deployment_log_forwarding").(bool) && !d.Get("audit_log_forwarding").(bool) {
		return errors.New("at least one of deployment_log_forwarding and audit_log_forwarding must be true")
	}

	return nil
}

func getLogForwardingPayload(d *schema.ResourceData) (*client.LogForwardingUpsertPayload, error) {
	var value map[string]any
	if err := json.Unmarshal([]byte(d.Get("value").(string)), &value); err != nil {
		return nil, fmt.Errorf("value must be a JSON object: %w", err)
	}

	return &client.LogForwardingUpsertPayload{
		Type:                    client.LogForwardingType(d.Get("type").(string)),
		Value:                   value,
		DeploymentLogForwarding: d.Get("deployment_log_forwarding").(bool),
		AuditLogForwarding:      d.Get("audit_log_forwarding").(bool),
	}, nil
}

func getLogForwarding(meta any, matches func(client.LogForwarding) bool) (*client.LogForwarding, error) {
	apiClient := meta.(client.ApiClientInterface)

	logForwardings, err := apiClient.LogForwardings()
	if err != nil {
		return nil, err
	}

	for _, logForwarding := range logForwardings {
		if matches(logForwarding) {
			return &logForwarding, nil
		}
	}

	return nil, ErrNotFound
}

func getLogForwardingById(id string, meta any) (*client.LogForwarding, error) {
	return getLogForwarding(meta, func(logForwarding client.LogForwarding) bool {
		return logForwarding.Id == id
	})
}

func getLogForwardingByType(logForwardingType client.LogForwardingType, meta any) (*client.LogForwarding, error) {
	return getLogForwarding(meta, func(logForwarding client.LogForwarding) bool {
		return logForwarding.Type == logForwardingType
	})
}

func resourceLogForwardingCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload, err := getLogForwardingPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API upserts by organization and type, so creating over an existing configuration would silently take it over.
	existing, err := getLogForwardingByType(payload.Type, meta)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return diag.Errorf("could not get log forwarding configurations: %v", err)
	}

	if existing != nil {
		return diag.Errorf("a log forwarding configuration of type %s already exists (id %s). Import it instead: terraform import <resource address> %s", payload.Type, existing.Id, payload.Type)
	}

	logForwarding, err := apiClient.LogForwardingUpsert(*payload)
	if err != nil {
		return diag.Errorf("could not create log forwarding configuration: %v", err)
	}

	d.SetId(logForwarding.Id)

	return nil
}

func resourceLogForwardingRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	logForwarding, err := getLogForwardingById(d.Id(), meta)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
			d.SetId("")

			return nil
		}

		return diag.Errorf("could not get log forwarding configuration: %v", err)
	}

	if err := writeResourceData(logForwarding, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceLogForwardingUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload, err := getLogForwardingPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}

	logForwarding, err := apiClient.LogForwardingUpsert(*payload)
	if err != nil {
		return diag.Errorf("could not update log forwarding configuration: %v", err)
	}

	// The upsert creates a new configuration, with a new id, when this one was deleted outside Terraform.
	d.SetId(logForwarding.Id)

	return nil
}

func resourceLogForwardingDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.LogForwardingDelete(d.Id()); err != nil {
		var failedResponseError *http.FailedResponseError
		if errors.As(err, &failedResponseError) && failedResponseError.NotFound() {
			return nil
		}

		return diag.Errorf("could not delete log forwarding configuration: %v", err)
	}

	return nil
}

func resourceLogForwardingImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id := d.Id()

	var (
		logForwarding *client.LogForwarding
		err           error
	)

	if _, uuidErr := uuid.Parse(id); uuidErr == nil {
		tflog.Info(ctx, "Resolving log forwarding configuration by id", map[string]any{"id": id})

		logForwarding, err = getLogForwardingById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving log forwarding configuration by type", map[string]any{"type": id})

		logForwarding, err = getLogForwardingByType(client.LogForwardingType(strings.ToUpper(id)), meta)
	}

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("log forwarding configuration %v not found", id)
		}

		return nil, err
	}

	if err := writeResourceData(logForwarding, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
