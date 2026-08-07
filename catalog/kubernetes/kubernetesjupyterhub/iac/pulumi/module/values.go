package module

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	kubernetesjupyterhubv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesjupyterhub/v1alpha1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// SECRET DISCIPLINE (load-bearing, and stricter here than on most Helm
// kinds): the chart embeds the ENTIRE rendered values document inside its
// own hub Secret (`values.yaml` key), so ANY credential placed in values
// becomes readable cluster state twice over. Nothing rendered here
// carries credential material — the external database password rides the
// module-owned hub.existingSecret, and the shared-password / OAuth client
// secrets reach JupyterHub as env vars (hub.extraEnv valueFrom) consumed
// by the module's hub.extraConfig python snippets.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values,
// yamlencode(re-pins)] and the provider merges the documents in exactly
// this order. Keep every typed mapping below in lockstep with the
// Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// ------------------------------ hub -----------------------------------
	hubSpec := spec.GetHub()
	hubBlock := map[string]interface{}{}

	// Database arm — the chart's own hub.db vocabulary. External arms
	// carry the CREDENTIAL-FREE url; the password rides the module's
	// existing-secret (chart truth: the hub exports PGPASSWORD/MYSQL_PWD
	// from the `hub.db.password` secret key at startup).
	dbBlock := map[string]interface{}{"type": locals.DbType}
	if locals.DbType == "sqlite-pvc" {
		pvcBlock := map[string]interface{}{}
		if sqlite := hubSpec.GetDatabase().GetSqlitePvc(); sqlite != nil {
			if sqlite.GetStorageSize() != "" {
				pvcBlock["storage"] = sqlite.GetStorageSize()
			}
			if sqlite.GetStorageClass() != "" {
				pvcBlock["storageClassName"] = sqlite.GetStorageClass()
			}
		}
		if len(pvcBlock) > 0 {
			dbBlock["pvc"] = pvcBlock
		}
	} else {
		dbBlock["url"] = locals.DbUrl
		hubBlock["existingSecret"] = locals.HubSecretName
	}
	hubBlock["db"] = dbBlock

	if hubSpec != nil && hubSpec.ConcurrentSpawnLimit != nil {
		hubBlock["concurrentSpawnLimit"] = int(hubSpec.GetConcurrentSpawnLimit())
	}
	if hubSpec != nil && hubSpec.ActiveServerLimit != nil {
		hubBlock["activeServerLimit"] = int(hubSpec.GetActiveServerLimit())
	}
	if hubSpec.GetAllowNamedServers() {
		hubBlock["allowNamedServers"] = true
		if hubSpec.NamedServerLimitPerUser != nil {
			hubBlock["namedServerLimitPerUser"] = int(hubSpec.GetNamedServerLimitPerUser())
		}
	}
	if hubSpec.GetShutdownOnLogout() {
		hubBlock["shutdownOnLogout"] = true
	}
	if resources := resourcesBlock(hubSpec.GetResources()); resources != nil {
		hubBlock["resources"] = resources
	}

	// Authentication — hub.config carries only PUBLIC identity settings
	// (class, client id, endpoints, rosters); every secret rides env
	// indirection through extraEnv + extraConfig.
	hubConfig, extraConfig, extraEnv := authBlocks(locals)
	hubBlock["config"] = hubConfig
	if len(extraConfig) > 0 {
		hubBlock["extraConfig"] = extraConfig
	}
	if len(extraEnv) > 0 {
		hubBlock["extraEnv"] = extraEnv
	}

	values["hub"] = hubBlock

	// ------------------------------ proxy ----------------------------------
	proxySpec := spec.GetProxy()
	serviceType := proxySpec.GetServiceType()
	if serviceType == "" {
		// DELIBERATE chart-default override: the chart ships
		// proxy.service.type LoadBalancer; this kind composes exposure
		// from first-class kinds, so the front door stays ClusterIP
		// unless the spec says otherwise.
		serviceType = "ClusterIP"
	}
	proxyService := map[string]interface{}{"type": serviceType}
	if len(proxySpec.GetServiceAnnotations()) > 0 {
		proxyService["annotations"] = toInterfaceMap(proxySpec.GetServiceAnnotations())
	}
	proxyBlock := map[string]interface{}{"service": proxyService}
	if resources := resourcesBlock(proxySpec.GetResources()); resources != nil {
		proxyBlock["chp"] = map[string]interface{}{"resources": resources}
	}
	values["proxy"] = proxyBlock

	// ---------------------------- singleuser --------------------------------
	singleUserBlock, err := singleUserValues(spec.GetSingleUser())
	if err != nil {
		return nil, err
	}
	if len(singleUserBlock) > 0 {
		values["singleuser"] = singleUserBlock
	}

	// ---------------------------- scheduling --------------------------------
	if scheduling := spec.GetScheduling(); scheduling != nil {
		schedulingBlock := map[string]interface{}{}
		if scheduling.UserSchedulerEnabled != nil {
			schedulingBlock["userScheduler"] = map[string]interface{}{
				"enabled": scheduling.GetUserSchedulerEnabled(),
			}
		}
		if scheduling.UserPlaceholderReplicas != nil {
			replicas := int(scheduling.GetUserPlaceholderReplicas())
			schedulingBlock["userPlaceholder"] = map[string]interface{}{
				"enabled":  replicas > 0,
				"replicas": replicas,
			}
			if replicas > 0 {
				// Placeholders only pre-warm capacity when real users
				// can EVICT them — that is the pod-priority machinery
				// (chart default off), so placeholders switch it on.
				schedulingBlock["podPriority"] = map[string]interface{}{"enabled": true}
			}
		}
		if len(schedulingBlock) > 0 {
			values["scheduling"] = schedulingBlock
		}
		if len(scheduling.GetCoreNodeSelector()) > 0 {
			coreSelector := toInterfaceMap(scheduling.GetCoreNodeSelector())
			hubBlock["nodeSelector"] = coreSelector
			if chpBlock, ok := proxyBlock["chp"].(map[string]interface{}); ok {
				chpBlock["nodeSelector"] = coreSelector
			} else {
				proxyBlock["chp"] = map[string]interface{}{"nodeSelector": coreSelector}
			}
		}
		if len(scheduling.GetUserNodeSelector()) > 0 {
			if suBlock, ok := values["singleuser"].(map[string]interface{}); ok {
				suBlock["nodeSelector"] = toInterfaceMap(scheduling.GetUserNodeSelector())
			} else {
				values["singleuser"] = map[string]interface{}{
					"nodeSelector": toInterfaceMap(scheduling.GetUserNodeSelector()),
				}
			}
		}
	}

	// ------------------------------ culling ---------------------------------
	if culling := spec.GetCulling(); culling != nil {
		cullBlock := map[string]interface{}{}
		if culling.Enabled != nil {
			cullBlock["enabled"] = culling.GetEnabled()
		}
		if culling.TimeoutSeconds != nil {
			cullBlock["timeout"] = int(culling.GetTimeoutSeconds())
		}
		if culling.EverySeconds != nil {
			cullBlock["every"] = int(culling.GetEverySeconds())
		}
		if culling.MaxAgeSeconds != nil {
			cullBlock["maxAge"] = int(culling.GetMaxAgeSeconds())
		}
		if culling.GetCullUsers() {
			cullBlock["users"] = true
		}
		if len(cullBlock) > 0 {
			values["cull"] = cullBlock
		}
	}

	// ----------------------------- pre-puller -------------------------------
	if prePuller := spec.GetPrePuller(); prePuller != nil {
		prePullerBlock := map[string]interface{}{}
		if prePuller.HookEnabled != nil {
			prePullerBlock["hook"] = map[string]interface{}{"enabled": prePuller.GetHookEnabled()}
		}
		if prePuller.ContinuousEnabled != nil {
			prePullerBlock["continuous"] = map[string]interface{}{"enabled": prePuller.GetContinuousEnabled()}
		}
		if len(prePullerBlock) > 0 {
			values["prePuller"] = prePullerBlock
		}
	}

	// --------------------------- network policies ---------------------------
	// One spec toggle drives the chart's three per-component
	// NetworkPolicy switches identically (they default true; only an
	// explicit false needs rendering).
	if spec.NetworkPolicyEnabled != nil && !spec.GetNetworkPolicyEnabled() {
		hubBlock["networkPolicy"] = map[string]interface{}{"enabled": false}
		if chpBlock, ok := proxyBlock["chp"].(map[string]interface{}); ok {
			chpBlock["networkPolicy"] = map[string]interface{}{"enabled": false}
		} else {
			proxyBlock["chp"] = map[string]interface{}{"networkPolicy": map[string]interface{}{"enabled": false}}
		}
		if suBlock, ok := values["singleuser"].(map[string]interface{}); ok {
			suBlock["networkPolicy"] = map[string]interface{}{"enabled": false}
		} else {
			values["singleuser"] = map[string]interface{}{
				"networkPolicy": map[string]interface{}{"enabled": false},
			}
		}
	}

	// ---- escape hatch (merged LAST, helm -f semantics) ----------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// Deliberate post-merge re-pin — the exception to the escape hatch's
	// last-word contract (twin of the Terraform module's third values
	// document): fullnameOverride stays "" — the chart-fixed bare names
	// (hub, proxy-public…) ARE the exported outputs; letting an
	// override prefix them would break every output and the per-
	// namespace singleton contract this kind documents.
	values["fullnameOverride"] = ""

	return values, nil
}

// authBlocks renders the authentication arm into the chart's hub.config,
// hub.extraConfig and hub.extraEnv values. The pattern is uniform across
// arms: PUBLIC settings ride hub.config; every SECRET rides an env var
// from a Secret (hub.extraEnv valueFrom) consumed by a python snippet
// (hub.extraConfig) — env values never land in the chart's readable
// values-embedding hub Secret.
func authBlocks(locals *Locals) (map[string]interface{}, map[string]interface{}, map[string]interface{}) {
	auth := locals.Spec.GetAuthentication()
	hubConfig := map[string]interface{}{}
	extraConfig := map[string]interface{}{}
	extraEnv := map[string]interface{}{}

	secretEnv := func(envVar, secretName, secretKey string) {
		extraEnv[envVar] = map[string]interface{}{
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"name": secretName,
					"key":  secretKey,
				},
			},
		}
	}

	// Roster settings apply to every authenticator class. JupyterHub 5
	// denies sign-in unless an allow rule matches — an empty roster
	// means "any authenticated identity", declared EXPLICITLY via
	// allow_all (never left to per-authenticator defaults).
	authenticatorConfig := map[string]interface{}{}
	if len(auth.GetAllowedUsers()) > 0 {
		authenticatorConfig["allowed_users"] = toInterfaceSlice(auth.GetAllowedUsers())
	} else {
		authenticatorConfig["allow_all"] = true
	}
	if len(auth.GetAdminUsers()) > 0 {
		authenticatorConfig["admin_users"] = toInterfaceSlice(auth.GetAdminUsers())
	}
	hubConfig["Authenticator"] = authenticatorConfig

	switch locals.AuthMethod {
	case "shared_password":
		// Full class path — entry-point shortnames are registration
		// details; the class path is unambiguous at any hub version.
		hubConfig["JupyterHub"] = map[string]interface{}{
			"authenticator_class": "jupyterhub.auth.DummyAuthenticator",
		}
		secretEnv(vars.SharedPasswordEnvVar, locals.SharedPasswordSecretName, locals.SharedPasswordSecretKey)
		extraConfig["plantonSharedPassword"] = fmt.Sprintf(
			"import os\nc.DummyAuthenticator.password = os.environ[%q]\n", vars.SharedPasswordEnvVar)

	case "native":
		native := auth.GetNative()
		hubConfig["JupyterHub"] = map[string]interface{}{
			"authenticator_class": "nativeauthenticator.NativeAuthenticator",
		}
		nativeConfig := map[string]interface{}{}
		if native.GetOpenSignup() {
			nativeConfig["open_signup"] = true
		}
		if native.MinimumPasswordLength != nil {
			nativeConfig["minimum_password_length"] = int(native.GetMinimumPasswordLength())
		}
		if len(nativeConfig) > 0 {
			hubConfig["NativeAuthenticator"] = nativeConfig
		}

	case "github":
		github := auth.GetGithub()
		hubConfig["JupyterHub"] = map[string]interface{}{
			"authenticator_class": "oauthenticator.github.GitHubOAuthenticator",
		}
		githubConfig := map[string]interface{}{
			"client_id":          github.GetClientId(),
			"oauth_callback_url": github.GetOauthCallbackUrl(),
		}
		if len(github.GetAllowedOrganizations()) > 0 {
			githubConfig["allowed_organizations"] = toInterfaceSlice(github.GetAllowedOrganizations())
			// Org/team membership checks need the read:org scope — the
			// authenticator cannot see memberships without it.
			githubConfig["scope"] = []interface{}{"read:org"}
		}
		hubConfig["GitHubOAuthenticator"] = githubConfig
		secretEnv(vars.OauthClientSecretEnvVar, locals.OauthClientSecretSecretName, locals.OauthClientSecretSecretKey)
		extraConfig["plantonOauthClientSecret"] = fmt.Sprintf(
			"import os\nc.GitHubOAuthenticator.client_secret = os.environ[%q]\n", vars.OauthClientSecretEnvVar)

	case "google":
		google := auth.GetGoogle()
		hubConfig["JupyterHub"] = map[string]interface{}{
			"authenticator_class": "oauthenticator.google.GoogleOAuthenticator",
		}
		googleConfig := map[string]interface{}{
			"client_id":          google.GetClientId(),
			"oauth_callback_url": google.GetOauthCallbackUrl(),
		}
		if len(google.GetHostedDomains()) > 0 {
			googleConfig["hosted_domain"] = toInterfaceSlice(google.GetHostedDomains())
		}
		hubConfig["GoogleOAuthenticator"] = googleConfig
		secretEnv(vars.OauthClientSecretEnvVar, locals.OauthClientSecretSecretName, locals.OauthClientSecretSecretKey)
		extraConfig["plantonOauthClientSecret"] = fmt.Sprintf(
			"import os\nc.GoogleOAuthenticator.client_secret = os.environ[%q]\n", vars.OauthClientSecretEnvVar)

	case "oidc":
		oidc := auth.GetOidc()
		hubConfig["JupyterHub"] = map[string]interface{}{
			"authenticator_class": "oauthenticator.generic.GenericOAuthenticator",
		}
		scopes := oidc.GetScopes()
		if len(scopes) == 0 {
			scopes = []string{"openid", "email", "profile"}
		}
		usernameClaim := oidc.GetUsernameClaim()
		if usernameClaim == "" {
			usernameClaim = "preferred_username"
		}
		loginService := oidc.GetLoginService()
		if loginService == "" {
			loginService = "OIDC"
		}
		hubConfig["GenericOAuthenticator"] = map[string]interface{}{
			"client_id":          oidc.GetClientId(),
			"oauth_callback_url": oidc.GetOauthCallbackUrl(),
			"authorize_url":      oidc.GetAuthorizeUrl(),
			"token_url":          oidc.GetTokenUrl(),
			"userdata_url":       oidc.GetUserdataUrl(),
			"scope":              toInterfaceSlice(scopes),
			"username_claim":     usernameClaim,
			"login_service":      loginService,
		}
		secretEnv(vars.OauthClientSecretEnvVar, locals.OauthClientSecretSecretName, locals.OauthClientSecretSecretKey)
		extraConfig["plantonOauthClientSecret"] = fmt.Sprintf(
			"import os\nc.GenericOAuthenticator.client_secret = os.environ[%q]\n", vars.OauthClientSecretEnvVar)
	}

	return hubConfig, extraConfig, extraEnv
}

// singleUserValues renders the per-user server configuration.
func singleUserValues(singleUser *kubernetesjupyterhubv1alpha1.KubernetesJupyterHubSingleUser) (map[string]interface{}, error) {
	block := map[string]interface{}{}
	if singleUser == nil {
		return block, nil
	}

	if image := singleUser.GetImage(); image != nil {
		block["image"] = map[string]interface{}{
			"name": image.GetRepository(),
			"tag":  image.GetTag(),
		}
	}

	memory := map[string]interface{}{}
	if singleUser.GetMemoryGuarantee() != "" {
		memory["guarantee"] = singleUser.GetMemoryGuarantee()
	}
	if singleUser.GetMemoryLimit() != "" {
		memory["limit"] = singleUser.GetMemoryLimit()
	}
	if len(memory) > 0 {
		block["memory"] = memory
	}

	cpu := map[string]interface{}{}
	if singleUser.GetCpuGuarantee() != "" {
		parsed, err := strconv.ParseFloat(singleUser.GetCpuGuarantee(), 64)
		if err != nil {
			return nil, errors.Wrapf(err, "single_user.cpu_guarantee %q is not a number — the chart takes CPU as a number (e.g. \"0.5\", \"2\")", singleUser.GetCpuGuarantee())
		}
		cpu["guarantee"] = parsed
	}
	if singleUser.GetCpuLimit() != "" {
		parsed, err := strconv.ParseFloat(singleUser.GetCpuLimit(), 64)
		if err != nil {
			return nil, errors.Wrapf(err, "single_user.cpu_limit %q is not a number — the chart takes CPU as a number (e.g. \"0.5\", \"2\")", singleUser.GetCpuLimit())
		}
		cpu["limit"] = parsed
	}
	if len(cpu) > 0 {
		block["cpu"] = cpu
	}

	if storage := singleUser.GetStorage(); storage != nil {
		storageBlock := map[string]interface{}{}
		switch {
		case storage.GetDynamic() != nil:
			dynamic := storage.GetDynamic()
			storageBlock["type"] = "dynamic"
			if dynamic.GetCapacity() != "" {
				storageBlock["capacity"] = dynamic.GetCapacity()
			}
			if dynamic.GetStorageClass() != "" {
				storageBlock["dynamic"] = map[string]interface{}{"storageClass": dynamic.GetStorageClass()}
			}
		case storage.GetStatic() != nil:
			static := storage.GetStatic()
			subPath := static.GetSubPath()
			if subPath == "" {
				subPath = "{username}"
			}
			storageBlock["type"] = "static"
			storageBlock["static"] = map[string]interface{}{
				"pvcName": static.GetPvcName(),
				"subPath": subPath,
			}
		case storage.GetNone() != nil:
			storageBlock["type"] = "none"
		}
		if len(storageBlock) > 0 {
			block["storage"] = storageBlock
		}
	}

	if singleUser.GetDefaultUrl() != "" {
		block["defaultUrl"] = singleUser.GetDefaultUrl()
	}
	if singleUser.StartTimeoutSeconds != nil {
		block["startTimeout"] = int(singleUser.GetStartTimeoutSeconds())
	}
	if len(singleUser.GetExtraEnv()) > 0 {
		block["extraEnv"] = toInterfaceMap(singleUser.GetExtraEnv())
	}

	if len(singleUser.GetProfiles()) > 0 {
		profiles := make([]interface{}, 0, len(singleUser.GetProfiles()))
		for _, profile := range singleUser.GetProfiles() {
			entry := map[string]interface{}{
				"display_name": profile.GetDisplayName(),
			}
			if profile.GetDescription() != "" {
				entry["description"] = profile.GetDescription()
			}
			if profile.GetDefault() {
				entry["default"] = true
			}
			// kubespawner_override keys are KubeSpawner trait names —
			// the chart passes profileList through to the spawner.
			override := map[string]interface{}{}
			if image := profile.GetImage(); image != nil {
				override["image"] = fmt.Sprintf("%s:%s", image.GetRepository(), image.GetTag())
			}
			if profile.GetMemoryGuarantee() != "" {
				override["mem_guarantee"] = profile.GetMemoryGuarantee()
			}
			if profile.GetMemoryLimit() != "" {
				override["mem_limit"] = profile.GetMemoryLimit()
			}
			if profile.GetCpuGuarantee() != "" {
				parsed, err := strconv.ParseFloat(profile.GetCpuGuarantee(), 64)
				if err != nil {
					return nil, errors.Wrapf(err, "profile %q cpu_guarantee %q is not a number", profile.GetDisplayName(), profile.GetCpuGuarantee())
				}
				override["cpu_guarantee"] = parsed
			}
			if profile.GetCpuLimit() != "" {
				parsed, err := strconv.ParseFloat(profile.GetCpuLimit(), 64)
				if err != nil {
					return nil, errors.Wrapf(err, "profile %q cpu_limit %q is not a number", profile.GetDisplayName(), profile.GetCpuLimit())
				}
				override["cpu_limit"] = parsed
			}
			if len(override) > 0 {
				entry["kubespawner_override"] = override
			}
			profiles = append(profiles, entry)
		}
		block["profileList"] = profiles
	}

	return block, nil
}
