# Overview

The **KubernetesDeployment** API resource streamlines and standardizes how developers deploy microservices onto
Kubernetes clusters. This Pulumi module interprets a `KubernetesDeploymentStackInput`, which includes core details
like Kubernetes credentials, Docker credentials, and the microservice configuration. From this input, the module
automatically creates and manages:

- **Kubernetes Namespaces**
- **Deployments**, including main containers and optional sidecars
- **Services** for exposing your microservice internally
- **Ingress** resources (optionally Istio-based) for external or mesh networking
- **Secrets** for storing sensitive environment variables or credentials

Developers simply provide a declarative resource specification – focusing on container images, resource requests,
environment variables, and ingress preferences – while the module handles the underlying Kubernetes constructs.

### Key Features

1. **Deployment Automation**  
   Eliminates the need to write manual Kubernetes definitions for namespaces, deployments, or services. The module
   compiles your specification into a robust, ready-to-run setup.

2. **Environment & Secrets**  
   Use `env.variables` for straightforward key-value pairs and `env.secrets` for sensitive data. Secrets are stored in
   Kubernetes as `"Opaque"` secrets and automatically injected into the container.

3. **Resource Allocation**  
   Specify CPU and memory requests/limits for containers (both main and sidecar). This ensures your microservice can run
   efficiently without jeopardizing the cluster’s stability.

4. **Sidecar Support**  
   Integrate additional containers (logging, monitoring, proxies, etc.) alongside your main container for advanced use
   cases like service meshes or data scrapers.

5. **Optional Ingress**  
   Leverage Istio or other ingress controllers. The module automatically creates gateway and routing resources to expose
   your microservice securely if `ingress.is_enabled` is set to `true`.

6. **Scalability**  
   Combine `availability.minReplicas` and optional horizontal pod autoscaling to dynamically scale your microservices
   based on CPU or memory utilization thresholds.

Overall, **KubernetesDeployment** helps you focus on application code and business logic. By delegating the Kubernetes
resource orchestration to this module, you gain a cleaner, more consistent deployment experience across development,
staging, and production environments.
