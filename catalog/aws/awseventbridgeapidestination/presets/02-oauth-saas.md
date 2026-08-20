# OAuth SaaS Integration

Client-credentials OAuth against a SaaS token endpoint: EventBridge fetches the token itself (grant_type and scope ride the token request's body — most OAuth servers require at least grant_type) and invokes the activity endpoint under it. The client secret is secret-typed and lands in the AWS-owned secret.
