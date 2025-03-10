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

package provider

import (
	"context"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &FilePoolPolicyResource{}
	_ resource.ResourceWithConfigure   = &FilePoolPolicyResource{}
	_ resource.ResourceWithImportState = &FilePoolPolicyResource{}
)

// NewFilePoolPolicyResource creates a new resource.
func NewFilePoolPolicyResource() resource.Resource {
	return &FilePoolPolicyResource{}
}

// FilePoolPolicyResource defines the resource implementation.
type FilePoolPolicyResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *FilePoolPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filepool_policy"
}

// Schema describes the resource arguments.
func (r *FilePoolPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the File Pool Policy entity of PowerScale Array. We can Create, Update and Delete the File Pool Policy using this resource. We can also import an existing File Pool Policy from PowerScale array. PowerScale File Pool Policy can identify logical groups of files and specify storage operations for these files.",
		Description:         "This resource is used to manage the File Pool Policy entity of PowerScale Array. We can Create, Update and Delete the File Pool Policy using this resource. We can also import an existing File Pool Policy from PowerScale array. PowerScale File Pool Policy can identify logical groups of files and specify storage operations for these files.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "A unique name for this policy. If the policy is default policy, its name should be \"Default policy\".",
				MarkdownDescription: "A unique name for this policy. If the policy is default policy, its name should be \"Default policy\".",
				Required:            true,
			},
			"is_default_policy": schema.BoolAttribute{
				Description: "Specifies if the policy is default policy. Default policy applies to all files not selected by higher-priority policies." +
					" Cannot be updated.",
				MarkdownDescription: "Specifies if the policy is default policy. Default policy applies to all files not selected by higher-priority policies." +
					" Cannot be updated.",
				Optional: true,
			},
			"file_matching_pattern": schema.SingleNestedAttribute{
				Description:         "Specifies the file matching rules for determining which files will be managed by this policy.",
				MarkdownDescription: "Specifies the file matching rules for determining which files will be managed by this policy.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"or_criteria": schema.ListNestedAttribute{
						Description:         "List of or_criteria file matching rules for this policy.",
						MarkdownDescription: "List of or_criteria file matching rules for this policy.",
						Required:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"and_criteria": schema.ListNestedAttribute{
									Description:         "List of and_criteria file matching rules for this policy.",
									MarkdownDescription: "List of and_criteria file matching rules for this policy.",
									Required:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"value": schema.StringAttribute{
												Description:         "The value to be compared against a file attribute.",
												MarkdownDescription: "The value to be compared against a file attribute.",
												Optional:            true,
											},
											"units": schema.StringAttribute{
												Description:         "Size unit value. One of 'B','KB','MB','GB','TB','PB','EB' (valid only with 'type' = 'size').",
												MarkdownDescription: "Size unit value. One of 'B','KB','MB','GB','TB','PB','EB' (valid only with 'type' = 'size').",
												Optional:            true,
												Validators:          []validator.String{stringvalidator.LengthBetween(1, 255)},
											},
											"type": schema.StringAttribute{
												Description:         "The file attribute to be compared to a given value.",
												MarkdownDescription: "The file attribute to be compared to a given value.",
												Required:            true,
											},
											"operator": schema.StringAttribute{
												Description:         "The comparison operator to use while comparing an attribute with its value.",
												MarkdownDescription: "The comparison operator to use while comparing an attribute with its value.",
												Optional:            true,
											},
											"field": schema.StringAttribute{
												Description:         "File attribute field name to be compared in a custom comparison (valid only with 'type' = 'custom_attribute').",
												MarkdownDescription: "File attribute field name to be compared in a custom comparison (valid only with 'type' = 'custom_attribute').",
												Optional:            true,
											},
											"use_relative_time": schema.BoolAttribute{
												Description:         "Whether time units refer to a calendar date and time (e.g., Jun 3, 2009) or a relative duration (e.g., 2 weeks) (valid only with 'type' in {accessed_time, birth_time, changed_time or metadata_changed_time}.",
												MarkdownDescription: "Whether time units refer to a calendar date and time (e.g., Jun 3, 2009) or a relative duration (e.g., 2 weeks) (valid only with 'type' in {accessed_time, birth_time, changed_time or metadata_changed_time}.",
												Optional:            true,
											},
											"case_sensitive": schema.BoolAttribute{
												Description:         "True to indicate case sensitivity when comparing file attributes (valid only with 'type' = 'name' or 'type' = 'path').",
												MarkdownDescription: "True to indicate case sensitivity when comparing file attributes (valid only with 'type' = 'name' or 'type' = 'path').",
												Optional:            true,
											},
											"begins_with": schema.BoolAttribute{
												Description:         "True to match the path exactly, False to match any subtree. (valid only with 'type' = 'path').",
												MarkdownDescription: "True to match the path exactly, False to match any subtree. (valid only with 'type' = 'path').",
												Optional:            true,
											},
											"attribute_exists": schema.BoolAttribute{
												Description:         "Indicates whether the existence of an attribute indicates a match (valid only with 'type' = 'custom_attribute').",
												MarkdownDescription: "Indicates whether the existence of an attribute indicates a match (valid only with 'type' = 'custom_attribute').",
												Optional:            true,
											},
										},
									},
									Validators: []validator.List{
										listvalidator.UniqueValues(),
										listvalidator.SizeBetween(1, 5),
									},
								},
							},
						},
						Validators: []validator.List{
							listvalidator.UniqueValues(),
							listvalidator.SizeBetween(1, 3),
						},
					},
				},
			},
			"actions": schema.ListNestedAttribute{
				Description:         "A list of actions to be taken for matching files.",
				MarkdownDescription: "A list of actions to be taken for matching files.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enable_packing_action": schema.BoolAttribute{
							Description:         "Action for enable_packing type. True to enable enable_packing action.",
							MarkdownDescription: "Action for enable_packing type. True to enable enable_packing action.",
							Optional:            true,
						},
						"enable_coalescer_action": schema.BoolAttribute{
							Description:         "Action for enable_coalescer type. Set write performance optimization. True to enable SmartCache action.",
							MarkdownDescription: "Action for enable_coalescer type. Set write performance optimization. True to enable SmartCache action.",
							Optional:            true,
						},
						"data_access_pattern_action": schema.StringAttribute{
							Description:         "Action for set_data_access_pattern type. Set data access pattern optimization. Acceptable values: random, concurrency, streaming.",
							MarkdownDescription: "Action for set_data_access_pattern type. Set data access pattern optimization. Acceptable values: random, concurrency, streaming.",
							Optional:            true,
						},
						"requested_protection_action": schema.StringAttribute{
							Description:         "Action for set_requested_protection type. Acceptable values: default, +1n, +2d:1n, +2n, +3d:1n, +3d:1n1d, +3n, +4d:1n, +4d:2n, +4n, 2x, 3x, 4x, 5x, 6x, 7x, 8x.",
							MarkdownDescription: "Action for set_requested_protection type. Acceptable values: default, +1n, +2d:1n, +2n, +3d:1n, +3d:1n1d, +3n, +4d:1n, +4d:2n, +4n, 2x, 3x, 4x, 5x, 6x, 7x, 8x.",
							Optional:            true,
						},
						"data_storage_policy_action": schema.SingleNestedAttribute{
							Description:         "Action for apply_data_storage_policy.",
							MarkdownDescription: "Action for apply_data_storage_policy.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"ssd_strategy": schema.StringAttribute{
									Description:         "Specifies the SSD strategy. Acceptable values: metadata, metadata-write, data, avoid.",
									MarkdownDescription: "Specifies the SSD strategy. Acceptable values: metadata, metadata-write, data, avoid.",
									Required:            true,
								},
								"storagepool": schema.StringAttribute{
									Description:         "Specifies the storage target.",
									MarkdownDescription: "Specifies the storage target.",
									Required:            true,
								},
							},
						},
						"snapshot_storage_policy_action": schema.SingleNestedAttribute{
							Description:         "Action for apply_snapshot_storage_policy.",
							MarkdownDescription: "Action for apply_snapshot_storage_policy.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"ssd_strategy": schema.StringAttribute{
									Description:         "Specifies the SSD strategy. Acceptable values: metadata, metadata-write, data, avoid.",
									MarkdownDescription: "Specifies the SSD strategy. Acceptable values: metadata, metadata-write, data, avoid.",
									Required:            true,
								},
								"storagepool": schema.StringAttribute{
									Description:         "Specifies the snapshot storage target.",
									MarkdownDescription: "Specifies the snapshot storage target.",
									Required:            true,
								},
							},
						},
						"cloudpool_policy_action": schema.SingleNestedAttribute{
							Description:         "Action for set_cloudpool_policy type.",
							MarkdownDescription: "Action for set_cloudpool_policy type.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"pool": schema.StringAttribute{
									Description:         "Specifies the cloudPool storage target.",
									MarkdownDescription: "Specifies the cloudPool storage target.",
									Required:            true,
								},
								"archive_snapshot_files": schema.BoolAttribute{
									Description:         "Specifies if files with snapshots should be archived.",
									MarkdownDescription: "Specifies if files with snapshots should be archived.",
									Computed:            true,
									Optional:            true,
								},
								"compression": schema.BoolAttribute{
									Description:         "Specifies if files should be compressed.",
									MarkdownDescription: "Specifies if files should be compressed.",
									Computed:            true,
									Optional:            true,
								},
								"encryption": schema.BoolAttribute{
									Description:         "Specifies if files should be encrypted.",
									MarkdownDescription: "Specifies if files should be encrypted.",
									Computed:            true,
									Optional:            true,
								},
								"data_retention": schema.Int64Attribute{
									Description:         "Specifies the minimum amount of time archived data will be retained in the cloud after deletion.",
									MarkdownDescription: "Specifies the minimum amount of time archived data will be retained in the cloud after deletion.",
									Computed:            true,
									Optional:            true,
								},
								"full_backup_retention": schema.Int64Attribute{
									Description:         "The minimum amount of time cloud files will be retained after the creation of a full NDMP backup. (Used with NDMP backups only.  Not applicable to SyncIQ.) ",
									MarkdownDescription: "The minimum amount of time cloud files will be retained after the creation of a full NDMP backup. (Used with NDMP backups only.  Not applicable to SyncIQ.) ",
									Computed:            true,
									Optional:            true,
								},
								"incremental_backup_retention": schema.Int64Attribute{
									Description:         "The minimum amount of time cloud files will be retained after the creation of a SyncIQ backup or an incremental NDMP backup. (Used with SyncIQ and NDMP backups.) ",
									MarkdownDescription: "The minimum amount of time cloud files will be retained after the creation of a SyncIQ backup or an incremental NDMP backup. (Used with SyncIQ and NDMP backups.) ",
									Computed:            true,
									Optional:            true,
								},
								"writeback_frequency": schema.Int64Attribute{
									Description:         "The minimum amount of time to wait before updating cloud data with local changes.",
									MarkdownDescription: "The minimum amount of time to wait before updating cloud data with local changes.",
									Computed:            true,
									Optional:            true,
								},
								"cache": schema.SingleNestedAttribute{
									Description:         "Specifies default cloudpool cache settings for new filepool policies.",
									MarkdownDescription: "Specifies default cloudpool cache settings for new filepool policies.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"expiration": schema.Int64Attribute{
											Description:         "Specifies cache expiration.",
											MarkdownDescription: "Specifies cache expiration.",
											Computed:            true,
											Optional:            true,
										},
										"read_ahead": schema.StringAttribute{
											Description:         "Specifies cache read ahead type. Acceptable values: partial, full.",
											MarkdownDescription: "Specifies cache read ahead type. Acceptable values: partial, full.",
											Computed:            true,
											Optional:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies cache type. Acceptable values: cached, no-cache.",
											MarkdownDescription: "Specifies cache type. Acceptable values: cached, no-cache.",
											Computed:            true,
											Optional:            true,
										},
									},
								},
							},
						},
						"action_type": schema.StringAttribute{
							Description:         "action_type Acceptable values: set_requested_protection, set_data_access_pattern, enable_coalescer, apply_data_storage_policy, apply_snapshot_storage_policy, set_cloudpool_policy, enable_packing.							",
							MarkdownDescription: "action_type Acceptable values: set_requested_protection, set_data_access_pattern, enable_coalescer, apply_data_storage_policy, apply_snapshot_storage_policy, set_cloudpool_policy, enable_packing.",
							Required:            true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.SizeBetween(1, 7),
				},
			},
			"apply_order": schema.Int64Attribute{
				Description:         "The order in which this policy should be applied (relative to other policies).",
				MarkdownDescription: "The order in which this policy should be applied (relative to other policies).",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "A description for this File Pool Policy.",
				MarkdownDescription: "A description for this File Pool Policy.",
				Optional:            true,
				Computed:            true,
			},
			"birth_cluster_id": schema.StringAttribute{
				Description:         "The guid assigned to the cluster on which the policy was created.",
				MarkdownDescription: "The guid assigned to the cluster on which the policy was created.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "A unique name for this File Pool Policy.",
				MarkdownDescription: "A unique name for this File Pool Policy.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "Indicates whether this policy is in a good state (\"OK\") or disabled (\"disabled\").",
				MarkdownDescription: "Indicates whether this policy is in a good state (\"OK\") or disabled (\"disabled\").",
				Computed:            true,
			},
			"state_details": schema.StringAttribute{
				Description:         "Gives further information to describe the state of this policy.",
				MarkdownDescription: "Gives further information to describe the state of this policy.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *FilePoolPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = pscaleClient
}

// Create allocates the resource.
func (r *FilePoolPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating File Pool Policy resource...")
	var plan models.FilePoolPolicyModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.IsPolicyParamInvalid(plan); err != nil {
		resp.Diagnostics.AddError("Error creating File Pool Policy", err.Error())
		return
	}

	if plan.IsDefaultPolicy.ValueBool() {
		if len(plan.Actions) > 0 {
			if err := helper.UpdateFilePoolDefaultPolicy(ctx, r.client, &plan); err != nil {
				resp.Diagnostics.AddError("Error creating Default File Pool Policy.", err.Error())
				return
			}
		}
		policyResponse, err := helper.GetFilePoolDefaultPolicy(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError("Error getting Default File Pool Policy", err.Error())
			return
		}
		if err = helper.UpdateFilePoolDefaultPolicyState(ctx, &plan, policyResponse); err != nil {
			resp.Diagnostics.AddError("Error reading Default File Pool Policy Resource",
				fmt.Sprintf("Error parsing Default File Pool Policy resource state: %s", err.Error()))
			return
		}
	} else {
		policyName := plan.Name.ValueString()
		if err := helper.CreateFilePoolPolicy(ctx, r.client, &plan); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error creating File Pool Policy - %s", policyName), err.Error())
			// if err, revert create
			if err = helper.DeleteFilePoolPolicy(ctx, r.client, policyName); err != nil {
				tflog.Error(ctx, fmt.Sprintf("Error deleting File Pool Policy when reverting creation - %s", err.Error()))
			}
			return
		}
		policyResponse, err := helper.GetFilePoolPolicy(ctx, r.client, plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error getting File Pool Policy after creation", err.Error())
			// if err, revert create
			if err = helper.DeleteFilePoolPolicy(ctx, r.client, policyName); err != nil {
				tflog.Error(ctx, fmt.Sprintf("Error deleting File Pool Policy when reverting creation - %s", err.Error()))
			}
			return
		}
		helper.UpdateFilePoolPolicyResourceState(ctx, &plan, policyResponse)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Create File Pool Policy resource")
}

// Read reads the resource state.
func (r *FilePoolPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading File Pool Policy resource")
	var state models.FilePoolPolicyModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.IsDefaultPolicy.ValueBool() {
		policyResponse, err := helper.GetFilePoolDefaultPolicy(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError("Error getting Default File Pool Policy", err.Error())
			return
		}
		if err = helper.UpdateFilePoolDefaultPolicyState(ctx, &state, policyResponse); err != nil {
			resp.Diagnostics.AddError("Error reading Default File Pool Policy Resource",
				fmt.Sprintf("Error parsing Default File Pool Policy resource state: %s", err.Error()))
			return
		}
	} else {
		policyResponse, err := helper.GetFilePoolPolicy(ctx, r.client, state.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error getting the File Pool Policy - %s", state.Name.ValueString()), err.Error())
			return
		}
		// parse response to state model
		helper.UpdateFilePoolPolicyResourceState(ctx, &state, policyResponse)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read File Pool Policy resource")
}

// Update updates the resource state.
func (r *FilePoolPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating File Pool Policy resource...")
	// Read Terraform plan into the model
	var plan models.FilePoolPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.FilePoolPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.IsDefaultPolicy.ValueBool() != state.IsDefaultPolicy.ValueBool() {
		resp.Diagnostics.AddError("Error updating File Pool Policy", "may not change is_default_policy")
		return
	}
	if err := helper.IsPolicyParamInvalid(plan); err != nil {
		resp.Diagnostics.AddError("Error updating File Pool Policy", err.Error())
		return
	}

	if plan.IsDefaultPolicy.ValueBool() {
		if err := helper.UpdateFilePoolDefaultPolicy(ctx, r.client, &plan); err != nil {
			resp.Diagnostics.AddError("Error updating Default File Pool Policy.", err.Error())
			return
		}
		policyResponse, err := helper.GetFilePoolDefaultPolicy(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError("Error getting Default File Pool Policy", err.Error())
			return
		}
		if err = helper.UpdateFilePoolDefaultPolicyState(ctx, &plan, policyResponse); err != nil {
			resp.Diagnostics.AddError("Error reading Default File Pool Policy Resource",
				fmt.Sprintf("Error parsing Default File Pool Policy resource state: %s", err.Error()))
			return
		}
	} else {
		if err := helper.UpdateFilePoolPolicy(ctx, r.client, &state, &plan); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error updating the File Pool Policy resource - %s", state.Name.ValueString()),
				err.Error(),
			)
			return
		}
		policyResponse, err := helper.GetFilePoolPolicy(ctx, r.client, plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error getting the File Pool Policy - %s", plan.Name.ValueString()),
				err.Error(),
			)
			return
		}
		helper.UpdateFilePoolPolicyResourceState(ctx, &plan, policyResponse)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Update File Pool Policy resource")
}

// Delete deletes the resource.
func (r *FilePoolPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting File Pool Policy resource")
	var state models.FilePoolPolicyModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if !state.IsDefaultPolicy.ValueBool() {
		if err := helper.DeleteFilePoolPolicy(ctx, r.client, state.Name.ValueString()); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error deleting the File Pool Policy - %s", state.Name.ValueString()), err.Error())
			return
		}
	}
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete File Pool Policy resource")
}

// ImportState imports the resource state.
func (r *FilePoolPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing File Pool Policy resource")
	var state models.FilePoolPolicyModel

	policyName := req.ID
	if policyName == "is_default_policy=true" {
		policyResponse, err := helper.GetFilePoolDefaultPolicy(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError("Error getting Default File Pool Policy", err.Error())
			return
		}
		if err = helper.UpdateFilePoolDefaultPolicyState(ctx, &state, policyResponse); err != nil {
			resp.Diagnostics.AddError("Error reading Default File Pool Policy Resource",
				fmt.Sprintf("Error parsing Default File Pool Policy resource state: %s", err.Error()))
			return
		}
	} else {
		policyResponse, err := helper.GetFilePoolPolicy(ctx, r.client, policyName)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error getting the File Pool Policy - %s", policyName), err.Error())
			return
		}
		if err := helper.UpdateFilePoolPolicyImportState(ctx, &state, policyResponse); err != nil {
			resp.Diagnostics.AddError("Error reading File Pool Policy Resource",
				fmt.Sprintf("Error parsing File Pool Policy resource state: %s", err.Error()))
			return
		}
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import File Pool Policy resource")
}
