package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrockfoundation"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// modelAccess accepts the foundation model's public marketplace agreement
// (and renders the account use-case form when declared) and exports
// outputs.
func modelAccess(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// The model's current public offer -- offer tokens are short-lived, so
	// the lookup happens on every deploy rather than ever living in a
	// manifest.
	offerType := "PUBLIC"
	offers, err := bedrockfoundation.GetModelAgreementOffers(ctx, &bedrockfoundation.GetModelAgreementOffersArgs{
		ModelId:   spec.ModelId,
		OfferType: &offerType,
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrapf(err, "look up public agreement offer for model %q", spec.ModelId)
	}
	if len(offers.Offers) == 0 {
		return errors.Errorf("model %q has no public agreement offer in this region", spec.ModelId)
	}

	// The account-global use-case form. PutUseCaseForModelAccess is
	// upsert-convergent (re-putting identical form data is a no-op) and
	// the provider imports an existing identical form instead of failing;
	// a DIFFERENT existing form fails the deploy loudly -- the form cannot
	// be updated through this API.
	var agreementDeps []pulumi.Resource
	if spec.UseCaseForm != nil {
		formJSON, err := json.Marshal(spec.UseCaseForm.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal use-case form to JSON")
		}
		createdUseCase, err := bedrock.NewUseCaseForModelAccess(ctx, "use-case-form", &bedrock.UseCaseForModelAccessArgs{
			FormData: pulumi.String(string(formJSON)),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create use-case form")
		}
		agreementDeps = append(agreementDeps, createdUseCase)
	}

	// The agreement itself. Create waits for the agreement to reach
	// AVAILABLE (the provider polls GetFoundationModelAvailability);
	// delete cancels the agreement and retries through AWS's transient
	// ConflictException window. Anthropic agreements activate only once
	// the account's use-case form is on file -- ordered after the form
	// whenever one is managed.
	//
	// IgnoreChanges on the token mirrors the Terraform module's lifecycle
	// guard: offer tokens are short-lived and minted FRESH on every
	// preview, and the argument is ForceNew -- without the guard every
	// re-apply would REPLACE the agreement (destroy = access revocation).
	// The token only matters at create; the accepted agreement is keyed
	// by model_id.
	createdAgreement, err := bedrockfoundation.NewModelAgreement(ctx, "agreement", &bedrockfoundation.ModelAgreementArgs{
		ModelId:    pulumi.String(spec.ModelId),
		OfferToken: pulumi.String(offers.Offers[0].OfferToken),
	}, pulumi.Provider(provider), pulumi.DependsOn(agreementDeps),
		pulumi.IgnoreChanges([]string{"offerToken"}))
	if err != nil {
		return errors.Wrap(err, "create foundation model agreement")
	}

	ctx.Export(OpModelId, createdAgreement.ModelId)

	return nil
}
