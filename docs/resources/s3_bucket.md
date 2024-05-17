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

title: "powerscale_s3_bucket resource"
linkTitle: "powerscale_s3_bucket"
page_title: "powerscale_s3_bucket Resource - terraform-provider-powerscale"
subcategory: ""
description: |-
  This resource is used to manage the S3 Bucket entity of PowerScale Array. PowerScale S3 Bucket map to the PowerScale file system as base directory for Objects. We can Create, Update and Delete the S3 Bucket using this resource. We can also import an existing S3 Bucket from PowerScale array.
---

# powerscale_s3_bucket (Resource)

This resource is used to manage the S3 Bucket entity of PowerScale Array. PowerScale S3 Bucket map to the PowerScale file system as base directory for Objects. We can Create, Update and Delete the S3 Bucket using this resource. We can also import an existing S3 Bucket from PowerScale array.


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
# After `terraform apply` of this example file it will create S3 Bucket on specified paths on the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale S3 Bucket enables access to file-based data that is stored on OneFS clusters as objects.

resource "powerscale_s3_bucket" "s3_bucket_example" {
  # Required attributes and update not supported
  name = "s3-bucket-example"
  path = "/ifs/s3_bucket_example"

  # Optional attributes and update not supported, 
  # Their default value shows as below if not provided during creation 
  # create_path = false
  # owner = "root"
  # zone = "System"

  # Optional attributes, can be updated
  #
  # By default acl is empty. To add an acl item, both grantee and permission are required.
  # Accepted values for permission are: READ, WRITE, READ_ACP, WRITE_ACP, FULL_CONTROL 
  # acl = [{
  #   grantee = {
  #     name = "root"
  #     type = "user"
  #   }
  #   permission = "FULL_CONTROL"
  # }]
  #
  # By default description is empty
  # description = ""
  #
  # Accepted values for object_acl_policy are: replace, deny.
  # The default value would be replace if unset.
  # object_acl_policy = "replace"
}

# After the execution of above resource block, a S3 Bucket would have been created on the PowerScale array.
# For more information, Please check the terraform state file.
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Bucket name.
- `path` (String) Path of bucket within /ifs.

### Optional

- `acl` (Attributes List) Specifies properties for an S3 Access Control Entry. (see [below for nested schema](#nestedatt--acl))
- `create_path` (Boolean) Create path if does not exist.
- `description` (String) Description for this S3 bucket.
- `object_acl_policy` (String) Set behavior of modifying object acls.
- `owner` (String) Specifies the name of the owner.
- `zone` (String) Zone Name.

### Read-Only

- `id` (String) Bucket ID.
- `zid` (Number) Zone ID.

<a id="nestedatt--acl"></a>
### Nested Schema for `acl`

Required:

- `grantee` (Attributes) Specifies the persona of the file group. (see [below for nested schema](#nestedatt--acl--grantee))
- `permission` (String) Specifies the S3 rights being allowed.

<a id="nestedatt--acl--grantee"></a>
### Nested Schema for `acl.grantee`

Required:

- `name` (String) Specifies the persona name, which must be combined with a type.
- `type` (String) Specifies the type of persona, which must be combined with a name.

Read-Only:

- `id` (String) Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.

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
# terraform import powerscale_s3_bucket.s3_bucket_example [<zoneID>]:<id>
# Example 1: <zoneID> is Optional, defaults to System:
terraform import powerscale_s3_bucket.s3_bucket_example example_s3_bucket_id
# Example 2:
terraform import powerscale_s3_bucket.s3_bucket_example zone_id:example_s3_bucket_id
# after running this command, populate the name field and other required parameters in the config file to start managing this resource.
# Note: running "terraform show" after importing shows the current config/state of the resource. You can copy/paste that config to make it easier to manage the resource.
```