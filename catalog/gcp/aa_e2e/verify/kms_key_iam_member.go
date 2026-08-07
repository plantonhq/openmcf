package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// kmsKeyIamMemberVerifier confirms an additive (role, member) pair is present
// in (or absent from) the CRYPTO KEY's IAM policy — the key-scoped policy
// where CMEK grants live. The grant has no standalone server-side object —
// the policy itself is the source of truth, so the verifier reads the whole
// policy and looks for the exact pair.
type kmsKeyIamMemberVerifier struct{}

func (v *kmsKeyIamMemberVerifier) IDOutputKey() string { return "member" }

func (v *kmsKeyIamMemberVerifier) grantInPolicy(ctx context.Context, svc *Services, outputs map[string]string) (bool, error) {
	cryptoKeyId := outputs["crypto_key_id"]
	if cryptoKeyId == "" {
		return false, errors.New("crypto_key_id output missing")
	}

	// Version 3 so bindings with IAM conditions are returned as distinct
	// entries rather than collapsed.
	policy, err := svc.CloudKms.Projects.Locations.KeyRings.CryptoKeys.GetIamPolicy(cryptoKeyId).
		OptionsRequestedPolicyVersion(3).Context(ctx).Do()
	if err != nil {
		return false, err
	}

	for _, binding := range policy.Bindings {
		if binding.Role != outputs["role"] {
			continue
		}
		for _, member := range binding.Members {
			// IAM treats most member emails case-insensitively; compare the
			// same way so server-side casing normalization never fails a probe.
			if strings.EqualFold(member, outputs["member"]) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (v *kmsKeyIamMemberVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	found, err := v.grantInPolicy(ctx, svc, outputs)
	if err != nil {
		return errors.Wrapf(err, "failed to read IAM policy of crypto key %s", outputs["crypto_key_id"])
	}
	if !found {
		return errors.Errorf("grant of %s to %s not present in crypto key IAM policy after deploy",
			outputs["role"], outputs["member"])
	}
	return nil
}

func (v *kmsKeyIamMemberVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	found, err := v.grantInPolicy(ctx, svc, outputs)
	if err != nil {
		// Crypto keys are undeletable, so a 404 is unexpected here — but if
		// the key is somehow gone, so is every grant on it.
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "failed to read IAM policy of crypto key %s", outputs["crypto_key_id"])
	}
	if found {
		return errors.Errorf("grant of %s to %s still present in crypto key IAM policy after destroy",
			outputs["role"], outputs["member"])
	}
	return nil
}
