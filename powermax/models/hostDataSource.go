// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// HostDataSourceModel describes the datasource model.
type HostDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	HostIDs types.List   `tfsdk:"host_ids"`
}
