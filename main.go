package main

import (
	"github.com/pkg/errors"
	route53zonev1 "github.com/project-planton/project-planton/apis/go/project/planton/provider/aws/route53zone/v1"
	"github.com/project-planton/project-planton/pkg/pulmod/stackinput"
	"github.com/project-planton/route53-zone-pulumi-module/pkg"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &route53zonev1.Route53ZoneStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return pkg.Resources(ctx, stackInput)
	})
}
