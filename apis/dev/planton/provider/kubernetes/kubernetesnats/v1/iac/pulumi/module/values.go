package module

import (
	"fmt"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesnatsv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesnats/v1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// SECRET DISCIPLINE (load-bearing): nothing rendered here carries
// credential material. Declared users' passwords live ONLY in the
// module-generated `<name>-auth` Secret; the values below wire each one as
// a secretKeyRef-backed env var, and the config carries `$NATS_PW_<i>`
// references (unquoted via the chart's `<< >>` syntax) the server
// resolves from environment at load — verified in the server's config
// parser at the pin.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every
// typed mapping below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// fullnameOverride pins every child name (`<name>`,
	// `<name>-headless`, `<name>-box`, ...) to the resource name; the
	// exported outputs are built from that contract.
	values["fullnameOverride"] = locals.ReleaseName

	// ---- config ------------------------------------------------------------
	config := map[string]interface{}{}

	// Clustering: with the block disabled the chart runs exactly one
	// server (its own template ignores replicas then).
	if cluster := spec.GetCluster(); cluster != nil && cluster.GetEnabled() {
		replicas := int32(3)
		if cluster.Replicas != nil {
			replicas = cluster.GetReplicas()
		}
		config["cluster"] = map[string]interface{}{
			"enabled":  true,
			"replicas": int(replicas),
		}
	}

	// JetStream: ON by default (the kind's persistent-messaging posture;
	// the chart's raw default is off) — rendered explicitly either way
	// so both engines state the posture.
	jetstream := map[string]interface{}{"enabled": locals.JetStreamEnabled}
	if locals.JetStreamEnabled {
		js := spec.GetJetStream()
		diskSize := js.GetDiskSize()
		if diskSize == "" {
			diskSize = "10Gi"
		}
		pvc := map[string]interface{}{"size": diskSize}
		if js.GetStorageClass().GetValue() != "" {
			pvc["storageClassName"] = js.GetStorageClass().GetValue()
		}
		fileStore := map[string]interface{}{"pvc": pvc}
		if js.GetMaxFileStore() != "" {
			fileStore["maxSize"] = js.GetMaxFileStore()
		}
		jetstream["fileStore"] = fileStore
		if js.GetMemoryStoreMaxSize() != "" {
			jetstream["memoryStore"] = map[string]interface{}{
				"enabled": true,
				"maxSize": js.GetMemoryStoreMaxSize(),
			}
		}
	}
	config["jetstream"] = jetstream

	// Client-listener TLS: the chart mounts the Secret and renders the
	// cert/key paths; verify_clients adds mutual-TLS enforcement with
	// the CA bundle from the same Secret (cert-manager includes ca.crt).
	if tls := spec.GetTls(); tls != nil {
		tlsBlock := map[string]interface{}{
			"enabled":    true,
			"secretName": tls.GetSecretName().GetValue(),
		}
		if tls.GetVerifyClients() {
			tlsBlock["merge"] = map[string]interface{}{
				"verify":  true,
				"ca_file": fmt.Sprintf("%s/%s", vars.ClientTlsDir, vars.CaFileName),
			}
		}
		config["nats"] = map[string]interface{}{"tls": tlsBlock}
	}

	// WebSocket: the SERVER refuses a websocket listener without TLS
	// unless no_tls is explicit (verified in the server source at the
	// pin) — the typed arm is plain-in-cluster, so no_tls rides the
	// block's merge; websocket TLS is helm_values territory (a TLS
	// config, when merged there, wins over no_tls server-side).
	if ws := spec.GetWebsocket(); ws != nil && ws.GetEnabled() {
		port := int32(8080)
		if ws.Port != nil {
			port = ws.GetPort()
		}
		config["websocket"] = map[string]interface{}{
			"enabled": true,
			"port":    int(port),
			"merge":   map[string]interface{}{"no_tls": true},
		}
	}

	// MQTT (JetStream-backed — CEL enforces the pairing on the spec).
	if mqtt := spec.GetMqtt(); mqtt != nil && mqtt.GetEnabled() {
		port := int32(1883)
		if mqtt.Port != nil {
			port = mqtt.GetPort()
		}
		config["mqtt"] = map[string]interface{}{
			"enabled": true,
			"port":    int(port),
		}
	}

	// Leafnodes (the hub side).
	if leaf := spec.GetLeafnodes(); leaf != nil && leaf.GetEnabled() {
		port := int32(7422)
		if leaf.Port != nil {
			port = leaf.GetPort()
		}
		config["leafnodes"] = map[string]interface{}{
			"enabled": true,
			"port":    int(port),
		}
	}

	// Auth: authorization users XOR accounts (CEL-enforced), rendered
	// through the chart's config merge. Passwords are `$NATS_PW_<i>`
	// env references — never values.
	if locals.AuthEnabled {
		merge := map[string]interface{}{}
		auth := spec.GetAuth()
		if len(locals.FlatUsers) > 0 {
			users := make([]interface{}, 0, len(locals.FlatUsers))
			for i, u := range auth.GetUsers() {
				users = append(users, userEntry(u, locals.FlatUsers[i].EnvVar))
			}
			merge["authorization"] = map[string]interface{}{"users": users}
		}
		if len(locals.AccountUsers) > 0 {
			accounts := map[string]interface{}{}
			for ai, account := range auth.GetAccounts() {
				users := make([]interface{}, 0, len(account.GetUsers()))
				for ui, u := range account.GetUsers() {
					users = append(users, userEntry(u, locals.AccountUsers[ai][ui].EnvVar))
				}
				accountBlock := map[string]interface{}{"users": users}
				if account.GetJetStreamEnabled() {
					accountBlock["jetstream"] = "enabled"
				}
				accounts[account.GetName()] = accountBlock
			}
			merge["accounts"] = accounts
		}
		if auth.GetNoAuthUser() != "" {
			merge["no_auth_user"] = auth.GetNoAuthUser()
		}
		config["merge"] = merge
	}

	values["config"] = config

	// ---- the nats container ---------------------------------------------------
	container := map[string]interface{}{}
	if locals.AuthEnabled {
		env := map[string]interface{}{}
		allUsers := append([]authUser{}, locals.FlatUsers...)
		for _, users := range locals.AccountUsers {
			allUsers = append(allUsers, users...)
		}
		for _, u := range allUsers {
			env[u.EnvVar] = map[string]interface{}{
				"valueFrom": map[string]interface{}{
					"secretKeyRef": map[string]interface{}{
						"name": locals.AuthSecretName,
						"key":  u.Username,
					},
				},
			}
		}
		container["env"] = env
	}
	if resources := resourcesBlock(spec.GetResources()); resources != nil {
		container["resources"] = resources
	}
	if img := imageBlock(spec.GetImages().GetNats()); img != nil {
		container["image"] = img
	}
	if len(container) > 0 {
		values["container"] = container
	}

	// ---- sidecars and satellites -------------------------------------------------
	if img := imageBlock(spec.GetImages().GetReloader()); img != nil {
		values["reloader"] = map[string]interface{}{"image": img}
	}

	if metrics := spec.GetMetrics(); metrics != nil && metrics.GetExporterEnabled() {
		promExporter := map[string]interface{}{"enabled": true}
		if metrics.GetPodMonitorEnabled() {
			promExporter["podMonitor"] = map[string]interface{}{"enabled": true}
		}
		if img := imageBlock(spec.GetImages().GetExporter()); img != nil {
			promExporter["image"] = img
		}
		values["promExporter"] = promExporter
	}

	natsBoxEnabled := true
	if spec.NatsBoxEnabled != nil {
		natsBoxEnabled = spec.GetNatsBoxEnabled()
	}
	natsBox := map[string]interface{}{}
	if !natsBoxEnabled {
		natsBox["enabled"] = false
	}
	if img := imageBlock(spec.GetImages().GetNatsBox()); img != nil {
		natsBox["container"] = map[string]interface{}{"image": img}
	}
	if len(natsBox) > 0 {
		values["natsBox"] = natsBox
	}

	// ---- client Service exposure ---------------------------------------------------
	// The chart has no first-class type/annotations values — both ride
	// its Service merge patch.
	if service := spec.GetService(); service != nil {
		serviceMerge := map[string]interface{}{}
		switch service.GetType() {
		case kubernetesnatsv1.KubernetesNatsService_load_balancer:
			serviceMerge["spec"] = map[string]interface{}{"type": "LoadBalancer"}
		case kubernetesnatsv1.KubernetesNatsService_node_port:
			serviceMerge["spec"] = map[string]interface{}{"type": "NodePort"}
		}
		if len(service.GetAnnotations()) > 0 {
			serviceMerge["metadata"] = map[string]interface{}{
				"annotations": stringMapToInterface(service.GetAnnotations()),
			}
		}
		if len(serviceMerge) > 0 {
			values["service"] = map[string]interface{}{"merge": serviceMerge}
		}
	}

	// ---- scheduling (the pod-template merge patch) -------------------------------------
	if sched := spec.GetScheduling(); sched != nil {
		podSpec := map[string]interface{}{}
		if len(sched.GetNodeSelector()) > 0 {
			podSpec["nodeSelector"] = stringMapToInterface(sched.GetNodeSelector())
		}
		if len(sched.GetTolerations()) > 0 {
			podSpec["tolerations"] = tolerationsSlice(sched.GetTolerations())
		}
		if len(podSpec) > 0 {
			values["podTemplate"] = map[string]interface{}{
				"merge": map[string]interface{}{"spec": podSpec},
			}
		}
	}

	// ---- image pull secrets (the chart's global list) ----------------------------------
	pullSecretNames := []interface{}{}
	seen := map[string]bool{}
	for _, img := range []*kubernetesprovider.ContainerImage{
		spec.GetImages().GetNats(), spec.GetImages().GetReloader(),
		spec.GetImages().GetExporter(), spec.GetImages().GetNatsBox(),
	} {
		name := img.GetPullSecretName()
		if name != "" && !seen[name] {
			seen[name] = true
			pullSecretNames = append(pullSecretNames, name)
		}
	}
	if len(pullSecretNames) > 0 {
		values["global"] = map[string]interface{}{
			"image": map[string]interface{}{"pullSecretNames": pullSecretNames},
		}
	}

	// ---- escape hatch (merged LAST, helm -f semantics) ------------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// fullnameOverride re-pinned AFTER the merge — the one deliberate
	// exception to the escape hatch's last-word contract (twin of the
	// Terraform module's third values document). Every child name — and
	// the exported outputs built from them — derive from the fullname;
	// letting an override move it would break every output.
	values["fullnameOverride"] = locals.ReleaseName

	return values, nil
}

// userEntry renders one nats.conf user entry: the password is a
// `<< $NATS_PW_<i> >>` reference the chart unquotes into `$NATS_PW_<i>`,
// which the server resolves from the container environment (the
// secretKeyRef env wiring in the container block).
func userEntry(u *kubernetesnatsv1.KubernetesNatsUser, envVar string) map[string]interface{} {
	entry := map[string]interface{}{
		"user":     u.GetUsername(),
		"password": fmt.Sprintf("<< $%s >>", envVar),
	}
	if p := u.GetPermissions(); p != nil {
		permissions := map[string]interface{}{}
		if pub := allowDeny(p.GetPublishAllow(), p.GetPublishDeny()); pub != nil {
			permissions["publish"] = pub
		}
		if sub := allowDeny(p.GetSubscribeAllow(), p.GetSubscribeDeny()); sub != nil {
			permissions["subscribe"] = sub
		}
		if len(permissions) > 0 {
			entry["permissions"] = permissions
		}
	}
	return entry
}

// allowDeny renders a permissions {allow, deny} pair (nil when empty).
func allowDeny(allow, deny []string) map[string]interface{} {
	block := map[string]interface{}{}
	if len(allow) > 0 {
		block["allow"] = stringsToInterface(allow)
	}
	if len(deny) > 0 {
		block["deny"] = stringsToInterface(deny)
	}
	if len(block) == 0 {
		return nil
	}
	return block
}
