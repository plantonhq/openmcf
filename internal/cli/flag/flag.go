package flag

import (
	log "github.com/sirupsen/logrus"
)

type Flag string

const (
	AtlasProviderConfig      Flag = "atlas-provider-config"
	Auth0ProviderConfig      Flag = "auth0-provider-config"
	AutoApprove              Flag = "auto-approve"
	AwsProviderConfig        Flag = "aws-provider-config"
	OpenFgaProviderConfig    Flag = "openfga-provider-config"
	Clipboard                Flag = "clipboard"
	AzureProviderConfig      Flag = "azure-provider-config"
	BackendConfig            Flag = "backend-config"
	BackendType              Flag = "backend-type"
	CloudflareProviderConfig Flag = "cloudflare-provider-config"
	ConfluentProviderConfig  Flag = "confluent-provider-config"
	Destroy                  Flag = "destroy"
	Diff                     Flag = "diff"
	Force                    Flag = "force"
	GcpProviderConfig        Flag = "gcp-provider-config"
	InputDir                 Flag = "input-dir"
	KubeContext              Flag = "kube-context"
	KubernetesProviderConfig Flag = "kubernetes-provider-config"
	KustomizeDir             Flag = "kustomize-dir"
	LocalModule              Flag = "local-module"
	Manifest                 Flag = "manifest"
	ModuleDir                Flag = "module-dir"
	ModuleVersion            Flag = "module-version"
	NoCleanup                Flag = "no-cleanup"
	OutputFile               Flag = "output-file"
	Overlay                  Flag = "overlay"
	ProjectPlantonGitRepo    Flag = "project-planton-git-repo"
	Reconfigure              Flag = "reconfigure"
	Set                      Flag = "set"
	SnowflakeProviderConfig  Flag = "snowflake-provider-config"
	Stack                    Flag = "stack"
	StackInput               Flag = "stack-input"
	Yes                      Flag = "yes"
)

func HandleFlagErrAndValue(err error, flag Flag, flagVal string) {
	if err != nil {
		log.Fatalf("error parsing %s flag. err %v", flag, err)
	}
	if flagVal == "" {
		log.Fatalf("please provide %s", flag)
	}
}

func HandleFlagErr(err error, flag Flag) {
	if err != nil {
		log.Fatalf("error parsing %s flag. err %v", flag, err)
	}
}
