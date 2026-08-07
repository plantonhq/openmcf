# AzureMonitorActionGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureMonitorActionGroupSpec** defines the configuration for creating an
Azure Monitor action group.

An action group is the notification and automation hub alerts fire into:
when a metric alert, log-query alert, activity-log alert, or Service Health
event triggers, Azure notifies every receiver in the referenced action
group. Receivers span human channels (email, SMS, voice, the Azure mobile
app), automation (webhooks, Azure Functions, Logic Apps, Automation
runbooks, Event Hubs), ITSM systems, and role-based fan-out (every member
of an ARM role).

Action groups are a GLOBAL resource: they live in a resource group but not
in a region, so notifications keep flowing during regional outages -- which
is exactly when they matter most.

One action group is typically shared by many alert rules (a "platform
on-call" group, a "database team" group), so treat the group as the stable
routing node and the alert rules as the volatile edge.

## Example

```yaml
# Offline-plan test manifest. Exercises every receiver type in one
# group -- email, SMS, voice, an Entra-authenticated webhook, app push,
# an automation runbook, a Logic App, an Azure Function, ARM-role
# fan-out, an Event Hub stream, and an ITSM connection -- plus tags.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMonitorActionGroup
metadata:
  name: test-action-group
  org: test-org
  env: dev
spec:
  resourceGroup:
    value: test-rg
  shortName: TestAG
  enabled: true
  emailReceivers:
    - name: oncall-email
      emailAddress: oncall@example.com
      useCommonAlertSchema: true
  smsReceivers:
    - name: oncall-sms
      countryCode: "1"
      phoneNumber: "5555550100"
  voiceReceivers:
    - name: oncall-voice
      countryCode: "1"
      phoneNumber: "5555550100"
  webhookReceivers:
    - name: pager
      # A domain-you-own shape: Azure rejects placeholder domains like
      # example.com server-side (WebhookServiceUriBlocked).
      serviceUri: https://hooks.yourcompany.com/pager
      useCommonAlertSchema: true
      aadAuth:
        objectId: 11111111-2222-3333-4444-555555555555
        tenantId: 99999999-8888-7777-6666-555555555555
  azureAppPushReceivers:
    - name: oncall-push
      emailAddress: oncall@example.com
  automationRunbookReceivers:
    - name: remediate
      automationAccountId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Automation/automationAccounts/test-aa
      runbookName: restart-app
      webhookResourceId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Automation/automationAccounts/test-aa/webHooks/restart-hook
      isGlobalRunbook: false
      serviceUri: https://s1events.azure-automation.net/webhooks?token=test
      useCommonAlertSchema: true
  logicAppReceivers:
    - name: workflow
      resourceId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Logic/workflows/test-wf
      callbackUrl: https://prod-1.eastus.logic.azure.com/workflows/test/triggers/manual/paths/invoke
      useCommonAlertSchema: true
  azureFunctionReceivers:
    - name: func
      functionAppResourceId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Web/sites/test-fn
      functionName: HandleAlert
      httpTriggerUrl: https://test-fn.azurewebsites.net/api/HandleAlert
      useCommonAlertSchema: true
  armRoleReceivers:
    - name: owners
      roleId:
        value: 8e3af657-a8ff-443c-a75c-2fe8c4bcb635
      useCommonAlertSchema: true
  eventHubReceivers:
    - name: siem
      eventHubName:
        value: alerts
      eventHubNamespace:
        value: test-alerts-ns
      subscriptionId: 00000000-0000-0000-0000-000000000000
      useCommonAlertSchema: true
  itsmReceivers:
    - name: servicenow
      workspaceId: 00000000-0000-0000-0000-000000000000|11111111-1111-1111-1111-111111111111
      connectionId: 22222222-2222-2222-2222-222222222222
      ticketConfiguration: '{"PayloadRevision":0,"WorkItemType":"Incident"}'
      region: eastus
  tags:
    team: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.shortName` | `string` | yes |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.emailReceivers` | `[]AzureMonitorActionGroupEmailReceiver` |  |  |  |
| `spec.emailReceivers[].name` | `string` | yes |  |  |
| `spec.emailReceivers[].emailAddress` | `string` | yes |  |  |
| `spec.emailReceivers[].useCommonAlertSchema` | `bool` |  |  |  |
| `spec.smsReceivers` | `[]AzureMonitorActionGroupSmsReceiver` |  |  |  |
| `spec.smsReceivers[].name` | `string` | yes |  |  |
| `spec.smsReceivers[].countryCode` | `string` | yes |  |  |
| `spec.smsReceivers[].phoneNumber` | `string` | yes |  |  |
| `spec.voiceReceivers` | `[]AzureMonitorActionGroupVoiceReceiver` |  |  |  |
| `spec.voiceReceivers[].name` | `string` | yes |  |  |
| `spec.voiceReceivers[].countryCode` | `string` | yes |  |  |
| `spec.voiceReceivers[].phoneNumber` | `string` | yes |  |  |
| `spec.webhookReceivers` | `[]AzureMonitorActionGroupWebhookReceiver` |  |  |  |
| `spec.webhookReceivers[].name` | `string` | yes |  |  |
| `spec.webhookReceivers[].serviceUri` | `string` | yes |  |  |
| `spec.webhookReceivers[].useCommonAlertSchema` | `bool` |  |  |  |
| `spec.webhookReceivers[].aadAuth` | `AzureMonitorActionGroupWebhookAadAuth` |  |  |  |
| `spec.webhookReceivers[].aadAuth.objectId` | `string` | yes |  |  |
| `spec.webhookReceivers[].aadAuth.identifierUri` | `string` |  |  |  |
| `spec.webhookReceivers[].aadAuth.tenantId` | `string` |  |  |  |
| `spec.azureAppPushReceivers` | `[]AzureMonitorActionGroupAzureAppPushReceiver` |  |  |  |
| `spec.azureAppPushReceivers[].name` | `string` | yes |  |  |
| `spec.azureAppPushReceivers[].emailAddress` | `string` | yes |  |  |
| `spec.automationRunbookReceivers` | `[]AzureMonitorActionGroupAutomationRunbookReceiver` |  |  |  |
| `spec.automationRunbookReceivers[].name` | `string` | yes |  |  |
| `spec.automationRunbookReceivers[].automationAccountId` | `string` | yes |  |  |
| `spec.automationRunbookReceivers[].runbookName` | `string` | yes |  |  |
| `spec.automationRunbookReceivers[].webhookResourceId` | `string` | yes |  |  |
| `spec.automationRunbookReceivers[].isGlobalRunbook` | `bool` |  |  |  |
| `spec.automationRunbookReceivers[].serviceUri` | `string` | yes |  |  |
| `spec.automationRunbookReceivers[].useCommonAlertSchema` | `bool` |  |  |  |
| `spec.logicAppReceivers` | `[]AzureMonitorActionGroupLogicAppReceiver` |  |  |  |
| `spec.logicAppReceivers[].name` | `string` | yes |  |  |
| `spec.logicAppReceivers[].resourceId` | `string` | yes |  |  |
| `spec.logicAppReceivers[].callbackUrl` | `string` | yes |  |  |
| `spec.logicAppReceivers[].useCommonAlertSchema` | `bool` |  |  |  |
| `spec.azureFunctionReceivers` | `[]AzureMonitorActionGroupAzureFunctionReceiver` |  |  |  |
| `spec.azureFunctionReceivers[].name` | `string` | yes |  |  |
| `spec.azureFunctionReceivers[].functionAppResourceId` | `string \| valueFrom` | yes |  | AzureFunctionApp (`status.outputs.function_app_id`) |
| `spec.azureFunctionReceivers[].functionName` | `string` | yes |  |  |
| `spec.azureFunctionReceivers[].httpTriggerUrl` | `string` | yes |  |  |
| `spec.azureFunctionReceivers[].useCommonAlertSchema` | `bool` |  |  |  |
| `spec.armRoleReceivers` | `[]AzureMonitorActionGroupArmRoleReceiver` |  |  |  |
| `spec.armRoleReceivers[].name` | `string` | yes |  |  |
| `spec.armRoleReceivers[].roleId` | `string \| valueFrom` | yes |  | AzureRoleDefinition (`status.outputs.role_definition_guid`) |
| `spec.armRoleReceivers[].useCommonAlertSchema` | `bool` |  |  |  |
| `spec.eventHubReceivers` | `[]AzureMonitorActionGroupEventHubReceiver` |  |  |  |
| `spec.eventHubReceivers[].name` | `string` | yes |  |  |
| `spec.eventHubReceivers[].eventHubName` | `string \| valueFrom` | yes |  | AzureEventHub (`status.outputs.event_hub_name`) |
| `spec.eventHubReceivers[].eventHubNamespace` | `string \| valueFrom` | yes |  | AzureEventHubNamespace (`status.outputs.namespace_name`) |
| `spec.eventHubReceivers[].tenantId` | `string` |  |  |  |
| `spec.eventHubReceivers[].subscriptionId` | `string` |  |  |  |
| `spec.eventHubReceivers[].useCommonAlertSchema` | `bool` |  |  |  |
| `spec.itsmReceivers` | `[]AzureMonitorActionGroupItsmReceiver` |  |  |  |
| `spec.itsmReceivers[].name` | `string` | yes |  |  |
| `spec.itsmReceivers[].workspaceId` | `string` | yes |  |  |
| `spec.itsmReceivers[].connectionId` | `string` | yes |  |  |
| `spec.itsmReceivers[].ticketConfiguration` | `string` | yes |  |  |
| `spec.itsmReceivers[].region` | `string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

The Azure Resource Group where the action group will be created.
Can be a literal string or a reference to an AzureResourceGroup output.
(Action groups are global -- there is no region to choose.)

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.shortName

`string` · required

The short name of the action group -- 1 to 12 characters, shown as the
sender identity in SMS messages and mobile push notifications. Make it
recognizable on a phone screen ("PltOnCall", "DBTeam").

- rule: {"required":true,"string":{"minLen":"1","maxLen":"12"}}

### spec.enabled

`bool` · optional (explicit presence)

Whether the action group is active. A disabled group silently swallows
every alert that fires into it -- useful as a maintenance-window kill
switch without editing alert rules.
Default: true

- default: `true`

### spec.emailReceivers

`[]AzureMonitorActionGroupEmailReceiver`

Email receivers -- each sends alert notifications to one address.

### spec.emailReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.emailReceivers[].emailAddress

`string` · required

The email address to notify.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.emailReceivers[].useCommonAlertSchema

`bool`

Whether to send the common alert schema payload -- one consistent JSON
shape across all alert types. Prefer true for anything parsed by
software; the legacy per-alert-type payloads exist for old integrations.

### spec.smsReceivers

`[]AzureMonitorActionGroupSmsReceiver`

SMS receivers -- each texts one phone number. SMS carries no payload
options (no common alert schema); rate limits apply per Azure.

### spec.smsReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.smsReceivers[].countryCode

`string` · required

The phone's country code (for example "1" for the US, "91" for India).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.smsReceivers[].phoneNumber

`string` · required

The phone number, without the country code.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.voiceReceivers

`[]AzureMonitorActionGroupVoiceReceiver`

Voice-call receivers -- each places an automated call to one number.

### spec.voiceReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.voiceReceivers[].countryCode

`string` · required

The phone's country code (for example "1" for the US).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.voiceReceivers[].phoneNumber

`string` · required

The phone number, without the country code.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.webhookReceivers

`[]AzureMonitorActionGroupWebhookReceiver`

Webhook receivers -- each POSTs the alert payload to an HTTP(S) endpoint,
optionally authenticating with a Microsoft Entra identity.

### spec.webhookReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.webhookReceivers[].serviceUri

`string` · required

The endpoint to POST the alert payload to (http:// or https://).

- rule: the webhook service_uri must be an http:// or https:// URL
- rule: {"required":true}

### spec.webhookReceivers[].useCommonAlertSchema

`bool`

Whether to send the common alert schema payload (recommended for
anything parsed by software).

### spec.webhookReceivers[].aadAuth

`AzureMonitorActionGroupWebhookAadAuth`

Authenticate the webhook call with a Microsoft Entra identity instead of
relying on a secret baked into the URL -- the keyless webhook posture.

### spec.webhookReceivers[].aadAuth.objectId

`string` · required

The object ID of the Entra application (service principal) the call
authenticates as.

- rule: {"required":true,"string":{"uuid":true}}

### spec.webhookReceivers[].aadAuth.identifierUri

`string`

The identifier URI of the Entra application (the token audience).
When empty, Azure derives it from the application.

### spec.webhookReceivers[].aadAuth.tenantId

`string`

The Entra tenant of the application. When empty, Azure uses the
subscription's home tenant.

- rule: tenant_id must be a UUID (the Entra tenant's directory ID)

### spec.azureAppPushReceivers

`[]AzureMonitorActionGroupAzureAppPushReceiver`

Azure mobile app push receivers -- each pushes to the Azure app account
signed in with the given email.

### spec.azureAppPushReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.azureAppPushReceivers[].emailAddress

`string` · required

The email address the Azure mobile app account is signed in with.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.automationRunbookReceivers

`[]AzureMonitorActionGroupAutomationRunbookReceiver`

Automation runbook receivers -- each starts an Azure Automation runbook
(remediation-as-a-response).

### spec.automationRunbookReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.automationRunbookReceivers[].automationAccountId

`string` · required

The ARM ID of the Automation Account hosting the runbook.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.automationRunbookReceivers[].runbookName

`string` · required

The name of the runbook to start.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.automationRunbookReceivers[].webhookResourceId

`string` · required

The ARM ID of the runbook's webhook resource (the Automation webhook
that actually starts the runbook).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.automationRunbookReceivers[].isGlobalRunbook

`bool`

Whether the runbook is a global runbook (true) or a user runbook
(false). Azure requires this to be stated explicitly.

### spec.automationRunbookReceivers[].serviceUri

`string` · required

The webhook's invocation URI.

- rule: the runbook service_uri must be an http:// or https:// URL
- rule: {"required":true}

### spec.automationRunbookReceivers[].useCommonAlertSchema

`bool`

Whether to send the common alert schema payload.

### spec.logicAppReceivers

`[]AzureMonitorActionGroupLogicAppReceiver`

Logic App receivers -- each triggers a Logic App workflow.

### spec.logicAppReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.logicAppReceivers[].resourceId

`string` · required

The ARM ID of the Logic App workflow.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.logicAppReceivers[].callbackUrl

`string` · required

The Logic App trigger's callback URL (from the workflow's HTTP
request trigger).

- rule: the Logic App callback_url must be an http:// or https:// URL
- rule: {"required":true}

### spec.logicAppReceivers[].useCommonAlertSchema

`bool`

Whether to send the common alert schema payload.

### spec.azureFunctionReceivers

`[]AzureMonitorActionGroupAzureFunctionReceiver`

Azure Function receivers -- each invokes an HTTP-triggered function.

### spec.azureFunctionReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.azureFunctionReceivers[].functionAppResourceId

`string | valueFrom` · required

The Function App hosting the function. Can be a literal ARM ID or a
reference to an AzureFunctionApp output.

- references: AzureFunctionApp (`status.outputs.function_app_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFunctionApp, name: <that resource's name>, fieldPath: status.outputs.function_app_id}} -- a bare string does not parse

### spec.azureFunctionReceivers[].functionName

`string` · required

The name of the function within the app.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.azureFunctionReceivers[].httpTriggerUrl

`string` · required

The function's HTTP trigger URL (including its code parameter when the
function uses key-based auth).

- rule: the function http_trigger_url must be an http:// or https:// URL
- rule: {"required":true}

### spec.azureFunctionReceivers[].useCommonAlertSchema

`bool`

Whether to send the common alert schema payload.

### spec.armRoleReceivers

`[]AzureMonitorActionGroupArmRoleReceiver`

ARM role receivers -- each notifies every user assigned the given role
on the subscription (role-based fan-out; no address list to maintain).

### spec.armRoleReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.armRoleReceivers[].roleId

`string | valueFrom` · required

The role definition GUID whose members are notified -- a built-in role's
well-known GUID (for example Owner:
8e3af657-a8ff-443c-a75c-2fe8c4bcb635) or a custom role. Can be a
literal GUID or a reference to an AzureRoleDefinition output.

- references: AzureRoleDefinition (`status.outputs.role_definition_guid`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureRoleDefinition, name: <that resource's name>, fieldPath: status.outputs.role_definition_guid}} -- a bare string does not parse

### spec.armRoleReceivers[].useCommonAlertSchema

`bool`

Whether to send the common alert schema payload.

### spec.eventHubReceivers

`[]AzureMonitorActionGroupEventHubReceiver`

Event Hub receivers -- each streams the alert payload to an Event Hub
(the SIEM / external-pipeline path).

### spec.eventHubReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.eventHubReceivers[].eventHubName

`string | valueFrom` · required

The name of the Event Hub to stream to. Can be a literal name or a
reference to an AzureEventHub output (the hub must live in the
namespace below).

- references: AzureEventHub (`status.outputs.event_hub_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHub, name: <that resource's name>, fieldPath: status.outputs.event_hub_name}} -- a bare string does not parse

### spec.eventHubReceivers[].eventHubNamespace

`string | valueFrom` · required

The name of the Event Hub namespace holding the hub. Can be a literal
name or a reference to an AzureEventHubNamespace output.

- references: AzureEventHubNamespace (`status.outputs.namespace_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHubNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_name}} -- a bare string does not parse

### spec.eventHubReceivers[].tenantId

`string`

The Entra tenant of the Event Hub. When empty, Azure uses the
subscription's home tenant.

- rule: tenant_id must be a UUID (the Entra tenant's directory ID)

### spec.eventHubReceivers[].subscriptionId

`string`

The subscription holding the Event Hub. When empty, Azure uses the
action group's own subscription.

- rule: subscription_id must be a UUID

### spec.eventHubReceivers[].useCommonAlertSchema

`bool`

Whether to send the common alert schema payload.

### spec.itsmReceivers

`[]AzureMonitorActionGroupItsmReceiver`

ITSM receivers -- each creates a work item in a connected IT Service
Management system through an ITSM connection on a Log Analytics
workspace.

### spec.itsmReceivers[].name

`string` · required

The receiver's name, unique within the action group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.itsmReceivers[].workspaceId

`string` · required

The composite ID of the Log Analytics workspace the ITSM connection is
defined on. Format: {subscription_id}|{workspace_customer_id} -- the
subscription GUID and the workspace's customer GUID joined by a pipe
(Azure's ITSM addressing, not an ARM resource ID).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.itsmReceivers[].connectionId

`string` · required

The ITSM connection's unique identifier (from the ITSM Connector
configured on the workspace).

- rule: {"required":true,"string":{"uuid":true}}

### spec.itsmReceivers[].ticketConfiguration

`string` · required

The ITSM ticket configuration as a JSON document. Azure requires at
least the PayloadRevision and WorkItemType keys (for example
{"PayloadRevision":0,"WorkItemType":"Incident"}).

- rule: ticket_configuration must be a JSON document containing the PayloadRevision and WorkItemType keys, e.g. {"PayloadRevision":0,"WorkItemType":"Incident"}
- rule: {"required":true}

### spec.itsmReceivers[].region

`string` · required

The Azure region of the ITSM connection (for example "eastus").

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.tags

`map<string, string>`

Tags to apply to the action group, merged over the Planton-derived
metadata tags (user values win on key conflicts).

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMonitorActionGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.action_group_id` | `string` | The Azure Resource Manager ID of the action group. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Insights/actionGroups/{name} Referenced by: AzureMonitorMetricAlert and AzureMonitorScheduledQueryAlert actions. |
| `status.outputs.action_group_name` | `string` | The name of the action group. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.azureFunctionReceivers[].functionAppResourceId` | AzureFunctionApp | `status.outputs.function_app_id` |
| `spec.armRoleReceivers[].roleId` | AzureRoleDefinition | `status.outputs.role_definition_guid` |
| `spec.eventHubReceivers[].eventHubName` | AzureEventHub | `status.outputs.event_hub_name` |
| `spec.eventHubReceivers[].eventHubNamespace` | AzureEventHubNamespace | `status.outputs.namespace_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureMonitorActivityLogAlert | `spec.actions[].actionGroupId` | `status.outputs.action_group_id` |
| AzureMonitorMetricAlert | `spec.actions[].actionGroupId` | `status.outputs.action_group_id` |
| AzureMonitorScheduledQueryAlert | `spec.action.actionGroupIds` | `status.outputs.action_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
