package powermax

const (
	// CreateSGDetailErrorMsg specifies error details occurred while creating storage group
	CreateSGDetailErrorMsg = "Could not create storage group "

	// ReadSGDetailsErrorMsg specifies error details occurred while reading storage group
	ReadSGDetailsErrorMsg = "Could not read storage group "

	// ImportSGDetailsErrorMsg specifies error details occurred while importing storage group
	ImportSGDetailsErrorMsg = "Could not import storage group "

	//UpdateSGDetailsErrorMsg specifies error details occurred while updating storage group
	UpdateSGDetailsErrorMsg = "Could not update storage group"

	// DeleteSGDetailsErrorMsg specifies error details occurred while deleting storage group
	DeleteSGDetailsErrorMsg = "Could not delete storage group "

	// CreateVolDetailErrorMsg specifies error details occurred while creating volume
	CreateVolDetailErrorMsg = "Could not create volume "

	// ReadVolDetailsErrorMsg specifies error details occurred while reading volume
	ReadVolDetailsErrorMsg = "Could not read volume "

	// UpdateVolDetailsErrorMsg specifies error details occured while updating volume
	UpdateVolDetailsErrorMsg = "Could not update volume "

	//ImportVolDetailsErrorMsg specifies error details occured while importing volume
	ImportVolDetailsErrorMsg = "Could not import volume "

	// CreateHostDetailErrorMsg specifies error details occurred while creating host
	CreateHostDetailErrorMsg = "Could not create host "

	// ReadHostDetailsErrorMsg specifies error details occurred while reading host
	ReadHostDetailsErrorMsg = "Could not read host "

	// DeleteHostDetailsErrorMsg specifies error details occurred while deleting host
	DeleteHostDetailsErrorMsg = "Could not delete host "

	// ImportHostDetailsErrorMsg specifies error details occurred while importing host
	ImportHostDetailsErrorMsg = "Could not import host "

	// DeleteVolumeDetailsErrorMsg specifies error details occurred while deleting volume
	DeleteVolumeDetailsErrorMsg = "Could not delete volume "

	// RemoveVolumeFromSGDetailsErrorMsg specifies error details occurred while removing volume from storage group
	RemoveVolumeFromSGDetailsErrorMsg = "Could not remove volume from storage group "

	// VolumeSetAddressing flag for host creation
	VolumeSetAddressing = "Volume_Set_Addressing(V)"

	// DisableQResetOnUa flag for host creation
	DisableQResetOnUa = "Disable_Q_Reset_on_UA(D)"

	// AvoidResetBroadcast flag for host creation
	AvoidResetBroadcast = "Avoid_Reset_Broadcast(ARB)"

	// EnvironSet flag for host creation
	EnvironSet = "Environ_Set(E)"

	// OpenVMS flag for host creation
	OpenVMS = "OpenVMS(OVMS)"

	// SCSISupport1 flag for host creation
	SCSISupport1 = "SCSI_Support1(OS2007)"

	// SCSI3 flag for host creation
	SCSI3 = "SCSI_3(SC3)"

	// SPC2ProtocolVersion flag for host creation
	SPC2ProtocolVersion = "SPC2_Protocol_Version(SPC2)"

	// CreatePGDetailErrorMsg specifies error details occurred while creating portgroup
	CreatePGDetailErrorMsg = "Could not create portgroup "

	// ReadPGDetailsErrorMsg specifies error details occurred while reading portgroup
	ReadPGDetailsErrorMsg = "Could not read portgroup "

	// UpdatePGDetailsErrMsg specifies error details occurred while updating portgroup
	UpdatePGDetailsErrMsg = "Could not update portgroup "

	// DeletePGDetailsErrorMsg specifies error details occurred while deleting portgroup
	DeletePGDetailsErrorMsg = "Could not delete portgroup "

	// CreateMVDetailErrorMsg specifies error details occured while creating maskingview
	CreateMVDetailErrorMsg = "Could not create maskingview "

	// RenameMVDetailErrorMsg specifies error details occured while renaming maskingview
	RenameMVDetailErrorMsg = "Could not rename maskingview "

	// ReadMVDetailsErrorMsg specifies error details occurred while reading maskingview
	ReadMVDetailsErrorMsg = "Could not read maskingview "

	// DeleteMVDetailsErrorMsg specifies error details occurred while deleting maskingview
	DeleteMVDetailsErrorMsg = "Could not delete maskingview "

	// ImportMVDetailsErrorMsg specifies error details occurred while importing masking view
	ImportMVDetailsErrorMsg = "Could not import masking view "

	// UpdateStorageGroupDetailsErrorMsg specifies error details occurred while updating storage group
	UpdateStorageGroupDetailsErrorMsg = "Unable to update all changes to StorageGroup"

	// CreateSgErrorMsg specifies error details occurred while creating storage group
	CreateSgErrorMsg = "Error creating storage group"

	// ValidCapUnits specifies the capacity unit supported while provisioning volumes
	ValidCapUnits = "CYL,GB,TB"

	// CreateSGAddVolumeErrMsg specifies error details during create SG with volume id already attached to another storage group
	CreateSGAddVolumeErrMsg = "could not add volumes to storageGroup"
)
