// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package iam

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	iam_sdkv1 "github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceAccessKeys,
			TypeName: "aws_iam_access_keys",
		},
		{
			Factory:  DataSourceAccountAlias,
			TypeName: "aws_iam_account_alias",
		},
		{
			Factory:  DataSourceGroup,
			TypeName: "aws_iam_group",
		},
		{
			Factory:  DataSourceInstanceProfile,
			TypeName: "aws_iam_instance_profile",
		},
		{
			Factory:  DataSourceInstanceProfiles,
			TypeName: "aws_iam_instance_profiles",
		},
		{
			Factory:  DataSourceOpenIDConnectProvider,
			TypeName: "aws_iam_openid_connect_provider",
		},
		{
			Factory:  DataSourcePolicy,
			TypeName: "aws_iam_policy",
		},
		{
			Factory:  DataSourcePolicyDocument,
			TypeName: "aws_iam_policy_document",
		},
		{
			Factory:  DataSourcePrincipalPolicySimulation,
			TypeName: "aws_iam_principal_policy_simulation",
		},
		{
			Factory:  DataSourceRole,
			TypeName: "aws_iam_role",
		},
		{
			Factory:  DataSourceRoles,
			TypeName: "aws_iam_roles",
		},
		{
			Factory:  DataSourceSAMLProvider,
			TypeName: "aws_iam_saml_provider",
		},
		{
			Factory:  DataSourceServerCertificate,
			TypeName: "aws_iam_server_certificate",
		},
		{
			Factory:  DataSourceSessionContext,
			TypeName: "aws_iam_session_context",
		},
		{
			Factory:  DataSourceUser,
			TypeName: "aws_iam_user",
		},
		{
			Factory:  DataSourceUserSSHKey,
			TypeName: "aws_iam_user_ssh_key",
		},
		{
			Factory:  DataSourceUsers,
			TypeName: "aws_iam_users",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceAccessKey,
			TypeName: "aws_iam_access_key",
		},
		{
			Factory:  ResourceAccountAlias,
			TypeName: "aws_iam_account_alias",
		},
		{
			Factory:  ResourceAccountPasswordPolicy,
			TypeName: "aws_iam_account_password_policy",
		},
		{
			Factory:  ResourceGroup,
			TypeName: "aws_iam_group",
		},
		{
			Factory:  ResourceGroupMembership,
			TypeName: "aws_iam_group_membership",
		},
		{
			Factory:  ResourceGroupPolicy,
			TypeName: "aws_iam_group_policy",
		},
		{
			Factory:  ResourceGroupPolicyAttachment,
			TypeName: "aws_iam_group_policy_attachment",
		},
		{
			Factory:  ResourceInstanceProfile,
			TypeName: "aws_iam_instance_profile",
			Name:     "Instance Profile",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceOpenIDConnectProvider,
			TypeName: "aws_iam_openid_connect_provider",
			Name:     "OIDC Provider",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourcePolicy,
			TypeName: "aws_iam_policy",
			Name:     "Policy",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourcePolicyAttachment,
			TypeName: "aws_iam_policy_attachment",
		},
		{
			Factory:  ResourceRole,
			TypeName: "aws_iam_role",
			Name:     "Role",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceRolePolicy,
			TypeName: "aws_iam_role_policy",
		},
		{
			Factory:  ResourceRolePolicyAttachment,
			TypeName: "aws_iam_role_policy_attachment",
		},
		{
			Factory:  ResourceSAMLProvider,
			TypeName: "aws_iam_saml_provider",
			Name:     "SAML Provider",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceServerCertificate,
			TypeName: "aws_iam_server_certificate",
			Name:     "Server Certificate",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceServiceLinkedRole,
			TypeName: "aws_iam_service_linked_role",
			Name:     "Service Linked Role",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceServiceSpecificCredential,
			TypeName: "aws_iam_service_specific_credential",
		},
		{
			Factory:  ResourceSigningCertificate,
			TypeName: "aws_iam_signing_certificate",
		},
		{
			Factory:  ResourceUser,
			TypeName: "aws_iam_user",
			Name:     "User",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceUserGroupMembership,
			TypeName: "aws_iam_user_group_membership",
		},
		{
			Factory:  ResourceUserLoginProfile,
			TypeName: "aws_iam_user_login_profile",
		},
		{
			Factory:  ResourceUserPolicy,
			TypeName: "aws_iam_user_policy",
		},
		{
			Factory:  ResourceUserPolicyAttachment,
			TypeName: "aws_iam_user_policy_attachment",
		},
		{
			Factory:  ResourceUserSSHKey,
			TypeName: "aws_iam_user_ssh_key",
		},
		{
			Factory:  ResourceVirtualMFADevice,
			TypeName: "aws_iam_virtual_mfa_device",
			Name:     "Virtual MFA Device",
			Tags:     &types.ServicePackageResourceTags{},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.IAM
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*iam_sdkv1.IAM, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return iam_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

var ServicePackage = &servicePackage{}
