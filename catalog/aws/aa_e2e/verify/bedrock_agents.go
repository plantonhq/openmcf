package verify

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagent"
	awstypes "github.com/aws/aws-sdk-go-v2/service/bedrockagent/types"
	"github.com/pkg/errors"
)

func bedrockAgentClient(cfg aws.Config, region string) *bedrockagent.Client {
	return bedrockagent.NewFromConfig(cfg, func(o *bedrockagent.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// isBedrockAgentNotFound matches the bedrock-agent control plane's
// resource-not-found class.
func isBedrockAgentNotFound(err error) bool {
	var nf *awstypes.ResourceNotFoundException
	return errors.As(err, &nf) || strings.Contains(err.Error(), "ResourceNotFoundException")
}

// bedrockAgentVerifier verifies AwsBedrockAgent components. The module
// prepares the agent after every change and waits for PREPARED, so a
// healthy deploy always lands there.
type bedrockAgentVerifier struct{}

func (*bedrockAgentVerifier) IDOutputKey() string { return "agent_id" }

func (*bedrockAgentVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockAgentClient(cfg, region).GetAgent(ctx, &bedrockagent.GetAgentInput{
		AgentId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "GetAgent(%s)", id)
	}
	status := out.Agent.AgentStatus
	if status == awstypes.AgentStatusFailed || status == awstypes.AgentStatusDeleting {
		return errors.Errorf("agent %s is in status %s", id, status)
	}
	if status != awstypes.AgentStatusPrepared {
		return errors.Errorf("agent %s is %s, expected PREPARED (the modules always prepare)", id, status)
	}
	return nil
}

func (*bedrockAgentVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := bedrockAgentClient(cfg, region).GetAgent(ctx, &bedrockagent.GetAgentInput{
		AgentId: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("agent %s still exists", id)
	}
	if isBedrockAgentNotFound(err) {
		return nil
	}
	return errors.Wrapf(err, "GetAgent(%s) during absence check", id)
}

// bedrockKnowledgeBaseVerifier verifies AwsBedrockKnowledgeBase
// components. failure_reasons surface in the error when the knowledge
// base is unhealthy.
type bedrockKnowledgeBaseVerifier struct{}

func (*bedrockKnowledgeBaseVerifier) IDOutputKey() string { return "knowledge_base_id" }

func (*bedrockKnowledgeBaseVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockAgentClient(cfg, region).GetKnowledgeBase(ctx, &bedrockagent.GetKnowledgeBaseInput{
		KnowledgeBaseId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "GetKnowledgeBase(%s)", id)
	}
	status := out.KnowledgeBase.Status
	if status != awstypes.KnowledgeBaseStatusActive {
		return errors.Errorf("knowledge base %s is %s (failure reasons: %s), expected ACTIVE",
			id, status, strings.Join(out.KnowledgeBase.FailureReasons, "; "))
	}
	return nil
}

func (*bedrockKnowledgeBaseVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := bedrockAgentClient(cfg, region).GetKnowledgeBase(ctx, &bedrockagent.GetKnowledgeBaseInput{
		KnowledgeBaseId: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("knowledge base %s still exists", id)
	}
	if isBedrockAgentNotFound(err) {
		return nil
	}
	return errors.Wrapf(err, "GetKnowledgeBase(%s) during absence check", id)
}

// bedrockFlowVerifier verifies AwsBedrockFlow components. The module does
// not prepare flows, so a healthy deploy is NotPrepared or Prepared --
// only Failed is a defect.
type bedrockFlowVerifier struct{}

func (*bedrockFlowVerifier) IDOutputKey() string { return "flow_id" }

func (*bedrockFlowVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockAgentClient(cfg, region).GetFlow(ctx, &bedrockagent.GetFlowInput{
		FlowIdentifier: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "GetFlow(%s)", id)
	}
	if out.Status == awstypes.FlowStatusFailed {
		return errors.Errorf("flow %s is in status %s", id, out.Status)
	}
	return nil
}

func (*bedrockFlowVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := bedrockAgentClient(cfg, region).GetFlow(ctx, &bedrockagent.GetFlowInput{
		FlowIdentifier: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("flow %s still exists", id)
	}
	if isBedrockAgentNotFound(err) {
		return nil
	}
	return errors.Wrapf(err, "GetFlow(%s) during absence check", id)
}

// bedrockPromptVerifier verifies AwsBedrockPrompt components. Prompts
// carry no status field -- existence of the DRAFT version is the
// assertion.
type bedrockPromptVerifier struct{}

func (*bedrockPromptVerifier) IDOutputKey() string { return "prompt_id" }

func (*bedrockPromptVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockAgentClient(cfg, region).GetPrompt(ctx, &bedrockagent.GetPromptInput{
		PromptIdentifier: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "GetPrompt(%s)", id)
	}
	if len(out.Variants) == 0 {
		return errors.Errorf("prompt %s exists but carries no variants", id)
	}
	return nil
}

func (*bedrockPromptVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := bedrockAgentClient(cfg, region).GetPrompt(ctx, &bedrockagent.GetPromptInput{
		PromptIdentifier: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("prompt %s still exists", id)
	}
	if isBedrockAgentNotFound(err) {
		return nil
	}
	return errors.Wrapf(err, "GetPrompt(%s) during absence check", id)
}
