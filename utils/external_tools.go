// This file is used to add external tools as dependencies
//go:build utils
// +build utils

package utils

import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
