package powermax

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributeValidator = &validCapUnitValidator{}

type validCapUnitValidator struct {
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validCapUnitValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Only one of the following values are supported for capacity: %s", ValidCapUnits)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validCapUnitValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Only one of the following values are supported for capacity: %s", ValidCapUnits)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v validCapUnitValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var capUnit types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &capUnit)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if capUnit.Unknown || capUnit.Null {
		return
	}

	validCapacityUnits := strings.Split(ValidCapUnits, ",")
	for _, validCapacityUnit := range validCapacityUnits {
		if capUnit.Value == validCapacityUnit {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		"Unsupported capacity unit for volume size",
		fmt.Sprintf("Allowed values for capacity unit are  :  %s", ValidCapUnits),
	)

}
