package route53zone

import (
	"github.com/pkg/errors"
	commonsdnszone "github.com/plantoncloud-inc/go-commons/network/dns/zone"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/route53zone/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/aws/pulumiawsprovider"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/dnsrecord"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/route53"
	awsclassicroute53 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

type ResourceStack struct {
	Input     *model.Route53ZoneStackInput
	AwsLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	awsNativeProvider, err := pulumiawsprovider.GetNative(ctx,
		s.Input.AwsCredential, s.Input.AwsCredential.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to setup aws provider")
	}

	route53Zone := s.Input.ApiResource

	zoneName := commonsdnszone.GetZoneName(route53Zone.Metadata.Name)

	newHostedZone, err := route53.NewHostedZone(ctx, zoneName, &route53.HostedZoneArgs{
		Name: pulumi.String(route53Zone.Metadata.Name),
		//HostedZoneTags: convertLabelsToTags(input.Labels),
	}, pulumi.Provider(awsNativeProvider))

	if err != nil {
		return errors.Wrapf(err, "failed to add zone for %s domain", route53Zone.Metadata.Name)
	}

	ctx.Export(GetManagedZoneNameOutputName(route53Zone.Metadata.Name), newHostedZone.Name)
	ctx.Export(GetManagedZoneNameserversOutputName(route53Zone.Metadata.Name), newHostedZone.NameServers)

	awsClassicProvider, err := pulumiawsprovider.GetClassic(ctx,
		s.Input.AwsCredential, s.Input.AwsCredential.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to setup aws provider")
	}

	for _, domainRecord := range route53Zone.Spec.Records {
		resName := dnsrecord.PulumiResourceName(domainRecord.Name, strings.ToLower(domainRecord.RecordType.String()))
		rs, err := awsclassicroute53.NewRecord(ctx, resName, &awsclassicroute53.RecordArgs{
			ZoneId:  newHostedZone.ID(),
			Name:    pulumi.String(domainRecord.Name),
			Ttl:     pulumi.IntPtr(int(domainRecord.TtlSeconds)),
			Type:    pulumi.String(domainRecord.RecordType),
			Records: pulumi.ToStringArray(domainRecord.Values),
		}, pulumi.Provider(awsClassicProvider))
		if err != nil {
			return errors.Wrapf(err, "failed to add %s rec", domainRecord)
		}

		ctx.Export(pulumiawsprovider.PulumiOutputName(rs, resName), rs.Records)
	}
	return nil
}

func GetManagedZoneNameOutputName(domainName string) string {
	return pulumiawsprovider.PulumiOutputName(route53.HostedZone{},
		commonsdnszone.GetZoneName(domainName), englishword.EnglishWord_name.String())
}

func GetManagedZoneNameserversOutputName(domainName string) string {
	return pulumiawsprovider.PulumiOutputName(route53.HostedZone{},
		commonsdnszone.GetZoneName(domainName), englishword.EnglishWord_nameservers.String())
}
