package powermax

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// To be removed in the upcoming release and use framework

var _ tfsdk.AttributeValidator = sizeAtLeastValidator{}

// sizeAtLeastValidator validates that list contains at least min elements.
type sizeAtLeastValidator struct {
	min int
}

// Description describes the validation in plain text formatting.
func (v sizeAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("%s%v", MinimumSizeValidationError, v.min)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeAtLeastValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	elems, ok := validateSet(ctx, req, resp)
	if !ok {
		return
	}

	if len(elems) < v.min {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			v.Description(ctx),
			fmt.Sprintf("current size : %d", len(elems)),
		)
	}

}

func validateSet(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) ([]attr.Value, bool) {
	var l types.Set

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &l)

	if diags.HasError() {
		response.Diagnostics.Append(diags...)

		return nil, false
	}

	if l.IsUnknown() || l.IsNull() {
		return nil, false
	}

	return l.Elems, true
}

// SizeAtLeast returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a List.
//   - Contains at least min elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeAtLeast(min int) tfsdk.AttributeValidator {
	return sizeAtLeastValidator{
		min: min,
	}
}
