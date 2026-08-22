package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// waitingRoom creates the waiting room and, when the spec declares bypass
// rules, the room's rules list. Cloudflare models the rules as a separate
// per-room list with full-replacement updates -- the rules resource owns that
// WHOLE list (rules added outside the manifest are overwritten on apply and
// cleared on destroy).
//
// Advanced-subscription fields fail at the API on plans without the add-on --
// the entitlement wall is Cloudflare's, mirrored in the spec comments.
func waitingRoom(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareWaitingRoom.Spec

	args := &cloudflare.WaitingRoomArgs{
		ZoneId:            pulumi.String(spec.ZoneId.GetValue()),
		Name:              pulumi.String(spec.Name),
		Host:              pulumi.String(spec.Host),
		NewUsersPerMinute: pulumi.Int(int(spec.NewUsersPerMinute)),
		TotalActiveUsers:  pulumi.Int(int(spec.TotalActiveUsers)),
	}

	if spec.Path != nil {
		args.Path = pulumi.StringPtr(spec.GetPath())
	}
	if spec.SessionDuration != nil {
		args.SessionDuration = pulumi.IntPtr(int(spec.GetSessionDuration()))
	}
	if spec.Suspended != nil {
		args.Suspended = pulumi.BoolPtr(spec.GetSuspended())
	}
	if spec.QueueAll != nil {
		args.QueueAll = pulumi.BoolPtr(spec.GetQueueAll())
	}
	if spec.QueueingMethod != nil {
		args.QueueingMethod = pulumi.StringPtr(spec.GetQueueingMethod())
	}
	if spec.QueueingStatusCode != nil {
		args.QueueingStatusCode = pulumi.IntPtr(int(spec.GetQueueingStatusCode()))
	}
	if spec.CookieAttributes != nil {
		cookieArgs := &cloudflare.WaitingRoomCookieAttributesArgs{}
		if spec.CookieAttributes.Samesite != nil {
			cookieArgs.Samesite = pulumi.StringPtr(spec.CookieAttributes.GetSamesite())
		}
		if spec.CookieAttributes.Secure != nil {
			cookieArgs.Secure = pulumi.StringPtr(spec.CookieAttributes.GetSecure())
		}
		args.CookieAttributes = cookieArgs
	}
	if spec.CookieSuffix != "" {
		args.CookieSuffix = pulumi.StringPtr(spec.CookieSuffix)
	}
	if spec.CustomPageHtml != "" {
		args.CustomPageHtml = pulumi.StringPtr(spec.CustomPageHtml)
	}
	if spec.DefaultTemplateLanguage != nil {
		args.DefaultTemplateLanguage = pulumi.StringPtr(spec.GetDefaultTemplateLanguage())
	}
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.DisableSessionRenewal != nil {
		args.DisableSessionRenewal = pulumi.BoolPtr(spec.GetDisableSessionRenewal())
	}
	if spec.JsonResponseEnabled != nil {
		args.JsonResponseEnabled = pulumi.BoolPtr(spec.GetJsonResponseEnabled())
	}
	if len(spec.AdditionalRoutes) > 0 {
		routes := make(cloudflare.WaitingRoomAdditionalRouteArray, 0, len(spec.AdditionalRoutes))
		for _, route := range spec.AdditionalRoutes {
			routeArgs := &cloudflare.WaitingRoomAdditionalRouteArgs{
				Host: pulumi.StringPtr(route.Host),
			}
			if route.Path != nil {
				routeArgs.Path = pulumi.StringPtr(route.GetPath())
			}
			routes = append(routes, routeArgs)
		}
		args.AdditionalRoutes = routes
	}
	if len(spec.EnabledOriginCommands) > 0 {
		args.EnabledOriginCommands = pulumi.ToStringArray(spec.EnabledOriginCommands)
	}
	if spec.TurnstileAction != nil {
		args.TurnstileAction = pulumi.StringPtr(spec.GetTurnstileAction())
	}
	if spec.TurnstileMode != nil {
		args.TurnstileMode = pulumi.StringPtr(spec.GetTurnstileMode())
	}

	createdRoom, err := cloudflare.NewWaitingRoom(
		ctx,
		"waiting_room",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create waiting room")
	}

	// The room's bypass rules. The action is fixed to bypass_waiting_room (the
	// only action Cloudflare supports here) -- the module supplies it so
	// manifests never repeat a constant. enabled defaults to TRUE upstream and
	// here.
	if len(spec.BypassRules) > 0 {
		rules := make(cloudflare.WaitingRoomRulesRuleArray, 0, len(spec.BypassRules))
		for _, rule := range spec.BypassRules {
			ruleArgs := &cloudflare.WaitingRoomRulesRuleArgs{
				Action:     pulumi.String("bypass_waiting_room"),
				Expression: pulumi.String(rule.Expression),
			}
			if rule.Description != "" {
				ruleArgs.Description = pulumi.StringPtr(rule.Description)
			}
			if rule.Enabled != nil {
				ruleArgs.Enabled = pulumi.BoolPtr(rule.GetEnabled())
			} else {
				ruleArgs.Enabled = pulumi.BoolPtr(true)
			}
			rules = append(rules, ruleArgs)
		}

		_, err := cloudflare.NewWaitingRoomRules(
			ctx,
			"waiting_room_rules",
			&cloudflare.WaitingRoomRulesArgs{
				ZoneId:        pulumi.String(spec.ZoneId.GetValue()),
				WaitingRoomId: createdRoom.ID(),
				Rules:         rules,
			},
			pulumi.Provider(cloudflareProvider),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create waiting room rules")
		}
	}

	ctx.Export(OpWaitingRoomId, createdRoom.ID())
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))

	return nil
}
