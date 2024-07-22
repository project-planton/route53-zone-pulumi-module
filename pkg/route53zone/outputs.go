package route53zone

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/route53zone/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"
)

func OutputMapTransformer(stackOutput auto.OutputMap,
	input *model.Route53ZoneStackInput) *model.Route53ZoneStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &model.Route53ZoneStackOutputs{}
	}
	return &model.Route53ZoneStackOutputs{
		Nameservers: autoapistackoutput.GetStringSliceVal(stackOutput,
			GetManagedZoneNameserversOutputName(input.ApiResource.Metadata.Name)),
	}
}
