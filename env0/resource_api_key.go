package env0

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	organizationRoleAdmin = "Admin"
	organizationRoleUser  = "User"
)

func resourceApiKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiKeyCreate,
		ReadContext:   resourceApiKeyRead,
		UpdateContext: resourceApiKeyUpdate,
		DeleteContext: resourceApiKeyDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceApiKeyImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the api key name",
				Required:    true,
				ForceNew:    true,
			},
			"organization_role": {
				Type:        schema.TypeString,
				Description: "the api key type. 'Admin', 'User' or a custom role id. Defaults to 'Admin'. For more details check https://docs.env0.com/docs/api-keys",
				Default:     organizationRoleAdmin,
				Optional:    true,
				ForceNew:    true,
				ValidateDiagFunc: func(i any, p cty.Path) diag.Diagnostics {
					val := i.(string)
					if val != organizationRoleAdmin && val != organizationRoleUser {
						_, err := uuid.Parse(val)
						if err != nil {
							return diag.Errorf("Organization role must be either 'Admin', 'User' or a custom role id, got: %s", val)
						}
					}

					return nil
				},
			},
			"project_permissions": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Project-specific permissions. Only valid when organization_role is 'User' or a custom role id",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The project ID to assign permissions to",
						},
						"project_role": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							Description:      "The role for this project. Must be one of: Planner, Viewer, Deployer, Admin",
							ValidateDiagFunc: NewStringInValidator([]string{"Planner", "Viewer", "Deployer", "Admin"}),
						},
					},
				},
			},
			"omit_api_key_secret": {
				Type:        schema.TypeBool,
				Description: "if set to 'true' will omit the api_key_secret from the state. This would mean that the api_key_secret cannot be used",
				Optional:    true,
				ForceNew:    true,
			},
			"api_key_secret": {
				Type:        schema.TypeString,
				Description: "the api key secret. This attribute is not computed for imported resources. Note that this will be written to the state file. To omit the secret: set 'omit_api_key_secret' to 'true'",
				Computed:    true,
				Sensitive:   true,
			},
			"api_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the api key id",
			},
			"ip_allowlist": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    50,
				Description: "the IP addresses this api key may be used from: an address ('198.51.100.7'), a CIDR subnet ('198.51.100.0/28') or a range ('198.51.100.40-198.51.100.60'), IPv4 or IPv6. Applies on top of the organization's IP allowlist, so an empty or omitted list only removes the per-key restriction. A change can take up to an hour to apply on the API. Terraform owns the list: one set in the env0 UI shows as drift and the next apply clears it",
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateIpAllowlistEntry,
				},
			},
		},
	}
}

func validateIpAllowlistEntry(i any, p cty.Path) diag.Diagnostics {
	entry := i.(string)
	if err := checkIpAllowlistEntry(entry); err != nil {
		return diag.Errorf("invalid ip_allowlist entry '%s': %v. Use an address (198.51.100.7), a CIDR subnet (198.51.100.0/28) or a range (198.51.100.40-198.51.100.60)", entry, err)
	}

	return nil
}

// netip matches the server's checks: it rejects IPv4 octets with leading zeros and any whitespace,
// so an entry that passes is stored and read back unchanged.
func checkIpAllowlistEntry(entry string) error {
	if strings.Contains(entry, "/") {
		_, err := netip.ParsePrefix(entry)

		return err
	}

	if startText, endText, isRange := strings.Cut(entry, "-"); isRange {
		start, err := parseIpAllowlistAddr(startText)
		if err != nil {
			return err
		}

		end, err := parseIpAllowlistAddr(endText)
		if err != nil {
			return err
		}

		if start.BitLen() != end.BitLen() {
			return errors.New("both ends of a range must be in the same IP family")
		}

		if start.Compare(end) > 0 {
			return errors.New("the range start is after its end")
		}

		return nil
	}

	_, err := parseIpAllowlistAddr(entry)

	return err
}

func parseIpAllowlistAddr(text string) (netip.Addr, error) {
	addr, err := netip.ParseAddr(text)
	if err == nil && addr.Zone() != "" {
		return addr, errors.New("IPv6 zones are not supported")
	}

	return addr, err
}

// Same conversion as the env0 UI. Entries are validated already.
func toIpAllowlistRules(entries []any) []client.IpAllowlistRule {
	rules := make([]client.IpAllowlistRule, 0, len(entries))

	for _, e := range entries {
		entry := e.(string)

		family := "ipv4"
		if strings.Contains(entry, ":") {
			family = "ipv6"
		}

		if network, prefixText, isSubnet := strings.Cut(entry, "/"); isSubnet {
			prefix, _ := strconv.Atoi(prefixText)
			rules = append(rules, client.IpAllowlistRule{Type: "subnet", Network: network, Prefix: &prefix, Family: family})
		} else if start, end, isRange := strings.Cut(entry, "-"); isRange {
			rules = append(rules, client.IpAllowlistRule{Type: "range", Start: start, End: end, Family: family})
		} else {
			rules = append(rules, client.IpAllowlistRule{Type: "address", Address: entry, Family: family})
		}
	}

	return rules
}

func fromIpAllowlistRules(rules []client.IpAllowlistRule) []string {
	entries := make([]string, 0, len(rules))

	for _, rule := range rules {
		switch rule.Type {
		case "subnet":
			prefix := 0
			if rule.Prefix != nil {
				prefix = *rule.Prefix
			}

			entries = append(entries, fmt.Sprintf("%s/%d", rule.Network, prefix))
		case "range":
			entries = append(entries, rule.Start+"-"+rule.End)
		default:
			entries = append(entries, rule.Address)
		}
	}

	return entries
}

func writeApiKeyResourceData(apiKey *client.ApiKey, d *schema.ResourceData) error {
	if err := writeResourceData(apiKey, d); err != nil {
		return err
	}

	return d.Set("ip_allowlist", fromIpAllowlistRules(apiKey.IpAllowlist))
}

func resourceApiKeyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	organizationRole := d.Get("organization_role").(string)
	if organizationRole == organizationRoleAdmin && len(d.Get("project_permissions").(*schema.Set).List()) > 0 {
		return diag.Errorf("project_permissions cannot be set when organization_role is Admin")
	}

	var payload client.ApiKeyCreatePayload

	payload.Name = d.Get("name").(string)

	if err := readResourceData(&payload.Permissions, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiKey, err := apiClient.ApiKeyCreate(payload)
	if err != nil {
		return diag.Errorf("could not create api key: %v", err)
	}

	var ipAllowlist []client.IpAllowlistRule

	if entries := d.Get("ip_allowlist").(*schema.Set).List(); len(entries) > 0 {
		withAllowlist, err := apiClient.ApiKeyUpdateIpAllowlist(apiKey.Id, toIpAllowlistRules(entries))
		if err != nil {
			// The key was never restricted, so it must not be left behind.
			if deleteErr := apiClient.ApiKeyDelete(apiKey.Id); deleteErr != nil {
				return diag.Errorf("could not set the api key ip_allowlist: %v. Deleting the new api key %s also failed, delete it by hand: %v", err, apiKey.Id, deleteErr)
			}

			return diag.Errorf("could not set the api key ip_allowlist: %v", err)
		}

		ipAllowlist = withAllowlist.IpAllowlist
	}

	if err := writeResourceData(apiKey, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	if err := d.Set("ip_allowlist", fromIpAllowlistRules(ipAllowlist)); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	if omit, ok := d.GetOk("omit_api_key_secret"); ok && omit.(bool) {
		d.Set("api_key_secret", "omitted")
	}

	return nil
}

func resourceApiKeyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiKey, err := getApiKeyById(d.Id(), meta)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
			d.SetId("")

			return nil
		}

		return diag.Errorf("could not get api key: %v", err)
	}

	apiKey.ApiKeySecret = "" // Don't override the api key secret currently in the state.

	if err := writeApiKeyResourceData(apiKey, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

// Every other attribute forces a new key, so only the ip_allowlist is updated in place.
func resourceApiKeyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	apiKey, err := apiClient.ApiKeyUpdateIpAllowlist(d.Id(), toIpAllowlistRules(d.Get("ip_allowlist").(*schema.Set).List()))
	if err != nil {
		return diag.Errorf("could not update api key ip_allowlist: %v", err)
	}

	if err := d.Set("ip_allowlist", fromIpAllowlistRules(apiKey.IpAllowlist)); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceApiKeyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.ApiKeyDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete api key: %v", err)
	}

	return nil
}

func getApiKeyById(id string, meta any) (*client.ApiKey, error) {
	apiClient := meta.(client.ApiClientInterface)

	apiKeys, err := apiClient.ApiKeys()
	if err != nil {
		return nil, err
	}

	for _, apiKey := range apiKeys {
		if apiKey.Id == id {
			return &apiKey, nil
		}
	}

	return nil, ErrNotFound
}

func getApiKeyByName(name string, meta any) (*client.ApiKey, error) {
	apiClient := meta.(client.ApiClientInterface)

	apiKeys, err := apiClient.ApiKeys()
	if err != nil {
		return nil, err
	}

	var foundApiKeys []client.ApiKey

	for _, apiKey := range apiKeys {
		if apiKey.Name == name {
			foundApiKeys = append(foundApiKeys, apiKey)
		}
	}

	if len(foundApiKeys) == 0 {
		return nil, fmt.Errorf("api key with name %v not found", name)
	}

	if len(foundApiKeys) > 1 {
		return nil, fmt.Errorf("found multiple api keys with name: %s. Use id instead or make sure api key names are unique %v", name, foundApiKeys)
	}

	return &foundApiKeys[0], nil
}

func getApiKey(ctx context.Context, id string, meta any) (*client.ApiKey, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		tflog.Info(ctx, "Resolving api key by id", map[string]any{"id": id})

		return getApiKeyById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving api key by name", map[string]any{"name": id})

		return getApiKeyByName(id, meta)
	}
}

func resourceApiKeyImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	apiKey, err := getApiKey(ctx, d.Id(), meta)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("api key with id %v not found", d.Id())
		}

		return nil, err
	}

	if err := writeApiKeyResourceData(apiKey, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
