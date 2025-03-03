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

package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"math"
	"strconv"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListSnapshotSchedules lists the snapshot schedules.
func ListSnapshotSchedules(ctx context.Context, client *client.Client, ssFilter *models.SnapshotScheduleFilter) ([]powerscale.V1SnapshotScheduleExtended, error) {
	listSsParam := client.PscaleOpenAPIClient.SnapshotApi.ListSnapshotv1SnapshotSchedules(ctx)
	if ssFilter != nil {
		if !ssFilter.Sort.IsNull() {
			listSsParam = listSsParam.Sort(ssFilter.Sort.ValueString())
		}
		if !ssFilter.Dir.IsNull() {
			listSsParam = listSsParam.Dir(ssFilter.Dir.ValueString())
		}
		if !ssFilter.Limit.IsNull() {
			listSsParam = listSsParam.Limit((ssFilter.Limit.ValueInt32()))
		}
	}
	snapshotSchedules, _, err := listSsParam.Execute()
	if err != nil {
		return nil, err
	}
	//pagination
	for snapshotSchedules.Resume != nil && (ssFilter == nil || ssFilter.Limit.IsNull()) {
		listSsParam = listSsParam.Resume(*snapshotSchedules.Resume)
		newresp, _, errAdd := listSsParam.Execute()
		if errAdd != nil {
			return snapshotSchedules.Schedules, err
		}
		snapshotSchedules.Resume = newresp.Resume
		snapshotSchedules.Schedules = append(snapshotSchedules.Schedules, newresp.Schedules...)
	}
	return snapshotSchedules.Schedules, nil
}

func ConvertTimeDurationToRetentionTime(time *int32) string {
	if time == nil {
		return "Never Expires"
	}
	timeValue := *time
	switch true {
	// Less then 60 use seconds
	case timeValue < int32(60):
		return fmt.Sprintf("%d %s", timeValue, "Second(s)")
	// Use Minutes
	case timeValue < int32(3600):
		return fmt.Sprintf("%d %s", timeValue/60, "Minute(s)")
	// Use hours
	case timeValue < int32(86400):
		return fmt.Sprintf("%d %s", timeValue/3600, "Hour(s)")
	// Use days
	case timeValue < int32(604800):
		return fmt.Sprintf("%d %s", timeValue/86400, "Day(s)")
	// Use weeks
	// This is not 29030400 1 year in seconds but it is what is returned by the API as equivalant to 1 year in seconds
	case timeValue < int32(29030400):
		return fmt.Sprintf("%d %s", timeValue/604800, "Week(s)")
	// Use years
	default:
		test := math.Ceil(float64(timeValue) / float64(31536000))
		rounded := strconv.FormatFloat(test, 'f', 0, 64)
		// Always round up for years
		return fmt.Sprintf("%s %s", rounded, "Year(s)")
	}
}

// ParseTimeStringToSeconds takes a string time value(in a specific format) and converts it to seconds.
func ParseTimeStringToSeconds(timeString string) (*int32, error) {
	if timeString == "Never Expires" {
		return nil, nil
	}
	parts := strings.Fields(timeString)

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid time format: %s", timeString)
	}

	value, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid time value: %s", parts[0])
	}

	unit := strings.ToLower(parts[1])
	var multiplier int64
	switch unit {
	case "second(s)":
		multiplier = 1
	case "minute(s)":
		multiplier = 60
	case "hour(s)":
		multiplier = 3600
	case "day(s)":
		multiplier = 86400
	case "week(s)":
		multiplier = 604800
	case "year(s)":
		multiplier = 31536000
	default:
		return nil, fmt.Errorf("unknown time unit: %s", unit)
	}

	result := value * multiplier

	if result > math.MaxInt32 || result < math.MinInt32 {
		return nil, fmt.Errorf("integer overflow when converting to int32")
	}

	seconds := int32(result)
	return &seconds, nil
}

// CreateSnapshotSchedule lists the snapshot schedules.
func CreateSnapshotSchedule(ctx context.Context, client *client.Client, plan *models.SnapshotScheduleResource) (*powerscale.Createv1SnapshotScheduleResponse, error) {

	ssBody := powerscale.V1SnapshotSchedule{
		Path:     plan.Path.ValueString(),
		Name:     plan.Name.ValueString(),
		Pattern:  plan.Pattern.ValueString(),
		Schedule: plan.Schedule.ValueString(),
	}
	if !plan.Alias.IsNull() {
		ssBody.SetAlias(plan.Alias.ValueString())
	}

	if !plan.RetentionTime.IsNull() {
		duration, err := ParseTimeStringToSeconds(plan.RetentionTime.ValueString())
		if err != nil {
			return nil, fmt.Errorf("error converting Retention time - %s", err)
		}

		ssBody.SetDuration(*duration)
	}

	createParam := client.PscaleOpenAPIClient.SnapshotApi.CreateSnapshotv1SnapshotSchedule(ctx)
	createParam = createParam.V1SnapshotSchedule(ssBody)
	result, _, err := createParam.Execute()
	return result, err
}

// SnapshotScheduleMapper Does the mapping from response to model.
func SnapshotScheduleMapper(ctx context.Context, snapshotSchedule powerscale.V1SnapshotScheduleExtendedExtendedExtended, model *models.SnapshotScheduleResource) error {
	err := CopyFields(ctx, &snapshotSchedule, model)
	if err != nil {
		return err
	}
	// Import case convert from duration to retention time
	if model.RetentionTime.IsNull() {
		model.RetentionTime = types.StringValue(ConvertTimeDurationToRetentionTime(snapshotSchedule.Duration))
	}
	model.ID = types.StringValue(fmt.Sprint(*snapshotSchedule.Id))
	model.NextRun = types.Int64Value(int64(*snapshotSchedule.NextRun))
	return nil
}

// UpdateSnapshotSchedule updates the Snapshot Schedule.
func UpdateSnapshotSchedule(ctx context.Context, client *client.Client, plan *models.SnapshotScheduleResource, state *models.SnapshotScheduleResource) error {
	ss := *powerscale.NewV1SnapshotScheduleExtendedExtended()

	if plan.Name.ValueString() != state.Name.ValueString() {
		ss.SetName(plan.Name.ValueString())
	}
	if plan.Path.ValueString() != state.Path.ValueString() {
		ss.SetPath(plan.Path.ValueString())
	}
	if plan.Pattern.ValueString() != state.Pattern.ValueString() {
		ss.SetPattern(plan.Pattern.ValueString())
	}
	if plan.Schedule.ValueString() != state.Schedule.ValueString() {
		ss.SetSchedule(plan.Schedule.ValueString())
	}

	if !plan.Alias.IsNull() && !plan.Alias.IsUnknown() && (state.Alias.IsNull() || plan.Alias.ValueString() != state.Alias.ValueString()) {
		ss.SetAlias(plan.Alias.ValueString())
	}
	if !plan.RetentionTime.IsNull() {
		duration, err := ParseTimeStringToSeconds(plan.RetentionTime.ValueString())
		if err != nil {
			return fmt.Errorf("error converting Retention time - %s", err)
		}
		ss.SetDuration(*duration)
	}
	updateReq := client.PscaleOpenAPIClient.SnapshotApi.UpdateSnapshotv1SnapshotSchedule(ctx, state.ID.ValueString())
	updateReq = updateReq.V1SnapshotSchedule(ss)
	_, err := updateReq.Execute()

	return err
}

// GetSpecificSnapshotSchedule returns a Snapshot Schedule.
func GetSpecificSnapshotSchedule(ctx context.Context, client *client.Client, id string) (powerscale.V1SnapshotScheduleExtendedExtendedExtended, error) {
	ss := powerscale.V1SnapshotScheduleExtendedExtendedExtended{}
	result, _, err := client.PscaleOpenAPIClient.SnapshotApi.GetSnapshotv1SnapshotSchedule(ctx, id).Execute()
	if result != nil && len(result.Schedules) > 0 {
		ss = result.Schedules[0]
	}
	return ss, err
}
