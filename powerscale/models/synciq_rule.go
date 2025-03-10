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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource models

// SyncIQRuleDataSource defines the data source implementation.
type SyncIQRuleDataSource struct {
	// ID is the unique identifier of the data source.
	ID types.String `tfsdk:"id"`
	// Rules is a list of SyncIQ rules.
	Rules []SyncIQRuleModel `tfsdk:"rules"`
}

// SyncIQRuleModel defines the model of a SyncIQ rule.
type SyncIQRuleModel struct {
	// Type is the type of the rule.
	Type types.String `tfsdk:"type"`
	// Description is the description of the rule.
	Description types.String `tfsdk:"description"`
	// Enabled indicates if the rule is enabled.
	Enabled types.Bool `tfsdk:"enabled"`
	// ID is the unique identifier of the rule.
	ID types.String `tfsdk:"id"`
	// Limit is the limit of the rule.
	Limit types.Int64 `tfsdk:"limit"`
	// Schedule is the schedule of the rule.
	Schedule Schedule `tfsdk:"schedule"`
}

// Schedule defines the schedule of a rule.
type Schedule struct {
	// Tuesday indicates if the rule is in effect on Tuesday.
	Tuesday types.Bool `tfsdk:"tuesday"`
	// Monday indicates if the rule is in effect on Monday.
	Monday types.Bool `tfsdk:"monday"`
	// Wednesday indicates if the rule is in effect on Wednesday.
	Wednesday types.Bool `tfsdk:"wednesday"`
	// Saturday indicates if the rule is in effect on Saturday.
	Saturday types.Bool `tfsdk:"saturday"`
	// Friday indicates if the rule is in effect on Friday.
	Friday types.Bool `tfsdk:"friday"`
	// End is the end time of the schedule.
	End types.String `tfsdk:"end"`
	// Sunday indicates if the rule is in effect on Sunday.
	Sunday types.Bool `tfsdk:"sunday"`
	// Begin is the beginning time of the schedule.
	Begin types.String `tfsdk:"begin"`
	// Thursday indicates if the rule is in effect on Thursday.
	Thursday types.Bool `tfsdk:"thursday"`
}

// Resource models

// SyncIQRulesResource defines the model of a SyncIQ rules resource.
type SyncIQRulesResource struct {
	ID             types.String `tfsdk:"id"`
	BandWidthRules types.List   `tfsdk:"bandwidth_rules"`
	FileCountRules types.List   `tfsdk:"file_count_rules"`
	CPURules       types.List   `tfsdk:"cpu_rules"`
	WorkerRules    types.List   `tfsdk:"worker_rules"`
}

// SyncIQRulesResourceRequest defines a model of with a list of rules for each stack
type SyncIQRulesResourceRequest struct {
	FileCountRules []SyncIQRuleResource
	CPURules       []SyncIQRuleResource
	WorkerRules    []SyncIQRuleResource
	BandWidthRules []SyncIQRuleResource
}

// SyncIQRuleResource defines the model of a SyncIQ rule.
type SyncIQRuleResource struct {
	// Description is the description of the rule.
	Description types.String `tfsdk:"description"`
	// Enabled indicates if the rule is enabled.
	Enabled types.Bool `tfsdk:"enabled"`
	// ID is the unique identifier of the rule.
	ID types.String `tfsdk:"id"`
	// Limit is the limit of the rule.
	Limit types.Int32 `tfsdk:"limit"`
	// Schedule is the schedule of the rule.
	Schedule types.Object `tfsdk:"schedule"`
}

// SyncIQRuleResourceSchedule defines the schedule of a SyncIQ rule.
type SyncIQRuleResourceSchedule struct {
	End        *string  `tfsdk:"end"`
	Begin      *string  `tfsdk:"begin"`
	DaysOfWeek []string `tfsdk:"days_of_week"`
}
