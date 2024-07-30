package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/route53zone/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	HostedZoneNameservers = "nameservers"
	HostedZoneName        = "hosted-zone-name"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *model.Route53ZoneStackInput) *model.Route53ZoneStackOutputs {
	return &model.Route53ZoneStackOutputs{
		Nameservers: autoapistackoutput.GetStringSliceVal(pulumiOutputs, HostedZoneNameservers),
	}
}
