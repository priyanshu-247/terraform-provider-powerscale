---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

title: "powerscale_s3_global_settings resource"
linkTitle: "powerscale_s3_global_settings"
page_title: "powerscale_s3_global_settings Resource - terraform-provider-powerscale"
subcategory: ""
description: |-
  This resource is used to manage the S3 Global Setting entity of PowerScale Array. PowerScale S3 Global Setting map to the PowerScale file system as base directory for Objects. We can Create, Update and Delete the S3 Global Setting using this resource. We can also import an existing S3 Global Setting from PowerScale array.
---

# powerscale_s3_global_settings (Resource)

This resource is used to manage the S3 Global Setting entity of PowerScale Array. PowerScale S3 Global Setting map to the PowerScale file system as base directory for Objects. We can Create, Update and Delete the S3 Global Setting using this resource. We can also import an existing S3 Global Setting from PowerScale array.


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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file it will modify S3 Global Settings on  the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale S3 global settings allows you to configure S3 global settings on PowerScale.
resource "powerscale_s3_global_settings" "s3_global_setting" {
  service    = true
  https_only = false
  http_port  = 9097
  https_port = 9098
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `http_port` (Number) Specifies the HTTP port.
- `https_only` (Boolean) Specifies if the service is HTTPS only.
- `https_port` (Number) Specifies the HTTPS port.
- `service` (Boolean) Specifies if the service is enabled.

Unless specified otherwise, all fields of this resource can be updated.

## Import

Import is supported using the following syntax:

```shell
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://mozilla.org/MPL/2.0/


# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The command is
# terraform import powerscale_s3_global_settings.s3_global_settings_example <any string>
terraform import powerscale_s3_global_settings.s3_global_settings_example ""

# after running this command, populate the name field and other required parameters in the config file to start managing this resource.
# Note: running "terraform show" after importing shows the current config/state of the resource. You can copy/paste that config to make it easier to manage the resource.
```