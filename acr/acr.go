package acr

import (
	"fmt"

	containerregistry "github.com/pulumi/pulumi-azure-native/sdk/go/azure/containerregistry"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Acr struct {
	pulumi.ResourceState
}

type AcrArgs struct {
	AdminUserEnabled pulumi.BoolInput
	SkuName          pulumi.StringInput
	RegistryName     pulumi.StringInput
	ResourceGroup    pulumi.StringInput
	Location         pulumi.StringInput
}

func CreateAcr(ctx *pulumi.Context, name string, args *AcrArgs, opts ...pulumi.ResourceOption) (*Acr, error) {
	acr := &Acr{}

	err := ctx.RegisterComponentResource("resource:index:Acr", name, acr, opts...)
	if err != nil {
		return nil, err
	}

	rg, err := core.NewResourceGroup(ctx, "resourceGroups", &core.ResourceGroupArgs{
		Name:     args.ResourceGroup,
		Location: args.Location,
	}, pulumi.Parent(acr))

	if err != nil {
		return nil, fmt.Errorf("error creating resource group: %v", err)
	}

	containerregistry.NewRegistry(ctx, "acr", &containerregistry.RegistryArgs{
		RegistryName:      args.RegistryName,
		AdminUserEnabled:  args.AdminUserEnabled,
		Location:          rg.Location,
		ResourceGroupName: rg.Name,
		Sku: &containerregistry.SkuArgs{
			Name: args.SkuName,
		},
	}, pulumi.Parent(rg))

	if err != nil {
		return nil, fmt.Errorf("error creating registry: %v", err)
	}
	return acr, nil
}
