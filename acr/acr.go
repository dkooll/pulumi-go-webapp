package acr

import (
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

	err := ctx.RegisterComponentResource("examples:dummy:Dummy", name, acr, opts...)
	if err != nil {
		return nil, err
	}

	rg, err := core.NewResourceGroup(ctx, "resourceGroup", &core.ResourceGroupArgs{
		Name:     args.ResourceGroup,
		Location: args.Location,
	}, nil)
	if err != nil {
		return nil, err
	}

	containerregistry.NewRegistry(ctx, "acr", &containerregistry.RegistryArgs{
		RegistryName:      args.RegistryName,
		AdminUserEnabled:  args.AdminUserEnabled,
		Location:          rg.Location,
		ResourceGroupName: rg.Name,
		Sku: &containerregistry.SkuArgs{
			Name: args.SkuName,
		},
	})
	return acr, nil
}
