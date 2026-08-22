package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// onOff maps the spec's managed booleans to the cache API's on/off strings.
func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

// cacheSettings emits one API object per managed cache setting.
func cacheSettings(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareCacheSettings.Spec
	zoneId := pulumi.String(spec.ZoneId.GetValue())

	// Smart Tiered Cache (the dashboard's Tiered Cache toggle). The one cache
	// setting with a real delete: destroy disables it.
	if spec.SmartTieredCache != nil {
		if _, err := cloudflare.NewTieredCache(
			ctx,
			"smart_tiered_cache",
			&cloudflare.TieredCacheArgs{
				ZoneId: zoneId,
				Value:  pulumi.String(onOff(*spec.SmartTieredCache)),
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create smart tiered cache")
		}
	}

	// Generic tiered caching. NO delete at Cloudflare: destroy abandons the
	// last-applied value.
	if spec.TieredCaching != nil {
		if _, err := cloudflare.NewArgoTieredCaching(
			ctx,
			"tiered_caching",
			&cloudflare.ArgoTieredCachingArgs{
				ZoneId: zoneId,
				Value:  pulumi.String(onOff(*spec.TieredCaching)),
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create tiered caching")
		}
	}

	// Cache Reserve (paid; storage keeps billing while on). NO delete at
	// Cloudflare: destroy abandons the last-applied value.
	if spec.CacheReserve != nil {
		if _, err := cloudflare.NewZoneCacheReserve(
			ctx,
			"cache_reserve",
			&cloudflare.ZoneCacheReserveArgs{
				ZoneId: zoneId,
				Value:  pulumi.String(onOff(*spec.CacheReserve)),
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create cache reserve")
		}
	}

	// Regional Tiered Cache. NO delete at Cloudflare.
	if spec.RegionalTieredCache != nil {
		if _, err := cloudflare.NewRegionalTieredCache(
			ctx,
			"regional_tiered_cache",
			&cloudflare.RegionalTieredCacheArgs{
				ZoneId: zoneId,
				Value:  pulumi.String(onOff(*spec.RegionalTieredCache)),
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create regional tiered cache")
		}
	}

	// Argo Smart Routing. PAID, and NO delete at Cloudflare -- destroying this
	// resource with the value still on KEEPS BILLING. Apply false first when
	// retiring the feature.
	if spec.ArgoSmartRouting != nil {
		if _, err := cloudflare.NewArgoSmartRouting(
			ctx,
			"argo_smart_routing",
			&cloudflare.ArgoSmartRoutingArgs{
				ZoneId: zoneId,
				Value:  pulumi.String(onOff(*spec.ArgoSmartRouting)),
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create argo smart routing")
		}
	}

	// Cache variants. Real delete: destroy resets variants to defaults. The
	// Pulumi SDK pluralizes the per-extension field names (avifs, jpgs, ...)
	// where the API and Terraform use the singular extension names. Only
	// extensions with at least one MIME type are sent (matching the Terraform
	// module: an unmanaged extension is omitted, never an empty list).
	if spec.CacheVariants != nil {
		variants := spec.CacheVariants
		variantsValue := cloudflare.ZoneCacheVariantsValueArgs{}
		setIfManaged := func(target *pulumi.StringArrayInput, mimeTypes []string) {
			if len(mimeTypes) > 0 {
				*target = pulumi.ToStringArray(mimeTypes)
			}
		}
		setIfManaged(&variantsValue.Avifs, variants.Avif)
		setIfManaged(&variantsValue.Bmps, variants.Bmp)
		setIfManaged(&variantsValue.Gifs, variants.Gif)
		setIfManaged(&variantsValue.Jp2s, variants.Jp2)
		setIfManaged(&variantsValue.Jpegs, variants.Jpeg)
		setIfManaged(&variantsValue.Jpgs, variants.Jpg)
		setIfManaged(&variantsValue.Jpg2s, variants.Jpg2)
		setIfManaged(&variantsValue.Pngs, variants.Png)
		setIfManaged(&variantsValue.Tifs, variants.Tif)
		setIfManaged(&variantsValue.Tiffs, variants.Tiff)
		setIfManaged(&variantsValue.Webps, variants.Webp)
		if _, err := cloudflare.NewZoneCacheVariants(
			ctx,
			"cache_variants",
			&cloudflare.ZoneCacheVariantsArgs{
				ZoneId: zoneId,
				Value:  variantsValue,
			},
			pulumi.Provider(cloudflareProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create cache variants")
		}
	}

	return nil
}
