**Overview:**

The Pulumi module provided automates the creation and management of Google Cloud Artifact Registry repositories using
Golang and Pulumi. It allows developers to define an API resource that specifies the creation of Docker, Maven, NPM, and
Python repositories within a Google Cloud project. The module handles the provisioning of service accounts with
appropriate permissions, repository creation, and access configurations based on the provided specifications.

By abstracting the complexity of setting up multiple types of repositories and managing access controls, this module
streamlines the process of publishing and consuming artifacts in a Google Cloud environment. It supports both internal
and external access configurations, making it suitable for private enterprise projects as well as open-source
initiatives that require public access to artifacts.

**StringValueOrRef Support:**

The `project_id` field uses the `StringValueOrRef` type which enables flexible resource references:

- **Literal value:** Provide a direct string value using `{value: "my-project-id"}`
- **Reference:** Reference another resource's output using `{value_from: {kind: GcpProject, name: "main-project", field_path: "status.outputs.project_id"}}`

**Current limitation**: Reference resolution is not yet implemented. Only literal `value` is used. References will fail silently (empty string).

**Future work**: Implement reference resolution in a shared library that all modules can use.
