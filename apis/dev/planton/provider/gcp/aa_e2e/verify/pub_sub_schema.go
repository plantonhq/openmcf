package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// pubSubSchemaVerifier probes a Pub/Sub schema via the Pub/Sub API. The
// schema_id output is the fully qualified resource path
// (projects/{p}/schemas/{name}) — exactly the handle a topic's
// schema_settings.schema reference consumes, so verifying with it doubles
// as proof the composition handle is honest. Posture assertions confirm
// the schema carries a type and a definition — proof GCP accepted the
// contract as configured.
type pubSubSchemaVerifier struct{}

func (v *pubSubSchemaVerifier) IDOutputKey() string { return "schema_id" }

func (v *pubSubSchemaVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	schemaID := outputs["schema_id"]
	if schemaID == "" {
		return errors.New("schema_id output missing after deploy")
	}

	// SchemaView FULL returns the definition; the default BASIC omits it.
	schema, err := svc.PubSub.Projects.Schemas.Get(schemaID).View("FULL").Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "pub/sub schema %s not found after deploy", schemaID)
	}

	if schema.Type == "" || schema.Type == "TYPE_UNSPECIFIED" {
		return errors.Errorf("pub/sub schema %s has no type after deploy (got %q)", schemaID, schema.Type)
	}
	if schema.Definition == "" {
		return errors.Errorf("pub/sub schema %s has no definition after deploy", schemaID)
	}
	return nil
}

func (v *pubSubSchemaVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	schemaID := outputs["schema_id"]
	if schemaID == "" {
		return nil
	}

	_, err := svc.PubSub.Projects.Schemas.Get(schemaID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("pub/sub schema %s still exists after destroy", schemaID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing pub/sub schema %s after destroy", schemaID)
}
