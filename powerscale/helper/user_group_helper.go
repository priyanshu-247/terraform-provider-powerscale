/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"strconv"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateUserGroupDataSourceState updates datasource state.
func UpdateUserGroupDataSourceState(model *models.UserGroupModel, groupResponse powerscale.V1AuthGroupExtended,
	groupMembers []powerscale.V1AuthAccessAccessItemFileGroup, roles []powerscale.V1AuthRoleExtended) {
	model.Dn = types.StringValue(groupResponse.Dn)
	model.Domain = types.StringValue(groupResponse.Domain)
	model.DNSDomain = types.StringValue(groupResponse.DnsDomain)
	model.ID = types.StringValue(groupResponse.Id)
	model.Name = types.StringValue(groupResponse.Name)
	model.Provider = types.StringValue(groupResponse.Provider)
	model.SamAccountName = types.StringValue(groupResponse.SamAccountName)
	model.Type = types.StringValue(groupResponse.Type)
	if groupResponse.Gid.Id != nil {
		model.GID = types.StringValue(*groupResponse.Gid.Id)
	}
	if groupResponse.Sid.Id != nil {
		model.SID = types.StringValue(*groupResponse.Sid.Id)
	}
	model.GeneratedGID = types.BoolValue(groupResponse.GeneratedGid)

	//parse roles
	var roleAttrs []attr.Value
	for _, r := range roles {
		for _, m := range r.Members {
			if *m.Id == *groupResponse.Gid.Id {
				roleAttrs = append(roleAttrs, types.StringValue(r.Name))
			}
		}
	}
	model.Roles, _ = types.ListValue(types.StringType, roleAttrs)

	// parse group members
	var members []models.V1AuthAccessAccessItemFileGroup
	for _, m := range groupMembers {
		members = append(members, models.V1AuthAccessAccessItemFileGroup{
			Name: types.StringValue(*m.Name),
			ID:   types.StringValue(*m.Id),
			Type: types.StringValue(*m.Type),
		})
	}
	model.Members = members
}

// GetAllGroupMembers returns all group members
func GetAllGroupMembers(ctx context.Context, client *client.Client, groupName string) (members []powerscale.V1AuthAccessAccessItemFileGroup, err error) {
	memberParams := client.PscaleOpenAPIClient.AuthGroupsApi.ListAuthGroupsv1GroupMembers(ctx, groupName)
	result, _, err := memberParams.Execute()
	if err != nil {
		errStr := constants.ReadUserGroupMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user group members: %s", message)
	}

	for {
		members = append(members, result.Members...)
		if result.Resume == nil || *result.Resume == "" {
			break
		}
		memberParams = client.PscaleOpenAPIClient.AuthGroupsApi.ListAuthGroupsv1GroupMembers(ctx, groupName).Resume(*result.Resume)
		if result, _, err = memberParams.Execute(); err != nil {
			errStr := constants.ReadUserGroupMemberErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return nil, fmt.Errorf("error getting user group members with resume: %s", message)
		}
	}

	return
}

// GetUserGroupsWithFilter returns all filtered groups.
func GetUserGroupsWithFilter(ctx context.Context, client *client.Client, filter *models.UserGroupFilterType) (groups []powerscale.V1AuthGroupExtended, err error) {
	groupParams := client.PscaleOpenAPIClient.AuthApi.ListAuthv1AuthGroups(ctx)
	if filter != nil {
		if !filter.NamePrefix.IsNull() {
			groupParams = groupParams.Filter(filter.NamePrefix.ValueString())
		}
		if !filter.Domain.IsNull() {
			groupParams = groupParams.Domain(filter.Domain.ValueString())
		}
		if !filter.Zone.IsNull() {
			groupParams = groupParams.Zone(filter.Zone.ValueString())
		}
		if !filter.Provider.IsNull() {
			groupParams = groupParams.Provider(filter.Provider.ValueString())
		}
		if !filter.Cached.IsNull() {
			groupParams = groupParams.Cached(filter.Cached.ValueBool())
		}
		if !filter.ResolveNames.IsNull() {
			groupParams = groupParams.ResolveNames(filter.ResolveNames.ValueBool())
		}
	}

	result, _, err := groupParams.Execute()
	if err != nil {
		errStr := constants.ReadUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user groups: %s", message)
	}

	for {
		groups = append(groups, result.Groups...)
		if result.Resume == nil || *result.Resume == "" {
			break
		}
		groupParams = client.PscaleOpenAPIClient.AuthApi.ListAuthv1AuthGroups(ctx).Resume(*result.Resume)
		if result, _, err = groupParams.Execute(); err != nil {
			errStr := constants.ReadUserGroupErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return nil, fmt.Errorf("error getting user groups with resume: %s", message)
		}
	}

	if filter != nil && len(filter.Names) > 0 {
		var validUserGroups []string
		var filteredUserGroups []powerscale.V1AuthGroupExtended

		for _, group := range groups {
			for _, name := range filter.Names {
				if (!name.Name.IsNull() && group.Name == name.Name.ValueString()) ||
					(!name.GID.IsNull() && fmt.Sprintf("GID:%d", name.GID.ValueInt64()) == *group.Gid.Id) {
					filteredUserGroups = append(filteredUserGroups, group)
					validUserGroups = append(validUserGroups, fmt.Sprintf("Name: %s, GID: %s", group.Name, *group.Gid.Id))
					continue
				}
			}
		}

		if len(filteredUserGroups) != len(filter.Names) {
			return nil, fmt.Errorf(
				"error one or more of the filtered user group names is not a valid powerscale user group. Valid user groups: [%v], filtered list: [%v]",
				validUserGroups, filter.Names)
		}

		groups = filteredUserGroups
	}
	return
}

// GetUserGroup Returns the User Group by user group name.
func GetUserGroup(ctx context.Context, client *client.Client, groupName string) (*powerscale.V1AuthGroupsExtended, error) {
	authID := fmt.Sprintf("GROUP:%s", groupName)
	getParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv1AuthGroup(ctx, groupName)
	result, _, err := getParam.Execute()
	if err != nil {
		errStr := constants.ReadUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user group - %s : %s", authID, message)
	}
	if len(result.Groups) < 1 {
		message := constants.ReadUserGroupErrorMsg + "with error: "
		return nil, fmt.Errorf("got empty user group - %s : %s", authID, message)
	}

	return result, err
}

// UpdateUserGroupResourceState updates resource state.
func UpdateUserGroupResourceState(model *models.UserGroupReourceModel, group powerscale.V1AuthGroupExtended,
	groupMembers []powerscale.V1AuthAccessAccessItemFileGroup, roles []powerscale.V1AuthRoleExtended) {
	model.Dn = types.StringValue(group.Dn)
	model.Domain = types.StringValue(group.Domain)
	model.DNSDomain = types.StringValue(group.DnsDomain)
	model.ID = types.StringValue(group.Id)
	model.Name = types.StringValue(group.Name)
	model.Provider = types.StringValue(group.Provider)
	model.SamAccountName = types.StringValue(group.SamAccountName)
	model.Type = types.StringValue(group.Type)
	model.GeneratedGID = types.BoolValue(group.GeneratedGid)

	if group.Sid.Id != nil {
		model.SID = types.StringValue(*group.Sid.Id)
	}
	if group.Gid.Id != nil && strings.HasPrefix(*group.Gid.Id, "GID:") {
		gidInt, _ := strconv.Atoi(strings.Trim(*group.Gid.Id, "GID:"))
		model.GID = types.Int64Value(int64(gidInt))
	}

	if roles != nil {
		var roleAttrs []attr.Value
		for _, r := range roles {
			for _, m := range r.Members {
				if *m.Id == *group.Gid.Id {
					roleAttrs = append(roleAttrs, types.StringValue(r.Name))
				}
			}
		}

		model.Roles, _ = types.ListValue(types.StringType, roleAttrs)
	}
	if groupMembers != nil {
		var users []attr.Value
		for _, m := range groupMembers {
			users = append(users, types.StringValue(*m.Name))
		}
		model.Users, _ = types.ListValue(types.StringType, users)
	}

}

// CreateUserGroup Creates an User Group.
func CreateUserGroup(ctx context.Context, client *client.Client, plan *models.UserGroupReourceModel) error {

	createParam := client.PscaleOpenAPIClient.AuthApi.CreateAuthv1AuthGroup(ctx)
	if !plan.QueryForce.IsNull() {
		createParam = createParam.Force(plan.QueryForce.ValueBool())
	}
	if !plan.QueryZone.IsNull() {
		createParam = createParam.Zone(plan.QueryZone.ValueString())
	}
	if !plan.QueryProvider.IsNull() {
		createParam = createParam.Provider(plan.QueryProvider.ValueString())
	}

	body := &powerscale.V1AuthGroup{Name: plan.Name.ValueString()}
	if !plan.GID.IsNull() && plan.GID.ValueInt64() > 0 {
		body.Gid = plan.GID.ValueInt64Pointer()
	}
	if !plan.SID.IsNull() && plan.SID.ValueString() != "" {
		body.Sid = plan.SID.ValueStringPointer()
	}

	for _, m := range plan.Users.Elements() {
		name := strings.Trim(m.String(), "\"")
		body.Members = append(body.Members, powerscale.V1AuthAccessAccessItemFileGroup{Name: &name})
	}

	createParam = createParam.V1AuthGroup(*body)
	if _, _, err := createParam.Execute(); err != nil {
		errStr := constants.CreateUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error creating user group: %s", message)
	}

	return nil
}

// UpdateUserGroup Updates an User Group GID.
func UpdateUserGroup(ctx context.Context, client *client.Client, state *models.UserGroupReourceModel, plan *models.UserGroupReourceModel) error {
	authID := fmt.Sprintf("GROUP:%s", plan.Name.ValueString())
	updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv1AuthGroup(ctx, authID)

	if !plan.QueryForce.IsNull() {
		updateParam = updateParam.Force(plan.QueryForce.ValueBool())
	}
	if !plan.QueryZone.IsNull() {
		updateParam = updateParam.Zone(plan.QueryZone.ValueString())
	}
	if !plan.QueryProvider.IsNull() {
		updateParam = updateParam.Provider(plan.QueryProvider.ValueString())
	}

	body := &powerscale.V1AuthGroupExtendedExtended{}
	if !state.Name.Equal(plan.Name) {
		return fmt.Errorf("may not change user group's name")
	}
	if !state.GID.Equal(plan.GID) && plan.GID.ValueInt64() > 0 {
		if !plan.QueryForce.ValueBool() {
			return fmt.Errorf("may not change user group's gid without using the force option")
		}
		body.Gid = plan.GID.ValueInt64Pointer()
	}

	updateParam = updateParam.V1AuthGroup(*body)
	if _, err := updateParam.Execute(); err != nil {
		errStr := constants.UpdateUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error updating user group - %s : %s", authID, message)
	}

	return nil
}

// UpdateUserGroupRoles Updates an User Group roles.
func UpdateUserGroupRoles(ctx context.Context, client *client.Client, state *models.UserGroupReourceModel, plan *models.UserGroupReourceModel) (diags diag.Diagnostics) {

	// get roles list changes
	toAdd, toRemove := GetElementsChanges(state.Roles.Elements(), plan.Roles.Elements())

	// if gid changed, should remove all roles firstly, then assign all roles.
	if !plan.GID.Equal(state.GID) {
		toAdd = plan.Roles.Elements()
		toRemove = state.Roles.Elements()
	}

	// remove roles from user group
	for _, i := range toRemove {
		roleID := strings.Trim(i.String(), "\"")
		if err := RemoveUserGroupRole(ctx, client, roleID, state.GID.ValueInt64()); err != nil {
			diags.AddError(fmt.Sprintf("Error remove User Group from Role - %s", roleID), err.Error())
		}
	}

	// assign roles to user group
	for _, i := range toAdd {
		roleID := strings.Trim(i.String(), "\"")
		if err := AddUserGroupRole(ctx, client, roleID, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error assign User Group to Role - %s", roleID), err.Error())
		}
	}

	return
}

// AddUserGroupRole Assigns role to user group.
func AddUserGroupRole(ctx context.Context, client *client.Client, roleID, userGroupName string) error {
	authID := userGroupName
	roleParam := client.PscaleOpenAPIClient.AuthRolesApi.CreateAuthRolesv1RoleMember(ctx, roleID).
		V1RoleMember(powerscale.V1AuthAccessAccessItemFileGroup{Name: &authID})
	if _, _, err := roleParam.Execute(); err != nil {
		errStr := constants.AddRoleMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error assign user group - %s to role - %s: %s", authID, roleID, message)
	}
	return nil
}

// RemoveUserGroupRole Removes role from user group.
func RemoveUserGroupRole(ctx context.Context, client *client.Client, roleID string, gid int64) error {
	authID := fmt.Sprintf("GID:%d", gid)
	roleParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv1RolesRoleMember(ctx, authID, roleID)
	if _, err := roleParam.Execute(); err != nil {
		errStr := constants.DeleteRoleMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error remove user group - %s from role - %s: %s", authID, roleID, message)
	}
	return nil
}

// UpdateUserGroupMembers Updates an User Group members.
func UpdateUserGroupMembers(ctx context.Context, client *client.Client, state *models.UserGroupReourceModel, plan *models.UserGroupReourceModel) (diags diag.Diagnostics) {

	// get members list changes
	toAdd, toRemove := GetElementsChanges(state.Users.Elements(), plan.Users.Elements())

	// remove users from user group by memberAuthID
	for _, i := range toRemove {
		memberAuthID := fmt.Sprintf("USER:%s", strings.Trim(i.String(), "\""))
		if err := RemoveUserGroupMember(ctx, client, memberAuthID, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error remove User - %s from User Group.", memberAuthID), err.Error())
		}
	}

	// add users to user group by memberID
	for _, i := range toAdd {
		memberID := strings.Trim(i.String(), "\"")
		if err := AddUserGroupMember(ctx, client, memberID, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error add User - %s to User Group.", memberID), err.Error())
		}
	}

	return
}

// RemoveUserGroupMember Removes member from user group by memberAuthID, like GROUP:groupName and USER:userName.
func RemoveUserGroupMember(ctx context.Context, client *client.Client, memberAuthID, userGroupName string) error {
	memberParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv1GroupsGroupMember(ctx, memberAuthID, userGroupName)
	if _, err := memberParam.Execute(); err != nil {
		errStr := constants.DeleteUserGroupMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error remove member - %s from user group - %s: %s", memberAuthID, userGroupName, message)
	}
	return nil
}

// AddUserGroupMember Adds member to user group by memberID.
func AddUserGroupMember(ctx context.Context, client *client.Client, memberID, userGroupName string) error {
	authID := memberID
	memberParam := client.PscaleOpenAPIClient.AuthGroupsApi.CreateAuthGroupsv1GroupMember(ctx, userGroupName).
		V1GroupMember(powerscale.V1AuthAccessAccessItemFileGroup{Name: &authID})
	if _, _, err := memberParam.Execute(); err != nil {
		errStr := constants.AddUserGroupMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error add member - %s to user group - %s: %s", memberID, userGroupName, message)
	}
	return nil
}

// DeleteUserGroup Deletes an User Group.
func DeleteUserGroup(ctx context.Context, client *client.Client, groupName string) error {
	authID := fmt.Sprintf("GROUP:%s", groupName)
	deleteParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv1AuthGroup(ctx, authID)
	if _, err := deleteParam.Execute(); err != nil {
		errStr := constants.DeleteUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error deleting user group - %s : %s", authID, message)
	}

	return nil
}