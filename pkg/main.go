package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/route53zone"
	"github.com/plantoncloud/route53-zone-pulumi-module/pkg/outputs"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/route53"
	awsclassic "github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	awsclassicroute53 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

type ResourceStack struct {
	StackInput *route53zone.Route53ZoneStackInput
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	//create a variable with descriptive name for the api-resource in the input
	route53Zone := s.StackInput.ApiResource

	awsCredential := s.StackInput.AwsCredential

	//create aws provider using the credentials from the input
	awsNativeProvider, err := aws.NewProvider(ctx,
		"native-provider",
		&aws.ProviderArgs{
			AccessKey: pulumi.String(awsCredential.Spec.AccessKeyId),
			SecretKey: pulumi.String(awsCredential.Spec.SecretAccessKey),
			Region:    pulumi.String(awsCredential.Spec.Region),
		})
	if err != nil {
		return errors.Wrap(err, "failed to create aws native provider")
	}

	//replace dots with hyphens to create valid managed-zone name
	managedZoneName := strings.ReplaceAll(route53Zone.Metadata.Name, ".", "-")

	//create new hosted-zone
	createdHostedZone, err := route53.NewHostedZone(ctx,
		managedZoneName,
		&route53.HostedZoneArgs{
			Name: pulumi.String(route53Zone.Metadata.Name),
			//HostedZoneTags: convertLabelsToTags(input.Labels),
		}, pulumi.Provider(awsNativeProvider))

	if err != nil {
		return errors.Wrapf(err, "failed to create hosted-zone for %s domain",
			route53Zone.Metadata.Name)
	}

	//export important information about created hosted-zone as outputs
	ctx.Export(outputs.HostedZoneName, createdHostedZone.Name)
	ctx.Export(outputs.HostedZoneNameservers, createdHostedZone.NameServers)

	//create aws-classic provider as the native provider does not yet support creating dns-records in hosted-zone
	awsClassicProvider, err := awsclassic.NewProvider(ctx,
		"classic-provider",
		&awsclassic.ProviderArgs{
			AccessKey: pulumi.String(awsCredential.Spec.AccessKeyId),
			SecretKey: pulumi.String(awsCredential.Spec.SecretAccessKey),
			Region:    pulumi.String(awsCredential.Spec.Region),
		})
	if err != nil {
		return errors.Wrap(err, "failed to create aws classic provider")
	}

	//for each dns-record in the input spec, insert the record in the created hosted-zone
	for index, dnsRecord := range route53Zone.Spec.Records {
		_, err := awsclassicroute53.NewRecord(ctx,
			fmt.Sprintf("dns-record-%d", index),
			&awsclassicroute53.RecordArgs{
				ZoneId:  createdHostedZone.ID(),
				Name:    pulumi.String(dnsRecord.Name),
				Ttl:     pulumi.IntPtr(int(dnsRecord.TtlSeconds)),
				Type:    pulumi.String(dnsRecord.RecordType),
				Records: pulumi.ToStringArray(dnsRecord.Values),
			}, pulumi.Provider(awsClassicProvider))
		if err != nil {
			return errors.Wrapf(err, "failed to add %s rec", dnsRecord)
		}
	}
	return nil
}
