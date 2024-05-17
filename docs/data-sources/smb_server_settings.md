---
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
#
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://mozilla.org/MPL/2.0/
#
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powerscale_smb_server_settings data source"
linkTitle: "powerscale_smb_server_settings"
page_title: "powerscale_smb_server_settings Data Source - terraform-provider-powerscale"
subcategory: ""
description: |-
  This datasource is used to query the SMB Server Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.
---

# powerscale_smb_server_settings (Data Source)

This datasource is used to query the SMB Server Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.

## Example Usage

```terraform
/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# Returns PowerScale SMB Server Settings based on filter
data "powerscale_smb_server_settings" "test" {
  filter {
    # Used for query parameter, supported by PowerScale Platform API
    scope = "effective"
  }
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_smb_server_settings.test
output "powerscale_smb_server_settings_test" {
  value = data.powerscale_smb_server_settings.test
}

# Returns SMB Server Settings
data "powerscale_smb_server_settings" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_smb_server_settings.all
output "powerscale_smb_server_settings_all" {
  value = data.powerscale_smb_server_settings.all
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `id` (String) ID of SMB Server Settings.
- `smb_server_settings` (Attributes) SMB Server Settings (see [below for nested schema](#nestedatt--smb_server_settings))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `scope` (String) If specified as "effective" or not specified, all fields are returned.  If specified as "user", only fields with non-default values are shown.  If specified as "default", the original values are returned.


<a id="nestedatt--smb_server_settings"></a>
### Nested Schema for `smb_server_settings`

Read-Only:

- `access_based_share_enum` (Boolean) Only enumerate files and folders the requesting user has access to.
- `audit_fileshare` (String) Specify level of file share audit events to log.
- `audit_logon` (String) Specify the level of logon audit events to log.
- `dot_snap_accessible_child` (Boolean) Allow access to .snapshot directories in share subdirectories.
- `dot_snap_accessible_root` (Boolean) Allow access to the .snapshot directory in the root of the share.
- `dot_snap_visible_child` (Boolean) Show .snapshot directories in share subdirectories.
- `dot_snap_visible_root` (Boolean) Show the .snapshot directory in the root of a share.
- `enable_security_signatures` (Boolean) Indicates whether the server supports signed SMB packets.
- `guest_user` (String) Specifies the fully-qualified user to use for guest access.
- `ignore_eas` (Boolean) Specify whether to ignore EAs on files.
- `onefs_cpu_multiplier` (Number) Specify the number of OneFS driver worker threads per CPU.
- `onefs_num_workers` (Number) Set the maximum number of OneFS driver worker threads.
- `reject_unencrypted_access` (Boolean) If SMB3 encryption is enabled, reject unencrypted access from clients.
- `require_security_signatures` (Boolean) Indicates whether the server requires signed SMB packets.
- `server_side_copy` (Boolean) Enable Server Side Copy.
- `server_string` (String) Provides a description of the server.
- `service` (Boolean) Specify whether service is enabled.
- `srv_cpu_multiplier` (Number) Specify the number of SRV service worker threads per CPU.
- `srv_num_workers` (Number) Set the maximum number of SRV service worker threads.
- `support_multichannel` (Boolean) Support multichannel.
- `support_netbios` (Boolean) Support NetBIOS.
- `support_smb2` (Boolean) Support the SMB2 protocol on the server.
- `support_smb3_encryption` (Boolean) Support the SMB3 encryption on the server.