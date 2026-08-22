package module

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// companions emits the zone-configuration resources that ride alongside the
// zone_setting fan-out: managed header transforms, URL normalization, origin
// cloud-region hints, and the zone-wide waiting-room crawler bypass. Each is
// emitted only when its spec surface is managed.
func companions(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZoneSettings.Spec
	zoneId := pulumi.String(spec.ZoneId.GetValue())

	// Managed transforms: one zone-wide object carrying both header lists. The
	// API treats a transform missing from the list as disabled, so the module
	// always sends both lists together when either is managed.
	if len(spec.ManagedRequestHeaders) > 0 || len(spec.ManagedResponseHeaders) > 0 {
		requestHeaders := cloudflare.ManagedTransformsManagedRequestHeaderArray{}
		for _, header := range spec.ManagedRequestHeaders {
			requestHeaders = append(requestHeaders, cloudflare.ManagedTransformsManagedRequestHeaderArgs{
				Id:      pulumi.String(header.Id),
				Enabled: pulumi.Bool(header.Enabled),
			})
		}
		responseHeaders := cloudflare.ManagedTransformsManagedResponseHeaderArray{}
		for _, header := range spec.ManagedResponseHeaders {
			responseHeaders = append(responseHeaders, cloudflare.ManagedTransformsManagedResponseHeaderArgs{
				Id:      pulumi.String(header.Id),
				Enabled: pulumi.Bool(header.Enabled),
			})
		}
		if _, err := cloudflare.NewManagedTransforms(
			ctx,
			"managed_transforms",
			&cloudflare.ManagedTransformsArgs{
				ZoneId:                 zoneId,
				ManagedRequestHeaders:  requestHeaders,
				ManagedResponseHeaders: responseHeaders,
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create managed transforms")
		}
	}

	// URL normalization: a zone singleton with a real delete (destroy resets
	// normalization to Cloudflare defaults).
	if spec.UrlNormalization != nil {
		if _, err := cloudflare.NewUrlNormalizationSettings(
			ctx,
			"url_normalization",
			&cloudflare.UrlNormalizationSettingsArgs{
				ZoneId: zoneId,
				Scope:  pulumi.String(spec.UrlNormalization.Scope),
				Type:   pulumi.String(spec.UrlNormalization.Type),
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create url normalization settings")
		}
	}

	// Origin cloud regions: one API object per origin IP (the IP is the row's
	// identity, so the resource name embeds it -- reordering rows never churns
	// unrelated resources).
	for _, region := range spec.OriginCloudRegions {
		resourceName := "origin_cloud_region_" + strings.NewReplacer(".", "-", ":", "-").Replace(region.OriginIp)
		if _, err := cloudflare.NewOriginCloudRegion(
			ctx,
			resourceName,
			&cloudflare.OriginCloudRegionArgs{
				ZoneId:   zoneId,
				OriginIp: pulumi.String(region.OriginIp),
				Region:   pulumi.String(region.Region),
				Vendor:   pulumi.String(region.Vendor),
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrapf(err, "failed to create origin cloud region for %q", region.OriginIp)
		}
	}

	// Waiting-room crawler bypass: a zone singleton with NO delete at
	// Cloudflare -- destroy abandons the last-applied value.
	if spec.WaitingRoomCrawlerBypass != nil {
		if _, err := cloudflare.NewWaitingRoomSettings(
			ctx,
			"waiting_room_settings",
			&cloudflare.WaitingRoomSettingsArgs{
				ZoneId:                    zoneId,
				SearchEngineCrawlerBypass: pulumi.Bool(*spec.WaitingRoomCrawlerBypass),
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create waiting room settings")
		}
	}

	return nil
}
