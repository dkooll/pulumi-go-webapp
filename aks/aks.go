package aks

import (
	"fmt"

	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/authorization"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/containerservice"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Aks struct {
	pulumi.ResourceState

	ResourceId pulumi.StringOutput `pulumi:"id"`
}

type AksArgs struct {
	EnableRBAC        pulumi.BoolInput
	ResourceName      pulumi.StringInput
	KubernetesVersion pulumi.StringInput
	NodeResourceGroup pulumi.StringInput
	ResourceGroup     pulumi.StringInput
	NodeCount         pulumi.IntInput
	Location          pulumi.StringInput
	Scope             pulumi.StringInput
}

func CreateAks(ctx *pulumi.Context, name string, args *AksArgs, opts ...pulumi.ResourceOption) (*Aks, error) {
	aks := &Aks{}

	err := ctx.RegisterComponentResource("resource:index:Aks", name, aks, opts...)
	if err != nil {
		return nil, err
	}

	rgs, err := core.NewResourceGroup(ctx, "resourceGroups", &core.ResourceGroupArgs{
		Name:     args.ResourceGroup,
		Location: args.Location,
	}, pulumi.Parent(aks))
	if err != nil {
		return nil, fmt.Errorf("error creating resource group: %v", err)
	}

	k8s, err := containerservice.NewKubernetesCluster(ctx, "akss", &containerservice.KubernetesClusterArgs{
		Name:                          args.ResourceName,
		Location:                      rgs.Location,
		ResourceGroupName:             rgs.Name,
		NodeResourceGroup:             args.NodeResourceGroup,
		RoleBasedAccessControlEnabled: args.EnableRBAC,
		DnsPrefix:                     pulumi.String("akss"),
		KubernetesVersion:             args.KubernetesVersion,
		SkuTier:                       pulumi.String("Free"),
		DefaultNodePool: &containerservice.KubernetesClusterDefaultNodePoolArgs{
			NodeCount:          args.NodeCount,
			Name:               pulumi.String("default"),
			EnableAutoScaling:  pulumi.Bool(false),
			EnableNodePublicIp: pulumi.Bool(true),
			OsSku:              pulumi.String("Ubuntu"),
			Type:               pulumi.String("VirtualMachineScaleSets"),
			VmSize:             pulumi.String("Standard_D2_v2"),
		},
		Identity: &containerservice.KubernetesClusterIdentityArgs{
			Type: pulumi.String("SystemAssigned"),
		},
	}, pulumi.Parent(rgs))
	if err != nil {
		return nil, fmt.Errorf("error creating cluster: %v", err)
	}

	//acr pull rights aks mi
	authorization.NewAssignment(ctx, "ra", &authorization.AssignmentArgs{
		PrincipalId: k8s.KubeletIdentity.ApplyT(func(kubeletIdentity containerservice.KubernetesClusterKubeletIdentity) (string, error) {
			return *kubeletIdentity.ObjectId, nil
		}).(pulumi.StringOutput),
		RoleDefinitionName:           pulumi.String("AcrPull"),
		Scope:                        args.Scope,
		SkipServicePrincipalAadCheck: pulumi.Bool(true),
	}, pulumi.Parent(k8s))
	if err != nil {
		return nil, fmt.Errorf("error creating role assignment: %v", err)
	}

	aks.ResourceId = k8s.ID().ToStringOutput()
	ctx.Export("id", aks.ResourceId)
	return aks, nil
}
