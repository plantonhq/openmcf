package module

const (
	// OpAccessKey is the access key ID (the resource's API identity).
	OpAccessKey = "access_key"
	// OpSecretKey is the secret access key -- returned ONLY at creation,
	// never retrievable again; exported as a Pulumi secret.
	OpSecretKey = "secret_key"
)
