---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powermax_host Resource - terraform-provider-powermax"
subcategory: ""
description: |-
  Resource to manage hosts in PowerMax array. Updates are supported for the following parameters: name, initiators, host_flags, consistent_lun.
---

# powermax_host (Resource)

Resource to manage hosts in PowerMax array. Updates are supported for the following parameters: `name`, `initiators`, `host_flags`, `consistent_lun`.

## Example Usage

```terraform
resource "powermax_host" "host_1" {
	name = "host_1"
	initiators = ["0000000000000001"]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
		consistent_lun = false
	}
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `host_flags` (Attributes) Flags set for the host. When host_flags = {} then default flags will be considered. (see [below for nested schema](#nestedatt--host_flags))
- `initiators` (Set of String) The initiators associated with the host.
- `name` (String) The name of the host.

### Optional

- `consistent_lun` (Boolean) It enables the rejection of any masking operation involving this host that would result in inconsistent LUN values.

### Read-Only

- `bw_limit` (Number) Specifies the bandwidth limit for a host.
- `id` (String) The ID of the host.
- `maskingview` (List of String) The masking views associated with the host.
- `numofhostgroups` (Number) The number of hostgroups associated with the host.
- `numofinitiators` (Number) The number of initiators associated with the host.
- `numofmaskingviews` (Number) The number of masking views associated with the host.
- `numofpowerpathhosts` (Number) The number of powerpath hosts associated with the host.
- `port_flags_override` (Boolean) States whether port flags override is enabled on the host.
- `powerpath_hosts` (List of String) The powerpath hosts associated with the host.
- `type` (String) Specifies the type of host.

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


