package aks

import (
	"fmt"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/containerservice"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Aks struct {
	pulumi.ResourceState

	resourceId pulumi.StringOutput
}

type AksArgs struct {
	EnableRBAC        pulumi.BoolInput
	ResourceName      pulumi.StringInput
	KubernetesVersion pulumi.StringInput
	NodeResourceGroup pulumi.StringInput
	ResourceGroup     pulumi.StringInput
	NodeCount         pulumi.IntInput
	Location          pulumi.StringInput
	resourceId        pulumi.StringOutput
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

	k8s, err := containerservice.NewManagedCluster(ctx, "akss", &containerservice.ManagedClusterArgs{
		Location: rgs.Location,
		AgentPoolProfiles: containerservice.ManagedClusterAgentPoolProfileArray{
			&containerservice.ManagedClusterAgentPoolProfileArgs{
				Count:              args.NodeCount,
				EnableAutoScaling:  pulumi.Bool(false),
				EnableNodePublicIP: pulumi.Bool(true),
				Mode:               pulumi.String("System"),
				Name:               pulumi.String("default"),
				OsType:             pulumi.String("Linux"),
				Type:               pulumi.String("VirtualMachineScaleSets"),
				VmSize:             pulumi.String("Standard_B2s"),
			},
		},
		DnsPrefix:         pulumi.String("akss"),
		EnableRBAC:        args.EnableRBAC,
		Identity:          &containerservice.ManagedClusterIdentityArgs{Type: containerservice.ResourceIdentityTypeSystemAssigned},
		KubernetesVersion: args.KubernetesVersion,
		NodeResourceGroup: args.NodeResourceGroup,
		ResourceGroupName: rgs.Name,
		ResourceName:      args.ResourceName,
		Sku: &containerservice.ManagedClusterSKUArgs{
			Name: pulumi.String("Basic"),
			Tier: pulumi.String("Free"),
		},
	}, pulumi.Parent(rgs))

	if err != nil {
		return nil, fmt.Errorf("error creating cluster: %v", err)
	}

	ctx.RegisterResourceOutputs(aks, pulumi.Map{
		"resourceId": k8s.ID(),
	})

	ctx.Export("resourceId", k8s.ID())
	return aks, nil
}
