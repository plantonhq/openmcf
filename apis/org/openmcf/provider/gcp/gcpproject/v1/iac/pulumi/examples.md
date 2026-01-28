# GcpProject Examples

Below are several YAML examples demonstrating how to create and configure Google Cloud projects using OpenMCF's
`GcpProject` resource. After creating a manifest, you can apply it with Pulumi or Terraform via the OpenMCF CLI,
just like any other resource in the OpenMCF ecosystem.

```shell
# Pulumi
openmcf pulumi up --manifest <yaml-path> --stack <stack-name>

# Terraform
openmcf terraform apply --manifest <yaml-path> --stack <stack-name>
```

---

## Basic Example (Organization as Parent)

This example creates a Google Cloud project under an organization ID, links a billing account, and enables one API.
The project_id is used directly without a random suffix.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: my-basic-gcp-project
spec:
  projectId: my-basic-12345
  parentType: organization
  parentId: "987654321012"
  billingAccountId: "0123AB-4567CD-89EFGH"
  enabledApis:
    - "compute.googleapis.com"
```

---

## Example with Folder as Parent

Instead of an organization, this project will be placed under a folder.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: my-folder-gcp-project
spec:
  projectId: my-folder-proj
  parentType: folder
  parentId: "345678901234"
  billingAccountId: "0123AB-4567CD-89EFGH"
  enabledApis:
    - "storage.googleapis.com"
```

---

## Example with Random Suffix

This example enables the add_suffix option to append a random 3-character suffix to the project_id,
useful for temporary or test projects where uniqueness needs to be guaranteed.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: test-project-with-suffix
spec:
  projectId: test-project
  addSuffix: true
  parentType: organization
  parentId: "123456789012"
  billingAccountId: "0123AB-4567CD-89EFGH"
```

---

## Example with Multiple APIs, Custom Labels, and Default Network Disabled

This project enables multiple Google Cloud APIs, adds metadata labels, and removes the default VPC network immediately
after creation.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: multi-api-and-labels
spec:
  projectId: multi-api-labels
  parentType: organization
  parentId: "123456789012"
  billingAccountId: "0123AB-4567CD-89EFGH"
  labels:
    env: "dev"
    costCenter: "finops"
  disableDefaultNetwork: true
  enabledApis:
    - "compute.googleapis.com"
    - "iam.googleapis.com"
    - "cloudkms.googleapis.com"
```

---

## Example with Deletion Protection

This example enables GCP-native deletion protection, which prevents the project from being accidentally deleted.
This is a critical safety feature for production projects.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: production-project-protected
spec:
  projectId: prod-protected-proj
  parentType: folder
  parentId: "345678901234"
  billingAccountId: "0123AB-4567CD-89EFGH"
  labels:
    env: "prod"
    criticality: "high"
  disableDefaultNetwork: true
  deleteProtection: true  # CRITICAL: Prevent accidental deletion
  enabledApis:
    - "compute.googleapis.com"
    - "storage.googleapis.com"
    - "container.googleapis.com"
```

---

## Example Granting Owner Role to a Specific Member

This project grants the Owner role to a specified user during creation (for automation or administrative tasks).

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: gcp-project-with-owner
spec:
  projectId: with-owner-123
  parentType: organization
  parentId: "123456789012"
  billingAccountId: "0123AB-4567CD-89EFGH"
  ownerMember: "devops@example.com"
  enabledApis:
    - "compute.googleapis.com"
```

---

After deploying any of these manifests, you can confirm the newly created project in your GCP account:

```shell
gcloud projects list
gcloud projects describe <your-project-id>
```

You should see the project with the specified parent organization or folder, optional billing account, labels, and
enabled APIs. From there, you can continue configuring additional resources in your GCP environment or integrate your
new project with other components in the OpenMCF ecosystem.
