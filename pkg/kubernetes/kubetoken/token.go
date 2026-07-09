// Package kubetoken mints short-lived Kubernetes API-server bearer tokens in-process
// for the managed-cluster providers whose control planes accept cloud-IAM credentials
// (EKS presigned-STS tokens, GKE OAuth2 access tokens). It is the single minting seam
// shared by the ExecCredential command the deploy engines re-invoke and by in-process
// Kubernetes clients; no cloud CLI is ever shelled out to.
package kubetoken

import "time"

// Token is a short-lived bearer credential for a Kubernetes API server, paired with
// the instant it stops being honored so callers can report an honest expiry (the
// ExecCredential protocol uses it to know when to re-invoke the credential command).
type Token struct {
	Value     string
	ExpiresAt time.Time
}
