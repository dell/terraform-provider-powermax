---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powermax_host_group Resource - terraform-provider-powermax"
subcategory: ""
description: |-
  Resource to manage hostgroup in PowerMax array. Updates are supported for the following parameters: name, host_ids, host_flags, consistent_lun.
---

# powermax_host_group (Resource)

Resource to manage hostgroup in PowerMax array. Updates are supported for the following parameters: `name`, `host_ids`, `host_flags`, `consistent_lun`.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `host_flags` (Attributes) Host Flags set for the hostgroup. When host_flags = {} then default flags will be considered. (see [below for nested schema](#nestedatt--host_flags))
- `host_ids` (Set of String) The masking views associated with the hostgroup.
- `name` (String) The name of the hostgroup.

### Optional

- `consistent_lun` (Boolean) It enables the rejection of any masking operation involving this hostgroup that would result in inconsistent LUN values.

### Read-Only

- `id` (String) The ID of the hostgroup.
- `maskingviews` (List of String) The masking views associated with the hostgroup.
- `numofhosts` (Number) The number of hosts associated with the hostgroup.
- `numofinitiators` (Number) The number of initiators associated with the hostgroup.
- `numofmaskingviews` (Number) The number of masking views associated with the hostgroup.
- `port_flags_override` (Boolean) States whether port flags override is enabled on the hostgroup.
- `type` (String) Specifies the type of hostgroup.

<a id="nestedatt--host_flags"></a>
### Nested Schema for `host_flags`

Optional:

- `avoid_reset_broadcast` (Attributes) It enables a SCSI bus reset to only occur to the port that received the reset. (see [below for nested schema](#nestedatt--host_flags--avoid_reset_broadcast))
- `disable_q_reset_on_ua` (Attributes) It is used for hosts that do not expect the queue to be flushed on a 0629 sense. (see [below for nested schema](#nestedatt--host_flags--disable_q_reset_on_ua))
- `environ_set` (Attributes) It enables the environmental error reporting by the storage system to the host on the specific port. (see [below for nested schema](#nestedatt--host_flags--environ_set))
- `openvms` (Attributes) This attribute enables an Open VMS fibre connection. (see [below for nested schema](#nestedatt--host_flags--openvms))
- `scsi_3` (Attributes) Alters the inquiry data to report that the storage system supports the SCSI-3 protocol. (see [below for nested schema](#nestedatt--host_flags--scsi_3))
- `scsi_support1` (Attributes) This attribute provides a stricter compliance with SCSI standards. (see [below for nested schema](#nestedatt--host_flags--scsi_support1))
- `spc2_protocol_version` (Attributes) When setting this flag, the port must be offline. (see [below for nested schema](#nestedatt--host_flags--spc2_protocol_version))
- `volume_set_addressing` (Attributes) It enables the volume set addressing mode. (see [below for nested schema](#nestedatt--host_flags--volume_set_addressing))

<a id="nestedatt--host_flags--avoid_reset_broadcast"></a>
### Nested Schema for `host_flags.avoid_reset_broadcast`

Optional:

- `enabled` (Boolean)
- `override` (Boolean)


<a id="nestedatt--host_flags--disable_q_reset_on_ua"></a>
### Nested Schema for `host_flags.disable_q_reset_on_ua`

Optional:

- `enabled` (Boolean)
- `override` (Boolean)


<a id="nestedatt--host_flags--environ_set"></a>
### Nested Schema for `host_flags.environ_set`

Optional:

- `enabled` (Boolean)
- `override` (Boolean)


<a id="nestedatt--host_flags--openvms"></a>
### Nested Schema for `host_flags.openvms`

Optional:

- `enabled` (Boolean)
- `override` (Boolean)


<a id="nestedatt--host_flags--scsi_3"></a>
### Nested Schema for `host_flags.scsi_3`

Optional:

- `enabled` (Boolean)
- `override` (Boolean)


<a id="nestedatt--host_flags--scsi_support1"></a>
### Nested Schema for `host_flags.scsi_support1`

Optional:

- `enabled` (Boolean)
- `override` (Boolean)


<a id="nestedatt--host_flags--spc2_protocol_version"></a>
### Nested Schema for `host_flags.spc2_protocol_version`

Optional:

- `enabled` (Boolean)
- `override` (Boolean)


<a id="nestedatt--host_flags--volume_set_addressing"></a>
### Nested Schema for `host_flags.volume_set_addressing`

Optional:

- `enabled` (Boolean)
- `override` (Boolean)


