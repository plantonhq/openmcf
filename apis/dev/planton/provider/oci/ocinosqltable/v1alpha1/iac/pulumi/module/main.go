package module

import (
	"github.com/pkg/errors"
	ocinosqltablev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/oci/ocinosqltable/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/oci/pulumiociprovider"
	"github.com/pulumi/pulumi-oci/sdk/v4/go/oci"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *ocinosqltablev1alpha1.OciNosqlTableStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	ociProvider, err := pulumiociprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup oci provider")
	}

	if err := nosqlTable(ctx, locals, ociProvider); err != nil {
		return errors.Wrap(err, "failed to create nosql table")
	}

	return nil
}

func pulumiOciOpt(provider *oci.Provider) pulumi.ResourceOption {
	return pulumi.Provider(provider)
}
