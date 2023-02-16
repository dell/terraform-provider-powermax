package helper

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"math/big"
	"reflect"
	"strings"
)

func CopyFields(ctx context.Context, source, destination interface{}) error {
	tflog.Debug(ctx, "Copy fields", map[string]interface{}{
		"source":      source,
		"destination": destination,
	})
	sourceValue := reflect.ValueOf(source)
	destinationValue := reflect.ValueOf(destination)

	// Check if destination is a pointer to a struct
	if destinationValue.Kind() != reflect.Ptr || destinationValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination is not a pointer to a struct")
	}

	// if source is a pointer, use the Elem() method to get the value that the pointer points to
	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}

	if sourceValue.Kind() != reflect.Struct {
		return fmt.Errorf("source is not a struct")
	}

	// Get the type of the destination struct
	//destinationType := destinationValue.Elem().Type()
	for i := 0; i < sourceValue.NumField(); i++ {
		sourceFieldName := sourceValue.Type().Field(i).Name

		tflog.Debug(ctx, "Converting source field", map[string]interface{}{
			"sourceFieldName": sourceFieldName,
			"sourceFieldKind": sourceValue.Field(i).Kind().String(),
		})

		sourceField := sourceValue.Field(i)
		if sourceField.Kind() == reflect.Ptr {
			sourceField = sourceField.Elem()
		}
		if !sourceField.IsValid() {
			tflog.Error(ctx, "source field is not valid", map[string]interface{}{
				"sourceFieldName": sourceFieldName,
				"sourceField":     sourceField,
			})
			continue
		}

		destinationField := destinationValue.Elem().FieldByName(sourceFieldName)
		if destinationField.IsValid() && destinationField.CanSet() {

			tflog.Debug(ctx, "debugging source field", map[string]interface{}{
				"sourceField Interface": sourceField.Interface(),
			})
			// Convert the source value to the type of the destination field dynamically
			var destinationFieldValue attr.Value

			switch sourceField.Kind() {
			case reflect.String:
				destinationFieldValue = types.StringValue(sourceField.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				destinationFieldValue = types.Int64Value(sourceField.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				destinationFieldValue = types.Int64Value(sourceField.Int())
			case reflect.Float32, reflect.Float64:
				//destinationFieldValue = types.Float64Value(sourceField.Float())
				destinationFieldValue = types.NumberValue(big.NewFloat(sourceField.Float()))
			case reflect.Bool:
				destinationFieldValue = types.BoolValue(sourceField.Bool())
			case reflect.Array, reflect.Slice:
				destinationFieldValue, _ = types.ListValueFrom(nil, types.StringType, sourceField.Interface())
			case reflect.Struct:
				if destinationField.Type().String() == "basetypes.ObjectValue" {
					err := CopyFields(ctx, sourceField.Interface(), &destinationFieldValue)
					if err != nil {
						return fmt.Errorf("failed to copy field object %v", sourceField.Type())
					}
				}
				if destinationField.Type().String() == "basetypes.MapValue" {
					destinationFieldValue, _ = types.MapValue(types.StringType, structToMap(sourceField.Interface()))
				}
			default:
				tflog.Error(ctx, "unsupported source field type", map[string]interface{}{
					"sourceField": sourceField,
				})
				continue
			}

			if destinationField.Type() == reflect.TypeOf(destinationFieldValue) {
				destinationField.Set(reflect.ValueOf(destinationFieldValue))
			}
		}
	}

	return nil
}

func structToMap(s interface{}) map[string]attr.Value {
	result := make(map[string]attr.Value)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			jsonTag := t.Field(i).Tag.Get("json")
			if len(jsonTag) > 0 {
				fieldTag := strings.Split(jsonTag, ",")[0]
				result[fieldTag] = types.StringValue(v.Field(i).String())
			} else {
				result[t.Field(i).Name] = types.StringValue(v.Field(i).String())
			}

		}
	}
	return result
}
