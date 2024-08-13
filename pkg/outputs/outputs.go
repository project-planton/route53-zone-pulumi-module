package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/route53zone"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	HostedZoneNameservers = "nameservers"
	HostedZoneName        = "hosted-zone-name"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *route53zone.Route53ZoneStackInput) *route53zone.Route53ZoneStackOutputs {
	return &route53zone.Route53ZoneStackOutputs{
		Nameservers: autoapistackoutput.GetStringSliceVal(pulumiOutputs, HostedZoneNameservers),
	}
}
