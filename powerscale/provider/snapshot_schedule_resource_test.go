/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

package provider

import (
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var createSnapshotScheduleMocker *mockey.Mocker
var mapSnapshotScheduleMocker *mockey.Mocker

func TestAccSnapshotScheduleResource(t *testing.T) {
	var snapshotScheduleResourceName = "powerscale_snapshot_schedule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SnapshotScheduleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapshotScheduleResourceName, "path", "/ifs/tfacc_file_system_test"),
					resource.TestCheckResourceAttr(snapshotScheduleResourceName, "name", "tfacc_snap_schedule_test"),
					resource.TestCheckResourceAttr(snapshotScheduleResourceName, "retention_time", "3 Hour(s)"),
					resource.TestCheckResourceAttr(snapshotScheduleResourceName, "alias", "test_alias"),
				),
			},
			// ImportState testing
			{
				ResourceName: snapshotScheduleResourceName,
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, "tfacc_snap_schedule_test", states[0].Attributes["name"])
					assert.Equal(t, "/ifs/tfacc_file_system_test", states[0].Attributes["path"])
					return nil
				},
			},
			// Update name, path ,alias then do Read testing
			{
				Config: ProviderConfig + SnapshotScheduleUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapshotScheduleResourceName, "path", "/ifs/tfacc_test"),
					resource.TestCheckResourceAttr(snapshotScheduleResourceName, "retention_time", "4 Hour(s)"),
					resource.TestCheckResourceAttr(snapshotScheduleResourceName, "alias", "test_alias_updated"),
					resource.TestCheckResourceAttr(snapshotScheduleResourceName, "name", "tfacc_snap_schedule_update"),
				),
			},
			// Update to error state
			{
				Config:      ProviderConfig + SnapshotScheduleResourceConfigUpdateError,
				ExpectError: regexp.MustCompile(".*Could not update snapshot schedule*."),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSnapshotScheduleResourceCreateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if createSnapshotScheduleMocker != nil {
						createSnapshotScheduleMocker.UnPatch()
					}
					if mapSnapshotScheduleMocker != nil {
						mapSnapshotScheduleMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.CreateSnapshotSchedule).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error creating snapshot schedule*.`),
			},
		},
	})
}

func TestAccSnapshotScheduleResourceReadError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: ProviderConfig + SnapshotScheduleResourceConfig,
			},
			// Read Error
			{
				PreConfig: func() {
					if createSnapshotScheduleMocker != nil {
						createSnapshotScheduleMocker.UnPatch()
					}
					if mapSnapshotScheduleMocker != nil {
						mapSnapshotScheduleMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.GetSpecificSnapshotSchedule).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting the snapshot schedule*.`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSnapshotScheduleResourceUpdateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: ProviderConfig + SnapshotScheduleResourceConfig,
			},
			// Update Error
			{
				PreConfig: func() {
					if createSnapshotScheduleMocker != nil {
						createSnapshotScheduleMocker.UnPatch()
					}
					if mapSnapshotScheduleMocker != nil {
						mapSnapshotScheduleMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSnapshotSchedule).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error Updating the snapshot schedule*.`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSnapshotScheduleResourceGetError(t *testing.T) {
	createResp := powerscale.Createv1SnapshotScheduleResponse{
		Id: 12,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if createSnapshotScheduleMocker != nil {
						createSnapshotScheduleMocker.UnPatch()
					}
					if mapSnapshotScheduleMocker != nil {
						mapSnapshotScheduleMocker.UnPatch()
					}
					createSnapshotScheduleMocker = mockey.Mock(helper.CreateSnapshotSchedule).Return(&createResp, nil).Build()
					FunctionMocker = mockey.Mock(helper.GetSpecificSnapshotSchedule).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error creating snapshot schedule*.`),
			},
		},
	})
}

func TestAccSnapshotScheduleResourceMappingError(t *testing.T) {
	createResp := powerscale.Createv1SnapshotScheduleResponse{
		Id: 12,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if createSnapshotScheduleMocker != nil {
						createSnapshotScheduleMocker.UnPatch()
					}
					if mapSnapshotScheduleMocker != nil {
						mapSnapshotScheduleMocker.UnPatch()
					}
					createSnapshotScheduleMocker = mockey.Mock(helper.CreateSnapshotSchedule).Return(&createResp, nil).Build()
					FunctionMocker = mockey.Mock(helper.GetSpecificSnapshotSchedule).Return(nil, nil).Build()
					mapSnapshotScheduleMocker = mockey.Mock(helper.SnapshotScheduleMapper).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var FileSystemResourceConfigCommon3 = `
resource "powerscale_filesystem" "file_system_test2" {
	directory_path         = "/ifs"	
	name = "tfacc_test"	
	  recursive = true
	  overwrite = true
	  group = {
		id   = "GID:0"
		name = "wheel"
		type = "group"
	  }
	  owner = {
		  id   = "UID:0",
		 name = "root",
		 type = "user"
	   }
	}
`

var SnapshotScheduleResourceConfig = FileSystemResourceConfigCommon + `
resource "powerscale_snapshot_schedule" "test" {
  depends_on = [powerscale_filesystem.file_system_test]
  # Required name of snapshot schedule
  name = "tfacc_snap_schedule_test"
  alias = "test_alias"
  retention_time = "3 Hour(s)"
  path = "/ifs/tfacc_file_system_test"
}
`

var SnapshotScheduleUpdateResourceConfig = FileSystemResourceConfigCommon3 + `
resource "powerscale_snapshot_schedule" "test" {
depends_on = [powerscale_filesystem.file_system_test2]
	# Required name of snapshot schedule
	name = "tfacc_snap_schedule_update"
	alias = "test_alias_updated"
	path = "/ifs/tfacc_test"
	retention_time = "4 Hour(s)"
  }
  `

var SnapshotScheduleResourceConfigUpdateError = `
resource "powerscale_snapshot_schedule" "test" {
	# Required name of snapshot schedule
	name = "tfacc_snap_schedule_update"
	alias = "test_alias_updated"
	path = "/ifs/tfacc_invalid"
	retention_time = "4 Hour(s)"
  }
`
