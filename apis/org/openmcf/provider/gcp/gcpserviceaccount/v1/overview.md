The **GcpServiceAccount** component lets you declaratively create and manage Google Cloud service accounts—plus their
optional keys and IAM role bindings—through OpenMCF’s uniform API.  
Instead of juggling `gcloud` commands or custom scripts, you describe the account once in a YAML manifest and
OpenMCF handles the rest, whether your stack uses Pulumi or Terraform under the hood.

---

## Why it matters

| Benefit                             | What it means for you                                                                                      |
|-------------------------------------|------------------------------------------------------------------------------------------------------------|
| **One-step provisioning**           | Create the service account, generate a key, and attach project- or org-level roles in a single manifest.   |
| **Consistent multi-cloud workflow** | Manage Google identities the same way you manage AWS, Azure, or Kubernetes resources in OpenMCF.    |
| **Safer IAM practices**             | Declarative YAML makes every permission explicit and auditable in Git.                                     |
| **Portable IaC**                    | The same YAML feeds either Pulumi or Terraform modules, so you can switch tools without rewriting configs. |

---

## Core capabilities

- **Service account creation** – Define the `serviceAccountId` (e.g. `my-ci-bot`) and the target `projectId`;
  OpenMCF ensures the account exists.
- **Optional key generation** – Set `createKey: true` to receive a JSON key (base64-encoded in stack outputs) for
  automation scenarios.
- **IAM bindings**
    - **Project roles** – Grant fine-grained roles in the specified project via `projectIamRoles`.
    - **Organization roles** – If you supply `orgId`, you can also attach org-level roles through `orgIamRoles`.

---

## Quick manifest example

The snippet below shows the minimum needed to spin up a CI service account with project-level logging and monitoring
roles. Notice every key is **camelCase** for consistency across all OpenMCF manifests.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpServiceAccount
metadata:
  name: ci-logs-writer
spec:
  serviceAccountId: ci-logs-writer
  projectId: my-gcp-project
  createKey: true
  projectIamRoles:
    - roles/logging.logWriter
    - roles/monitoring.metricWriter
```

Apply it with either workflow:

```bash
# Pulumi
openmcf pulumi up --manifest ci-logs-writer.yaml --stack org/project/dev

# Terraform
openmcf terraform apply --manifest ci-logs-writer.yaml --stack org/project/dev
```

---

## Next steps

1. Commit your manifest to your infrastructure repo.
2. Wire the generated key (if any) into your CI secrets manager.
3. Use additional fields like `orgIamRoles` when you need organization-wide permissions.

The **GcpServiceAccount** component brings Google Cloud identity management into the same streamlined, Git-driven
workflow you already use for the rest of your multi-cloud infrastructure with OpenMCF.
