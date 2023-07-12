/*
Copyright (c) 2022-2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package constants

const (

	// ReadPortDetailErrorMsg specifies error details while reading ports.
	ReadPortDetailErrorMsg = "Could not read ports "

	// CreateSGDetailErrorMsg specifies error details occurred while creating storage group.
	CreateSGDetailErrorMsg = "Could not create storage group "

	// ReadSGDetailsErrorMsg specifies error details occurred while reading storage group.
	ReadSGDetailsErrorMsg = "Could not read storage group "

	// ImportSGDetailsErrorMsg specifies error details occurred while importing storage group.
	ImportSGDetailsErrorMsg = "Could not import storage group "

	//UpdateSGDetailsErrorMsg specifies error details occurred while updating storage group.
	UpdateSGDetailsErrorMsg = "Could not update storage group"

	// DeleteSGDetailsErrorMsg specifies error details occurred while deleting storage group.
	DeleteSGDetailsErrorMsg = "Could not delete storage group "

	// CreateVolDetailErrorMsg specifies error details occurred while creating volume.
	CreateVolDetailErrorMsg = "Could not create volume "

	// ReadVolDetailsErrorMsg specifies error details occurred while reading volume.
	ReadVolDetailsErrorMsg = "Could not read volume "

	// UpdateVolDetailsErrorMsg specifies error details occurred while updating volume.
	UpdateVolDetailsErrorMsg = "Could not update volume "

	//ImportVolDetailsErrorMsg specifies error details occurred while importing volume.
	ImportVolDetailsErrorMsg = "Could not import volume "

	// CreateHostDetailErrorMsg specifies error details occurred while creating host.
	CreateHostDetailErrorMsg = "Could not create host "

	// ReadHostDetailsErrorMsg specifies error details occurred while reading host.
	ReadHostDetailsErrorMsg = "Could not read host "

	// UpdateHostDetailsErrorMsg specifies error details occurred while updating host.
	UpdateHostDetailsErrorMsg = "Could not update host "

	// DeleteHostDetailsErrorMsg specifies error details occurred while deleting host.
	DeleteHostDetailsErrorMsg = "Could not delete host "

	// ImportHostDetailsErrorMsg specifies error details occurred while importing host.
	ImportHostDetailsErrorMsg = "Could not import host "

	// CreateSnapshot specifies error while creating snapshot.
	CreateSnapshot = "Could not create snapshot"

	// ReadSnapshots specifies error while reading snapshot.
	ReadSnapshots = "Could not read snapshots"

	// UpdateSnapshot specifies error while updating snapshot.
	UpdateSnapshot = "Could not update the snapshot"

	// DeleteSnapshot specifies error while deleting snapshot.
	DeleteSnapshot = "Could not delete snapshot"

	// CreateHostGroupDetailErrorMsg specifies error details occurred while creating hostgroup.
	CreateHostGroupDetailErrorMsg = "Could not create hostgroup "

	// ReadHostGroupDetailsErrorMsg specifies error details occurred while reading hostgroup.
	ReadHostGroupDetailsErrorMsg = "Could not read hostgroup "

	// ReadHostGroupListDetailsErrorMsg specifies error details occurred while reading hostgroup.
	ReadHostGroupListDetailsErrorMsg = "Could not read hostgroups "

	// UpdateHostGroupDetailsErrorMsg specifies error details occurred while updating hostgroup.
	UpdateHostGroupDetailsErrorMsg = "Could not update hostgroup "

	// DeleteHostGroupDetailsErrorMsg specifies error details occurred while deleting hostgroup.
	DeleteHostGroupDetailsErrorMsg = "Could not delete hostgroup "

	// ImportHostGroupDetailsErrorMsg specifies error details occurred while importing hostgroup.
	ImportHostGroupDetailsErrorMsg = "Could not import hostgroup "

	// DeleteVolumeDetailsErrorMsg specifies error details occurred while deleting volume.
	DeleteVolumeDetailsErrorMsg = "Could not delete volume "

	// RemoveVolumeFromSGDetailsErrorMsg specifies error details occurred while removing volume from storage group.
	RemoveVolumeFromSGDetailsErrorMsg = "Could not remove volume from storage group "

	// VolumeSetAddressing flag for host creation.
	VolumeSetAddressing = "Volume_Set_Addressing(V)"

	// DisableQResetOnUa flag for host creation.
	DisableQResetOnUa = "Disable_Q_Reset_on_UA(D)"

	// AvoidResetBroadcast flag for host creation.
	AvoidResetBroadcast = "Avoid_Reset_Broadcast(ARB)"

	// EnvironSet flag for host creation.
	EnvironSet = "Environ_Set(E)"

	// OpenVMS flag for host creation.
	OpenVMS = "OpenVMS(OVMS)"

	// SCSISupport1 flag for host creation.
	SCSISupport1 = "SCSI_Support1(OS2007)"

	// SCSI3 flag for host creation.
	SCSI3 = "SCSI_3(SC3)"

	// SPC2ProtocolVersion flag for host creation.
	SPC2ProtocolVersion = "SPC2_Protocol_Version(SPC2)"

	// CreatePGDetailErrorMsg specifies error details occurred while creating portgroup.
	CreatePGDetailErrorMsg = "Could not create portgroup "

	// ReadPGDetailsErrorMsg specifies error details occurred while reading portgroup.
	ReadPGDetailsErrorMsg = "Could not read portgroup "

	// ImportPGDetailsErrorMsg specifies error details occurred while importing portgroup.
	ImportPGDetailsErrorMsg = "Could not import portgroup "

	// UpdatePGDetailsErrMsg specifies error details occurred while updating portgroup.
	UpdatePGDetailsErrMsg = "Could not update portgroup "

	// DeletePGDetailsErrorMsg specifies error details occurred while deleting portgroup.
	DeletePGDetailsErrorMsg = "Could not delete portgroup "

	// CreateMVDetailErrorMsg specifies error details occurred while creating maskingview.
	CreateMVDetailErrorMsg = "Could not create maskingview "

	// RenameMVDetailErrorMsg specifies error details occurred while renaming maskingview.
	RenameMVDetailErrorMsg = "Could not rename maskingview "

	// ReadMVDetailsErrorMsg specifies error details occurred while reading maskingview.
	ReadMVDetailsErrorMsg = "Could not read maskingview "

	// DeleteMVDetailsErrorMsg specifies error details occurred while deleting maskingview.
	DeleteMVDetailsErrorMsg = "Could not delete maskingview "

	// ImportMVDetailsErrorMsg specifies error details occurred while importing masking view.
	ImportMVDetailsErrorMsg = "Could not import masking view "

	// CreateSgErrorMsg specifies error details occurred while creating storage group.
	CreateSgErrorMsg = "Error creating storage group"

	// ValidCapUnits specifies the capacity unit supported while provisioning volumes.
	ValidCapUnits = "CYL,GB,TB"

	// CreateSGAddVolumeErrMsg specifies error details during create SG with volume id already attached to another storage group.
	CreateSGAddVolumeErrMsg = "could not add volumes to storageGroup"

	// SweepTestsTemplateIdentifier specifies the string match for all the dangling resources for sweepers.
	SweepTestsTemplateIdentifier = "test_acc_"
	// MinimumSizeValidationError specifies error details returned if the length of the collection is lesser than the specified min size.
	MinimumSizeValidationError = "Required size of the parameter is less than the minimum size: "
)
