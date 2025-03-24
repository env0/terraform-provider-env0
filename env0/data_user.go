package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUserRead,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Description: "the email of the user",
				Required:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "id of the user",
				Computed:    true,
			},
		},
	}
}

func getUserByEmail(email string, meta any) (*client.User, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)

	organizationUsers, err := apiClient.Users()
	if err != nil {
		return nil, diag.Errorf("Could not get users: %v", err)
	}

	var usersByEmail []client.User

	for _, organizationUser := range organizationUsers {
		if organizationUser.User.Email == email {
			usersByEmail = append(usersByEmail, organizationUser.User)
		}
	}

	if len(usersByEmail) > 1 {
		return nil, diag.Errorf("Found multiple users with the same email: %s", email)
	}

	if len(usersByEmail) == 0 {
		return nil, diag.Errorf("Could not find a user with the email: %s", email)
	}

	return &usersByEmail[0], nil
}

func dataUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	email := d.Get("email").(string)

	user, err := getUserByEmail(email, meta)
	if err != nil {
		return err
	}

	d.SetId(user.UserId)

	return nil
}
