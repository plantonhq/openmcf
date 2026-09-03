package resources

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// GatewayDefaultImageRepo is the official nginx image the front-door
	// gateway runs. Overridable via spec.gateway.image for mirrored/air-gapped
	// registries, like every other image the operator pulls.
	GatewayDefaultImageRepo = "nginx"
	// GatewayDefaultImageTag pins a specific stable line rather than
	// "latest", so installs are reproducible.
	GatewayDefaultImageTag = "1.27-alpine"

	// GatewayDefaultLocalPort is the workstation port the port-forward is
	// expected to serve the platform at when spec.gateway.localPort is unset.
	// Sign-in bakes this port into the advertised issuer and auth callback
	// URLs, which is why it is spec-configurable rather than free-form.
	GatewayDefaultLocalPort = int32(8080)

	gatewayContainerPort = 8080
	gatewayServicePort   = 80

	// GatewayConfigKey is the data key of the nginx config TEMPLATE in the
	// ConfigMap. It is a template (mounted under /etc/nginx/templates) rather
	// than a literal conf.d file because the image's entrypoint substitutes
	// ${NGINX_LOCAL_RESOLVERS} -- the cluster DNS server from the pod's own
	// resolv.conf -- which the config needs for request-time upstream
	// resolution (see GatewayNginxConfig).
	GatewayConfigKey = "default.conf.template"

	gatewayTemplatesMountPath = "/etc/nginx/templates"
)

// GatewayConfig bundles all inputs needed to build the front-door gateway.
type GatewayConfig struct {
	CRName    string
	Namespace string
	OwnerRef  *metav1.OwnerReference

	ImageRepository string
	ImageTag        string

	// ConfigHash rolls the pod when the rendered nginx config changes
	// (nginx only reads its config at startup).
	ConfigHash string
}

// GatewayDeploymentName returns the Deployment name: "{crName}-gateway".
func GatewayDeploymentName(crName string) string {
	return fmt.Sprintf("%s-gateway", crName)
}

// GatewayServiceName returns the Service name: "{crName}-gateway".
func GatewayServiceName(crName string) string {
	return fmt.Sprintf("%s-gateway", crName)
}

// GatewayConfigMapName returns the nginx config ConfigMap name:
// "{crName}-gateway-config".
func GatewayConfigMapName(crName string) string {
	return fmt.Sprintf("%s-gateway-config", crName)
}

// GatewayPortForwardCommand renders the exact command that opens the platform
// on the workstation. Printed in the gateway component's status so nobody
// composes it by hand -- and the local port matters: sign-in URLs are pinned
// to it.
func GatewayPortForwardCommand(crName, namespace string, localPort int32) string {
	return fmt.Sprintf("kubectl port-forward -n %s svc/%s %d:%d",
		namespace, GatewayServiceName(crName), localPort, gatewayServicePort)
}

// GatewayNginxConfig renders the nginx server block that mirrors the ingress
// path layout on one origin: the gRPC-Web API prefix to the control plane,
// the identity path to the identity server, everything else to the console.
// Keeping the layout identical to the Ingress rules means the two front-door
// modes are the same architecture at different addresses -- URLs, sign-in
// callbacks, and docs carry over unchanged when an install graduates from
// port-forward to ingress.
func GatewayNginxConfig(crName, namespace string) string {
	controlPlaneUpstream := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d",
		ControlPlaneServiceName(crName), namespace, controlPlaneGrpcWebPort)
	identityUpstream := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d",
		IdentityServiceName(crName), namespace, identityServicePort)
	consoleUpstream := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d",
		ConsoleServiceName(crName), namespace, consoleServicePort)

	// Notes on the directives:
	//   - Upstreams are variables + a resolver, NOT literal proxy_pass hosts:
	//     nginx resolves literal upstreams once at startup and refuses to
	//     boot if any is absent -- but the gateway deliberately starts before
	//     the platform Services exist (it is the front door, like an ingress
	//     controller). Variables defer resolution to request time, so missing
	//     backends are an honest 502, never a crash loop. The resolver is the
	//     cluster DNS server, injected by the nginx image's own entrypoint
	//     (${NGINX_LOCAL_RESOLVERS} from the pod's resolv.conf) -- which is
	//     also why this renders as an entrypoint TEMPLATE, not a literal
	//     conf.d file. Upstream names are FQDNs because request-time
	//     resolution bypasses resolv.conf search domains.
	//   - The API location matches by string prefix (regex): gRPC-Web request
	//     paths are single segments like /ai.planton.iam...Service/Method, so
	//     path-element prefix matching would reject them (the same nuance the
	//     Ingress builder handles with use-regex on nginx controllers).
	//   - proxy_buffering off + long read timeout: gRPC-Web server-streaming
	//     (deploy progress, log tails) is a long-lived chunked response.
	//   - proxy_set_header Host preserves the browser's localhost:port origin
	//     for the identity server, whose sign-in pages and OIDC responses
	//     derive URLs from forwarded headers (KC_PROXY_HEADERS=xforwarded).
	//   - client_max_body_size: nginx defaults to 1m, which would 413 any
	//     browser payload over a megabyte while BOTH backend servers accept
	//     100m -- large manifest applies and state-file uploads are real
	//     traffic on these paths. The API and storage locations raise the
	//     cap to the servers' own limit; the servers stay the authority
	//     (the relay additionally enforces its 50MB transfer cap itself).
	return fmt.Sprintf(`resolver ${NGINX_LOCAL_RESOLVERS} valid=10s;

server {
    listen %d;

    set $controlplane_upstream %s;
    set $identity_upstream %s;
    set $console_upstream %s;

    # Browser-facing gRPC-Web API -> the control plane's gRPC-Web port.
    location ~* ^/ai\.planton\. {
        proxy_pass $controlplane_upstream;
        proxy_http_version 1.1;
        proxy_buffering off;
        proxy_request_buffering off;
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
        client_max_body_size 100m;
    }

    # Storage relay: expiring transfer URLs (state files) served by the
    # control plane on the same port as the API.
    location /storage/ {
        proxy_pass $controlplane_upstream;
        proxy_http_version 1.1;
        proxy_buffering off;
        proxy_request_buffering off;
        client_max_body_size 100m;
    }

    # Sign-in: the identity server's pages and OIDC endpoints.
    location %s {
        proxy_pass $identity_upstream;
        proxy_http_version 1.1;
        proxy_set_header Host $http_host;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $http_host;
        proxy_buffering off;
    }

    # Everything else is the web console.
    location / {
        proxy_pass $console_upstream;
        proxy_http_version 1.1;
        proxy_set_header Host $http_host;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_buffering off;
    }
}
`, gatewayContainerPort, controlPlaneUpstream, identityUpstream, consoleUpstream, IdentityPathPrefix)
}

// GatewayConfigMap builds the ConfigMap carrying the rendered nginx config
// template.
func GatewayConfigMap(crName, namespace string, ownerRef *metav1.OwnerReference) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayConfigMapName(crName),
			Namespace: namespace,
			Labels:    gatewayLabels(crName),
		},
		Data: map[string]string{
			GatewayConfigKey: GatewayNginxConfig(crName, namespace),
		},
	}
	if ownerRef != nil {
		cm.OwnerReferences = []metav1.OwnerReference{*ownerRef}
	}
	return cm
}

// GatewayDeployment builds the front-door gateway Deployment: a single nginx
// replica proxying the three platform surfaces. Deliberately no readiness
// dependency on its upstreams -- like an ingress controller, the gateway comes
// up immediately and serves 502s for backends that are still booting, so the
// port-forward command works (and shows honest progress) from the first
// minute of an install.
func GatewayDeployment(cfg GatewayConfig) *appsv1.Deployment {
	imageRepo := cfg.ImageRepository
	if imageRepo == "" {
		imageRepo = GatewayDefaultImageRepo
	}
	imageTag := cfg.ImageTag
	if imageTag == "" {
		imageTag = GatewayDefaultImageTag
	}

	labels := gatewayLabels(cfg.CRName)
	replicas := int32(1)

	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayDeploymentName(cfg.CRName),
			Namespace: cfg.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						"planton.ai/gateway-config-hash": cfg.ConfigHash,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "gateway",
						Image: fmt.Sprintf("%s:%s", imageRepo, imageTag),
						Env: []corev1.EnvVar{{
							// Opt-in switch for the image's 15-local-resolvers
							// entrypoint script: without it NGINX_LOCAL_RESOLVERS
							// is never exported and the template's resolver
							// directive survives unsubstituted (observed live as
							// a crash loop).
							Name:  "NGINX_ENTRYPOINT_LOCAL_RESOLVERS",
							Value: "true",
						}},
						Ports: []corev1.ContainerPort{{
							Name:          "http",
							ContainerPort: gatewayContainerPort,
							Protocol:      corev1.ProtocolTCP,
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "nginx-config",
							MountPath: gatewayTemplatesMountPath,
							ReadOnly:  true,
						}},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromInt32(gatewayContainerPort),
								},
							},
							InitialDelaySeconds: 2,
							PeriodSeconds:       5,
							TimeoutSeconds:      3,
							FailureThreshold:    3,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromInt32(gatewayContainerPort),
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       30,
							TimeoutSeconds:      3,
							FailureThreshold:    3,
						},
					}},
					Volumes: []corev1.Volume{{
						Name: "nginx-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: GatewayConfigMapName(cfg.CRName),
								},
							},
						},
					}},
				},
			},
		},
	}

	if cfg.OwnerRef != nil {
		deploy.OwnerReferences = []metav1.OwnerReference{*cfg.OwnerRef}
	}

	return deploy
}

// GatewayService builds the ClusterIP Service the port-forward targets,
// exposing the gateway on port 80 (mapped to the nginx container port).
func GatewayService(crName, namespace string, ownerRef *metav1.OwnerReference) *corev1.Service {
	labels := gatewayLabels(crName)

	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayServiceName(crName),
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Port:       gatewayServicePort,
				TargetPort: intstr.FromInt32(gatewayContainerPort),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}

	if ownerRef != nil {
		svc.OwnerReferences = []metav1.OwnerReference{*ownerRef}
	}

	return svc
}

func gatewayLabels(crName string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "gateway",
		"app.kubernetes.io/instance":   crName,
		"app.kubernetes.io/managed-by": ManagedByLabel,
		"app.kubernetes.io/component":  "networking",
	}
}
