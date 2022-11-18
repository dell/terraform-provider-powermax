package powermax

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func updateVolState(volState *models.Volume, volResponse *pmaxTypes.Volume, volPlan *models.Volume, operation string) {
	volState.ID.Value = volResponse.VolumeID
	if operation == "create" || operation == "update" {
		volState.CapUnit.Value = volPlan.CapUnit.Value
		volState.SGName.Value = volPlan.SGName.Value
	} else if operation == "import" {
		volState.CapUnit.Value = "GB"
		volState.SGName.Value = ""
	}
	volState.Name.Value = volResponse.VolumeIdentifier
	volState.AllocatedPercent.Value = int64(volResponse.AllocatedPercent)
	volState.EffectiveWWN.Value = volResponse.EffectiveWWN
	volState.Emulation.Value = volResponse.Emulation
	volState.Encapsulated.Value = volResponse.Encapsulated
	volState.EncapsulatedWWN.Value = volResponse.EncapsulatedWWN
	volState.HasEffectiveWWN.Value = volResponse.HasEffectiveWWN
	volState.EnableMobilityID.Value = volResponse.MobilityIDEnabled
	volState.NGUID.Value = volResponse.NGUID
	volState.NumOfFrontEndPaths.Value = int64(volResponse.NumberOfFrontEndPaths)
	volState.NumOfStorageGroups.Value = int64(volResponse.NumberOfStorageGroups)
	volState.OracleInstanceName.Value = volResponse.OracleInstanceName
	volState.Pinned.Value = volResponse.Pinned
	volState.Reserved.Value = volResponse.Reserved
	volState.SSID.Value = volResponse.SSID
	switch volState.CapUnit.Value {
	case "TB":
		volState.Size = types.Number{Value: big.NewFloat(volResponse.CapacityGB / 1024)}
	case "CYL":
		volState.Size = types.Number{Value: big.NewFloat(float64(volResponse.CapacityCYL))}
	default:
		volState.Size = types.Number{Value: big.NewFloat(volResponse.CapacityGB)}
	}

	volState.SnapSource.Value = volResponse.SnapSource
	volState.SnapTarget.Value = volResponse.SnapTarget
	volState.Status.Value = volResponse.Status
	volState.Type.Value = volResponse.Type
	volState.UnreducibleDataGB.Value = volResponse.UnreducibleDataGB
	volState.WWN.Value = volResponse.WWN

	var sgList []attr.Value
	for _, sg := range volResponse.StorageGroupIDList {
		sgList = append(sgList, types.String{Value: sg})
	}
	volState.StorageGroupIDs = types.List{
		ElemType: types.StringType,
		Elems:    sgList,
	}

	symmetrixPortKeysTfsdk := types.List{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"director_id": types.StringType,
				"port_id":     types.StringType,
			},
		},
	}

	symmetrixPortKeysObjects := []attr.Value{}

	if len(volResponse.SymmetrixPortKey) > 0 {
		symmetrixPortkeyDetails := make(map[string]attr.Value)

		for _, symmetrixPortkey := range volResponse.SymmetrixPortKey {
			symmetrixPortkeyDetails["director_id"] = types.String{Value: symmetrixPortkey.DirectorID}
			symmetrixPortkeyDetails["port_id"] = types.String{Value: symmetrixPortkey.PortID}
			symmetrixPortkeyObject := types.Object{
				Attrs: symmetrixPortkeyDetails,
				AttrTypes: map[string]attr.Type{
					"director_id": types.StringType,
					"port_id":     types.StringType,
				},
			}
			symmetrixPortKeysObjects = append(symmetrixPortKeysObjects, symmetrixPortkeyObject)
		}
	}

	symmetrixPortKeysTfsdk.Elems = symmetrixPortKeysObjects
	volState.SymmetrixPortKeys = symmetrixPortKeysTfsdk
	rdfGroupsTfsdk := types.List{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"rdf_group_number": types.Int64Type,
				"label":            types.StringType,
			},
		},
	}

	rdfGroupsObjects := []attr.Value{}

	if len(volResponse.RDFGroupIDList) > 0 {
		rdfGroupDetails := make(map[string]attr.Value)

		for _, rdfGroup := range volResponse.RDFGroupIDList {
			rdfGroupDetails["rdf_group_number"] = types.Int64{Value: int64(rdfGroup.RDFGroupNumber)}
			rdfGroupDetails["label"] = types.String{Value: rdfGroup.Label}
			rdfGroupObject := types.Object{
				Attrs: rdfGroupDetails,
				AttrTypes: map[string]attr.Type{
					"rdf_group_number": types.Int64Type,
					"label":            types.StringType,
				},
			}
			rdfGroupsObjects = append(rdfGroupsObjects, rdfGroupObject)
		}
	}

	rdfGroupsTfsdk.Elems = rdfGroupsObjects
	volState.RDFGroupIDs = rdfGroupsTfsdk
}

func updateVol(ctx context.Context, client client.Client, planVol, stateVol models.Volume) ([]string, []string, []string) {
	updatedParameters := []string{}
	updateFailedParameters := []string{}
	errorMessages := []string{}

	if planVol.Name.Value != stateVol.Name.Value {
		_, err := client.PmaxClient.RenameVolume(ctx, client.SymmetrixID, stateVol.ID.Value, planVol.Name.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename volume: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

	if planVol.EnableMobilityID.Value != stateVol.EnableMobilityID.Value {
		_, err := client.PmaxClient.ModifyMobilityForVolume(ctx, client.SymmetrixID, stateVol.ID.Value, planVol.EnableMobilityID.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "enable_mobility_id")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify mobility: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "enable_mobility_id")
		}
	}

	if planVol.Size.Value.String() != stateVol.Size.Value.String() {
		size, err := getVolumeSize(planVol)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "size")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the volume size: %s", err.Error()))
			return updatedParameters, updateFailedParameters, errorMessages
		}
		_, err = client.PmaxClient.ExpandVolume(ctx, client.SymmetrixID, stateVol.ID.Value, 0, size, planVol.CapUnit.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "size")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the volume size: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "size")
		}
	}
	return updatedParameters, updateFailedParameters, errorMessages
}

func getVolumeSize(volume models.Volume) (interface{}, error) {
	var size interface{}
	if volume.CapUnit.Value == "CYL" {
		if volume.Size.Value.IsInt() {
			intVal, err := strconv.Atoi(volume.Size.Value.String())
			if err != nil {
				return size, err
			}
			size = intVal
		} else {
			return size, fmt.Errorf("when cap_unit is 'CYL', size should be defined only as an integer")
		}

	} else {
		size = volume.Size.Value.String()
	}
	return size, nil
}
