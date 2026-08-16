package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// waitingRoomEvent creates the scheduled event. Every override field is
// null-means-inherit -- unset fields are never sent, so the room's value stays
// in charge during the window. Time ordering is validated by the spec's CEL up
// front; the API enforces the same rules with opaque errors.
func waitingRoomEvent(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareWaitingRoomEvent.Spec

	args := &cloudflare.WaitingRoomEventArgs{
		ZoneId:         pulumi.String(spec.ZoneId.GetValue()),
		WaitingRoomId:  pulumi.String(spec.WaitingRoomId.GetValue()),
		Name:           pulumi.String(spec.Name),
		EventStartTime: pulumi.String(spec.EventStartTime),
		EventEndTime:   pulumi.String(spec.EventEndTime),
	}

	if spec.PrequeueStartTime != "" {
		args.PrequeueStartTime = pulumi.StringPtr(spec.PrequeueStartTime)
	}
	if spec.ShuffleAtEventStart != nil {
		args.ShuffleAtEventStart = pulumi.BoolPtr(spec.GetShuffleAtEventStart())
	}
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.Suspended != nil {
		args.Suspended = pulumi.BoolPtr(spec.GetSuspended())
	}
	if spec.CustomPageHtml != "" {
		args.CustomPageHtml = pulumi.StringPtr(spec.CustomPageHtml)
	}
	if spec.DisableSessionRenewal != nil {
		args.DisableSessionRenewal = pulumi.BoolPtr(spec.GetDisableSessionRenewal())
	}
	if spec.NewUsersPerMinute != nil {
		args.NewUsersPerMinute = pulumi.IntPtr(int(spec.GetNewUsersPerMinute()))
	}
	if spec.TotalActiveUsers != nil {
		args.TotalActiveUsers = pulumi.IntPtr(int(spec.GetTotalActiveUsers()))
	}
	if spec.QueueingMethod != nil {
		args.QueueingMethod = pulumi.StringPtr(spec.GetQueueingMethod())
	}
	if spec.SessionDuration != nil {
		args.SessionDuration = pulumi.IntPtr(int(spec.GetSessionDuration()))
	}
	if spec.TurnstileAction != nil {
		args.TurnstileAction = pulumi.StringPtr(spec.GetTurnstileAction())
	}
	if spec.TurnstileMode != nil {
		args.TurnstileMode = pulumi.StringPtr(spec.GetTurnstileMode())
	}

	createdEvent, err := cloudflare.NewWaitingRoomEvent(
		ctx,
		"waiting_room_event",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create waiting room event")
	}

	ctx.Export(OpEventId, createdEvent.ID())
	ctx.Export(OpWaitingRoomId, pulumi.String(spec.WaitingRoomId.GetValue()))
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))

	return nil
}
