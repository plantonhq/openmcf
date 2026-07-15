package module

import (
	awscognitouserpoolclientv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awscognitouserpoolclient/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
//
// App clients carry no AwsTags map: the aws_cognito_user_pool_client resource
// is not taggable (identity tagging lives on the pool).
type Locals struct {
	Target *awscognitouserpoolclientv1.AwsCognitoUserPoolClient
	Spec   *awscognitouserpoolclientv1.AwsCognitoUserPoolClientSpec
}

func initializeLocals(_ *pulumi.Context, in *awscognitouserpoolclientv1.AwsCognitoUserPoolClientStackInput) *Locals {
	return &Locals{
		Target: in.Target,
		Spec:   in.Target.Spec,
	}
}
