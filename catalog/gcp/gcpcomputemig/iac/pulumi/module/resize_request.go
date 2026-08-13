package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// resizeRequests creates one queued capacity request per spec entry —
// zonal or regional per the location selector — each wired to the group
// manager by NAME. Requests are immutable one-shot asks: any change
// replaces the request (queues a new ask); destroying an ACTIVE request
// cancels it.
func resizeRequests(
	ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider,
	groupManager *igmResult,
) error {

	spec := locals.GcpComputeMig.Spec

	for _, request := range spec.ResizeRequests {
		if locals.IsRegional {
			args := &compute.RegionResizeRequestArgs{
				Name:                 pulumi.String(request.RequestName),
				Region:               pulumi.StringPtr(spec.Region),
				InstanceGroupManager: groupManager.Name,
				ResizeBy:             pulumi.Int(int(request.ResizeBy)),
			}
			if spec.ProjectId.GetValue() != "" {
				args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
			}
			if request.Description != "" {
				args.Description = pulumi.StringPtr(request.Description)
			}
			// The provider models the duration's seconds as a STRING —
			// rendered from the spec's int64 identically on both engines.
			if request.RequestedRunDurationSeconds != nil {
				args.RequestedRunDuration = &compute.RegionResizeRequestRequestedRunDurationArgs{
					Seconds: pulumi.String(strconv.FormatInt(request.GetRequestedRunDurationSeconds(), 10)),
				}
			}
			if spec.DeletionPolicy != "" {
				args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
			}
			if _, err := compute.NewRegionResizeRequest(ctx,
				request.RequestName,
				args,
				pulumi.Provider(gcpProvider),
			); err != nil {
				return errors.Wrapf(err, "failed to create regional resize request %s", request.RequestName)
			}
			continue
		}

		args := &compute.ResizeRequestArgs{
			Name:                 pulumi.String(request.RequestName),
			Zone:                 pulumi.StringPtr(spec.Zone),
			InstanceGroupManager: groupManager.Name,
			ResizeBy:             pulumi.Int(int(request.ResizeBy)),
		}
		if spec.ProjectId.GetValue() != "" {
			args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
		}
		if request.Description != "" {
			args.Description = pulumi.StringPtr(request.Description)
		}
		if request.RequestedRunDurationSeconds != nil {
			args.RequestedRunDuration = &compute.ResizeRequestRequestedRunDurationArgs{
				Seconds: pulumi.String(strconv.FormatInt(request.GetRequestedRunDurationSeconds(), 10)),
			}
		}
		if spec.DeletionPolicy != "" {
			args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
		}
		if _, err := compute.NewResizeRequest(ctx,
			request.RequestName,
			args,
			pulumi.Provider(gcpProvider),
		); err != nil {
			return errors.Wrapf(err, "failed to create resize request %s", request.RequestName)
		}
	}

	return nil
}
