package env0

import (
	"context"
	"errors"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ErrNoChanges = errors.New("no changes")
var ErrNotFound = errors.New("not found")

func driftDetected(err error) bool {
	var failedResponseError *http.FailedResponseError
	if errors.As(err, &failedResponseError) && failedResponseError.NotFound() {
		return true
	}

	var notfoundError *client.NotFoundError

	return errors.As(err, &notfoundError)
}

func ResourceGetFailure(ctx context.Context, resourceName string, d *schema.ResourceData, err error) diag.Diagnostics {
	if driftDetected(err) {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
		d.SetId("")

		return nil
	}

	return diag.Errorf("could not get %s: %v", resourceName, err)
}

func DataGetFailure(dataName string, id any, err error) diag.Diagnostics {
	if frerr, ok := err.(*http.FailedResponseError); ok && frerr.NotFound() {
		return diag.Errorf("could not read %s: id %v not found", dataName, id)
	}

	return diag.Errorf("could not read %s: %v", dataName, err)
}
