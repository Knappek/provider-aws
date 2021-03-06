package v1beta1

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/reference"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
)

// QueueARN returns ARN of the Queue resource.
func QueueARN() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		cr, ok := mg.(*Queue)
		if !ok {
			return ""
		}
		return cr.Status.AtProvider.ARN
	}
}

// ResolveReferences of this Queue
func (mg *Queue) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	if mg.Spec.ForProvider.RedrivePolicy != nil {
		// Resolve spec.forProvider.redrivePolicy.deadLetterQueueArn
		rsp, err := r.Resolve(ctx, reference.ResolutionRequest{
			CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.RedrivePolicy.DeadLetterQueueARN),
			Reference:    mg.Spec.ForProvider.RedrivePolicy.DeadLetterQueueARNRef,
			Selector:     mg.Spec.ForProvider.RedrivePolicy.DeadLetterQueueARNSelector,
			To:           reference.To{Managed: &Queue{}, List: &QueueList{}},
			Extract:      QueueARN(),
		})
		if err != nil {
			return errors.Wrap(err, "spec.forProvider.redrivePolicy.deadLetterQueueArn")
		}
		mg.Spec.ForProvider.RedrivePolicy.DeadLetterQueueARN = aws.String(rsp.ResolvedValue)
		mg.Spec.ForProvider.RedrivePolicy.DeadLetterQueueARNRef = rsp.ResolvedReference
	}
	return nil
}
