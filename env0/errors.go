package env0

import (
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func driftDetected(d *schema.ResourceData, err error) bool {
	if frerr, ok := err.(*http.FailedResponseError); ok && frerr.NotFound() {
		return true
	}

	if _, ok := err.(*client.NotFoundError); ok {
		return true
	}

	return false
}

func ResourceGetFailure(resourceName string, d *schema.ResourceData, err error) diag.Diagnostics {
	if driftDetected(d, err) {
		log.Printf("[WARN] Drift Detected: Terraform will remove %s from state", d.Id())
		d.SetId("")
		return nil
	}

	return diag.Errorf("could not get %s: %v", resourceName, err)
}

func DataGetFailure(dataName string, id interface{}, err error) diag.Diagnostics {
	if frerr, ok := err.(*http.FailedResponseError); ok && frerr.NotFound() {
		return diag.Errorf("could not read %s: id %v not found", dataName, id)
	}

	return diag.Errorf("could not read %s: %v", dataName, err)
}
