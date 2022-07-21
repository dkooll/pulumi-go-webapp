package main

import (
	"math/rand"
	"packages/acr"
	"packages/aks"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type AcrParams struct {
	RgName           string
	Location         string
	AdminUserEnabled bool
	SkuName          string
	EnableRBAC       bool
}

type AksParams struct {
	EnableRBAC        bool
	KubernetesVersion string
	NodeResourceGroup string
	NodeCount         int
	RgName            string
	Location          string
	resourceId        pulumi.StringOutput
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		var x AcrParams
		cfg := config.New(ctx, "")
		cfg.RequireObject("acr", &x)

		acr.CreateAcr(ctx, "acr", &acr.AcrArgs{
			RegistryName:     pulumi.String("acr" + RandomString(5)),
			AdminUserEnabled: pulumi.Bool(x.AdminUserEnabled),
			SkuName:          pulumi.String(x.SkuName),
			ResourceGroup:    pulumi.String(x.RgName),
			Location:         pulumi.String(x.Location),
		})

		var y AksParams
		cfg.RequireObject("aks", &y)
		aks.CreateAks(ctx, "aks", &aks.AksArgs{
			ResourceName:      pulumi.String("aks" + RandomString(5)),
			EnableRBAC:        pulumi.Bool(y.EnableRBAC),
			KubernetesVersion: pulumi.String(y.KubernetesVersion),
			NodeResourceGroup: pulumi.String(y.NodeResourceGroup),
			NodeCount:         pulumi.Int(y.NodeCount),
			ResourceGroup:     pulumi.String(y.RgName),
			Location:          pulumi.String(y.Location),
		})
		ctx.Export("resourceId", y.resourceId)

		return nil
	})
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
