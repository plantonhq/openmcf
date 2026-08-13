# KubernetesStatefulSet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

**KubernetesStatefulSetSpec** deploys a stateful application on a Kubernetes
cluster as an apps/v1 StatefulSet: every replica gets a stable name
(<name>-0, <name>-1, ...), stable per-replica DNS through a headless Service the
module derives from the resource name, and its own PersistentVolumeClaim stamped
from `volume_claim_templates`. This is the kind for databases, message brokers,
and consensus systems — anything where replicas are NOT interchangeable. For
stateless services use KubernetesDeployment; for run-to-completion work use
KubernetesJob/KubernetesCronJob.

External exposure is composed, never embedded: this kind exports its governing
Service and selector labels, and first-class exposure kinds (KubernetesIngress,
KubernetesHttpRoute and the other Gateway API route kinds, with certificates)
reference them — so every piece of exposure infrastructure is a visible node
in the resource graph.

DEPLOY-TARGET CONTRACT: this kind is a Service Hub deployment target. Deployment
pipelines inject the freshly built artifact at exactly one stable path —
`spec.container.app.image` (repo + tag). That path is part of the kind's public
contract and must survive any future spec evolution.

## Example

```yaml
# Local development / offline-proof manifest exercising the full spec surface,
# including arms the live E2E scenarios exclude — the offline tofu plan and
# pulumi preview proofs run against this file.
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesStatefulSet
metadata:
  name: hack-statefulset
  id: hack-statefulset-id
  org: hack-org
  env: hack-env
spec:
  namespace:
    value: hack-statefulset-ns
  create_namespace: true
  container:
    app:
      image:
        repo: nginx
        tag: "1.27-alpine"
      image_pull_policy: IfNotPresent
      working_dir: /usr/share/nginx/html
      resources:
        requests:
          cpu: "10m"
          memory: "32Mi"
        limits:
          cpu: "100m"
          memory: "64Mi"
      ports:
        - name: http
          container_port: 80
          network_protocol: TCP
          app_protocol: http
          service_port: 80
        - name: metrics
          container_port: 9113
          network_protocol: TCP
          service_port: 9113
      env:
        variables:
          - name: LOG_LEVEL
            value: info
          - name: POD_NAME
            field_ref:
              field_path: metadata.name
          - name: CPU_LIMIT
            resource_field_ref:
              resource: limits.cpu
        secrets:
          - name: DB_PASSWORD
            value: hack-password
        env_from:
          - config_map_ref:
              name: hack-env-config
              optional: true
      liveness_probe:
        http_get:
          path: /
          port_number: 80
        initial_delay_seconds: 5
        period_seconds: 10
      readiness_probe:
        tcp_socket:
          port_number: 80
        period_seconds: 5
      startup_probe:
        http_get:
          path: /
          port_number: 80
        failure_threshold: 30
        period_seconds: 2
      volume_mounts:
        - name: data
          mount_path: /var/lib/data
          pvc:
            claim_name: data
        - name: wal
          mount_path: /var/lib/wal
          pvc:
            claim_name: wal
        - name: scratch
          mount_path: /scratch
          empty_dir:
            medium: Memory
            size_limit: 32Mi
      lifecycle:
        post_start:
          exec:
            command: ["/bin/sh", "-c", "true"]
        pre_stop:
          exec:
            command: ["/bin/sleep", "5"]
      security_context:
        run_as_non_root: true
        read_only_root_filesystem: false
        allow_privilege_escalation: false
        capabilities:
          drop: ["ALL"]
          add: ["NET_BIND_SERVICE"]
        seccomp_profile:
          type: RuntimeDefault
    sidecars:
      - name: metrics-exporter
        image:
          repo: busybox
          tag: "1.36"
        command: ["/bin/sh", "-c", "sleep infinity"]
        resources:
          requests:
            cpu: "5m"
            memory: "16Mi"
          limits:
            cpu: "50m"
            memory: "32Mi"
        volume_mounts:
          - name: data
            mount_path: /data
            read_only: true
  pod:
    service_account:
      value: hack-service-account
    automount_service_account_token: false
    image_pull_secrets:
      - value: hack-registry-secret
    init_containers:
      - name: init-datadir
        image:
          repo: busybox
          tag: "1.36"
        command: ["/bin/sh", "-c", "echo init done"]
    labels:
      team: hack
    annotations:
      example.org/owner: hack
    scheduling:
      node_selector:
        kubernetes.io/os: linux
      tolerations:
        - key: dedicated
          operator: Equal
          value: hack
          effect: NoSchedule
      node_affinity:
        preferred:
          - weight: 50
            term:
              match_expressions:
                - key: kubernetes.io/arch
                  operator: In
                  values: ["arm64", "amd64"]
      pod_anti_affinity:
        required:
          - match_labels:
              app: hack-statefulset
            topology_key: kubernetes.io/hostname
      topology_spread_constraints:
        - max_skew: 1
          topology_key: topology.kubernetes.io/zone
          when_unsatisfiable: ScheduleAnyway
    security_context:
      run_as_user: 101
      run_as_group: 101
      run_as_non_root: true
      fs_group: 101
      fs_group_change_policy: OnRootMismatch
      sysctls:
        - name: net.core.somaxconn
          value: "1024"
    termination_grace_period_seconds: 120
    dns_policy: ClusterFirst
    dns_config:
      options:
        - name: ndots
          value: "2"
    host_aliases:
      - ip: 203.0.113.10
        hostnames: ["legacy.internal.example"]
    priority_class_name: ""
    runtime_class_name: ""
  availability:
    replicas: 3
    pod_disruption_budget:
      enabled: true
      min_available: "2"
    min_ready_seconds: 5
    revision_history_limit: 5
  volume_claim_templates:
    - name: data
      storage_class: standard
      size: "1Gi"
      access_modes:
        - ReadWriteOnce
      volume_mode: Filesystem
    - name: wal
      size: "500Mi"
      access_modes:
        - ReadWriteOnce
  update_strategy:
    type: RollingUpdate
    partition: 1
  pod_management_policy: Parallel
  pvc_retention_policy:
    when_deleted: Delete
    when_scaled: Retain
  # ordinals is deliberately absent here: the offline plan proof runs through
  # BOTH engines from this one manifest, and the Terraform provider cannot
  # express ordinals (see the PARITY-EXCEPTION in the modules).
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.container` | `KubernetesStatefulSetContainer` | yes |  |  |
| `spec.container.app` | `WorkloadContainer` | yes |  |  |
| `spec.container.app.name` | `string` |  |  |  |
| `spec.container.app.image` | `ContainerImage` | yes |  |  |
| `spec.container.app.image.repo` | `string` |  |  |  |
| `spec.container.app.image.tag` | `string` |  |  |  |
| `spec.container.app.image.pullSecretName` | `string` |  |  |  |
| `spec.container.app.imagePullPolicy` | `string` |  |  |  |
| `spec.container.app.command` | `[]string` |  |  |  |
| `spec.container.app.args` | `[]string` |  |  |  |
| `spec.container.app.workingDir` | `string` |  |  |  |
| `spec.container.app.ports` | `[]WorkloadContainerPort` |  |  |  |
| `spec.container.app.ports[].name` | `string` | yes |  |  |
| `spec.container.app.ports[].containerPort` | `int32` | yes |  |  |
| `spec.container.app.ports[].networkProtocol` | `string` |  |  |  |
| `spec.container.app.ports[].appProtocol` | `string` |  |  |  |
| `spec.container.app.ports[].servicePort` | `int32` |  |  |  |
| `spec.container.app.ports[].hostPort` | `int32` |  |  |  |
| `spec.container.app.env` | `ContainerEnv` |  |  |  |
| `spec.container.app.env.variables` | `[]EnvVar` |  |  |  |
| `spec.container.app.env.variables[].name` | `string` | yes |  |  |
| `spec.container.app.env.variables[].value` | `string` |  |  |  |
| `spec.container.app.env.variables[].valueFrom` | `ValueFromRef` |  |  |  |
| `spec.container.app.env.variables[].valueFrom.kind` | `enum` |  |  |  |
| `spec.container.app.env.variables[].valueFrom.env` | `string` |  |  |  |
| `spec.container.app.env.variables[].valueFrom.name` | `string` | yes |  |  |
| `spec.container.app.env.variables[].valueFrom.fieldPath` | `string` |  |  |  |
| `spec.container.app.env.variables[].configMapKeyRef` | `ConfigMapKeyRef` |  |  |  |
| `spec.container.app.env.variables[].configMapKeyRef.name` | `string` | yes |  |  |
| `spec.container.app.env.variables[].configMapKeyRef.key` | `string` | yes |  |  |
| `spec.container.app.env.variables[].configMapKeyRef.optional` | `bool` |  |  |  |
| `spec.container.app.env.variables[].fieldRef` | `ObjectFieldRef` |  |  |  |
| `spec.container.app.env.variables[].fieldRef.apiVersion` | `string` |  |  |  |
| `spec.container.app.env.variables[].fieldRef.fieldPath` | `string` | yes |  |  |
| `spec.container.app.env.variables[].resourceFieldRef` | `ResourceFieldRef` |  |  |  |
| `spec.container.app.env.variables[].resourceFieldRef.containerName` | `string` |  |  |  |
| `spec.container.app.env.variables[].resourceFieldRef.resource` | `string` | yes |  |  |
| `spec.container.app.env.variables[].resourceFieldRef.divisor` | `string` |  |  |  |
| `spec.container.app.env.secrets` | `[]SecretEnvVar` |  |  |  |
| `spec.container.app.env.secrets[].name` | `string` | yes |  |  |
| `spec.container.app.env.secrets[].value` | `string` |  |  |  |
| `spec.container.app.env.secrets[].secretRef` | `KubernetesSecretKeyRef` |  |  |  |
| `spec.container.app.env.secrets[].secretRef.namespace` | `string` |  |  |  |
| `spec.container.app.env.secrets[].secretRef.name` | `string` | yes |  |  |
| `spec.container.app.env.secrets[].secretRef.key` | `string` | yes |  |  |
| `spec.container.app.env.secrets[].secretRef.optional` | `bool` |  |  |  |
| `spec.container.app.env.secrets[].valueFrom` | `ValueFromRef` |  |  |  |
| `spec.container.app.env.secrets[].valueFrom.kind` | `enum` |  |  |  |
| `spec.container.app.env.secrets[].valueFrom.env` | `string` |  |  |  |
| `spec.container.app.env.secrets[].valueFrom.name` | `string` | yes |  |  |
| `spec.container.app.env.secrets[].valueFrom.fieldPath` | `string` |  |  |  |
| `spec.container.app.env.envFrom` | `[]EnvFromSource` |  |  |  |
| `spec.container.app.env.envFrom[].prefix` | `string` |  |  |  |
| `spec.container.app.env.envFrom[].configMapRef` | `ConfigMapRef` |  |  |  |
| `spec.container.app.env.envFrom[].configMapRef.name` | `string` | yes |  |  |
| `spec.container.app.env.envFrom[].configMapRef.optional` | `bool` |  |  |  |
| `spec.container.app.env.envFrom[].secretRef` | `SecretRef` |  |  |  |
| `spec.container.app.env.envFrom[].secretRef.name` | `string` | yes |  |  |
| `spec.container.app.env.envFrom[].secretRef.optional` | `bool` |  |  |  |
| `spec.container.app.resources` | `ContainerResources` |  |  |  |
| `spec.container.app.resources.limits` | `CpuMemory` |  |  |  |
| `spec.container.app.resources.limits.cpu` | `string` |  |  |  |
| `spec.container.app.resources.limits.memory` | `string` |  |  |  |
| `spec.container.app.resources.requests` | `CpuMemory` |  |  |  |
| `spec.container.app.resources.requests.cpu` | `string` |  |  |  |
| `spec.container.app.resources.requests.memory` | `string` |  |  |  |
| `spec.container.app.livenessProbe` | `Probe` |  |  |  |
| `spec.container.app.livenessProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.container.app.livenessProbe.periodSeconds` | `int32` |  |  |  |
| `spec.container.app.livenessProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.container.app.livenessProbe.successThreshold` | `int32` |  |  |  |
| `spec.container.app.livenessProbe.failureThreshold` | `int32` |  |  |  |
| `spec.container.app.livenessProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.app.livenessProbe.httpGet.path` | `string` |  |  |  |
| `spec.container.app.livenessProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.app.livenessProbe.httpGet.portName` | `string` |  |  |  |
| `spec.container.app.livenessProbe.httpGet.host` | `string` |  |  |  |
| `spec.container.app.livenessProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.container.app.livenessProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.app.livenessProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.app.livenessProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.app.livenessProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.container.app.livenessProbe.grpc.port` | `int32` |  |  |  |
| `spec.container.app.livenessProbe.grpc.service` | `string` |  |  |  |
| `spec.container.app.livenessProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.app.livenessProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.app.livenessProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.app.livenessProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.container.app.livenessProbe.exec` | `ExecAction` |  |  |  |
| `spec.container.app.livenessProbe.exec.command` | `[]string` |  |  |  |
| `spec.container.app.readinessProbe` | `Probe` |  |  |  |
| `spec.container.app.readinessProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.container.app.readinessProbe.periodSeconds` | `int32` |  |  |  |
| `spec.container.app.readinessProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.container.app.readinessProbe.successThreshold` | `int32` |  |  |  |
| `spec.container.app.readinessProbe.failureThreshold` | `int32` |  |  |  |
| `spec.container.app.readinessProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.app.readinessProbe.httpGet.path` | `string` |  |  |  |
| `spec.container.app.readinessProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.app.readinessProbe.httpGet.portName` | `string` |  |  |  |
| `spec.container.app.readinessProbe.httpGet.host` | `string` |  |  |  |
| `spec.container.app.readinessProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.container.app.readinessProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.app.readinessProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.app.readinessProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.app.readinessProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.container.app.readinessProbe.grpc.port` | `int32` |  |  |  |
| `spec.container.app.readinessProbe.grpc.service` | `string` |  |  |  |
| `spec.container.app.readinessProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.app.readinessProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.app.readinessProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.app.readinessProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.container.app.readinessProbe.exec` | `ExecAction` |  |  |  |
| `spec.container.app.readinessProbe.exec.command` | `[]string` |  |  |  |
| `spec.container.app.startupProbe` | `Probe` |  |  |  |
| `spec.container.app.startupProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.container.app.startupProbe.periodSeconds` | `int32` |  |  |  |
| `spec.container.app.startupProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.container.app.startupProbe.successThreshold` | `int32` |  |  |  |
| `spec.container.app.startupProbe.failureThreshold` | `int32` |  |  |  |
| `spec.container.app.startupProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.app.startupProbe.httpGet.path` | `string` |  |  |  |
| `spec.container.app.startupProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.app.startupProbe.httpGet.portName` | `string` |  |  |  |
| `spec.container.app.startupProbe.httpGet.host` | `string` |  |  |  |
| `spec.container.app.startupProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.container.app.startupProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.app.startupProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.app.startupProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.app.startupProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.container.app.startupProbe.grpc.port` | `int32` |  |  |  |
| `spec.container.app.startupProbe.grpc.service` | `string` |  |  |  |
| `spec.container.app.startupProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.app.startupProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.app.startupProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.app.startupProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.container.app.startupProbe.exec` | `ExecAction` |  |  |  |
| `spec.container.app.startupProbe.exec.command` | `[]string` |  |  |  |
| `spec.container.app.volumeMounts` | `[]VolumeMount` |  |  |  |
| `spec.container.app.volumeMounts[].name` | `string` | yes |  |  |
| `spec.container.app.volumeMounts[].mountPath` | `string` | yes |  |  |
| `spec.container.app.volumeMounts[].readOnly` | `bool` |  |  |  |
| `spec.container.app.volumeMounts[].subPath` | `string` |  |  |  |
| `spec.container.app.volumeMounts[].configMap` | `ConfigMapVolumeSource` |  |  |  |
| `spec.container.app.volumeMounts[].configMap.name` | `string` | yes |  |  |
| `spec.container.app.volumeMounts[].configMap.key` | `string` |  |  |  |
| `spec.container.app.volumeMounts[].configMap.path` | `string` |  |  |  |
| `spec.container.app.volumeMounts[].configMap.defaultMode` | `int32` |  |  |  |
| `spec.container.app.volumeMounts[].secret` | `SecretVolumeSource` |  |  |  |
| `spec.container.app.volumeMounts[].secret.name` | `string` | yes |  |  |
| `spec.container.app.volumeMounts[].secret.key` | `string` |  |  |  |
| `spec.container.app.volumeMounts[].secret.path` | `string` |  |  |  |
| `spec.container.app.volumeMounts[].secret.defaultMode` | `int32` |  |  |  |
| `spec.container.app.volumeMounts[].hostPath` | `HostPathVolumeSource` |  |  |  |
| `spec.container.app.volumeMounts[].hostPath.path` | `string` | yes |  |  |
| `spec.container.app.volumeMounts[].hostPath.type` | `string` |  |  |  |
| `spec.container.app.volumeMounts[].emptyDir` | `EmptyDirVolumeSource` |  |  |  |
| `spec.container.app.volumeMounts[].emptyDir.medium` | `string` |  |  |  |
| `spec.container.app.volumeMounts[].emptyDir.sizeLimit` | `string` |  |  |  |
| `spec.container.app.volumeMounts[].pvc` | `PvcVolumeSource` |  |  |  |
| `spec.container.app.volumeMounts[].pvc.claimName` | `string` | yes |  |  |
| `spec.container.app.volumeMounts[].pvc.readOnly` | `bool` |  |  |  |
| `spec.container.app.lifecycle` | `WorkloadContainerLifecycle` |  |  |  |
| `spec.container.app.lifecycle.postStart` | `WorkloadLifecycleHandler` |  |  |  |
| `spec.container.app.lifecycle.postStart.exec` | `ExecAction` |  |  |  |
| `spec.container.app.lifecycle.postStart.exec.command` | `[]string` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet.path` | `string` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet.portName` | `string` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet.host` | `string` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet.scheme` | `string` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.app.lifecycle.postStart.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.app.lifecycle.postStart.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.app.lifecycle.postStart.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.app.lifecycle.postStart.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.app.lifecycle.postStart.tcpSocket.host` | `string` |  |  |  |
| `spec.container.app.lifecycle.postStart.sleep` | `SleepAction` |  |  |  |
| `spec.container.app.lifecycle.postStart.sleep.seconds` | `int64` |  |  |  |
| `spec.container.app.lifecycle.preStop` | `WorkloadLifecycleHandler` |  |  |  |
| `spec.container.app.lifecycle.preStop.exec` | `ExecAction` |  |  |  |
| `spec.container.app.lifecycle.preStop.exec.command` | `[]string` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet.path` | `string` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet.portName` | `string` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet.host` | `string` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet.scheme` | `string` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.app.lifecycle.preStop.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.app.lifecycle.preStop.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.app.lifecycle.preStop.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.app.lifecycle.preStop.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.app.lifecycle.preStop.tcpSocket.host` | `string` |  |  |  |
| `spec.container.app.lifecycle.preStop.sleep` | `SleepAction` |  |  |  |
| `spec.container.app.lifecycle.preStop.sleep.seconds` | `int64` |  |  |  |
| `spec.container.app.securityContext` | `WorkloadContainerSecurityContext` |  |  |  |
| `spec.container.app.securityContext.privileged` | `bool` |  |  |  |
| `spec.container.app.securityContext.runAsUser` | `int64` |  |  |  |
| `spec.container.app.securityContext.runAsGroup` | `int64` |  |  |  |
| `spec.container.app.securityContext.runAsNonRoot` | `bool` |  |  |  |
| `spec.container.app.securityContext.readOnlyRootFilesystem` | `bool` |  |  |  |
| `spec.container.app.securityContext.allowPrivilegeEscalation` | `bool` |  |  |  |
| `spec.container.app.securityContext.capabilities` | `WorkloadCapabilities` |  |  |  |
| `spec.container.app.securityContext.capabilities.add` | `[]string` |  |  |  |
| `spec.container.app.securityContext.capabilities.drop` | `[]string` |  |  |  |
| `spec.container.app.securityContext.seccompProfile` | `WorkloadSeccompProfile` |  |  |  |
| `spec.container.app.securityContext.seccompProfile.type` | `string` | yes |  |  |
| `spec.container.app.securityContext.seccompProfile.localhostProfile` | `string` |  |  |  |
| `spec.container.sidecars` | `[]WorkloadContainer` |  |  |  |
| `spec.container.sidecars[].name` | `string` |  |  |  |
| `spec.container.sidecars[].image` | `ContainerImage` | yes |  |  |
| `spec.container.sidecars[].image.repo` | `string` |  |  |  |
| `spec.container.sidecars[].image.tag` | `string` |  |  |  |
| `spec.container.sidecars[].image.pullSecretName` | `string` |  |  |  |
| `spec.container.sidecars[].imagePullPolicy` | `string` |  |  |  |
| `spec.container.sidecars[].command` | `[]string` |  |  |  |
| `spec.container.sidecars[].args` | `[]string` |  |  |  |
| `spec.container.sidecars[].workingDir` | `string` |  |  |  |
| `spec.container.sidecars[].ports` | `[]WorkloadContainerPort` |  |  |  |
| `spec.container.sidecars[].ports[].name` | `string` | yes |  |  |
| `spec.container.sidecars[].ports[].containerPort` | `int32` | yes |  |  |
| `spec.container.sidecars[].ports[].networkProtocol` | `string` |  |  |  |
| `spec.container.sidecars[].ports[].appProtocol` | `string` |  |  |  |
| `spec.container.sidecars[].ports[].servicePort` | `int32` |  |  |  |
| `spec.container.sidecars[].ports[].hostPort` | `int32` |  |  |  |
| `spec.container.sidecars[].env` | `ContainerEnv` |  |  |  |
| `spec.container.sidecars[].env.variables` | `[]EnvVar` |  |  |  |
| `spec.container.sidecars[].env.variables[].name` | `string` | yes |  |  |
| `spec.container.sidecars[].env.variables[].value` | `string` |  |  |  |
| `spec.container.sidecars[].env.variables[].valueFrom` | `ValueFromRef` |  |  |  |
| `spec.container.sidecars[].env.variables[].valueFrom.kind` | `enum` |  |  |  |
| `spec.container.sidecars[].env.variables[].valueFrom.env` | `string` |  |  |  |
| `spec.container.sidecars[].env.variables[].valueFrom.name` | `string` | yes |  |  |
| `spec.container.sidecars[].env.variables[].valueFrom.fieldPath` | `string` |  |  |  |
| `spec.container.sidecars[].env.variables[].configMapKeyRef` | `ConfigMapKeyRef` |  |  |  |
| `spec.container.sidecars[].env.variables[].configMapKeyRef.name` | `string` | yes |  |  |
| `spec.container.sidecars[].env.variables[].configMapKeyRef.key` | `string` | yes |  |  |
| `spec.container.sidecars[].env.variables[].configMapKeyRef.optional` | `bool` |  |  |  |
| `spec.container.sidecars[].env.variables[].fieldRef` | `ObjectFieldRef` |  |  |  |
| `spec.container.sidecars[].env.variables[].fieldRef.apiVersion` | `string` |  |  |  |
| `spec.container.sidecars[].env.variables[].fieldRef.fieldPath` | `string` | yes |  |  |
| `spec.container.sidecars[].env.variables[].resourceFieldRef` | `ResourceFieldRef` |  |  |  |
| `spec.container.sidecars[].env.variables[].resourceFieldRef.containerName` | `string` |  |  |  |
| `spec.container.sidecars[].env.variables[].resourceFieldRef.resource` | `string` | yes |  |  |
| `spec.container.sidecars[].env.variables[].resourceFieldRef.divisor` | `string` |  |  |  |
| `spec.container.sidecars[].env.secrets` | `[]SecretEnvVar` |  |  |  |
| `spec.container.sidecars[].env.secrets[].name` | `string` | yes |  |  |
| `spec.container.sidecars[].env.secrets[].value` | `string` |  |  |  |
| `spec.container.sidecars[].env.secrets[].secretRef` | `KubernetesSecretKeyRef` |  |  |  |
| `spec.container.sidecars[].env.secrets[].secretRef.namespace` | `string` |  |  |  |
| `spec.container.sidecars[].env.secrets[].secretRef.name` | `string` | yes |  |  |
| `spec.container.sidecars[].env.secrets[].secretRef.key` | `string` | yes |  |  |
| `spec.container.sidecars[].env.secrets[].secretRef.optional` | `bool` |  |  |  |
| `spec.container.sidecars[].env.secrets[].valueFrom` | `ValueFromRef` |  |  |  |
| `spec.container.sidecars[].env.secrets[].valueFrom.kind` | `enum` |  |  |  |
| `spec.container.sidecars[].env.secrets[].valueFrom.env` | `string` |  |  |  |
| `spec.container.sidecars[].env.secrets[].valueFrom.name` | `string` | yes |  |  |
| `spec.container.sidecars[].env.secrets[].valueFrom.fieldPath` | `string` |  |  |  |
| `spec.container.sidecars[].env.envFrom` | `[]EnvFromSource` |  |  |  |
| `spec.container.sidecars[].env.envFrom[].prefix` | `string` |  |  |  |
| `spec.container.sidecars[].env.envFrom[].configMapRef` | `ConfigMapRef` |  |  |  |
| `spec.container.sidecars[].env.envFrom[].configMapRef.name` | `string` | yes |  |  |
| `spec.container.sidecars[].env.envFrom[].configMapRef.optional` | `bool` |  |  |  |
| `spec.container.sidecars[].env.envFrom[].secretRef` | `SecretRef` |  |  |  |
| `spec.container.sidecars[].env.envFrom[].secretRef.name` | `string` | yes |  |  |
| `spec.container.sidecars[].env.envFrom[].secretRef.optional` | `bool` |  |  |  |
| `spec.container.sidecars[].resources` | `ContainerResources` |  |  |  |
| `spec.container.sidecars[].resources.limits` | `CpuMemory` |  |  |  |
| `spec.container.sidecars[].resources.limits.cpu` | `string` |  |  |  |
| `spec.container.sidecars[].resources.limits.memory` | `string` |  |  |  |
| `spec.container.sidecars[].resources.requests` | `CpuMemory` |  |  |  |
| `spec.container.sidecars[].resources.requests.cpu` | `string` |  |  |  |
| `spec.container.sidecars[].resources.requests.memory` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe` | `Probe` |  |  |  |
| `spec.container.sidecars[].livenessProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].livenessProbe.periodSeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].livenessProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].livenessProbe.successThreshold` | `int32` |  |  |  |
| `spec.container.sidecars[].livenessProbe.failureThreshold` | `int32` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet.path` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet.portName` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet.host` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.container.sidecars[].livenessProbe.grpc.port` | `int32` |  |  |  |
| `spec.container.sidecars[].livenessProbe.grpc.service` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.sidecars[].livenessProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].livenessProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.container.sidecars[].livenessProbe.exec` | `ExecAction` |  |  |  |
| `spec.container.sidecars[].livenessProbe.exec.command` | `[]string` |  |  |  |
| `spec.container.sidecars[].readinessProbe` | `Probe` |  |  |  |
| `spec.container.sidecars[].readinessProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].readinessProbe.periodSeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].readinessProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].readinessProbe.successThreshold` | `int32` |  |  |  |
| `spec.container.sidecars[].readinessProbe.failureThreshold` | `int32` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet.path` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet.portName` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet.host` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.container.sidecars[].readinessProbe.grpc.port` | `int32` |  |  |  |
| `spec.container.sidecars[].readinessProbe.grpc.service` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.sidecars[].readinessProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].readinessProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.container.sidecars[].readinessProbe.exec` | `ExecAction` |  |  |  |
| `spec.container.sidecars[].readinessProbe.exec.command` | `[]string` |  |  |  |
| `spec.container.sidecars[].startupProbe` | `Probe` |  |  |  |
| `spec.container.sidecars[].startupProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].startupProbe.periodSeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].startupProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.container.sidecars[].startupProbe.successThreshold` | `int32` |  |  |  |
| `spec.container.sidecars[].startupProbe.failureThreshold` | `int32` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet.path` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet.portName` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet.host` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.container.sidecars[].startupProbe.grpc.port` | `int32` |  |  |  |
| `spec.container.sidecars[].startupProbe.grpc.service` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.sidecars[].startupProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].startupProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.container.sidecars[].startupProbe.exec` | `ExecAction` |  |  |  |
| `spec.container.sidecars[].startupProbe.exec.command` | `[]string` |  |  |  |
| `spec.container.sidecars[].volumeMounts` | `[]VolumeMount` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].name` | `string` | yes |  |  |
| `spec.container.sidecars[].volumeMounts[].mountPath` | `string` | yes |  |  |
| `spec.container.sidecars[].volumeMounts[].readOnly` | `bool` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].subPath` | `string` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].configMap` | `ConfigMapVolumeSource` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].configMap.name` | `string` | yes |  |  |
| `spec.container.sidecars[].volumeMounts[].configMap.key` | `string` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].configMap.path` | `string` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].configMap.defaultMode` | `int32` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].secret` | `SecretVolumeSource` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].secret.name` | `string` | yes |  |  |
| `spec.container.sidecars[].volumeMounts[].secret.key` | `string` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].secret.path` | `string` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].secret.defaultMode` | `int32` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].hostPath` | `HostPathVolumeSource` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].hostPath.path` | `string` | yes |  |  |
| `spec.container.sidecars[].volumeMounts[].hostPath.type` | `string` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].emptyDir` | `EmptyDirVolumeSource` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].emptyDir.medium` | `string` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].emptyDir.sizeLimit` | `string` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].pvc` | `PvcVolumeSource` |  |  |  |
| `spec.container.sidecars[].volumeMounts[].pvc.claimName` | `string` | yes |  |  |
| `spec.container.sidecars[].volumeMounts[].pvc.readOnly` | `bool` |  |  |  |
| `spec.container.sidecars[].lifecycle` | `WorkloadContainerLifecycle` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart` | `WorkloadLifecycleHandler` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.exec` | `ExecAction` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.exec.command` | `[]string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet.path` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet.portName` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet.host` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet.scheme` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.tcpSocket.host` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.sleep` | `SleepAction` |  |  |  |
| `spec.container.sidecars[].lifecycle.postStart.sleep.seconds` | `int64` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop` | `WorkloadLifecycleHandler` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.exec` | `ExecAction` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.exec.command` | `[]string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet.path` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet.portName` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet.host` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet.scheme` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.tcpSocket.portName` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.tcpSocket.host` | `string` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.sleep` | `SleepAction` |  |  |  |
| `spec.container.sidecars[].lifecycle.preStop.sleep.seconds` | `int64` |  |  |  |
| `spec.container.sidecars[].securityContext` | `WorkloadContainerSecurityContext` |  |  |  |
| `spec.container.sidecars[].securityContext.privileged` | `bool` |  |  |  |
| `spec.container.sidecars[].securityContext.runAsUser` | `int64` |  |  |  |
| `spec.container.sidecars[].securityContext.runAsGroup` | `int64` |  |  |  |
| `spec.container.sidecars[].securityContext.runAsNonRoot` | `bool` |  |  |  |
| `spec.container.sidecars[].securityContext.readOnlyRootFilesystem` | `bool` |  |  |  |
| `spec.container.sidecars[].securityContext.allowPrivilegeEscalation` | `bool` |  |  |  |
| `spec.container.sidecars[].securityContext.capabilities` | `WorkloadCapabilities` |  |  |  |
| `spec.container.sidecars[].securityContext.capabilities.add` | `[]string` |  |  |  |
| `spec.container.sidecars[].securityContext.capabilities.drop` | `[]string` |  |  |  |
| `spec.container.sidecars[].securityContext.seccompProfile` | `WorkloadSeccompProfile` |  |  |  |
| `spec.container.sidecars[].securityContext.seccompProfile.type` | `string` | yes |  |  |
| `spec.container.sidecars[].securityContext.seccompProfile.localhostProfile` | `string` |  |  |  |
| `spec.pod` | `WorkloadPod` |  |  |  |
| `spec.pod.serviceAccount` | `string \| valueFrom` |  |  | KubernetesServiceAccount (`status.outputs.service_account_name`) |
| `spec.pod.automountServiceAccountToken` | `bool` |  |  |  |
| `spec.pod.imagePullSecrets` | `[]string \| valueFrom` |  |  | KubernetesSecret (`spec.name`) |
| `spec.pod.initContainers` | `[]WorkloadContainer` |  |  |  |
| `spec.pod.initContainers[].name` | `string` |  |  |  |
| `spec.pod.initContainers[].image` | `ContainerImage` | yes |  |  |
| `spec.pod.initContainers[].image.repo` | `string` |  |  |  |
| `spec.pod.initContainers[].image.tag` | `string` |  |  |  |
| `spec.pod.initContainers[].image.pullSecretName` | `string` |  |  |  |
| `spec.pod.initContainers[].imagePullPolicy` | `string` |  |  |  |
| `spec.pod.initContainers[].command` | `[]string` |  |  |  |
| `spec.pod.initContainers[].args` | `[]string` |  |  |  |
| `spec.pod.initContainers[].workingDir` | `string` |  |  |  |
| `spec.pod.initContainers[].ports` | `[]WorkloadContainerPort` |  |  |  |
| `spec.pod.initContainers[].ports[].name` | `string` | yes |  |  |
| `spec.pod.initContainers[].ports[].containerPort` | `int32` | yes |  |  |
| `spec.pod.initContainers[].ports[].networkProtocol` | `string` |  |  |  |
| `spec.pod.initContainers[].ports[].appProtocol` | `string` |  |  |  |
| `spec.pod.initContainers[].ports[].servicePort` | `int32` |  |  |  |
| `spec.pod.initContainers[].ports[].hostPort` | `int32` |  |  |  |
| `spec.pod.initContainers[].env` | `ContainerEnv` |  |  |  |
| `spec.pod.initContainers[].env.variables` | `[]EnvVar` |  |  |  |
| `spec.pod.initContainers[].env.variables[].name` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.variables[].value` | `string` |  |  |  |
| `spec.pod.initContainers[].env.variables[].valueFrom` | `ValueFromRef` |  |  |  |
| `spec.pod.initContainers[].env.variables[].valueFrom.kind` | `enum` |  |  |  |
| `spec.pod.initContainers[].env.variables[].valueFrom.env` | `string` |  |  |  |
| `spec.pod.initContainers[].env.variables[].valueFrom.name` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.variables[].valueFrom.fieldPath` | `string` |  |  |  |
| `spec.pod.initContainers[].env.variables[].configMapKeyRef` | `ConfigMapKeyRef` |  |  |  |
| `spec.pod.initContainers[].env.variables[].configMapKeyRef.name` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.variables[].configMapKeyRef.key` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.variables[].configMapKeyRef.optional` | `bool` |  |  |  |
| `spec.pod.initContainers[].env.variables[].fieldRef` | `ObjectFieldRef` |  |  |  |
| `spec.pod.initContainers[].env.variables[].fieldRef.apiVersion` | `string` |  |  |  |
| `spec.pod.initContainers[].env.variables[].fieldRef.fieldPath` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.variables[].resourceFieldRef` | `ResourceFieldRef` |  |  |  |
| `spec.pod.initContainers[].env.variables[].resourceFieldRef.containerName` | `string` |  |  |  |
| `spec.pod.initContainers[].env.variables[].resourceFieldRef.resource` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.variables[].resourceFieldRef.divisor` | `string` |  |  |  |
| `spec.pod.initContainers[].env.secrets` | `[]SecretEnvVar` |  |  |  |
| `spec.pod.initContainers[].env.secrets[].name` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.secrets[].value` | `string` |  |  |  |
| `spec.pod.initContainers[].env.secrets[].secretRef` | `KubernetesSecretKeyRef` |  |  |  |
| `spec.pod.initContainers[].env.secrets[].secretRef.namespace` | `string` |  |  |  |
| `spec.pod.initContainers[].env.secrets[].secretRef.name` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.secrets[].secretRef.key` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.secrets[].secretRef.optional` | `bool` |  |  |  |
| `spec.pod.initContainers[].env.secrets[].valueFrom` | `ValueFromRef` |  |  |  |
| `spec.pod.initContainers[].env.secrets[].valueFrom.kind` | `enum` |  |  |  |
| `spec.pod.initContainers[].env.secrets[].valueFrom.env` | `string` |  |  |  |
| `spec.pod.initContainers[].env.secrets[].valueFrom.name` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.secrets[].valueFrom.fieldPath` | `string` |  |  |  |
| `spec.pod.initContainers[].env.envFrom` | `[]EnvFromSource` |  |  |  |
| `spec.pod.initContainers[].env.envFrom[].prefix` | `string` |  |  |  |
| `spec.pod.initContainers[].env.envFrom[].configMapRef` | `ConfigMapRef` |  |  |  |
| `spec.pod.initContainers[].env.envFrom[].configMapRef.name` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.envFrom[].configMapRef.optional` | `bool` |  |  |  |
| `spec.pod.initContainers[].env.envFrom[].secretRef` | `SecretRef` |  |  |  |
| `spec.pod.initContainers[].env.envFrom[].secretRef.name` | `string` | yes |  |  |
| `spec.pod.initContainers[].env.envFrom[].secretRef.optional` | `bool` |  |  |  |
| `spec.pod.initContainers[].resources` | `ContainerResources` |  |  |  |
| `spec.pod.initContainers[].resources.limits` | `CpuMemory` |  |  |  |
| `spec.pod.initContainers[].resources.limits.cpu` | `string` |  |  |  |
| `spec.pod.initContainers[].resources.limits.memory` | `string` |  |  |  |
| `spec.pod.initContainers[].resources.requests` | `CpuMemory` |  |  |  |
| `spec.pod.initContainers[].resources.requests.cpu` | `string` |  |  |  |
| `spec.pod.initContainers[].resources.requests.memory` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe` | `Probe` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.periodSeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.successThreshold` | `int32` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.failureThreshold` | `int32` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet.path` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet.host` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.grpc.port` | `int32` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.grpc.service` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.exec` | `ExecAction` |  |  |  |
| `spec.pod.initContainers[].livenessProbe.exec.command` | `[]string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe` | `Probe` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.periodSeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.successThreshold` | `int32` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.failureThreshold` | `int32` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet.path` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet.host` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.grpc.port` | `int32` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.grpc.service` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.exec` | `ExecAction` |  |  |  |
| `spec.pod.initContainers[].readinessProbe.exec.command` | `[]string` |  |  |  |
| `spec.pod.initContainers[].startupProbe` | `Probe` |  |  |  |
| `spec.pod.initContainers[].startupProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].startupProbe.periodSeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].startupProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.pod.initContainers[].startupProbe.successThreshold` | `int32` |  |  |  |
| `spec.pod.initContainers[].startupProbe.failureThreshold` | `int32` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet.path` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet.host` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.grpc` | `GRPCAction` |  |  |  |
| `spec.pod.initContainers[].startupProbe.grpc.port` | `int32` |  |  |  |
| `spec.pod.initContainers[].startupProbe.grpc.service` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.pod.initContainers[].startupProbe.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].startupProbe.tcpSocket.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.tcpSocket.host` | `string` |  |  |  |
| `spec.pod.initContainers[].startupProbe.exec` | `ExecAction` |  |  |  |
| `spec.pod.initContainers[].startupProbe.exec.command` | `[]string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts` | `[]VolumeMount` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].name` | `string` | yes |  |  |
| `spec.pod.initContainers[].volumeMounts[].mountPath` | `string` | yes |  |  |
| `spec.pod.initContainers[].volumeMounts[].readOnly` | `bool` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].subPath` | `string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].configMap` | `ConfigMapVolumeSource` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].configMap.name` | `string` | yes |  |  |
| `spec.pod.initContainers[].volumeMounts[].configMap.key` | `string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].configMap.path` | `string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].configMap.defaultMode` | `int32` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].secret` | `SecretVolumeSource` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].secret.name` | `string` | yes |  |  |
| `spec.pod.initContainers[].volumeMounts[].secret.key` | `string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].secret.path` | `string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].secret.defaultMode` | `int32` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].hostPath` | `HostPathVolumeSource` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].hostPath.path` | `string` | yes |  |  |
| `spec.pod.initContainers[].volumeMounts[].hostPath.type` | `string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].emptyDir` | `EmptyDirVolumeSource` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].emptyDir.medium` | `string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].emptyDir.sizeLimit` | `string` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].pvc` | `PvcVolumeSource` |  |  |  |
| `spec.pod.initContainers[].volumeMounts[].pvc.claimName` | `string` | yes |  |  |
| `spec.pod.initContainers[].volumeMounts[].pvc.readOnly` | `bool` |  |  |  |
| `spec.pod.initContainers[].lifecycle` | `WorkloadContainerLifecycle` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart` | `WorkloadLifecycleHandler` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.exec` | `ExecAction` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.exec.command` | `[]string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet.path` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet.host` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet.scheme` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.tcpSocket.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.tcpSocket.host` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.sleep` | `SleepAction` |  |  |  |
| `spec.pod.initContainers[].lifecycle.postStart.sleep.seconds` | `int64` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop` | `WorkloadLifecycleHandler` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.exec` | `ExecAction` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.exec.command` | `[]string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet` | `HTTPGetAction` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet.path` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet.host` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet.scheme` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet.httpHeaders` | `[]HTTPHeader` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet.httpHeaders[].name` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.httpGet.httpHeaders[].value` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.tcpSocket` | `TCPSocketAction` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.tcpSocket.portNumber` | `int32` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.tcpSocket.portName` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.tcpSocket.host` | `string` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.sleep` | `SleepAction` |  |  |  |
| `spec.pod.initContainers[].lifecycle.preStop.sleep.seconds` | `int64` |  |  |  |
| `spec.pod.initContainers[].securityContext` | `WorkloadContainerSecurityContext` |  |  |  |
| `spec.pod.initContainers[].securityContext.privileged` | `bool` |  |  |  |
| `spec.pod.initContainers[].securityContext.runAsUser` | `int64` |  |  |  |
| `spec.pod.initContainers[].securityContext.runAsGroup` | `int64` |  |  |  |
| `spec.pod.initContainers[].securityContext.runAsNonRoot` | `bool` |  |  |  |
| `spec.pod.initContainers[].securityContext.readOnlyRootFilesystem` | `bool` |  |  |  |
| `spec.pod.initContainers[].securityContext.allowPrivilegeEscalation` | `bool` |  |  |  |
| `spec.pod.initContainers[].securityContext.capabilities` | `WorkloadCapabilities` |  |  |  |
| `spec.pod.initContainers[].securityContext.capabilities.add` | `[]string` |  |  |  |
| `spec.pod.initContainers[].securityContext.capabilities.drop` | `[]string` |  |  |  |
| `spec.pod.initContainers[].securityContext.seccompProfile` | `WorkloadSeccompProfile` |  |  |  |
| `spec.pod.initContainers[].securityContext.seccompProfile.type` | `string` | yes |  |  |
| `spec.pod.initContainers[].securityContext.seccompProfile.localhostProfile` | `string` |  |  |  |
| `spec.pod.labels` | `map<string, string>` |  |  |  |
| `spec.pod.annotations` | `map<string, string>` |  |  |  |
| `spec.pod.scheduling` | `WorkloadScheduling` |  |  |  |
| `spec.pod.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.pod.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.pod.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.pod.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.pod.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.pod.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.pod.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.pod.scheduling.nodeAffinity` | `WorkloadNodeAffinity` |  |  |  |
| `spec.pod.scheduling.nodeAffinity.required` | `[]WorkloadNodeSelectorTerm` |  |  |  |
| `spec.pod.scheduling.nodeAffinity.required[].matchExpressions` | `[]WorkloadNodeSelectorRequirement` | yes |  |  |
| `spec.pod.scheduling.nodeAffinity.required[].matchExpressions[].key` | `string` | yes |  |  |
| `spec.pod.scheduling.nodeAffinity.required[].matchExpressions[].operator` | `string` | yes |  |  |
| `spec.pod.scheduling.nodeAffinity.required[].matchExpressions[].values` | `[]string` |  |  |  |
| `spec.pod.scheduling.nodeAffinity.preferred` | `[]WorkloadPreferredNodeSelectorTerm` |  |  |  |
| `spec.pod.scheduling.nodeAffinity.preferred[].weight` | `int32` |  |  |  |
| `spec.pod.scheduling.nodeAffinity.preferred[].term` | `WorkloadNodeSelectorTerm` | yes |  |  |
| `spec.pod.scheduling.nodeAffinity.preferred[].term.matchExpressions` | `[]WorkloadNodeSelectorRequirement` | yes |  |  |
| `spec.pod.scheduling.nodeAffinity.preferred[].term.matchExpressions[].key` | `string` | yes |  |  |
| `spec.pod.scheduling.nodeAffinity.preferred[].term.matchExpressions[].operator` | `string` | yes |  |  |
| `spec.pod.scheduling.nodeAffinity.preferred[].term.matchExpressions[].values` | `[]string` |  |  |  |
| `spec.pod.scheduling.podAffinity` | `WorkloadPodAffinity` |  |  |  |
| `spec.pod.scheduling.podAffinity.required` | `[]WorkloadPodAffinityTerm` |  |  |  |
| `spec.pod.scheduling.podAffinity.required[].matchLabels` | `map<string, string>` | yes |  |  |
| `spec.pod.scheduling.podAffinity.required[].topologyKey` | `string` | yes |  |  |
| `spec.pod.scheduling.podAffinity.required[].namespaces` | `[]string` |  |  |  |
| `spec.pod.scheduling.podAffinity.preferred` | `[]WorkloadWeightedPodAffinityTerm` |  |  |  |
| `spec.pod.scheduling.podAffinity.preferred[].weight` | `int32` |  |  |  |
| `spec.pod.scheduling.podAffinity.preferred[].term` | `WorkloadPodAffinityTerm` | yes |  |  |
| `spec.pod.scheduling.podAffinity.preferred[].term.matchLabels` | `map<string, string>` | yes |  |  |
| `spec.pod.scheduling.podAffinity.preferred[].term.topologyKey` | `string` | yes |  |  |
| `spec.pod.scheduling.podAffinity.preferred[].term.namespaces` | `[]string` |  |  |  |
| `spec.pod.scheduling.podAntiAffinity` | `WorkloadPodAffinity` |  |  |  |
| `spec.pod.scheduling.podAntiAffinity.required` | `[]WorkloadPodAffinityTerm` |  |  |  |
| `spec.pod.scheduling.podAntiAffinity.required[].matchLabels` | `map<string, string>` | yes |  |  |
| `spec.pod.scheduling.podAntiAffinity.required[].topologyKey` | `string` | yes |  |  |
| `spec.pod.scheduling.podAntiAffinity.required[].namespaces` | `[]string` |  |  |  |
| `spec.pod.scheduling.podAntiAffinity.preferred` | `[]WorkloadWeightedPodAffinityTerm` |  |  |  |
| `spec.pod.scheduling.podAntiAffinity.preferred[].weight` | `int32` |  |  |  |
| `spec.pod.scheduling.podAntiAffinity.preferred[].term` | `WorkloadPodAffinityTerm` | yes |  |  |
| `spec.pod.scheduling.podAntiAffinity.preferred[].term.matchLabels` | `map<string, string>` | yes |  |  |
| `spec.pod.scheduling.podAntiAffinity.preferred[].term.topologyKey` | `string` | yes |  |  |
| `spec.pod.scheduling.podAntiAffinity.preferred[].term.namespaces` | `[]string` |  |  |  |
| `spec.pod.scheduling.topologySpreadConstraints` | `[]WorkloadTopologySpreadConstraint` |  |  |  |
| `spec.pod.scheduling.topologySpreadConstraints[].maxSkew` | `int32` |  |  |  |
| `spec.pod.scheduling.topologySpreadConstraints[].topologyKey` | `string` | yes |  |  |
| `spec.pod.scheduling.topologySpreadConstraints[].whenUnsatisfiable` | `string` | yes |  |  |
| `spec.pod.scheduling.topologySpreadConstraints[].matchLabels` | `map<string, string>` |  |  |  |
| `spec.pod.scheduling.schedulerName` | `string` |  |  |  |
| `spec.pod.securityContext` | `WorkloadPodSecurityContext` |  |  |  |
| `spec.pod.securityContext.runAsUser` | `int64` |  |  |  |
| `spec.pod.securityContext.runAsGroup` | `int64` |  |  |  |
| `spec.pod.securityContext.runAsNonRoot` | `bool` |  |  |  |
| `spec.pod.securityContext.fsGroup` | `int64` |  |  |  |
| `spec.pod.securityContext.fsGroupChangePolicy` | `string` |  |  |  |
| `spec.pod.securityContext.supplementalGroups` | `[]int64` |  |  |  |
| `spec.pod.securityContext.sysctls` | `[]WorkloadSysctl` |  |  |  |
| `spec.pod.securityContext.sysctls[].name` | `string` | yes |  |  |
| `spec.pod.securityContext.sysctls[].value` | `string` | yes |  |  |
| `spec.pod.securityContext.seccompProfile` | `WorkloadSeccompProfile` |  |  |  |
| `spec.pod.securityContext.seccompProfile.type` | `string` | yes |  |  |
| `spec.pod.securityContext.seccompProfile.localhostProfile` | `string` |  |  |  |
| `spec.pod.terminationGracePeriodSeconds` | `int64` |  |  |  |
| `spec.pod.dnsPolicy` | `string` |  |  |  |
| `spec.pod.dnsConfig` | `WorkloadPodDnsConfig` |  |  |  |
| `spec.pod.dnsConfig.nameservers` | `[]string` |  |  |  |
| `spec.pod.dnsConfig.searches` | `[]string` |  |  |  |
| `spec.pod.dnsConfig.options` | `[]WorkloadPodDnsConfigOption` |  |  |  |
| `spec.pod.dnsConfig.options[].name` | `string` | yes |  |  |
| `spec.pod.dnsConfig.options[].value` | `string` |  |  |  |
| `spec.pod.hostAliases` | `[]WorkloadHostAlias` |  |  |  |
| `spec.pod.hostAliases[].ip` | `string` | yes |  |  |
| `spec.pod.hostAliases[].hostnames` | `[]string` | yes |  |  |
| `spec.pod.hostNetwork` | `bool` |  |  |  |
| `spec.pod.hostPid` | `bool` |  |  |  |
| `spec.pod.priorityClassName` | `string` |  |  |  |
| `spec.pod.runtimeClassName` | `string` |  |  |  |
| `spec.availability` | `KubernetesStatefulSetAvailability` |  |  |  |
| `spec.availability.replicas` | `int32` |  | `1` |  |
| `spec.availability.podDisruptionBudget` | `KubernetesStatefulSetPodDisruptionBudget` |  |  |  |
| `spec.availability.podDisruptionBudget.enabled` | `bool` |  |  |  |
| `spec.availability.podDisruptionBudget.minAvailable` | `string` |  |  |  |
| `spec.availability.podDisruptionBudget.maxUnavailable` | `string` |  |  |  |
| `spec.availability.minReadySeconds` | `int32` |  |  |  |
| `spec.availability.revisionHistoryLimit` | `int32` |  |  |  |
| `spec.volumeClaimTemplates` | `[]KubernetesStatefulSetVolumeClaimTemplate` |  |  |  |
| `spec.volumeClaimTemplates[].name` | `string` | yes |  |  |
| `spec.volumeClaimTemplates[].storageClass` | `string` |  |  |  |
| `spec.volumeClaimTemplates[].size` | `string` | yes |  |  |
| `spec.volumeClaimTemplates[].accessModes` | `[]string` |  |  |  |
| `spec.volumeClaimTemplates[].volumeMode` | `string` |  |  |  |
| `spec.updateStrategy` | `KubernetesStatefulSetUpdateStrategy` |  |  |  |
| `spec.updateStrategy.type` | `string` |  |  |  |
| `spec.updateStrategy.partition` | `int32` |  |  |  |
| `spec.updateStrategy.maxUnavailable` | `string` |  |  |  |
| `spec.podManagementPolicy` | `string` |  |  |  |
| `spec.pvcRetentionPolicy` | `KubernetesStatefulSetPvcRetentionPolicy` |  |  |  |
| `spec.pvcRetentionPolicy.whenDeleted` | `string` |  |  |  |
| `spec.pvcRetentionPolicy.whenScaled` | `string` |  |  |  |
| `spec.ordinals` | `KubernetesStatefulSetOrdinals` |  |  |  |
| `spec.ordinals.start` | `int32` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

The namespace to deploy into. Accepts a literal namespace name or a reference to a
KubernetesNamespace resource, so an infra chart creates the namespace and the
workload in one run.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.createNamespace

`bool`

When true, the module creates the namespace if it does not exist. Leave false when
the namespace is owned by a KubernetesNamespace resource or pre-exists.

### spec.container

`KubernetesStatefulSetContainer` · required

The containers of every replica: one main application container and any sidecars.
All containers share the pod's network namespace and volumes. Mount the
per-replica storage by declaring a volume mount whose PVC claim name matches a
`volume_claim_templates` entry.

- rule: {"required":true}

### spec.container.app

`WorkloadContainer` · required

The main application container. Its ports drive the workload's Services (both the
headless governing Service and the client-facing one), and its image is the
pipeline injection point.

- rule: {"required":true}

### spec.container.app.name

`string`

The container's name, unique within the pod. Required for sidecars and init
containers (Kubernetes rejects unnamed containers); for the main app container the
module defaults it when omitted, so minimal manifests stay minimal. Must be a valid
DNS label: lowercase alphanumeric and hyphens, starting and ending alphanumeric.

- rule: Container name must be a lowercase DNS label (alphanumeric and hyphens, starting and ending with an alphanumeric character)

### spec.container.app.image

`ContainerImage` · required

The container image, split into repository and tag so deployment pipelines can
inject a freshly built tag without rewriting the whole reference. The optional
`pull_secret_name` names an existing docker-registry secret; prefer attaching pull
secrets on the ServiceAccount (or `pod.image_pull_secrets`) so they apply pod-wide.

- rule: Image repo is required — the repository half of the image reference (e.g. "nginx" or "ghcr.io/acme/api")
- rule: Image tag is required — pin a version (e.g. "1.27.1"); avoid "latest" for anything you intend to roll back
- rule: {"required":true}

### spec.container.app.image.repo

`string`

The repository of the image (e.g., "gcr.io/project/image").

### spec.container.app.image.tag

`string`

The tag of the image (e.g., "latest" or "1.0.0").

### spec.container.app.image.pullSecretName

`string`

The name of the image pull secret for private image repositories.

### spec.container.app.imagePullPolicy

`string`

When the kubelet pulls the image. "IfNotPresent" (the Kubernetes default for tagged
images) reuses a cached copy; "Always" re-resolves the tag on every pod start —
required when a mutable tag like a branch name is reused across builds; "Never"
only uses pre-loaded images (air-gapped nodes, kind-loaded test images).

- rule: Image pull policy must be one of "Always", "IfNotPresent", or "Never"

### spec.container.app.command

`[]string`

Entrypoint override (Kubernetes `command`, Docker ENTRYPOINT). The image's
entrypoint runs when omitted. Not executed in a shell — provide argv elements,
e.g. ["/bin/sh", "-c", "exec my-server"].

### spec.container.app.args

`[]string`

Arguments to the entrypoint (Kubernetes `args`, Docker CMD). The image's CMD is
used when omitted. Variable references like $(VAR_NAME) are expanded from the
container's environment by the kubelet.

### spec.container.app.workingDir

`string`

Working directory for the entrypoint. Defaults to the image's configured WORKDIR.

### spec.container.app.ports

`[]WorkloadContainerPort`

Network ports this container exposes. Purely informational to Kubernetes for plain
pod-to-pod traffic, but load-bearing here: named ports are referenced by probes,
and `service_port` drives the Service wiring on kinds that create one
(Deployment, StatefulSet).

### spec.container.app.ports[].name

`string` · required

Port name, e.g. "http", "grpc", "metrics". Must be a lowercase DNS label that
starts and ends alphanumeric. Named ports are referenced by probes and become the
Service port names on service-fronted kinds.

- rule: Port name must contain only lowercase alphanumeric characters and hyphens, and start and end with an alphanumeric character (e.g. "http", "grpc-web")
- rule: {"required":true}

### spec.container.app.ports[].containerPort

`int32` · required

The port number the container listens on (1–65535).

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.container.app.ports[].networkProtocol

`string`

L4 protocol of the port. Defaults to "TCP" when omitted — the overwhelmingly
common case, so minimal manifests need not repeat it.

- rule: The network protocol must be one of "TCP", "UDP", or "SCTP"

### spec.container.app.ports[].appProtocol

`string`

Application protocol hint (e.g. "http", "grpc", "https"). Propagated to the
Service port's appProtocol on service-fronted kinds, where meshes and L7 load
balancers use it to pick the right protocol handling.

### spec.container.app.ports[].servicePort

`int32`

The port the workload's Kubernetes Service exposes for this container port.
Only meaningful on kinds that create a Service (Deployment, StatefulSet); other
kinds ignore it. E.g. containerPort 8080 with servicePort 80 serves the app on
the conventional port while the process binds an unprivileged one. External
exposure is composed separately with first-class ingress kinds referencing the
workload's exported Service handle — workloads never create ingress themselves.

- rule: Service port must be between 1 and 65535

### spec.container.app.ports[].hostPort

`int32`

Exposes the container port directly on the node's IP (hostPort). Chiefly a
DaemonSet pattern (node-level agents that must be reachable on every node);
on other kinds it constrains scheduling to one pod per node per port — prefer
a Service unless node-local reachability is the point.

- rule: Host port must be between 1 and 65535

### spec.container.app.env

`ContainerEnv`

Environment configuration: plain variables (with Kubernetes-native value sources
and Planton cross-resource references), secret variables (materialized into a
managed Kubernetes Secret), and bulk envFrom imports.

### spec.container.app.env.variables

`[]EnvVar`

Individual environment variables (non-sensitive).

### spec.container.app.env.variables[].name

`string` · required

The environment variable name.
Must be a valid C_IDENTIFIER: starts with a letter or underscore,
followed by letters, digits, or underscores.

- rule: Must be a valid C_IDENTIFIER (start with letter/underscore, contain only letters, digits, underscores)
- rule: {"required":true}

### spec.container.app.env.variables[].value

`string`

Direct literal value.

### spec.container.app.env.variables[].valueFrom

`ValueFromRef`

Reference to another Planton resource's field.
The orchestrator resolves this and populates the value before invoking IaC modules.

### spec.container.app.env.variables[].valueFrom.kind

`enum`

Allowed values (use exactly as shown):

- `unspecified` -- 0: Default/unspecified
- `TestCloudResourceGeneric` -- 1–49: Test/dev/custom
- `TestCloudResourceKubernetes`
- `ConfluentKafka` -- 50–199: saas platform resources
- `AtlasMongodb`
- `SnowflakeDatabase`
- `AwsAlb` -- 1000–1999: AWS resources AwsSubnet is a prerequisite because an ALB requires at least two subnets in different availability zones -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsCertManagerCert`
- `AwsCloudFront`
- `AwsDynamodb`
- `AwsEcrRepo`
- `AwsEcsCluster`
- `AwsEcsService` -- AwsEcsCluster, AwsEcsTaskDefinition, and AwsSubnet are prerequisites because a service schedules a referenced task-definition revision into a referenced live cluster and places task network interfaces into referenced subnets -- all three references must resolve first.
- `AwsEksCluster` -- AwsSubnet and AwsIamRole are prerequisites because the control plane attaches its network interfaces into referenced subnets and assumes a referenced cluster role that must already carry AmazonEKSClusterPolicy.
- `AwsIamRole`
- `AwsLambda`
- `AwsRdsCluster`
- `AwsRdsInstance`
- `AwsRoute53Zone`
- `AwsS3Bucket`
- `AwsLbTargetGroup` -- AwsVpc is a prerequisite because a target group's health checks and target registrations live inside one VPC -- the spec's vpc_id reference must resolve before the group can be created.
- `AwsSecurityGroup` -- AwsVpc is a prerequisite because every security group is created in a VPC; the E2E install profile resolves vpc_id against the VPC prerequisite.
- `AwsVpc`
- `AwsEksNodeGroup` -- AwsEksCluster is a prerequisite because nodes register with a live control plane; AwsIamRole and AwsSubnet back the node role and worker subnet references.
- `AwsIamUser`
- `AwsKmsKey`
- `AwsEc2Instance`
- `AwsClientVpn` -- Every Client VPN endpoint requires an ACM server certificate at create time; the imported self-signed fixture satisfies it. Subnets/VPC are optional composition (a zero-association endpoint is valid) -- composed scenarios declare them via the e2e-prerequisites annotation.
- `AwsDocumentDb`
- `AwsRoute53DnsRecord` -- AwsRoute53Zone is a prerequisite because every record lives inside a hosted zone -- the spec's zone_id reference must resolve before the record can be created.
- `AwsS3ObjectSet` -- AwsS3Bucket is a prerequisite because the object set's bucket reference is required -- objects cannot exist without the bucket that holds them.
- `AwsSqsQueue`
- `AwsSnsTopic`
- `AwsEventBridgeBus`
- `AwsEventBridgeRule`
- `AwsIamOidcProvider`
- `AwsIamPolicy`
- `AwsIamInstanceProfile` -- AwsIamRole is a prerequisite because an instance profile is a wrapper that must contain a role to be useful -- the profile's spec requires a role reference, so the role must be deployed first.
- `AwsLbListener` -- AwsAlb and AwsLbTargetGroup are prerequisites because a listener is an attachment point on a load balancer and its default action almost always forwards to a target group -- both references must resolve before the listener can be created.
- `AwsLbListenerRule` -- AwsLbListener is a prerequisite because a rule only exists as an attachment on a listener -- the listener_arn reference must resolve before the rule can be created.
- `AwsLaunchTemplate`
- `AwsAutoScalingGroup` -- AwsSubnet and AwsLaunchTemplate are prerequisites because a group cannot exist without subnets to place capacity in and a launch template to launch from -- the spec's subnets and launch_template references must resolve before the group can be created.
- `AwsEksAddon` -- AwsEksCluster is a prerequisite because an add-on installs onto a live control plane -- the spec's cluster_name reference must resolve before the add-on can be created.
- `AwsEksFargateProfile` -- AwsEksCluster, AwsIamRole, and AwsSubnet are prerequisites because a Fargate profile attaches to a live control plane, runs pods as a referenced pod-execution role, and launches them into referenced private subnets -- all three references must resolve first.
- `AwsEksAccessEntry` -- AwsEksCluster and AwsIamRole are prerequisites because an access entry grants a referenced IAM principal access to a live control plane -- both references must resolve before the entry can be created.
- `AwsEcsTaskDefinition` -- AwsIamRole is a prerequisite because the kind's default posture -- Fargate with the awslogs logging default -- is rejected by AWS at registration time without an execution role the agent can assume.
- `AwsHttpApiGateway`
- `AwsStepFunction` -- AwsIamRole is a prerequisite because a state machine cannot be created without an execution role it can assume -- the spec's role_arn reference must resolve before the CreateStateMachine call.
- `AwsHttpApiVpcLink` -- AwsSubnet is a prerequisite because a VPC link is a set of managed ENIs provisioned into referenced subnets -- the subnet references must resolve before the link can be created. Security groups are optional on the link, so they compose per-scenario rather than as a registry prerequisite.
- `AwsHttpApiDomain` -- AwsCertManagerCert is a prerequisite because a custom domain cannot be created without a TLS certificate in the same region covering the domain -- the spec's certificate_arn reference must resolve first.
- `AwsVpcEndpoint` -- AwsVpcEndpoint's composed E2E scenarios reference the AwsVpc prerequisite's outputs (vpc_id + default_route_table_id for gateway endpoints) and the AwsSubnet pair's subnet_id outputs (interface endpoints), so both are genuine deploy-order prerequisites.
- `AwsElasticacheUser`
- `AwsElasticacheUserGroup` -- AwsElasticacheUser is a genuine prerequisite: AWS refuses to create a user group that does not contain a user named "default", so a group's composed E2E scenario must resolve a deployed user's outputs.
- `AwsRedshiftServerlessNamespace`
- `AwsRedshiftServerlessWorkgroup` -- The namespace is a genuine prerequisite: a workgroup attaches to exactly one namespace by name at create time, so its composed E2E scenario must resolve a deployed namespace's outputs. AwsSubnet is a prerequisite because Redshift Serverless requires the workgroup's subnets to span three availability zones.
- `AwsRedisElasticache` -- AwsSubnet is a prerequisite because the module builds an ElastiCache subnet group from referenced subnets -- the spec's subnet references must resolve before the replication group can deploy.
- `AwsOpenSearchDomain`
- `AwsMemcachedElasticache`
- `AwsServerlessElasticache`
- `AwsNlb` -- AwsSubnet is a prerequisite because an NLB requires at least one subnet mapping -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsElasticIp`
- `AwsTransitGateway`
- `AwsGlobalAccelerator`
- `AwsSubnet`
- `AwsInternetGateway`
- `AwsNatGateway` -- AwsInternetGateway is a prerequisite because a public NAT gateway can only become available once the VPC it sits in has an internet gateway attached (AWS rejects the create otherwise) -- so the gateway must be deployed first. AwsVpc is a prerequisite because a REGIONAL NAT gateway (availability_mode = regional) references the VPC directly instead of a subnet.
- `AwsEgressOnlyInternetGateway`
- `AwsElasticFileSystem` -- AwsSubnet and AwsSecurityGroup are prerequisites because mount targets (required, min 1) place the file system's NFS endpoints into subnets and attach security groups -- both references must resolve before the CreateMountTarget calls.
- `AwsEfsAccessPoint` -- AwsElasticFileSystem is a prerequisite because an access point is created INTO a file system -- the spec's required file_system_id reference must resolve before the CreateAccessPoint call.
- `AwsFsxLustreFileSystem`
- `AwsFsxOpenzfsFileSystem`
- `AwsFsxWindowsFileSystem` -- Every Windows file system must join an Active Directory domain; the directory itself is external infrastructure (AWS Managed Microsoft AD or a self-managed domain), so only the network dependency is a declarable prerequisite.
- `AwsFsxOntapFileSystem`
- `AwsFsxOntapStorageVirtualMachine`
- `AwsFsxOntapVolume`
- `AwsFsxDataRepositoryAssociation`
- `AwsCognitoUserPool`
- `AwsCognitoIdentityProvider` -- AwsCognitoUserPool is a prerequisite because an identity provider is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateIdentityProvider call.
- `AwsCognitoUserPoolClient` -- AwsCognitoUserPool is a prerequisite because an app client is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateUserPoolClient call.
- `AwsCognitoResourceServer` -- AwsCognitoUserPool is a prerequisite because a resource server is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateResourceServer call.
- `AwsWafWebAcl`
- `AwsWafIpSet`
- `AwsWafRegexPatternSet`
- `AwsCloudwatchLogGroup`
- `AwsCloudwatchAlarm`
- `AwsCloudwatchCompositeAlarm`
- `AwsKinesisStream`
- `AwsKinesisFirehose` -- Every Firehose destination requires an S3 configuration (the primary target for extended_s3; the failed/all-document backup for the rest) and an IAM role Firehose assumes to write to it, so both are hard deploy prerequisites.
- `AwsKinesisStreamConsumer` -- A consumer registers against exactly one stream and cannot exist without it.
- `AwsAthenaWorkgroup`
- `AwsGlueCatalogDatabase`
- `AwsRedshiftCluster`
- `AwsSagemakerDomain` -- AI/ML A domain cannot exist without VPC subnets and a SageMaker execution role (default_user_settings.execution_role_arn is required), so both are hard deploy prerequisites.
- `AwsAppRunnerService` -- A service can run entirely on companion defaults, so the App Runner family's kinds are dependency-free leaves except the VPC connector (which cannot exist without subnets and security groups). A service's companion references (auto scaling / VPC connector / observability / WAF) are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsAppRunnerAutoScalingConfiguration`
- `AwsAppRunnerVpcConnector`
- `AwsAppRunnerObservabilityConfiguration`
- `AwsTransitGatewayVpcAttachment` -- AwsTransitGateway is a prerequisite because an attachment cannot exist without the gateway it attaches to; AwsSubnet because the attachment provisions an ENI into at least one subnet (the VPC arrives transitively through the subnet's own prerequisites).
- `AwsTransitGatewayRouteTable` -- Only the gateway is a hard prerequisite: a route table can exist empty. Associations, propagations, and routes referencing attachments are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsBatchComputeEnvironment` -- A MANAGED compute environment always launches into VPC subnets, so the subnet is a hard deploy prerequisite (security groups are required only for the Fargate types -- scenario-declared, not a registry edge).
- `AwsBatchJobQueue` -- A job queue cannot exist without at least one VALID compute environment to map onto.
- `AwsBatchSchedulingPolicy`
- `AwsBatchJobDefinition`
- `AwsCodeBuildProject` -- CI/CD
- `AwsCodePipeline`
- `AwsMwaaEnvironment` -- Workflow / Orchestration AwsSubnet and AwsSecurityGroup are prerequisites because the environment's network interfaces are placed in referenced private subnets and AWS requires at least one attached security group at creation.
- `AwsNeptuneCluster` -- Graph Database
- `AwsMemorydbCluster` -- A cluster always launches into a subnet group; the subnets are the hard deploy prerequisite. The ACL it attaches is optional composition (the built-in "open-access" ACL needs no resource) -- scenarios declare the ACL/user chain via the e2e-prerequisites annotation.
- `AwsMemorydbUser`
- `AwsMemorydbAcl` -- An empty ACL is valid (MemoryDB has no mandatory "default" member), so the user is optional composition -- the composed scenario declares it via the e2e-prerequisites annotation, never a registry edge.
- `AwsMskCluster` -- Streaming AwsSubnet and AwsSecurityGroup are prerequisites because brokers are placed in referenced subnets and AWS requires at least one attached security group at creation.
- `AwsMskServerlessCluster` -- AwsSubnet is a prerequisite because the serverless cluster's network interfaces are placed in referenced subnets (security groups are optional -- AWS attaches the VPC default group when none are referenced).
- `AwsLambdaEventSourceMapping` -- AwsLambda is a prerequisite because a mapping cannot exist without the function it invokes (a required reference). Event sources (SQS, Kinesis, DynamoDB, MSK) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsSnsSubscription` -- AwsSnsTopic is a prerequisite because a subscription cannot exist without the topic it subscribes to (a required reference). Endpoints (SQS queues, Lambda functions, Firehose streams) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsPlantonRunner` -- AwsSubnet is a prerequisite because the runner appliance places its network interfaces into referenced subnets -- the placement reference must resolve before the appliance can deploy.
- `AwsRoute53HealthCheck`
- `AwsSesConfigurationSet` -- Both SES kinds are dependency-free leaves: an identity's configuration set is optional composition (scenarios declare it via the e2e-prerequisites annotation), and a configuration set's event destinations reference other kinds only optionally.
- `AwsSesEmailIdentity`
- `AwsSecretsManagerSecret` -- A dependency-free leaf: the KMS key, rotation Lambda, and external rotation role references are all optional composition -- scenarios declare them via the e2e-prerequisites annotation, never registry edges.
- `AwsOpenSearchServerlessCollection` -- A dependency-free leaf: the collection-scoped encryption/network/ data-access/retention policies are module-rendered, and the KMS key and data-access principal references are optional composition (e2e-prerequisites annotation).
- `AwsBedrockGuardrail` -- A dependency-free leaf: the KMS key reference is optional composition (e2e-prerequisites annotation); published versions are folded satellites of the guardrail itself.
- `AwsBedrockCustomModel` -- AwsIamRole is a prerequisite because Bedrock assumes the job role to read training data and write outputs; the S3 locations and KMS key are optional composition (e2e-prerequisites annotation).
- `AwsBedrockInferenceProfile` -- A dependency-free leaf: the model source is a foundation model or an AWS system-defined cross-region profile, never a customer resource.
- `AwsBedrockProvisionedThroughput` -- A dependency-free leaf in the registry: capacity is typically bought for an AwsBedrockCustomModel (the default reference), but foundation model ARNs are equally legal, so the edge is optional composition.
- `AwsBedrockModelAccess` -- A dependency-free leaf: the agreement covers an AWS-listed foundation model, never a customer resource.
- `AwsBedrockAgent` -- AwsIamRole is a prerequisite because the Bedrock service assumes the agent resource role to invoke models, action-group Lambdas, and knowledge bases; the guardrail, KMS key, provisioned throughput, and collaborator/knowledge-base edges are optional composition (e2e-prerequisites annotation). Action groups, aliases, collaborators, and knowledge-base associations are folded satellites of the agent.
- `AwsBedrockKnowledgeBase` -- AwsIamRole is a prerequisite because the Bedrock service assumes the knowledge-base role to read data sources, call the embedding model, and read/write the vector store; the vector-store and data-source reference edges (OpenSearch, S3, Secrets Manager, ...) are optional composition (e2e-prerequisites annotation). Data sources are folded satellites of the knowledge base.
- `AwsBedrockFlow` -- AwsIamRole is a prerequisite because the Bedrock service assumes the flow execution role to invoke the models, agents, knowledge bases, and Lambdas its nodes reference; the node-level reference edges are optional composition (e2e-prerequisites annotation).
- `AwsBedrockPrompt` -- A dependency-free leaf: variants target AWS-listed foundation models by ID; targeting another agent's alias is optional composition (e2e-prerequisites annotation).
- `AzureResourceGroup` -- 2000–2999: Azure resources
- `AzureAksCluster` -- AzureResourceGroup is the only required parent: the cluster is created inside a referenced resource group. Subnet is optional on the default node pool (AKS provisions managed networking when unset).
- `AzureAksNodePool` -- AzureAksCluster is a prerequisite because a node pool attaches to an existing cluster by ARM ID; the resource group chains transitively.
- `AzureContainerRegistry` -- AzureResourceGroup is a prerequisite because a container registry is created inside a resource group.
- `AzureDnsZone` -- AzureResourceGroup is a prerequisite because the DNS zone is created inside a referenced resource group that must already exist.
- `AzureKeyVault` -- AzureResourceGroup is a prerequisite because a key vault is created inside a referenced resource group in composed environments.
- `AzureVirtualNetwork` -- AzureResourceGroup is a prerequisite because a virtual network is created inside a referenced resource group in composed environments.
- `AzureNatGateway` -- AzureResourceGroup is a prerequisite because a NAT gateway is created inside a referenced resource group in composed environments.
- `AzureVirtualMachine` -- AzureNetworkInterface is a prerequisite because a virtual machine attaches at least one NIC (the subnet, network, and resource group chain transitively through the NIC's own prerequisites).
- `AzureStorageAccount` -- AzureResourceGroup is a prerequisite because a storage account is created inside a referenced resource group in composed environments.
- `AzureDnsRecord` -- AzureDnsZone is a prerequisite because a record set is created inside a referenced zone (the resource group chains transitively through the zone). Public DNS zone names are not globally unique, so a shared zone fixture is safe to recreate across scenarios.
- `AzureSubnet` -- AzureVirtualNetwork is a prerequisite because a subnet is an ARM child of a referenced network -- the network must exist before the subnet can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzureNetworkSecurityGroup` -- AzureResourceGroup is a prerequisite because a network security group is created inside a referenced resource group in composed environments.
- `AzurePublicIp` -- AzureResourceGroup is a prerequisite because a public IP is created inside a referenced resource group in composed environments.
- `AzurePrivateEndpoint` -- AzureSubnet is a prerequisite because a private endpoint draws its private IP from a referenced subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite). The connection target is polymorphic and the DNS zones / ASGs are optional, so none of those are prerequisites.
- `AzurePrivateDnsZone` -- AzureResourceGroup is a prerequisite because a private DNS zone is created inside a referenced resource group in composed environments.
- `AzureApplicationGateway` -- AzureSubnet is a prerequisite because a gateway cannot exist without its dedicated gateway_ip_configuration subnet (the network and resource group chain transitively through the subnet's own prerequisites); public frontends additionally reference a public IP, but private-only gateways are legal, so it is not a registry prerequisite.
- `AzureLoadBalancer` -- AzureResourceGroup is a prerequisite because a load balancer is created inside a referenced resource group (frontends additionally reference subnets or public IPs, but neither is universally required, so they are not registry prerequisites).
- `AzureRouteTable` -- AzureResourceGroup is a prerequisite because a route table is created inside a referenced resource group in composed environments.
- `AzurePrivateDnsZoneVirtualNetworkLink` -- AzurePrivateDnsZone and AzureVirtualNetwork are prerequisites because a virtual network link is a child resource of a referenced zone and binds it to a referenced network -- both must exist before the link can be written. (The resource group arrives transitively through the zone's and network's own prerequisite declarations.)
- `AzureVirtualNetworkPeering` -- AzureVirtualNetwork is a prerequisite because a peering is an ARM child of its local network and binds it to a remote network -- the local network must exist before the peering can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzurePublicIpPrefix` -- AzureResourceGroup is a prerequisite because a public IP prefix is created inside a referenced resource group in composed environments.
- `AzureNetworkInterface` -- AzureSubnet is a prerequisite because a network interface's IP configurations deploy into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureManagedDisk` -- AzureResourceGroup is a prerequisite because a managed disk is created inside a resource group.
- `AzureVirtualMachineScaleSet` -- AzureSubnet is a prerequisite because every scale-set instance's network interface deploys into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureKeyVaultKey` -- AzureKeyVault is a prerequisite because a key is a data-plane object inside a referenced vault -- the vault must exist before the key can be written (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultCertificate` -- AzureKeyVault is a prerequisite because a certificate is a data-plane object inside a referenced vault -- the vault must exist before the certificate can be enrolled or imported (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultSecret` -- AzureKeyVault is a prerequisite because a secret is a data-plane object inside a referenced vault -- the vault must exist before the secret can be written (the resource group chains transitively through the vault's own prerequisite). Part of the Key Vault family (2005, 2025-2026) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureWebApplicationFirewallPolicy` -- AzureResourceGroup is a prerequisite because a WAF policy is created inside a referenced resource group; the Application Gateways that attach the policy reference it, never the reverse.
- `AzureApplicationSecurityGroup` -- AzureResourceGroup is a prerequisite because an application security group is created inside a referenced resource group; network interfaces, scale-set IP configurations, and NSG security rules reference the group, never the reverse.
- `AzureDiskEncryptionSet` -- AzureKeyVaultKey is a prerequisite because a disk encryption set wraps customer data with a referenced key -- the key (and its vault, which chains transitively) must exist before the set can resolve the key URL at create time.
- `AzurePostgresqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureRedisCache` -- AzureResourceGroup is a prerequisite because the cache is created inside a referenced resource group (VNet injection additionally references a dedicated subnet, but only the Premium tier supports it, so it is not a registry prerequisite).
- `AzureCosmosdbAccount` -- AzureResourceGroup is a prerequisite because the account is created inside a referenced resource group.
- `AzureMssqlServer` -- AzureResourceGroup is a prerequisite because the logical server is created inside a referenced resource group.
- `AzureMysqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureMssqlDatabase` -- The parent logical server is referenced via server_id, not auto-deployed: E2E scenarios declare their own server fixture (minimal-server.yaml or the pool-attach chain through AzureMssqlElasticPool) so sequential subtests never destroy and recreate the same globally unique server_name.
- `AzureMssqlElasticPool` -- AzureMssqlServer is a prerequisite because every elastic pool lives on a referenced logical server (the server's resource group is transitive).
- `AzureRedisLinkedServer` -- The target and linked caches are referenced via ARM ids, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicy` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicyAssignment` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureContainerAppEnvironment` -- AzureResourceGroup is a prerequisite because the environment is created inside a referenced resource group that must already exist.
- `AzureContainerApp` -- AzureContainerAppEnvironment is a prerequisite because every app runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureServicePlan` -- AzureResourceGroup is a prerequisite because the plan is created inside a referenced resource group that must already exist.
- `AzureFunctionApp` -- AzureServicePlan is a prerequisite because a function app runs on a referenced plan (the resource group arrives transitively through the plan). The required storage account is deliberately NOT a registry prerequisite: storage-account names are globally unique, so scenarios bring their own scenario-local account fixtures.
- `AzureLinuxWebApp` -- AzureServicePlan is a prerequisite because a web app runs on a referenced plan (the resource group arrives transitively through the plan).
- `AzureContainerAppJob` -- AzureContainerAppEnvironment is a prerequisite because a job runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureContainerAppEnvironmentStorage` -- AzureContainerAppEnvironment is a prerequisite because the storage registration lives on a referenced environment. The Azure Files share and storage account are deliberately NOT registry prerequisites: storage-account names are globally unique, so scenarios bring their own scenario-local account + share fixtures.
- `AzureContainerAppEnvironmentDaprComponent` -- AzureContainerAppEnvironment is a prerequisite because the Dapr component is registered on a referenced environment.
- `AzureContainerAppEnvironmentCertificate` -- AzureContainerAppEnvironment is a prerequisite because the certificate is stored on a referenced environment.
- `AzureContainerAppEnvironmentManagedCertificate` -- AzureContainerAppEnvironment is a prerequisite because the managed certificate is provisioned on a referenced environment.
- `AzureLogAnalyticsWorkspace` -- AzureResourceGroup is a prerequisite because the workspace is created inside a referenced resource group that must already exist.
- `AzureApplicationInsights` -- AzureLogAnalyticsWorkspace is a prerequisite because workspace-based Application Insights stores its telemetry in a referenced workspace (the resource group chains transitively through the workspace).
- `AzureMonitorDiagnosticSetting` -- AzureLogAnalyticsWorkspace is a prerequisite because the setting's scenarios route a fixture workspace's telemetry (the workspace doubles as target and destination); the target itself is polymorphic.
- `AzureMonitorActionGroup` -- AzureResourceGroup is a prerequisite because the action group is created inside a referenced resource group that must already exist.
- `AzureMonitorMetricAlert` -- AzureMonitorActionGroup is a prerequisite because a metric alert's actions fire into a referenced action group (the resource group chains transitively); alert scopes are polymorphic.
- `AzureMonitorScheduledQueryAlert` -- AzureLogAnalyticsWorkspace is a prerequisite because the rule queries a referenced workspace scope; AzureMonitorActionGroup because its action fires into a referenced action group.
- `AzureMonitorActivityLogAlert` -- AzureMonitorActionGroup is a prerequisite because an activity log alert's actions fire into a referenced action group (the resource group chains transitively). The alert itself is subscription-global and its scopes are polymorphic.
- `AzureApplicationInsightsStandardWebTest` -- AzureApplicationInsights is a prerequisite because a standard web test binds to a referenced Application Insights component (the resource group chains transitively through the component).
- `AzureUserAssignedIdentity` -- AzureResourceGroup is a prerequisite because the identity is created inside a referenced resource group that must already exist.
- `AzureRoleAssignment` -- AzureResourceGroup and AzureUserAssignedIdentity are prerequisites because an assignment grants a role at a referenced scope (most commonly a resource group) to a referenced principal (most commonly a managed identity) -- both must exist before the grant can be written.
- `AzureRoleDefinition` -- AzureResourceGroup is a prerequisite because a custom role definition is created at a referenced scope, most commonly a resource group in composed environments -- the scope must exist before the definition can be written.
- `AzureFederatedIdentityCredential` -- AzureUserAssignedIdentity is the prerequisite because a federated identity credential is a child resource of a referenced managed identity -- the identity must exist before the credential can be written on it. (The resource group arrives transitively through the identity's own prerequisite declaration.)
- `AzureServiceBusNamespace` -- AzureResourceGroup is a prerequisite because a Service Bus namespace is created inside a referenced resource group in composed environments. The namespace is the container every Service Bus messaging entity (queue, topic, subscription, authorization rule, geo-DR pairing) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureEventHubNamespace` -- AzureResourceGroup is a prerequisite because an Event Hub namespace is created inside a referenced resource group in composed environments. The namespace is the container every Event Hubs entity (event hub, consumer group, authorization rule, schema group, geo-DR pairing, customer-managed key) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureServiceBusQueue`
- `AzureServiceBusTopic`
- `AzureServiceBusSubscription`
- `AzureServiceBusAuthorizationRule`
- `AzureServiceBusDisasterRecoveryConfig`
- `AzureEventHub`
- `AzureEventHubConsumerGroup`
- `AzureEventHubAuthorizationRule`
- `AzureFrontDoorProfile` -- AzureResourceGroup is a prerequisite because a Front Door profile is created inside a referenced resource group in composed environments. The profile is the container every Front Door delivery resource (endpoint, origin group, origin, route) nests under.
- `AzureFrontDoorEndpoint` -- AzureFrontDoorProfile is a prerequisite because an endpoint is an ARM child of a referenced profile -- the profile must exist before the endpoint can be written. (The resource group arrives transitively through the profile's own prerequisite declaration.)
- `AzureFrontDoorOriginGroup` -- AzureFrontDoorProfile is a prerequisite because an origin group is an ARM child of a referenced profile.
- `AzureFrontDoorOrigin` -- AzureFrontDoorOriginGroup is a prerequisite because an origin is an ARM child of a referenced origin group (the profile and resource group chain transitively).
- `AzureFrontDoorRoute` -- A route attaches to an endpoint (its ARM parent) and forwards to an origin group whose origins must exist before ARM accepts the route -- so both the endpoint and the origin chain are genuine deploy-order prerequisites.
- `AzureFrontDoorRuleSet` -- AzureFrontDoorProfile is a prerequisite because a rule set is an ARM child of a referenced profile. The rules live inside the set (they form one ordered policy document); routes attach the set by ARM ID.
- `AzureFrontDoorCustomDomain` -- AzureFrontDoorProfile is a prerequisite because a custom domain is an ARM child of a referenced profile. The DNS zone and (for bring-your-own certificates) the Front Door secret are optional references, not deploy-order prerequisites.
- `AzureFrontDoorSecret` -- AzureFrontDoorSecret is a prerequisite-light kind: only the profile (its ARM parent) must exist. The Key Vault certificate it wraps is a reference resolved before the module runs; its vault chain is exercised through scenario-local fixtures in E2E.
- `AzureFrontDoorFirewallPolicy` -- AzureResourceGroup is a prerequisite because the Front Door WAF policy is created inside a referenced resource group -- it is a GLOBAL resource, not a profile child (a different ARM type than the regional Application Gateway WAF policy). Security policies attach it to profiles; the policy itself depends on nothing else.
- `AzureFrontDoorSecurityPolicy` -- A security policy is an ARM child of a profile that associates a referenced WAF policy with referenced domains -- so the endpoint (the default-domain association target; the profile arrives transitively through it) and the WAF policy are genuine deploy-order prerequisites.
- `AzureStorageContainer` -- None of the storage data-service kinds declares a registry prerequisite on AzureStorageAccount: account names are GLOBALLY unique and Azure holds a just-deleted name, so a recreate-per-scenario fixture would hang -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the storage_account_id reference itself.
- `AzureStorageShare`
- `AzureStorageQueue`
- `AzureStorageTable`
- `AzureStorageEncryptionScope`
- `AzureStorageDataLakeGen2Filesystem`
- `AzureStorageLocalUser`
- `AzureStorageObjectReplication`
- `AzureCosmosdbSqlDatabase` -- None of the Cosmos DB data-service kinds declares a registry prerequisite on AzureCosmosdbAccount: account names are GLOBALLY unique DNS labels, so a recreate-per-scenario fixture would risk name-reuse hangs -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the cosmosdb_account_id / parent-database references themselves.
- `AzureCosmosdbSqlContainer`
- `AzureCosmosdbMongoDatabase`
- `AzureCosmosdbMongoCollection`
- `AzureCosmosdbSqlRoleDefinition`
- `AzureCosmosdbSqlRoleAssignment`
- `AzureManagedRedis` -- AzureResourceGroup is the cluster's only registry prerequisite: the cluster is created inside a referenced resource group. The geo-replication and access-policy-assignment children declare NO prerequisite on AzureManagedRedis: clusters are expensive, slow-provisioning parents, so their E2E scenarios declare scenario-local cluster fixtures instead of recreating a shared one per scenario. Deploy ordering in composed environments still flows from the managed_redis_id references themselves.
- `AzureManagedRedisGeoReplication`
- `AzureManagedRedisAccessPolicyAssignment`
- `AzureEventHubDisasterRecoveryConfig`
- `AzureEventHubSchemaGroup`
- `AzureEventHubCluster` -- AzureResourceGroup is a prerequisite because a dedicated Event Hubs cluster is created inside a referenced resource group in composed environments. Note: clusters cannot be deleted for 4 hours after creation (Azure's moratorium), so E2E treats this kind as offline-gated.
- `AzureEventHubNamespaceCustomerManagedKey`
- `AzureMssqlFailoverGroup` -- AzureMssqlServer is a prerequisite because a failover group is created on a referenced primary logical server and points at a partner server; the primary (and its resource group, which chains transitively) must exist before the group can be written.
- `AzureContainerAppCustomDomain` -- AzureContainerApp is a prerequisite because the domain binding lives in a referenced app's ingress configuration (the environment and resource group chain transitively through the app).
- `AzureFirewallPolicy`
- `AzureFirewallPolicyRuleCollectionGroup` -- AzureFirewallPolicy is a prerequisite because a rule collection group is a child document of a referenced policy (the resource group chains transitively through the policy).
- `AzureFirewall` -- AzureSubnet is a prerequisite because a VNet-deployed firewall's data path lives in a dedicated subnet that must be named exactly "AzureFirewallSubnet" (the virtual network and resource group chain transitively through the subnet). The E2E install profile publishes a fixture subnet with that exact name and a /26 prefix.
- `AzureIpGroup`
- `AzureVirtualNetworkGateway` -- AzureSubnet is a prerequisite because every virtual network gateway lives in a dedicated subnet named exactly "GatewaySubnet" (the virtual network and resource group chain transitively through the subnet); the subnet install profile publishes a fixture instance with that exact ARM name. AzurePublicIp is a prerequisite because a VPN-type gateway (the default shape) requires a public IP per ip configuration; the address install profile publishes a dedicated zone-redundant instance (a gateway binds its address exclusively, and the AZ gateway SKUs require zones on it).
- `AzureVirtualNetworkGatewayConnection` -- Both gateways are prerequisites: a connection joins a virtual network gateway to a far side, and the site-to-site far side is a local network gateway (the GatewaySubnet, VNet, and resource group chain transitively through the virtual network gateway).
- `AzureLocalNetworkGateway`
- `AzurePrivateLinkService` -- AzureSubnet is the sole prerequisite: every NAT ip configuration draws its address from a subnet with private-link-service network policies disabled (the subnet install profile publishes a fixture instance with that flag off). The Standard load balancer whose frontend the service typically fronts is NOT a registry prerequisite -- the spec's destination is an exactly-one-of (load balancer frontend OR fixed destination IP), so scenarios that use the load-balancer shape declare it via the planton.dev/e2e-prerequisites annotation instead.
- `AzureExpressRouteCircuit`
- `AzureExpressRouteCircuitPeering` -- The circuit is the prerequisite: a peering is an ARM child of the circuit, addressed by the circuit's name (the resource group chains transitively through the circuit).
- `AzureExpressRouteGateway` -- The hub is the prerequisite: ARM requires an ExpressRoute Gateway to be deployed INTO a Virtual WAN hub (the WAN and resource group chain transitively through the hub).
- `AzureExpressRoutePort` -- ExpressRoute Port: your own physical port pair on a Microsoft edge router (ExpressRoute Direct), from whose bandwidth circuits are carved. Self-contained -- only the resource group is required.
- `AzureVirtualWan` -- Virtual WAN: the umbrella of Azure's managed hub-and-spoke networking, under which virtual hubs and their gateways are created. Self-contained -- only the resource group is required.
- `AzureVirtualHub` -- The WAN is the prerequisite: this kind models the Virtual WAN hub (virtual_wan_id is required; standalone hubs are the legacy Route Server construction, which has its own ARM surface). The resource group chains transitively through the WAN.
- `AzureVirtualHubConnection` -- Both sides of the attachment are prerequisites: the hub being joined and the spoke virtual network being attached.
- `AzureVpnGateway` -- The hub is the prerequisite: ARM deploys a Virtual WAN VPN gateway INTO a virtual hub (virtual_hub_id is required and immutable; the WAN and resource group chain transitively through the hub). ARM allows one VPN gateway per hub.
- `AzureVpnGatewayConnection` -- Both ends of the tunnel are prerequisites: a connection is an ARM child of the VPN gateway and pins each of its links to a specific link of the remote VPN site (the hub, WAN, and resource group chain transitively through the gateway).
- `AzureVpnSite` -- The WAN is the prerequisite: a VPN site is the Virtual WAN world's address-book entry for one branch location (virtual_wan_id is required; the classic-world sibling without a WAN is AzureLocalNetworkGateway). The resource group chains transitively through the WAN.
- `AzurePointToSiteVpnGateway` -- The hub and the server configuration are both prerequisites: a point-to-site VPN gateway deploys INTO a virtual hub (one P2S gateway per hub, a slot separate from the hub's site-to-site VPN gateway) and is born pointing at the VPN server configuration that defines how its users authenticate -- both ARM-required and fixed at creation. The WAN and resource group chain transitively through the hub.
- `AzureVpnServerConfiguration` -- Self-contained -- only the resource group is required: a VPN server configuration is the reusable "who may connect and how" authentication policy (Entra ID / certificate / RADIUS) that point-to-site VPN gateways attach to; it references no other Azure resource.
- `AzureCognitiveAccount` -- Self-contained -- only the resource group is required: an Azure AI services account (Azure OpenAI, the multi-service AIServices account, the single-service accounts) needs no other Azure resource; subnets (network rules), Key Vault keys (CMK), storage accounts and user-assigned identities are optional references.
- `AzureCognitiveDeployment` -- An ARM child of its account: a model deployment (which model runs, at which throughput class) exists only on an Azure AI services account of kind "OpenAI" or "AIServices".
- `AzureCognitiveAccountProject` -- An ARM child of its account: an AI Foundry project exists only on an "AIServices"-kind account with project management enabled.
- `AzureMachineLearningWorkspace` -- The workspace REQUIRES all three companion services at creation (default storage, secrets vault, telemetry) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureMachineLearningDatastore` -- An ARM child of its workspace. The storage target (container, filesystem or share) is scenario-declared via the e2e-prerequisites annotation -- only the blob scenario needs a container, so it is not a kind-wide prerequisite.
- `AzureMachineLearningComputeCluster` -- An ARM child of its workspace (.../computes/{name}) -- the auto-scaling pool of VMs training jobs run on.
- `AzureMachineLearningComputeInstance` -- An ARM child of its workspace (.../computes/{name}) -- a single always-on VM serving as one data scientist's cloud workstation.
- `AzureAiFoundry` -- The hub REQUIRES both companion services at creation (secrets vault, default storage) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureAiFoundryProject` -- Deploys into its hub's resource group (the provider derives the group from the hub reference -- the project spec carries none).
- `AzureSearchService`
- `AzureMachineLearningOnlineEndpoint` -- An ARM child of its workspace (.../onlineEndpoints/{name}) -- the stable scoring address applications call. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningOnlineDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the running copy of a model the endpoint's traffic map routes to.
- `AzureMachineLearningBatchEndpoint` -- An ARM child of its workspace (.../batchEndpoints/{name}) -- the stable address batch scoring jobs are submitted to. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningBatchDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the job recipe (model, compute, batching behavior) the endpoint's default-deployment pointer routes submissions to.
- `AzureRecoveryServicesVault` -- The Recovery Services vault (Microsoft.RecoveryServices/vaults) -- the safe that classic Azure Backup data and Site Recovery configuration live in. Backup policies and protected items are ARM children of a vault.
- `AzureBackupPolicyVm` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern IaaS VM backups.
- `AzureBackupProtectedVm` -- An ARM child of its vault (.../protectedItems/...) -- the binding that puts one virtual machine under a backup policy's protection.
- `AzureBackupPolicyFileShare` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern Azure Files share backups (snapshot or vaulted).
- `AzureBackupProtectedFileShare` -- An ARM child of its vault (.../protectedItems/AzureFileShare;...) -- the binding that puts one Azure Files share under a backup policy's protection. The share's storage account must already be registered with the vault (AzureBackupContainerStorageAccount).
- `AzureDataProtectionBackupVault` -- The Data Protection backup vault (Microsoft.DataProtection/ backupVaults) -- the safe that MODERN Azure Backup data lives in (managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, Data Lake storage). Backup policies and backup instances are ARM children of a vault.
- `AzureDataProtectionBackupPolicy` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules for ONE Data Protection datasource type (blob storage, disk, Kubernetes cluster, MySQL/PostgreSQL flexible server, or Data Lake storage), modeled as one kind with variant blocks.
- `AzureDataProtectionBackupInstance` -- An ARM child of its vault (.../backupInstances/{name}) -- the binding that puts ONE datasource (a managed disk, a storage account's blob services, an AKS cluster, a MySQL/PostgreSQL flexible server, or a Data Lake storage account) under a Data Protection backup policy, modeled as one kind with variant blocks. The vault's managed identity must hold the datasource roles Azure Backup requires BEFORE the instance is created.
- `AzureBastionHost` -- AzureSubnet and AzurePublicIp are prerequisites because a dedicated-infrastructure Bastion host (Basic/Standard/Premium -- the default shapes) deploys into a subnet named exactly "AzureBastionSubnet" and binds a Standard static public IP EXCLUSIVELY (the virtual network and resource group chain transitively through the subnet). The Developer SKU instead attaches to a virtual network directly and uses neither.
- `AzureNetworkWatcherFlowLog` -- AzureVirtualNetwork and AzureStorageAccount are prerequisites because a flow log records a network-scoped target (a virtual network in the common case; subnets and network interfaces chain through the network) into a referenced storage account. The regional Network Watcher parent is NOT a prerequisite: Azure auto-creates it ("NetworkWatcher_{region}" in "NetworkWatcherRG") the moment the region hosts a virtual network, and the flow log references it by name. Traffic Analytics' Log Analytics workspace is an optional arm, declared by scenarios that use it.
- `AzurePrivateDnsResolver` -- AzureVirtualNetwork and AzureSubnet are prerequisites because a DNS Private Resolver anchors to a referenced virtual network (at most ONE resolver per network -- Azure enforces it) and each of its inbound/outbound endpoints occupies its own dedicated subnet delegated to "Microsoft.Network/dnsResolvers" (the resource group chains transitively through the network and subnets).
- `AzurePrivateDnsResolverForwardingRuleset` -- AzurePrivateDnsResolver is a prerequisite because a DNS forwarding ruleset steers a resolver's OUTBOUND endpoints -- it binds their ARM ids (at most 2, same resolver) at creation. (The resource group and network chain transitively through the resolver's own prerequisite declarations.)
- `AzureBackupContainerStorageAccount` -- Registers a storage account with a Recovery Services vault as a backup container (.../protectionContainers/StorageContainer;...) -- one registration per storage-account-and-vault pair, required BEFORE any of the account's file shares can be protected. Part of the backup family (2175-2179) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureDataProtectionResourceGuard` -- The Data Protection Resource Guard (Microsoft.DataProtection/ resourceGuards) -- the approval gate behind Multi-User Authorization: privileged vault operations (disabling soft delete, reducing retention) require an approval through a guard, which typically lives in a DIFFERENT administrator's scope. Vaults reference a guard by its ARM ID. Part of the backup family (2175-2182) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzurePrivateDnsResolverVirtualNetworkLink` -- The attachment that makes a DNS forwarding ruleset take effect in one virtual network ({ruleset_id}/virtualNetworkLinks/{name}) -- one link per ruleset-network pair, up to 500 per ruleset, spokes joining and leaving independently (which is why the link is a standalone kind, exactly like AzurePrivateDnsZoneVirtualNetworkLink). Part of the DNS Private Resolver family (2186-2187) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `GcpArtifactRegistryRepo` -- 3000–3999: GCP resources
- `GcpTargetHttpsProxy` -- The URL map is the parent a proxy cannot exist without; the classic compute certificate kinds and the SSL policy are the fixture parents the committed scenarios attach. The Certificate Manager certificate list (certificate_manager_certificates, honored only by the cross-region internal ALB) is optional composition -- a scenario that arms it declares GcpCertManagerCert via the e2e-prerequisites annotation, never a registry edge that would tax every proxy and forwarding-rule chain.
- `GcpCloudFunction`
- `GcpCloudRun`
- `GcpCloudSql`
- `GcpDnsZone`
- `GcpGcsBucket`
- `GcpGkeCluster`
- `GcpIamCustomRole`
- `GcpProject`
- `GcpVpcNetwork`
- `GcpSubnetwork`
- `GcpRouterNat`
- `GcpGkeNodePool`
- `GcpServiceAccount`
- `GcpGkeWorkloadIdentityBinding`
- `GcpCertManagerCert`
- `GcpComputeInstance`
- `GcpDnsRecord`
- `GcpProjectIamMember`
- `GcpFirewallRule`
- `GcpGlobalAddress`
- `GcpCloudArmorPolicy`
- `GcpHealthCheck`
- `GcpBackendBucket`
- `GcpBackendService`
- `GcpRegionNetworkEndpointGroup`
- `GcpUrlMap`
- `GcpManagedSslCertificate`
- `GcpTargetHttpProxy`
- `GcpAlloydbCluster`
- `GcpRedisInstance`
- `GcpFirestoreDatabase`
- `GcpSpannerInstance`
- `GcpSpannerDatabase`
- `GcpBigtableInstance`
- `GcpMemorystoreInstance`
- `GcpCloudSqlDatabase`
- `GcpCloudSqlUser`
- `GcpAlloydbInstance`
- `GcpAlloydbUser`
- `GcpSpannerBackupSchedule`
- `GcpBigtableTable`
- `GcpFirestoreBackupSchedule`
- `GcpFirestoreIndex`
- `GcpBigQueryDataset`
- `GcpDataprocCluster`
- `GcpDataprocAutoscalingPolicy`
- `GcpBigQueryTable`
- `GcpPubSubTopic`
- `GcpPubSubSubscription`
- `GcpCloudTasksQueue`
- `GcpCloudSchedulerJob`
- `GcpPubSubSchema`
- `GcpVertexAiNotebook`
- `GcpVertexAiEndpoint`
- `GcpVertexAiIndex`
- `GcpVertexAiIndexEndpoint` -- Vector Search IndexEndpoint — distinct from the online-prediction GcpVertexAiEndpoint (671); different GCP resources, different kinds.
- `GcpVertexAiDeployedIndex`
- `GcpCloudComposerEnvironment`
- `GcpCloudComposerUserWorkloadsSecret`
- `GcpCloudComposerUserWorkloadsConfigMap`
- `GcpKmsKeyRing`
- `GcpKmsKey`
- `GcpKmsKeyIamMember`
- `GcpFilestoreInstance`
- `GcpWorkloadIdentityPool` -- 3101–3109: IAM/identity family (overflow block; the 3000–3022 foundation/security sub-band is fully allocated)
- `GcpWorkloadIdentityPoolProvider`
- `GcpServiceAccountIamMember`
- `GcpGlobalForwardingRule` -- 3110–3119: networking/load-balancer family (overflow block; the 3023–3029 LB sub-band is fully allocated)
- `GcpSslPolicy`
- `GcpSslCertificate`
- `GcpServiceNetworkingConnection`
- `GcpAddress`
- `GcpServiceConnectionPolicy`
- `GcpCertManagerDnsAuthorization`
- `GcpCertificateMap` -- GcpCertManagerCert is a prerequisite because a map entry binds hostnames to EXISTING certificates — the canonical map references a certificate fixture's resource name.
- `GcpCloudRunJob` -- 3120–3129: GCP serverless overflow
- `GcpServerlessVpcConnector`
- `GcpComputeDisk` -- 3130–3139: GCP compute overflow (the 3000–3022 foundation sub-band that holds GcpComputeInstance is fully allocated)
- `GcpComputeMig` -- GcpVpcNetwork is a prerequisite because the canonical group runs its fleet on a dedicated custom-mode VPC — a managed instance group's template must attach every VM to a network, and the default VPC is never assumed.
- `GcpMonitoringNotificationChannel` -- 3140–3149: GCP observability & log routing
- `GcpMonitoringAlertPolicy` -- GcpMonitoringNotificationChannel is a prerequisite because the policy's canonical shape references a channel to notify — a policy without a delivery endpoint measures but never pages.
- `GcpMonitoringUptimeCheck`
- `GcpLoggingSink` -- GcpGcsBucket is a prerequisite because the canonical sink exports to a Cloud Storage bucket — the cheapest destination that proves the whole writer-identity grant flow.
- `GcpMonitoringDashboard`
- `GcpMonitoringSlo`
- `GcpLogBucket`
- `GcpLogMetric`
- `GcpSecretManagerSecret` -- 3150–3159: GCP security & identity GcpServiceAccount is a prerequisite because the canonical secret grants secretAccessor to a workload service account — the access story the kind exists to model.
- `GcpIdentityPlatformConfig`
- `GcpIdentityPlatformTenant` -- GcpIdentityPlatformConfig is a prerequisite because tenants exist only in projects whose Identity Platform config enables multi_tenant.allow_tenants — a tenant without the initialized, tenant-enabled project config cannot be created at all.
- `GcpIamOauthClient`
- `GcpIamDenyPolicy`
- `GcpCloudRunDomainMapping` -- 3160–3169: GCP serverless edge GcpCloudRun is a prerequisite because a domain mapping exists only to point a verified domain at a running Cloud Run service — the route it maps must already exist for the mapping to be created at all.
- `GcpWorkflow`
- `GcpEventarcTrigger` -- GcpCloudRun is a prerequisite because the canonical trigger routes a Pub/Sub messagePublished event to a Cloud Run service — the destination story the kind exists to model.
- `GcpEventarcMessageBus`
- `KubernetesNamespace` -- 4000–4999: Kubernetes resources, organized in family sub-bands (4030–4069 also hosts CNI/autoscaling/DR addons; 4130–4149 hosts analytics & ML; 4190–4199 reserved for growth) 4000–4029: Kubernetes building blocks (core API primitives)
- `KubernetesDeployment`
- `KubernetesStatefulSet`
- `KubernetesDaemonSet`
- `KubernetesJob`
- `KubernetesCronJob`
- `KubernetesService`
- `KubernetesSecret`
- `KubernetesManifest`
- `KubernetesHelmRelease`
- `KubernetesConfigMap`
- `KubernetesServiceAccount`
- `KubernetesRbac` -- Bundles the RBAC grant grain (Role/ClusterRole + its binding) into one component: "grant these permissions to these subjects in this scope".
- `KubernetesIngress`
- `KubernetesNetworkPolicy`
- `KubernetesPersistentVolumeClaim`
- `KubernetesStorageClass`
- `KubernetesResourceQuota` -- Manages the namespace-governance pair: the ResourceQuota plus an optional companion LimitRange (per-object defaults/bounds) — two API objects, one governance story.
- `KubernetesPriorityClass`
- `KubernetesPodDisruptionBudget`
- `KubernetesHorizontalPodAutoscaler`
- `KubernetesCertManager` -- 4030–4069: Kubernetes foundation addons (certs, DNS, secrets, ingress, Gateway API, mesh, CNI/autoscaling/DR)
- `KubernetesClusterIssuer` -- KubernetesCertManager is a prerequisite for the three cert-manager CR kinds below: ClusterIssuer/Issuer/Certificate are cert-manager custom resources — without the controller and its CRDs they cannot be applied.
- `KubernetesIssuer`
- `KubernetesCertificate`
- `KubernetesExternalDns`
- `KubernetesExternalSecretsOperator`
- `KubernetesClusterSecretStore` -- KubernetesExternalSecretsOperator is a prerequisite for the three external-secrets CR kinds below: ClusterSecretStore/SecretStore/ ExternalSecret are external-secrets custom resources — without the operator and its CRDs they cannot be applied.
- `KubernetesSecretStore`
- `KubernetesExternalSecret`
- `KubernetesIngressNginx`
- `KubernetesGatewayApiCrds`
- `KubernetesGatewayClass`
- `KubernetesGateway`
- `KubernetesListenerSet`
- `KubernetesHttpRoute`
- `KubernetesGrpcRoute`
- `KubernetesTcpRoute`
- `KubernetesUdpRoute`
- `KubernetesTlsRoute`
- `KubernetesReferenceGrant`
- `KubernetesBackendTlsPolicy`
- `KubernetesIstioBaseCrds`
- `KubernetesIstio`
- `KubernetesDestinationRule` -- Istio API components (mesh traffic policy, security, telemetry). The seven typed resources below (4053–4059) require the Istio CRDs on the cluster, provided by the lightweight CRDs-only KubernetesIstioBaseCrds (851) — NOT the full mesh KubernetesIstio (852).
- `KubernetesServiceEntry`
- `KubernetesPeerAuthentication`
- `KubernetesRequestAuthentication`
- `KubernetesAuthorizationPolicy`
- `KubernetesTelemetry`
- `KubernetesEnvoyFilter`
- `KubernetesMetricsServer`
- `KubernetesCilium`
- `KubernetesKeda`
- `KubernetesKarpenter`
- `KubernetesKarpenterNodePool`
- `KubernetesKarpenterEc2NodeClass`
- `KubernetesClusterAutoscaler`
- `KubernetesVelero`
- `KubernetesKubePrometheusStack` -- 4070–4089: Kubernetes observability
- `KubernetesGrafana`
- `KubernetesSignoz` -- KubernetesClickHouse is a prerequisite because SigNoz stores every trace, metric and log in ClickHouse and deploys none of its own — the telemetry store is composed, never bundled.
- `KubernetesLoki`
- `KubernetesTempo`
- `KubernetesOtelOperator` -- The operator's admission webhooks (failurePolicy Fail) are served with a cert-manager Certificate in the default posture — cert-manager must be running before the operator installs.
- `KubernetesOtelCollector`
- `KubernetesKyverno` -- 4080–4099: Kubernetes security, policy, and identity
- `KubernetesGatekeeper`
- `KubernetesKeycloak` -- Keycloak declarations compose the official Keycloak Operator (which reconciles the Keycloak CR this kind renders) and, on the recommended postgres vendor, a KubernetesPostgres database — both must resolve before the CR can converge.
- `KubernetesOpenBao`
- `KubernetesOpenFga` -- OpenFGA requires a datastore; the recommended arm composes a KubernetesPostgres database (the sandbox memory arm needs nothing, but the registry declares the shape real deployments require).
- `KubernetesKeycloakOperator`
- `KubernetesCloudNativePgOperator` -- 4100–4129: Kubernetes data platforms
- `KubernetesPostgres`
- `KubernetesValkey`
- `KubernetesPerconaMysqlOperator`
- `KubernetesMysql`
- `KubernetesPerconaMongoOperator`
- `KubernetesMongodb`
- `KubernetesStrimziKafkaOperator`
- `KubernetesKafka` -- container_kind: a Strimzi Kafka cluster is a place in the provider's own model — KafkaTopic and KafkaUser declarations BELONG to one cluster (the strimzi.io/cluster label) and are drawn inside its box. Clients that merely talk to the cluster (Connect, MirrorMaker2, UI, Karapace) carry containment_exempt on their bootstrap/trust references.
- `KubernetesKafkaTopic`
- `KubernetesKafkaUser`
- `KubernetesKafkaConnect` -- container_kind: a Connect cluster hosts the connectors deployed INTO it (KafkaConnector's strimzi.io/cluster label names its Connect cluster) — the same room shape as KubernetesKafka above.
- `KubernetesKafkaConnector`
- `KubernetesKafkaMirrorMaker2`
- `KubernetesKarapace`
- `KubernetesKafkaUi`
- `KubernetesOpenSearchOperator`
- `KubernetesOpenSearch`
- `KubernetesAltinityOperator`
- `KubernetesClickHouse`
- `KubernetesSolrOperator`
- `KubernetesSolr`
- `KubernetesNeo4j`
- `KubernetesSeaweedFs`
- `KubernetesQdrant`
- `KubernetesRabbitMqOperator` -- The RabbitMQ Cluster Operator's release manifest ships admission webhooks whose serving certificate is a cert-manager Certificate — cert-manager must be running before the operator installs.
- `KubernetesRabbitMq`
- `KubernetesAirflow` -- 4130–4149: Kubernetes analytics and ML KubernetesPostgres is a prerequisite because Airflow's metadata database composes a KubernetesPostgres by default (the spec's FK defaults resolve onto its outputs) and the migration Job needs the database reachable before the server components start.
- `KubernetesSparkOperator`
- `KubernetesKubeRayOperator`
- `KubernetesRayCluster` -- KubernetesKubeRayOperator is a prerequisite because this kind declares the RayCluster custom resource that only the operator's CRDs admit and only the operator reconciles into head and worker pods.
- `KubernetesFlinkOperator` -- KubernetesCertManager is a prerequisite because the Flink operator's chart, with its default-on admission webhook, renders cert-manager Issuer/Certificate resources and trusts the API server through cert-manager's CA injection — there is no self-signed fallback at the pinned chart, and the webhooks are fail-closed.
- `KubernetesFlinkDeployment` -- KubernetesFlinkOperator is a prerequisite because this kind declares the FlinkDeployment custom resource that only the operator's CRDs admit and only the operator reconciles into a running Flink cluster.
- `KubernetesJupyterHub` -- KubernetesPostgres is a prerequisite because JupyterHub's hub database composes a KubernetesPostgres in its external-database arm (the spec's FK defaults resolve onto its outputs) and the hub pod mounts that database's credential Secret before it can start.
- `KubernetesMlflow` -- KubernetesPostgres is a prerequisite because MLflow's backend store composes a KubernetesPostgres in its production arm (FK defaults onto its outputs; the module composes the connection URI from its credential Secret), and KubernetesSeaweedFs because the artifact store's S3-compatible arm FK-defaults onto the SeaweedFS endpoint and credential Secret.
- `KubernetesTrino` -- KubernetesPostgres is a prerequisite because Trino's postgres catalogs compose a KubernetesPostgres (the catalog host and credential FK-default onto its outputs), and the pods read that database's credential Secret to resolve catalog passwords from environment.
- `KubernetesSuperset` -- KubernetesPostgres is a prerequisite because Superset's REQUIRED metadata database composes a KubernetesPostgres (FK defaults onto its outputs; the module composes the environment Secret from its credential Secret), and KubernetesValkey because the cache/broker arm FK-defaults onto a KubernetesValkey's service and password Secret.
- `KubernetesArgocd` -- 4150–4169: Kubernetes GitOps and CI/CD
- `KubernetesArgoWorkflows`
- `KubernetesTektonOperator`
- `KubernetesTekton` -- KubernetesTektonOperator is a prerequisite because this kind declares the TektonConfig custom resource that only the operator's CRDs admit and only the operator reconciles into running components.
- `KubernetesGhaRunnerScaleSetController`
- `KubernetesGhaRunnerScaleSet` -- KubernetesGhaRunnerScaleSetController is a prerequisite because this kind renders an AutoscalingRunnerSet custom resource that only the controller's CRDs admit and only the controller reconciles into listener and runner pods.
- `KubernetesHarbor`
- `KubernetesJenkins`
- `KubernetesTemporal` -- 4170–4189: Kubernetes app platforms KubernetesPostgres is a prerequisite because the recommended (and E2E-proven) database composition backs Temporal's default and visibility stores with a CloudNativePG cluster.
- `KubernetesNats`
- `KubernetesLocust`
- `DigitalOceanAppPlatformService` -- 5000–5999: DigitalOcean resources
- `DigitalOceanBucket`
- `DigitalOceanContainerRegistry`
- `DigitalOceanDatabaseCluster`
- `DigitalOceanDnsZone`
- `DigitalOceanDroplet`
- `DigitalOceanFirewall`
- `DigitalOceanFunction`
- `DigitalOceanKubernetesCluster`
- `DigitalOceanKubernetesNodePool`
- `DigitalOceanLoadBalancer`
- `DigitalOceanVolume`
- `DigitalOceanVpc`
- `DigitalOceanCertificate`
- `DigitalOceanDnsRecord`
- `CivoBucket` -- 6000–6999: Civo resources
- `CivoCertificate`
- `CivoComputeInstance`
- `CivoDatabase`
- `CivoDnsZone`
- `CivoFirewall`
- `CivoIpAddress`
- `CivoKubernetesCluster`
- `CivoKubernetesNodePool`
- `CivoVolume`
- `CivoVpc`
- `CivoDnsRecord`
- `CloudflareDnsZone` -- 7000–7999: Cloudflare resources
- `CloudflareKvNamespace`
- `CloudflareR2Bucket`
- `CloudflareWorker`
- `CloudflareLoadBalancer`
- `CloudflareD1Database`
- `CloudflareZeroTrustAccessApplication`
- `CloudflareDnsRecord`
- `CloudflareRuleset`
- `CloudflareWorkersKvPair`
- `CloudflareHyperdriveConfig`
- `CloudflareLoadBalancerPool`
- `CloudflareLoadBalancerMonitor`
- `CloudflareZeroTrustAccessPolicy`
- `CloudflareZeroTrustAccessGroup`
- `CloudflareQueue`
- `CloudflarePagesProject`
- `CloudflareZeroTrustTunnel`
- `CloudflareZeroTrustTunnelVirtualNetwork`
- `CloudflareZeroTrustTunnelRoute`
- `CloudflareList`
- `CloudflareListItem`
- `CloudflareTurnstileWidget`
- `CloudflareEmailRoutingZone`
- `CloudflareEmailRoutingRule`
- `CloudflareEmailRoutingAddress`
- `CloudflareOriginCaCertificate`
- `CloudflareCertificatePack`
- `CloudflareCustomHostname`
- `CloudflareCustomHostnameFallbackOrigin`
- `Auth0Connection` -- 8000–8999: Auth0 resources
- `Auth0Client`
- `Auth0EventStream`
- `Auth0ResourceServer`
- `Auth0Action`
- `Auth0Role`
- `OpenFgaStore` -- 9000–9999: OpenFGA resources Note: OpenFGA is Terraform-only - there is no Pulumi provider available. Pulumi modules for OpenFGA resources are pass-through placeholders.
- `OpenFgaAuthorizationModel`
- `OpenFgaRelationshipTuple`
- `OpenStackKeypair` -- 10000–10999: OpenStack resources
- `OpenStackNetwork`
- `OpenStackSubnet`
- `OpenStackRouter`
- `OpenStackRouterInterface`
- `OpenStackSecurityGroup`
- `OpenStackFloatingIp`
- `OpenStackNetworkPort`
- `OpenStackSecurityGroupRule`
- `OpenStackFloatingIpAssociate`
- `OpenStackInstance`
- `OpenStackServerGroup`
- `OpenStackVolume`
- `OpenStackVolumeAttach`
- `OpenStackProject`
- `OpenStackApplicationCredential`
- `OpenStackImage`
- `OpenStackRoleAssignment`
- `OpenStackLoadBalancer`
- `OpenStackLoadBalancerListener`
- `OpenStackLoadBalancerPool`
- `OpenStackLoadBalancerMember`
- `OpenStackLoadBalancerMonitor`
- `OpenStackDnsZone`
- `OpenStackDnsRecord`
- `ScalewayVpc`
- `ScalewayPrivateNetwork`
- `ScalewayPublicGateway`
- `ScalewayLoadBalancer`
- `ScalewayInstanceSecurityGroup`
- `ScalewayInstance`
- `ScalewayKapsuleCluster`
- `ScalewayKapsulePool`
- `ScalewayRdbInstance`
- `ScalewayRedisCluster`
- `ScalewayMongodbInstance`
- `ScalewayObjectBucket`
- `ScalewayBlockVolume`
- `ScalewayContainerRegistry`
- `ScalewayDnsZone`
- `ScalewayDnsRecord`
- `ScalewayServerlessFunction`
- `ScalewayServerlessContainer`
- `AliCloudLogProject`
- `AliCloudRamRole`
- `AliCloudRamPolicy`
- `AliCloudVpc`
- `AliCloudVswitch`
- `AliCloudSecurityGroup`
- `AliCloudEipAddress`
- `AliCloudNatGateway`
- `AliCloudApplicationLoadBalancer`
- `AliCloudNetworkLoadBalancer`
- `AliCloudVpnGateway`
- `AliCloudDnsZone`
- `AliCloudDnsRecord`
- `AliCloudPrivateDnsZone`
- `AliCloudStorageBucket`
- `AliCloudNasFileSystem`
- `AliCloudKmsKey`
- `AliCloudRdsInstance`
- `AliCloudPolardbCluster`
- `AliCloudRedisInstance`
- `AliCloudMongodbInstance`
- `AliCloudEcsInstance`
- `AliCloudContainerRegistry`
- `AliCloudKubernetesCluster`
- `AliCloudKubernetesNodePool`
- `AliCloudCdnDomain`
- `AliCloudFunction`
- `AliCloudSaeApplication`
- `AliCloudRocketmqInstance`
- `AliCloudCenInstance`
- `OciVcn`
- `OciSubnet`
- `OciSecurityGroup`
- `OciCompartment`
- `OciIdentityPolicy`
- `OciDynamicGroup`
- `OciComputeInstance`
- `OciContainerEngineCluster`
- `OciContainerEngineNodePool`
- `OciContainerInstance`
- `OciApplicationLoadBalancer`
- `OciNetworkLoadBalancer`
- `OciDynamicRoutingGateway`
- `OciPublicIp`
- `OciAutonomousDatabase`
- `OciDbSystem`
- `OciMysqlDbSystem`
- `OciPostgresqlDbSystem`
- `OciRedisCluster`
- `OciNosqlTable`
- `OciObjectStorageBucket`
- `OciFileSystem`
- `OciBlockVolume`
- `OciKmsVault`
- `OciKmsKey`
- `OciVaultSecret`
- `OciBastion`
- `OciFunctionsApplication`
- `OciApiGateway`
- `OciStreamPool`
- `OciQueue`
- `OciAlarm`
- `OciLogGroup`
- `OciDnsZone`
- `OciDnsRecord`
- `OciNetworkFirewall`
- `OciDevopsProject`
- `HetznerCloudSshKey`
- `HetznerCloudPlacementGroup`
- `HetznerCloudFirewall`
- `HetznerCloudNetwork`
- `HetznerCloudPrimaryIp`
- `HetznerCloudFloatingIp`
- `HetznerCloudServer`
- `HetznerCloudVolume`
- `HetznerCloudSnapshot`
- `HetznerCloudCertificate`
- `HetznerCloudLoadBalancer`
- `HetznerCloudDnsZone`

### spec.container.app.env.variables[].valueFrom.env

`string`

### spec.container.app.env.variables[].valueFrom.name

`string` · required

- rule: {"required":true}

### spec.container.app.env.variables[].valueFrom.fieldPath

`string`

### spec.container.app.env.variables[].configMapKeyRef

`ConfigMapKeyRef`

Reference to a key in a Kubernetes ConfigMap.

### spec.container.app.env.variables[].configMapKeyRef.name

`string` · required

Name of the ConfigMap.

- rule: {"required":true}

### spec.container.app.env.variables[].configMapKeyRef.key

`string` · required

Key within the ConfigMap.

- rule: {"required":true}

### spec.container.app.env.variables[].configMapKeyRef.optional

`bool`

If true, the env var is silently skipped when the ConfigMap or key does not exist
(instead of blocking pod startup).

### spec.container.app.env.variables[].fieldRef

`ObjectFieldRef`

Reference to a pod-level field (metadata.name, status.podIP, etc.).

### spec.container.app.env.variables[].fieldRef.apiVersion

`string`

Version of the schema. Defaults to "v1".

### spec.container.app.env.variables[].fieldRef.fieldPath

`string` · required

Path of the field to select (e.g., "metadata.name", "status.podIP").

- rule: {"required":true}

### spec.container.app.env.variables[].resourceFieldRef

`ResourceFieldRef`

Reference to container resource limits or requests (limits.cpu, requests.memory, etc.).

### spec.container.app.env.variables[].resourceFieldRef.containerName

`string`

Container name. Required for init containers; defaults to the current
container for regular containers.

### spec.container.app.env.variables[].resourceFieldRef.resource

`string` · required

Resource to select (e.g., "limits.cpu", "requests.memory").

- rule: {"required":true}

### spec.container.app.env.variables[].resourceFieldRef.divisor

`string`

Specifies the output format of the exposed resource.
For CPU: "1" means cores. For memory: "1", "1Ki", "1Mi", "1Gi".

### spec.container.app.env.secrets

`[]SecretEnvVar`

Individual secret environment variables (sensitive).

### spec.container.app.env.secrets[].name

`string` · required

The environment variable name.

- rule: Must be a valid C_IDENTIFIER (start with letter/underscore, contain only letters, digits, underscores)
- rule: {"required":true}

### spec.container.app.env.secrets[].value

`string`

Literal string value.
A Kubernetes Secret is automatically created and the environment variable
references that secret.

### spec.container.app.env.secrets[].secretRef

`KubernetesSecretKeyRef`

Reference to a key within an existing Kubernetes Secret.

### spec.container.app.env.secrets[].secretRef.namespace

`string`

The namespace of the Kubernetes Secret.
If not specified, defaults to the namespace where the component is deployed.
Note: Cross-namespace secret references may not be supported by all Helm charts.

### spec.container.app.env.secrets[].secretRef.name

`string` · required

The name of the Kubernetes Secret.

- rule: {"required":true}

### spec.container.app.env.secrets[].secretRef.key

`string` · required

The key within the Kubernetes Secret that contains the value.

- rule: {"required":true}

### spec.container.app.env.secrets[].secretRef.optional

`bool`

If true, the env var is silently skipped when the Secret or key does not exist
(instead of blocking pod startup).

### spec.container.app.env.secrets[].valueFrom

`ValueFromRef`

Reference to another Planton resource's secret output field.
The orchestrator resolves this before invoking IaC modules.

### spec.container.app.env.secrets[].valueFrom.kind

`enum`

Allowed values (use exactly as shown):

- `unspecified` -- 0: Default/unspecified
- `TestCloudResourceGeneric` -- 1–49: Test/dev/custom
- `TestCloudResourceKubernetes`
- `ConfluentKafka` -- 50–199: saas platform resources
- `AtlasMongodb`
- `SnowflakeDatabase`
- `AwsAlb` -- 1000–1999: AWS resources AwsSubnet is a prerequisite because an ALB requires at least two subnets in different availability zones -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsCertManagerCert`
- `AwsCloudFront`
- `AwsDynamodb`
- `AwsEcrRepo`
- `AwsEcsCluster`
- `AwsEcsService` -- AwsEcsCluster, AwsEcsTaskDefinition, and AwsSubnet are prerequisites because a service schedules a referenced task-definition revision into a referenced live cluster and places task network interfaces into referenced subnets -- all three references must resolve first.
- `AwsEksCluster` -- AwsSubnet and AwsIamRole are prerequisites because the control plane attaches its network interfaces into referenced subnets and assumes a referenced cluster role that must already carry AmazonEKSClusterPolicy.
- `AwsIamRole`
- `AwsLambda`
- `AwsRdsCluster`
- `AwsRdsInstance`
- `AwsRoute53Zone`
- `AwsS3Bucket`
- `AwsLbTargetGroup` -- AwsVpc is a prerequisite because a target group's health checks and target registrations live inside one VPC -- the spec's vpc_id reference must resolve before the group can be created.
- `AwsSecurityGroup` -- AwsVpc is a prerequisite because every security group is created in a VPC; the E2E install profile resolves vpc_id against the VPC prerequisite.
- `AwsVpc`
- `AwsEksNodeGroup` -- AwsEksCluster is a prerequisite because nodes register with a live control plane; AwsIamRole and AwsSubnet back the node role and worker subnet references.
- `AwsIamUser`
- `AwsKmsKey`
- `AwsEc2Instance`
- `AwsClientVpn` -- Every Client VPN endpoint requires an ACM server certificate at create time; the imported self-signed fixture satisfies it. Subnets/VPC are optional composition (a zero-association endpoint is valid) -- composed scenarios declare them via the e2e-prerequisites annotation.
- `AwsDocumentDb`
- `AwsRoute53DnsRecord` -- AwsRoute53Zone is a prerequisite because every record lives inside a hosted zone -- the spec's zone_id reference must resolve before the record can be created.
- `AwsS3ObjectSet` -- AwsS3Bucket is a prerequisite because the object set's bucket reference is required -- objects cannot exist without the bucket that holds them.
- `AwsSqsQueue`
- `AwsSnsTopic`
- `AwsEventBridgeBus`
- `AwsEventBridgeRule`
- `AwsIamOidcProvider`
- `AwsIamPolicy`
- `AwsIamInstanceProfile` -- AwsIamRole is a prerequisite because an instance profile is a wrapper that must contain a role to be useful -- the profile's spec requires a role reference, so the role must be deployed first.
- `AwsLbListener` -- AwsAlb and AwsLbTargetGroup are prerequisites because a listener is an attachment point on a load balancer and its default action almost always forwards to a target group -- both references must resolve before the listener can be created.
- `AwsLbListenerRule` -- AwsLbListener is a prerequisite because a rule only exists as an attachment on a listener -- the listener_arn reference must resolve before the rule can be created.
- `AwsLaunchTemplate`
- `AwsAutoScalingGroup` -- AwsSubnet and AwsLaunchTemplate are prerequisites because a group cannot exist without subnets to place capacity in and a launch template to launch from -- the spec's subnets and launch_template references must resolve before the group can be created.
- `AwsEksAddon` -- AwsEksCluster is a prerequisite because an add-on installs onto a live control plane -- the spec's cluster_name reference must resolve before the add-on can be created.
- `AwsEksFargateProfile` -- AwsEksCluster, AwsIamRole, and AwsSubnet are prerequisites because a Fargate profile attaches to a live control plane, runs pods as a referenced pod-execution role, and launches them into referenced private subnets -- all three references must resolve first.
- `AwsEksAccessEntry` -- AwsEksCluster and AwsIamRole are prerequisites because an access entry grants a referenced IAM principal access to a live control plane -- both references must resolve before the entry can be created.
- `AwsEcsTaskDefinition` -- AwsIamRole is a prerequisite because the kind's default posture -- Fargate with the awslogs logging default -- is rejected by AWS at registration time without an execution role the agent can assume.
- `AwsHttpApiGateway`
- `AwsStepFunction` -- AwsIamRole is a prerequisite because a state machine cannot be created without an execution role it can assume -- the spec's role_arn reference must resolve before the CreateStateMachine call.
- `AwsHttpApiVpcLink` -- AwsSubnet is a prerequisite because a VPC link is a set of managed ENIs provisioned into referenced subnets -- the subnet references must resolve before the link can be created. Security groups are optional on the link, so they compose per-scenario rather than as a registry prerequisite.
- `AwsHttpApiDomain` -- AwsCertManagerCert is a prerequisite because a custom domain cannot be created without a TLS certificate in the same region covering the domain -- the spec's certificate_arn reference must resolve first.
- `AwsVpcEndpoint` -- AwsVpcEndpoint's composed E2E scenarios reference the AwsVpc prerequisite's outputs (vpc_id + default_route_table_id for gateway endpoints) and the AwsSubnet pair's subnet_id outputs (interface endpoints), so both are genuine deploy-order prerequisites.
- `AwsElasticacheUser`
- `AwsElasticacheUserGroup` -- AwsElasticacheUser is a genuine prerequisite: AWS refuses to create a user group that does not contain a user named "default", so a group's composed E2E scenario must resolve a deployed user's outputs.
- `AwsRedshiftServerlessNamespace`
- `AwsRedshiftServerlessWorkgroup` -- The namespace is a genuine prerequisite: a workgroup attaches to exactly one namespace by name at create time, so its composed E2E scenario must resolve a deployed namespace's outputs. AwsSubnet is a prerequisite because Redshift Serverless requires the workgroup's subnets to span three availability zones.
- `AwsRedisElasticache` -- AwsSubnet is a prerequisite because the module builds an ElastiCache subnet group from referenced subnets -- the spec's subnet references must resolve before the replication group can deploy.
- `AwsOpenSearchDomain`
- `AwsMemcachedElasticache`
- `AwsServerlessElasticache`
- `AwsNlb` -- AwsSubnet is a prerequisite because an NLB requires at least one subnet mapping -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsElasticIp`
- `AwsTransitGateway`
- `AwsGlobalAccelerator`
- `AwsSubnet`
- `AwsInternetGateway`
- `AwsNatGateway` -- AwsInternetGateway is a prerequisite because a public NAT gateway can only become available once the VPC it sits in has an internet gateway attached (AWS rejects the create otherwise) -- so the gateway must be deployed first. AwsVpc is a prerequisite because a REGIONAL NAT gateway (availability_mode = regional) references the VPC directly instead of a subnet.
- `AwsEgressOnlyInternetGateway`
- `AwsElasticFileSystem` -- AwsSubnet and AwsSecurityGroup are prerequisites because mount targets (required, min 1) place the file system's NFS endpoints into subnets and attach security groups -- both references must resolve before the CreateMountTarget calls.
- `AwsEfsAccessPoint` -- AwsElasticFileSystem is a prerequisite because an access point is created INTO a file system -- the spec's required file_system_id reference must resolve before the CreateAccessPoint call.
- `AwsFsxLustreFileSystem`
- `AwsFsxOpenzfsFileSystem`
- `AwsFsxWindowsFileSystem` -- Every Windows file system must join an Active Directory domain; the directory itself is external infrastructure (AWS Managed Microsoft AD or a self-managed domain), so only the network dependency is a declarable prerequisite.
- `AwsFsxOntapFileSystem`
- `AwsFsxOntapStorageVirtualMachine`
- `AwsFsxOntapVolume`
- `AwsFsxDataRepositoryAssociation`
- `AwsCognitoUserPool`
- `AwsCognitoIdentityProvider` -- AwsCognitoUserPool is a prerequisite because an identity provider is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateIdentityProvider call.
- `AwsCognitoUserPoolClient` -- AwsCognitoUserPool is a prerequisite because an app client is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateUserPoolClient call.
- `AwsCognitoResourceServer` -- AwsCognitoUserPool is a prerequisite because a resource server is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateResourceServer call.
- `AwsWafWebAcl`
- `AwsWafIpSet`
- `AwsWafRegexPatternSet`
- `AwsCloudwatchLogGroup`
- `AwsCloudwatchAlarm`
- `AwsCloudwatchCompositeAlarm`
- `AwsKinesisStream`
- `AwsKinesisFirehose` -- Every Firehose destination requires an S3 configuration (the primary target for extended_s3; the failed/all-document backup for the rest) and an IAM role Firehose assumes to write to it, so both are hard deploy prerequisites.
- `AwsKinesisStreamConsumer` -- A consumer registers against exactly one stream and cannot exist without it.
- `AwsAthenaWorkgroup`
- `AwsGlueCatalogDatabase`
- `AwsRedshiftCluster`
- `AwsSagemakerDomain` -- AI/ML A domain cannot exist without VPC subnets and a SageMaker execution role (default_user_settings.execution_role_arn is required), so both are hard deploy prerequisites.
- `AwsAppRunnerService` -- A service can run entirely on companion defaults, so the App Runner family's kinds are dependency-free leaves except the VPC connector (which cannot exist without subnets and security groups). A service's companion references (auto scaling / VPC connector / observability / WAF) are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsAppRunnerAutoScalingConfiguration`
- `AwsAppRunnerVpcConnector`
- `AwsAppRunnerObservabilityConfiguration`
- `AwsTransitGatewayVpcAttachment` -- AwsTransitGateway is a prerequisite because an attachment cannot exist without the gateway it attaches to; AwsSubnet because the attachment provisions an ENI into at least one subnet (the VPC arrives transitively through the subnet's own prerequisites).
- `AwsTransitGatewayRouteTable` -- Only the gateway is a hard prerequisite: a route table can exist empty. Associations, propagations, and routes referencing attachments are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsBatchComputeEnvironment` -- A MANAGED compute environment always launches into VPC subnets, so the subnet is a hard deploy prerequisite (security groups are required only for the Fargate types -- scenario-declared, not a registry edge).
- `AwsBatchJobQueue` -- A job queue cannot exist without at least one VALID compute environment to map onto.
- `AwsBatchSchedulingPolicy`
- `AwsBatchJobDefinition`
- `AwsCodeBuildProject` -- CI/CD
- `AwsCodePipeline`
- `AwsMwaaEnvironment` -- Workflow / Orchestration AwsSubnet and AwsSecurityGroup are prerequisites because the environment's network interfaces are placed in referenced private subnets and AWS requires at least one attached security group at creation.
- `AwsNeptuneCluster` -- Graph Database
- `AwsMemorydbCluster` -- A cluster always launches into a subnet group; the subnets are the hard deploy prerequisite. The ACL it attaches is optional composition (the built-in "open-access" ACL needs no resource) -- scenarios declare the ACL/user chain via the e2e-prerequisites annotation.
- `AwsMemorydbUser`
- `AwsMemorydbAcl` -- An empty ACL is valid (MemoryDB has no mandatory "default" member), so the user is optional composition -- the composed scenario declares it via the e2e-prerequisites annotation, never a registry edge.
- `AwsMskCluster` -- Streaming AwsSubnet and AwsSecurityGroup are prerequisites because brokers are placed in referenced subnets and AWS requires at least one attached security group at creation.
- `AwsMskServerlessCluster` -- AwsSubnet is a prerequisite because the serverless cluster's network interfaces are placed in referenced subnets (security groups are optional -- AWS attaches the VPC default group when none are referenced).
- `AwsLambdaEventSourceMapping` -- AwsLambda is a prerequisite because a mapping cannot exist without the function it invokes (a required reference). Event sources (SQS, Kinesis, DynamoDB, MSK) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsSnsSubscription` -- AwsSnsTopic is a prerequisite because a subscription cannot exist without the topic it subscribes to (a required reference). Endpoints (SQS queues, Lambda functions, Firehose streams) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsPlantonRunner` -- AwsSubnet is a prerequisite because the runner appliance places its network interfaces into referenced subnets -- the placement reference must resolve before the appliance can deploy.
- `AwsRoute53HealthCheck`
- `AwsSesConfigurationSet` -- Both SES kinds are dependency-free leaves: an identity's configuration set is optional composition (scenarios declare it via the e2e-prerequisites annotation), and a configuration set's event destinations reference other kinds only optionally.
- `AwsSesEmailIdentity`
- `AwsSecretsManagerSecret` -- A dependency-free leaf: the KMS key, rotation Lambda, and external rotation role references are all optional composition -- scenarios declare them via the e2e-prerequisites annotation, never registry edges.
- `AwsOpenSearchServerlessCollection` -- A dependency-free leaf: the collection-scoped encryption/network/ data-access/retention policies are module-rendered, and the KMS key and data-access principal references are optional composition (e2e-prerequisites annotation).
- `AwsBedrockGuardrail` -- A dependency-free leaf: the KMS key reference is optional composition (e2e-prerequisites annotation); published versions are folded satellites of the guardrail itself.
- `AwsBedrockCustomModel` -- AwsIamRole is a prerequisite because Bedrock assumes the job role to read training data and write outputs; the S3 locations and KMS key are optional composition (e2e-prerequisites annotation).
- `AwsBedrockInferenceProfile` -- A dependency-free leaf: the model source is a foundation model or an AWS system-defined cross-region profile, never a customer resource.
- `AwsBedrockProvisionedThroughput` -- A dependency-free leaf in the registry: capacity is typically bought for an AwsBedrockCustomModel (the default reference), but foundation model ARNs are equally legal, so the edge is optional composition.
- `AwsBedrockModelAccess` -- A dependency-free leaf: the agreement covers an AWS-listed foundation model, never a customer resource.
- `AwsBedrockAgent` -- AwsIamRole is a prerequisite because the Bedrock service assumes the agent resource role to invoke models, action-group Lambdas, and knowledge bases; the guardrail, KMS key, provisioned throughput, and collaborator/knowledge-base edges are optional composition (e2e-prerequisites annotation). Action groups, aliases, collaborators, and knowledge-base associations are folded satellites of the agent.
- `AwsBedrockKnowledgeBase` -- AwsIamRole is a prerequisite because the Bedrock service assumes the knowledge-base role to read data sources, call the embedding model, and read/write the vector store; the vector-store and data-source reference edges (OpenSearch, S3, Secrets Manager, ...) are optional composition (e2e-prerequisites annotation). Data sources are folded satellites of the knowledge base.
- `AwsBedrockFlow` -- AwsIamRole is a prerequisite because the Bedrock service assumes the flow execution role to invoke the models, agents, knowledge bases, and Lambdas its nodes reference; the node-level reference edges are optional composition (e2e-prerequisites annotation).
- `AwsBedrockPrompt` -- A dependency-free leaf: variants target AWS-listed foundation models by ID; targeting another agent's alias is optional composition (e2e-prerequisites annotation).
- `AzureResourceGroup` -- 2000–2999: Azure resources
- `AzureAksCluster` -- AzureResourceGroup is the only required parent: the cluster is created inside a referenced resource group. Subnet is optional on the default node pool (AKS provisions managed networking when unset).
- `AzureAksNodePool` -- AzureAksCluster is a prerequisite because a node pool attaches to an existing cluster by ARM ID; the resource group chains transitively.
- `AzureContainerRegistry` -- AzureResourceGroup is a prerequisite because a container registry is created inside a resource group.
- `AzureDnsZone` -- AzureResourceGroup is a prerequisite because the DNS zone is created inside a referenced resource group that must already exist.
- `AzureKeyVault` -- AzureResourceGroup is a prerequisite because a key vault is created inside a referenced resource group in composed environments.
- `AzureVirtualNetwork` -- AzureResourceGroup is a prerequisite because a virtual network is created inside a referenced resource group in composed environments.
- `AzureNatGateway` -- AzureResourceGroup is a prerequisite because a NAT gateway is created inside a referenced resource group in composed environments.
- `AzureVirtualMachine` -- AzureNetworkInterface is a prerequisite because a virtual machine attaches at least one NIC (the subnet, network, and resource group chain transitively through the NIC's own prerequisites).
- `AzureStorageAccount` -- AzureResourceGroup is a prerequisite because a storage account is created inside a referenced resource group in composed environments.
- `AzureDnsRecord` -- AzureDnsZone is a prerequisite because a record set is created inside a referenced zone (the resource group chains transitively through the zone). Public DNS zone names are not globally unique, so a shared zone fixture is safe to recreate across scenarios.
- `AzureSubnet` -- AzureVirtualNetwork is a prerequisite because a subnet is an ARM child of a referenced network -- the network must exist before the subnet can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzureNetworkSecurityGroup` -- AzureResourceGroup is a prerequisite because a network security group is created inside a referenced resource group in composed environments.
- `AzurePublicIp` -- AzureResourceGroup is a prerequisite because a public IP is created inside a referenced resource group in composed environments.
- `AzurePrivateEndpoint` -- AzureSubnet is a prerequisite because a private endpoint draws its private IP from a referenced subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite). The connection target is polymorphic and the DNS zones / ASGs are optional, so none of those are prerequisites.
- `AzurePrivateDnsZone` -- AzureResourceGroup is a prerequisite because a private DNS zone is created inside a referenced resource group in composed environments.
- `AzureApplicationGateway` -- AzureSubnet is a prerequisite because a gateway cannot exist without its dedicated gateway_ip_configuration subnet (the network and resource group chain transitively through the subnet's own prerequisites); public frontends additionally reference a public IP, but private-only gateways are legal, so it is not a registry prerequisite.
- `AzureLoadBalancer` -- AzureResourceGroup is a prerequisite because a load balancer is created inside a referenced resource group (frontends additionally reference subnets or public IPs, but neither is universally required, so they are not registry prerequisites).
- `AzureRouteTable` -- AzureResourceGroup is a prerequisite because a route table is created inside a referenced resource group in composed environments.
- `AzurePrivateDnsZoneVirtualNetworkLink` -- AzurePrivateDnsZone and AzureVirtualNetwork are prerequisites because a virtual network link is a child resource of a referenced zone and binds it to a referenced network -- both must exist before the link can be written. (The resource group arrives transitively through the zone's and network's own prerequisite declarations.)
- `AzureVirtualNetworkPeering` -- AzureVirtualNetwork is a prerequisite because a peering is an ARM child of its local network and binds it to a remote network -- the local network must exist before the peering can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzurePublicIpPrefix` -- AzureResourceGroup is a prerequisite because a public IP prefix is created inside a referenced resource group in composed environments.
- `AzureNetworkInterface` -- AzureSubnet is a prerequisite because a network interface's IP configurations deploy into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureManagedDisk` -- AzureResourceGroup is a prerequisite because a managed disk is created inside a resource group.
- `AzureVirtualMachineScaleSet` -- AzureSubnet is a prerequisite because every scale-set instance's network interface deploys into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureKeyVaultKey` -- AzureKeyVault is a prerequisite because a key is a data-plane object inside a referenced vault -- the vault must exist before the key can be written (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultCertificate` -- AzureKeyVault is a prerequisite because a certificate is a data-plane object inside a referenced vault -- the vault must exist before the certificate can be enrolled or imported (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultSecret` -- AzureKeyVault is a prerequisite because a secret is a data-plane object inside a referenced vault -- the vault must exist before the secret can be written (the resource group chains transitively through the vault's own prerequisite). Part of the Key Vault family (2005, 2025-2026) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureWebApplicationFirewallPolicy` -- AzureResourceGroup is a prerequisite because a WAF policy is created inside a referenced resource group; the Application Gateways that attach the policy reference it, never the reverse.
- `AzureApplicationSecurityGroup` -- AzureResourceGroup is a prerequisite because an application security group is created inside a referenced resource group; network interfaces, scale-set IP configurations, and NSG security rules reference the group, never the reverse.
- `AzureDiskEncryptionSet` -- AzureKeyVaultKey is a prerequisite because a disk encryption set wraps customer data with a referenced key -- the key (and its vault, which chains transitively) must exist before the set can resolve the key URL at create time.
- `AzurePostgresqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureRedisCache` -- AzureResourceGroup is a prerequisite because the cache is created inside a referenced resource group (VNet injection additionally references a dedicated subnet, but only the Premium tier supports it, so it is not a registry prerequisite).
- `AzureCosmosdbAccount` -- AzureResourceGroup is a prerequisite because the account is created inside a referenced resource group.
- `AzureMssqlServer` -- AzureResourceGroup is a prerequisite because the logical server is created inside a referenced resource group.
- `AzureMysqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureMssqlDatabase` -- The parent logical server is referenced via server_id, not auto-deployed: E2E scenarios declare their own server fixture (minimal-server.yaml or the pool-attach chain through AzureMssqlElasticPool) so sequential subtests never destroy and recreate the same globally unique server_name.
- `AzureMssqlElasticPool` -- AzureMssqlServer is a prerequisite because every elastic pool lives on a referenced logical server (the server's resource group is transitive).
- `AzureRedisLinkedServer` -- The target and linked caches are referenced via ARM ids, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicy` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicyAssignment` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureContainerAppEnvironment` -- AzureResourceGroup is a prerequisite because the environment is created inside a referenced resource group that must already exist.
- `AzureContainerApp` -- AzureContainerAppEnvironment is a prerequisite because every app runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureServicePlan` -- AzureResourceGroup is a prerequisite because the plan is created inside a referenced resource group that must already exist.
- `AzureFunctionApp` -- AzureServicePlan is a prerequisite because a function app runs on a referenced plan (the resource group arrives transitively through the plan). The required storage account is deliberately NOT a registry prerequisite: storage-account names are globally unique, so scenarios bring their own scenario-local account fixtures.
- `AzureLinuxWebApp` -- AzureServicePlan is a prerequisite because a web app runs on a referenced plan (the resource group arrives transitively through the plan).
- `AzureContainerAppJob` -- AzureContainerAppEnvironment is a prerequisite because a job runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureContainerAppEnvironmentStorage` -- AzureContainerAppEnvironment is a prerequisite because the storage registration lives on a referenced environment. The Azure Files share and storage account are deliberately NOT registry prerequisites: storage-account names are globally unique, so scenarios bring their own scenario-local account + share fixtures.
- `AzureContainerAppEnvironmentDaprComponent` -- AzureContainerAppEnvironment is a prerequisite because the Dapr component is registered on a referenced environment.
- `AzureContainerAppEnvironmentCertificate` -- AzureContainerAppEnvironment is a prerequisite because the certificate is stored on a referenced environment.
- `AzureContainerAppEnvironmentManagedCertificate` -- AzureContainerAppEnvironment is a prerequisite because the managed certificate is provisioned on a referenced environment.
- `AzureLogAnalyticsWorkspace` -- AzureResourceGroup is a prerequisite because the workspace is created inside a referenced resource group that must already exist.
- `AzureApplicationInsights` -- AzureLogAnalyticsWorkspace is a prerequisite because workspace-based Application Insights stores its telemetry in a referenced workspace (the resource group chains transitively through the workspace).
- `AzureMonitorDiagnosticSetting` -- AzureLogAnalyticsWorkspace is a prerequisite because the setting's scenarios route a fixture workspace's telemetry (the workspace doubles as target and destination); the target itself is polymorphic.
- `AzureMonitorActionGroup` -- AzureResourceGroup is a prerequisite because the action group is created inside a referenced resource group that must already exist.
- `AzureMonitorMetricAlert` -- AzureMonitorActionGroup is a prerequisite because a metric alert's actions fire into a referenced action group (the resource group chains transitively); alert scopes are polymorphic.
- `AzureMonitorScheduledQueryAlert` -- AzureLogAnalyticsWorkspace is a prerequisite because the rule queries a referenced workspace scope; AzureMonitorActionGroup because its action fires into a referenced action group.
- `AzureMonitorActivityLogAlert` -- AzureMonitorActionGroup is a prerequisite because an activity log alert's actions fire into a referenced action group (the resource group chains transitively). The alert itself is subscription-global and its scopes are polymorphic.
- `AzureApplicationInsightsStandardWebTest` -- AzureApplicationInsights is a prerequisite because a standard web test binds to a referenced Application Insights component (the resource group chains transitively through the component).
- `AzureUserAssignedIdentity` -- AzureResourceGroup is a prerequisite because the identity is created inside a referenced resource group that must already exist.
- `AzureRoleAssignment` -- AzureResourceGroup and AzureUserAssignedIdentity are prerequisites because an assignment grants a role at a referenced scope (most commonly a resource group) to a referenced principal (most commonly a managed identity) -- both must exist before the grant can be written.
- `AzureRoleDefinition` -- AzureResourceGroup is a prerequisite because a custom role definition is created at a referenced scope, most commonly a resource group in composed environments -- the scope must exist before the definition can be written.
- `AzureFederatedIdentityCredential` -- AzureUserAssignedIdentity is the prerequisite because a federated identity credential is a child resource of a referenced managed identity -- the identity must exist before the credential can be written on it. (The resource group arrives transitively through the identity's own prerequisite declaration.)
- `AzureServiceBusNamespace` -- AzureResourceGroup is a prerequisite because a Service Bus namespace is created inside a referenced resource group in composed environments. The namespace is the container every Service Bus messaging entity (queue, topic, subscription, authorization rule, geo-DR pairing) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureEventHubNamespace` -- AzureResourceGroup is a prerequisite because an Event Hub namespace is created inside a referenced resource group in composed environments. The namespace is the container every Event Hubs entity (event hub, consumer group, authorization rule, schema group, geo-DR pairing, customer-managed key) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureServiceBusQueue`
- `AzureServiceBusTopic`
- `AzureServiceBusSubscription`
- `AzureServiceBusAuthorizationRule`
- `AzureServiceBusDisasterRecoveryConfig`
- `AzureEventHub`
- `AzureEventHubConsumerGroup`
- `AzureEventHubAuthorizationRule`
- `AzureFrontDoorProfile` -- AzureResourceGroup is a prerequisite because a Front Door profile is created inside a referenced resource group in composed environments. The profile is the container every Front Door delivery resource (endpoint, origin group, origin, route) nests under.
- `AzureFrontDoorEndpoint` -- AzureFrontDoorProfile is a prerequisite because an endpoint is an ARM child of a referenced profile -- the profile must exist before the endpoint can be written. (The resource group arrives transitively through the profile's own prerequisite declaration.)
- `AzureFrontDoorOriginGroup` -- AzureFrontDoorProfile is a prerequisite because an origin group is an ARM child of a referenced profile.
- `AzureFrontDoorOrigin` -- AzureFrontDoorOriginGroup is a prerequisite because an origin is an ARM child of a referenced origin group (the profile and resource group chain transitively).
- `AzureFrontDoorRoute` -- A route attaches to an endpoint (its ARM parent) and forwards to an origin group whose origins must exist before ARM accepts the route -- so both the endpoint and the origin chain are genuine deploy-order prerequisites.
- `AzureFrontDoorRuleSet` -- AzureFrontDoorProfile is a prerequisite because a rule set is an ARM child of a referenced profile. The rules live inside the set (they form one ordered policy document); routes attach the set by ARM ID.
- `AzureFrontDoorCustomDomain` -- AzureFrontDoorProfile is a prerequisite because a custom domain is an ARM child of a referenced profile. The DNS zone and (for bring-your-own certificates) the Front Door secret are optional references, not deploy-order prerequisites.
- `AzureFrontDoorSecret` -- AzureFrontDoorSecret is a prerequisite-light kind: only the profile (its ARM parent) must exist. The Key Vault certificate it wraps is a reference resolved before the module runs; its vault chain is exercised through scenario-local fixtures in E2E.
- `AzureFrontDoorFirewallPolicy` -- AzureResourceGroup is a prerequisite because the Front Door WAF policy is created inside a referenced resource group -- it is a GLOBAL resource, not a profile child (a different ARM type than the regional Application Gateway WAF policy). Security policies attach it to profiles; the policy itself depends on nothing else.
- `AzureFrontDoorSecurityPolicy` -- A security policy is an ARM child of a profile that associates a referenced WAF policy with referenced domains -- so the endpoint (the default-domain association target; the profile arrives transitively through it) and the WAF policy are genuine deploy-order prerequisites.
- `AzureStorageContainer` -- None of the storage data-service kinds declares a registry prerequisite on AzureStorageAccount: account names are GLOBALLY unique and Azure holds a just-deleted name, so a recreate-per-scenario fixture would hang -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the storage_account_id reference itself.
- `AzureStorageShare`
- `AzureStorageQueue`
- `AzureStorageTable`
- `AzureStorageEncryptionScope`
- `AzureStorageDataLakeGen2Filesystem`
- `AzureStorageLocalUser`
- `AzureStorageObjectReplication`
- `AzureCosmosdbSqlDatabase` -- None of the Cosmos DB data-service kinds declares a registry prerequisite on AzureCosmosdbAccount: account names are GLOBALLY unique DNS labels, so a recreate-per-scenario fixture would risk name-reuse hangs -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the cosmosdb_account_id / parent-database references themselves.
- `AzureCosmosdbSqlContainer`
- `AzureCosmosdbMongoDatabase`
- `AzureCosmosdbMongoCollection`
- `AzureCosmosdbSqlRoleDefinition`
- `AzureCosmosdbSqlRoleAssignment`
- `AzureManagedRedis` -- AzureResourceGroup is the cluster's only registry prerequisite: the cluster is created inside a referenced resource group. The geo-replication and access-policy-assignment children declare NO prerequisite on AzureManagedRedis: clusters are expensive, slow-provisioning parents, so their E2E scenarios declare scenario-local cluster fixtures instead of recreating a shared one per scenario. Deploy ordering in composed environments still flows from the managed_redis_id references themselves.
- `AzureManagedRedisGeoReplication`
- `AzureManagedRedisAccessPolicyAssignment`
- `AzureEventHubDisasterRecoveryConfig`
- `AzureEventHubSchemaGroup`
- `AzureEventHubCluster` -- AzureResourceGroup is a prerequisite because a dedicated Event Hubs cluster is created inside a referenced resource group in composed environments. Note: clusters cannot be deleted for 4 hours after creation (Azure's moratorium), so E2E treats this kind as offline-gated.
- `AzureEventHubNamespaceCustomerManagedKey`
- `AzureMssqlFailoverGroup` -- AzureMssqlServer is a prerequisite because a failover group is created on a referenced primary logical server and points at a partner server; the primary (and its resource group, which chains transitively) must exist before the group can be written.
- `AzureContainerAppCustomDomain` -- AzureContainerApp is a prerequisite because the domain binding lives in a referenced app's ingress configuration (the environment and resource group chain transitively through the app).
- `AzureFirewallPolicy`
- `AzureFirewallPolicyRuleCollectionGroup` -- AzureFirewallPolicy is a prerequisite because a rule collection group is a child document of a referenced policy (the resource group chains transitively through the policy).
- `AzureFirewall` -- AzureSubnet is a prerequisite because a VNet-deployed firewall's data path lives in a dedicated subnet that must be named exactly "AzureFirewallSubnet" (the virtual network and resource group chain transitively through the subnet). The E2E install profile publishes a fixture subnet with that exact name and a /26 prefix.
- `AzureIpGroup`
- `AzureVirtualNetworkGateway` -- AzureSubnet is a prerequisite because every virtual network gateway lives in a dedicated subnet named exactly "GatewaySubnet" (the virtual network and resource group chain transitively through the subnet); the subnet install profile publishes a fixture instance with that exact ARM name. AzurePublicIp is a prerequisite because a VPN-type gateway (the default shape) requires a public IP per ip configuration; the address install profile publishes a dedicated zone-redundant instance (a gateway binds its address exclusively, and the AZ gateway SKUs require zones on it).
- `AzureVirtualNetworkGatewayConnection` -- Both gateways are prerequisites: a connection joins a virtual network gateway to a far side, and the site-to-site far side is a local network gateway (the GatewaySubnet, VNet, and resource group chain transitively through the virtual network gateway).
- `AzureLocalNetworkGateway`
- `AzurePrivateLinkService` -- AzureSubnet is the sole prerequisite: every NAT ip configuration draws its address from a subnet with private-link-service network policies disabled (the subnet install profile publishes a fixture instance with that flag off). The Standard load balancer whose frontend the service typically fronts is NOT a registry prerequisite -- the spec's destination is an exactly-one-of (load balancer frontend OR fixed destination IP), so scenarios that use the load-balancer shape declare it via the planton.dev/e2e-prerequisites annotation instead.
- `AzureExpressRouteCircuit`
- `AzureExpressRouteCircuitPeering` -- The circuit is the prerequisite: a peering is an ARM child of the circuit, addressed by the circuit's name (the resource group chains transitively through the circuit).
- `AzureExpressRouteGateway` -- The hub is the prerequisite: ARM requires an ExpressRoute Gateway to be deployed INTO a Virtual WAN hub (the WAN and resource group chain transitively through the hub).
- `AzureExpressRoutePort` -- ExpressRoute Port: your own physical port pair on a Microsoft edge router (ExpressRoute Direct), from whose bandwidth circuits are carved. Self-contained -- only the resource group is required.
- `AzureVirtualWan` -- Virtual WAN: the umbrella of Azure's managed hub-and-spoke networking, under which virtual hubs and their gateways are created. Self-contained -- only the resource group is required.
- `AzureVirtualHub` -- The WAN is the prerequisite: this kind models the Virtual WAN hub (virtual_wan_id is required; standalone hubs are the legacy Route Server construction, which has its own ARM surface). The resource group chains transitively through the WAN.
- `AzureVirtualHubConnection` -- Both sides of the attachment are prerequisites: the hub being joined and the spoke virtual network being attached.
- `AzureVpnGateway` -- The hub is the prerequisite: ARM deploys a Virtual WAN VPN gateway INTO a virtual hub (virtual_hub_id is required and immutable; the WAN and resource group chain transitively through the hub). ARM allows one VPN gateway per hub.
- `AzureVpnGatewayConnection` -- Both ends of the tunnel are prerequisites: a connection is an ARM child of the VPN gateway and pins each of its links to a specific link of the remote VPN site (the hub, WAN, and resource group chain transitively through the gateway).
- `AzureVpnSite` -- The WAN is the prerequisite: a VPN site is the Virtual WAN world's address-book entry for one branch location (virtual_wan_id is required; the classic-world sibling without a WAN is AzureLocalNetworkGateway). The resource group chains transitively through the WAN.
- `AzurePointToSiteVpnGateway` -- The hub and the server configuration are both prerequisites: a point-to-site VPN gateway deploys INTO a virtual hub (one P2S gateway per hub, a slot separate from the hub's site-to-site VPN gateway) and is born pointing at the VPN server configuration that defines how its users authenticate -- both ARM-required and fixed at creation. The WAN and resource group chain transitively through the hub.
- `AzureVpnServerConfiguration` -- Self-contained -- only the resource group is required: a VPN server configuration is the reusable "who may connect and how" authentication policy (Entra ID / certificate / RADIUS) that point-to-site VPN gateways attach to; it references no other Azure resource.
- `AzureCognitiveAccount` -- Self-contained -- only the resource group is required: an Azure AI services account (Azure OpenAI, the multi-service AIServices account, the single-service accounts) needs no other Azure resource; subnets (network rules), Key Vault keys (CMK), storage accounts and user-assigned identities are optional references.
- `AzureCognitiveDeployment` -- An ARM child of its account: a model deployment (which model runs, at which throughput class) exists only on an Azure AI services account of kind "OpenAI" or "AIServices".
- `AzureCognitiveAccountProject` -- An ARM child of its account: an AI Foundry project exists only on an "AIServices"-kind account with project management enabled.
- `AzureMachineLearningWorkspace` -- The workspace REQUIRES all three companion services at creation (default storage, secrets vault, telemetry) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureMachineLearningDatastore` -- An ARM child of its workspace. The storage target (container, filesystem or share) is scenario-declared via the e2e-prerequisites annotation -- only the blob scenario needs a container, so it is not a kind-wide prerequisite.
- `AzureMachineLearningComputeCluster` -- An ARM child of its workspace (.../computes/{name}) -- the auto-scaling pool of VMs training jobs run on.
- `AzureMachineLearningComputeInstance` -- An ARM child of its workspace (.../computes/{name}) -- a single always-on VM serving as one data scientist's cloud workstation.
- `AzureAiFoundry` -- The hub REQUIRES both companion services at creation (secrets vault, default storage) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureAiFoundryProject` -- Deploys into its hub's resource group (the provider derives the group from the hub reference -- the project spec carries none).
- `AzureSearchService`
- `AzureMachineLearningOnlineEndpoint` -- An ARM child of its workspace (.../onlineEndpoints/{name}) -- the stable scoring address applications call. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningOnlineDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the running copy of a model the endpoint's traffic map routes to.
- `AzureMachineLearningBatchEndpoint` -- An ARM child of its workspace (.../batchEndpoints/{name}) -- the stable address batch scoring jobs are submitted to. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningBatchDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the job recipe (model, compute, batching behavior) the endpoint's default-deployment pointer routes submissions to.
- `AzureRecoveryServicesVault` -- The Recovery Services vault (Microsoft.RecoveryServices/vaults) -- the safe that classic Azure Backup data and Site Recovery configuration live in. Backup policies and protected items are ARM children of a vault.
- `AzureBackupPolicyVm` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern IaaS VM backups.
- `AzureBackupProtectedVm` -- An ARM child of its vault (.../protectedItems/...) -- the binding that puts one virtual machine under a backup policy's protection.
- `AzureBackupPolicyFileShare` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern Azure Files share backups (snapshot or vaulted).
- `AzureBackupProtectedFileShare` -- An ARM child of its vault (.../protectedItems/AzureFileShare;...) -- the binding that puts one Azure Files share under a backup policy's protection. The share's storage account must already be registered with the vault (AzureBackupContainerStorageAccount).
- `AzureDataProtectionBackupVault` -- The Data Protection backup vault (Microsoft.DataProtection/ backupVaults) -- the safe that MODERN Azure Backup data lives in (managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, Data Lake storage). Backup policies and backup instances are ARM children of a vault.
- `AzureDataProtectionBackupPolicy` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules for ONE Data Protection datasource type (blob storage, disk, Kubernetes cluster, MySQL/PostgreSQL flexible server, or Data Lake storage), modeled as one kind with variant blocks.
- `AzureDataProtectionBackupInstance` -- An ARM child of its vault (.../backupInstances/{name}) -- the binding that puts ONE datasource (a managed disk, a storage account's blob services, an AKS cluster, a MySQL/PostgreSQL flexible server, or a Data Lake storage account) under a Data Protection backup policy, modeled as one kind with variant blocks. The vault's managed identity must hold the datasource roles Azure Backup requires BEFORE the instance is created.
- `AzureBastionHost` -- AzureSubnet and AzurePublicIp are prerequisites because a dedicated-infrastructure Bastion host (Basic/Standard/Premium -- the default shapes) deploys into a subnet named exactly "AzureBastionSubnet" and binds a Standard static public IP EXCLUSIVELY (the virtual network and resource group chain transitively through the subnet). The Developer SKU instead attaches to a virtual network directly and uses neither.
- `AzureNetworkWatcherFlowLog` -- AzureVirtualNetwork and AzureStorageAccount are prerequisites because a flow log records a network-scoped target (a virtual network in the common case; subnets and network interfaces chain through the network) into a referenced storage account. The regional Network Watcher parent is NOT a prerequisite: Azure auto-creates it ("NetworkWatcher_{region}" in "NetworkWatcherRG") the moment the region hosts a virtual network, and the flow log references it by name. Traffic Analytics' Log Analytics workspace is an optional arm, declared by scenarios that use it.
- `AzurePrivateDnsResolver` -- AzureVirtualNetwork and AzureSubnet are prerequisites because a DNS Private Resolver anchors to a referenced virtual network (at most ONE resolver per network -- Azure enforces it) and each of its inbound/outbound endpoints occupies its own dedicated subnet delegated to "Microsoft.Network/dnsResolvers" (the resource group chains transitively through the network and subnets).
- `AzurePrivateDnsResolverForwardingRuleset` -- AzurePrivateDnsResolver is a prerequisite because a DNS forwarding ruleset steers a resolver's OUTBOUND endpoints -- it binds their ARM ids (at most 2, same resolver) at creation. (The resource group and network chain transitively through the resolver's own prerequisite declarations.)
- `AzureBackupContainerStorageAccount` -- Registers a storage account with a Recovery Services vault as a backup container (.../protectionContainers/StorageContainer;...) -- one registration per storage-account-and-vault pair, required BEFORE any of the account's file shares can be protected. Part of the backup family (2175-2179) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureDataProtectionResourceGuard` -- The Data Protection Resource Guard (Microsoft.DataProtection/ resourceGuards) -- the approval gate behind Multi-User Authorization: privileged vault operations (disabling soft delete, reducing retention) require an approval through a guard, which typically lives in a DIFFERENT administrator's scope. Vaults reference a guard by its ARM ID. Part of the backup family (2175-2182) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzurePrivateDnsResolverVirtualNetworkLink` -- The attachment that makes a DNS forwarding ruleset take effect in one virtual network ({ruleset_id}/virtualNetworkLinks/{name}) -- one link per ruleset-network pair, up to 500 per ruleset, spokes joining and leaving independently (which is why the link is a standalone kind, exactly like AzurePrivateDnsZoneVirtualNetworkLink). Part of the DNS Private Resolver family (2186-2187) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `GcpArtifactRegistryRepo` -- 3000–3999: GCP resources
- `GcpTargetHttpsProxy` -- The URL map is the parent a proxy cannot exist without; the classic compute certificate kinds and the SSL policy are the fixture parents the committed scenarios attach. The Certificate Manager certificate list (certificate_manager_certificates, honored only by the cross-region internal ALB) is optional composition -- a scenario that arms it declares GcpCertManagerCert via the e2e-prerequisites annotation, never a registry edge that would tax every proxy and forwarding-rule chain.
- `GcpCloudFunction`
- `GcpCloudRun`
- `GcpCloudSql`
- `GcpDnsZone`
- `GcpGcsBucket`
- `GcpGkeCluster`
- `GcpIamCustomRole`
- `GcpProject`
- `GcpVpcNetwork`
- `GcpSubnetwork`
- `GcpRouterNat`
- `GcpGkeNodePool`
- `GcpServiceAccount`
- `GcpGkeWorkloadIdentityBinding`
- `GcpCertManagerCert`
- `GcpComputeInstance`
- `GcpDnsRecord`
- `GcpProjectIamMember`
- `GcpFirewallRule`
- `GcpGlobalAddress`
- `GcpCloudArmorPolicy`
- `GcpHealthCheck`
- `GcpBackendBucket`
- `GcpBackendService`
- `GcpRegionNetworkEndpointGroup`
- `GcpUrlMap`
- `GcpManagedSslCertificate`
- `GcpTargetHttpProxy`
- `GcpAlloydbCluster`
- `GcpRedisInstance`
- `GcpFirestoreDatabase`
- `GcpSpannerInstance`
- `GcpSpannerDatabase`
- `GcpBigtableInstance`
- `GcpMemorystoreInstance`
- `GcpCloudSqlDatabase`
- `GcpCloudSqlUser`
- `GcpAlloydbInstance`
- `GcpAlloydbUser`
- `GcpSpannerBackupSchedule`
- `GcpBigtableTable`
- `GcpFirestoreBackupSchedule`
- `GcpFirestoreIndex`
- `GcpBigQueryDataset`
- `GcpDataprocCluster`
- `GcpDataprocAutoscalingPolicy`
- `GcpBigQueryTable`
- `GcpPubSubTopic`
- `GcpPubSubSubscription`
- `GcpCloudTasksQueue`
- `GcpCloudSchedulerJob`
- `GcpPubSubSchema`
- `GcpVertexAiNotebook`
- `GcpVertexAiEndpoint`
- `GcpVertexAiIndex`
- `GcpVertexAiIndexEndpoint` -- Vector Search IndexEndpoint — distinct from the online-prediction GcpVertexAiEndpoint (671); different GCP resources, different kinds.
- `GcpVertexAiDeployedIndex`
- `GcpCloudComposerEnvironment`
- `GcpCloudComposerUserWorkloadsSecret`
- `GcpCloudComposerUserWorkloadsConfigMap`
- `GcpKmsKeyRing`
- `GcpKmsKey`
- `GcpKmsKeyIamMember`
- `GcpFilestoreInstance`
- `GcpWorkloadIdentityPool` -- 3101–3109: IAM/identity family (overflow block; the 3000–3022 foundation/security sub-band is fully allocated)
- `GcpWorkloadIdentityPoolProvider`
- `GcpServiceAccountIamMember`
- `GcpGlobalForwardingRule` -- 3110–3119: networking/load-balancer family (overflow block; the 3023–3029 LB sub-band is fully allocated)
- `GcpSslPolicy`
- `GcpSslCertificate`
- `GcpServiceNetworkingConnection`
- `GcpAddress`
- `GcpServiceConnectionPolicy`
- `GcpCertManagerDnsAuthorization`
- `GcpCertificateMap` -- GcpCertManagerCert is a prerequisite because a map entry binds hostnames to EXISTING certificates — the canonical map references a certificate fixture's resource name.
- `GcpCloudRunJob` -- 3120–3129: GCP serverless overflow
- `GcpServerlessVpcConnector`
- `GcpComputeDisk` -- 3130–3139: GCP compute overflow (the 3000–3022 foundation sub-band that holds GcpComputeInstance is fully allocated)
- `GcpComputeMig` -- GcpVpcNetwork is a prerequisite because the canonical group runs its fleet on a dedicated custom-mode VPC — a managed instance group's template must attach every VM to a network, and the default VPC is never assumed.
- `GcpMonitoringNotificationChannel` -- 3140–3149: GCP observability & log routing
- `GcpMonitoringAlertPolicy` -- GcpMonitoringNotificationChannel is a prerequisite because the policy's canonical shape references a channel to notify — a policy without a delivery endpoint measures but never pages.
- `GcpMonitoringUptimeCheck`
- `GcpLoggingSink` -- GcpGcsBucket is a prerequisite because the canonical sink exports to a Cloud Storage bucket — the cheapest destination that proves the whole writer-identity grant flow.
- `GcpMonitoringDashboard`
- `GcpMonitoringSlo`
- `GcpLogBucket`
- `GcpLogMetric`
- `GcpSecretManagerSecret` -- 3150–3159: GCP security & identity GcpServiceAccount is a prerequisite because the canonical secret grants secretAccessor to a workload service account — the access story the kind exists to model.
- `GcpIdentityPlatformConfig`
- `GcpIdentityPlatformTenant` -- GcpIdentityPlatformConfig is a prerequisite because tenants exist only in projects whose Identity Platform config enables multi_tenant.allow_tenants — a tenant without the initialized, tenant-enabled project config cannot be created at all.
- `GcpIamOauthClient`
- `GcpIamDenyPolicy`
- `GcpCloudRunDomainMapping` -- 3160–3169: GCP serverless edge GcpCloudRun is a prerequisite because a domain mapping exists only to point a verified domain at a running Cloud Run service — the route it maps must already exist for the mapping to be created at all.
- `GcpWorkflow`
- `GcpEventarcTrigger` -- GcpCloudRun is a prerequisite because the canonical trigger routes a Pub/Sub messagePublished event to a Cloud Run service — the destination story the kind exists to model.
- `GcpEventarcMessageBus`
- `KubernetesNamespace` -- 4000–4999: Kubernetes resources, organized in family sub-bands (4030–4069 also hosts CNI/autoscaling/DR addons; 4130–4149 hosts analytics & ML; 4190–4199 reserved for growth) 4000–4029: Kubernetes building blocks (core API primitives)
- `KubernetesDeployment`
- `KubernetesStatefulSet`
- `KubernetesDaemonSet`
- `KubernetesJob`
- `KubernetesCronJob`
- `KubernetesService`
- `KubernetesSecret`
- `KubernetesManifest`
- `KubernetesHelmRelease`
- `KubernetesConfigMap`
- `KubernetesServiceAccount`
- `KubernetesRbac` -- Bundles the RBAC grant grain (Role/ClusterRole + its binding) into one component: "grant these permissions to these subjects in this scope".
- `KubernetesIngress`
- `KubernetesNetworkPolicy`
- `KubernetesPersistentVolumeClaim`
- `KubernetesStorageClass`
- `KubernetesResourceQuota` -- Manages the namespace-governance pair: the ResourceQuota plus an optional companion LimitRange (per-object defaults/bounds) — two API objects, one governance story.
- `KubernetesPriorityClass`
- `KubernetesPodDisruptionBudget`
- `KubernetesHorizontalPodAutoscaler`
- `KubernetesCertManager` -- 4030–4069: Kubernetes foundation addons (certs, DNS, secrets, ingress, Gateway API, mesh, CNI/autoscaling/DR)
- `KubernetesClusterIssuer` -- KubernetesCertManager is a prerequisite for the three cert-manager CR kinds below: ClusterIssuer/Issuer/Certificate are cert-manager custom resources — without the controller and its CRDs they cannot be applied.
- `KubernetesIssuer`
- `KubernetesCertificate`
- `KubernetesExternalDns`
- `KubernetesExternalSecretsOperator`
- `KubernetesClusterSecretStore` -- KubernetesExternalSecretsOperator is a prerequisite for the three external-secrets CR kinds below: ClusterSecretStore/SecretStore/ ExternalSecret are external-secrets custom resources — without the operator and its CRDs they cannot be applied.
- `KubernetesSecretStore`
- `KubernetesExternalSecret`
- `KubernetesIngressNginx`
- `KubernetesGatewayApiCrds`
- `KubernetesGatewayClass`
- `KubernetesGateway`
- `KubernetesListenerSet`
- `KubernetesHttpRoute`
- `KubernetesGrpcRoute`
- `KubernetesTcpRoute`
- `KubernetesUdpRoute`
- `KubernetesTlsRoute`
- `KubernetesReferenceGrant`
- `KubernetesBackendTlsPolicy`
- `KubernetesIstioBaseCrds`
- `KubernetesIstio`
- `KubernetesDestinationRule` -- Istio API components (mesh traffic policy, security, telemetry). The seven typed resources below (4053–4059) require the Istio CRDs on the cluster, provided by the lightweight CRDs-only KubernetesIstioBaseCrds (851) — NOT the full mesh KubernetesIstio (852).
- `KubernetesServiceEntry`
- `KubernetesPeerAuthentication`
- `KubernetesRequestAuthentication`
- `KubernetesAuthorizationPolicy`
- `KubernetesTelemetry`
- `KubernetesEnvoyFilter`
- `KubernetesMetricsServer`
- `KubernetesCilium`
- `KubernetesKeda`
- `KubernetesKarpenter`
- `KubernetesKarpenterNodePool`
- `KubernetesKarpenterEc2NodeClass`
- `KubernetesClusterAutoscaler`
- `KubernetesVelero`
- `KubernetesKubePrometheusStack` -- 4070–4089: Kubernetes observability
- `KubernetesGrafana`
- `KubernetesSignoz` -- KubernetesClickHouse is a prerequisite because SigNoz stores every trace, metric and log in ClickHouse and deploys none of its own — the telemetry store is composed, never bundled.
- `KubernetesLoki`
- `KubernetesTempo`
- `KubernetesOtelOperator` -- The operator's admission webhooks (failurePolicy Fail) are served with a cert-manager Certificate in the default posture — cert-manager must be running before the operator installs.
- `KubernetesOtelCollector`
- `KubernetesKyverno` -- 4080–4099: Kubernetes security, policy, and identity
- `KubernetesGatekeeper`
- `KubernetesKeycloak` -- Keycloak declarations compose the official Keycloak Operator (which reconciles the Keycloak CR this kind renders) and, on the recommended postgres vendor, a KubernetesPostgres database — both must resolve before the CR can converge.
- `KubernetesOpenBao`
- `KubernetesOpenFga` -- OpenFGA requires a datastore; the recommended arm composes a KubernetesPostgres database (the sandbox memory arm needs nothing, but the registry declares the shape real deployments require).
- `KubernetesKeycloakOperator`
- `KubernetesCloudNativePgOperator` -- 4100–4129: Kubernetes data platforms
- `KubernetesPostgres`
- `KubernetesValkey`
- `KubernetesPerconaMysqlOperator`
- `KubernetesMysql`
- `KubernetesPerconaMongoOperator`
- `KubernetesMongodb`
- `KubernetesStrimziKafkaOperator`
- `KubernetesKafka` -- container_kind: a Strimzi Kafka cluster is a place in the provider's own model — KafkaTopic and KafkaUser declarations BELONG to one cluster (the strimzi.io/cluster label) and are drawn inside its box. Clients that merely talk to the cluster (Connect, MirrorMaker2, UI, Karapace) carry containment_exempt on their bootstrap/trust references.
- `KubernetesKafkaTopic`
- `KubernetesKafkaUser`
- `KubernetesKafkaConnect` -- container_kind: a Connect cluster hosts the connectors deployed INTO it (KafkaConnector's strimzi.io/cluster label names its Connect cluster) — the same room shape as KubernetesKafka above.
- `KubernetesKafkaConnector`
- `KubernetesKafkaMirrorMaker2`
- `KubernetesKarapace`
- `KubernetesKafkaUi`
- `KubernetesOpenSearchOperator`
- `KubernetesOpenSearch`
- `KubernetesAltinityOperator`
- `KubernetesClickHouse`
- `KubernetesSolrOperator`
- `KubernetesSolr`
- `KubernetesNeo4j`
- `KubernetesSeaweedFs`
- `KubernetesQdrant`
- `KubernetesRabbitMqOperator` -- The RabbitMQ Cluster Operator's release manifest ships admission webhooks whose serving certificate is a cert-manager Certificate — cert-manager must be running before the operator installs.
- `KubernetesRabbitMq`
- `KubernetesAirflow` -- 4130–4149: Kubernetes analytics and ML KubernetesPostgres is a prerequisite because Airflow's metadata database composes a KubernetesPostgres by default (the spec's FK defaults resolve onto its outputs) and the migration Job needs the database reachable before the server components start.
- `KubernetesSparkOperator`
- `KubernetesKubeRayOperator`
- `KubernetesRayCluster` -- KubernetesKubeRayOperator is a prerequisite because this kind declares the RayCluster custom resource that only the operator's CRDs admit and only the operator reconciles into head and worker pods.
- `KubernetesFlinkOperator` -- KubernetesCertManager is a prerequisite because the Flink operator's chart, with its default-on admission webhook, renders cert-manager Issuer/Certificate resources and trusts the API server through cert-manager's CA injection — there is no self-signed fallback at the pinned chart, and the webhooks are fail-closed.
- `KubernetesFlinkDeployment` -- KubernetesFlinkOperator is a prerequisite because this kind declares the FlinkDeployment custom resource that only the operator's CRDs admit and only the operator reconciles into a running Flink cluster.
- `KubernetesJupyterHub` -- KubernetesPostgres is a prerequisite because JupyterHub's hub database composes a KubernetesPostgres in its external-database arm (the spec's FK defaults resolve onto its outputs) and the hub pod mounts that database's credential Secret before it can start.
- `KubernetesMlflow` -- KubernetesPostgres is a prerequisite because MLflow's backend store composes a KubernetesPostgres in its production arm (FK defaults onto its outputs; the module composes the connection URI from its credential Secret), and KubernetesSeaweedFs because the artifact store's S3-compatible arm FK-defaults onto the SeaweedFS endpoint and credential Secret.
- `KubernetesTrino` -- KubernetesPostgres is a prerequisite because Trino's postgres catalogs compose a KubernetesPostgres (the catalog host and credential FK-default onto its outputs), and the pods read that database's credential Secret to resolve catalog passwords from environment.
- `KubernetesSuperset` -- KubernetesPostgres is a prerequisite because Superset's REQUIRED metadata database composes a KubernetesPostgres (FK defaults onto its outputs; the module composes the environment Secret from its credential Secret), and KubernetesValkey because the cache/broker arm FK-defaults onto a KubernetesValkey's service and password Secret.
- `KubernetesArgocd` -- 4150–4169: Kubernetes GitOps and CI/CD
- `KubernetesArgoWorkflows`
- `KubernetesTektonOperator`
- `KubernetesTekton` -- KubernetesTektonOperator is a prerequisite because this kind declares the TektonConfig custom resource that only the operator's CRDs admit and only the operator reconciles into running components.
- `KubernetesGhaRunnerScaleSetController`
- `KubernetesGhaRunnerScaleSet` -- KubernetesGhaRunnerScaleSetController is a prerequisite because this kind renders an AutoscalingRunnerSet custom resource that only the controller's CRDs admit and only the controller reconciles into listener and runner pods.
- `KubernetesHarbor`
- `KubernetesJenkins`
- `KubernetesTemporal` -- 4170–4189: Kubernetes app platforms KubernetesPostgres is a prerequisite because the recommended (and E2E-proven) database composition backs Temporal's default and visibility stores with a CloudNativePG cluster.
- `KubernetesNats`
- `KubernetesLocust`
- `DigitalOceanAppPlatformService` -- 5000–5999: DigitalOcean resources
- `DigitalOceanBucket`
- `DigitalOceanContainerRegistry`
- `DigitalOceanDatabaseCluster`
- `DigitalOceanDnsZone`
- `DigitalOceanDroplet`
- `DigitalOceanFirewall`
- `DigitalOceanFunction`
- `DigitalOceanKubernetesCluster`
- `DigitalOceanKubernetesNodePool`
- `DigitalOceanLoadBalancer`
- `DigitalOceanVolume`
- `DigitalOceanVpc`
- `DigitalOceanCertificate`
- `DigitalOceanDnsRecord`
- `CivoBucket` -- 6000–6999: Civo resources
- `CivoCertificate`
- `CivoComputeInstance`
- `CivoDatabase`
- `CivoDnsZone`
- `CivoFirewall`
- `CivoIpAddress`
- `CivoKubernetesCluster`
- `CivoKubernetesNodePool`
- `CivoVolume`
- `CivoVpc`
- `CivoDnsRecord`
- `CloudflareDnsZone` -- 7000–7999: Cloudflare resources
- `CloudflareKvNamespace`
- `CloudflareR2Bucket`
- `CloudflareWorker`
- `CloudflareLoadBalancer`
- `CloudflareD1Database`
- `CloudflareZeroTrustAccessApplication`
- `CloudflareDnsRecord`
- `CloudflareRuleset`
- `CloudflareWorkersKvPair`
- `CloudflareHyperdriveConfig`
- `CloudflareLoadBalancerPool`
- `CloudflareLoadBalancerMonitor`
- `CloudflareZeroTrustAccessPolicy`
- `CloudflareZeroTrustAccessGroup`
- `CloudflareQueue`
- `CloudflarePagesProject`
- `CloudflareZeroTrustTunnel`
- `CloudflareZeroTrustTunnelVirtualNetwork`
- `CloudflareZeroTrustTunnelRoute`
- `CloudflareList`
- `CloudflareListItem`
- `CloudflareTurnstileWidget`
- `CloudflareEmailRoutingZone`
- `CloudflareEmailRoutingRule`
- `CloudflareEmailRoutingAddress`
- `CloudflareOriginCaCertificate`
- `CloudflareCertificatePack`
- `CloudflareCustomHostname`
- `CloudflareCustomHostnameFallbackOrigin`
- `Auth0Connection` -- 8000–8999: Auth0 resources
- `Auth0Client`
- `Auth0EventStream`
- `Auth0ResourceServer`
- `Auth0Action`
- `Auth0Role`
- `OpenFgaStore` -- 9000–9999: OpenFGA resources Note: OpenFGA is Terraform-only - there is no Pulumi provider available. Pulumi modules for OpenFGA resources are pass-through placeholders.
- `OpenFgaAuthorizationModel`
- `OpenFgaRelationshipTuple`
- `OpenStackKeypair` -- 10000–10999: OpenStack resources
- `OpenStackNetwork`
- `OpenStackSubnet`
- `OpenStackRouter`
- `OpenStackRouterInterface`
- `OpenStackSecurityGroup`
- `OpenStackFloatingIp`
- `OpenStackNetworkPort`
- `OpenStackSecurityGroupRule`
- `OpenStackFloatingIpAssociate`
- `OpenStackInstance`
- `OpenStackServerGroup`
- `OpenStackVolume`
- `OpenStackVolumeAttach`
- `OpenStackProject`
- `OpenStackApplicationCredential`
- `OpenStackImage`
- `OpenStackRoleAssignment`
- `OpenStackLoadBalancer`
- `OpenStackLoadBalancerListener`
- `OpenStackLoadBalancerPool`
- `OpenStackLoadBalancerMember`
- `OpenStackLoadBalancerMonitor`
- `OpenStackDnsZone`
- `OpenStackDnsRecord`
- `ScalewayVpc`
- `ScalewayPrivateNetwork`
- `ScalewayPublicGateway`
- `ScalewayLoadBalancer`
- `ScalewayInstanceSecurityGroup`
- `ScalewayInstance`
- `ScalewayKapsuleCluster`
- `ScalewayKapsulePool`
- `ScalewayRdbInstance`
- `ScalewayRedisCluster`
- `ScalewayMongodbInstance`
- `ScalewayObjectBucket`
- `ScalewayBlockVolume`
- `ScalewayContainerRegistry`
- `ScalewayDnsZone`
- `ScalewayDnsRecord`
- `ScalewayServerlessFunction`
- `ScalewayServerlessContainer`
- `AliCloudLogProject`
- `AliCloudRamRole`
- `AliCloudRamPolicy`
- `AliCloudVpc`
- `AliCloudVswitch`
- `AliCloudSecurityGroup`
- `AliCloudEipAddress`
- `AliCloudNatGateway`
- `AliCloudApplicationLoadBalancer`
- `AliCloudNetworkLoadBalancer`
- `AliCloudVpnGateway`
- `AliCloudDnsZone`
- `AliCloudDnsRecord`
- `AliCloudPrivateDnsZone`
- `AliCloudStorageBucket`
- `AliCloudNasFileSystem`
- `AliCloudKmsKey`
- `AliCloudRdsInstance`
- `AliCloudPolardbCluster`
- `AliCloudRedisInstance`
- `AliCloudMongodbInstance`
- `AliCloudEcsInstance`
- `AliCloudContainerRegistry`
- `AliCloudKubernetesCluster`
- `AliCloudKubernetesNodePool`
- `AliCloudCdnDomain`
- `AliCloudFunction`
- `AliCloudSaeApplication`
- `AliCloudRocketmqInstance`
- `AliCloudCenInstance`
- `OciVcn`
- `OciSubnet`
- `OciSecurityGroup`
- `OciCompartment`
- `OciIdentityPolicy`
- `OciDynamicGroup`
- `OciComputeInstance`
- `OciContainerEngineCluster`
- `OciContainerEngineNodePool`
- `OciContainerInstance`
- `OciApplicationLoadBalancer`
- `OciNetworkLoadBalancer`
- `OciDynamicRoutingGateway`
- `OciPublicIp`
- `OciAutonomousDatabase`
- `OciDbSystem`
- `OciMysqlDbSystem`
- `OciPostgresqlDbSystem`
- `OciRedisCluster`
- `OciNosqlTable`
- `OciObjectStorageBucket`
- `OciFileSystem`
- `OciBlockVolume`
- `OciKmsVault`
- `OciKmsKey`
- `OciVaultSecret`
- `OciBastion`
- `OciFunctionsApplication`
- `OciApiGateway`
- `OciStreamPool`
- `OciQueue`
- `OciAlarm`
- `OciLogGroup`
- `OciDnsZone`
- `OciDnsRecord`
- `OciNetworkFirewall`
- `OciDevopsProject`
- `HetznerCloudSshKey`
- `HetznerCloudPlacementGroup`
- `HetznerCloudFirewall`
- `HetznerCloudNetwork`
- `HetznerCloudPrimaryIp`
- `HetznerCloudFloatingIp`
- `HetznerCloudServer`
- `HetznerCloudVolume`
- `HetznerCloudSnapshot`
- `HetznerCloudCertificate`
- `HetznerCloudLoadBalancer`
- `HetznerCloudDnsZone`

### spec.container.app.env.secrets[].valueFrom.env

`string`

### spec.container.app.env.secrets[].valueFrom.name

`string` · required

- rule: {"required":true}

### spec.container.app.env.secrets[].valueFrom.fieldPath

`string`

### spec.container.app.env.envFrom

`[]EnvFromSource`

Bulk import of environment variables from ConfigMaps or Secrets.

### spec.container.app.env.envFrom[].prefix

`string`

Optional prefix prepended to each imported key name.
For example, prefix "APP_" with key "PORT" produces env var "APP_PORT".

### spec.container.app.env.envFrom[].configMapRef

`ConfigMapRef`

Import all keys from a ConfigMap.

### spec.container.app.env.envFrom[].configMapRef.name

`string` · required

Name of the ConfigMap.

- rule: {"required":true}

### spec.container.app.env.envFrom[].configMapRef.optional

`bool`

If true, the ConfigMap is allowed to not exist without blocking pod startup.

### spec.container.app.env.envFrom[].secretRef

`SecretRef`

Import all keys from a Secret.

### spec.container.app.env.envFrom[].secretRef.name

`string` · required

Name of the Secret.

- rule: {"required":true}

### spec.container.app.env.envFrom[].secretRef.optional

`bool`

If true, the Secret is allowed to not exist without blocking pod startup.

### spec.container.app.resources

`ContainerResources`

CPU and memory requests and limits. Requests drive scheduling and are what the
pod is guaranteed; limits are the ceiling enforced at runtime (CPU is throttled,
memory overage is OOM-killed). Omitting limits entirely leaves the container
unbounded — acceptable for batch work on dedicated nodes, risky on shared ones.

### spec.container.app.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.container.app.resources.limits.cpu

`string`

### spec.container.app.resources.limits.memory

`string`

### spec.container.app.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.container.app.resources.requests.cpu

`string`

### spec.container.app.resources.requests.memory

`string`

### spec.container.app.livenessProbe

`Probe`

Liveness probe: restarts the container when it fails. Detects deadlocked or
wedged processes. Keep it strictly about "is the process alive" — checking
downstream dependencies here turns a dependency blip into a restart storm.

### spec.container.app.livenessProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.container.app.livenessProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.container.app.livenessProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.container.app.livenessProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.container.app.livenessProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.container.app.livenessProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.container.app.livenessProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.app.livenessProbe.httpGet.portNumber

`int32`

### spec.container.app.livenessProbe.httpGet.portName

`string`

### spec.container.app.livenessProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.app.livenessProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.app.livenessProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.app.livenessProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.app.livenessProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.app.livenessProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.container.app.livenessProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.container.app.livenessProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.container.app.livenessProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.container.app.livenessProbe.tcpSocket.portNumber

`int32`

### spec.container.app.livenessProbe.tcpSocket.portName

`string`

### spec.container.app.livenessProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.app.livenessProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.container.app.livenessProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.app.readinessProbe

`Probe`

Readiness probe: removes the pod from Service endpoints while it fails. This is
the probe that makes rolling updates zero-downtime — traffic only reaches pods
that report ready.

### spec.container.app.readinessProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.container.app.readinessProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.container.app.readinessProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.container.app.readinessProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.container.app.readinessProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.container.app.readinessProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.container.app.readinessProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.app.readinessProbe.httpGet.portNumber

`int32`

### spec.container.app.readinessProbe.httpGet.portName

`string`

### spec.container.app.readinessProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.app.readinessProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.app.readinessProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.app.readinessProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.app.readinessProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.app.readinessProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.container.app.readinessProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.container.app.readinessProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.container.app.readinessProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.container.app.readinessProbe.tcpSocket.portNumber

`int32`

### spec.container.app.readinessProbe.tcpSocket.portName

`string`

### spec.container.app.readinessProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.app.readinessProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.container.app.readinessProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.app.startupProbe

`Probe`

Startup probe: holds off liveness and readiness checking until the app has
started, so slow-booting applications are not killed mid-initialization. Size
`failure_threshold × period_seconds` to the worst-case startup time.

### spec.container.app.startupProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.container.app.startupProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.container.app.startupProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.container.app.startupProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.container.app.startupProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.container.app.startupProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.container.app.startupProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.app.startupProbe.httpGet.portNumber

`int32`

### spec.container.app.startupProbe.httpGet.portName

`string`

### spec.container.app.startupProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.app.startupProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.app.startupProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.app.startupProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.app.startupProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.app.startupProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.container.app.startupProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.container.app.startupProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.container.app.startupProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.container.app.startupProbe.tcpSocket.portNumber

`int32`

### spec.container.app.startupProbe.tcpSocket.portName

`string`

### spec.container.app.startupProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.app.startupProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.container.app.startupProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.app.volumeMounts

`[]VolumeMount`

Volume mounts for this container. Each entry both declares the mount path and
carries its volume source (ConfigMap, Secret, HostPath, EmptyDir, or PVC); the
module derives the pod-level volume list from the union of all containers'
mounts, de-duplicating by name — so two containers sharing an EmptyDir simply
declare the same mount name and source.

### spec.container.app.volumeMounts[].name

`string` · required

Name of the volume mount. Must be unique within the container.
Used to correlate with the volume definition.

- rule: {"required":true}

### spec.container.app.volumeMounts[].mountPath

`string` · required

Path within the container at which the volume should be mounted.
Must be an absolute path.

- rule: {"required":true}

### spec.container.app.volumeMounts[].readOnly

`bool`

Whether the volume should be mounted read-only.
Default is false.

### spec.container.app.volumeMounts[].subPath

`string`

Path within the volume from which the container's volume should be mounted.
Defaults to "" (volume's root).
Useful for mounting a subdirectory of a volume.

### spec.container.app.volumeMounts[].configMap

`ConfigMapVolumeSource`

ConfigMap volume source.
Use this to mount a ConfigMap as a file or directory.

### spec.container.app.volumeMounts[].configMap.name

`string` · required

Name of the ConfigMap to mount.
Can reference a ConfigMap defined in spec.config_maps or an existing one in the namespace.

- rule: {"required":true}

### spec.container.app.volumeMounts[].configMap.key

`string`

Specific key from the ConfigMap to mount as a single file.
If not specified, all keys are mounted as files in the directory.

### spec.container.app.volumeMounts[].configMap.path

`string`

If key is specified, this is the filename to use for the mounted file.
Defaults to the key name if not specified.
Example: key="config" path="app.yaml" mounts the "config" key as "app.yaml"

### spec.container.app.volumeMounts[].configMap.defaultMode

`int32`

Mode bits to use on created files. Must be a value between 0 and 0777.
Defaults to 0644.
Use 0755 (493 in decimal) for executable scripts.

### spec.container.app.volumeMounts[].secret

`SecretVolumeSource`

Secret volume source.
Use this to mount a Secret as a file or directory.

### spec.container.app.volumeMounts[].secret.name

`string` · required

Name of the Secret to mount.

- rule: {"required":true}

### spec.container.app.volumeMounts[].secret.key

`string`

Specific key from the Secret to mount as a single file.
If not specified, all keys are mounted as files in the directory.

### spec.container.app.volumeMounts[].secret.path

`string`

If key is specified, this is the filename to use for the mounted file.
Defaults to the key name if not specified.

### spec.container.app.volumeMounts[].secret.defaultMode

`int32`

Mode bits to use on created files. Must be a value between 0 and 0777.
Defaults to 0644.

### spec.container.app.volumeMounts[].hostPath

`HostPathVolumeSource`

HostPath volume source.
Use this to mount a file or directory from the host node's filesystem.
Common for DaemonSets that need access to node-level resources.

### spec.container.app.volumeMounts[].hostPath.path

`string` · required

Path on the host to mount.

- rule: {"required":true}

### spec.container.app.volumeMounts[].hostPath.type

`string`

Type of the host path.
Valid values:
  "" - Empty string (default) means no check is performed before mounting
  "DirectoryOrCreate" - Create directory if it doesn't exist
  "Directory" - Directory must exist
  "FileOrCreate" - Create file if it doesn't exist
  "File" - File must exist
  "Socket" - UNIX socket must exist
  "CharDevice" - Character device must exist
  "BlockDevice" - Block device must exist

- rule: Type must be one of: "", "DirectoryOrCreate", "Directory", "FileOrCreate", "File", "Socket", "CharDevice", "BlockDevice"

### spec.container.app.volumeMounts[].emptyDir

`EmptyDirVolumeSource`

EmptyDir volume source.
Use this for temporary storage that is erased when the pod is removed.
Useful for scratch space, caching, or sharing data between containers.

### spec.container.app.volumeMounts[].emptyDir.medium

`string`

Medium for the empty directory.
"" (default) uses the node's default medium (typically disk).
"Memory" uses a tmpfs (RAM-backed filesystem).

Memory-backed volumes are faster but:
- Count against container memory limits
- Are lost on node restart
- Should have sizeLimit set to prevent OOM

- rule: Medium must be either "" or "Memory"

### spec.container.app.volumeMounts[].emptyDir.sizeLimit

`string`

Size limit for the empty directory.
Format: Kubernetes quantity (e.g., "1Gi", "500Mi").
Only strictly enforced when medium is "Memory".
For disk-backed volumes, this is a best-effort limit.

### spec.container.app.volumeMounts[].pvc

`PvcVolumeSource`

PersistentVolumeClaim volume source.
Use this to mount an existing PVC.
For StatefulSets, this can reference a volumeClaimTemplate.

### spec.container.app.volumeMounts[].pvc.claimName

`string` · required

Name of the PersistentVolumeClaim to mount.
For StatefulSets, this can be the name of a volumeClaimTemplate.

- rule: {"required":true}

### spec.container.app.volumeMounts[].pvc.readOnly

`bool`

Whether the PVC should be mounted read-only.
Default is false.

### spec.container.app.lifecycle

`WorkloadContainerLifecycle`

Lifecycle hooks. `post_start` runs immediately after the container starts (the
container is not Running until it completes); `pre_stop` runs before termination
and is the standard lever for connection draining — e.g. a short sleep that keeps
the endpoint serving while load balancers converge on the terminating state.

### spec.container.app.lifecycle.postStart

`WorkloadLifecycleHandler`

Runs immediately after the container is created. The container does not reach
Running until the hook completes; a failing post_start kills the container per
its restart policy.

### spec.container.app.lifecycle.postStart.exec

`ExecAction`

Execute a command inside the container; non-zero exit fails the hook.

### spec.container.app.lifecycle.postStart.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.app.lifecycle.postStart.httpGet

`HTTPGetAction`

Perform an HTTP GET against the container.

### spec.container.app.lifecycle.postStart.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.app.lifecycle.postStart.httpGet.portNumber

`int32`

### spec.container.app.lifecycle.postStart.httpGet.portName

`string`

### spec.container.app.lifecycle.postStart.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.app.lifecycle.postStart.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.app.lifecycle.postStart.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.app.lifecycle.postStart.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.app.lifecycle.postStart.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.app.lifecycle.postStart.tcpSocket

`TCPSocketAction`

Open a TCP connection against the container.

### spec.container.app.lifecycle.postStart.tcpSocket.portNumber

`int32`

### spec.container.app.lifecycle.postStart.tcpSocket.portName

`string`

### spec.container.app.lifecycle.postStart.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.app.lifecycle.postStart.sleep

`SleepAction`

Pause for a fixed number of seconds — the idiomatic pre_stop drain hook,
with no shell or sleep binary required in the image.

### spec.container.app.lifecycle.postStart.sleep.seconds

`int64`

Seconds to sleep.

- rule: {"int64":{"gte":"0"}}

### spec.container.app.lifecycle.preStop

`WorkloadLifecycleHandler`

Runs before the container is terminated by the kubelet (pod deletion, rolling
update, eviction). The termination grace period starts BEFORE the hook runs, so
keep `pod.termination_grace_period_seconds` larger than the hook's worst-case
duration. The classic zero-downtime pattern is a short sleep here so the pod
keeps serving while endpoint removal propagates.

### spec.container.app.lifecycle.preStop.exec

`ExecAction`

Execute a command inside the container; non-zero exit fails the hook.

### spec.container.app.lifecycle.preStop.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.app.lifecycle.preStop.httpGet

`HTTPGetAction`

Perform an HTTP GET against the container.

### spec.container.app.lifecycle.preStop.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.app.lifecycle.preStop.httpGet.portNumber

`int32`

### spec.container.app.lifecycle.preStop.httpGet.portName

`string`

### spec.container.app.lifecycle.preStop.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.app.lifecycle.preStop.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.app.lifecycle.preStop.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.app.lifecycle.preStop.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.app.lifecycle.preStop.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.app.lifecycle.preStop.tcpSocket

`TCPSocketAction`

Open a TCP connection against the container.

### spec.container.app.lifecycle.preStop.tcpSocket.portNumber

`int32`

### spec.container.app.lifecycle.preStop.tcpSocket.portName

`string`

### spec.container.app.lifecycle.preStop.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.app.lifecycle.preStop.sleep

`SleepAction`

Pause for a fixed number of seconds — the idiomatic pre_stop drain hook,
with no shell or sleep binary required in the image.

### spec.container.app.lifecycle.preStop.sleep.seconds

`int64`

Seconds to sleep.

- rule: {"int64":{"gte":"0"}}

### spec.container.app.securityContext

`WorkloadContainerSecurityContext`

Container-level security hardening. Settings here override the pod-level
security context for this container only.

### spec.container.app.securityContext.privileged

`bool`

Runs the container with full host access — equivalent to root on the node.
Required by some node-level agents (device managers, network plugins). Never
combine with untrusted images.

### spec.container.app.securityContext.runAsUser

`int64` · optional (explicit presence)

UID the container process runs as. Overrides the image's USER directive.

### spec.container.app.securityContext.runAsGroup

`int64` · optional (explicit presence)

Primary GID the container process runs as.

### spec.container.app.securityContext.runAsNonRoot

`bool` · optional (explicit presence)

Refuses to start the container if its effective user is root. The standard
baseline hardening — it catches images that silently default to UID 0.

### spec.container.app.securityContext.readOnlyRootFilesystem

`bool` · optional (explicit presence)

Mounts the container's root filesystem read-only. Pair with EmptyDir mounts for
paths the app must write (e.g. /tmp).

### spec.container.app.securityContext.allowPrivilegeEscalation

`bool` · optional (explicit presence)

Whether the process can gain more privileges than its parent (setuid binaries,
file capabilities). The restricted Pod Security Standard requires this to be
false. Always true when `privileged` is set, so leave it unset in that case.

### spec.container.app.securityContext.capabilities

`WorkloadCapabilities`

Linux capabilities to add or drop. The restricted profile drops ALL and adds
back only NET_BIND_SERVICE when needed. Capability names are uppercase without
the CAP_ prefix (e.g. "NET_ADMIN", "SYS_TIME").

### spec.container.app.securityContext.capabilities.add

`[]string`

Capabilities to add (e.g. "NET_BIND_SERVICE").

### spec.container.app.securityContext.capabilities.drop

`[]string`

Capabilities to drop. Use ["ALL"] as the hardened baseline.

### spec.container.app.securityContext.seccompProfile

`WorkloadSeccompProfile`

Seccomp syscall filter for the container. "RuntimeDefault" is the hardened
baseline; "Localhost" selects a node-local profile file via `localhost_profile`.

- rule: localhost_profile is required when type is "Localhost" and must be empty otherwise

### spec.container.app.securityContext.seccompProfile.type

`string` · required

Profile type: "RuntimeDefault" (the container runtime's default filter — the
recommended baseline), "Unconfined" (no filtering), or "Localhost" (a profile
file installed on the node, named via localhost_profile).

- rule: Seccomp profile type must be one of "RuntimeDefault", "Unconfined", or "Localhost"
- rule: {"required":true}

### spec.container.app.securityContext.seccompProfile.localhostProfile

`string`

Path of the profile file relative to the node's seccomp profile root. Required
when (and only meaningful when) type is "Localhost".

### spec.container.sidecars

`[]WorkloadContainer`

Sidecar containers running alongside the app in every replica — metrics
exporters, backup agents, connection poolers. Sidecars are full containers:
probes, mounts, security context, and lifecycle hooks all apply, and they may
mount the same volume claim templates as the app. Each sidecar must be named.

- rule: Every sidecar container must have a name

### spec.container.sidecars[].name

`string`

The container's name, unique within the pod. Required for sidecars and init
containers (Kubernetes rejects unnamed containers); for the main app container the
module defaults it when omitted, so minimal manifests stay minimal. Must be a valid
DNS label: lowercase alphanumeric and hyphens, starting and ending alphanumeric.

- rule: Container name must be a lowercase DNS label (alphanumeric and hyphens, starting and ending with an alphanumeric character)

### spec.container.sidecars[].image

`ContainerImage` · required

The container image, split into repository and tag so deployment pipelines can
inject a freshly built tag without rewriting the whole reference. The optional
`pull_secret_name` names an existing docker-registry secret; prefer attaching pull
secrets on the ServiceAccount (or `pod.image_pull_secrets`) so they apply pod-wide.

- rule: Image repo is required — the repository half of the image reference (e.g. "nginx" or "ghcr.io/acme/api")
- rule: Image tag is required — pin a version (e.g. "1.27.1"); avoid "latest" for anything you intend to roll back
- rule: {"required":true}

### spec.container.sidecars[].image.repo

`string`

The repository of the image (e.g., "gcr.io/project/image").

### spec.container.sidecars[].image.tag

`string`

The tag of the image (e.g., "latest" or "1.0.0").

### spec.container.sidecars[].image.pullSecretName

`string`

The name of the image pull secret for private image repositories.

### spec.container.sidecars[].imagePullPolicy

`string`

When the kubelet pulls the image. "IfNotPresent" (the Kubernetes default for tagged
images) reuses a cached copy; "Always" re-resolves the tag on every pod start —
required when a mutable tag like a branch name is reused across builds; "Never"
only uses pre-loaded images (air-gapped nodes, kind-loaded test images).

- rule: Image pull policy must be one of "Always", "IfNotPresent", or "Never"

### spec.container.sidecars[].command

`[]string`

Entrypoint override (Kubernetes `command`, Docker ENTRYPOINT). The image's
entrypoint runs when omitted. Not executed in a shell — provide argv elements,
e.g. ["/bin/sh", "-c", "exec my-server"].

### spec.container.sidecars[].args

`[]string`

Arguments to the entrypoint (Kubernetes `args`, Docker CMD). The image's CMD is
used when omitted. Variable references like $(VAR_NAME) are expanded from the
container's environment by the kubelet.

### spec.container.sidecars[].workingDir

`string`

Working directory for the entrypoint. Defaults to the image's configured WORKDIR.

### spec.container.sidecars[].ports

`[]WorkloadContainerPort`

Network ports this container exposes. Purely informational to Kubernetes for plain
pod-to-pod traffic, but load-bearing here: named ports are referenced by probes,
and `service_port` drives the Service wiring on kinds that create one
(Deployment, StatefulSet).

### spec.container.sidecars[].ports[].name

`string` · required

Port name, e.g. "http", "grpc", "metrics". Must be a lowercase DNS label that
starts and ends alphanumeric. Named ports are referenced by probes and become the
Service port names on service-fronted kinds.

- rule: Port name must contain only lowercase alphanumeric characters and hyphens, and start and end with an alphanumeric character (e.g. "http", "grpc-web")
- rule: {"required":true}

### spec.container.sidecars[].ports[].containerPort

`int32` · required

The port number the container listens on (1–65535).

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.container.sidecars[].ports[].networkProtocol

`string`

L4 protocol of the port. Defaults to "TCP" when omitted — the overwhelmingly
common case, so minimal manifests need not repeat it.

- rule: The network protocol must be one of "TCP", "UDP", or "SCTP"

### spec.container.sidecars[].ports[].appProtocol

`string`

Application protocol hint (e.g. "http", "grpc", "https"). Propagated to the
Service port's appProtocol on service-fronted kinds, where meshes and L7 load
balancers use it to pick the right protocol handling.

### spec.container.sidecars[].ports[].servicePort

`int32`

The port the workload's Kubernetes Service exposes for this container port.
Only meaningful on kinds that create a Service (Deployment, StatefulSet); other
kinds ignore it. E.g. containerPort 8080 with servicePort 80 serves the app on
the conventional port while the process binds an unprivileged one. External
exposure is composed separately with first-class ingress kinds referencing the
workload's exported Service handle — workloads never create ingress themselves.

- rule: Service port must be between 1 and 65535

### spec.container.sidecars[].ports[].hostPort

`int32`

Exposes the container port directly on the node's IP (hostPort). Chiefly a
DaemonSet pattern (node-level agents that must be reachable on every node);
on other kinds it constrains scheduling to one pod per node per port — prefer
a Service unless node-local reachability is the point.

- rule: Host port must be between 1 and 65535

### spec.container.sidecars[].env

`ContainerEnv`

Environment configuration: plain variables (with Kubernetes-native value sources
and Planton cross-resource references), secret variables (materialized into a
managed Kubernetes Secret), and bulk envFrom imports.

### spec.container.sidecars[].env.variables

`[]EnvVar`

Individual environment variables (non-sensitive).

### spec.container.sidecars[].env.variables[].name

`string` · required

The environment variable name.
Must be a valid C_IDENTIFIER: starts with a letter or underscore,
followed by letters, digits, or underscores.

- rule: Must be a valid C_IDENTIFIER (start with letter/underscore, contain only letters, digits, underscores)
- rule: {"required":true}

### spec.container.sidecars[].env.variables[].value

`string`

Direct literal value.

### spec.container.sidecars[].env.variables[].valueFrom

`ValueFromRef`

Reference to another Planton resource's field.
The orchestrator resolves this and populates the value before invoking IaC modules.

### spec.container.sidecars[].env.variables[].valueFrom.kind

`enum`

Allowed values (use exactly as shown):

- `unspecified` -- 0: Default/unspecified
- `TestCloudResourceGeneric` -- 1–49: Test/dev/custom
- `TestCloudResourceKubernetes`
- `ConfluentKafka` -- 50–199: saas platform resources
- `AtlasMongodb`
- `SnowflakeDatabase`
- `AwsAlb` -- 1000–1999: AWS resources AwsSubnet is a prerequisite because an ALB requires at least two subnets in different availability zones -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsCertManagerCert`
- `AwsCloudFront`
- `AwsDynamodb`
- `AwsEcrRepo`
- `AwsEcsCluster`
- `AwsEcsService` -- AwsEcsCluster, AwsEcsTaskDefinition, and AwsSubnet are prerequisites because a service schedules a referenced task-definition revision into a referenced live cluster and places task network interfaces into referenced subnets -- all three references must resolve first.
- `AwsEksCluster` -- AwsSubnet and AwsIamRole are prerequisites because the control plane attaches its network interfaces into referenced subnets and assumes a referenced cluster role that must already carry AmazonEKSClusterPolicy.
- `AwsIamRole`
- `AwsLambda`
- `AwsRdsCluster`
- `AwsRdsInstance`
- `AwsRoute53Zone`
- `AwsS3Bucket`
- `AwsLbTargetGroup` -- AwsVpc is a prerequisite because a target group's health checks and target registrations live inside one VPC -- the spec's vpc_id reference must resolve before the group can be created.
- `AwsSecurityGroup` -- AwsVpc is a prerequisite because every security group is created in a VPC; the E2E install profile resolves vpc_id against the VPC prerequisite.
- `AwsVpc`
- `AwsEksNodeGroup` -- AwsEksCluster is a prerequisite because nodes register with a live control plane; AwsIamRole and AwsSubnet back the node role and worker subnet references.
- `AwsIamUser`
- `AwsKmsKey`
- `AwsEc2Instance`
- `AwsClientVpn` -- Every Client VPN endpoint requires an ACM server certificate at create time; the imported self-signed fixture satisfies it. Subnets/VPC are optional composition (a zero-association endpoint is valid) -- composed scenarios declare them via the e2e-prerequisites annotation.
- `AwsDocumentDb`
- `AwsRoute53DnsRecord` -- AwsRoute53Zone is a prerequisite because every record lives inside a hosted zone -- the spec's zone_id reference must resolve before the record can be created.
- `AwsS3ObjectSet` -- AwsS3Bucket is a prerequisite because the object set's bucket reference is required -- objects cannot exist without the bucket that holds them.
- `AwsSqsQueue`
- `AwsSnsTopic`
- `AwsEventBridgeBus`
- `AwsEventBridgeRule`
- `AwsIamOidcProvider`
- `AwsIamPolicy`
- `AwsIamInstanceProfile` -- AwsIamRole is a prerequisite because an instance profile is a wrapper that must contain a role to be useful -- the profile's spec requires a role reference, so the role must be deployed first.
- `AwsLbListener` -- AwsAlb and AwsLbTargetGroup are prerequisites because a listener is an attachment point on a load balancer and its default action almost always forwards to a target group -- both references must resolve before the listener can be created.
- `AwsLbListenerRule` -- AwsLbListener is a prerequisite because a rule only exists as an attachment on a listener -- the listener_arn reference must resolve before the rule can be created.
- `AwsLaunchTemplate`
- `AwsAutoScalingGroup` -- AwsSubnet and AwsLaunchTemplate are prerequisites because a group cannot exist without subnets to place capacity in and a launch template to launch from -- the spec's subnets and launch_template references must resolve before the group can be created.
- `AwsEksAddon` -- AwsEksCluster is a prerequisite because an add-on installs onto a live control plane -- the spec's cluster_name reference must resolve before the add-on can be created.
- `AwsEksFargateProfile` -- AwsEksCluster, AwsIamRole, and AwsSubnet are prerequisites because a Fargate profile attaches to a live control plane, runs pods as a referenced pod-execution role, and launches them into referenced private subnets -- all three references must resolve first.
- `AwsEksAccessEntry` -- AwsEksCluster and AwsIamRole are prerequisites because an access entry grants a referenced IAM principal access to a live control plane -- both references must resolve before the entry can be created.
- `AwsEcsTaskDefinition` -- AwsIamRole is a prerequisite because the kind's default posture -- Fargate with the awslogs logging default -- is rejected by AWS at registration time without an execution role the agent can assume.
- `AwsHttpApiGateway`
- `AwsStepFunction` -- AwsIamRole is a prerequisite because a state machine cannot be created without an execution role it can assume -- the spec's role_arn reference must resolve before the CreateStateMachine call.
- `AwsHttpApiVpcLink` -- AwsSubnet is a prerequisite because a VPC link is a set of managed ENIs provisioned into referenced subnets -- the subnet references must resolve before the link can be created. Security groups are optional on the link, so they compose per-scenario rather than as a registry prerequisite.
- `AwsHttpApiDomain` -- AwsCertManagerCert is a prerequisite because a custom domain cannot be created without a TLS certificate in the same region covering the domain -- the spec's certificate_arn reference must resolve first.
- `AwsVpcEndpoint` -- AwsVpcEndpoint's composed E2E scenarios reference the AwsVpc prerequisite's outputs (vpc_id + default_route_table_id for gateway endpoints) and the AwsSubnet pair's subnet_id outputs (interface endpoints), so both are genuine deploy-order prerequisites.
- `AwsElasticacheUser`
- `AwsElasticacheUserGroup` -- AwsElasticacheUser is a genuine prerequisite: AWS refuses to create a user group that does not contain a user named "default", so a group's composed E2E scenario must resolve a deployed user's outputs.
- `AwsRedshiftServerlessNamespace`
- `AwsRedshiftServerlessWorkgroup` -- The namespace is a genuine prerequisite: a workgroup attaches to exactly one namespace by name at create time, so its composed E2E scenario must resolve a deployed namespace's outputs. AwsSubnet is a prerequisite because Redshift Serverless requires the workgroup's subnets to span three availability zones.
- `AwsRedisElasticache` -- AwsSubnet is a prerequisite because the module builds an ElastiCache subnet group from referenced subnets -- the spec's subnet references must resolve before the replication group can deploy.
- `AwsOpenSearchDomain`
- `AwsMemcachedElasticache`
- `AwsServerlessElasticache`
- `AwsNlb` -- AwsSubnet is a prerequisite because an NLB requires at least one subnet mapping -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsElasticIp`
- `AwsTransitGateway`
- `AwsGlobalAccelerator`
- `AwsSubnet`
- `AwsInternetGateway`
- `AwsNatGateway` -- AwsInternetGateway is a prerequisite because a public NAT gateway can only become available once the VPC it sits in has an internet gateway attached (AWS rejects the create otherwise) -- so the gateway must be deployed first. AwsVpc is a prerequisite because a REGIONAL NAT gateway (availability_mode = regional) references the VPC directly instead of a subnet.
- `AwsEgressOnlyInternetGateway`
- `AwsElasticFileSystem` -- AwsSubnet and AwsSecurityGroup are prerequisites because mount targets (required, min 1) place the file system's NFS endpoints into subnets and attach security groups -- both references must resolve before the CreateMountTarget calls.
- `AwsEfsAccessPoint` -- AwsElasticFileSystem is a prerequisite because an access point is created INTO a file system -- the spec's required file_system_id reference must resolve before the CreateAccessPoint call.
- `AwsFsxLustreFileSystem`
- `AwsFsxOpenzfsFileSystem`
- `AwsFsxWindowsFileSystem` -- Every Windows file system must join an Active Directory domain; the directory itself is external infrastructure (AWS Managed Microsoft AD or a self-managed domain), so only the network dependency is a declarable prerequisite.
- `AwsFsxOntapFileSystem`
- `AwsFsxOntapStorageVirtualMachine`
- `AwsFsxOntapVolume`
- `AwsFsxDataRepositoryAssociation`
- `AwsCognitoUserPool`
- `AwsCognitoIdentityProvider` -- AwsCognitoUserPool is a prerequisite because an identity provider is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateIdentityProvider call.
- `AwsCognitoUserPoolClient` -- AwsCognitoUserPool is a prerequisite because an app client is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateUserPoolClient call.
- `AwsCognitoResourceServer` -- AwsCognitoUserPool is a prerequisite because a resource server is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateResourceServer call.
- `AwsWafWebAcl`
- `AwsWafIpSet`
- `AwsWafRegexPatternSet`
- `AwsCloudwatchLogGroup`
- `AwsCloudwatchAlarm`
- `AwsCloudwatchCompositeAlarm`
- `AwsKinesisStream`
- `AwsKinesisFirehose` -- Every Firehose destination requires an S3 configuration (the primary target for extended_s3; the failed/all-document backup for the rest) and an IAM role Firehose assumes to write to it, so both are hard deploy prerequisites.
- `AwsKinesisStreamConsumer` -- A consumer registers against exactly one stream and cannot exist without it.
- `AwsAthenaWorkgroup`
- `AwsGlueCatalogDatabase`
- `AwsRedshiftCluster`
- `AwsSagemakerDomain` -- AI/ML A domain cannot exist without VPC subnets and a SageMaker execution role (default_user_settings.execution_role_arn is required), so both are hard deploy prerequisites.
- `AwsAppRunnerService` -- A service can run entirely on companion defaults, so the App Runner family's kinds are dependency-free leaves except the VPC connector (which cannot exist without subnets and security groups). A service's companion references (auto scaling / VPC connector / observability / WAF) are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsAppRunnerAutoScalingConfiguration`
- `AwsAppRunnerVpcConnector`
- `AwsAppRunnerObservabilityConfiguration`
- `AwsTransitGatewayVpcAttachment` -- AwsTransitGateway is a prerequisite because an attachment cannot exist without the gateway it attaches to; AwsSubnet because the attachment provisions an ENI into at least one subnet (the VPC arrives transitively through the subnet's own prerequisites).
- `AwsTransitGatewayRouteTable` -- Only the gateway is a hard prerequisite: a route table can exist empty. Associations, propagations, and routes referencing attachments are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsBatchComputeEnvironment` -- A MANAGED compute environment always launches into VPC subnets, so the subnet is a hard deploy prerequisite (security groups are required only for the Fargate types -- scenario-declared, not a registry edge).
- `AwsBatchJobQueue` -- A job queue cannot exist without at least one VALID compute environment to map onto.
- `AwsBatchSchedulingPolicy`
- `AwsBatchJobDefinition`
- `AwsCodeBuildProject` -- CI/CD
- `AwsCodePipeline`
- `AwsMwaaEnvironment` -- Workflow / Orchestration AwsSubnet and AwsSecurityGroup are prerequisites because the environment's network interfaces are placed in referenced private subnets and AWS requires at least one attached security group at creation.
- `AwsNeptuneCluster` -- Graph Database
- `AwsMemorydbCluster` -- A cluster always launches into a subnet group; the subnets are the hard deploy prerequisite. The ACL it attaches is optional composition (the built-in "open-access" ACL needs no resource) -- scenarios declare the ACL/user chain via the e2e-prerequisites annotation.
- `AwsMemorydbUser`
- `AwsMemorydbAcl` -- An empty ACL is valid (MemoryDB has no mandatory "default" member), so the user is optional composition -- the composed scenario declares it via the e2e-prerequisites annotation, never a registry edge.
- `AwsMskCluster` -- Streaming AwsSubnet and AwsSecurityGroup are prerequisites because brokers are placed in referenced subnets and AWS requires at least one attached security group at creation.
- `AwsMskServerlessCluster` -- AwsSubnet is a prerequisite because the serverless cluster's network interfaces are placed in referenced subnets (security groups are optional -- AWS attaches the VPC default group when none are referenced).
- `AwsLambdaEventSourceMapping` -- AwsLambda is a prerequisite because a mapping cannot exist without the function it invokes (a required reference). Event sources (SQS, Kinesis, DynamoDB, MSK) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsSnsSubscription` -- AwsSnsTopic is a prerequisite because a subscription cannot exist without the topic it subscribes to (a required reference). Endpoints (SQS queues, Lambda functions, Firehose streams) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsPlantonRunner` -- AwsSubnet is a prerequisite because the runner appliance places its network interfaces into referenced subnets -- the placement reference must resolve before the appliance can deploy.
- `AwsRoute53HealthCheck`
- `AwsSesConfigurationSet` -- Both SES kinds are dependency-free leaves: an identity's configuration set is optional composition (scenarios declare it via the e2e-prerequisites annotation), and a configuration set's event destinations reference other kinds only optionally.
- `AwsSesEmailIdentity`
- `AwsSecretsManagerSecret` -- A dependency-free leaf: the KMS key, rotation Lambda, and external rotation role references are all optional composition -- scenarios declare them via the e2e-prerequisites annotation, never registry edges.
- `AwsOpenSearchServerlessCollection` -- A dependency-free leaf: the collection-scoped encryption/network/ data-access/retention policies are module-rendered, and the KMS key and data-access principal references are optional composition (e2e-prerequisites annotation).
- `AwsBedrockGuardrail` -- A dependency-free leaf: the KMS key reference is optional composition (e2e-prerequisites annotation); published versions are folded satellites of the guardrail itself.
- `AwsBedrockCustomModel` -- AwsIamRole is a prerequisite because Bedrock assumes the job role to read training data and write outputs; the S3 locations and KMS key are optional composition (e2e-prerequisites annotation).
- `AwsBedrockInferenceProfile` -- A dependency-free leaf: the model source is a foundation model or an AWS system-defined cross-region profile, never a customer resource.
- `AwsBedrockProvisionedThroughput` -- A dependency-free leaf in the registry: capacity is typically bought for an AwsBedrockCustomModel (the default reference), but foundation model ARNs are equally legal, so the edge is optional composition.
- `AwsBedrockModelAccess` -- A dependency-free leaf: the agreement covers an AWS-listed foundation model, never a customer resource.
- `AwsBedrockAgent` -- AwsIamRole is a prerequisite because the Bedrock service assumes the agent resource role to invoke models, action-group Lambdas, and knowledge bases; the guardrail, KMS key, provisioned throughput, and collaborator/knowledge-base edges are optional composition (e2e-prerequisites annotation). Action groups, aliases, collaborators, and knowledge-base associations are folded satellites of the agent.
- `AwsBedrockKnowledgeBase` -- AwsIamRole is a prerequisite because the Bedrock service assumes the knowledge-base role to read data sources, call the embedding model, and read/write the vector store; the vector-store and data-source reference edges (OpenSearch, S3, Secrets Manager, ...) are optional composition (e2e-prerequisites annotation). Data sources are folded satellites of the knowledge base.
- `AwsBedrockFlow` -- AwsIamRole is a prerequisite because the Bedrock service assumes the flow execution role to invoke the models, agents, knowledge bases, and Lambdas its nodes reference; the node-level reference edges are optional composition (e2e-prerequisites annotation).
- `AwsBedrockPrompt` -- A dependency-free leaf: variants target AWS-listed foundation models by ID; targeting another agent's alias is optional composition (e2e-prerequisites annotation).
- `AzureResourceGroup` -- 2000–2999: Azure resources
- `AzureAksCluster` -- AzureResourceGroup is the only required parent: the cluster is created inside a referenced resource group. Subnet is optional on the default node pool (AKS provisions managed networking when unset).
- `AzureAksNodePool` -- AzureAksCluster is a prerequisite because a node pool attaches to an existing cluster by ARM ID; the resource group chains transitively.
- `AzureContainerRegistry` -- AzureResourceGroup is a prerequisite because a container registry is created inside a resource group.
- `AzureDnsZone` -- AzureResourceGroup is a prerequisite because the DNS zone is created inside a referenced resource group that must already exist.
- `AzureKeyVault` -- AzureResourceGroup is a prerequisite because a key vault is created inside a referenced resource group in composed environments.
- `AzureVirtualNetwork` -- AzureResourceGroup is a prerequisite because a virtual network is created inside a referenced resource group in composed environments.
- `AzureNatGateway` -- AzureResourceGroup is a prerequisite because a NAT gateway is created inside a referenced resource group in composed environments.
- `AzureVirtualMachine` -- AzureNetworkInterface is a prerequisite because a virtual machine attaches at least one NIC (the subnet, network, and resource group chain transitively through the NIC's own prerequisites).
- `AzureStorageAccount` -- AzureResourceGroup is a prerequisite because a storage account is created inside a referenced resource group in composed environments.
- `AzureDnsRecord` -- AzureDnsZone is a prerequisite because a record set is created inside a referenced zone (the resource group chains transitively through the zone). Public DNS zone names are not globally unique, so a shared zone fixture is safe to recreate across scenarios.
- `AzureSubnet` -- AzureVirtualNetwork is a prerequisite because a subnet is an ARM child of a referenced network -- the network must exist before the subnet can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzureNetworkSecurityGroup` -- AzureResourceGroup is a prerequisite because a network security group is created inside a referenced resource group in composed environments.
- `AzurePublicIp` -- AzureResourceGroup is a prerequisite because a public IP is created inside a referenced resource group in composed environments.
- `AzurePrivateEndpoint` -- AzureSubnet is a prerequisite because a private endpoint draws its private IP from a referenced subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite). The connection target is polymorphic and the DNS zones / ASGs are optional, so none of those are prerequisites.
- `AzurePrivateDnsZone` -- AzureResourceGroup is a prerequisite because a private DNS zone is created inside a referenced resource group in composed environments.
- `AzureApplicationGateway` -- AzureSubnet is a prerequisite because a gateway cannot exist without its dedicated gateway_ip_configuration subnet (the network and resource group chain transitively through the subnet's own prerequisites); public frontends additionally reference a public IP, but private-only gateways are legal, so it is not a registry prerequisite.
- `AzureLoadBalancer` -- AzureResourceGroup is a prerequisite because a load balancer is created inside a referenced resource group (frontends additionally reference subnets or public IPs, but neither is universally required, so they are not registry prerequisites).
- `AzureRouteTable` -- AzureResourceGroup is a prerequisite because a route table is created inside a referenced resource group in composed environments.
- `AzurePrivateDnsZoneVirtualNetworkLink` -- AzurePrivateDnsZone and AzureVirtualNetwork are prerequisites because a virtual network link is a child resource of a referenced zone and binds it to a referenced network -- both must exist before the link can be written. (The resource group arrives transitively through the zone's and network's own prerequisite declarations.)
- `AzureVirtualNetworkPeering` -- AzureVirtualNetwork is a prerequisite because a peering is an ARM child of its local network and binds it to a remote network -- the local network must exist before the peering can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzurePublicIpPrefix` -- AzureResourceGroup is a prerequisite because a public IP prefix is created inside a referenced resource group in composed environments.
- `AzureNetworkInterface` -- AzureSubnet is a prerequisite because a network interface's IP configurations deploy into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureManagedDisk` -- AzureResourceGroup is a prerequisite because a managed disk is created inside a resource group.
- `AzureVirtualMachineScaleSet` -- AzureSubnet is a prerequisite because every scale-set instance's network interface deploys into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureKeyVaultKey` -- AzureKeyVault is a prerequisite because a key is a data-plane object inside a referenced vault -- the vault must exist before the key can be written (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultCertificate` -- AzureKeyVault is a prerequisite because a certificate is a data-plane object inside a referenced vault -- the vault must exist before the certificate can be enrolled or imported (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultSecret` -- AzureKeyVault is a prerequisite because a secret is a data-plane object inside a referenced vault -- the vault must exist before the secret can be written (the resource group chains transitively through the vault's own prerequisite). Part of the Key Vault family (2005, 2025-2026) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureWebApplicationFirewallPolicy` -- AzureResourceGroup is a prerequisite because a WAF policy is created inside a referenced resource group; the Application Gateways that attach the policy reference it, never the reverse.
- `AzureApplicationSecurityGroup` -- AzureResourceGroup is a prerequisite because an application security group is created inside a referenced resource group; network interfaces, scale-set IP configurations, and NSG security rules reference the group, never the reverse.
- `AzureDiskEncryptionSet` -- AzureKeyVaultKey is a prerequisite because a disk encryption set wraps customer data with a referenced key -- the key (and its vault, which chains transitively) must exist before the set can resolve the key URL at create time.
- `AzurePostgresqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureRedisCache` -- AzureResourceGroup is a prerequisite because the cache is created inside a referenced resource group (VNet injection additionally references a dedicated subnet, but only the Premium tier supports it, so it is not a registry prerequisite).
- `AzureCosmosdbAccount` -- AzureResourceGroup is a prerequisite because the account is created inside a referenced resource group.
- `AzureMssqlServer` -- AzureResourceGroup is a prerequisite because the logical server is created inside a referenced resource group.
- `AzureMysqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureMssqlDatabase` -- The parent logical server is referenced via server_id, not auto-deployed: E2E scenarios declare their own server fixture (minimal-server.yaml or the pool-attach chain through AzureMssqlElasticPool) so sequential subtests never destroy and recreate the same globally unique server_name.
- `AzureMssqlElasticPool` -- AzureMssqlServer is a prerequisite because every elastic pool lives on a referenced logical server (the server's resource group is transitive).
- `AzureRedisLinkedServer` -- The target and linked caches are referenced via ARM ids, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicy` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicyAssignment` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureContainerAppEnvironment` -- AzureResourceGroup is a prerequisite because the environment is created inside a referenced resource group that must already exist.
- `AzureContainerApp` -- AzureContainerAppEnvironment is a prerequisite because every app runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureServicePlan` -- AzureResourceGroup is a prerequisite because the plan is created inside a referenced resource group that must already exist.
- `AzureFunctionApp` -- AzureServicePlan is a prerequisite because a function app runs on a referenced plan (the resource group arrives transitively through the plan). The required storage account is deliberately NOT a registry prerequisite: storage-account names are globally unique, so scenarios bring their own scenario-local account fixtures.
- `AzureLinuxWebApp` -- AzureServicePlan is a prerequisite because a web app runs on a referenced plan (the resource group arrives transitively through the plan).
- `AzureContainerAppJob` -- AzureContainerAppEnvironment is a prerequisite because a job runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureContainerAppEnvironmentStorage` -- AzureContainerAppEnvironment is a prerequisite because the storage registration lives on a referenced environment. The Azure Files share and storage account are deliberately NOT registry prerequisites: storage-account names are globally unique, so scenarios bring their own scenario-local account + share fixtures.
- `AzureContainerAppEnvironmentDaprComponent` -- AzureContainerAppEnvironment is a prerequisite because the Dapr component is registered on a referenced environment.
- `AzureContainerAppEnvironmentCertificate` -- AzureContainerAppEnvironment is a prerequisite because the certificate is stored on a referenced environment.
- `AzureContainerAppEnvironmentManagedCertificate` -- AzureContainerAppEnvironment is a prerequisite because the managed certificate is provisioned on a referenced environment.
- `AzureLogAnalyticsWorkspace` -- AzureResourceGroup is a prerequisite because the workspace is created inside a referenced resource group that must already exist.
- `AzureApplicationInsights` -- AzureLogAnalyticsWorkspace is a prerequisite because workspace-based Application Insights stores its telemetry in a referenced workspace (the resource group chains transitively through the workspace).
- `AzureMonitorDiagnosticSetting` -- AzureLogAnalyticsWorkspace is a prerequisite because the setting's scenarios route a fixture workspace's telemetry (the workspace doubles as target and destination); the target itself is polymorphic.
- `AzureMonitorActionGroup` -- AzureResourceGroup is a prerequisite because the action group is created inside a referenced resource group that must already exist.
- `AzureMonitorMetricAlert` -- AzureMonitorActionGroup is a prerequisite because a metric alert's actions fire into a referenced action group (the resource group chains transitively); alert scopes are polymorphic.
- `AzureMonitorScheduledQueryAlert` -- AzureLogAnalyticsWorkspace is a prerequisite because the rule queries a referenced workspace scope; AzureMonitorActionGroup because its action fires into a referenced action group.
- `AzureMonitorActivityLogAlert` -- AzureMonitorActionGroup is a prerequisite because an activity log alert's actions fire into a referenced action group (the resource group chains transitively). The alert itself is subscription-global and its scopes are polymorphic.
- `AzureApplicationInsightsStandardWebTest` -- AzureApplicationInsights is a prerequisite because a standard web test binds to a referenced Application Insights component (the resource group chains transitively through the component).
- `AzureUserAssignedIdentity` -- AzureResourceGroup is a prerequisite because the identity is created inside a referenced resource group that must already exist.
- `AzureRoleAssignment` -- AzureResourceGroup and AzureUserAssignedIdentity are prerequisites because an assignment grants a role at a referenced scope (most commonly a resource group) to a referenced principal (most commonly a managed identity) -- both must exist before the grant can be written.
- `AzureRoleDefinition` -- AzureResourceGroup is a prerequisite because a custom role definition is created at a referenced scope, most commonly a resource group in composed environments -- the scope must exist before the definition can be written.
- `AzureFederatedIdentityCredential` -- AzureUserAssignedIdentity is the prerequisite because a federated identity credential is a child resource of a referenced managed identity -- the identity must exist before the credential can be written on it. (The resource group arrives transitively through the identity's own prerequisite declaration.)
- `AzureServiceBusNamespace` -- AzureResourceGroup is a prerequisite because a Service Bus namespace is created inside a referenced resource group in composed environments. The namespace is the container every Service Bus messaging entity (queue, topic, subscription, authorization rule, geo-DR pairing) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureEventHubNamespace` -- AzureResourceGroup is a prerequisite because an Event Hub namespace is created inside a referenced resource group in composed environments. The namespace is the container every Event Hubs entity (event hub, consumer group, authorization rule, schema group, geo-DR pairing, customer-managed key) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureServiceBusQueue`
- `AzureServiceBusTopic`
- `AzureServiceBusSubscription`
- `AzureServiceBusAuthorizationRule`
- `AzureServiceBusDisasterRecoveryConfig`
- `AzureEventHub`
- `AzureEventHubConsumerGroup`
- `AzureEventHubAuthorizationRule`
- `AzureFrontDoorProfile` -- AzureResourceGroup is a prerequisite because a Front Door profile is created inside a referenced resource group in composed environments. The profile is the container every Front Door delivery resource (endpoint, origin group, origin, route) nests under.
- `AzureFrontDoorEndpoint` -- AzureFrontDoorProfile is a prerequisite because an endpoint is an ARM child of a referenced profile -- the profile must exist before the endpoint can be written. (The resource group arrives transitively through the profile's own prerequisite declaration.)
- `AzureFrontDoorOriginGroup` -- AzureFrontDoorProfile is a prerequisite because an origin group is an ARM child of a referenced profile.
- `AzureFrontDoorOrigin` -- AzureFrontDoorOriginGroup is a prerequisite because an origin is an ARM child of a referenced origin group (the profile and resource group chain transitively).
- `AzureFrontDoorRoute` -- A route attaches to an endpoint (its ARM parent) and forwards to an origin group whose origins must exist before ARM accepts the route -- so both the endpoint and the origin chain are genuine deploy-order prerequisites.
- `AzureFrontDoorRuleSet` -- AzureFrontDoorProfile is a prerequisite because a rule set is an ARM child of a referenced profile. The rules live inside the set (they form one ordered policy document); routes attach the set by ARM ID.
- `AzureFrontDoorCustomDomain` -- AzureFrontDoorProfile is a prerequisite because a custom domain is an ARM child of a referenced profile. The DNS zone and (for bring-your-own certificates) the Front Door secret are optional references, not deploy-order prerequisites.
- `AzureFrontDoorSecret` -- AzureFrontDoorSecret is a prerequisite-light kind: only the profile (its ARM parent) must exist. The Key Vault certificate it wraps is a reference resolved before the module runs; its vault chain is exercised through scenario-local fixtures in E2E.
- `AzureFrontDoorFirewallPolicy` -- AzureResourceGroup is a prerequisite because the Front Door WAF policy is created inside a referenced resource group -- it is a GLOBAL resource, not a profile child (a different ARM type than the regional Application Gateway WAF policy). Security policies attach it to profiles; the policy itself depends on nothing else.
- `AzureFrontDoorSecurityPolicy` -- A security policy is an ARM child of a profile that associates a referenced WAF policy with referenced domains -- so the endpoint (the default-domain association target; the profile arrives transitively through it) and the WAF policy are genuine deploy-order prerequisites.
- `AzureStorageContainer` -- None of the storage data-service kinds declares a registry prerequisite on AzureStorageAccount: account names are GLOBALLY unique and Azure holds a just-deleted name, so a recreate-per-scenario fixture would hang -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the storage_account_id reference itself.
- `AzureStorageShare`
- `AzureStorageQueue`
- `AzureStorageTable`
- `AzureStorageEncryptionScope`
- `AzureStorageDataLakeGen2Filesystem`
- `AzureStorageLocalUser`
- `AzureStorageObjectReplication`
- `AzureCosmosdbSqlDatabase` -- None of the Cosmos DB data-service kinds declares a registry prerequisite on AzureCosmosdbAccount: account names are GLOBALLY unique DNS labels, so a recreate-per-scenario fixture would risk name-reuse hangs -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the cosmosdb_account_id / parent-database references themselves.
- `AzureCosmosdbSqlContainer`
- `AzureCosmosdbMongoDatabase`
- `AzureCosmosdbMongoCollection`
- `AzureCosmosdbSqlRoleDefinition`
- `AzureCosmosdbSqlRoleAssignment`
- `AzureManagedRedis` -- AzureResourceGroup is the cluster's only registry prerequisite: the cluster is created inside a referenced resource group. The geo-replication and access-policy-assignment children declare NO prerequisite on AzureManagedRedis: clusters are expensive, slow-provisioning parents, so their E2E scenarios declare scenario-local cluster fixtures instead of recreating a shared one per scenario. Deploy ordering in composed environments still flows from the managed_redis_id references themselves.
- `AzureManagedRedisGeoReplication`
- `AzureManagedRedisAccessPolicyAssignment`
- `AzureEventHubDisasterRecoveryConfig`
- `AzureEventHubSchemaGroup`
- `AzureEventHubCluster` -- AzureResourceGroup is a prerequisite because a dedicated Event Hubs cluster is created inside a referenced resource group in composed environments. Note: clusters cannot be deleted for 4 hours after creation (Azure's moratorium), so E2E treats this kind as offline-gated.
- `AzureEventHubNamespaceCustomerManagedKey`
- `AzureMssqlFailoverGroup` -- AzureMssqlServer is a prerequisite because a failover group is created on a referenced primary logical server and points at a partner server; the primary (and its resource group, which chains transitively) must exist before the group can be written.
- `AzureContainerAppCustomDomain` -- AzureContainerApp is a prerequisite because the domain binding lives in a referenced app's ingress configuration (the environment and resource group chain transitively through the app).
- `AzureFirewallPolicy`
- `AzureFirewallPolicyRuleCollectionGroup` -- AzureFirewallPolicy is a prerequisite because a rule collection group is a child document of a referenced policy (the resource group chains transitively through the policy).
- `AzureFirewall` -- AzureSubnet is a prerequisite because a VNet-deployed firewall's data path lives in a dedicated subnet that must be named exactly "AzureFirewallSubnet" (the virtual network and resource group chain transitively through the subnet). The E2E install profile publishes a fixture subnet with that exact name and a /26 prefix.
- `AzureIpGroup`
- `AzureVirtualNetworkGateway` -- AzureSubnet is a prerequisite because every virtual network gateway lives in a dedicated subnet named exactly "GatewaySubnet" (the virtual network and resource group chain transitively through the subnet); the subnet install profile publishes a fixture instance with that exact ARM name. AzurePublicIp is a prerequisite because a VPN-type gateway (the default shape) requires a public IP per ip configuration; the address install profile publishes a dedicated zone-redundant instance (a gateway binds its address exclusively, and the AZ gateway SKUs require zones on it).
- `AzureVirtualNetworkGatewayConnection` -- Both gateways are prerequisites: a connection joins a virtual network gateway to a far side, and the site-to-site far side is a local network gateway (the GatewaySubnet, VNet, and resource group chain transitively through the virtual network gateway).
- `AzureLocalNetworkGateway`
- `AzurePrivateLinkService` -- AzureSubnet is the sole prerequisite: every NAT ip configuration draws its address from a subnet with private-link-service network policies disabled (the subnet install profile publishes a fixture instance with that flag off). The Standard load balancer whose frontend the service typically fronts is NOT a registry prerequisite -- the spec's destination is an exactly-one-of (load balancer frontend OR fixed destination IP), so scenarios that use the load-balancer shape declare it via the planton.dev/e2e-prerequisites annotation instead.
- `AzureExpressRouteCircuit`
- `AzureExpressRouteCircuitPeering` -- The circuit is the prerequisite: a peering is an ARM child of the circuit, addressed by the circuit's name (the resource group chains transitively through the circuit).
- `AzureExpressRouteGateway` -- The hub is the prerequisite: ARM requires an ExpressRoute Gateway to be deployed INTO a Virtual WAN hub (the WAN and resource group chain transitively through the hub).
- `AzureExpressRoutePort` -- ExpressRoute Port: your own physical port pair on a Microsoft edge router (ExpressRoute Direct), from whose bandwidth circuits are carved. Self-contained -- only the resource group is required.
- `AzureVirtualWan` -- Virtual WAN: the umbrella of Azure's managed hub-and-spoke networking, under which virtual hubs and their gateways are created. Self-contained -- only the resource group is required.
- `AzureVirtualHub` -- The WAN is the prerequisite: this kind models the Virtual WAN hub (virtual_wan_id is required; standalone hubs are the legacy Route Server construction, which has its own ARM surface). The resource group chains transitively through the WAN.
- `AzureVirtualHubConnection` -- Both sides of the attachment are prerequisites: the hub being joined and the spoke virtual network being attached.
- `AzureVpnGateway` -- The hub is the prerequisite: ARM deploys a Virtual WAN VPN gateway INTO a virtual hub (virtual_hub_id is required and immutable; the WAN and resource group chain transitively through the hub). ARM allows one VPN gateway per hub.
- `AzureVpnGatewayConnection` -- Both ends of the tunnel are prerequisites: a connection is an ARM child of the VPN gateway and pins each of its links to a specific link of the remote VPN site (the hub, WAN, and resource group chain transitively through the gateway).
- `AzureVpnSite` -- The WAN is the prerequisite: a VPN site is the Virtual WAN world's address-book entry for one branch location (virtual_wan_id is required; the classic-world sibling without a WAN is AzureLocalNetworkGateway). The resource group chains transitively through the WAN.
- `AzurePointToSiteVpnGateway` -- The hub and the server configuration are both prerequisites: a point-to-site VPN gateway deploys INTO a virtual hub (one P2S gateway per hub, a slot separate from the hub's site-to-site VPN gateway) and is born pointing at the VPN server configuration that defines how its users authenticate -- both ARM-required and fixed at creation. The WAN and resource group chain transitively through the hub.
- `AzureVpnServerConfiguration` -- Self-contained -- only the resource group is required: a VPN server configuration is the reusable "who may connect and how" authentication policy (Entra ID / certificate / RADIUS) that point-to-site VPN gateways attach to; it references no other Azure resource.
- `AzureCognitiveAccount` -- Self-contained -- only the resource group is required: an Azure AI services account (Azure OpenAI, the multi-service AIServices account, the single-service accounts) needs no other Azure resource; subnets (network rules), Key Vault keys (CMK), storage accounts and user-assigned identities are optional references.
- `AzureCognitiveDeployment` -- An ARM child of its account: a model deployment (which model runs, at which throughput class) exists only on an Azure AI services account of kind "OpenAI" or "AIServices".
- `AzureCognitiveAccountProject` -- An ARM child of its account: an AI Foundry project exists only on an "AIServices"-kind account with project management enabled.
- `AzureMachineLearningWorkspace` -- The workspace REQUIRES all three companion services at creation (default storage, secrets vault, telemetry) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureMachineLearningDatastore` -- An ARM child of its workspace. The storage target (container, filesystem or share) is scenario-declared via the e2e-prerequisites annotation -- only the blob scenario needs a container, so it is not a kind-wide prerequisite.
- `AzureMachineLearningComputeCluster` -- An ARM child of its workspace (.../computes/{name}) -- the auto-scaling pool of VMs training jobs run on.
- `AzureMachineLearningComputeInstance` -- An ARM child of its workspace (.../computes/{name}) -- a single always-on VM serving as one data scientist's cloud workstation.
- `AzureAiFoundry` -- The hub REQUIRES both companion services at creation (secrets vault, default storage) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureAiFoundryProject` -- Deploys into its hub's resource group (the provider derives the group from the hub reference -- the project spec carries none).
- `AzureSearchService`
- `AzureMachineLearningOnlineEndpoint` -- An ARM child of its workspace (.../onlineEndpoints/{name}) -- the stable scoring address applications call. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningOnlineDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the running copy of a model the endpoint's traffic map routes to.
- `AzureMachineLearningBatchEndpoint` -- An ARM child of its workspace (.../batchEndpoints/{name}) -- the stable address batch scoring jobs are submitted to. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningBatchDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the job recipe (model, compute, batching behavior) the endpoint's default-deployment pointer routes submissions to.
- `AzureRecoveryServicesVault` -- The Recovery Services vault (Microsoft.RecoveryServices/vaults) -- the safe that classic Azure Backup data and Site Recovery configuration live in. Backup policies and protected items are ARM children of a vault.
- `AzureBackupPolicyVm` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern IaaS VM backups.
- `AzureBackupProtectedVm` -- An ARM child of its vault (.../protectedItems/...) -- the binding that puts one virtual machine under a backup policy's protection.
- `AzureBackupPolicyFileShare` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern Azure Files share backups (snapshot or vaulted).
- `AzureBackupProtectedFileShare` -- An ARM child of its vault (.../protectedItems/AzureFileShare;...) -- the binding that puts one Azure Files share under a backup policy's protection. The share's storage account must already be registered with the vault (AzureBackupContainerStorageAccount).
- `AzureDataProtectionBackupVault` -- The Data Protection backup vault (Microsoft.DataProtection/ backupVaults) -- the safe that MODERN Azure Backup data lives in (managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, Data Lake storage). Backup policies and backup instances are ARM children of a vault.
- `AzureDataProtectionBackupPolicy` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules for ONE Data Protection datasource type (blob storage, disk, Kubernetes cluster, MySQL/PostgreSQL flexible server, or Data Lake storage), modeled as one kind with variant blocks.
- `AzureDataProtectionBackupInstance` -- An ARM child of its vault (.../backupInstances/{name}) -- the binding that puts ONE datasource (a managed disk, a storage account's blob services, an AKS cluster, a MySQL/PostgreSQL flexible server, or a Data Lake storage account) under a Data Protection backup policy, modeled as one kind with variant blocks. The vault's managed identity must hold the datasource roles Azure Backup requires BEFORE the instance is created.
- `AzureBastionHost` -- AzureSubnet and AzurePublicIp are prerequisites because a dedicated-infrastructure Bastion host (Basic/Standard/Premium -- the default shapes) deploys into a subnet named exactly "AzureBastionSubnet" and binds a Standard static public IP EXCLUSIVELY (the virtual network and resource group chain transitively through the subnet). The Developer SKU instead attaches to a virtual network directly and uses neither.
- `AzureNetworkWatcherFlowLog` -- AzureVirtualNetwork and AzureStorageAccount are prerequisites because a flow log records a network-scoped target (a virtual network in the common case; subnets and network interfaces chain through the network) into a referenced storage account. The regional Network Watcher parent is NOT a prerequisite: Azure auto-creates it ("NetworkWatcher_{region}" in "NetworkWatcherRG") the moment the region hosts a virtual network, and the flow log references it by name. Traffic Analytics' Log Analytics workspace is an optional arm, declared by scenarios that use it.
- `AzurePrivateDnsResolver` -- AzureVirtualNetwork and AzureSubnet are prerequisites because a DNS Private Resolver anchors to a referenced virtual network (at most ONE resolver per network -- Azure enforces it) and each of its inbound/outbound endpoints occupies its own dedicated subnet delegated to "Microsoft.Network/dnsResolvers" (the resource group chains transitively through the network and subnets).
- `AzurePrivateDnsResolverForwardingRuleset` -- AzurePrivateDnsResolver is a prerequisite because a DNS forwarding ruleset steers a resolver's OUTBOUND endpoints -- it binds their ARM ids (at most 2, same resolver) at creation. (The resource group and network chain transitively through the resolver's own prerequisite declarations.)
- `AzureBackupContainerStorageAccount` -- Registers a storage account with a Recovery Services vault as a backup container (.../protectionContainers/StorageContainer;...) -- one registration per storage-account-and-vault pair, required BEFORE any of the account's file shares can be protected. Part of the backup family (2175-2179) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureDataProtectionResourceGuard` -- The Data Protection Resource Guard (Microsoft.DataProtection/ resourceGuards) -- the approval gate behind Multi-User Authorization: privileged vault operations (disabling soft delete, reducing retention) require an approval through a guard, which typically lives in a DIFFERENT administrator's scope. Vaults reference a guard by its ARM ID. Part of the backup family (2175-2182) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzurePrivateDnsResolverVirtualNetworkLink` -- The attachment that makes a DNS forwarding ruleset take effect in one virtual network ({ruleset_id}/virtualNetworkLinks/{name}) -- one link per ruleset-network pair, up to 500 per ruleset, spokes joining and leaving independently (which is why the link is a standalone kind, exactly like AzurePrivateDnsZoneVirtualNetworkLink). Part of the DNS Private Resolver family (2186-2187) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `GcpArtifactRegistryRepo` -- 3000–3999: GCP resources
- `GcpTargetHttpsProxy` -- The URL map is the parent a proxy cannot exist without; the classic compute certificate kinds and the SSL policy are the fixture parents the committed scenarios attach. The Certificate Manager certificate list (certificate_manager_certificates, honored only by the cross-region internal ALB) is optional composition -- a scenario that arms it declares GcpCertManagerCert via the e2e-prerequisites annotation, never a registry edge that would tax every proxy and forwarding-rule chain.
- `GcpCloudFunction`
- `GcpCloudRun`
- `GcpCloudSql`
- `GcpDnsZone`
- `GcpGcsBucket`
- `GcpGkeCluster`
- `GcpIamCustomRole`
- `GcpProject`
- `GcpVpcNetwork`
- `GcpSubnetwork`
- `GcpRouterNat`
- `GcpGkeNodePool`
- `GcpServiceAccount`
- `GcpGkeWorkloadIdentityBinding`
- `GcpCertManagerCert`
- `GcpComputeInstance`
- `GcpDnsRecord`
- `GcpProjectIamMember`
- `GcpFirewallRule`
- `GcpGlobalAddress`
- `GcpCloudArmorPolicy`
- `GcpHealthCheck`
- `GcpBackendBucket`
- `GcpBackendService`
- `GcpRegionNetworkEndpointGroup`
- `GcpUrlMap`
- `GcpManagedSslCertificate`
- `GcpTargetHttpProxy`
- `GcpAlloydbCluster`
- `GcpRedisInstance`
- `GcpFirestoreDatabase`
- `GcpSpannerInstance`
- `GcpSpannerDatabase`
- `GcpBigtableInstance`
- `GcpMemorystoreInstance`
- `GcpCloudSqlDatabase`
- `GcpCloudSqlUser`
- `GcpAlloydbInstance`
- `GcpAlloydbUser`
- `GcpSpannerBackupSchedule`
- `GcpBigtableTable`
- `GcpFirestoreBackupSchedule`
- `GcpFirestoreIndex`
- `GcpBigQueryDataset`
- `GcpDataprocCluster`
- `GcpDataprocAutoscalingPolicy`
- `GcpBigQueryTable`
- `GcpPubSubTopic`
- `GcpPubSubSubscription`
- `GcpCloudTasksQueue`
- `GcpCloudSchedulerJob`
- `GcpPubSubSchema`
- `GcpVertexAiNotebook`
- `GcpVertexAiEndpoint`
- `GcpVertexAiIndex`
- `GcpVertexAiIndexEndpoint` -- Vector Search IndexEndpoint — distinct from the online-prediction GcpVertexAiEndpoint (671); different GCP resources, different kinds.
- `GcpVertexAiDeployedIndex`
- `GcpCloudComposerEnvironment`
- `GcpCloudComposerUserWorkloadsSecret`
- `GcpCloudComposerUserWorkloadsConfigMap`
- `GcpKmsKeyRing`
- `GcpKmsKey`
- `GcpKmsKeyIamMember`
- `GcpFilestoreInstance`
- `GcpWorkloadIdentityPool` -- 3101–3109: IAM/identity family (overflow block; the 3000–3022 foundation/security sub-band is fully allocated)
- `GcpWorkloadIdentityPoolProvider`
- `GcpServiceAccountIamMember`
- `GcpGlobalForwardingRule` -- 3110–3119: networking/load-balancer family (overflow block; the 3023–3029 LB sub-band is fully allocated)
- `GcpSslPolicy`
- `GcpSslCertificate`
- `GcpServiceNetworkingConnection`
- `GcpAddress`
- `GcpServiceConnectionPolicy`
- `GcpCertManagerDnsAuthorization`
- `GcpCertificateMap` -- GcpCertManagerCert is a prerequisite because a map entry binds hostnames to EXISTING certificates — the canonical map references a certificate fixture's resource name.
- `GcpCloudRunJob` -- 3120–3129: GCP serverless overflow
- `GcpServerlessVpcConnector`
- `GcpComputeDisk` -- 3130–3139: GCP compute overflow (the 3000–3022 foundation sub-band that holds GcpComputeInstance is fully allocated)
- `GcpComputeMig` -- GcpVpcNetwork is a prerequisite because the canonical group runs its fleet on a dedicated custom-mode VPC — a managed instance group's template must attach every VM to a network, and the default VPC is never assumed.
- `GcpMonitoringNotificationChannel` -- 3140–3149: GCP observability & log routing
- `GcpMonitoringAlertPolicy` -- GcpMonitoringNotificationChannel is a prerequisite because the policy's canonical shape references a channel to notify — a policy without a delivery endpoint measures but never pages.
- `GcpMonitoringUptimeCheck`
- `GcpLoggingSink` -- GcpGcsBucket is a prerequisite because the canonical sink exports to a Cloud Storage bucket — the cheapest destination that proves the whole writer-identity grant flow.
- `GcpMonitoringDashboard`
- `GcpMonitoringSlo`
- `GcpLogBucket`
- `GcpLogMetric`
- `GcpSecretManagerSecret` -- 3150–3159: GCP security & identity GcpServiceAccount is a prerequisite because the canonical secret grants secretAccessor to a workload service account — the access story the kind exists to model.
- `GcpIdentityPlatformConfig`
- `GcpIdentityPlatformTenant` -- GcpIdentityPlatformConfig is a prerequisite because tenants exist only in projects whose Identity Platform config enables multi_tenant.allow_tenants — a tenant without the initialized, tenant-enabled project config cannot be created at all.
- `GcpIamOauthClient`
- `GcpIamDenyPolicy`
- `GcpCloudRunDomainMapping` -- 3160–3169: GCP serverless edge GcpCloudRun is a prerequisite because a domain mapping exists only to point a verified domain at a running Cloud Run service — the route it maps must already exist for the mapping to be created at all.
- `GcpWorkflow`
- `GcpEventarcTrigger` -- GcpCloudRun is a prerequisite because the canonical trigger routes a Pub/Sub messagePublished event to a Cloud Run service — the destination story the kind exists to model.
- `GcpEventarcMessageBus`
- `KubernetesNamespace` -- 4000–4999: Kubernetes resources, organized in family sub-bands (4030–4069 also hosts CNI/autoscaling/DR addons; 4130–4149 hosts analytics & ML; 4190–4199 reserved for growth) 4000–4029: Kubernetes building blocks (core API primitives)
- `KubernetesDeployment`
- `KubernetesStatefulSet`
- `KubernetesDaemonSet`
- `KubernetesJob`
- `KubernetesCronJob`
- `KubernetesService`
- `KubernetesSecret`
- `KubernetesManifest`
- `KubernetesHelmRelease`
- `KubernetesConfigMap`
- `KubernetesServiceAccount`
- `KubernetesRbac` -- Bundles the RBAC grant grain (Role/ClusterRole + its binding) into one component: "grant these permissions to these subjects in this scope".
- `KubernetesIngress`
- `KubernetesNetworkPolicy`
- `KubernetesPersistentVolumeClaim`
- `KubernetesStorageClass`
- `KubernetesResourceQuota` -- Manages the namespace-governance pair: the ResourceQuota plus an optional companion LimitRange (per-object defaults/bounds) — two API objects, one governance story.
- `KubernetesPriorityClass`
- `KubernetesPodDisruptionBudget`
- `KubernetesHorizontalPodAutoscaler`
- `KubernetesCertManager` -- 4030–4069: Kubernetes foundation addons (certs, DNS, secrets, ingress, Gateway API, mesh, CNI/autoscaling/DR)
- `KubernetesClusterIssuer` -- KubernetesCertManager is a prerequisite for the three cert-manager CR kinds below: ClusterIssuer/Issuer/Certificate are cert-manager custom resources — without the controller and its CRDs they cannot be applied.
- `KubernetesIssuer`
- `KubernetesCertificate`
- `KubernetesExternalDns`
- `KubernetesExternalSecretsOperator`
- `KubernetesClusterSecretStore` -- KubernetesExternalSecretsOperator is a prerequisite for the three external-secrets CR kinds below: ClusterSecretStore/SecretStore/ ExternalSecret are external-secrets custom resources — without the operator and its CRDs they cannot be applied.
- `KubernetesSecretStore`
- `KubernetesExternalSecret`
- `KubernetesIngressNginx`
- `KubernetesGatewayApiCrds`
- `KubernetesGatewayClass`
- `KubernetesGateway`
- `KubernetesListenerSet`
- `KubernetesHttpRoute`
- `KubernetesGrpcRoute`
- `KubernetesTcpRoute`
- `KubernetesUdpRoute`
- `KubernetesTlsRoute`
- `KubernetesReferenceGrant`
- `KubernetesBackendTlsPolicy`
- `KubernetesIstioBaseCrds`
- `KubernetesIstio`
- `KubernetesDestinationRule` -- Istio API components (mesh traffic policy, security, telemetry). The seven typed resources below (4053–4059) require the Istio CRDs on the cluster, provided by the lightweight CRDs-only KubernetesIstioBaseCrds (851) — NOT the full mesh KubernetesIstio (852).
- `KubernetesServiceEntry`
- `KubernetesPeerAuthentication`
- `KubernetesRequestAuthentication`
- `KubernetesAuthorizationPolicy`
- `KubernetesTelemetry`
- `KubernetesEnvoyFilter`
- `KubernetesMetricsServer`
- `KubernetesCilium`
- `KubernetesKeda`
- `KubernetesKarpenter`
- `KubernetesKarpenterNodePool`
- `KubernetesKarpenterEc2NodeClass`
- `KubernetesClusterAutoscaler`
- `KubernetesVelero`
- `KubernetesKubePrometheusStack` -- 4070–4089: Kubernetes observability
- `KubernetesGrafana`
- `KubernetesSignoz` -- KubernetesClickHouse is a prerequisite because SigNoz stores every trace, metric and log in ClickHouse and deploys none of its own — the telemetry store is composed, never bundled.
- `KubernetesLoki`
- `KubernetesTempo`
- `KubernetesOtelOperator` -- The operator's admission webhooks (failurePolicy Fail) are served with a cert-manager Certificate in the default posture — cert-manager must be running before the operator installs.
- `KubernetesOtelCollector`
- `KubernetesKyverno` -- 4080–4099: Kubernetes security, policy, and identity
- `KubernetesGatekeeper`
- `KubernetesKeycloak` -- Keycloak declarations compose the official Keycloak Operator (which reconciles the Keycloak CR this kind renders) and, on the recommended postgres vendor, a KubernetesPostgres database — both must resolve before the CR can converge.
- `KubernetesOpenBao`
- `KubernetesOpenFga` -- OpenFGA requires a datastore; the recommended arm composes a KubernetesPostgres database (the sandbox memory arm needs nothing, but the registry declares the shape real deployments require).
- `KubernetesKeycloakOperator`
- `KubernetesCloudNativePgOperator` -- 4100–4129: Kubernetes data platforms
- `KubernetesPostgres`
- `KubernetesValkey`
- `KubernetesPerconaMysqlOperator`
- `KubernetesMysql`
- `KubernetesPerconaMongoOperator`
- `KubernetesMongodb`
- `KubernetesStrimziKafkaOperator`
- `KubernetesKafka` -- container_kind: a Strimzi Kafka cluster is a place in the provider's own model — KafkaTopic and KafkaUser declarations BELONG to one cluster (the strimzi.io/cluster label) and are drawn inside its box. Clients that merely talk to the cluster (Connect, MirrorMaker2, UI, Karapace) carry containment_exempt on their bootstrap/trust references.
- `KubernetesKafkaTopic`
- `KubernetesKafkaUser`
- `KubernetesKafkaConnect` -- container_kind: a Connect cluster hosts the connectors deployed INTO it (KafkaConnector's strimzi.io/cluster label names its Connect cluster) — the same room shape as KubernetesKafka above.
- `KubernetesKafkaConnector`
- `KubernetesKafkaMirrorMaker2`
- `KubernetesKarapace`
- `KubernetesKafkaUi`
- `KubernetesOpenSearchOperator`
- `KubernetesOpenSearch`
- `KubernetesAltinityOperator`
- `KubernetesClickHouse`
- `KubernetesSolrOperator`
- `KubernetesSolr`
- `KubernetesNeo4j`
- `KubernetesSeaweedFs`
- `KubernetesQdrant`
- `KubernetesRabbitMqOperator` -- The RabbitMQ Cluster Operator's release manifest ships admission webhooks whose serving certificate is a cert-manager Certificate — cert-manager must be running before the operator installs.
- `KubernetesRabbitMq`
- `KubernetesAirflow` -- 4130–4149: Kubernetes analytics and ML KubernetesPostgres is a prerequisite because Airflow's metadata database composes a KubernetesPostgres by default (the spec's FK defaults resolve onto its outputs) and the migration Job needs the database reachable before the server components start.
- `KubernetesSparkOperator`
- `KubernetesKubeRayOperator`
- `KubernetesRayCluster` -- KubernetesKubeRayOperator is a prerequisite because this kind declares the RayCluster custom resource that only the operator's CRDs admit and only the operator reconciles into head and worker pods.
- `KubernetesFlinkOperator` -- KubernetesCertManager is a prerequisite because the Flink operator's chart, with its default-on admission webhook, renders cert-manager Issuer/Certificate resources and trusts the API server through cert-manager's CA injection — there is no self-signed fallback at the pinned chart, and the webhooks are fail-closed.
- `KubernetesFlinkDeployment` -- KubernetesFlinkOperator is a prerequisite because this kind declares the FlinkDeployment custom resource that only the operator's CRDs admit and only the operator reconciles into a running Flink cluster.
- `KubernetesJupyterHub` -- KubernetesPostgres is a prerequisite because JupyterHub's hub database composes a KubernetesPostgres in its external-database arm (the spec's FK defaults resolve onto its outputs) and the hub pod mounts that database's credential Secret before it can start.
- `KubernetesMlflow` -- KubernetesPostgres is a prerequisite because MLflow's backend store composes a KubernetesPostgres in its production arm (FK defaults onto its outputs; the module composes the connection URI from its credential Secret), and KubernetesSeaweedFs because the artifact store's S3-compatible arm FK-defaults onto the SeaweedFS endpoint and credential Secret.
- `KubernetesTrino` -- KubernetesPostgres is a prerequisite because Trino's postgres catalogs compose a KubernetesPostgres (the catalog host and credential FK-default onto its outputs), and the pods read that database's credential Secret to resolve catalog passwords from environment.
- `KubernetesSuperset` -- KubernetesPostgres is a prerequisite because Superset's REQUIRED metadata database composes a KubernetesPostgres (FK defaults onto its outputs; the module composes the environment Secret from its credential Secret), and KubernetesValkey because the cache/broker arm FK-defaults onto a KubernetesValkey's service and password Secret.
- `KubernetesArgocd` -- 4150–4169: Kubernetes GitOps and CI/CD
- `KubernetesArgoWorkflows`
- `KubernetesTektonOperator`
- `KubernetesTekton` -- KubernetesTektonOperator is a prerequisite because this kind declares the TektonConfig custom resource that only the operator's CRDs admit and only the operator reconciles into running components.
- `KubernetesGhaRunnerScaleSetController`
- `KubernetesGhaRunnerScaleSet` -- KubernetesGhaRunnerScaleSetController is a prerequisite because this kind renders an AutoscalingRunnerSet custom resource that only the controller's CRDs admit and only the controller reconciles into listener and runner pods.
- `KubernetesHarbor`
- `KubernetesJenkins`
- `KubernetesTemporal` -- 4170–4189: Kubernetes app platforms KubernetesPostgres is a prerequisite because the recommended (and E2E-proven) database composition backs Temporal's default and visibility stores with a CloudNativePG cluster.
- `KubernetesNats`
- `KubernetesLocust`
- `DigitalOceanAppPlatformService` -- 5000–5999: DigitalOcean resources
- `DigitalOceanBucket`
- `DigitalOceanContainerRegistry`
- `DigitalOceanDatabaseCluster`
- `DigitalOceanDnsZone`
- `DigitalOceanDroplet`
- `DigitalOceanFirewall`
- `DigitalOceanFunction`
- `DigitalOceanKubernetesCluster`
- `DigitalOceanKubernetesNodePool`
- `DigitalOceanLoadBalancer`
- `DigitalOceanVolume`
- `DigitalOceanVpc`
- `DigitalOceanCertificate`
- `DigitalOceanDnsRecord`
- `CivoBucket` -- 6000–6999: Civo resources
- `CivoCertificate`
- `CivoComputeInstance`
- `CivoDatabase`
- `CivoDnsZone`
- `CivoFirewall`
- `CivoIpAddress`
- `CivoKubernetesCluster`
- `CivoKubernetesNodePool`
- `CivoVolume`
- `CivoVpc`
- `CivoDnsRecord`
- `CloudflareDnsZone` -- 7000–7999: Cloudflare resources
- `CloudflareKvNamespace`
- `CloudflareR2Bucket`
- `CloudflareWorker`
- `CloudflareLoadBalancer`
- `CloudflareD1Database`
- `CloudflareZeroTrustAccessApplication`
- `CloudflareDnsRecord`
- `CloudflareRuleset`
- `CloudflareWorkersKvPair`
- `CloudflareHyperdriveConfig`
- `CloudflareLoadBalancerPool`
- `CloudflareLoadBalancerMonitor`
- `CloudflareZeroTrustAccessPolicy`
- `CloudflareZeroTrustAccessGroup`
- `CloudflareQueue`
- `CloudflarePagesProject`
- `CloudflareZeroTrustTunnel`
- `CloudflareZeroTrustTunnelVirtualNetwork`
- `CloudflareZeroTrustTunnelRoute`
- `CloudflareList`
- `CloudflareListItem`
- `CloudflareTurnstileWidget`
- `CloudflareEmailRoutingZone`
- `CloudflareEmailRoutingRule`
- `CloudflareEmailRoutingAddress`
- `CloudflareOriginCaCertificate`
- `CloudflareCertificatePack`
- `CloudflareCustomHostname`
- `CloudflareCustomHostnameFallbackOrigin`
- `Auth0Connection` -- 8000–8999: Auth0 resources
- `Auth0Client`
- `Auth0EventStream`
- `Auth0ResourceServer`
- `Auth0Action`
- `Auth0Role`
- `OpenFgaStore` -- 9000–9999: OpenFGA resources Note: OpenFGA is Terraform-only - there is no Pulumi provider available. Pulumi modules for OpenFGA resources are pass-through placeholders.
- `OpenFgaAuthorizationModel`
- `OpenFgaRelationshipTuple`
- `OpenStackKeypair` -- 10000–10999: OpenStack resources
- `OpenStackNetwork`
- `OpenStackSubnet`
- `OpenStackRouter`
- `OpenStackRouterInterface`
- `OpenStackSecurityGroup`
- `OpenStackFloatingIp`
- `OpenStackNetworkPort`
- `OpenStackSecurityGroupRule`
- `OpenStackFloatingIpAssociate`
- `OpenStackInstance`
- `OpenStackServerGroup`
- `OpenStackVolume`
- `OpenStackVolumeAttach`
- `OpenStackProject`
- `OpenStackApplicationCredential`
- `OpenStackImage`
- `OpenStackRoleAssignment`
- `OpenStackLoadBalancer`
- `OpenStackLoadBalancerListener`
- `OpenStackLoadBalancerPool`
- `OpenStackLoadBalancerMember`
- `OpenStackLoadBalancerMonitor`
- `OpenStackDnsZone`
- `OpenStackDnsRecord`
- `ScalewayVpc`
- `ScalewayPrivateNetwork`
- `ScalewayPublicGateway`
- `ScalewayLoadBalancer`
- `ScalewayInstanceSecurityGroup`
- `ScalewayInstance`
- `ScalewayKapsuleCluster`
- `ScalewayKapsulePool`
- `ScalewayRdbInstance`
- `ScalewayRedisCluster`
- `ScalewayMongodbInstance`
- `ScalewayObjectBucket`
- `ScalewayBlockVolume`
- `ScalewayContainerRegistry`
- `ScalewayDnsZone`
- `ScalewayDnsRecord`
- `ScalewayServerlessFunction`
- `ScalewayServerlessContainer`
- `AliCloudLogProject`
- `AliCloudRamRole`
- `AliCloudRamPolicy`
- `AliCloudVpc`
- `AliCloudVswitch`
- `AliCloudSecurityGroup`
- `AliCloudEipAddress`
- `AliCloudNatGateway`
- `AliCloudApplicationLoadBalancer`
- `AliCloudNetworkLoadBalancer`
- `AliCloudVpnGateway`
- `AliCloudDnsZone`
- `AliCloudDnsRecord`
- `AliCloudPrivateDnsZone`
- `AliCloudStorageBucket`
- `AliCloudNasFileSystem`
- `AliCloudKmsKey`
- `AliCloudRdsInstance`
- `AliCloudPolardbCluster`
- `AliCloudRedisInstance`
- `AliCloudMongodbInstance`
- `AliCloudEcsInstance`
- `AliCloudContainerRegistry`
- `AliCloudKubernetesCluster`
- `AliCloudKubernetesNodePool`
- `AliCloudCdnDomain`
- `AliCloudFunction`
- `AliCloudSaeApplication`
- `AliCloudRocketmqInstance`
- `AliCloudCenInstance`
- `OciVcn`
- `OciSubnet`
- `OciSecurityGroup`
- `OciCompartment`
- `OciIdentityPolicy`
- `OciDynamicGroup`
- `OciComputeInstance`
- `OciContainerEngineCluster`
- `OciContainerEngineNodePool`
- `OciContainerInstance`
- `OciApplicationLoadBalancer`
- `OciNetworkLoadBalancer`
- `OciDynamicRoutingGateway`
- `OciPublicIp`
- `OciAutonomousDatabase`
- `OciDbSystem`
- `OciMysqlDbSystem`
- `OciPostgresqlDbSystem`
- `OciRedisCluster`
- `OciNosqlTable`
- `OciObjectStorageBucket`
- `OciFileSystem`
- `OciBlockVolume`
- `OciKmsVault`
- `OciKmsKey`
- `OciVaultSecret`
- `OciBastion`
- `OciFunctionsApplication`
- `OciApiGateway`
- `OciStreamPool`
- `OciQueue`
- `OciAlarm`
- `OciLogGroup`
- `OciDnsZone`
- `OciDnsRecord`
- `OciNetworkFirewall`
- `OciDevopsProject`
- `HetznerCloudSshKey`
- `HetznerCloudPlacementGroup`
- `HetznerCloudFirewall`
- `HetznerCloudNetwork`
- `HetznerCloudPrimaryIp`
- `HetznerCloudFloatingIp`
- `HetznerCloudServer`
- `HetznerCloudVolume`
- `HetznerCloudSnapshot`
- `HetznerCloudCertificate`
- `HetznerCloudLoadBalancer`
- `HetznerCloudDnsZone`

### spec.container.sidecars[].env.variables[].valueFrom.env

`string`

### spec.container.sidecars[].env.variables[].valueFrom.name

`string` · required

- rule: {"required":true}

### spec.container.sidecars[].env.variables[].valueFrom.fieldPath

`string`

### spec.container.sidecars[].env.variables[].configMapKeyRef

`ConfigMapKeyRef`

Reference to a key in a Kubernetes ConfigMap.

### spec.container.sidecars[].env.variables[].configMapKeyRef.name

`string` · required

Name of the ConfigMap.

- rule: {"required":true}

### spec.container.sidecars[].env.variables[].configMapKeyRef.key

`string` · required

Key within the ConfigMap.

- rule: {"required":true}

### spec.container.sidecars[].env.variables[].configMapKeyRef.optional

`bool`

If true, the env var is silently skipped when the ConfigMap or key does not exist
(instead of blocking pod startup).

### spec.container.sidecars[].env.variables[].fieldRef

`ObjectFieldRef`

Reference to a pod-level field (metadata.name, status.podIP, etc.).

### spec.container.sidecars[].env.variables[].fieldRef.apiVersion

`string`

Version of the schema. Defaults to "v1".

### spec.container.sidecars[].env.variables[].fieldRef.fieldPath

`string` · required

Path of the field to select (e.g., "metadata.name", "status.podIP").

- rule: {"required":true}

### spec.container.sidecars[].env.variables[].resourceFieldRef

`ResourceFieldRef`

Reference to container resource limits or requests (limits.cpu, requests.memory, etc.).

### spec.container.sidecars[].env.variables[].resourceFieldRef.containerName

`string`

Container name. Required for init containers; defaults to the current
container for regular containers.

### spec.container.sidecars[].env.variables[].resourceFieldRef.resource

`string` · required

Resource to select (e.g., "limits.cpu", "requests.memory").

- rule: {"required":true}

### spec.container.sidecars[].env.variables[].resourceFieldRef.divisor

`string`

Specifies the output format of the exposed resource.
For CPU: "1" means cores. For memory: "1", "1Ki", "1Mi", "1Gi".

### spec.container.sidecars[].env.secrets

`[]SecretEnvVar`

Individual secret environment variables (sensitive).

### spec.container.sidecars[].env.secrets[].name

`string` · required

The environment variable name.

- rule: Must be a valid C_IDENTIFIER (start with letter/underscore, contain only letters, digits, underscores)
- rule: {"required":true}

### spec.container.sidecars[].env.secrets[].value

`string`

Literal string value.
A Kubernetes Secret is automatically created and the environment variable
references that secret.

### spec.container.sidecars[].env.secrets[].secretRef

`KubernetesSecretKeyRef`

Reference to a key within an existing Kubernetes Secret.

### spec.container.sidecars[].env.secrets[].secretRef.namespace

`string`

The namespace of the Kubernetes Secret.
If not specified, defaults to the namespace where the component is deployed.
Note: Cross-namespace secret references may not be supported by all Helm charts.

### spec.container.sidecars[].env.secrets[].secretRef.name

`string` · required

The name of the Kubernetes Secret.

- rule: {"required":true}

### spec.container.sidecars[].env.secrets[].secretRef.key

`string` · required

The key within the Kubernetes Secret that contains the value.

- rule: {"required":true}

### spec.container.sidecars[].env.secrets[].secretRef.optional

`bool`

If true, the env var is silently skipped when the Secret or key does not exist
(instead of blocking pod startup).

### spec.container.sidecars[].env.secrets[].valueFrom

`ValueFromRef`

Reference to another Planton resource's secret output field.
The orchestrator resolves this before invoking IaC modules.

### spec.container.sidecars[].env.secrets[].valueFrom.kind

`enum`

Allowed values (use exactly as shown):

- `unspecified` -- 0: Default/unspecified
- `TestCloudResourceGeneric` -- 1–49: Test/dev/custom
- `TestCloudResourceKubernetes`
- `ConfluentKafka` -- 50–199: saas platform resources
- `AtlasMongodb`
- `SnowflakeDatabase`
- `AwsAlb` -- 1000–1999: AWS resources AwsSubnet is a prerequisite because an ALB requires at least two subnets in different availability zones -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsCertManagerCert`
- `AwsCloudFront`
- `AwsDynamodb`
- `AwsEcrRepo`
- `AwsEcsCluster`
- `AwsEcsService` -- AwsEcsCluster, AwsEcsTaskDefinition, and AwsSubnet are prerequisites because a service schedules a referenced task-definition revision into a referenced live cluster and places task network interfaces into referenced subnets -- all three references must resolve first.
- `AwsEksCluster` -- AwsSubnet and AwsIamRole are prerequisites because the control plane attaches its network interfaces into referenced subnets and assumes a referenced cluster role that must already carry AmazonEKSClusterPolicy.
- `AwsIamRole`
- `AwsLambda`
- `AwsRdsCluster`
- `AwsRdsInstance`
- `AwsRoute53Zone`
- `AwsS3Bucket`
- `AwsLbTargetGroup` -- AwsVpc is a prerequisite because a target group's health checks and target registrations live inside one VPC -- the spec's vpc_id reference must resolve before the group can be created.
- `AwsSecurityGroup` -- AwsVpc is a prerequisite because every security group is created in a VPC; the E2E install profile resolves vpc_id against the VPC prerequisite.
- `AwsVpc`
- `AwsEksNodeGroup` -- AwsEksCluster is a prerequisite because nodes register with a live control plane; AwsIamRole and AwsSubnet back the node role and worker subnet references.
- `AwsIamUser`
- `AwsKmsKey`
- `AwsEc2Instance`
- `AwsClientVpn` -- Every Client VPN endpoint requires an ACM server certificate at create time; the imported self-signed fixture satisfies it. Subnets/VPC are optional composition (a zero-association endpoint is valid) -- composed scenarios declare them via the e2e-prerequisites annotation.
- `AwsDocumentDb`
- `AwsRoute53DnsRecord` -- AwsRoute53Zone is a prerequisite because every record lives inside a hosted zone -- the spec's zone_id reference must resolve before the record can be created.
- `AwsS3ObjectSet` -- AwsS3Bucket is a prerequisite because the object set's bucket reference is required -- objects cannot exist without the bucket that holds them.
- `AwsSqsQueue`
- `AwsSnsTopic`
- `AwsEventBridgeBus`
- `AwsEventBridgeRule`
- `AwsIamOidcProvider`
- `AwsIamPolicy`
- `AwsIamInstanceProfile` -- AwsIamRole is a prerequisite because an instance profile is a wrapper that must contain a role to be useful -- the profile's spec requires a role reference, so the role must be deployed first.
- `AwsLbListener` -- AwsAlb and AwsLbTargetGroup are prerequisites because a listener is an attachment point on a load balancer and its default action almost always forwards to a target group -- both references must resolve before the listener can be created.
- `AwsLbListenerRule` -- AwsLbListener is a prerequisite because a rule only exists as an attachment on a listener -- the listener_arn reference must resolve before the rule can be created.
- `AwsLaunchTemplate`
- `AwsAutoScalingGroup` -- AwsSubnet and AwsLaunchTemplate are prerequisites because a group cannot exist without subnets to place capacity in and a launch template to launch from -- the spec's subnets and launch_template references must resolve before the group can be created.
- `AwsEksAddon` -- AwsEksCluster is a prerequisite because an add-on installs onto a live control plane -- the spec's cluster_name reference must resolve before the add-on can be created.
- `AwsEksFargateProfile` -- AwsEksCluster, AwsIamRole, and AwsSubnet are prerequisites because a Fargate profile attaches to a live control plane, runs pods as a referenced pod-execution role, and launches them into referenced private subnets -- all three references must resolve first.
- `AwsEksAccessEntry` -- AwsEksCluster and AwsIamRole are prerequisites because an access entry grants a referenced IAM principal access to a live control plane -- both references must resolve before the entry can be created.
- `AwsEcsTaskDefinition` -- AwsIamRole is a prerequisite because the kind's default posture -- Fargate with the awslogs logging default -- is rejected by AWS at registration time without an execution role the agent can assume.
- `AwsHttpApiGateway`
- `AwsStepFunction` -- AwsIamRole is a prerequisite because a state machine cannot be created without an execution role it can assume -- the spec's role_arn reference must resolve before the CreateStateMachine call.
- `AwsHttpApiVpcLink` -- AwsSubnet is a prerequisite because a VPC link is a set of managed ENIs provisioned into referenced subnets -- the subnet references must resolve before the link can be created. Security groups are optional on the link, so they compose per-scenario rather than as a registry prerequisite.
- `AwsHttpApiDomain` -- AwsCertManagerCert is a prerequisite because a custom domain cannot be created without a TLS certificate in the same region covering the domain -- the spec's certificate_arn reference must resolve first.
- `AwsVpcEndpoint` -- AwsVpcEndpoint's composed E2E scenarios reference the AwsVpc prerequisite's outputs (vpc_id + default_route_table_id for gateway endpoints) and the AwsSubnet pair's subnet_id outputs (interface endpoints), so both are genuine deploy-order prerequisites.
- `AwsElasticacheUser`
- `AwsElasticacheUserGroup` -- AwsElasticacheUser is a genuine prerequisite: AWS refuses to create a user group that does not contain a user named "default", so a group's composed E2E scenario must resolve a deployed user's outputs.
- `AwsRedshiftServerlessNamespace`
- `AwsRedshiftServerlessWorkgroup` -- The namespace is a genuine prerequisite: a workgroup attaches to exactly one namespace by name at create time, so its composed E2E scenario must resolve a deployed namespace's outputs. AwsSubnet is a prerequisite because Redshift Serverless requires the workgroup's subnets to span three availability zones.
- `AwsRedisElasticache` -- AwsSubnet is a prerequisite because the module builds an ElastiCache subnet group from referenced subnets -- the spec's subnet references must resolve before the replication group can deploy.
- `AwsOpenSearchDomain`
- `AwsMemcachedElasticache`
- `AwsServerlessElasticache`
- `AwsNlb` -- AwsSubnet is a prerequisite because an NLB requires at least one subnet mapping -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsElasticIp`
- `AwsTransitGateway`
- `AwsGlobalAccelerator`
- `AwsSubnet`
- `AwsInternetGateway`
- `AwsNatGateway` -- AwsInternetGateway is a prerequisite because a public NAT gateway can only become available once the VPC it sits in has an internet gateway attached (AWS rejects the create otherwise) -- so the gateway must be deployed first. AwsVpc is a prerequisite because a REGIONAL NAT gateway (availability_mode = regional) references the VPC directly instead of a subnet.
- `AwsEgressOnlyInternetGateway`
- `AwsElasticFileSystem` -- AwsSubnet and AwsSecurityGroup are prerequisites because mount targets (required, min 1) place the file system's NFS endpoints into subnets and attach security groups -- both references must resolve before the CreateMountTarget calls.
- `AwsEfsAccessPoint` -- AwsElasticFileSystem is a prerequisite because an access point is created INTO a file system -- the spec's required file_system_id reference must resolve before the CreateAccessPoint call.
- `AwsFsxLustreFileSystem`
- `AwsFsxOpenzfsFileSystem`
- `AwsFsxWindowsFileSystem` -- Every Windows file system must join an Active Directory domain; the directory itself is external infrastructure (AWS Managed Microsoft AD or a self-managed domain), so only the network dependency is a declarable prerequisite.
- `AwsFsxOntapFileSystem`
- `AwsFsxOntapStorageVirtualMachine`
- `AwsFsxOntapVolume`
- `AwsFsxDataRepositoryAssociation`
- `AwsCognitoUserPool`
- `AwsCognitoIdentityProvider` -- AwsCognitoUserPool is a prerequisite because an identity provider is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateIdentityProvider call.
- `AwsCognitoUserPoolClient` -- AwsCognitoUserPool is a prerequisite because an app client is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateUserPoolClient call.
- `AwsCognitoResourceServer` -- AwsCognitoUserPool is a prerequisite because a resource server is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateResourceServer call.
- `AwsWafWebAcl`
- `AwsWafIpSet`
- `AwsWafRegexPatternSet`
- `AwsCloudwatchLogGroup`
- `AwsCloudwatchAlarm`
- `AwsCloudwatchCompositeAlarm`
- `AwsKinesisStream`
- `AwsKinesisFirehose` -- Every Firehose destination requires an S3 configuration (the primary target for extended_s3; the failed/all-document backup for the rest) and an IAM role Firehose assumes to write to it, so both are hard deploy prerequisites.
- `AwsKinesisStreamConsumer` -- A consumer registers against exactly one stream and cannot exist without it.
- `AwsAthenaWorkgroup`
- `AwsGlueCatalogDatabase`
- `AwsRedshiftCluster`
- `AwsSagemakerDomain` -- AI/ML A domain cannot exist without VPC subnets and a SageMaker execution role (default_user_settings.execution_role_arn is required), so both are hard deploy prerequisites.
- `AwsAppRunnerService` -- A service can run entirely on companion defaults, so the App Runner family's kinds are dependency-free leaves except the VPC connector (which cannot exist without subnets and security groups). A service's companion references (auto scaling / VPC connector / observability / WAF) are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsAppRunnerAutoScalingConfiguration`
- `AwsAppRunnerVpcConnector`
- `AwsAppRunnerObservabilityConfiguration`
- `AwsTransitGatewayVpcAttachment` -- AwsTransitGateway is a prerequisite because an attachment cannot exist without the gateway it attaches to; AwsSubnet because the attachment provisions an ENI into at least one subnet (the VPC arrives transitively through the subnet's own prerequisites).
- `AwsTransitGatewayRouteTable` -- Only the gateway is a hard prerequisite: a route table can exist empty. Associations, propagations, and routes referencing attachments are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsBatchComputeEnvironment` -- A MANAGED compute environment always launches into VPC subnets, so the subnet is a hard deploy prerequisite (security groups are required only for the Fargate types -- scenario-declared, not a registry edge).
- `AwsBatchJobQueue` -- A job queue cannot exist without at least one VALID compute environment to map onto.
- `AwsBatchSchedulingPolicy`
- `AwsBatchJobDefinition`
- `AwsCodeBuildProject` -- CI/CD
- `AwsCodePipeline`
- `AwsMwaaEnvironment` -- Workflow / Orchestration AwsSubnet and AwsSecurityGroup are prerequisites because the environment's network interfaces are placed in referenced private subnets and AWS requires at least one attached security group at creation.
- `AwsNeptuneCluster` -- Graph Database
- `AwsMemorydbCluster` -- A cluster always launches into a subnet group; the subnets are the hard deploy prerequisite. The ACL it attaches is optional composition (the built-in "open-access" ACL needs no resource) -- scenarios declare the ACL/user chain via the e2e-prerequisites annotation.
- `AwsMemorydbUser`
- `AwsMemorydbAcl` -- An empty ACL is valid (MemoryDB has no mandatory "default" member), so the user is optional composition -- the composed scenario declares it via the e2e-prerequisites annotation, never a registry edge.
- `AwsMskCluster` -- Streaming AwsSubnet and AwsSecurityGroup are prerequisites because brokers are placed in referenced subnets and AWS requires at least one attached security group at creation.
- `AwsMskServerlessCluster` -- AwsSubnet is a prerequisite because the serverless cluster's network interfaces are placed in referenced subnets (security groups are optional -- AWS attaches the VPC default group when none are referenced).
- `AwsLambdaEventSourceMapping` -- AwsLambda is a prerequisite because a mapping cannot exist without the function it invokes (a required reference). Event sources (SQS, Kinesis, DynamoDB, MSK) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsSnsSubscription` -- AwsSnsTopic is a prerequisite because a subscription cannot exist without the topic it subscribes to (a required reference). Endpoints (SQS queues, Lambda functions, Firehose streams) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsPlantonRunner` -- AwsSubnet is a prerequisite because the runner appliance places its network interfaces into referenced subnets -- the placement reference must resolve before the appliance can deploy.
- `AwsRoute53HealthCheck`
- `AwsSesConfigurationSet` -- Both SES kinds are dependency-free leaves: an identity's configuration set is optional composition (scenarios declare it via the e2e-prerequisites annotation), and a configuration set's event destinations reference other kinds only optionally.
- `AwsSesEmailIdentity`
- `AwsSecretsManagerSecret` -- A dependency-free leaf: the KMS key, rotation Lambda, and external rotation role references are all optional composition -- scenarios declare them via the e2e-prerequisites annotation, never registry edges.
- `AwsOpenSearchServerlessCollection` -- A dependency-free leaf: the collection-scoped encryption/network/ data-access/retention policies are module-rendered, and the KMS key and data-access principal references are optional composition (e2e-prerequisites annotation).
- `AwsBedrockGuardrail` -- A dependency-free leaf: the KMS key reference is optional composition (e2e-prerequisites annotation); published versions are folded satellites of the guardrail itself.
- `AwsBedrockCustomModel` -- AwsIamRole is a prerequisite because Bedrock assumes the job role to read training data and write outputs; the S3 locations and KMS key are optional composition (e2e-prerequisites annotation).
- `AwsBedrockInferenceProfile` -- A dependency-free leaf: the model source is a foundation model or an AWS system-defined cross-region profile, never a customer resource.
- `AwsBedrockProvisionedThroughput` -- A dependency-free leaf in the registry: capacity is typically bought for an AwsBedrockCustomModel (the default reference), but foundation model ARNs are equally legal, so the edge is optional composition.
- `AwsBedrockModelAccess` -- A dependency-free leaf: the agreement covers an AWS-listed foundation model, never a customer resource.
- `AwsBedrockAgent` -- AwsIamRole is a prerequisite because the Bedrock service assumes the agent resource role to invoke models, action-group Lambdas, and knowledge bases; the guardrail, KMS key, provisioned throughput, and collaborator/knowledge-base edges are optional composition (e2e-prerequisites annotation). Action groups, aliases, collaborators, and knowledge-base associations are folded satellites of the agent.
- `AwsBedrockKnowledgeBase` -- AwsIamRole is a prerequisite because the Bedrock service assumes the knowledge-base role to read data sources, call the embedding model, and read/write the vector store; the vector-store and data-source reference edges (OpenSearch, S3, Secrets Manager, ...) are optional composition (e2e-prerequisites annotation). Data sources are folded satellites of the knowledge base.
- `AwsBedrockFlow` -- AwsIamRole is a prerequisite because the Bedrock service assumes the flow execution role to invoke the models, agents, knowledge bases, and Lambdas its nodes reference; the node-level reference edges are optional composition (e2e-prerequisites annotation).
- `AwsBedrockPrompt` -- A dependency-free leaf: variants target AWS-listed foundation models by ID; targeting another agent's alias is optional composition (e2e-prerequisites annotation).
- `AzureResourceGroup` -- 2000–2999: Azure resources
- `AzureAksCluster` -- AzureResourceGroup is the only required parent: the cluster is created inside a referenced resource group. Subnet is optional on the default node pool (AKS provisions managed networking when unset).
- `AzureAksNodePool` -- AzureAksCluster is a prerequisite because a node pool attaches to an existing cluster by ARM ID; the resource group chains transitively.
- `AzureContainerRegistry` -- AzureResourceGroup is a prerequisite because a container registry is created inside a resource group.
- `AzureDnsZone` -- AzureResourceGroup is a prerequisite because the DNS zone is created inside a referenced resource group that must already exist.
- `AzureKeyVault` -- AzureResourceGroup is a prerequisite because a key vault is created inside a referenced resource group in composed environments.
- `AzureVirtualNetwork` -- AzureResourceGroup is a prerequisite because a virtual network is created inside a referenced resource group in composed environments.
- `AzureNatGateway` -- AzureResourceGroup is a prerequisite because a NAT gateway is created inside a referenced resource group in composed environments.
- `AzureVirtualMachine` -- AzureNetworkInterface is a prerequisite because a virtual machine attaches at least one NIC (the subnet, network, and resource group chain transitively through the NIC's own prerequisites).
- `AzureStorageAccount` -- AzureResourceGroup is a prerequisite because a storage account is created inside a referenced resource group in composed environments.
- `AzureDnsRecord` -- AzureDnsZone is a prerequisite because a record set is created inside a referenced zone (the resource group chains transitively through the zone). Public DNS zone names are not globally unique, so a shared zone fixture is safe to recreate across scenarios.
- `AzureSubnet` -- AzureVirtualNetwork is a prerequisite because a subnet is an ARM child of a referenced network -- the network must exist before the subnet can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzureNetworkSecurityGroup` -- AzureResourceGroup is a prerequisite because a network security group is created inside a referenced resource group in composed environments.
- `AzurePublicIp` -- AzureResourceGroup is a prerequisite because a public IP is created inside a referenced resource group in composed environments.
- `AzurePrivateEndpoint` -- AzureSubnet is a prerequisite because a private endpoint draws its private IP from a referenced subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite). The connection target is polymorphic and the DNS zones / ASGs are optional, so none of those are prerequisites.
- `AzurePrivateDnsZone` -- AzureResourceGroup is a prerequisite because a private DNS zone is created inside a referenced resource group in composed environments.
- `AzureApplicationGateway` -- AzureSubnet is a prerequisite because a gateway cannot exist without its dedicated gateway_ip_configuration subnet (the network and resource group chain transitively through the subnet's own prerequisites); public frontends additionally reference a public IP, but private-only gateways are legal, so it is not a registry prerequisite.
- `AzureLoadBalancer` -- AzureResourceGroup is a prerequisite because a load balancer is created inside a referenced resource group (frontends additionally reference subnets or public IPs, but neither is universally required, so they are not registry prerequisites).
- `AzureRouteTable` -- AzureResourceGroup is a prerequisite because a route table is created inside a referenced resource group in composed environments.
- `AzurePrivateDnsZoneVirtualNetworkLink` -- AzurePrivateDnsZone and AzureVirtualNetwork are prerequisites because a virtual network link is a child resource of a referenced zone and binds it to a referenced network -- both must exist before the link can be written. (The resource group arrives transitively through the zone's and network's own prerequisite declarations.)
- `AzureVirtualNetworkPeering` -- AzureVirtualNetwork is a prerequisite because a peering is an ARM child of its local network and binds it to a remote network -- the local network must exist before the peering can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzurePublicIpPrefix` -- AzureResourceGroup is a prerequisite because a public IP prefix is created inside a referenced resource group in composed environments.
- `AzureNetworkInterface` -- AzureSubnet is a prerequisite because a network interface's IP configurations deploy into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureManagedDisk` -- AzureResourceGroup is a prerequisite because a managed disk is created inside a resource group.
- `AzureVirtualMachineScaleSet` -- AzureSubnet is a prerequisite because every scale-set instance's network interface deploys into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureKeyVaultKey` -- AzureKeyVault is a prerequisite because a key is a data-plane object inside a referenced vault -- the vault must exist before the key can be written (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultCertificate` -- AzureKeyVault is a prerequisite because a certificate is a data-plane object inside a referenced vault -- the vault must exist before the certificate can be enrolled or imported (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultSecret` -- AzureKeyVault is a prerequisite because a secret is a data-plane object inside a referenced vault -- the vault must exist before the secret can be written (the resource group chains transitively through the vault's own prerequisite). Part of the Key Vault family (2005, 2025-2026) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureWebApplicationFirewallPolicy` -- AzureResourceGroup is a prerequisite because a WAF policy is created inside a referenced resource group; the Application Gateways that attach the policy reference it, never the reverse.
- `AzureApplicationSecurityGroup` -- AzureResourceGroup is a prerequisite because an application security group is created inside a referenced resource group; network interfaces, scale-set IP configurations, and NSG security rules reference the group, never the reverse.
- `AzureDiskEncryptionSet` -- AzureKeyVaultKey is a prerequisite because a disk encryption set wraps customer data with a referenced key -- the key (and its vault, which chains transitively) must exist before the set can resolve the key URL at create time.
- `AzurePostgresqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureRedisCache` -- AzureResourceGroup is a prerequisite because the cache is created inside a referenced resource group (VNet injection additionally references a dedicated subnet, but only the Premium tier supports it, so it is not a registry prerequisite).
- `AzureCosmosdbAccount` -- AzureResourceGroup is a prerequisite because the account is created inside a referenced resource group.
- `AzureMssqlServer` -- AzureResourceGroup is a prerequisite because the logical server is created inside a referenced resource group.
- `AzureMysqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureMssqlDatabase` -- The parent logical server is referenced via server_id, not auto-deployed: E2E scenarios declare their own server fixture (minimal-server.yaml or the pool-attach chain through AzureMssqlElasticPool) so sequential subtests never destroy and recreate the same globally unique server_name.
- `AzureMssqlElasticPool` -- AzureMssqlServer is a prerequisite because every elastic pool lives on a referenced logical server (the server's resource group is transitive).
- `AzureRedisLinkedServer` -- The target and linked caches are referenced via ARM ids, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicy` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicyAssignment` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureContainerAppEnvironment` -- AzureResourceGroup is a prerequisite because the environment is created inside a referenced resource group that must already exist.
- `AzureContainerApp` -- AzureContainerAppEnvironment is a prerequisite because every app runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureServicePlan` -- AzureResourceGroup is a prerequisite because the plan is created inside a referenced resource group that must already exist.
- `AzureFunctionApp` -- AzureServicePlan is a prerequisite because a function app runs on a referenced plan (the resource group arrives transitively through the plan). The required storage account is deliberately NOT a registry prerequisite: storage-account names are globally unique, so scenarios bring their own scenario-local account fixtures.
- `AzureLinuxWebApp` -- AzureServicePlan is a prerequisite because a web app runs on a referenced plan (the resource group arrives transitively through the plan).
- `AzureContainerAppJob` -- AzureContainerAppEnvironment is a prerequisite because a job runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureContainerAppEnvironmentStorage` -- AzureContainerAppEnvironment is a prerequisite because the storage registration lives on a referenced environment. The Azure Files share and storage account are deliberately NOT registry prerequisites: storage-account names are globally unique, so scenarios bring their own scenario-local account + share fixtures.
- `AzureContainerAppEnvironmentDaprComponent` -- AzureContainerAppEnvironment is a prerequisite because the Dapr component is registered on a referenced environment.
- `AzureContainerAppEnvironmentCertificate` -- AzureContainerAppEnvironment is a prerequisite because the certificate is stored on a referenced environment.
- `AzureContainerAppEnvironmentManagedCertificate` -- AzureContainerAppEnvironment is a prerequisite because the managed certificate is provisioned on a referenced environment.
- `AzureLogAnalyticsWorkspace` -- AzureResourceGroup is a prerequisite because the workspace is created inside a referenced resource group that must already exist.
- `AzureApplicationInsights` -- AzureLogAnalyticsWorkspace is a prerequisite because workspace-based Application Insights stores its telemetry in a referenced workspace (the resource group chains transitively through the workspace).
- `AzureMonitorDiagnosticSetting` -- AzureLogAnalyticsWorkspace is a prerequisite because the setting's scenarios route a fixture workspace's telemetry (the workspace doubles as target and destination); the target itself is polymorphic.
- `AzureMonitorActionGroup` -- AzureResourceGroup is a prerequisite because the action group is created inside a referenced resource group that must already exist.
- `AzureMonitorMetricAlert` -- AzureMonitorActionGroup is a prerequisite because a metric alert's actions fire into a referenced action group (the resource group chains transitively); alert scopes are polymorphic.
- `AzureMonitorScheduledQueryAlert` -- AzureLogAnalyticsWorkspace is a prerequisite because the rule queries a referenced workspace scope; AzureMonitorActionGroup because its action fires into a referenced action group.
- `AzureMonitorActivityLogAlert` -- AzureMonitorActionGroup is a prerequisite because an activity log alert's actions fire into a referenced action group (the resource group chains transitively). The alert itself is subscription-global and its scopes are polymorphic.
- `AzureApplicationInsightsStandardWebTest` -- AzureApplicationInsights is a prerequisite because a standard web test binds to a referenced Application Insights component (the resource group chains transitively through the component).
- `AzureUserAssignedIdentity` -- AzureResourceGroup is a prerequisite because the identity is created inside a referenced resource group that must already exist.
- `AzureRoleAssignment` -- AzureResourceGroup and AzureUserAssignedIdentity are prerequisites because an assignment grants a role at a referenced scope (most commonly a resource group) to a referenced principal (most commonly a managed identity) -- both must exist before the grant can be written.
- `AzureRoleDefinition` -- AzureResourceGroup is a prerequisite because a custom role definition is created at a referenced scope, most commonly a resource group in composed environments -- the scope must exist before the definition can be written.
- `AzureFederatedIdentityCredential` -- AzureUserAssignedIdentity is the prerequisite because a federated identity credential is a child resource of a referenced managed identity -- the identity must exist before the credential can be written on it. (The resource group arrives transitively through the identity's own prerequisite declaration.)
- `AzureServiceBusNamespace` -- AzureResourceGroup is a prerequisite because a Service Bus namespace is created inside a referenced resource group in composed environments. The namespace is the container every Service Bus messaging entity (queue, topic, subscription, authorization rule, geo-DR pairing) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureEventHubNamespace` -- AzureResourceGroup is a prerequisite because an Event Hub namespace is created inside a referenced resource group in composed environments. The namespace is the container every Event Hubs entity (event hub, consumer group, authorization rule, schema group, geo-DR pairing, customer-managed key) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureServiceBusQueue`
- `AzureServiceBusTopic`
- `AzureServiceBusSubscription`
- `AzureServiceBusAuthorizationRule`
- `AzureServiceBusDisasterRecoveryConfig`
- `AzureEventHub`
- `AzureEventHubConsumerGroup`
- `AzureEventHubAuthorizationRule`
- `AzureFrontDoorProfile` -- AzureResourceGroup is a prerequisite because a Front Door profile is created inside a referenced resource group in composed environments. The profile is the container every Front Door delivery resource (endpoint, origin group, origin, route) nests under.
- `AzureFrontDoorEndpoint` -- AzureFrontDoorProfile is a prerequisite because an endpoint is an ARM child of a referenced profile -- the profile must exist before the endpoint can be written. (The resource group arrives transitively through the profile's own prerequisite declaration.)
- `AzureFrontDoorOriginGroup` -- AzureFrontDoorProfile is a prerequisite because an origin group is an ARM child of a referenced profile.
- `AzureFrontDoorOrigin` -- AzureFrontDoorOriginGroup is a prerequisite because an origin is an ARM child of a referenced origin group (the profile and resource group chain transitively).
- `AzureFrontDoorRoute` -- A route attaches to an endpoint (its ARM parent) and forwards to an origin group whose origins must exist before ARM accepts the route -- so both the endpoint and the origin chain are genuine deploy-order prerequisites.
- `AzureFrontDoorRuleSet` -- AzureFrontDoorProfile is a prerequisite because a rule set is an ARM child of a referenced profile. The rules live inside the set (they form one ordered policy document); routes attach the set by ARM ID.
- `AzureFrontDoorCustomDomain` -- AzureFrontDoorProfile is a prerequisite because a custom domain is an ARM child of a referenced profile. The DNS zone and (for bring-your-own certificates) the Front Door secret are optional references, not deploy-order prerequisites.
- `AzureFrontDoorSecret` -- AzureFrontDoorSecret is a prerequisite-light kind: only the profile (its ARM parent) must exist. The Key Vault certificate it wraps is a reference resolved before the module runs; its vault chain is exercised through scenario-local fixtures in E2E.
- `AzureFrontDoorFirewallPolicy` -- AzureResourceGroup is a prerequisite because the Front Door WAF policy is created inside a referenced resource group -- it is a GLOBAL resource, not a profile child (a different ARM type than the regional Application Gateway WAF policy). Security policies attach it to profiles; the policy itself depends on nothing else.
- `AzureFrontDoorSecurityPolicy` -- A security policy is an ARM child of a profile that associates a referenced WAF policy with referenced domains -- so the endpoint (the default-domain association target; the profile arrives transitively through it) and the WAF policy are genuine deploy-order prerequisites.
- `AzureStorageContainer` -- None of the storage data-service kinds declares a registry prerequisite on AzureStorageAccount: account names are GLOBALLY unique and Azure holds a just-deleted name, so a recreate-per-scenario fixture would hang -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the storage_account_id reference itself.
- `AzureStorageShare`
- `AzureStorageQueue`
- `AzureStorageTable`
- `AzureStorageEncryptionScope`
- `AzureStorageDataLakeGen2Filesystem`
- `AzureStorageLocalUser`
- `AzureStorageObjectReplication`
- `AzureCosmosdbSqlDatabase` -- None of the Cosmos DB data-service kinds declares a registry prerequisite on AzureCosmosdbAccount: account names are GLOBALLY unique DNS labels, so a recreate-per-scenario fixture would risk name-reuse hangs -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the cosmosdb_account_id / parent-database references themselves.
- `AzureCosmosdbSqlContainer`
- `AzureCosmosdbMongoDatabase`
- `AzureCosmosdbMongoCollection`
- `AzureCosmosdbSqlRoleDefinition`
- `AzureCosmosdbSqlRoleAssignment`
- `AzureManagedRedis` -- AzureResourceGroup is the cluster's only registry prerequisite: the cluster is created inside a referenced resource group. The geo-replication and access-policy-assignment children declare NO prerequisite on AzureManagedRedis: clusters are expensive, slow-provisioning parents, so their E2E scenarios declare scenario-local cluster fixtures instead of recreating a shared one per scenario. Deploy ordering in composed environments still flows from the managed_redis_id references themselves.
- `AzureManagedRedisGeoReplication`
- `AzureManagedRedisAccessPolicyAssignment`
- `AzureEventHubDisasterRecoveryConfig`
- `AzureEventHubSchemaGroup`
- `AzureEventHubCluster` -- AzureResourceGroup is a prerequisite because a dedicated Event Hubs cluster is created inside a referenced resource group in composed environments. Note: clusters cannot be deleted for 4 hours after creation (Azure's moratorium), so E2E treats this kind as offline-gated.
- `AzureEventHubNamespaceCustomerManagedKey`
- `AzureMssqlFailoverGroup` -- AzureMssqlServer is a prerequisite because a failover group is created on a referenced primary logical server and points at a partner server; the primary (and its resource group, which chains transitively) must exist before the group can be written.
- `AzureContainerAppCustomDomain` -- AzureContainerApp is a prerequisite because the domain binding lives in a referenced app's ingress configuration (the environment and resource group chain transitively through the app).
- `AzureFirewallPolicy`
- `AzureFirewallPolicyRuleCollectionGroup` -- AzureFirewallPolicy is a prerequisite because a rule collection group is a child document of a referenced policy (the resource group chains transitively through the policy).
- `AzureFirewall` -- AzureSubnet is a prerequisite because a VNet-deployed firewall's data path lives in a dedicated subnet that must be named exactly "AzureFirewallSubnet" (the virtual network and resource group chain transitively through the subnet). The E2E install profile publishes a fixture subnet with that exact name and a /26 prefix.
- `AzureIpGroup`
- `AzureVirtualNetworkGateway` -- AzureSubnet is a prerequisite because every virtual network gateway lives in a dedicated subnet named exactly "GatewaySubnet" (the virtual network and resource group chain transitively through the subnet); the subnet install profile publishes a fixture instance with that exact ARM name. AzurePublicIp is a prerequisite because a VPN-type gateway (the default shape) requires a public IP per ip configuration; the address install profile publishes a dedicated zone-redundant instance (a gateway binds its address exclusively, and the AZ gateway SKUs require zones on it).
- `AzureVirtualNetworkGatewayConnection` -- Both gateways are prerequisites: a connection joins a virtual network gateway to a far side, and the site-to-site far side is a local network gateway (the GatewaySubnet, VNet, and resource group chain transitively through the virtual network gateway).
- `AzureLocalNetworkGateway`
- `AzurePrivateLinkService` -- AzureSubnet is the sole prerequisite: every NAT ip configuration draws its address from a subnet with private-link-service network policies disabled (the subnet install profile publishes a fixture instance with that flag off). The Standard load balancer whose frontend the service typically fronts is NOT a registry prerequisite -- the spec's destination is an exactly-one-of (load balancer frontend OR fixed destination IP), so scenarios that use the load-balancer shape declare it via the planton.dev/e2e-prerequisites annotation instead.
- `AzureExpressRouteCircuit`
- `AzureExpressRouteCircuitPeering` -- The circuit is the prerequisite: a peering is an ARM child of the circuit, addressed by the circuit's name (the resource group chains transitively through the circuit).
- `AzureExpressRouteGateway` -- The hub is the prerequisite: ARM requires an ExpressRoute Gateway to be deployed INTO a Virtual WAN hub (the WAN and resource group chain transitively through the hub).
- `AzureExpressRoutePort` -- ExpressRoute Port: your own physical port pair on a Microsoft edge router (ExpressRoute Direct), from whose bandwidth circuits are carved. Self-contained -- only the resource group is required.
- `AzureVirtualWan` -- Virtual WAN: the umbrella of Azure's managed hub-and-spoke networking, under which virtual hubs and their gateways are created. Self-contained -- only the resource group is required.
- `AzureVirtualHub` -- The WAN is the prerequisite: this kind models the Virtual WAN hub (virtual_wan_id is required; standalone hubs are the legacy Route Server construction, which has its own ARM surface). The resource group chains transitively through the WAN.
- `AzureVirtualHubConnection` -- Both sides of the attachment are prerequisites: the hub being joined and the spoke virtual network being attached.
- `AzureVpnGateway` -- The hub is the prerequisite: ARM deploys a Virtual WAN VPN gateway INTO a virtual hub (virtual_hub_id is required and immutable; the WAN and resource group chain transitively through the hub). ARM allows one VPN gateway per hub.
- `AzureVpnGatewayConnection` -- Both ends of the tunnel are prerequisites: a connection is an ARM child of the VPN gateway and pins each of its links to a specific link of the remote VPN site (the hub, WAN, and resource group chain transitively through the gateway).
- `AzureVpnSite` -- The WAN is the prerequisite: a VPN site is the Virtual WAN world's address-book entry for one branch location (virtual_wan_id is required; the classic-world sibling without a WAN is AzureLocalNetworkGateway). The resource group chains transitively through the WAN.
- `AzurePointToSiteVpnGateway` -- The hub and the server configuration are both prerequisites: a point-to-site VPN gateway deploys INTO a virtual hub (one P2S gateway per hub, a slot separate from the hub's site-to-site VPN gateway) and is born pointing at the VPN server configuration that defines how its users authenticate -- both ARM-required and fixed at creation. The WAN and resource group chain transitively through the hub.
- `AzureVpnServerConfiguration` -- Self-contained -- only the resource group is required: a VPN server configuration is the reusable "who may connect and how" authentication policy (Entra ID / certificate / RADIUS) that point-to-site VPN gateways attach to; it references no other Azure resource.
- `AzureCognitiveAccount` -- Self-contained -- only the resource group is required: an Azure AI services account (Azure OpenAI, the multi-service AIServices account, the single-service accounts) needs no other Azure resource; subnets (network rules), Key Vault keys (CMK), storage accounts and user-assigned identities are optional references.
- `AzureCognitiveDeployment` -- An ARM child of its account: a model deployment (which model runs, at which throughput class) exists only on an Azure AI services account of kind "OpenAI" or "AIServices".
- `AzureCognitiveAccountProject` -- An ARM child of its account: an AI Foundry project exists only on an "AIServices"-kind account with project management enabled.
- `AzureMachineLearningWorkspace` -- The workspace REQUIRES all three companion services at creation (default storage, secrets vault, telemetry) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureMachineLearningDatastore` -- An ARM child of its workspace. The storage target (container, filesystem or share) is scenario-declared via the e2e-prerequisites annotation -- only the blob scenario needs a container, so it is not a kind-wide prerequisite.
- `AzureMachineLearningComputeCluster` -- An ARM child of its workspace (.../computes/{name}) -- the auto-scaling pool of VMs training jobs run on.
- `AzureMachineLearningComputeInstance` -- An ARM child of its workspace (.../computes/{name}) -- a single always-on VM serving as one data scientist's cloud workstation.
- `AzureAiFoundry` -- The hub REQUIRES both companion services at creation (secrets vault, default storage) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureAiFoundryProject` -- Deploys into its hub's resource group (the provider derives the group from the hub reference -- the project spec carries none).
- `AzureSearchService`
- `AzureMachineLearningOnlineEndpoint` -- An ARM child of its workspace (.../onlineEndpoints/{name}) -- the stable scoring address applications call. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningOnlineDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the running copy of a model the endpoint's traffic map routes to.
- `AzureMachineLearningBatchEndpoint` -- An ARM child of its workspace (.../batchEndpoints/{name}) -- the stable address batch scoring jobs are submitted to. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningBatchDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the job recipe (model, compute, batching behavior) the endpoint's default-deployment pointer routes submissions to.
- `AzureRecoveryServicesVault` -- The Recovery Services vault (Microsoft.RecoveryServices/vaults) -- the safe that classic Azure Backup data and Site Recovery configuration live in. Backup policies and protected items are ARM children of a vault.
- `AzureBackupPolicyVm` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern IaaS VM backups.
- `AzureBackupProtectedVm` -- An ARM child of its vault (.../protectedItems/...) -- the binding that puts one virtual machine under a backup policy's protection.
- `AzureBackupPolicyFileShare` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern Azure Files share backups (snapshot or vaulted).
- `AzureBackupProtectedFileShare` -- An ARM child of its vault (.../protectedItems/AzureFileShare;...) -- the binding that puts one Azure Files share under a backup policy's protection. The share's storage account must already be registered with the vault (AzureBackupContainerStorageAccount).
- `AzureDataProtectionBackupVault` -- The Data Protection backup vault (Microsoft.DataProtection/ backupVaults) -- the safe that MODERN Azure Backup data lives in (managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, Data Lake storage). Backup policies and backup instances are ARM children of a vault.
- `AzureDataProtectionBackupPolicy` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules for ONE Data Protection datasource type (blob storage, disk, Kubernetes cluster, MySQL/PostgreSQL flexible server, or Data Lake storage), modeled as one kind with variant blocks.
- `AzureDataProtectionBackupInstance` -- An ARM child of its vault (.../backupInstances/{name}) -- the binding that puts ONE datasource (a managed disk, a storage account's blob services, an AKS cluster, a MySQL/PostgreSQL flexible server, or a Data Lake storage account) under a Data Protection backup policy, modeled as one kind with variant blocks. The vault's managed identity must hold the datasource roles Azure Backup requires BEFORE the instance is created.
- `AzureBastionHost` -- AzureSubnet and AzurePublicIp are prerequisites because a dedicated-infrastructure Bastion host (Basic/Standard/Premium -- the default shapes) deploys into a subnet named exactly "AzureBastionSubnet" and binds a Standard static public IP EXCLUSIVELY (the virtual network and resource group chain transitively through the subnet). The Developer SKU instead attaches to a virtual network directly and uses neither.
- `AzureNetworkWatcherFlowLog` -- AzureVirtualNetwork and AzureStorageAccount are prerequisites because a flow log records a network-scoped target (a virtual network in the common case; subnets and network interfaces chain through the network) into a referenced storage account. The regional Network Watcher parent is NOT a prerequisite: Azure auto-creates it ("NetworkWatcher_{region}" in "NetworkWatcherRG") the moment the region hosts a virtual network, and the flow log references it by name. Traffic Analytics' Log Analytics workspace is an optional arm, declared by scenarios that use it.
- `AzurePrivateDnsResolver` -- AzureVirtualNetwork and AzureSubnet are prerequisites because a DNS Private Resolver anchors to a referenced virtual network (at most ONE resolver per network -- Azure enforces it) and each of its inbound/outbound endpoints occupies its own dedicated subnet delegated to "Microsoft.Network/dnsResolvers" (the resource group chains transitively through the network and subnets).
- `AzurePrivateDnsResolverForwardingRuleset` -- AzurePrivateDnsResolver is a prerequisite because a DNS forwarding ruleset steers a resolver's OUTBOUND endpoints -- it binds their ARM ids (at most 2, same resolver) at creation. (The resource group and network chain transitively through the resolver's own prerequisite declarations.)
- `AzureBackupContainerStorageAccount` -- Registers a storage account with a Recovery Services vault as a backup container (.../protectionContainers/StorageContainer;...) -- one registration per storage-account-and-vault pair, required BEFORE any of the account's file shares can be protected. Part of the backup family (2175-2179) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureDataProtectionResourceGuard` -- The Data Protection Resource Guard (Microsoft.DataProtection/ resourceGuards) -- the approval gate behind Multi-User Authorization: privileged vault operations (disabling soft delete, reducing retention) require an approval through a guard, which typically lives in a DIFFERENT administrator's scope. Vaults reference a guard by its ARM ID. Part of the backup family (2175-2182) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzurePrivateDnsResolverVirtualNetworkLink` -- The attachment that makes a DNS forwarding ruleset take effect in one virtual network ({ruleset_id}/virtualNetworkLinks/{name}) -- one link per ruleset-network pair, up to 500 per ruleset, spokes joining and leaving independently (which is why the link is a standalone kind, exactly like AzurePrivateDnsZoneVirtualNetworkLink). Part of the DNS Private Resolver family (2186-2187) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `GcpArtifactRegistryRepo` -- 3000–3999: GCP resources
- `GcpTargetHttpsProxy` -- The URL map is the parent a proxy cannot exist without; the classic compute certificate kinds and the SSL policy are the fixture parents the committed scenarios attach. The Certificate Manager certificate list (certificate_manager_certificates, honored only by the cross-region internal ALB) is optional composition -- a scenario that arms it declares GcpCertManagerCert via the e2e-prerequisites annotation, never a registry edge that would tax every proxy and forwarding-rule chain.
- `GcpCloudFunction`
- `GcpCloudRun`
- `GcpCloudSql`
- `GcpDnsZone`
- `GcpGcsBucket`
- `GcpGkeCluster`
- `GcpIamCustomRole`
- `GcpProject`
- `GcpVpcNetwork`
- `GcpSubnetwork`
- `GcpRouterNat`
- `GcpGkeNodePool`
- `GcpServiceAccount`
- `GcpGkeWorkloadIdentityBinding`
- `GcpCertManagerCert`
- `GcpComputeInstance`
- `GcpDnsRecord`
- `GcpProjectIamMember`
- `GcpFirewallRule`
- `GcpGlobalAddress`
- `GcpCloudArmorPolicy`
- `GcpHealthCheck`
- `GcpBackendBucket`
- `GcpBackendService`
- `GcpRegionNetworkEndpointGroup`
- `GcpUrlMap`
- `GcpManagedSslCertificate`
- `GcpTargetHttpProxy`
- `GcpAlloydbCluster`
- `GcpRedisInstance`
- `GcpFirestoreDatabase`
- `GcpSpannerInstance`
- `GcpSpannerDatabase`
- `GcpBigtableInstance`
- `GcpMemorystoreInstance`
- `GcpCloudSqlDatabase`
- `GcpCloudSqlUser`
- `GcpAlloydbInstance`
- `GcpAlloydbUser`
- `GcpSpannerBackupSchedule`
- `GcpBigtableTable`
- `GcpFirestoreBackupSchedule`
- `GcpFirestoreIndex`
- `GcpBigQueryDataset`
- `GcpDataprocCluster`
- `GcpDataprocAutoscalingPolicy`
- `GcpBigQueryTable`
- `GcpPubSubTopic`
- `GcpPubSubSubscription`
- `GcpCloudTasksQueue`
- `GcpCloudSchedulerJob`
- `GcpPubSubSchema`
- `GcpVertexAiNotebook`
- `GcpVertexAiEndpoint`
- `GcpVertexAiIndex`
- `GcpVertexAiIndexEndpoint` -- Vector Search IndexEndpoint — distinct from the online-prediction GcpVertexAiEndpoint (671); different GCP resources, different kinds.
- `GcpVertexAiDeployedIndex`
- `GcpCloudComposerEnvironment`
- `GcpCloudComposerUserWorkloadsSecret`
- `GcpCloudComposerUserWorkloadsConfigMap`
- `GcpKmsKeyRing`
- `GcpKmsKey`
- `GcpKmsKeyIamMember`
- `GcpFilestoreInstance`
- `GcpWorkloadIdentityPool` -- 3101–3109: IAM/identity family (overflow block; the 3000–3022 foundation/security sub-band is fully allocated)
- `GcpWorkloadIdentityPoolProvider`
- `GcpServiceAccountIamMember`
- `GcpGlobalForwardingRule` -- 3110–3119: networking/load-balancer family (overflow block; the 3023–3029 LB sub-band is fully allocated)
- `GcpSslPolicy`
- `GcpSslCertificate`
- `GcpServiceNetworkingConnection`
- `GcpAddress`
- `GcpServiceConnectionPolicy`
- `GcpCertManagerDnsAuthorization`
- `GcpCertificateMap` -- GcpCertManagerCert is a prerequisite because a map entry binds hostnames to EXISTING certificates — the canonical map references a certificate fixture's resource name.
- `GcpCloudRunJob` -- 3120–3129: GCP serverless overflow
- `GcpServerlessVpcConnector`
- `GcpComputeDisk` -- 3130–3139: GCP compute overflow (the 3000–3022 foundation sub-band that holds GcpComputeInstance is fully allocated)
- `GcpComputeMig` -- GcpVpcNetwork is a prerequisite because the canonical group runs its fleet on a dedicated custom-mode VPC — a managed instance group's template must attach every VM to a network, and the default VPC is never assumed.
- `GcpMonitoringNotificationChannel` -- 3140–3149: GCP observability & log routing
- `GcpMonitoringAlertPolicy` -- GcpMonitoringNotificationChannel is a prerequisite because the policy's canonical shape references a channel to notify — a policy without a delivery endpoint measures but never pages.
- `GcpMonitoringUptimeCheck`
- `GcpLoggingSink` -- GcpGcsBucket is a prerequisite because the canonical sink exports to a Cloud Storage bucket — the cheapest destination that proves the whole writer-identity grant flow.
- `GcpMonitoringDashboard`
- `GcpMonitoringSlo`
- `GcpLogBucket`
- `GcpLogMetric`
- `GcpSecretManagerSecret` -- 3150–3159: GCP security & identity GcpServiceAccount is a prerequisite because the canonical secret grants secretAccessor to a workload service account — the access story the kind exists to model.
- `GcpIdentityPlatformConfig`
- `GcpIdentityPlatformTenant` -- GcpIdentityPlatformConfig is a prerequisite because tenants exist only in projects whose Identity Platform config enables multi_tenant.allow_tenants — a tenant without the initialized, tenant-enabled project config cannot be created at all.
- `GcpIamOauthClient`
- `GcpIamDenyPolicy`
- `GcpCloudRunDomainMapping` -- 3160–3169: GCP serverless edge GcpCloudRun is a prerequisite because a domain mapping exists only to point a verified domain at a running Cloud Run service — the route it maps must already exist for the mapping to be created at all.
- `GcpWorkflow`
- `GcpEventarcTrigger` -- GcpCloudRun is a prerequisite because the canonical trigger routes a Pub/Sub messagePublished event to a Cloud Run service — the destination story the kind exists to model.
- `GcpEventarcMessageBus`
- `KubernetesNamespace` -- 4000–4999: Kubernetes resources, organized in family sub-bands (4030–4069 also hosts CNI/autoscaling/DR addons; 4130–4149 hosts analytics & ML; 4190–4199 reserved for growth) 4000–4029: Kubernetes building blocks (core API primitives)
- `KubernetesDeployment`
- `KubernetesStatefulSet`
- `KubernetesDaemonSet`
- `KubernetesJob`
- `KubernetesCronJob`
- `KubernetesService`
- `KubernetesSecret`
- `KubernetesManifest`
- `KubernetesHelmRelease`
- `KubernetesConfigMap`
- `KubernetesServiceAccount`
- `KubernetesRbac` -- Bundles the RBAC grant grain (Role/ClusterRole + its binding) into one component: "grant these permissions to these subjects in this scope".
- `KubernetesIngress`
- `KubernetesNetworkPolicy`
- `KubernetesPersistentVolumeClaim`
- `KubernetesStorageClass`
- `KubernetesResourceQuota` -- Manages the namespace-governance pair: the ResourceQuota plus an optional companion LimitRange (per-object defaults/bounds) — two API objects, one governance story.
- `KubernetesPriorityClass`
- `KubernetesPodDisruptionBudget`
- `KubernetesHorizontalPodAutoscaler`
- `KubernetesCertManager` -- 4030–4069: Kubernetes foundation addons (certs, DNS, secrets, ingress, Gateway API, mesh, CNI/autoscaling/DR)
- `KubernetesClusterIssuer` -- KubernetesCertManager is a prerequisite for the three cert-manager CR kinds below: ClusterIssuer/Issuer/Certificate are cert-manager custom resources — without the controller and its CRDs they cannot be applied.
- `KubernetesIssuer`
- `KubernetesCertificate`
- `KubernetesExternalDns`
- `KubernetesExternalSecretsOperator`
- `KubernetesClusterSecretStore` -- KubernetesExternalSecretsOperator is a prerequisite for the three external-secrets CR kinds below: ClusterSecretStore/SecretStore/ ExternalSecret are external-secrets custom resources — without the operator and its CRDs they cannot be applied.
- `KubernetesSecretStore`
- `KubernetesExternalSecret`
- `KubernetesIngressNginx`
- `KubernetesGatewayApiCrds`
- `KubernetesGatewayClass`
- `KubernetesGateway`
- `KubernetesListenerSet`
- `KubernetesHttpRoute`
- `KubernetesGrpcRoute`
- `KubernetesTcpRoute`
- `KubernetesUdpRoute`
- `KubernetesTlsRoute`
- `KubernetesReferenceGrant`
- `KubernetesBackendTlsPolicy`
- `KubernetesIstioBaseCrds`
- `KubernetesIstio`
- `KubernetesDestinationRule` -- Istio API components (mesh traffic policy, security, telemetry). The seven typed resources below (4053–4059) require the Istio CRDs on the cluster, provided by the lightweight CRDs-only KubernetesIstioBaseCrds (851) — NOT the full mesh KubernetesIstio (852).
- `KubernetesServiceEntry`
- `KubernetesPeerAuthentication`
- `KubernetesRequestAuthentication`
- `KubernetesAuthorizationPolicy`
- `KubernetesTelemetry`
- `KubernetesEnvoyFilter`
- `KubernetesMetricsServer`
- `KubernetesCilium`
- `KubernetesKeda`
- `KubernetesKarpenter`
- `KubernetesKarpenterNodePool`
- `KubernetesKarpenterEc2NodeClass`
- `KubernetesClusterAutoscaler`
- `KubernetesVelero`
- `KubernetesKubePrometheusStack` -- 4070–4089: Kubernetes observability
- `KubernetesGrafana`
- `KubernetesSignoz` -- KubernetesClickHouse is a prerequisite because SigNoz stores every trace, metric and log in ClickHouse and deploys none of its own — the telemetry store is composed, never bundled.
- `KubernetesLoki`
- `KubernetesTempo`
- `KubernetesOtelOperator` -- The operator's admission webhooks (failurePolicy Fail) are served with a cert-manager Certificate in the default posture — cert-manager must be running before the operator installs.
- `KubernetesOtelCollector`
- `KubernetesKyverno` -- 4080–4099: Kubernetes security, policy, and identity
- `KubernetesGatekeeper`
- `KubernetesKeycloak` -- Keycloak declarations compose the official Keycloak Operator (which reconciles the Keycloak CR this kind renders) and, on the recommended postgres vendor, a KubernetesPostgres database — both must resolve before the CR can converge.
- `KubernetesOpenBao`
- `KubernetesOpenFga` -- OpenFGA requires a datastore; the recommended arm composes a KubernetesPostgres database (the sandbox memory arm needs nothing, but the registry declares the shape real deployments require).
- `KubernetesKeycloakOperator`
- `KubernetesCloudNativePgOperator` -- 4100–4129: Kubernetes data platforms
- `KubernetesPostgres`
- `KubernetesValkey`
- `KubernetesPerconaMysqlOperator`
- `KubernetesMysql`
- `KubernetesPerconaMongoOperator`
- `KubernetesMongodb`
- `KubernetesStrimziKafkaOperator`
- `KubernetesKafka` -- container_kind: a Strimzi Kafka cluster is a place in the provider's own model — KafkaTopic and KafkaUser declarations BELONG to one cluster (the strimzi.io/cluster label) and are drawn inside its box. Clients that merely talk to the cluster (Connect, MirrorMaker2, UI, Karapace) carry containment_exempt on their bootstrap/trust references.
- `KubernetesKafkaTopic`
- `KubernetesKafkaUser`
- `KubernetesKafkaConnect` -- container_kind: a Connect cluster hosts the connectors deployed INTO it (KafkaConnector's strimzi.io/cluster label names its Connect cluster) — the same room shape as KubernetesKafka above.
- `KubernetesKafkaConnector`
- `KubernetesKafkaMirrorMaker2`
- `KubernetesKarapace`
- `KubernetesKafkaUi`
- `KubernetesOpenSearchOperator`
- `KubernetesOpenSearch`
- `KubernetesAltinityOperator`
- `KubernetesClickHouse`
- `KubernetesSolrOperator`
- `KubernetesSolr`
- `KubernetesNeo4j`
- `KubernetesSeaweedFs`
- `KubernetesQdrant`
- `KubernetesRabbitMqOperator` -- The RabbitMQ Cluster Operator's release manifest ships admission webhooks whose serving certificate is a cert-manager Certificate — cert-manager must be running before the operator installs.
- `KubernetesRabbitMq`
- `KubernetesAirflow` -- 4130–4149: Kubernetes analytics and ML KubernetesPostgres is a prerequisite because Airflow's metadata database composes a KubernetesPostgres by default (the spec's FK defaults resolve onto its outputs) and the migration Job needs the database reachable before the server components start.
- `KubernetesSparkOperator`
- `KubernetesKubeRayOperator`
- `KubernetesRayCluster` -- KubernetesKubeRayOperator is a prerequisite because this kind declares the RayCluster custom resource that only the operator's CRDs admit and only the operator reconciles into head and worker pods.
- `KubernetesFlinkOperator` -- KubernetesCertManager is a prerequisite because the Flink operator's chart, with its default-on admission webhook, renders cert-manager Issuer/Certificate resources and trusts the API server through cert-manager's CA injection — there is no self-signed fallback at the pinned chart, and the webhooks are fail-closed.
- `KubernetesFlinkDeployment` -- KubernetesFlinkOperator is a prerequisite because this kind declares the FlinkDeployment custom resource that only the operator's CRDs admit and only the operator reconciles into a running Flink cluster.
- `KubernetesJupyterHub` -- KubernetesPostgres is a prerequisite because JupyterHub's hub database composes a KubernetesPostgres in its external-database arm (the spec's FK defaults resolve onto its outputs) and the hub pod mounts that database's credential Secret before it can start.
- `KubernetesMlflow` -- KubernetesPostgres is a prerequisite because MLflow's backend store composes a KubernetesPostgres in its production arm (FK defaults onto its outputs; the module composes the connection URI from its credential Secret), and KubernetesSeaweedFs because the artifact store's S3-compatible arm FK-defaults onto the SeaweedFS endpoint and credential Secret.
- `KubernetesTrino` -- KubernetesPostgres is a prerequisite because Trino's postgres catalogs compose a KubernetesPostgres (the catalog host and credential FK-default onto its outputs), and the pods read that database's credential Secret to resolve catalog passwords from environment.
- `KubernetesSuperset` -- KubernetesPostgres is a prerequisite because Superset's REQUIRED metadata database composes a KubernetesPostgres (FK defaults onto its outputs; the module composes the environment Secret from its credential Secret), and KubernetesValkey because the cache/broker arm FK-defaults onto a KubernetesValkey's service and password Secret.
- `KubernetesArgocd` -- 4150–4169: Kubernetes GitOps and CI/CD
- `KubernetesArgoWorkflows`
- `KubernetesTektonOperator`
- `KubernetesTekton` -- KubernetesTektonOperator is a prerequisite because this kind declares the TektonConfig custom resource that only the operator's CRDs admit and only the operator reconciles into running components.
- `KubernetesGhaRunnerScaleSetController`
- `KubernetesGhaRunnerScaleSet` -- KubernetesGhaRunnerScaleSetController is a prerequisite because this kind renders an AutoscalingRunnerSet custom resource that only the controller's CRDs admit and only the controller reconciles into listener and runner pods.
- `KubernetesHarbor`
- `KubernetesJenkins`
- `KubernetesTemporal` -- 4170–4189: Kubernetes app platforms KubernetesPostgres is a prerequisite because the recommended (and E2E-proven) database composition backs Temporal's default and visibility stores with a CloudNativePG cluster.
- `KubernetesNats`
- `KubernetesLocust`
- `DigitalOceanAppPlatformService` -- 5000–5999: DigitalOcean resources
- `DigitalOceanBucket`
- `DigitalOceanContainerRegistry`
- `DigitalOceanDatabaseCluster`
- `DigitalOceanDnsZone`
- `DigitalOceanDroplet`
- `DigitalOceanFirewall`
- `DigitalOceanFunction`
- `DigitalOceanKubernetesCluster`
- `DigitalOceanKubernetesNodePool`
- `DigitalOceanLoadBalancer`
- `DigitalOceanVolume`
- `DigitalOceanVpc`
- `DigitalOceanCertificate`
- `DigitalOceanDnsRecord`
- `CivoBucket` -- 6000–6999: Civo resources
- `CivoCertificate`
- `CivoComputeInstance`
- `CivoDatabase`
- `CivoDnsZone`
- `CivoFirewall`
- `CivoIpAddress`
- `CivoKubernetesCluster`
- `CivoKubernetesNodePool`
- `CivoVolume`
- `CivoVpc`
- `CivoDnsRecord`
- `CloudflareDnsZone` -- 7000–7999: Cloudflare resources
- `CloudflareKvNamespace`
- `CloudflareR2Bucket`
- `CloudflareWorker`
- `CloudflareLoadBalancer`
- `CloudflareD1Database`
- `CloudflareZeroTrustAccessApplication`
- `CloudflareDnsRecord`
- `CloudflareRuleset`
- `CloudflareWorkersKvPair`
- `CloudflareHyperdriveConfig`
- `CloudflareLoadBalancerPool`
- `CloudflareLoadBalancerMonitor`
- `CloudflareZeroTrustAccessPolicy`
- `CloudflareZeroTrustAccessGroup`
- `CloudflareQueue`
- `CloudflarePagesProject`
- `CloudflareZeroTrustTunnel`
- `CloudflareZeroTrustTunnelVirtualNetwork`
- `CloudflareZeroTrustTunnelRoute`
- `CloudflareList`
- `CloudflareListItem`
- `CloudflareTurnstileWidget`
- `CloudflareEmailRoutingZone`
- `CloudflareEmailRoutingRule`
- `CloudflareEmailRoutingAddress`
- `CloudflareOriginCaCertificate`
- `CloudflareCertificatePack`
- `CloudflareCustomHostname`
- `CloudflareCustomHostnameFallbackOrigin`
- `Auth0Connection` -- 8000–8999: Auth0 resources
- `Auth0Client`
- `Auth0EventStream`
- `Auth0ResourceServer`
- `Auth0Action`
- `Auth0Role`
- `OpenFgaStore` -- 9000–9999: OpenFGA resources Note: OpenFGA is Terraform-only - there is no Pulumi provider available. Pulumi modules for OpenFGA resources are pass-through placeholders.
- `OpenFgaAuthorizationModel`
- `OpenFgaRelationshipTuple`
- `OpenStackKeypair` -- 10000–10999: OpenStack resources
- `OpenStackNetwork`
- `OpenStackSubnet`
- `OpenStackRouter`
- `OpenStackRouterInterface`
- `OpenStackSecurityGroup`
- `OpenStackFloatingIp`
- `OpenStackNetworkPort`
- `OpenStackSecurityGroupRule`
- `OpenStackFloatingIpAssociate`
- `OpenStackInstance`
- `OpenStackServerGroup`
- `OpenStackVolume`
- `OpenStackVolumeAttach`
- `OpenStackProject`
- `OpenStackApplicationCredential`
- `OpenStackImage`
- `OpenStackRoleAssignment`
- `OpenStackLoadBalancer`
- `OpenStackLoadBalancerListener`
- `OpenStackLoadBalancerPool`
- `OpenStackLoadBalancerMember`
- `OpenStackLoadBalancerMonitor`
- `OpenStackDnsZone`
- `OpenStackDnsRecord`
- `ScalewayVpc`
- `ScalewayPrivateNetwork`
- `ScalewayPublicGateway`
- `ScalewayLoadBalancer`
- `ScalewayInstanceSecurityGroup`
- `ScalewayInstance`
- `ScalewayKapsuleCluster`
- `ScalewayKapsulePool`
- `ScalewayRdbInstance`
- `ScalewayRedisCluster`
- `ScalewayMongodbInstance`
- `ScalewayObjectBucket`
- `ScalewayBlockVolume`
- `ScalewayContainerRegistry`
- `ScalewayDnsZone`
- `ScalewayDnsRecord`
- `ScalewayServerlessFunction`
- `ScalewayServerlessContainer`
- `AliCloudLogProject`
- `AliCloudRamRole`
- `AliCloudRamPolicy`
- `AliCloudVpc`
- `AliCloudVswitch`
- `AliCloudSecurityGroup`
- `AliCloudEipAddress`
- `AliCloudNatGateway`
- `AliCloudApplicationLoadBalancer`
- `AliCloudNetworkLoadBalancer`
- `AliCloudVpnGateway`
- `AliCloudDnsZone`
- `AliCloudDnsRecord`
- `AliCloudPrivateDnsZone`
- `AliCloudStorageBucket`
- `AliCloudNasFileSystem`
- `AliCloudKmsKey`
- `AliCloudRdsInstance`
- `AliCloudPolardbCluster`
- `AliCloudRedisInstance`
- `AliCloudMongodbInstance`
- `AliCloudEcsInstance`
- `AliCloudContainerRegistry`
- `AliCloudKubernetesCluster`
- `AliCloudKubernetesNodePool`
- `AliCloudCdnDomain`
- `AliCloudFunction`
- `AliCloudSaeApplication`
- `AliCloudRocketmqInstance`
- `AliCloudCenInstance`
- `OciVcn`
- `OciSubnet`
- `OciSecurityGroup`
- `OciCompartment`
- `OciIdentityPolicy`
- `OciDynamicGroup`
- `OciComputeInstance`
- `OciContainerEngineCluster`
- `OciContainerEngineNodePool`
- `OciContainerInstance`
- `OciApplicationLoadBalancer`
- `OciNetworkLoadBalancer`
- `OciDynamicRoutingGateway`
- `OciPublicIp`
- `OciAutonomousDatabase`
- `OciDbSystem`
- `OciMysqlDbSystem`
- `OciPostgresqlDbSystem`
- `OciRedisCluster`
- `OciNosqlTable`
- `OciObjectStorageBucket`
- `OciFileSystem`
- `OciBlockVolume`
- `OciKmsVault`
- `OciKmsKey`
- `OciVaultSecret`
- `OciBastion`
- `OciFunctionsApplication`
- `OciApiGateway`
- `OciStreamPool`
- `OciQueue`
- `OciAlarm`
- `OciLogGroup`
- `OciDnsZone`
- `OciDnsRecord`
- `OciNetworkFirewall`
- `OciDevopsProject`
- `HetznerCloudSshKey`
- `HetznerCloudPlacementGroup`
- `HetznerCloudFirewall`
- `HetznerCloudNetwork`
- `HetznerCloudPrimaryIp`
- `HetznerCloudFloatingIp`
- `HetznerCloudServer`
- `HetznerCloudVolume`
- `HetznerCloudSnapshot`
- `HetznerCloudCertificate`
- `HetznerCloudLoadBalancer`
- `HetznerCloudDnsZone`

### spec.container.sidecars[].env.secrets[].valueFrom.env

`string`

### spec.container.sidecars[].env.secrets[].valueFrom.name

`string` · required

- rule: {"required":true}

### spec.container.sidecars[].env.secrets[].valueFrom.fieldPath

`string`

### spec.container.sidecars[].env.envFrom

`[]EnvFromSource`

Bulk import of environment variables from ConfigMaps or Secrets.

### spec.container.sidecars[].env.envFrom[].prefix

`string`

Optional prefix prepended to each imported key name.
For example, prefix "APP_" with key "PORT" produces env var "APP_PORT".

### spec.container.sidecars[].env.envFrom[].configMapRef

`ConfigMapRef`

Import all keys from a ConfigMap.

### spec.container.sidecars[].env.envFrom[].configMapRef.name

`string` · required

Name of the ConfigMap.

- rule: {"required":true}

### spec.container.sidecars[].env.envFrom[].configMapRef.optional

`bool`

If true, the ConfigMap is allowed to not exist without blocking pod startup.

### spec.container.sidecars[].env.envFrom[].secretRef

`SecretRef`

Import all keys from a Secret.

### spec.container.sidecars[].env.envFrom[].secretRef.name

`string` · required

Name of the Secret.

- rule: {"required":true}

### spec.container.sidecars[].env.envFrom[].secretRef.optional

`bool`

If true, the Secret is allowed to not exist without blocking pod startup.

### spec.container.sidecars[].resources

`ContainerResources`

CPU and memory requests and limits. Requests drive scheduling and are what the
pod is guaranteed; limits are the ceiling enforced at runtime (CPU is throttled,
memory overage is OOM-killed). Omitting limits entirely leaves the container
unbounded — acceptable for batch work on dedicated nodes, risky on shared ones.

### spec.container.sidecars[].resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.container.sidecars[].resources.limits.cpu

`string`

### spec.container.sidecars[].resources.limits.memory

`string`

### spec.container.sidecars[].resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.container.sidecars[].resources.requests.cpu

`string`

### spec.container.sidecars[].resources.requests.memory

`string`

### spec.container.sidecars[].livenessProbe

`Probe`

Liveness probe: restarts the container when it fails. Detects deadlocked or
wedged processes. Keep it strictly about "is the process alive" — checking
downstream dependencies here turns a dependency blip into a restart storm.

### spec.container.sidecars[].livenessProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.container.sidecars[].livenessProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.container.sidecars[].livenessProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.container.sidecars[].livenessProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.container.sidecars[].livenessProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.container.sidecars[].livenessProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.container.sidecars[].livenessProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.sidecars[].livenessProbe.httpGet.portNumber

`int32`

### spec.container.sidecars[].livenessProbe.httpGet.portName

`string`

### spec.container.sidecars[].livenessProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.sidecars[].livenessProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.sidecars[].livenessProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.sidecars[].livenessProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.sidecars[].livenessProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.sidecars[].livenessProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.container.sidecars[].livenessProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.container.sidecars[].livenessProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.container.sidecars[].livenessProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.container.sidecars[].livenessProbe.tcpSocket.portNumber

`int32`

### spec.container.sidecars[].livenessProbe.tcpSocket.portName

`string`

### spec.container.sidecars[].livenessProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.sidecars[].livenessProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.container.sidecars[].livenessProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.sidecars[].readinessProbe

`Probe`

Readiness probe: removes the pod from Service endpoints while it fails. This is
the probe that makes rolling updates zero-downtime — traffic only reaches pods
that report ready.

### spec.container.sidecars[].readinessProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.container.sidecars[].readinessProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.container.sidecars[].readinessProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.container.sidecars[].readinessProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.container.sidecars[].readinessProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.container.sidecars[].readinessProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.container.sidecars[].readinessProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.sidecars[].readinessProbe.httpGet.portNumber

`int32`

### spec.container.sidecars[].readinessProbe.httpGet.portName

`string`

### spec.container.sidecars[].readinessProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.sidecars[].readinessProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.sidecars[].readinessProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.sidecars[].readinessProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.sidecars[].readinessProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.sidecars[].readinessProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.container.sidecars[].readinessProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.container.sidecars[].readinessProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.container.sidecars[].readinessProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.container.sidecars[].readinessProbe.tcpSocket.portNumber

`int32`

### spec.container.sidecars[].readinessProbe.tcpSocket.portName

`string`

### spec.container.sidecars[].readinessProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.sidecars[].readinessProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.container.sidecars[].readinessProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.sidecars[].startupProbe

`Probe`

Startup probe: holds off liveness and readiness checking until the app has
started, so slow-booting applications are not killed mid-initialization. Size
`failure_threshold × period_seconds` to the worst-case startup time.

### spec.container.sidecars[].startupProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.container.sidecars[].startupProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.container.sidecars[].startupProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.container.sidecars[].startupProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.container.sidecars[].startupProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.container.sidecars[].startupProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.container.sidecars[].startupProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.sidecars[].startupProbe.httpGet.portNumber

`int32`

### spec.container.sidecars[].startupProbe.httpGet.portName

`string`

### spec.container.sidecars[].startupProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.sidecars[].startupProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.sidecars[].startupProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.sidecars[].startupProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.sidecars[].startupProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.sidecars[].startupProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.container.sidecars[].startupProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.container.sidecars[].startupProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.container.sidecars[].startupProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.container.sidecars[].startupProbe.tcpSocket.portNumber

`int32`

### spec.container.sidecars[].startupProbe.tcpSocket.portName

`string`

### spec.container.sidecars[].startupProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.sidecars[].startupProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.container.sidecars[].startupProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.sidecars[].volumeMounts

`[]VolumeMount`

Volume mounts for this container. Each entry both declares the mount path and
carries its volume source (ConfigMap, Secret, HostPath, EmptyDir, or PVC); the
module derives the pod-level volume list from the union of all containers'
mounts, de-duplicating by name — so two containers sharing an EmptyDir simply
declare the same mount name and source.

### spec.container.sidecars[].volumeMounts[].name

`string` · required

Name of the volume mount. Must be unique within the container.
Used to correlate with the volume definition.

- rule: {"required":true}

### spec.container.sidecars[].volumeMounts[].mountPath

`string` · required

Path within the container at which the volume should be mounted.
Must be an absolute path.

- rule: {"required":true}

### spec.container.sidecars[].volumeMounts[].readOnly

`bool`

Whether the volume should be mounted read-only.
Default is false.

### spec.container.sidecars[].volumeMounts[].subPath

`string`

Path within the volume from which the container's volume should be mounted.
Defaults to "" (volume's root).
Useful for mounting a subdirectory of a volume.

### spec.container.sidecars[].volumeMounts[].configMap

`ConfigMapVolumeSource`

ConfigMap volume source.
Use this to mount a ConfigMap as a file or directory.

### spec.container.sidecars[].volumeMounts[].configMap.name

`string` · required

Name of the ConfigMap to mount.
Can reference a ConfigMap defined in spec.config_maps or an existing one in the namespace.

- rule: {"required":true}

### spec.container.sidecars[].volumeMounts[].configMap.key

`string`

Specific key from the ConfigMap to mount as a single file.
If not specified, all keys are mounted as files in the directory.

### spec.container.sidecars[].volumeMounts[].configMap.path

`string`

If key is specified, this is the filename to use for the mounted file.
Defaults to the key name if not specified.
Example: key="config" path="app.yaml" mounts the "config" key as "app.yaml"

### spec.container.sidecars[].volumeMounts[].configMap.defaultMode

`int32`

Mode bits to use on created files. Must be a value between 0 and 0777.
Defaults to 0644.
Use 0755 (493 in decimal) for executable scripts.

### spec.container.sidecars[].volumeMounts[].secret

`SecretVolumeSource`

Secret volume source.
Use this to mount a Secret as a file or directory.

### spec.container.sidecars[].volumeMounts[].secret.name

`string` · required

Name of the Secret to mount.

- rule: {"required":true}

### spec.container.sidecars[].volumeMounts[].secret.key

`string`

Specific key from the Secret to mount as a single file.
If not specified, all keys are mounted as files in the directory.

### spec.container.sidecars[].volumeMounts[].secret.path

`string`

If key is specified, this is the filename to use for the mounted file.
Defaults to the key name if not specified.

### spec.container.sidecars[].volumeMounts[].secret.defaultMode

`int32`

Mode bits to use on created files. Must be a value between 0 and 0777.
Defaults to 0644.

### spec.container.sidecars[].volumeMounts[].hostPath

`HostPathVolumeSource`

HostPath volume source.
Use this to mount a file or directory from the host node's filesystem.
Common for DaemonSets that need access to node-level resources.

### spec.container.sidecars[].volumeMounts[].hostPath.path

`string` · required

Path on the host to mount.

- rule: {"required":true}

### spec.container.sidecars[].volumeMounts[].hostPath.type

`string`

Type of the host path.
Valid values:
  "" - Empty string (default) means no check is performed before mounting
  "DirectoryOrCreate" - Create directory if it doesn't exist
  "Directory" - Directory must exist
  "FileOrCreate" - Create file if it doesn't exist
  "File" - File must exist
  "Socket" - UNIX socket must exist
  "CharDevice" - Character device must exist
  "BlockDevice" - Block device must exist

- rule: Type must be one of: "", "DirectoryOrCreate", "Directory", "FileOrCreate", "File", "Socket", "CharDevice", "BlockDevice"

### spec.container.sidecars[].volumeMounts[].emptyDir

`EmptyDirVolumeSource`

EmptyDir volume source.
Use this for temporary storage that is erased when the pod is removed.
Useful for scratch space, caching, or sharing data between containers.

### spec.container.sidecars[].volumeMounts[].emptyDir.medium

`string`

Medium for the empty directory.
"" (default) uses the node's default medium (typically disk).
"Memory" uses a tmpfs (RAM-backed filesystem).

Memory-backed volumes are faster but:
- Count against container memory limits
- Are lost on node restart
- Should have sizeLimit set to prevent OOM

- rule: Medium must be either "" or "Memory"

### spec.container.sidecars[].volumeMounts[].emptyDir.sizeLimit

`string`

Size limit for the empty directory.
Format: Kubernetes quantity (e.g., "1Gi", "500Mi").
Only strictly enforced when medium is "Memory".
For disk-backed volumes, this is a best-effort limit.

### spec.container.sidecars[].volumeMounts[].pvc

`PvcVolumeSource`

PersistentVolumeClaim volume source.
Use this to mount an existing PVC.
For StatefulSets, this can reference a volumeClaimTemplate.

### spec.container.sidecars[].volumeMounts[].pvc.claimName

`string` · required

Name of the PersistentVolumeClaim to mount.
For StatefulSets, this can be the name of a volumeClaimTemplate.

- rule: {"required":true}

### spec.container.sidecars[].volumeMounts[].pvc.readOnly

`bool`

Whether the PVC should be mounted read-only.
Default is false.

### spec.container.sidecars[].lifecycle

`WorkloadContainerLifecycle`

Lifecycle hooks. `post_start` runs immediately after the container starts (the
container is not Running until it completes); `pre_stop` runs before termination
and is the standard lever for connection draining — e.g. a short sleep that keeps
the endpoint serving while load balancers converge on the terminating state.

### spec.container.sidecars[].lifecycle.postStart

`WorkloadLifecycleHandler`

Runs immediately after the container is created. The container does not reach
Running until the hook completes; a failing post_start kills the container per
its restart policy.

### spec.container.sidecars[].lifecycle.postStart.exec

`ExecAction`

Execute a command inside the container; non-zero exit fails the hook.

### spec.container.sidecars[].lifecycle.postStart.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.sidecars[].lifecycle.postStart.httpGet

`HTTPGetAction`

Perform an HTTP GET against the container.

### spec.container.sidecars[].lifecycle.postStart.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.sidecars[].lifecycle.postStart.httpGet.portNumber

`int32`

### spec.container.sidecars[].lifecycle.postStart.httpGet.portName

`string`

### spec.container.sidecars[].lifecycle.postStart.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.sidecars[].lifecycle.postStart.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.sidecars[].lifecycle.postStart.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.sidecars[].lifecycle.postStart.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.sidecars[].lifecycle.postStart.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.sidecars[].lifecycle.postStart.tcpSocket

`TCPSocketAction`

Open a TCP connection against the container.

### spec.container.sidecars[].lifecycle.postStart.tcpSocket.portNumber

`int32`

### spec.container.sidecars[].lifecycle.postStart.tcpSocket.portName

`string`

### spec.container.sidecars[].lifecycle.postStart.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.sidecars[].lifecycle.postStart.sleep

`SleepAction`

Pause for a fixed number of seconds — the idiomatic pre_stop drain hook,
with no shell or sleep binary required in the image.

### spec.container.sidecars[].lifecycle.postStart.sleep.seconds

`int64`

Seconds to sleep.

- rule: {"int64":{"gte":"0"}}

### spec.container.sidecars[].lifecycle.preStop

`WorkloadLifecycleHandler`

Runs before the container is terminated by the kubelet (pod deletion, rolling
update, eviction). The termination grace period starts BEFORE the hook runs, so
keep `pod.termination_grace_period_seconds` larger than the hook's worst-case
duration. The classic zero-downtime pattern is a short sleep here so the pod
keeps serving while endpoint removal propagates.

### spec.container.sidecars[].lifecycle.preStop.exec

`ExecAction`

Execute a command inside the container; non-zero exit fails the hook.

### spec.container.sidecars[].lifecycle.preStop.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.container.sidecars[].lifecycle.preStop.httpGet

`HTTPGetAction`

Perform an HTTP GET against the container.

### spec.container.sidecars[].lifecycle.preStop.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.container.sidecars[].lifecycle.preStop.httpGet.portNumber

`int32`

### spec.container.sidecars[].lifecycle.preStop.httpGet.portName

`string`

### spec.container.sidecars[].lifecycle.preStop.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.container.sidecars[].lifecycle.preStop.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.container.sidecars[].lifecycle.preStop.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.container.sidecars[].lifecycle.preStop.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.container.sidecars[].lifecycle.preStop.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.container.sidecars[].lifecycle.preStop.tcpSocket

`TCPSocketAction`

Open a TCP connection against the container.

### spec.container.sidecars[].lifecycle.preStop.tcpSocket.portNumber

`int32`

### spec.container.sidecars[].lifecycle.preStop.tcpSocket.portName

`string`

### spec.container.sidecars[].lifecycle.preStop.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.container.sidecars[].lifecycle.preStop.sleep

`SleepAction`

Pause for a fixed number of seconds — the idiomatic pre_stop drain hook,
with no shell or sleep binary required in the image.

### spec.container.sidecars[].lifecycle.preStop.sleep.seconds

`int64`

Seconds to sleep.

- rule: {"int64":{"gte":"0"}}

### spec.container.sidecars[].securityContext

`WorkloadContainerSecurityContext`

Container-level security hardening. Settings here override the pod-level
security context for this container only.

### spec.container.sidecars[].securityContext.privileged

`bool`

Runs the container with full host access — equivalent to root on the node.
Required by some node-level agents (device managers, network plugins). Never
combine with untrusted images.

### spec.container.sidecars[].securityContext.runAsUser

`int64` · optional (explicit presence)

UID the container process runs as. Overrides the image's USER directive.

### spec.container.sidecars[].securityContext.runAsGroup

`int64` · optional (explicit presence)

Primary GID the container process runs as.

### spec.container.sidecars[].securityContext.runAsNonRoot

`bool` · optional (explicit presence)

Refuses to start the container if its effective user is root. The standard
baseline hardening — it catches images that silently default to UID 0.

### spec.container.sidecars[].securityContext.readOnlyRootFilesystem

`bool` · optional (explicit presence)

Mounts the container's root filesystem read-only. Pair with EmptyDir mounts for
paths the app must write (e.g. /tmp).

### spec.container.sidecars[].securityContext.allowPrivilegeEscalation

`bool` · optional (explicit presence)

Whether the process can gain more privileges than its parent (setuid binaries,
file capabilities). The restricted Pod Security Standard requires this to be
false. Always true when `privileged` is set, so leave it unset in that case.

### spec.container.sidecars[].securityContext.capabilities

`WorkloadCapabilities`

Linux capabilities to add or drop. The restricted profile drops ALL and adds
back only NET_BIND_SERVICE when needed. Capability names are uppercase without
the CAP_ prefix (e.g. "NET_ADMIN", "SYS_TIME").

### spec.container.sidecars[].securityContext.capabilities.add

`[]string`

Capabilities to add (e.g. "NET_BIND_SERVICE").

### spec.container.sidecars[].securityContext.capabilities.drop

`[]string`

Capabilities to drop. Use ["ALL"] as the hardened baseline.

### spec.container.sidecars[].securityContext.seccompProfile

`WorkloadSeccompProfile`

Seccomp syscall filter for the container. "RuntimeDefault" is the hardened
baseline; "Localhost" selects a node-local profile file via `localhost_profile`.

- rule: localhost_profile is required when type is "Localhost" and must be empty otherwise

### spec.container.sidecars[].securityContext.seccompProfile.type

`string` · required

Profile type: "RuntimeDefault" (the container runtime's default filter — the
recommended baseline), "Unconfined" (no filtering), or "Localhost" (a profile
file installed on the node, named via localhost_profile).

- rule: Seccomp profile type must be one of "RuntimeDefault", "Unconfined", or "Localhost"
- rule: {"required":true}

### spec.container.sidecars[].securityContext.seccompProfile.localhostProfile

`string`

Path of the profile file relative to the node's seccomp profile root. Required
when (and only meaningful when) type is "Localhost".

### spec.pod

`WorkloadPod`

Pod-level configuration shared by all replicas: identity (ServiceAccount
reference), init containers, scheduling, security hardening, DNS, and termination
behavior. For clustered stateful systems, size
`termination_grace_period_seconds` to cover a clean member handoff.

### spec.pod.serviceAccount

`string | valueFrom`

The ServiceAccount pods run as. Accepts a literal ServiceAccount name or a
reference to a KubernetesServiceAccount resource, so an infra chart deploys the
identity (with its workload-identity binding and pull secrets) and the workload
in one run. When omitted, pods run as the namespace's `default` ServiceAccount.
Permissions attach to this identity through KubernetesRbac grants; cloud access
federates through the identity's own workload_identity configuration.

- references: KubernetesServiceAccount (`status.outputs.service_account_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesServiceAccount, name: <that resource's name>, fieldPath: status.outputs.service_account_name}} -- a bare string does not parse

### spec.pod.automountServiceAccountToken

`bool` · optional (explicit presence)

Whether pods receive a projected kube-apiserver token mount. Tri-state like the
Kubernetes API: unset defers to the ServiceAccount/cluster default, false hardens
pods that never call the Kubernetes API (a security-baseline recommendation for
ordinary app workloads), true forces the mount.

### spec.pod.imagePullSecrets

`[]string | valueFrom`

Docker-registry secret names the kubelet uses to pull this workload's images.
Each entry accepts a literal secret name or a reference to a KubernetesSecret.
Prefer attaching pull secrets to the ServiceAccount when several workloads share
a registry; use this field for workload-specific registries.

- references: KubernetesSecret (`spec.name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesSecret, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.pod.initContainers

`[]WorkloadContainer`

Init containers, run to completion in order before app containers start.
Standard uses: schema migrations, config templating, waiting on dependencies.
A failing init container blocks the pod per its restart policy.

### spec.pod.initContainers[].name

`string`

The container's name, unique within the pod. Required for sidecars and init
containers (Kubernetes rejects unnamed containers); for the main app container the
module defaults it when omitted, so minimal manifests stay minimal. Must be a valid
DNS label: lowercase alphanumeric and hyphens, starting and ending alphanumeric.

- rule: Container name must be a lowercase DNS label (alphanumeric and hyphens, starting and ending with an alphanumeric character)

### spec.pod.initContainers[].image

`ContainerImage` · required

The container image, split into repository and tag so deployment pipelines can
inject a freshly built tag without rewriting the whole reference. The optional
`pull_secret_name` names an existing docker-registry secret; prefer attaching pull
secrets on the ServiceAccount (or `pod.image_pull_secrets`) so they apply pod-wide.

- rule: Image repo is required — the repository half of the image reference (e.g. "nginx" or "ghcr.io/acme/api")
- rule: Image tag is required — pin a version (e.g. "1.27.1"); avoid "latest" for anything you intend to roll back
- rule: {"required":true}

### spec.pod.initContainers[].image.repo

`string`

The repository of the image (e.g., "gcr.io/project/image").

### spec.pod.initContainers[].image.tag

`string`

The tag of the image (e.g., "latest" or "1.0.0").

### spec.pod.initContainers[].image.pullSecretName

`string`

The name of the image pull secret for private image repositories.

### spec.pod.initContainers[].imagePullPolicy

`string`

When the kubelet pulls the image. "IfNotPresent" (the Kubernetes default for tagged
images) reuses a cached copy; "Always" re-resolves the tag on every pod start —
required when a mutable tag like a branch name is reused across builds; "Never"
only uses pre-loaded images (air-gapped nodes, kind-loaded test images).

- rule: Image pull policy must be one of "Always", "IfNotPresent", or "Never"

### spec.pod.initContainers[].command

`[]string`

Entrypoint override (Kubernetes `command`, Docker ENTRYPOINT). The image's
entrypoint runs when omitted. Not executed in a shell — provide argv elements,
e.g. ["/bin/sh", "-c", "exec my-server"].

### spec.pod.initContainers[].args

`[]string`

Arguments to the entrypoint (Kubernetes `args`, Docker CMD). The image's CMD is
used when omitted. Variable references like $(VAR_NAME) are expanded from the
container's environment by the kubelet.

### spec.pod.initContainers[].workingDir

`string`

Working directory for the entrypoint. Defaults to the image's configured WORKDIR.

### spec.pod.initContainers[].ports

`[]WorkloadContainerPort`

Network ports this container exposes. Purely informational to Kubernetes for plain
pod-to-pod traffic, but load-bearing here: named ports are referenced by probes,
and `service_port` drives the Service wiring on kinds that create one
(Deployment, StatefulSet).

### spec.pod.initContainers[].ports[].name

`string` · required

Port name, e.g. "http", "grpc", "metrics". Must be a lowercase DNS label that
starts and ends alphanumeric. Named ports are referenced by probes and become the
Service port names on service-fronted kinds.

- rule: Port name must contain only lowercase alphanumeric characters and hyphens, and start and end with an alphanumeric character (e.g. "http", "grpc-web")
- rule: {"required":true}

### spec.pod.initContainers[].ports[].containerPort

`int32` · required

The port number the container listens on (1–65535).

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.pod.initContainers[].ports[].networkProtocol

`string`

L4 protocol of the port. Defaults to "TCP" when omitted — the overwhelmingly
common case, so minimal manifests need not repeat it.

- rule: The network protocol must be one of "TCP", "UDP", or "SCTP"

### spec.pod.initContainers[].ports[].appProtocol

`string`

Application protocol hint (e.g. "http", "grpc", "https"). Propagated to the
Service port's appProtocol on service-fronted kinds, where meshes and L7 load
balancers use it to pick the right protocol handling.

### spec.pod.initContainers[].ports[].servicePort

`int32`

The port the workload's Kubernetes Service exposes for this container port.
Only meaningful on kinds that create a Service (Deployment, StatefulSet); other
kinds ignore it. E.g. containerPort 8080 with servicePort 80 serves the app on
the conventional port while the process binds an unprivileged one. External
exposure is composed separately with first-class ingress kinds referencing the
workload's exported Service handle — workloads never create ingress themselves.

- rule: Service port must be between 1 and 65535

### spec.pod.initContainers[].ports[].hostPort

`int32`

Exposes the container port directly on the node's IP (hostPort). Chiefly a
DaemonSet pattern (node-level agents that must be reachable on every node);
on other kinds it constrains scheduling to one pod per node per port — prefer
a Service unless node-local reachability is the point.

- rule: Host port must be between 1 and 65535

### spec.pod.initContainers[].env

`ContainerEnv`

Environment configuration: plain variables (with Kubernetes-native value sources
and Planton cross-resource references), secret variables (materialized into a
managed Kubernetes Secret), and bulk envFrom imports.

### spec.pod.initContainers[].env.variables

`[]EnvVar`

Individual environment variables (non-sensitive).

### spec.pod.initContainers[].env.variables[].name

`string` · required

The environment variable name.
Must be a valid C_IDENTIFIER: starts with a letter or underscore,
followed by letters, digits, or underscores.

- rule: Must be a valid C_IDENTIFIER (start with letter/underscore, contain only letters, digits, underscores)
- rule: {"required":true}

### spec.pod.initContainers[].env.variables[].value

`string`

Direct literal value.

### spec.pod.initContainers[].env.variables[].valueFrom

`ValueFromRef`

Reference to another Planton resource's field.
The orchestrator resolves this and populates the value before invoking IaC modules.

### spec.pod.initContainers[].env.variables[].valueFrom.kind

`enum`

Allowed values (use exactly as shown):

- `unspecified` -- 0: Default/unspecified
- `TestCloudResourceGeneric` -- 1–49: Test/dev/custom
- `TestCloudResourceKubernetes`
- `ConfluentKafka` -- 50–199: saas platform resources
- `AtlasMongodb`
- `SnowflakeDatabase`
- `AwsAlb` -- 1000–1999: AWS resources AwsSubnet is a prerequisite because an ALB requires at least two subnets in different availability zones -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsCertManagerCert`
- `AwsCloudFront`
- `AwsDynamodb`
- `AwsEcrRepo`
- `AwsEcsCluster`
- `AwsEcsService` -- AwsEcsCluster, AwsEcsTaskDefinition, and AwsSubnet are prerequisites because a service schedules a referenced task-definition revision into a referenced live cluster and places task network interfaces into referenced subnets -- all three references must resolve first.
- `AwsEksCluster` -- AwsSubnet and AwsIamRole are prerequisites because the control plane attaches its network interfaces into referenced subnets and assumes a referenced cluster role that must already carry AmazonEKSClusterPolicy.
- `AwsIamRole`
- `AwsLambda`
- `AwsRdsCluster`
- `AwsRdsInstance`
- `AwsRoute53Zone`
- `AwsS3Bucket`
- `AwsLbTargetGroup` -- AwsVpc is a prerequisite because a target group's health checks and target registrations live inside one VPC -- the spec's vpc_id reference must resolve before the group can be created.
- `AwsSecurityGroup` -- AwsVpc is a prerequisite because every security group is created in a VPC; the E2E install profile resolves vpc_id against the VPC prerequisite.
- `AwsVpc`
- `AwsEksNodeGroup` -- AwsEksCluster is a prerequisite because nodes register with a live control plane; AwsIamRole and AwsSubnet back the node role and worker subnet references.
- `AwsIamUser`
- `AwsKmsKey`
- `AwsEc2Instance`
- `AwsClientVpn` -- Every Client VPN endpoint requires an ACM server certificate at create time; the imported self-signed fixture satisfies it. Subnets/VPC are optional composition (a zero-association endpoint is valid) -- composed scenarios declare them via the e2e-prerequisites annotation.
- `AwsDocumentDb`
- `AwsRoute53DnsRecord` -- AwsRoute53Zone is a prerequisite because every record lives inside a hosted zone -- the spec's zone_id reference must resolve before the record can be created.
- `AwsS3ObjectSet` -- AwsS3Bucket is a prerequisite because the object set's bucket reference is required -- objects cannot exist without the bucket that holds them.
- `AwsSqsQueue`
- `AwsSnsTopic`
- `AwsEventBridgeBus`
- `AwsEventBridgeRule`
- `AwsIamOidcProvider`
- `AwsIamPolicy`
- `AwsIamInstanceProfile` -- AwsIamRole is a prerequisite because an instance profile is a wrapper that must contain a role to be useful -- the profile's spec requires a role reference, so the role must be deployed first.
- `AwsLbListener` -- AwsAlb and AwsLbTargetGroup are prerequisites because a listener is an attachment point on a load balancer and its default action almost always forwards to a target group -- both references must resolve before the listener can be created.
- `AwsLbListenerRule` -- AwsLbListener is a prerequisite because a rule only exists as an attachment on a listener -- the listener_arn reference must resolve before the rule can be created.
- `AwsLaunchTemplate`
- `AwsAutoScalingGroup` -- AwsSubnet and AwsLaunchTemplate are prerequisites because a group cannot exist without subnets to place capacity in and a launch template to launch from -- the spec's subnets and launch_template references must resolve before the group can be created.
- `AwsEksAddon` -- AwsEksCluster is a prerequisite because an add-on installs onto a live control plane -- the spec's cluster_name reference must resolve before the add-on can be created.
- `AwsEksFargateProfile` -- AwsEksCluster, AwsIamRole, and AwsSubnet are prerequisites because a Fargate profile attaches to a live control plane, runs pods as a referenced pod-execution role, and launches them into referenced private subnets -- all three references must resolve first.
- `AwsEksAccessEntry` -- AwsEksCluster and AwsIamRole are prerequisites because an access entry grants a referenced IAM principal access to a live control plane -- both references must resolve before the entry can be created.
- `AwsEcsTaskDefinition` -- AwsIamRole is a prerequisite because the kind's default posture -- Fargate with the awslogs logging default -- is rejected by AWS at registration time without an execution role the agent can assume.
- `AwsHttpApiGateway`
- `AwsStepFunction` -- AwsIamRole is a prerequisite because a state machine cannot be created without an execution role it can assume -- the spec's role_arn reference must resolve before the CreateStateMachine call.
- `AwsHttpApiVpcLink` -- AwsSubnet is a prerequisite because a VPC link is a set of managed ENIs provisioned into referenced subnets -- the subnet references must resolve before the link can be created. Security groups are optional on the link, so they compose per-scenario rather than as a registry prerequisite.
- `AwsHttpApiDomain` -- AwsCertManagerCert is a prerequisite because a custom domain cannot be created without a TLS certificate in the same region covering the domain -- the spec's certificate_arn reference must resolve first.
- `AwsVpcEndpoint` -- AwsVpcEndpoint's composed E2E scenarios reference the AwsVpc prerequisite's outputs (vpc_id + default_route_table_id for gateway endpoints) and the AwsSubnet pair's subnet_id outputs (interface endpoints), so both are genuine deploy-order prerequisites.
- `AwsElasticacheUser`
- `AwsElasticacheUserGroup` -- AwsElasticacheUser is a genuine prerequisite: AWS refuses to create a user group that does not contain a user named "default", so a group's composed E2E scenario must resolve a deployed user's outputs.
- `AwsRedshiftServerlessNamespace`
- `AwsRedshiftServerlessWorkgroup` -- The namespace is a genuine prerequisite: a workgroup attaches to exactly one namespace by name at create time, so its composed E2E scenario must resolve a deployed namespace's outputs. AwsSubnet is a prerequisite because Redshift Serverless requires the workgroup's subnets to span three availability zones.
- `AwsRedisElasticache` -- AwsSubnet is a prerequisite because the module builds an ElastiCache subnet group from referenced subnets -- the spec's subnet references must resolve before the replication group can deploy.
- `AwsOpenSearchDomain`
- `AwsMemcachedElasticache`
- `AwsServerlessElasticache`
- `AwsNlb` -- AwsSubnet is a prerequisite because an NLB requires at least one subnet mapping -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsElasticIp`
- `AwsTransitGateway`
- `AwsGlobalAccelerator`
- `AwsSubnet`
- `AwsInternetGateway`
- `AwsNatGateway` -- AwsInternetGateway is a prerequisite because a public NAT gateway can only become available once the VPC it sits in has an internet gateway attached (AWS rejects the create otherwise) -- so the gateway must be deployed first. AwsVpc is a prerequisite because a REGIONAL NAT gateway (availability_mode = regional) references the VPC directly instead of a subnet.
- `AwsEgressOnlyInternetGateway`
- `AwsElasticFileSystem` -- AwsSubnet and AwsSecurityGroup are prerequisites because mount targets (required, min 1) place the file system's NFS endpoints into subnets and attach security groups -- both references must resolve before the CreateMountTarget calls.
- `AwsEfsAccessPoint` -- AwsElasticFileSystem is a prerequisite because an access point is created INTO a file system -- the spec's required file_system_id reference must resolve before the CreateAccessPoint call.
- `AwsFsxLustreFileSystem`
- `AwsFsxOpenzfsFileSystem`
- `AwsFsxWindowsFileSystem` -- Every Windows file system must join an Active Directory domain; the directory itself is external infrastructure (AWS Managed Microsoft AD or a self-managed domain), so only the network dependency is a declarable prerequisite.
- `AwsFsxOntapFileSystem`
- `AwsFsxOntapStorageVirtualMachine`
- `AwsFsxOntapVolume`
- `AwsFsxDataRepositoryAssociation`
- `AwsCognitoUserPool`
- `AwsCognitoIdentityProvider` -- AwsCognitoUserPool is a prerequisite because an identity provider is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateIdentityProvider call.
- `AwsCognitoUserPoolClient` -- AwsCognitoUserPool is a prerequisite because an app client is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateUserPoolClient call.
- `AwsCognitoResourceServer` -- AwsCognitoUserPool is a prerequisite because a resource server is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateResourceServer call.
- `AwsWafWebAcl`
- `AwsWafIpSet`
- `AwsWafRegexPatternSet`
- `AwsCloudwatchLogGroup`
- `AwsCloudwatchAlarm`
- `AwsCloudwatchCompositeAlarm`
- `AwsKinesisStream`
- `AwsKinesisFirehose` -- Every Firehose destination requires an S3 configuration (the primary target for extended_s3; the failed/all-document backup for the rest) and an IAM role Firehose assumes to write to it, so both are hard deploy prerequisites.
- `AwsKinesisStreamConsumer` -- A consumer registers against exactly one stream and cannot exist without it.
- `AwsAthenaWorkgroup`
- `AwsGlueCatalogDatabase`
- `AwsRedshiftCluster`
- `AwsSagemakerDomain` -- AI/ML A domain cannot exist without VPC subnets and a SageMaker execution role (default_user_settings.execution_role_arn is required), so both are hard deploy prerequisites.
- `AwsAppRunnerService` -- A service can run entirely on companion defaults, so the App Runner family's kinds are dependency-free leaves except the VPC connector (which cannot exist without subnets and security groups). A service's companion references (auto scaling / VPC connector / observability / WAF) are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsAppRunnerAutoScalingConfiguration`
- `AwsAppRunnerVpcConnector`
- `AwsAppRunnerObservabilityConfiguration`
- `AwsTransitGatewayVpcAttachment` -- AwsTransitGateway is a prerequisite because an attachment cannot exist without the gateway it attaches to; AwsSubnet because the attachment provisions an ENI into at least one subnet (the VPC arrives transitively through the subnet's own prerequisites).
- `AwsTransitGatewayRouteTable` -- Only the gateway is a hard prerequisite: a route table can exist empty. Associations, propagations, and routes referencing attachments are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsBatchComputeEnvironment` -- A MANAGED compute environment always launches into VPC subnets, so the subnet is a hard deploy prerequisite (security groups are required only for the Fargate types -- scenario-declared, not a registry edge).
- `AwsBatchJobQueue` -- A job queue cannot exist without at least one VALID compute environment to map onto.
- `AwsBatchSchedulingPolicy`
- `AwsBatchJobDefinition`
- `AwsCodeBuildProject` -- CI/CD
- `AwsCodePipeline`
- `AwsMwaaEnvironment` -- Workflow / Orchestration AwsSubnet and AwsSecurityGroup are prerequisites because the environment's network interfaces are placed in referenced private subnets and AWS requires at least one attached security group at creation.
- `AwsNeptuneCluster` -- Graph Database
- `AwsMemorydbCluster` -- A cluster always launches into a subnet group; the subnets are the hard deploy prerequisite. The ACL it attaches is optional composition (the built-in "open-access" ACL needs no resource) -- scenarios declare the ACL/user chain via the e2e-prerequisites annotation.
- `AwsMemorydbUser`
- `AwsMemorydbAcl` -- An empty ACL is valid (MemoryDB has no mandatory "default" member), so the user is optional composition -- the composed scenario declares it via the e2e-prerequisites annotation, never a registry edge.
- `AwsMskCluster` -- Streaming AwsSubnet and AwsSecurityGroup are prerequisites because brokers are placed in referenced subnets and AWS requires at least one attached security group at creation.
- `AwsMskServerlessCluster` -- AwsSubnet is a prerequisite because the serverless cluster's network interfaces are placed in referenced subnets (security groups are optional -- AWS attaches the VPC default group when none are referenced).
- `AwsLambdaEventSourceMapping` -- AwsLambda is a prerequisite because a mapping cannot exist without the function it invokes (a required reference). Event sources (SQS, Kinesis, DynamoDB, MSK) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsSnsSubscription` -- AwsSnsTopic is a prerequisite because a subscription cannot exist without the topic it subscribes to (a required reference). Endpoints (SQS queues, Lambda functions, Firehose streams) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsPlantonRunner` -- AwsSubnet is a prerequisite because the runner appliance places its network interfaces into referenced subnets -- the placement reference must resolve before the appliance can deploy.
- `AwsRoute53HealthCheck`
- `AwsSesConfigurationSet` -- Both SES kinds are dependency-free leaves: an identity's configuration set is optional composition (scenarios declare it via the e2e-prerequisites annotation), and a configuration set's event destinations reference other kinds only optionally.
- `AwsSesEmailIdentity`
- `AwsSecretsManagerSecret` -- A dependency-free leaf: the KMS key, rotation Lambda, and external rotation role references are all optional composition -- scenarios declare them via the e2e-prerequisites annotation, never registry edges.
- `AwsOpenSearchServerlessCollection` -- A dependency-free leaf: the collection-scoped encryption/network/ data-access/retention policies are module-rendered, and the KMS key and data-access principal references are optional composition (e2e-prerequisites annotation).
- `AwsBedrockGuardrail` -- A dependency-free leaf: the KMS key reference is optional composition (e2e-prerequisites annotation); published versions are folded satellites of the guardrail itself.
- `AwsBedrockCustomModel` -- AwsIamRole is a prerequisite because Bedrock assumes the job role to read training data and write outputs; the S3 locations and KMS key are optional composition (e2e-prerequisites annotation).
- `AwsBedrockInferenceProfile` -- A dependency-free leaf: the model source is a foundation model or an AWS system-defined cross-region profile, never a customer resource.
- `AwsBedrockProvisionedThroughput` -- A dependency-free leaf in the registry: capacity is typically bought for an AwsBedrockCustomModel (the default reference), but foundation model ARNs are equally legal, so the edge is optional composition.
- `AwsBedrockModelAccess` -- A dependency-free leaf: the agreement covers an AWS-listed foundation model, never a customer resource.
- `AwsBedrockAgent` -- AwsIamRole is a prerequisite because the Bedrock service assumes the agent resource role to invoke models, action-group Lambdas, and knowledge bases; the guardrail, KMS key, provisioned throughput, and collaborator/knowledge-base edges are optional composition (e2e-prerequisites annotation). Action groups, aliases, collaborators, and knowledge-base associations are folded satellites of the agent.
- `AwsBedrockKnowledgeBase` -- AwsIamRole is a prerequisite because the Bedrock service assumes the knowledge-base role to read data sources, call the embedding model, and read/write the vector store; the vector-store and data-source reference edges (OpenSearch, S3, Secrets Manager, ...) are optional composition (e2e-prerequisites annotation). Data sources are folded satellites of the knowledge base.
- `AwsBedrockFlow` -- AwsIamRole is a prerequisite because the Bedrock service assumes the flow execution role to invoke the models, agents, knowledge bases, and Lambdas its nodes reference; the node-level reference edges are optional composition (e2e-prerequisites annotation).
- `AwsBedrockPrompt` -- A dependency-free leaf: variants target AWS-listed foundation models by ID; targeting another agent's alias is optional composition (e2e-prerequisites annotation).
- `AzureResourceGroup` -- 2000–2999: Azure resources
- `AzureAksCluster` -- AzureResourceGroup is the only required parent: the cluster is created inside a referenced resource group. Subnet is optional on the default node pool (AKS provisions managed networking when unset).
- `AzureAksNodePool` -- AzureAksCluster is a prerequisite because a node pool attaches to an existing cluster by ARM ID; the resource group chains transitively.
- `AzureContainerRegistry` -- AzureResourceGroup is a prerequisite because a container registry is created inside a resource group.
- `AzureDnsZone` -- AzureResourceGroup is a prerequisite because the DNS zone is created inside a referenced resource group that must already exist.
- `AzureKeyVault` -- AzureResourceGroup is a prerequisite because a key vault is created inside a referenced resource group in composed environments.
- `AzureVirtualNetwork` -- AzureResourceGroup is a prerequisite because a virtual network is created inside a referenced resource group in composed environments.
- `AzureNatGateway` -- AzureResourceGroup is a prerequisite because a NAT gateway is created inside a referenced resource group in composed environments.
- `AzureVirtualMachine` -- AzureNetworkInterface is a prerequisite because a virtual machine attaches at least one NIC (the subnet, network, and resource group chain transitively through the NIC's own prerequisites).
- `AzureStorageAccount` -- AzureResourceGroup is a prerequisite because a storage account is created inside a referenced resource group in composed environments.
- `AzureDnsRecord` -- AzureDnsZone is a prerequisite because a record set is created inside a referenced zone (the resource group chains transitively through the zone). Public DNS zone names are not globally unique, so a shared zone fixture is safe to recreate across scenarios.
- `AzureSubnet` -- AzureVirtualNetwork is a prerequisite because a subnet is an ARM child of a referenced network -- the network must exist before the subnet can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzureNetworkSecurityGroup` -- AzureResourceGroup is a prerequisite because a network security group is created inside a referenced resource group in composed environments.
- `AzurePublicIp` -- AzureResourceGroup is a prerequisite because a public IP is created inside a referenced resource group in composed environments.
- `AzurePrivateEndpoint` -- AzureSubnet is a prerequisite because a private endpoint draws its private IP from a referenced subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite). The connection target is polymorphic and the DNS zones / ASGs are optional, so none of those are prerequisites.
- `AzurePrivateDnsZone` -- AzureResourceGroup is a prerequisite because a private DNS zone is created inside a referenced resource group in composed environments.
- `AzureApplicationGateway` -- AzureSubnet is a prerequisite because a gateway cannot exist without its dedicated gateway_ip_configuration subnet (the network and resource group chain transitively through the subnet's own prerequisites); public frontends additionally reference a public IP, but private-only gateways are legal, so it is not a registry prerequisite.
- `AzureLoadBalancer` -- AzureResourceGroup is a prerequisite because a load balancer is created inside a referenced resource group (frontends additionally reference subnets or public IPs, but neither is universally required, so they are not registry prerequisites).
- `AzureRouteTable` -- AzureResourceGroup is a prerequisite because a route table is created inside a referenced resource group in composed environments.
- `AzurePrivateDnsZoneVirtualNetworkLink` -- AzurePrivateDnsZone and AzureVirtualNetwork are prerequisites because a virtual network link is a child resource of a referenced zone and binds it to a referenced network -- both must exist before the link can be written. (The resource group arrives transitively through the zone's and network's own prerequisite declarations.)
- `AzureVirtualNetworkPeering` -- AzureVirtualNetwork is a prerequisite because a peering is an ARM child of its local network and binds it to a remote network -- the local network must exist before the peering can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzurePublicIpPrefix` -- AzureResourceGroup is a prerequisite because a public IP prefix is created inside a referenced resource group in composed environments.
- `AzureNetworkInterface` -- AzureSubnet is a prerequisite because a network interface's IP configurations deploy into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureManagedDisk` -- AzureResourceGroup is a prerequisite because a managed disk is created inside a resource group.
- `AzureVirtualMachineScaleSet` -- AzureSubnet is a prerequisite because every scale-set instance's network interface deploys into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureKeyVaultKey` -- AzureKeyVault is a prerequisite because a key is a data-plane object inside a referenced vault -- the vault must exist before the key can be written (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultCertificate` -- AzureKeyVault is a prerequisite because a certificate is a data-plane object inside a referenced vault -- the vault must exist before the certificate can be enrolled or imported (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultSecret` -- AzureKeyVault is a prerequisite because a secret is a data-plane object inside a referenced vault -- the vault must exist before the secret can be written (the resource group chains transitively through the vault's own prerequisite). Part of the Key Vault family (2005, 2025-2026) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureWebApplicationFirewallPolicy` -- AzureResourceGroup is a prerequisite because a WAF policy is created inside a referenced resource group; the Application Gateways that attach the policy reference it, never the reverse.
- `AzureApplicationSecurityGroup` -- AzureResourceGroup is a prerequisite because an application security group is created inside a referenced resource group; network interfaces, scale-set IP configurations, and NSG security rules reference the group, never the reverse.
- `AzureDiskEncryptionSet` -- AzureKeyVaultKey is a prerequisite because a disk encryption set wraps customer data with a referenced key -- the key (and its vault, which chains transitively) must exist before the set can resolve the key URL at create time.
- `AzurePostgresqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureRedisCache` -- AzureResourceGroup is a prerequisite because the cache is created inside a referenced resource group (VNet injection additionally references a dedicated subnet, but only the Premium tier supports it, so it is not a registry prerequisite).
- `AzureCosmosdbAccount` -- AzureResourceGroup is a prerequisite because the account is created inside a referenced resource group.
- `AzureMssqlServer` -- AzureResourceGroup is a prerequisite because the logical server is created inside a referenced resource group.
- `AzureMysqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureMssqlDatabase` -- The parent logical server is referenced via server_id, not auto-deployed: E2E scenarios declare their own server fixture (minimal-server.yaml or the pool-attach chain through AzureMssqlElasticPool) so sequential subtests never destroy and recreate the same globally unique server_name.
- `AzureMssqlElasticPool` -- AzureMssqlServer is a prerequisite because every elastic pool lives on a referenced logical server (the server's resource group is transitive).
- `AzureRedisLinkedServer` -- The target and linked caches are referenced via ARM ids, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicy` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicyAssignment` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureContainerAppEnvironment` -- AzureResourceGroup is a prerequisite because the environment is created inside a referenced resource group that must already exist.
- `AzureContainerApp` -- AzureContainerAppEnvironment is a prerequisite because every app runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureServicePlan` -- AzureResourceGroup is a prerequisite because the plan is created inside a referenced resource group that must already exist.
- `AzureFunctionApp` -- AzureServicePlan is a prerequisite because a function app runs on a referenced plan (the resource group arrives transitively through the plan). The required storage account is deliberately NOT a registry prerequisite: storage-account names are globally unique, so scenarios bring their own scenario-local account fixtures.
- `AzureLinuxWebApp` -- AzureServicePlan is a prerequisite because a web app runs on a referenced plan (the resource group arrives transitively through the plan).
- `AzureContainerAppJob` -- AzureContainerAppEnvironment is a prerequisite because a job runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureContainerAppEnvironmentStorage` -- AzureContainerAppEnvironment is a prerequisite because the storage registration lives on a referenced environment. The Azure Files share and storage account are deliberately NOT registry prerequisites: storage-account names are globally unique, so scenarios bring their own scenario-local account + share fixtures.
- `AzureContainerAppEnvironmentDaprComponent` -- AzureContainerAppEnvironment is a prerequisite because the Dapr component is registered on a referenced environment.
- `AzureContainerAppEnvironmentCertificate` -- AzureContainerAppEnvironment is a prerequisite because the certificate is stored on a referenced environment.
- `AzureContainerAppEnvironmentManagedCertificate` -- AzureContainerAppEnvironment is a prerequisite because the managed certificate is provisioned on a referenced environment.
- `AzureLogAnalyticsWorkspace` -- AzureResourceGroup is a prerequisite because the workspace is created inside a referenced resource group that must already exist.
- `AzureApplicationInsights` -- AzureLogAnalyticsWorkspace is a prerequisite because workspace-based Application Insights stores its telemetry in a referenced workspace (the resource group chains transitively through the workspace).
- `AzureMonitorDiagnosticSetting` -- AzureLogAnalyticsWorkspace is a prerequisite because the setting's scenarios route a fixture workspace's telemetry (the workspace doubles as target and destination); the target itself is polymorphic.
- `AzureMonitorActionGroup` -- AzureResourceGroup is a prerequisite because the action group is created inside a referenced resource group that must already exist.
- `AzureMonitorMetricAlert` -- AzureMonitorActionGroup is a prerequisite because a metric alert's actions fire into a referenced action group (the resource group chains transitively); alert scopes are polymorphic.
- `AzureMonitorScheduledQueryAlert` -- AzureLogAnalyticsWorkspace is a prerequisite because the rule queries a referenced workspace scope; AzureMonitorActionGroup because its action fires into a referenced action group.
- `AzureMonitorActivityLogAlert` -- AzureMonitorActionGroup is a prerequisite because an activity log alert's actions fire into a referenced action group (the resource group chains transitively). The alert itself is subscription-global and its scopes are polymorphic.
- `AzureApplicationInsightsStandardWebTest` -- AzureApplicationInsights is a prerequisite because a standard web test binds to a referenced Application Insights component (the resource group chains transitively through the component).
- `AzureUserAssignedIdentity` -- AzureResourceGroup is a prerequisite because the identity is created inside a referenced resource group that must already exist.
- `AzureRoleAssignment` -- AzureResourceGroup and AzureUserAssignedIdentity are prerequisites because an assignment grants a role at a referenced scope (most commonly a resource group) to a referenced principal (most commonly a managed identity) -- both must exist before the grant can be written.
- `AzureRoleDefinition` -- AzureResourceGroup is a prerequisite because a custom role definition is created at a referenced scope, most commonly a resource group in composed environments -- the scope must exist before the definition can be written.
- `AzureFederatedIdentityCredential` -- AzureUserAssignedIdentity is the prerequisite because a federated identity credential is a child resource of a referenced managed identity -- the identity must exist before the credential can be written on it. (The resource group arrives transitively through the identity's own prerequisite declaration.)
- `AzureServiceBusNamespace` -- AzureResourceGroup is a prerequisite because a Service Bus namespace is created inside a referenced resource group in composed environments. The namespace is the container every Service Bus messaging entity (queue, topic, subscription, authorization rule, geo-DR pairing) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureEventHubNamespace` -- AzureResourceGroup is a prerequisite because an Event Hub namespace is created inside a referenced resource group in composed environments. The namespace is the container every Event Hubs entity (event hub, consumer group, authorization rule, schema group, geo-DR pairing, customer-managed key) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureServiceBusQueue`
- `AzureServiceBusTopic`
- `AzureServiceBusSubscription`
- `AzureServiceBusAuthorizationRule`
- `AzureServiceBusDisasterRecoveryConfig`
- `AzureEventHub`
- `AzureEventHubConsumerGroup`
- `AzureEventHubAuthorizationRule`
- `AzureFrontDoorProfile` -- AzureResourceGroup is a prerequisite because a Front Door profile is created inside a referenced resource group in composed environments. The profile is the container every Front Door delivery resource (endpoint, origin group, origin, route) nests under.
- `AzureFrontDoorEndpoint` -- AzureFrontDoorProfile is a prerequisite because an endpoint is an ARM child of a referenced profile -- the profile must exist before the endpoint can be written. (The resource group arrives transitively through the profile's own prerequisite declaration.)
- `AzureFrontDoorOriginGroup` -- AzureFrontDoorProfile is a prerequisite because an origin group is an ARM child of a referenced profile.
- `AzureFrontDoorOrigin` -- AzureFrontDoorOriginGroup is a prerequisite because an origin is an ARM child of a referenced origin group (the profile and resource group chain transitively).
- `AzureFrontDoorRoute` -- A route attaches to an endpoint (its ARM parent) and forwards to an origin group whose origins must exist before ARM accepts the route -- so both the endpoint and the origin chain are genuine deploy-order prerequisites.
- `AzureFrontDoorRuleSet` -- AzureFrontDoorProfile is a prerequisite because a rule set is an ARM child of a referenced profile. The rules live inside the set (they form one ordered policy document); routes attach the set by ARM ID.
- `AzureFrontDoorCustomDomain` -- AzureFrontDoorProfile is a prerequisite because a custom domain is an ARM child of a referenced profile. The DNS zone and (for bring-your-own certificates) the Front Door secret are optional references, not deploy-order prerequisites.
- `AzureFrontDoorSecret` -- AzureFrontDoorSecret is a prerequisite-light kind: only the profile (its ARM parent) must exist. The Key Vault certificate it wraps is a reference resolved before the module runs; its vault chain is exercised through scenario-local fixtures in E2E.
- `AzureFrontDoorFirewallPolicy` -- AzureResourceGroup is a prerequisite because the Front Door WAF policy is created inside a referenced resource group -- it is a GLOBAL resource, not a profile child (a different ARM type than the regional Application Gateway WAF policy). Security policies attach it to profiles; the policy itself depends on nothing else.
- `AzureFrontDoorSecurityPolicy` -- A security policy is an ARM child of a profile that associates a referenced WAF policy with referenced domains -- so the endpoint (the default-domain association target; the profile arrives transitively through it) and the WAF policy are genuine deploy-order prerequisites.
- `AzureStorageContainer` -- None of the storage data-service kinds declares a registry prerequisite on AzureStorageAccount: account names are GLOBALLY unique and Azure holds a just-deleted name, so a recreate-per-scenario fixture would hang -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the storage_account_id reference itself.
- `AzureStorageShare`
- `AzureStorageQueue`
- `AzureStorageTable`
- `AzureStorageEncryptionScope`
- `AzureStorageDataLakeGen2Filesystem`
- `AzureStorageLocalUser`
- `AzureStorageObjectReplication`
- `AzureCosmosdbSqlDatabase` -- None of the Cosmos DB data-service kinds declares a registry prerequisite on AzureCosmosdbAccount: account names are GLOBALLY unique DNS labels, so a recreate-per-scenario fixture would risk name-reuse hangs -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the cosmosdb_account_id / parent-database references themselves.
- `AzureCosmosdbSqlContainer`
- `AzureCosmosdbMongoDatabase`
- `AzureCosmosdbMongoCollection`
- `AzureCosmosdbSqlRoleDefinition`
- `AzureCosmosdbSqlRoleAssignment`
- `AzureManagedRedis` -- AzureResourceGroup is the cluster's only registry prerequisite: the cluster is created inside a referenced resource group. The geo-replication and access-policy-assignment children declare NO prerequisite on AzureManagedRedis: clusters are expensive, slow-provisioning parents, so their E2E scenarios declare scenario-local cluster fixtures instead of recreating a shared one per scenario. Deploy ordering in composed environments still flows from the managed_redis_id references themselves.
- `AzureManagedRedisGeoReplication`
- `AzureManagedRedisAccessPolicyAssignment`
- `AzureEventHubDisasterRecoveryConfig`
- `AzureEventHubSchemaGroup`
- `AzureEventHubCluster` -- AzureResourceGroup is a prerequisite because a dedicated Event Hubs cluster is created inside a referenced resource group in composed environments. Note: clusters cannot be deleted for 4 hours after creation (Azure's moratorium), so E2E treats this kind as offline-gated.
- `AzureEventHubNamespaceCustomerManagedKey`
- `AzureMssqlFailoverGroup` -- AzureMssqlServer is a prerequisite because a failover group is created on a referenced primary logical server and points at a partner server; the primary (and its resource group, which chains transitively) must exist before the group can be written.
- `AzureContainerAppCustomDomain` -- AzureContainerApp is a prerequisite because the domain binding lives in a referenced app's ingress configuration (the environment and resource group chain transitively through the app).
- `AzureFirewallPolicy`
- `AzureFirewallPolicyRuleCollectionGroup` -- AzureFirewallPolicy is a prerequisite because a rule collection group is a child document of a referenced policy (the resource group chains transitively through the policy).
- `AzureFirewall` -- AzureSubnet is a prerequisite because a VNet-deployed firewall's data path lives in a dedicated subnet that must be named exactly "AzureFirewallSubnet" (the virtual network and resource group chain transitively through the subnet). The E2E install profile publishes a fixture subnet with that exact name and a /26 prefix.
- `AzureIpGroup`
- `AzureVirtualNetworkGateway` -- AzureSubnet is a prerequisite because every virtual network gateway lives in a dedicated subnet named exactly "GatewaySubnet" (the virtual network and resource group chain transitively through the subnet); the subnet install profile publishes a fixture instance with that exact ARM name. AzurePublicIp is a prerequisite because a VPN-type gateway (the default shape) requires a public IP per ip configuration; the address install profile publishes a dedicated zone-redundant instance (a gateway binds its address exclusively, and the AZ gateway SKUs require zones on it).
- `AzureVirtualNetworkGatewayConnection` -- Both gateways are prerequisites: a connection joins a virtual network gateway to a far side, and the site-to-site far side is a local network gateway (the GatewaySubnet, VNet, and resource group chain transitively through the virtual network gateway).
- `AzureLocalNetworkGateway`
- `AzurePrivateLinkService` -- AzureSubnet is the sole prerequisite: every NAT ip configuration draws its address from a subnet with private-link-service network policies disabled (the subnet install profile publishes a fixture instance with that flag off). The Standard load balancer whose frontend the service typically fronts is NOT a registry prerequisite -- the spec's destination is an exactly-one-of (load balancer frontend OR fixed destination IP), so scenarios that use the load-balancer shape declare it via the planton.dev/e2e-prerequisites annotation instead.
- `AzureExpressRouteCircuit`
- `AzureExpressRouteCircuitPeering` -- The circuit is the prerequisite: a peering is an ARM child of the circuit, addressed by the circuit's name (the resource group chains transitively through the circuit).
- `AzureExpressRouteGateway` -- The hub is the prerequisite: ARM requires an ExpressRoute Gateway to be deployed INTO a Virtual WAN hub (the WAN and resource group chain transitively through the hub).
- `AzureExpressRoutePort` -- ExpressRoute Port: your own physical port pair on a Microsoft edge router (ExpressRoute Direct), from whose bandwidth circuits are carved. Self-contained -- only the resource group is required.
- `AzureVirtualWan` -- Virtual WAN: the umbrella of Azure's managed hub-and-spoke networking, under which virtual hubs and their gateways are created. Self-contained -- only the resource group is required.
- `AzureVirtualHub` -- The WAN is the prerequisite: this kind models the Virtual WAN hub (virtual_wan_id is required; standalone hubs are the legacy Route Server construction, which has its own ARM surface). The resource group chains transitively through the WAN.
- `AzureVirtualHubConnection` -- Both sides of the attachment are prerequisites: the hub being joined and the spoke virtual network being attached.
- `AzureVpnGateway` -- The hub is the prerequisite: ARM deploys a Virtual WAN VPN gateway INTO a virtual hub (virtual_hub_id is required and immutable; the WAN and resource group chain transitively through the hub). ARM allows one VPN gateway per hub.
- `AzureVpnGatewayConnection` -- Both ends of the tunnel are prerequisites: a connection is an ARM child of the VPN gateway and pins each of its links to a specific link of the remote VPN site (the hub, WAN, and resource group chain transitively through the gateway).
- `AzureVpnSite` -- The WAN is the prerequisite: a VPN site is the Virtual WAN world's address-book entry for one branch location (virtual_wan_id is required; the classic-world sibling without a WAN is AzureLocalNetworkGateway). The resource group chains transitively through the WAN.
- `AzurePointToSiteVpnGateway` -- The hub and the server configuration are both prerequisites: a point-to-site VPN gateway deploys INTO a virtual hub (one P2S gateway per hub, a slot separate from the hub's site-to-site VPN gateway) and is born pointing at the VPN server configuration that defines how its users authenticate -- both ARM-required and fixed at creation. The WAN and resource group chain transitively through the hub.
- `AzureVpnServerConfiguration` -- Self-contained -- only the resource group is required: a VPN server configuration is the reusable "who may connect and how" authentication policy (Entra ID / certificate / RADIUS) that point-to-site VPN gateways attach to; it references no other Azure resource.
- `AzureCognitiveAccount` -- Self-contained -- only the resource group is required: an Azure AI services account (Azure OpenAI, the multi-service AIServices account, the single-service accounts) needs no other Azure resource; subnets (network rules), Key Vault keys (CMK), storage accounts and user-assigned identities are optional references.
- `AzureCognitiveDeployment` -- An ARM child of its account: a model deployment (which model runs, at which throughput class) exists only on an Azure AI services account of kind "OpenAI" or "AIServices".
- `AzureCognitiveAccountProject` -- An ARM child of its account: an AI Foundry project exists only on an "AIServices"-kind account with project management enabled.
- `AzureMachineLearningWorkspace` -- The workspace REQUIRES all three companion services at creation (default storage, secrets vault, telemetry) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureMachineLearningDatastore` -- An ARM child of its workspace. The storage target (container, filesystem or share) is scenario-declared via the e2e-prerequisites annotation -- only the blob scenario needs a container, so it is not a kind-wide prerequisite.
- `AzureMachineLearningComputeCluster` -- An ARM child of its workspace (.../computes/{name}) -- the auto-scaling pool of VMs training jobs run on.
- `AzureMachineLearningComputeInstance` -- An ARM child of its workspace (.../computes/{name}) -- a single always-on VM serving as one data scientist's cloud workstation.
- `AzureAiFoundry` -- The hub REQUIRES both companion services at creation (secrets vault, default storage) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureAiFoundryProject` -- Deploys into its hub's resource group (the provider derives the group from the hub reference -- the project spec carries none).
- `AzureSearchService`
- `AzureMachineLearningOnlineEndpoint` -- An ARM child of its workspace (.../onlineEndpoints/{name}) -- the stable scoring address applications call. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningOnlineDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the running copy of a model the endpoint's traffic map routes to.
- `AzureMachineLearningBatchEndpoint` -- An ARM child of its workspace (.../batchEndpoints/{name}) -- the stable address batch scoring jobs are submitted to. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningBatchDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the job recipe (model, compute, batching behavior) the endpoint's default-deployment pointer routes submissions to.
- `AzureRecoveryServicesVault` -- The Recovery Services vault (Microsoft.RecoveryServices/vaults) -- the safe that classic Azure Backup data and Site Recovery configuration live in. Backup policies and protected items are ARM children of a vault.
- `AzureBackupPolicyVm` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern IaaS VM backups.
- `AzureBackupProtectedVm` -- An ARM child of its vault (.../protectedItems/...) -- the binding that puts one virtual machine under a backup policy's protection.
- `AzureBackupPolicyFileShare` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern Azure Files share backups (snapshot or vaulted).
- `AzureBackupProtectedFileShare` -- An ARM child of its vault (.../protectedItems/AzureFileShare;...) -- the binding that puts one Azure Files share under a backup policy's protection. The share's storage account must already be registered with the vault (AzureBackupContainerStorageAccount).
- `AzureDataProtectionBackupVault` -- The Data Protection backup vault (Microsoft.DataProtection/ backupVaults) -- the safe that MODERN Azure Backup data lives in (managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, Data Lake storage). Backup policies and backup instances are ARM children of a vault.
- `AzureDataProtectionBackupPolicy` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules for ONE Data Protection datasource type (blob storage, disk, Kubernetes cluster, MySQL/PostgreSQL flexible server, or Data Lake storage), modeled as one kind with variant blocks.
- `AzureDataProtectionBackupInstance` -- An ARM child of its vault (.../backupInstances/{name}) -- the binding that puts ONE datasource (a managed disk, a storage account's blob services, an AKS cluster, a MySQL/PostgreSQL flexible server, or a Data Lake storage account) under a Data Protection backup policy, modeled as one kind with variant blocks. The vault's managed identity must hold the datasource roles Azure Backup requires BEFORE the instance is created.
- `AzureBastionHost` -- AzureSubnet and AzurePublicIp are prerequisites because a dedicated-infrastructure Bastion host (Basic/Standard/Premium -- the default shapes) deploys into a subnet named exactly "AzureBastionSubnet" and binds a Standard static public IP EXCLUSIVELY (the virtual network and resource group chain transitively through the subnet). The Developer SKU instead attaches to a virtual network directly and uses neither.
- `AzureNetworkWatcherFlowLog` -- AzureVirtualNetwork and AzureStorageAccount are prerequisites because a flow log records a network-scoped target (a virtual network in the common case; subnets and network interfaces chain through the network) into a referenced storage account. The regional Network Watcher parent is NOT a prerequisite: Azure auto-creates it ("NetworkWatcher_{region}" in "NetworkWatcherRG") the moment the region hosts a virtual network, and the flow log references it by name. Traffic Analytics' Log Analytics workspace is an optional arm, declared by scenarios that use it.
- `AzurePrivateDnsResolver` -- AzureVirtualNetwork and AzureSubnet are prerequisites because a DNS Private Resolver anchors to a referenced virtual network (at most ONE resolver per network -- Azure enforces it) and each of its inbound/outbound endpoints occupies its own dedicated subnet delegated to "Microsoft.Network/dnsResolvers" (the resource group chains transitively through the network and subnets).
- `AzurePrivateDnsResolverForwardingRuleset` -- AzurePrivateDnsResolver is a prerequisite because a DNS forwarding ruleset steers a resolver's OUTBOUND endpoints -- it binds their ARM ids (at most 2, same resolver) at creation. (The resource group and network chain transitively through the resolver's own prerequisite declarations.)
- `AzureBackupContainerStorageAccount` -- Registers a storage account with a Recovery Services vault as a backup container (.../protectionContainers/StorageContainer;...) -- one registration per storage-account-and-vault pair, required BEFORE any of the account's file shares can be protected. Part of the backup family (2175-2179) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureDataProtectionResourceGuard` -- The Data Protection Resource Guard (Microsoft.DataProtection/ resourceGuards) -- the approval gate behind Multi-User Authorization: privileged vault operations (disabling soft delete, reducing retention) require an approval through a guard, which typically lives in a DIFFERENT administrator's scope. Vaults reference a guard by its ARM ID. Part of the backup family (2175-2182) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzurePrivateDnsResolverVirtualNetworkLink` -- The attachment that makes a DNS forwarding ruleset take effect in one virtual network ({ruleset_id}/virtualNetworkLinks/{name}) -- one link per ruleset-network pair, up to 500 per ruleset, spokes joining and leaving independently (which is why the link is a standalone kind, exactly like AzurePrivateDnsZoneVirtualNetworkLink). Part of the DNS Private Resolver family (2186-2187) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `GcpArtifactRegistryRepo` -- 3000–3999: GCP resources
- `GcpTargetHttpsProxy` -- The URL map is the parent a proxy cannot exist without; the classic compute certificate kinds and the SSL policy are the fixture parents the committed scenarios attach. The Certificate Manager certificate list (certificate_manager_certificates, honored only by the cross-region internal ALB) is optional composition -- a scenario that arms it declares GcpCertManagerCert via the e2e-prerequisites annotation, never a registry edge that would tax every proxy and forwarding-rule chain.
- `GcpCloudFunction`
- `GcpCloudRun`
- `GcpCloudSql`
- `GcpDnsZone`
- `GcpGcsBucket`
- `GcpGkeCluster`
- `GcpIamCustomRole`
- `GcpProject`
- `GcpVpcNetwork`
- `GcpSubnetwork`
- `GcpRouterNat`
- `GcpGkeNodePool`
- `GcpServiceAccount`
- `GcpGkeWorkloadIdentityBinding`
- `GcpCertManagerCert`
- `GcpComputeInstance`
- `GcpDnsRecord`
- `GcpProjectIamMember`
- `GcpFirewallRule`
- `GcpGlobalAddress`
- `GcpCloudArmorPolicy`
- `GcpHealthCheck`
- `GcpBackendBucket`
- `GcpBackendService`
- `GcpRegionNetworkEndpointGroup`
- `GcpUrlMap`
- `GcpManagedSslCertificate`
- `GcpTargetHttpProxy`
- `GcpAlloydbCluster`
- `GcpRedisInstance`
- `GcpFirestoreDatabase`
- `GcpSpannerInstance`
- `GcpSpannerDatabase`
- `GcpBigtableInstance`
- `GcpMemorystoreInstance`
- `GcpCloudSqlDatabase`
- `GcpCloudSqlUser`
- `GcpAlloydbInstance`
- `GcpAlloydbUser`
- `GcpSpannerBackupSchedule`
- `GcpBigtableTable`
- `GcpFirestoreBackupSchedule`
- `GcpFirestoreIndex`
- `GcpBigQueryDataset`
- `GcpDataprocCluster`
- `GcpDataprocAutoscalingPolicy`
- `GcpBigQueryTable`
- `GcpPubSubTopic`
- `GcpPubSubSubscription`
- `GcpCloudTasksQueue`
- `GcpCloudSchedulerJob`
- `GcpPubSubSchema`
- `GcpVertexAiNotebook`
- `GcpVertexAiEndpoint`
- `GcpVertexAiIndex`
- `GcpVertexAiIndexEndpoint` -- Vector Search IndexEndpoint — distinct from the online-prediction GcpVertexAiEndpoint (671); different GCP resources, different kinds.
- `GcpVertexAiDeployedIndex`
- `GcpCloudComposerEnvironment`
- `GcpCloudComposerUserWorkloadsSecret`
- `GcpCloudComposerUserWorkloadsConfigMap`
- `GcpKmsKeyRing`
- `GcpKmsKey`
- `GcpKmsKeyIamMember`
- `GcpFilestoreInstance`
- `GcpWorkloadIdentityPool` -- 3101–3109: IAM/identity family (overflow block; the 3000–3022 foundation/security sub-band is fully allocated)
- `GcpWorkloadIdentityPoolProvider`
- `GcpServiceAccountIamMember`
- `GcpGlobalForwardingRule` -- 3110–3119: networking/load-balancer family (overflow block; the 3023–3029 LB sub-band is fully allocated)
- `GcpSslPolicy`
- `GcpSslCertificate`
- `GcpServiceNetworkingConnection`
- `GcpAddress`
- `GcpServiceConnectionPolicy`
- `GcpCertManagerDnsAuthorization`
- `GcpCertificateMap` -- GcpCertManagerCert is a prerequisite because a map entry binds hostnames to EXISTING certificates — the canonical map references a certificate fixture's resource name.
- `GcpCloudRunJob` -- 3120–3129: GCP serverless overflow
- `GcpServerlessVpcConnector`
- `GcpComputeDisk` -- 3130–3139: GCP compute overflow (the 3000–3022 foundation sub-band that holds GcpComputeInstance is fully allocated)
- `GcpComputeMig` -- GcpVpcNetwork is a prerequisite because the canonical group runs its fleet on a dedicated custom-mode VPC — a managed instance group's template must attach every VM to a network, and the default VPC is never assumed.
- `GcpMonitoringNotificationChannel` -- 3140–3149: GCP observability & log routing
- `GcpMonitoringAlertPolicy` -- GcpMonitoringNotificationChannel is a prerequisite because the policy's canonical shape references a channel to notify — a policy without a delivery endpoint measures but never pages.
- `GcpMonitoringUptimeCheck`
- `GcpLoggingSink` -- GcpGcsBucket is a prerequisite because the canonical sink exports to a Cloud Storage bucket — the cheapest destination that proves the whole writer-identity grant flow.
- `GcpMonitoringDashboard`
- `GcpMonitoringSlo`
- `GcpLogBucket`
- `GcpLogMetric`
- `GcpSecretManagerSecret` -- 3150–3159: GCP security & identity GcpServiceAccount is a prerequisite because the canonical secret grants secretAccessor to a workload service account — the access story the kind exists to model.
- `GcpIdentityPlatformConfig`
- `GcpIdentityPlatformTenant` -- GcpIdentityPlatformConfig is a prerequisite because tenants exist only in projects whose Identity Platform config enables multi_tenant.allow_tenants — a tenant without the initialized, tenant-enabled project config cannot be created at all.
- `GcpIamOauthClient`
- `GcpIamDenyPolicy`
- `GcpCloudRunDomainMapping` -- 3160–3169: GCP serverless edge GcpCloudRun is a prerequisite because a domain mapping exists only to point a verified domain at a running Cloud Run service — the route it maps must already exist for the mapping to be created at all.
- `GcpWorkflow`
- `GcpEventarcTrigger` -- GcpCloudRun is a prerequisite because the canonical trigger routes a Pub/Sub messagePublished event to a Cloud Run service — the destination story the kind exists to model.
- `GcpEventarcMessageBus`
- `KubernetesNamespace` -- 4000–4999: Kubernetes resources, organized in family sub-bands (4030–4069 also hosts CNI/autoscaling/DR addons; 4130–4149 hosts analytics & ML; 4190–4199 reserved for growth) 4000–4029: Kubernetes building blocks (core API primitives)
- `KubernetesDeployment`
- `KubernetesStatefulSet`
- `KubernetesDaemonSet`
- `KubernetesJob`
- `KubernetesCronJob`
- `KubernetesService`
- `KubernetesSecret`
- `KubernetesManifest`
- `KubernetesHelmRelease`
- `KubernetesConfigMap`
- `KubernetesServiceAccount`
- `KubernetesRbac` -- Bundles the RBAC grant grain (Role/ClusterRole + its binding) into one component: "grant these permissions to these subjects in this scope".
- `KubernetesIngress`
- `KubernetesNetworkPolicy`
- `KubernetesPersistentVolumeClaim`
- `KubernetesStorageClass`
- `KubernetesResourceQuota` -- Manages the namespace-governance pair: the ResourceQuota plus an optional companion LimitRange (per-object defaults/bounds) — two API objects, one governance story.
- `KubernetesPriorityClass`
- `KubernetesPodDisruptionBudget`
- `KubernetesHorizontalPodAutoscaler`
- `KubernetesCertManager` -- 4030–4069: Kubernetes foundation addons (certs, DNS, secrets, ingress, Gateway API, mesh, CNI/autoscaling/DR)
- `KubernetesClusterIssuer` -- KubernetesCertManager is a prerequisite for the three cert-manager CR kinds below: ClusterIssuer/Issuer/Certificate are cert-manager custom resources — without the controller and its CRDs they cannot be applied.
- `KubernetesIssuer`
- `KubernetesCertificate`
- `KubernetesExternalDns`
- `KubernetesExternalSecretsOperator`
- `KubernetesClusterSecretStore` -- KubernetesExternalSecretsOperator is a prerequisite for the three external-secrets CR kinds below: ClusterSecretStore/SecretStore/ ExternalSecret are external-secrets custom resources — without the operator and its CRDs they cannot be applied.
- `KubernetesSecretStore`
- `KubernetesExternalSecret`
- `KubernetesIngressNginx`
- `KubernetesGatewayApiCrds`
- `KubernetesGatewayClass`
- `KubernetesGateway`
- `KubernetesListenerSet`
- `KubernetesHttpRoute`
- `KubernetesGrpcRoute`
- `KubernetesTcpRoute`
- `KubernetesUdpRoute`
- `KubernetesTlsRoute`
- `KubernetesReferenceGrant`
- `KubernetesBackendTlsPolicy`
- `KubernetesIstioBaseCrds`
- `KubernetesIstio`
- `KubernetesDestinationRule` -- Istio API components (mesh traffic policy, security, telemetry). The seven typed resources below (4053–4059) require the Istio CRDs on the cluster, provided by the lightweight CRDs-only KubernetesIstioBaseCrds (851) — NOT the full mesh KubernetesIstio (852).
- `KubernetesServiceEntry`
- `KubernetesPeerAuthentication`
- `KubernetesRequestAuthentication`
- `KubernetesAuthorizationPolicy`
- `KubernetesTelemetry`
- `KubernetesEnvoyFilter`
- `KubernetesMetricsServer`
- `KubernetesCilium`
- `KubernetesKeda`
- `KubernetesKarpenter`
- `KubernetesKarpenterNodePool`
- `KubernetesKarpenterEc2NodeClass`
- `KubernetesClusterAutoscaler`
- `KubernetesVelero`
- `KubernetesKubePrometheusStack` -- 4070–4089: Kubernetes observability
- `KubernetesGrafana`
- `KubernetesSignoz` -- KubernetesClickHouse is a prerequisite because SigNoz stores every trace, metric and log in ClickHouse and deploys none of its own — the telemetry store is composed, never bundled.
- `KubernetesLoki`
- `KubernetesTempo`
- `KubernetesOtelOperator` -- The operator's admission webhooks (failurePolicy Fail) are served with a cert-manager Certificate in the default posture — cert-manager must be running before the operator installs.
- `KubernetesOtelCollector`
- `KubernetesKyverno` -- 4080–4099: Kubernetes security, policy, and identity
- `KubernetesGatekeeper`
- `KubernetesKeycloak` -- Keycloak declarations compose the official Keycloak Operator (which reconciles the Keycloak CR this kind renders) and, on the recommended postgres vendor, a KubernetesPostgres database — both must resolve before the CR can converge.
- `KubernetesOpenBao`
- `KubernetesOpenFga` -- OpenFGA requires a datastore; the recommended arm composes a KubernetesPostgres database (the sandbox memory arm needs nothing, but the registry declares the shape real deployments require).
- `KubernetesKeycloakOperator`
- `KubernetesCloudNativePgOperator` -- 4100–4129: Kubernetes data platforms
- `KubernetesPostgres`
- `KubernetesValkey`
- `KubernetesPerconaMysqlOperator`
- `KubernetesMysql`
- `KubernetesPerconaMongoOperator`
- `KubernetesMongodb`
- `KubernetesStrimziKafkaOperator`
- `KubernetesKafka` -- container_kind: a Strimzi Kafka cluster is a place in the provider's own model — KafkaTopic and KafkaUser declarations BELONG to one cluster (the strimzi.io/cluster label) and are drawn inside its box. Clients that merely talk to the cluster (Connect, MirrorMaker2, UI, Karapace) carry containment_exempt on their bootstrap/trust references.
- `KubernetesKafkaTopic`
- `KubernetesKafkaUser`
- `KubernetesKafkaConnect` -- container_kind: a Connect cluster hosts the connectors deployed INTO it (KafkaConnector's strimzi.io/cluster label names its Connect cluster) — the same room shape as KubernetesKafka above.
- `KubernetesKafkaConnector`
- `KubernetesKafkaMirrorMaker2`
- `KubernetesKarapace`
- `KubernetesKafkaUi`
- `KubernetesOpenSearchOperator`
- `KubernetesOpenSearch`
- `KubernetesAltinityOperator`
- `KubernetesClickHouse`
- `KubernetesSolrOperator`
- `KubernetesSolr`
- `KubernetesNeo4j`
- `KubernetesSeaweedFs`
- `KubernetesQdrant`
- `KubernetesRabbitMqOperator` -- The RabbitMQ Cluster Operator's release manifest ships admission webhooks whose serving certificate is a cert-manager Certificate — cert-manager must be running before the operator installs.
- `KubernetesRabbitMq`
- `KubernetesAirflow` -- 4130–4149: Kubernetes analytics and ML KubernetesPostgres is a prerequisite because Airflow's metadata database composes a KubernetesPostgres by default (the spec's FK defaults resolve onto its outputs) and the migration Job needs the database reachable before the server components start.
- `KubernetesSparkOperator`
- `KubernetesKubeRayOperator`
- `KubernetesRayCluster` -- KubernetesKubeRayOperator is a prerequisite because this kind declares the RayCluster custom resource that only the operator's CRDs admit and only the operator reconciles into head and worker pods.
- `KubernetesFlinkOperator` -- KubernetesCertManager is a prerequisite because the Flink operator's chart, with its default-on admission webhook, renders cert-manager Issuer/Certificate resources and trusts the API server through cert-manager's CA injection — there is no self-signed fallback at the pinned chart, and the webhooks are fail-closed.
- `KubernetesFlinkDeployment` -- KubernetesFlinkOperator is a prerequisite because this kind declares the FlinkDeployment custom resource that only the operator's CRDs admit and only the operator reconciles into a running Flink cluster.
- `KubernetesJupyterHub` -- KubernetesPostgres is a prerequisite because JupyterHub's hub database composes a KubernetesPostgres in its external-database arm (the spec's FK defaults resolve onto its outputs) and the hub pod mounts that database's credential Secret before it can start.
- `KubernetesMlflow` -- KubernetesPostgres is a prerequisite because MLflow's backend store composes a KubernetesPostgres in its production arm (FK defaults onto its outputs; the module composes the connection URI from its credential Secret), and KubernetesSeaweedFs because the artifact store's S3-compatible arm FK-defaults onto the SeaweedFS endpoint and credential Secret.
- `KubernetesTrino` -- KubernetesPostgres is a prerequisite because Trino's postgres catalogs compose a KubernetesPostgres (the catalog host and credential FK-default onto its outputs), and the pods read that database's credential Secret to resolve catalog passwords from environment.
- `KubernetesSuperset` -- KubernetesPostgres is a prerequisite because Superset's REQUIRED metadata database composes a KubernetesPostgres (FK defaults onto its outputs; the module composes the environment Secret from its credential Secret), and KubernetesValkey because the cache/broker arm FK-defaults onto a KubernetesValkey's service and password Secret.
- `KubernetesArgocd` -- 4150–4169: Kubernetes GitOps and CI/CD
- `KubernetesArgoWorkflows`
- `KubernetesTektonOperator`
- `KubernetesTekton` -- KubernetesTektonOperator is a prerequisite because this kind declares the TektonConfig custom resource that only the operator's CRDs admit and only the operator reconciles into running components.
- `KubernetesGhaRunnerScaleSetController`
- `KubernetesGhaRunnerScaleSet` -- KubernetesGhaRunnerScaleSetController is a prerequisite because this kind renders an AutoscalingRunnerSet custom resource that only the controller's CRDs admit and only the controller reconciles into listener and runner pods.
- `KubernetesHarbor`
- `KubernetesJenkins`
- `KubernetesTemporal` -- 4170–4189: Kubernetes app platforms KubernetesPostgres is a prerequisite because the recommended (and E2E-proven) database composition backs Temporal's default and visibility stores with a CloudNativePG cluster.
- `KubernetesNats`
- `KubernetesLocust`
- `DigitalOceanAppPlatformService` -- 5000–5999: DigitalOcean resources
- `DigitalOceanBucket`
- `DigitalOceanContainerRegistry`
- `DigitalOceanDatabaseCluster`
- `DigitalOceanDnsZone`
- `DigitalOceanDroplet`
- `DigitalOceanFirewall`
- `DigitalOceanFunction`
- `DigitalOceanKubernetesCluster`
- `DigitalOceanKubernetesNodePool`
- `DigitalOceanLoadBalancer`
- `DigitalOceanVolume`
- `DigitalOceanVpc`
- `DigitalOceanCertificate`
- `DigitalOceanDnsRecord`
- `CivoBucket` -- 6000–6999: Civo resources
- `CivoCertificate`
- `CivoComputeInstance`
- `CivoDatabase`
- `CivoDnsZone`
- `CivoFirewall`
- `CivoIpAddress`
- `CivoKubernetesCluster`
- `CivoKubernetesNodePool`
- `CivoVolume`
- `CivoVpc`
- `CivoDnsRecord`
- `CloudflareDnsZone` -- 7000–7999: Cloudflare resources
- `CloudflareKvNamespace`
- `CloudflareR2Bucket`
- `CloudflareWorker`
- `CloudflareLoadBalancer`
- `CloudflareD1Database`
- `CloudflareZeroTrustAccessApplication`
- `CloudflareDnsRecord`
- `CloudflareRuleset`
- `CloudflareWorkersKvPair`
- `CloudflareHyperdriveConfig`
- `CloudflareLoadBalancerPool`
- `CloudflareLoadBalancerMonitor`
- `CloudflareZeroTrustAccessPolicy`
- `CloudflareZeroTrustAccessGroup`
- `CloudflareQueue`
- `CloudflarePagesProject`
- `CloudflareZeroTrustTunnel`
- `CloudflareZeroTrustTunnelVirtualNetwork`
- `CloudflareZeroTrustTunnelRoute`
- `CloudflareList`
- `CloudflareListItem`
- `CloudflareTurnstileWidget`
- `CloudflareEmailRoutingZone`
- `CloudflareEmailRoutingRule`
- `CloudflareEmailRoutingAddress`
- `CloudflareOriginCaCertificate`
- `CloudflareCertificatePack`
- `CloudflareCustomHostname`
- `CloudflareCustomHostnameFallbackOrigin`
- `Auth0Connection` -- 8000–8999: Auth0 resources
- `Auth0Client`
- `Auth0EventStream`
- `Auth0ResourceServer`
- `Auth0Action`
- `Auth0Role`
- `OpenFgaStore` -- 9000–9999: OpenFGA resources Note: OpenFGA is Terraform-only - there is no Pulumi provider available. Pulumi modules for OpenFGA resources are pass-through placeholders.
- `OpenFgaAuthorizationModel`
- `OpenFgaRelationshipTuple`
- `OpenStackKeypair` -- 10000–10999: OpenStack resources
- `OpenStackNetwork`
- `OpenStackSubnet`
- `OpenStackRouter`
- `OpenStackRouterInterface`
- `OpenStackSecurityGroup`
- `OpenStackFloatingIp`
- `OpenStackNetworkPort`
- `OpenStackSecurityGroupRule`
- `OpenStackFloatingIpAssociate`
- `OpenStackInstance`
- `OpenStackServerGroup`
- `OpenStackVolume`
- `OpenStackVolumeAttach`
- `OpenStackProject`
- `OpenStackApplicationCredential`
- `OpenStackImage`
- `OpenStackRoleAssignment`
- `OpenStackLoadBalancer`
- `OpenStackLoadBalancerListener`
- `OpenStackLoadBalancerPool`
- `OpenStackLoadBalancerMember`
- `OpenStackLoadBalancerMonitor`
- `OpenStackDnsZone`
- `OpenStackDnsRecord`
- `ScalewayVpc`
- `ScalewayPrivateNetwork`
- `ScalewayPublicGateway`
- `ScalewayLoadBalancer`
- `ScalewayInstanceSecurityGroup`
- `ScalewayInstance`
- `ScalewayKapsuleCluster`
- `ScalewayKapsulePool`
- `ScalewayRdbInstance`
- `ScalewayRedisCluster`
- `ScalewayMongodbInstance`
- `ScalewayObjectBucket`
- `ScalewayBlockVolume`
- `ScalewayContainerRegistry`
- `ScalewayDnsZone`
- `ScalewayDnsRecord`
- `ScalewayServerlessFunction`
- `ScalewayServerlessContainer`
- `AliCloudLogProject`
- `AliCloudRamRole`
- `AliCloudRamPolicy`
- `AliCloudVpc`
- `AliCloudVswitch`
- `AliCloudSecurityGroup`
- `AliCloudEipAddress`
- `AliCloudNatGateway`
- `AliCloudApplicationLoadBalancer`
- `AliCloudNetworkLoadBalancer`
- `AliCloudVpnGateway`
- `AliCloudDnsZone`
- `AliCloudDnsRecord`
- `AliCloudPrivateDnsZone`
- `AliCloudStorageBucket`
- `AliCloudNasFileSystem`
- `AliCloudKmsKey`
- `AliCloudRdsInstance`
- `AliCloudPolardbCluster`
- `AliCloudRedisInstance`
- `AliCloudMongodbInstance`
- `AliCloudEcsInstance`
- `AliCloudContainerRegistry`
- `AliCloudKubernetesCluster`
- `AliCloudKubernetesNodePool`
- `AliCloudCdnDomain`
- `AliCloudFunction`
- `AliCloudSaeApplication`
- `AliCloudRocketmqInstance`
- `AliCloudCenInstance`
- `OciVcn`
- `OciSubnet`
- `OciSecurityGroup`
- `OciCompartment`
- `OciIdentityPolicy`
- `OciDynamicGroup`
- `OciComputeInstance`
- `OciContainerEngineCluster`
- `OciContainerEngineNodePool`
- `OciContainerInstance`
- `OciApplicationLoadBalancer`
- `OciNetworkLoadBalancer`
- `OciDynamicRoutingGateway`
- `OciPublicIp`
- `OciAutonomousDatabase`
- `OciDbSystem`
- `OciMysqlDbSystem`
- `OciPostgresqlDbSystem`
- `OciRedisCluster`
- `OciNosqlTable`
- `OciObjectStorageBucket`
- `OciFileSystem`
- `OciBlockVolume`
- `OciKmsVault`
- `OciKmsKey`
- `OciVaultSecret`
- `OciBastion`
- `OciFunctionsApplication`
- `OciApiGateway`
- `OciStreamPool`
- `OciQueue`
- `OciAlarm`
- `OciLogGroup`
- `OciDnsZone`
- `OciDnsRecord`
- `OciNetworkFirewall`
- `OciDevopsProject`
- `HetznerCloudSshKey`
- `HetznerCloudPlacementGroup`
- `HetznerCloudFirewall`
- `HetznerCloudNetwork`
- `HetznerCloudPrimaryIp`
- `HetznerCloudFloatingIp`
- `HetznerCloudServer`
- `HetznerCloudVolume`
- `HetznerCloudSnapshot`
- `HetznerCloudCertificate`
- `HetznerCloudLoadBalancer`
- `HetznerCloudDnsZone`

### spec.pod.initContainers[].env.variables[].valueFrom.env

`string`

### spec.pod.initContainers[].env.variables[].valueFrom.name

`string` · required

- rule: {"required":true}

### spec.pod.initContainers[].env.variables[].valueFrom.fieldPath

`string`

### spec.pod.initContainers[].env.variables[].configMapKeyRef

`ConfigMapKeyRef`

Reference to a key in a Kubernetes ConfigMap.

### spec.pod.initContainers[].env.variables[].configMapKeyRef.name

`string` · required

Name of the ConfigMap.

- rule: {"required":true}

### spec.pod.initContainers[].env.variables[].configMapKeyRef.key

`string` · required

Key within the ConfigMap.

- rule: {"required":true}

### spec.pod.initContainers[].env.variables[].configMapKeyRef.optional

`bool`

If true, the env var is silently skipped when the ConfigMap or key does not exist
(instead of blocking pod startup).

### spec.pod.initContainers[].env.variables[].fieldRef

`ObjectFieldRef`

Reference to a pod-level field (metadata.name, status.podIP, etc.).

### spec.pod.initContainers[].env.variables[].fieldRef.apiVersion

`string`

Version of the schema. Defaults to "v1".

### spec.pod.initContainers[].env.variables[].fieldRef.fieldPath

`string` · required

Path of the field to select (e.g., "metadata.name", "status.podIP").

- rule: {"required":true}

### spec.pod.initContainers[].env.variables[].resourceFieldRef

`ResourceFieldRef`

Reference to container resource limits or requests (limits.cpu, requests.memory, etc.).

### spec.pod.initContainers[].env.variables[].resourceFieldRef.containerName

`string`

Container name. Required for init containers; defaults to the current
container for regular containers.

### spec.pod.initContainers[].env.variables[].resourceFieldRef.resource

`string` · required

Resource to select (e.g., "limits.cpu", "requests.memory").

- rule: {"required":true}

### spec.pod.initContainers[].env.variables[].resourceFieldRef.divisor

`string`

Specifies the output format of the exposed resource.
For CPU: "1" means cores. For memory: "1", "1Ki", "1Mi", "1Gi".

### spec.pod.initContainers[].env.secrets

`[]SecretEnvVar`

Individual secret environment variables (sensitive).

### spec.pod.initContainers[].env.secrets[].name

`string` · required

The environment variable name.

- rule: Must be a valid C_IDENTIFIER (start with letter/underscore, contain only letters, digits, underscores)
- rule: {"required":true}

### spec.pod.initContainers[].env.secrets[].value

`string`

Literal string value.
A Kubernetes Secret is automatically created and the environment variable
references that secret.

### spec.pod.initContainers[].env.secrets[].secretRef

`KubernetesSecretKeyRef`

Reference to a key within an existing Kubernetes Secret.

### spec.pod.initContainers[].env.secrets[].secretRef.namespace

`string`

The namespace of the Kubernetes Secret.
If not specified, defaults to the namespace where the component is deployed.
Note: Cross-namespace secret references may not be supported by all Helm charts.

### spec.pod.initContainers[].env.secrets[].secretRef.name

`string` · required

The name of the Kubernetes Secret.

- rule: {"required":true}

### spec.pod.initContainers[].env.secrets[].secretRef.key

`string` · required

The key within the Kubernetes Secret that contains the value.

- rule: {"required":true}

### spec.pod.initContainers[].env.secrets[].secretRef.optional

`bool`

If true, the env var is silently skipped when the Secret or key does not exist
(instead of blocking pod startup).

### spec.pod.initContainers[].env.secrets[].valueFrom

`ValueFromRef`

Reference to another Planton resource's secret output field.
The orchestrator resolves this before invoking IaC modules.

### spec.pod.initContainers[].env.secrets[].valueFrom.kind

`enum`

Allowed values (use exactly as shown):

- `unspecified` -- 0: Default/unspecified
- `TestCloudResourceGeneric` -- 1–49: Test/dev/custom
- `TestCloudResourceKubernetes`
- `ConfluentKafka` -- 50–199: saas platform resources
- `AtlasMongodb`
- `SnowflakeDatabase`
- `AwsAlb` -- 1000–1999: AWS resources AwsSubnet is a prerequisite because an ALB requires at least two subnets in different availability zones -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsCertManagerCert`
- `AwsCloudFront`
- `AwsDynamodb`
- `AwsEcrRepo`
- `AwsEcsCluster`
- `AwsEcsService` -- AwsEcsCluster, AwsEcsTaskDefinition, and AwsSubnet are prerequisites because a service schedules a referenced task-definition revision into a referenced live cluster and places task network interfaces into referenced subnets -- all three references must resolve first.
- `AwsEksCluster` -- AwsSubnet and AwsIamRole are prerequisites because the control plane attaches its network interfaces into referenced subnets and assumes a referenced cluster role that must already carry AmazonEKSClusterPolicy.
- `AwsIamRole`
- `AwsLambda`
- `AwsRdsCluster`
- `AwsRdsInstance`
- `AwsRoute53Zone`
- `AwsS3Bucket`
- `AwsLbTargetGroup` -- AwsVpc is a prerequisite because a target group's health checks and target registrations live inside one VPC -- the spec's vpc_id reference must resolve before the group can be created.
- `AwsSecurityGroup` -- AwsVpc is a prerequisite because every security group is created in a VPC; the E2E install profile resolves vpc_id against the VPC prerequisite.
- `AwsVpc`
- `AwsEksNodeGroup` -- AwsEksCluster is a prerequisite because nodes register with a live control plane; AwsIamRole and AwsSubnet back the node role and worker subnet references.
- `AwsIamUser`
- `AwsKmsKey`
- `AwsEc2Instance`
- `AwsClientVpn` -- Every Client VPN endpoint requires an ACM server certificate at create time; the imported self-signed fixture satisfies it. Subnets/VPC are optional composition (a zero-association endpoint is valid) -- composed scenarios declare them via the e2e-prerequisites annotation.
- `AwsDocumentDb`
- `AwsRoute53DnsRecord` -- AwsRoute53Zone is a prerequisite because every record lives inside a hosted zone -- the spec's zone_id reference must resolve before the record can be created.
- `AwsS3ObjectSet` -- AwsS3Bucket is a prerequisite because the object set's bucket reference is required -- objects cannot exist without the bucket that holds them.
- `AwsSqsQueue`
- `AwsSnsTopic`
- `AwsEventBridgeBus`
- `AwsEventBridgeRule`
- `AwsIamOidcProvider`
- `AwsIamPolicy`
- `AwsIamInstanceProfile` -- AwsIamRole is a prerequisite because an instance profile is a wrapper that must contain a role to be useful -- the profile's spec requires a role reference, so the role must be deployed first.
- `AwsLbListener` -- AwsAlb and AwsLbTargetGroup are prerequisites because a listener is an attachment point on a load balancer and its default action almost always forwards to a target group -- both references must resolve before the listener can be created.
- `AwsLbListenerRule` -- AwsLbListener is a prerequisite because a rule only exists as an attachment on a listener -- the listener_arn reference must resolve before the rule can be created.
- `AwsLaunchTemplate`
- `AwsAutoScalingGroup` -- AwsSubnet and AwsLaunchTemplate are prerequisites because a group cannot exist without subnets to place capacity in and a launch template to launch from -- the spec's subnets and launch_template references must resolve before the group can be created.
- `AwsEksAddon` -- AwsEksCluster is a prerequisite because an add-on installs onto a live control plane -- the spec's cluster_name reference must resolve before the add-on can be created.
- `AwsEksFargateProfile` -- AwsEksCluster, AwsIamRole, and AwsSubnet are prerequisites because a Fargate profile attaches to a live control plane, runs pods as a referenced pod-execution role, and launches them into referenced private subnets -- all three references must resolve first.
- `AwsEksAccessEntry` -- AwsEksCluster and AwsIamRole are prerequisites because an access entry grants a referenced IAM principal access to a live control plane -- both references must resolve before the entry can be created.
- `AwsEcsTaskDefinition` -- AwsIamRole is a prerequisite because the kind's default posture -- Fargate with the awslogs logging default -- is rejected by AWS at registration time without an execution role the agent can assume.
- `AwsHttpApiGateway`
- `AwsStepFunction` -- AwsIamRole is a prerequisite because a state machine cannot be created without an execution role it can assume -- the spec's role_arn reference must resolve before the CreateStateMachine call.
- `AwsHttpApiVpcLink` -- AwsSubnet is a prerequisite because a VPC link is a set of managed ENIs provisioned into referenced subnets -- the subnet references must resolve before the link can be created. Security groups are optional on the link, so they compose per-scenario rather than as a registry prerequisite.
- `AwsHttpApiDomain` -- AwsCertManagerCert is a prerequisite because a custom domain cannot be created without a TLS certificate in the same region covering the domain -- the spec's certificate_arn reference must resolve first.
- `AwsVpcEndpoint` -- AwsVpcEndpoint's composed E2E scenarios reference the AwsVpc prerequisite's outputs (vpc_id + default_route_table_id for gateway endpoints) and the AwsSubnet pair's subnet_id outputs (interface endpoints), so both are genuine deploy-order prerequisites.
- `AwsElasticacheUser`
- `AwsElasticacheUserGroup` -- AwsElasticacheUser is a genuine prerequisite: AWS refuses to create a user group that does not contain a user named "default", so a group's composed E2E scenario must resolve a deployed user's outputs.
- `AwsRedshiftServerlessNamespace`
- `AwsRedshiftServerlessWorkgroup` -- The namespace is a genuine prerequisite: a workgroup attaches to exactly one namespace by name at create time, so its composed E2E scenario must resolve a deployed namespace's outputs. AwsSubnet is a prerequisite because Redshift Serverless requires the workgroup's subnets to span three availability zones.
- `AwsRedisElasticache` -- AwsSubnet is a prerequisite because the module builds an ElastiCache subnet group from referenced subnets -- the spec's subnet references must resolve before the replication group can deploy.
- `AwsOpenSearchDomain`
- `AwsMemcachedElasticache`
- `AwsServerlessElasticache`
- `AwsNlb` -- AwsSubnet is a prerequisite because an NLB requires at least one subnet mapping -- the spec's subnet references must resolve before the load balancer can be created.
- `AwsElasticIp`
- `AwsTransitGateway`
- `AwsGlobalAccelerator`
- `AwsSubnet`
- `AwsInternetGateway`
- `AwsNatGateway` -- AwsInternetGateway is a prerequisite because a public NAT gateway can only become available once the VPC it sits in has an internet gateway attached (AWS rejects the create otherwise) -- so the gateway must be deployed first. AwsVpc is a prerequisite because a REGIONAL NAT gateway (availability_mode = regional) references the VPC directly instead of a subnet.
- `AwsEgressOnlyInternetGateway`
- `AwsElasticFileSystem` -- AwsSubnet and AwsSecurityGroup are prerequisites because mount targets (required, min 1) place the file system's NFS endpoints into subnets and attach security groups -- both references must resolve before the CreateMountTarget calls.
- `AwsEfsAccessPoint` -- AwsElasticFileSystem is a prerequisite because an access point is created INTO a file system -- the spec's required file_system_id reference must resolve before the CreateAccessPoint call.
- `AwsFsxLustreFileSystem`
- `AwsFsxOpenzfsFileSystem`
- `AwsFsxWindowsFileSystem` -- Every Windows file system must join an Active Directory domain; the directory itself is external infrastructure (AWS Managed Microsoft AD or a self-managed domain), so only the network dependency is a declarable prerequisite.
- `AwsFsxOntapFileSystem`
- `AwsFsxOntapStorageVirtualMachine`
- `AwsFsxOntapVolume`
- `AwsFsxDataRepositoryAssociation`
- `AwsCognitoUserPool`
- `AwsCognitoIdentityProvider` -- AwsCognitoUserPool is a prerequisite because an identity provider is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateIdentityProvider call.
- `AwsCognitoUserPoolClient` -- AwsCognitoUserPool is a prerequisite because an app client is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateUserPoolClient call.
- `AwsCognitoResourceServer` -- AwsCognitoUserPool is a prerequisite because a resource server is created INTO a pool -- the spec's required user_pool_id reference must resolve before the CreateResourceServer call.
- `AwsWafWebAcl`
- `AwsWafIpSet`
- `AwsWafRegexPatternSet`
- `AwsCloudwatchLogGroup`
- `AwsCloudwatchAlarm`
- `AwsCloudwatchCompositeAlarm`
- `AwsKinesisStream`
- `AwsKinesisFirehose` -- Every Firehose destination requires an S3 configuration (the primary target for extended_s3; the failed/all-document backup for the rest) and an IAM role Firehose assumes to write to it, so both are hard deploy prerequisites.
- `AwsKinesisStreamConsumer` -- A consumer registers against exactly one stream and cannot exist without it.
- `AwsAthenaWorkgroup`
- `AwsGlueCatalogDatabase`
- `AwsRedshiftCluster`
- `AwsSagemakerDomain` -- AI/ML A domain cannot exist without VPC subnets and a SageMaker execution role (default_user_settings.execution_role_arn is required), so both are hard deploy prerequisites.
- `AwsAppRunnerService` -- A service can run entirely on companion defaults, so the App Runner family's kinds are dependency-free leaves except the VPC connector (which cannot exist without subnets and security groups). A service's companion references (auto scaling / VPC connector / observability / WAF) are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsAppRunnerAutoScalingConfiguration`
- `AwsAppRunnerVpcConnector`
- `AwsAppRunnerObservabilityConfiguration`
- `AwsTransitGatewayVpcAttachment` -- AwsTransitGateway is a prerequisite because an attachment cannot exist without the gateway it attaches to; AwsSubnet because the attachment provisions an ENI into at least one subnet (the VPC arrives transitively through the subnet's own prerequisites).
- `AwsTransitGatewayRouteTable` -- Only the gateway is a hard prerequisite: a route table can exist empty. Associations, propagations, and routes referencing attachments are optional composition -- scenarios declare them via the e2e-prerequisites annotation.
- `AwsBatchComputeEnvironment` -- A MANAGED compute environment always launches into VPC subnets, so the subnet is a hard deploy prerequisite (security groups are required only for the Fargate types -- scenario-declared, not a registry edge).
- `AwsBatchJobQueue` -- A job queue cannot exist without at least one VALID compute environment to map onto.
- `AwsBatchSchedulingPolicy`
- `AwsBatchJobDefinition`
- `AwsCodeBuildProject` -- CI/CD
- `AwsCodePipeline`
- `AwsMwaaEnvironment` -- Workflow / Orchestration AwsSubnet and AwsSecurityGroup are prerequisites because the environment's network interfaces are placed in referenced private subnets and AWS requires at least one attached security group at creation.
- `AwsNeptuneCluster` -- Graph Database
- `AwsMemorydbCluster` -- A cluster always launches into a subnet group; the subnets are the hard deploy prerequisite. The ACL it attaches is optional composition (the built-in "open-access" ACL needs no resource) -- scenarios declare the ACL/user chain via the e2e-prerequisites annotation.
- `AwsMemorydbUser`
- `AwsMemorydbAcl` -- An empty ACL is valid (MemoryDB has no mandatory "default" member), so the user is optional composition -- the composed scenario declares it via the e2e-prerequisites annotation, never a registry edge.
- `AwsMskCluster` -- Streaming AwsSubnet and AwsSecurityGroup are prerequisites because brokers are placed in referenced subnets and AWS requires at least one attached security group at creation.
- `AwsMskServerlessCluster` -- AwsSubnet is a prerequisite because the serverless cluster's network interfaces are placed in referenced subnets (security groups are optional -- AWS attaches the VPC default group when none are referenced).
- `AwsLambdaEventSourceMapping` -- AwsLambda is a prerequisite because a mapping cannot exist without the function it invokes (a required reference). Event sources (SQS, Kinesis, DynamoDB, MSK) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsSnsSubscription` -- AwsSnsTopic is a prerequisite because a subscription cannot exist without the topic it subscribes to (a required reference). Endpoints (SQS queues, Lambda functions, Firehose streams) are optional composition -- scenarios declare them via the e2e-prerequisites annotation rather than taxing every consumer's chain.
- `AwsPlantonRunner` -- AwsSubnet is a prerequisite because the runner appliance places its network interfaces into referenced subnets -- the placement reference must resolve before the appliance can deploy.
- `AwsRoute53HealthCheck`
- `AwsSesConfigurationSet` -- Both SES kinds are dependency-free leaves: an identity's configuration set is optional composition (scenarios declare it via the e2e-prerequisites annotation), and a configuration set's event destinations reference other kinds only optionally.
- `AwsSesEmailIdentity`
- `AwsSecretsManagerSecret` -- A dependency-free leaf: the KMS key, rotation Lambda, and external rotation role references are all optional composition -- scenarios declare them via the e2e-prerequisites annotation, never registry edges.
- `AwsOpenSearchServerlessCollection` -- A dependency-free leaf: the collection-scoped encryption/network/ data-access/retention policies are module-rendered, and the KMS key and data-access principal references are optional composition (e2e-prerequisites annotation).
- `AwsBedrockGuardrail` -- A dependency-free leaf: the KMS key reference is optional composition (e2e-prerequisites annotation); published versions are folded satellites of the guardrail itself.
- `AwsBedrockCustomModel` -- AwsIamRole is a prerequisite because Bedrock assumes the job role to read training data and write outputs; the S3 locations and KMS key are optional composition (e2e-prerequisites annotation).
- `AwsBedrockInferenceProfile` -- A dependency-free leaf: the model source is a foundation model or an AWS system-defined cross-region profile, never a customer resource.
- `AwsBedrockProvisionedThroughput` -- A dependency-free leaf in the registry: capacity is typically bought for an AwsBedrockCustomModel (the default reference), but foundation model ARNs are equally legal, so the edge is optional composition.
- `AwsBedrockModelAccess` -- A dependency-free leaf: the agreement covers an AWS-listed foundation model, never a customer resource.
- `AwsBedrockAgent` -- AwsIamRole is a prerequisite because the Bedrock service assumes the agent resource role to invoke models, action-group Lambdas, and knowledge bases; the guardrail, KMS key, provisioned throughput, and collaborator/knowledge-base edges are optional composition (e2e-prerequisites annotation). Action groups, aliases, collaborators, and knowledge-base associations are folded satellites of the agent.
- `AwsBedrockKnowledgeBase` -- AwsIamRole is a prerequisite because the Bedrock service assumes the knowledge-base role to read data sources, call the embedding model, and read/write the vector store; the vector-store and data-source reference edges (OpenSearch, S3, Secrets Manager, ...) are optional composition (e2e-prerequisites annotation). Data sources are folded satellites of the knowledge base.
- `AwsBedrockFlow` -- AwsIamRole is a prerequisite because the Bedrock service assumes the flow execution role to invoke the models, agents, knowledge bases, and Lambdas its nodes reference; the node-level reference edges are optional composition (e2e-prerequisites annotation).
- `AwsBedrockPrompt` -- A dependency-free leaf: variants target AWS-listed foundation models by ID; targeting another agent's alias is optional composition (e2e-prerequisites annotation).
- `AzureResourceGroup` -- 2000–2999: Azure resources
- `AzureAksCluster` -- AzureResourceGroup is the only required parent: the cluster is created inside a referenced resource group. Subnet is optional on the default node pool (AKS provisions managed networking when unset).
- `AzureAksNodePool` -- AzureAksCluster is a prerequisite because a node pool attaches to an existing cluster by ARM ID; the resource group chains transitively.
- `AzureContainerRegistry` -- AzureResourceGroup is a prerequisite because a container registry is created inside a resource group.
- `AzureDnsZone` -- AzureResourceGroup is a prerequisite because the DNS zone is created inside a referenced resource group that must already exist.
- `AzureKeyVault` -- AzureResourceGroup is a prerequisite because a key vault is created inside a referenced resource group in composed environments.
- `AzureVirtualNetwork` -- AzureResourceGroup is a prerequisite because a virtual network is created inside a referenced resource group in composed environments.
- `AzureNatGateway` -- AzureResourceGroup is a prerequisite because a NAT gateway is created inside a referenced resource group in composed environments.
- `AzureVirtualMachine` -- AzureNetworkInterface is a prerequisite because a virtual machine attaches at least one NIC (the subnet, network, and resource group chain transitively through the NIC's own prerequisites).
- `AzureStorageAccount` -- AzureResourceGroup is a prerequisite because a storage account is created inside a referenced resource group in composed environments.
- `AzureDnsRecord` -- AzureDnsZone is a prerequisite because a record set is created inside a referenced zone (the resource group chains transitively through the zone). Public DNS zone names are not globally unique, so a shared zone fixture is safe to recreate across scenarios.
- `AzureSubnet` -- AzureVirtualNetwork is a prerequisite because a subnet is an ARM child of a referenced network -- the network must exist before the subnet can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzureNetworkSecurityGroup` -- AzureResourceGroup is a prerequisite because a network security group is created inside a referenced resource group in composed environments.
- `AzurePublicIp` -- AzureResourceGroup is a prerequisite because a public IP is created inside a referenced resource group in composed environments.
- `AzurePrivateEndpoint` -- AzureSubnet is a prerequisite because a private endpoint draws its private IP from a referenced subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite). The connection target is polymorphic and the DNS zones / ASGs are optional, so none of those are prerequisites.
- `AzurePrivateDnsZone` -- AzureResourceGroup is a prerequisite because a private DNS zone is created inside a referenced resource group in composed environments.
- `AzureApplicationGateway` -- AzureSubnet is a prerequisite because a gateway cannot exist without its dedicated gateway_ip_configuration subnet (the network and resource group chain transitively through the subnet's own prerequisites); public frontends additionally reference a public IP, but private-only gateways are legal, so it is not a registry prerequisite.
- `AzureLoadBalancer` -- AzureResourceGroup is a prerequisite because a load balancer is created inside a referenced resource group (frontends additionally reference subnets or public IPs, but neither is universally required, so they are not registry prerequisites).
- `AzureRouteTable` -- AzureResourceGroup is a prerequisite because a route table is created inside a referenced resource group in composed environments.
- `AzurePrivateDnsZoneVirtualNetworkLink` -- AzurePrivateDnsZone and AzureVirtualNetwork are prerequisites because a virtual network link is a child resource of a referenced zone and binds it to a referenced network -- both must exist before the link can be written. (The resource group arrives transitively through the zone's and network's own prerequisite declarations.)
- `AzureVirtualNetworkPeering` -- AzureVirtualNetwork is a prerequisite because a peering is an ARM child of its local network and binds it to a remote network -- the local network must exist before the peering can be written. (The resource group arrives transitively through the network's own prerequisite declaration.)
- `AzurePublicIpPrefix` -- AzureResourceGroup is a prerequisite because a public IP prefix is created inside a referenced resource group in composed environments.
- `AzureNetworkInterface` -- AzureSubnet is a prerequisite because a network interface's IP configurations deploy into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureManagedDisk` -- AzureResourceGroup is a prerequisite because a managed disk is created inside a resource group.
- `AzureVirtualMachineScaleSet` -- AzureSubnet is a prerequisite because every scale-set instance's network interface deploys into a subnet (the virtual network and resource group chain transitively through the subnet's own prerequisite).
- `AzureKeyVaultKey` -- AzureKeyVault is a prerequisite because a key is a data-plane object inside a referenced vault -- the vault must exist before the key can be written (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultCertificate` -- AzureKeyVault is a prerequisite because a certificate is a data-plane object inside a referenced vault -- the vault must exist before the certificate can be enrolled or imported (the resource group chains transitively through the vault's own prerequisite).
- `AzureKeyVaultSecret` -- AzureKeyVault is a prerequisite because a secret is a data-plane object inside a referenced vault -- the vault must exist before the secret can be written (the resource group chains transitively through the vault's own prerequisite). Part of the Key Vault family (2005, 2025-2026) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureWebApplicationFirewallPolicy` -- AzureResourceGroup is a prerequisite because a WAF policy is created inside a referenced resource group; the Application Gateways that attach the policy reference it, never the reverse.
- `AzureApplicationSecurityGroup` -- AzureResourceGroup is a prerequisite because an application security group is created inside a referenced resource group; network interfaces, scale-set IP configurations, and NSG security rules reference the group, never the reverse.
- `AzureDiskEncryptionSet` -- AzureKeyVaultKey is a prerequisite because a disk encryption set wraps customer data with a referenced key -- the key (and its vault, which chains transitively) must exist before the set can resolve the key URL at create time.
- `AzurePostgresqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureRedisCache` -- AzureResourceGroup is a prerequisite because the cache is created inside a referenced resource group (VNet injection additionally references a dedicated subnet, but only the Premium tier supports it, so it is not a registry prerequisite).
- `AzureCosmosdbAccount` -- AzureResourceGroup is a prerequisite because the account is created inside a referenced resource group.
- `AzureMssqlServer` -- AzureResourceGroup is a prerequisite because the logical server is created inside a referenced resource group.
- `AzureMysqlFlexibleServer` -- AzureResourceGroup is a prerequisite because a server is created inside a referenced resource group (VNet injection additionally references a delegated subnet and a private DNS zone, but neither is universally required, so they are not registry prerequisites).
- `AzureMssqlDatabase` -- The parent logical server is referenced via server_id, not auto-deployed: E2E scenarios declare their own server fixture (minimal-server.yaml or the pool-attach chain through AzureMssqlElasticPool) so sequential subtests never destroy and recreate the same globally unique server_name.
- `AzureMssqlElasticPool` -- AzureMssqlServer is a prerequisite because every elastic pool lives on a referenced logical server (the server's resource group is transitive).
- `AzureRedisLinkedServer` -- The target and linked caches are referenced via ARM ids, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicy` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureRedisCacheAccessPolicyAssignment` -- The parent cache is referenced via redis_cache_id, not auto-deployed: caches are the slowest-provisioning resources in the Azure catalog and their names are globally unique, so E2E scenarios declare their own cache fixtures instead of a registry prerequisite recreating a shared one per run.
- `AzureContainerAppEnvironment` -- AzureResourceGroup is a prerequisite because the environment is created inside a referenced resource group that must already exist.
- `AzureContainerApp` -- AzureContainerAppEnvironment is a prerequisite because every app runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureServicePlan` -- AzureResourceGroup is a prerequisite because the plan is created inside a referenced resource group that must already exist.
- `AzureFunctionApp` -- AzureServicePlan is a prerequisite because a function app runs on a referenced plan (the resource group arrives transitively through the plan). The required storage account is deliberately NOT a registry prerequisite: storage-account names are globally unique, so scenarios bring their own scenario-local account fixtures.
- `AzureLinuxWebApp` -- AzureServicePlan is a prerequisite because a web app runs on a referenced plan (the resource group arrives transitively through the plan).
- `AzureContainerAppJob` -- AzureContainerAppEnvironment is a prerequisite because a job runs inside a referenced environment (the resource group arrives transitively through it).
- `AzureContainerAppEnvironmentStorage` -- AzureContainerAppEnvironment is a prerequisite because the storage registration lives on a referenced environment. The Azure Files share and storage account are deliberately NOT registry prerequisites: storage-account names are globally unique, so scenarios bring their own scenario-local account + share fixtures.
- `AzureContainerAppEnvironmentDaprComponent` -- AzureContainerAppEnvironment is a prerequisite because the Dapr component is registered on a referenced environment.
- `AzureContainerAppEnvironmentCertificate` -- AzureContainerAppEnvironment is a prerequisite because the certificate is stored on a referenced environment.
- `AzureContainerAppEnvironmentManagedCertificate` -- AzureContainerAppEnvironment is a prerequisite because the managed certificate is provisioned on a referenced environment.
- `AzureLogAnalyticsWorkspace` -- AzureResourceGroup is a prerequisite because the workspace is created inside a referenced resource group that must already exist.
- `AzureApplicationInsights` -- AzureLogAnalyticsWorkspace is a prerequisite because workspace-based Application Insights stores its telemetry in a referenced workspace (the resource group chains transitively through the workspace).
- `AzureMonitorDiagnosticSetting` -- AzureLogAnalyticsWorkspace is a prerequisite because the setting's scenarios route a fixture workspace's telemetry (the workspace doubles as target and destination); the target itself is polymorphic.
- `AzureMonitorActionGroup` -- AzureResourceGroup is a prerequisite because the action group is created inside a referenced resource group that must already exist.
- `AzureMonitorMetricAlert` -- AzureMonitorActionGroup is a prerequisite because a metric alert's actions fire into a referenced action group (the resource group chains transitively); alert scopes are polymorphic.
- `AzureMonitorScheduledQueryAlert` -- AzureLogAnalyticsWorkspace is a prerequisite because the rule queries a referenced workspace scope; AzureMonitorActionGroup because its action fires into a referenced action group.
- `AzureMonitorActivityLogAlert` -- AzureMonitorActionGroup is a prerequisite because an activity log alert's actions fire into a referenced action group (the resource group chains transitively). The alert itself is subscription-global and its scopes are polymorphic.
- `AzureApplicationInsightsStandardWebTest` -- AzureApplicationInsights is a prerequisite because a standard web test binds to a referenced Application Insights component (the resource group chains transitively through the component).
- `AzureUserAssignedIdentity` -- AzureResourceGroup is a prerequisite because the identity is created inside a referenced resource group that must already exist.
- `AzureRoleAssignment` -- AzureResourceGroup and AzureUserAssignedIdentity are prerequisites because an assignment grants a role at a referenced scope (most commonly a resource group) to a referenced principal (most commonly a managed identity) -- both must exist before the grant can be written.
- `AzureRoleDefinition` -- AzureResourceGroup is a prerequisite because a custom role definition is created at a referenced scope, most commonly a resource group in composed environments -- the scope must exist before the definition can be written.
- `AzureFederatedIdentityCredential` -- AzureUserAssignedIdentity is the prerequisite because a federated identity credential is a child resource of a referenced managed identity -- the identity must exist before the credential can be written on it. (The resource group arrives transitively through the identity's own prerequisite declaration.)
- `AzureServiceBusNamespace` -- AzureResourceGroup is a prerequisite because a Service Bus namespace is created inside a referenced resource group in composed environments. The namespace is the container every Service Bus messaging entity (queue, topic, subscription, authorization rule, geo-DR pairing) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureEventHubNamespace` -- AzureResourceGroup is a prerequisite because an Event Hub namespace is created inside a referenced resource group in composed environments. The namespace is the container every Event Hubs entity (event hub, consumer group, authorization rule, schema group, geo-DR pairing, customer-managed key) nests under. The child kinds deliberately declare NO registry prerequisite: namespace names are globally unique with a post-delete name hold, so E2E composes them with scenario-local namespace fixtures instead of a shared recreate-per-scenario prerequisite.
- `AzureServiceBusQueue`
- `AzureServiceBusTopic`
- `AzureServiceBusSubscription`
- `AzureServiceBusAuthorizationRule`
- `AzureServiceBusDisasterRecoveryConfig`
- `AzureEventHub`
- `AzureEventHubConsumerGroup`
- `AzureEventHubAuthorizationRule`
- `AzureFrontDoorProfile` -- AzureResourceGroup is a prerequisite because a Front Door profile is created inside a referenced resource group in composed environments. The profile is the container every Front Door delivery resource (endpoint, origin group, origin, route) nests under.
- `AzureFrontDoorEndpoint` -- AzureFrontDoorProfile is a prerequisite because an endpoint is an ARM child of a referenced profile -- the profile must exist before the endpoint can be written. (The resource group arrives transitively through the profile's own prerequisite declaration.)
- `AzureFrontDoorOriginGroup` -- AzureFrontDoorProfile is a prerequisite because an origin group is an ARM child of a referenced profile.
- `AzureFrontDoorOrigin` -- AzureFrontDoorOriginGroup is a prerequisite because an origin is an ARM child of a referenced origin group (the profile and resource group chain transitively).
- `AzureFrontDoorRoute` -- A route attaches to an endpoint (its ARM parent) and forwards to an origin group whose origins must exist before ARM accepts the route -- so both the endpoint and the origin chain are genuine deploy-order prerequisites.
- `AzureFrontDoorRuleSet` -- AzureFrontDoorProfile is a prerequisite because a rule set is an ARM child of a referenced profile. The rules live inside the set (they form one ordered policy document); routes attach the set by ARM ID.
- `AzureFrontDoorCustomDomain` -- AzureFrontDoorProfile is a prerequisite because a custom domain is an ARM child of a referenced profile. The DNS zone and (for bring-your-own certificates) the Front Door secret are optional references, not deploy-order prerequisites.
- `AzureFrontDoorSecret` -- AzureFrontDoorSecret is a prerequisite-light kind: only the profile (its ARM parent) must exist. The Key Vault certificate it wraps is a reference resolved before the module runs; its vault chain is exercised through scenario-local fixtures in E2E.
- `AzureFrontDoorFirewallPolicy` -- AzureResourceGroup is a prerequisite because the Front Door WAF policy is created inside a referenced resource group -- it is a GLOBAL resource, not a profile child (a different ARM type than the regional Application Gateway WAF policy). Security policies attach it to profiles; the policy itself depends on nothing else.
- `AzureFrontDoorSecurityPolicy` -- A security policy is an ARM child of a profile that associates a referenced WAF policy with referenced domains -- so the endpoint (the default-domain association target; the profile arrives transitively through it) and the WAF policy are genuine deploy-order prerequisites.
- `AzureStorageContainer` -- None of the storage data-service kinds declares a registry prerequisite on AzureStorageAccount: account names are GLOBALLY unique and Azure holds a just-deleted name, so a recreate-per-scenario fixture would hang -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the storage_account_id reference itself.
- `AzureStorageShare`
- `AzureStorageQueue`
- `AzureStorageTable`
- `AzureStorageEncryptionScope`
- `AzureStorageDataLakeGen2Filesystem`
- `AzureStorageLocalUser`
- `AzureStorageObjectReplication`
- `AzureCosmosdbSqlDatabase` -- None of the Cosmos DB data-service kinds declares a registry prerequisite on AzureCosmosdbAccount: account names are GLOBALLY unique DNS labels, so a recreate-per-scenario fixture would risk name-reuse hangs -- their E2E scenarios declare scenario-local account fixtures instead. Deploy ordering in composed environments still flows from the cosmosdb_account_id / parent-database references themselves.
- `AzureCosmosdbSqlContainer`
- `AzureCosmosdbMongoDatabase`
- `AzureCosmosdbMongoCollection`
- `AzureCosmosdbSqlRoleDefinition`
- `AzureCosmosdbSqlRoleAssignment`
- `AzureManagedRedis` -- AzureResourceGroup is the cluster's only registry prerequisite: the cluster is created inside a referenced resource group. The geo-replication and access-policy-assignment children declare NO prerequisite on AzureManagedRedis: clusters are expensive, slow-provisioning parents, so their E2E scenarios declare scenario-local cluster fixtures instead of recreating a shared one per scenario. Deploy ordering in composed environments still flows from the managed_redis_id references themselves.
- `AzureManagedRedisGeoReplication`
- `AzureManagedRedisAccessPolicyAssignment`
- `AzureEventHubDisasterRecoveryConfig`
- `AzureEventHubSchemaGroup`
- `AzureEventHubCluster` -- AzureResourceGroup is a prerequisite because a dedicated Event Hubs cluster is created inside a referenced resource group in composed environments. Note: clusters cannot be deleted for 4 hours after creation (Azure's moratorium), so E2E treats this kind as offline-gated.
- `AzureEventHubNamespaceCustomerManagedKey`
- `AzureMssqlFailoverGroup` -- AzureMssqlServer is a prerequisite because a failover group is created on a referenced primary logical server and points at a partner server; the primary (and its resource group, which chains transitively) must exist before the group can be written.
- `AzureContainerAppCustomDomain` -- AzureContainerApp is a prerequisite because the domain binding lives in a referenced app's ingress configuration (the environment and resource group chain transitively through the app).
- `AzureFirewallPolicy`
- `AzureFirewallPolicyRuleCollectionGroup` -- AzureFirewallPolicy is a prerequisite because a rule collection group is a child document of a referenced policy (the resource group chains transitively through the policy).
- `AzureFirewall` -- AzureSubnet is a prerequisite because a VNet-deployed firewall's data path lives in a dedicated subnet that must be named exactly "AzureFirewallSubnet" (the virtual network and resource group chain transitively through the subnet). The E2E install profile publishes a fixture subnet with that exact name and a /26 prefix.
- `AzureIpGroup`
- `AzureVirtualNetworkGateway` -- AzureSubnet is a prerequisite because every virtual network gateway lives in a dedicated subnet named exactly "GatewaySubnet" (the virtual network and resource group chain transitively through the subnet); the subnet install profile publishes a fixture instance with that exact ARM name. AzurePublicIp is a prerequisite because a VPN-type gateway (the default shape) requires a public IP per ip configuration; the address install profile publishes a dedicated zone-redundant instance (a gateway binds its address exclusively, and the AZ gateway SKUs require zones on it).
- `AzureVirtualNetworkGatewayConnection` -- Both gateways are prerequisites: a connection joins a virtual network gateway to a far side, and the site-to-site far side is a local network gateway (the GatewaySubnet, VNet, and resource group chain transitively through the virtual network gateway).
- `AzureLocalNetworkGateway`
- `AzurePrivateLinkService` -- AzureSubnet is the sole prerequisite: every NAT ip configuration draws its address from a subnet with private-link-service network policies disabled (the subnet install profile publishes a fixture instance with that flag off). The Standard load balancer whose frontend the service typically fronts is NOT a registry prerequisite -- the spec's destination is an exactly-one-of (load balancer frontend OR fixed destination IP), so scenarios that use the load-balancer shape declare it via the planton.dev/e2e-prerequisites annotation instead.
- `AzureExpressRouteCircuit`
- `AzureExpressRouteCircuitPeering` -- The circuit is the prerequisite: a peering is an ARM child of the circuit, addressed by the circuit's name (the resource group chains transitively through the circuit).
- `AzureExpressRouteGateway` -- The hub is the prerequisite: ARM requires an ExpressRoute Gateway to be deployed INTO a Virtual WAN hub (the WAN and resource group chain transitively through the hub).
- `AzureExpressRoutePort` -- ExpressRoute Port: your own physical port pair on a Microsoft edge router (ExpressRoute Direct), from whose bandwidth circuits are carved. Self-contained -- only the resource group is required.
- `AzureVirtualWan` -- Virtual WAN: the umbrella of Azure's managed hub-and-spoke networking, under which virtual hubs and their gateways are created. Self-contained -- only the resource group is required.
- `AzureVirtualHub` -- The WAN is the prerequisite: this kind models the Virtual WAN hub (virtual_wan_id is required; standalone hubs are the legacy Route Server construction, which has its own ARM surface). The resource group chains transitively through the WAN.
- `AzureVirtualHubConnection` -- Both sides of the attachment are prerequisites: the hub being joined and the spoke virtual network being attached.
- `AzureVpnGateway` -- The hub is the prerequisite: ARM deploys a Virtual WAN VPN gateway INTO a virtual hub (virtual_hub_id is required and immutable; the WAN and resource group chain transitively through the hub). ARM allows one VPN gateway per hub.
- `AzureVpnGatewayConnection` -- Both ends of the tunnel are prerequisites: a connection is an ARM child of the VPN gateway and pins each of its links to a specific link of the remote VPN site (the hub, WAN, and resource group chain transitively through the gateway).
- `AzureVpnSite` -- The WAN is the prerequisite: a VPN site is the Virtual WAN world's address-book entry for one branch location (virtual_wan_id is required; the classic-world sibling without a WAN is AzureLocalNetworkGateway). The resource group chains transitively through the WAN.
- `AzurePointToSiteVpnGateway` -- The hub and the server configuration are both prerequisites: a point-to-site VPN gateway deploys INTO a virtual hub (one P2S gateway per hub, a slot separate from the hub's site-to-site VPN gateway) and is born pointing at the VPN server configuration that defines how its users authenticate -- both ARM-required and fixed at creation. The WAN and resource group chain transitively through the hub.
- `AzureVpnServerConfiguration` -- Self-contained -- only the resource group is required: a VPN server configuration is the reusable "who may connect and how" authentication policy (Entra ID / certificate / RADIUS) that point-to-site VPN gateways attach to; it references no other Azure resource.
- `AzureCognitiveAccount` -- Self-contained -- only the resource group is required: an Azure AI services account (Azure OpenAI, the multi-service AIServices account, the single-service accounts) needs no other Azure resource; subnets (network rules), Key Vault keys (CMK), storage accounts and user-assigned identities are optional references.
- `AzureCognitiveDeployment` -- An ARM child of its account: a model deployment (which model runs, at which throughput class) exists only on an Azure AI services account of kind "OpenAI" or "AIServices".
- `AzureCognitiveAccountProject` -- An ARM child of its account: an AI Foundry project exists only on an "AIServices"-kind account with project management enabled.
- `AzureMachineLearningWorkspace` -- The workspace REQUIRES all three companion services at creation (default storage, secrets vault, telemetry) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureMachineLearningDatastore` -- An ARM child of its workspace. The storage target (container, filesystem or share) is scenario-declared via the e2e-prerequisites annotation -- only the blob scenario needs a container, so it is not a kind-wide prerequisite.
- `AzureMachineLearningComputeCluster` -- An ARM child of its workspace (.../computes/{name}) -- the auto-scaling pool of VMs training jobs run on.
- `AzureMachineLearningComputeInstance` -- An ARM child of its workspace (.../computes/{name}) -- a single always-on VM serving as one data scientist's cloud workstation.
- `AzureAiFoundry` -- The hub REQUIRES both companion services at creation (secrets vault, default storage) -- genuine deploy-order prerequisites, each with its own fixture profile.
- `AzureAiFoundryProject` -- Deploys into its hub's resource group (the provider derives the group from the hub reference -- the project spec carries none).
- `AzureSearchService`
- `AzureMachineLearningOnlineEndpoint` -- An ARM child of its workspace (.../onlineEndpoints/{name}) -- the stable scoring address applications call. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningOnlineDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the running copy of a model the endpoint's traffic map routes to.
- `AzureMachineLearningBatchEndpoint` -- An ARM child of its workspace (.../batchEndpoints/{name}) -- the stable address batch scoring jobs are submitted to. azurerm carries no ML endpoint resources; the modules write the raw ARM shape at a pinned api-version (azapi / azure-native).
- `AzureMachineLearningBatchDeployment` -- An ARM child of its endpoint (.../deployments/{name}) -- the job recipe (model, compute, batching behavior) the endpoint's default-deployment pointer routes submissions to.
- `AzureRecoveryServicesVault` -- The Recovery Services vault (Microsoft.RecoveryServices/vaults) -- the safe that classic Azure Backup data and Site Recovery configuration live in. Backup policies and protected items are ARM children of a vault.
- `AzureBackupPolicyVm` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern IaaS VM backups.
- `AzureBackupProtectedVm` -- An ARM child of its vault (.../protectedItems/...) -- the binding that puts one virtual machine under a backup policy's protection.
- `AzureBackupPolicyFileShare` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules that govern Azure Files share backups (snapshot or vaulted).
- `AzureBackupProtectedFileShare` -- An ARM child of its vault (.../protectedItems/AzureFileShare;...) -- the binding that puts one Azure Files share under a backup policy's protection. The share's storage account must already be registered with the vault (AzureBackupContainerStorageAccount).
- `AzureDataProtectionBackupVault` -- The Data Protection backup vault (Microsoft.DataProtection/ backupVaults) -- the safe that MODERN Azure Backup data lives in (managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, Data Lake storage). Backup policies and backup instances are ARM children of a vault.
- `AzureDataProtectionBackupPolicy` -- An ARM child of its vault (.../backupPolicies/{name}) -- the schedule and retention rules for ONE Data Protection datasource type (blob storage, disk, Kubernetes cluster, MySQL/PostgreSQL flexible server, or Data Lake storage), modeled as one kind with variant blocks.
- `AzureDataProtectionBackupInstance` -- An ARM child of its vault (.../backupInstances/{name}) -- the binding that puts ONE datasource (a managed disk, a storage account's blob services, an AKS cluster, a MySQL/PostgreSQL flexible server, or a Data Lake storage account) under a Data Protection backup policy, modeled as one kind with variant blocks. The vault's managed identity must hold the datasource roles Azure Backup requires BEFORE the instance is created.
- `AzureBastionHost` -- AzureSubnet and AzurePublicIp are prerequisites because a dedicated-infrastructure Bastion host (Basic/Standard/Premium -- the default shapes) deploys into a subnet named exactly "AzureBastionSubnet" and binds a Standard static public IP EXCLUSIVELY (the virtual network and resource group chain transitively through the subnet). The Developer SKU instead attaches to a virtual network directly and uses neither.
- `AzureNetworkWatcherFlowLog` -- AzureVirtualNetwork and AzureStorageAccount are prerequisites because a flow log records a network-scoped target (a virtual network in the common case; subnets and network interfaces chain through the network) into a referenced storage account. The regional Network Watcher parent is NOT a prerequisite: Azure auto-creates it ("NetworkWatcher_{region}" in "NetworkWatcherRG") the moment the region hosts a virtual network, and the flow log references it by name. Traffic Analytics' Log Analytics workspace is an optional arm, declared by scenarios that use it.
- `AzurePrivateDnsResolver` -- AzureVirtualNetwork and AzureSubnet are prerequisites because a DNS Private Resolver anchors to a referenced virtual network (at most ONE resolver per network -- Azure enforces it) and each of its inbound/outbound endpoints occupies its own dedicated subnet delegated to "Microsoft.Network/dnsResolvers" (the resource group chains transitively through the network and subnets).
- `AzurePrivateDnsResolverForwardingRuleset` -- AzurePrivateDnsResolver is a prerequisite because a DNS forwarding ruleset steers a resolver's OUTBOUND endpoints -- it binds their ARM ids (at most 2, same resolver) at creation. (The resource group and network chain transitively through the resolver's own prerequisite declarations.)
- `AzureBackupContainerStorageAccount` -- Registers a storage account with a Recovery Services vault as a backup container (.../protectionContainers/StorageContainer;...) -- one registration per storage-account-and-vault pair, required BEFORE any of the account's file shares can be protected. Part of the backup family (2175-2179) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzureDataProtectionResourceGuard` -- The Data Protection Resource Guard (Microsoft.DataProtection/ resourceGuards) -- the approval gate behind Multi-User Authorization: privileged vault operations (disabling soft delete, reducing retention) require an approval through a guard, which typically lives in a DIFFERENT administrator's scope. Vaults reference a guard by its ARM ID. Part of the backup family (2175-2182) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `AzurePrivateDnsResolverVirtualNetworkLink` -- The attachment that makes a DNS forwarding ruleset take effect in one virtual network ({ruleset_id}/virtualNetworkLinks/{name}) -- one link per ruleset-network pair, up to 500 per ruleset, spokes joining and leaving independently (which is why the link is a standalone kind, exactly like AzurePrivateDnsZoneVirtualNetworkLink). Part of the DNS Private Resolver family (2186-2187) despite the out-of-run number -- enum numbers are pinned by the registry snapshot; never renumber.
- `GcpArtifactRegistryRepo` -- 3000–3999: GCP resources
- `GcpTargetHttpsProxy` -- The URL map is the parent a proxy cannot exist without; the classic compute certificate kinds and the SSL policy are the fixture parents the committed scenarios attach. The Certificate Manager certificate list (certificate_manager_certificates, honored only by the cross-region internal ALB) is optional composition -- a scenario that arms it declares GcpCertManagerCert via the e2e-prerequisites annotation, never a registry edge that would tax every proxy and forwarding-rule chain.
- `GcpCloudFunction`
- `GcpCloudRun`
- `GcpCloudSql`
- `GcpDnsZone`
- `GcpGcsBucket`
- `GcpGkeCluster`
- `GcpIamCustomRole`
- `GcpProject`
- `GcpVpcNetwork`
- `GcpSubnetwork`
- `GcpRouterNat`
- `GcpGkeNodePool`
- `GcpServiceAccount`
- `GcpGkeWorkloadIdentityBinding`
- `GcpCertManagerCert`
- `GcpComputeInstance`
- `GcpDnsRecord`
- `GcpProjectIamMember`
- `GcpFirewallRule`
- `GcpGlobalAddress`
- `GcpCloudArmorPolicy`
- `GcpHealthCheck`
- `GcpBackendBucket`
- `GcpBackendService`
- `GcpRegionNetworkEndpointGroup`
- `GcpUrlMap`
- `GcpManagedSslCertificate`
- `GcpTargetHttpProxy`
- `GcpAlloydbCluster`
- `GcpRedisInstance`
- `GcpFirestoreDatabase`
- `GcpSpannerInstance`
- `GcpSpannerDatabase`
- `GcpBigtableInstance`
- `GcpMemorystoreInstance`
- `GcpCloudSqlDatabase`
- `GcpCloudSqlUser`
- `GcpAlloydbInstance`
- `GcpAlloydbUser`
- `GcpSpannerBackupSchedule`
- `GcpBigtableTable`
- `GcpFirestoreBackupSchedule`
- `GcpFirestoreIndex`
- `GcpBigQueryDataset`
- `GcpDataprocCluster`
- `GcpDataprocAutoscalingPolicy`
- `GcpBigQueryTable`
- `GcpPubSubTopic`
- `GcpPubSubSubscription`
- `GcpCloudTasksQueue`
- `GcpCloudSchedulerJob`
- `GcpPubSubSchema`
- `GcpVertexAiNotebook`
- `GcpVertexAiEndpoint`
- `GcpVertexAiIndex`
- `GcpVertexAiIndexEndpoint` -- Vector Search IndexEndpoint — distinct from the online-prediction GcpVertexAiEndpoint (671); different GCP resources, different kinds.
- `GcpVertexAiDeployedIndex`
- `GcpCloudComposerEnvironment`
- `GcpCloudComposerUserWorkloadsSecret`
- `GcpCloudComposerUserWorkloadsConfigMap`
- `GcpKmsKeyRing`
- `GcpKmsKey`
- `GcpKmsKeyIamMember`
- `GcpFilestoreInstance`
- `GcpWorkloadIdentityPool` -- 3101–3109: IAM/identity family (overflow block; the 3000–3022 foundation/security sub-band is fully allocated)
- `GcpWorkloadIdentityPoolProvider`
- `GcpServiceAccountIamMember`
- `GcpGlobalForwardingRule` -- 3110–3119: networking/load-balancer family (overflow block; the 3023–3029 LB sub-band is fully allocated)
- `GcpSslPolicy`
- `GcpSslCertificate`
- `GcpServiceNetworkingConnection`
- `GcpAddress`
- `GcpServiceConnectionPolicy`
- `GcpCertManagerDnsAuthorization`
- `GcpCertificateMap` -- GcpCertManagerCert is a prerequisite because a map entry binds hostnames to EXISTING certificates — the canonical map references a certificate fixture's resource name.
- `GcpCloudRunJob` -- 3120–3129: GCP serverless overflow
- `GcpServerlessVpcConnector`
- `GcpComputeDisk` -- 3130–3139: GCP compute overflow (the 3000–3022 foundation sub-band that holds GcpComputeInstance is fully allocated)
- `GcpComputeMig` -- GcpVpcNetwork is a prerequisite because the canonical group runs its fleet on a dedicated custom-mode VPC — a managed instance group's template must attach every VM to a network, and the default VPC is never assumed.
- `GcpMonitoringNotificationChannel` -- 3140–3149: GCP observability & log routing
- `GcpMonitoringAlertPolicy` -- GcpMonitoringNotificationChannel is a prerequisite because the policy's canonical shape references a channel to notify — a policy without a delivery endpoint measures but never pages.
- `GcpMonitoringUptimeCheck`
- `GcpLoggingSink` -- GcpGcsBucket is a prerequisite because the canonical sink exports to a Cloud Storage bucket — the cheapest destination that proves the whole writer-identity grant flow.
- `GcpMonitoringDashboard`
- `GcpMonitoringSlo`
- `GcpLogBucket`
- `GcpLogMetric`
- `GcpSecretManagerSecret` -- 3150–3159: GCP security & identity GcpServiceAccount is a prerequisite because the canonical secret grants secretAccessor to a workload service account — the access story the kind exists to model.
- `GcpIdentityPlatformConfig`
- `GcpIdentityPlatformTenant` -- GcpIdentityPlatformConfig is a prerequisite because tenants exist only in projects whose Identity Platform config enables multi_tenant.allow_tenants — a tenant without the initialized, tenant-enabled project config cannot be created at all.
- `GcpIamOauthClient`
- `GcpIamDenyPolicy`
- `GcpCloudRunDomainMapping` -- 3160–3169: GCP serverless edge GcpCloudRun is a prerequisite because a domain mapping exists only to point a verified domain at a running Cloud Run service — the route it maps must already exist for the mapping to be created at all.
- `GcpWorkflow`
- `GcpEventarcTrigger` -- GcpCloudRun is a prerequisite because the canonical trigger routes a Pub/Sub messagePublished event to a Cloud Run service — the destination story the kind exists to model.
- `GcpEventarcMessageBus`
- `KubernetesNamespace` -- 4000–4999: Kubernetes resources, organized in family sub-bands (4030–4069 also hosts CNI/autoscaling/DR addons; 4130–4149 hosts analytics & ML; 4190–4199 reserved for growth) 4000–4029: Kubernetes building blocks (core API primitives)
- `KubernetesDeployment`
- `KubernetesStatefulSet`
- `KubernetesDaemonSet`
- `KubernetesJob`
- `KubernetesCronJob`
- `KubernetesService`
- `KubernetesSecret`
- `KubernetesManifest`
- `KubernetesHelmRelease`
- `KubernetesConfigMap`
- `KubernetesServiceAccount`
- `KubernetesRbac` -- Bundles the RBAC grant grain (Role/ClusterRole + its binding) into one component: "grant these permissions to these subjects in this scope".
- `KubernetesIngress`
- `KubernetesNetworkPolicy`
- `KubernetesPersistentVolumeClaim`
- `KubernetesStorageClass`
- `KubernetesResourceQuota` -- Manages the namespace-governance pair: the ResourceQuota plus an optional companion LimitRange (per-object defaults/bounds) — two API objects, one governance story.
- `KubernetesPriorityClass`
- `KubernetesPodDisruptionBudget`
- `KubernetesHorizontalPodAutoscaler`
- `KubernetesCertManager` -- 4030–4069: Kubernetes foundation addons (certs, DNS, secrets, ingress, Gateway API, mesh, CNI/autoscaling/DR)
- `KubernetesClusterIssuer` -- KubernetesCertManager is a prerequisite for the three cert-manager CR kinds below: ClusterIssuer/Issuer/Certificate are cert-manager custom resources — without the controller and its CRDs they cannot be applied.
- `KubernetesIssuer`
- `KubernetesCertificate`
- `KubernetesExternalDns`
- `KubernetesExternalSecretsOperator`
- `KubernetesClusterSecretStore` -- KubernetesExternalSecretsOperator is a prerequisite for the three external-secrets CR kinds below: ClusterSecretStore/SecretStore/ ExternalSecret are external-secrets custom resources — without the operator and its CRDs they cannot be applied.
- `KubernetesSecretStore`
- `KubernetesExternalSecret`
- `KubernetesIngressNginx`
- `KubernetesGatewayApiCrds`
- `KubernetesGatewayClass`
- `KubernetesGateway`
- `KubernetesListenerSet`
- `KubernetesHttpRoute`
- `KubernetesGrpcRoute`
- `KubernetesTcpRoute`
- `KubernetesUdpRoute`
- `KubernetesTlsRoute`
- `KubernetesReferenceGrant`
- `KubernetesBackendTlsPolicy`
- `KubernetesIstioBaseCrds`
- `KubernetesIstio`
- `KubernetesDestinationRule` -- Istio API components (mesh traffic policy, security, telemetry). The seven typed resources below (4053–4059) require the Istio CRDs on the cluster, provided by the lightweight CRDs-only KubernetesIstioBaseCrds (851) — NOT the full mesh KubernetesIstio (852).
- `KubernetesServiceEntry`
- `KubernetesPeerAuthentication`
- `KubernetesRequestAuthentication`
- `KubernetesAuthorizationPolicy`
- `KubernetesTelemetry`
- `KubernetesEnvoyFilter`
- `KubernetesMetricsServer`
- `KubernetesCilium`
- `KubernetesKeda`
- `KubernetesKarpenter`
- `KubernetesKarpenterNodePool`
- `KubernetesKarpenterEc2NodeClass`
- `KubernetesClusterAutoscaler`
- `KubernetesVelero`
- `KubernetesKubePrometheusStack` -- 4070–4089: Kubernetes observability
- `KubernetesGrafana`
- `KubernetesSignoz` -- KubernetesClickHouse is a prerequisite because SigNoz stores every trace, metric and log in ClickHouse and deploys none of its own — the telemetry store is composed, never bundled.
- `KubernetesLoki`
- `KubernetesTempo`
- `KubernetesOtelOperator` -- The operator's admission webhooks (failurePolicy Fail) are served with a cert-manager Certificate in the default posture — cert-manager must be running before the operator installs.
- `KubernetesOtelCollector`
- `KubernetesKyverno` -- 4080–4099: Kubernetes security, policy, and identity
- `KubernetesGatekeeper`
- `KubernetesKeycloak` -- Keycloak declarations compose the official Keycloak Operator (which reconciles the Keycloak CR this kind renders) and, on the recommended postgres vendor, a KubernetesPostgres database — both must resolve before the CR can converge.
- `KubernetesOpenBao`
- `KubernetesOpenFga` -- OpenFGA requires a datastore; the recommended arm composes a KubernetesPostgres database (the sandbox memory arm needs nothing, but the registry declares the shape real deployments require).
- `KubernetesKeycloakOperator`
- `KubernetesCloudNativePgOperator` -- 4100–4129: Kubernetes data platforms
- `KubernetesPostgres`
- `KubernetesValkey`
- `KubernetesPerconaMysqlOperator`
- `KubernetesMysql`
- `KubernetesPerconaMongoOperator`
- `KubernetesMongodb`
- `KubernetesStrimziKafkaOperator`
- `KubernetesKafka` -- container_kind: a Strimzi Kafka cluster is a place in the provider's own model — KafkaTopic and KafkaUser declarations BELONG to one cluster (the strimzi.io/cluster label) and are drawn inside its box. Clients that merely talk to the cluster (Connect, MirrorMaker2, UI, Karapace) carry containment_exempt on their bootstrap/trust references.
- `KubernetesKafkaTopic`
- `KubernetesKafkaUser`
- `KubernetesKafkaConnect` -- container_kind: a Connect cluster hosts the connectors deployed INTO it (KafkaConnector's strimzi.io/cluster label names its Connect cluster) — the same room shape as KubernetesKafka above.
- `KubernetesKafkaConnector`
- `KubernetesKafkaMirrorMaker2`
- `KubernetesKarapace`
- `KubernetesKafkaUi`
- `KubernetesOpenSearchOperator`
- `KubernetesOpenSearch`
- `KubernetesAltinityOperator`
- `KubernetesClickHouse`
- `KubernetesSolrOperator`
- `KubernetesSolr`
- `KubernetesNeo4j`
- `KubernetesSeaweedFs`
- `KubernetesQdrant`
- `KubernetesRabbitMqOperator` -- The RabbitMQ Cluster Operator's release manifest ships admission webhooks whose serving certificate is a cert-manager Certificate — cert-manager must be running before the operator installs.
- `KubernetesRabbitMq`
- `KubernetesAirflow` -- 4130–4149: Kubernetes analytics and ML KubernetesPostgres is a prerequisite because Airflow's metadata database composes a KubernetesPostgres by default (the spec's FK defaults resolve onto its outputs) and the migration Job needs the database reachable before the server components start.
- `KubernetesSparkOperator`
- `KubernetesKubeRayOperator`
- `KubernetesRayCluster` -- KubernetesKubeRayOperator is a prerequisite because this kind declares the RayCluster custom resource that only the operator's CRDs admit and only the operator reconciles into head and worker pods.
- `KubernetesFlinkOperator` -- KubernetesCertManager is a prerequisite because the Flink operator's chart, with its default-on admission webhook, renders cert-manager Issuer/Certificate resources and trusts the API server through cert-manager's CA injection — there is no self-signed fallback at the pinned chart, and the webhooks are fail-closed.
- `KubernetesFlinkDeployment` -- KubernetesFlinkOperator is a prerequisite because this kind declares the FlinkDeployment custom resource that only the operator's CRDs admit and only the operator reconciles into a running Flink cluster.
- `KubernetesJupyterHub` -- KubernetesPostgres is a prerequisite because JupyterHub's hub database composes a KubernetesPostgres in its external-database arm (the spec's FK defaults resolve onto its outputs) and the hub pod mounts that database's credential Secret before it can start.
- `KubernetesMlflow` -- KubernetesPostgres is a prerequisite because MLflow's backend store composes a KubernetesPostgres in its production arm (FK defaults onto its outputs; the module composes the connection URI from its credential Secret), and KubernetesSeaweedFs because the artifact store's S3-compatible arm FK-defaults onto the SeaweedFS endpoint and credential Secret.
- `KubernetesTrino` -- KubernetesPostgres is a prerequisite because Trino's postgres catalogs compose a KubernetesPostgres (the catalog host and credential FK-default onto its outputs), and the pods read that database's credential Secret to resolve catalog passwords from environment.
- `KubernetesSuperset` -- KubernetesPostgres is a prerequisite because Superset's REQUIRED metadata database composes a KubernetesPostgres (FK defaults onto its outputs; the module composes the environment Secret from its credential Secret), and KubernetesValkey because the cache/broker arm FK-defaults onto a KubernetesValkey's service and password Secret.
- `KubernetesArgocd` -- 4150–4169: Kubernetes GitOps and CI/CD
- `KubernetesArgoWorkflows`
- `KubernetesTektonOperator`
- `KubernetesTekton` -- KubernetesTektonOperator is a prerequisite because this kind declares the TektonConfig custom resource that only the operator's CRDs admit and only the operator reconciles into running components.
- `KubernetesGhaRunnerScaleSetController`
- `KubernetesGhaRunnerScaleSet` -- KubernetesGhaRunnerScaleSetController is a prerequisite because this kind renders an AutoscalingRunnerSet custom resource that only the controller's CRDs admit and only the controller reconciles into listener and runner pods.
- `KubernetesHarbor`
- `KubernetesJenkins`
- `KubernetesTemporal` -- 4170–4189: Kubernetes app platforms KubernetesPostgres is a prerequisite because the recommended (and E2E-proven) database composition backs Temporal's default and visibility stores with a CloudNativePG cluster.
- `KubernetesNats`
- `KubernetesLocust`
- `DigitalOceanAppPlatformService` -- 5000–5999: DigitalOcean resources
- `DigitalOceanBucket`
- `DigitalOceanContainerRegistry`
- `DigitalOceanDatabaseCluster`
- `DigitalOceanDnsZone`
- `DigitalOceanDroplet`
- `DigitalOceanFirewall`
- `DigitalOceanFunction`
- `DigitalOceanKubernetesCluster`
- `DigitalOceanKubernetesNodePool`
- `DigitalOceanLoadBalancer`
- `DigitalOceanVolume`
- `DigitalOceanVpc`
- `DigitalOceanCertificate`
- `DigitalOceanDnsRecord`
- `CivoBucket` -- 6000–6999: Civo resources
- `CivoCertificate`
- `CivoComputeInstance`
- `CivoDatabase`
- `CivoDnsZone`
- `CivoFirewall`
- `CivoIpAddress`
- `CivoKubernetesCluster`
- `CivoKubernetesNodePool`
- `CivoVolume`
- `CivoVpc`
- `CivoDnsRecord`
- `CloudflareDnsZone` -- 7000–7999: Cloudflare resources
- `CloudflareKvNamespace`
- `CloudflareR2Bucket`
- `CloudflareWorker`
- `CloudflareLoadBalancer`
- `CloudflareD1Database`
- `CloudflareZeroTrustAccessApplication`
- `CloudflareDnsRecord`
- `CloudflareRuleset`
- `CloudflareWorkersKvPair`
- `CloudflareHyperdriveConfig`
- `CloudflareLoadBalancerPool`
- `CloudflareLoadBalancerMonitor`
- `CloudflareZeroTrustAccessPolicy`
- `CloudflareZeroTrustAccessGroup`
- `CloudflareQueue`
- `CloudflarePagesProject`
- `CloudflareZeroTrustTunnel`
- `CloudflareZeroTrustTunnelVirtualNetwork`
- `CloudflareZeroTrustTunnelRoute`
- `CloudflareList`
- `CloudflareListItem`
- `CloudflareTurnstileWidget`
- `CloudflareEmailRoutingZone`
- `CloudflareEmailRoutingRule`
- `CloudflareEmailRoutingAddress`
- `CloudflareOriginCaCertificate`
- `CloudflareCertificatePack`
- `CloudflareCustomHostname`
- `CloudflareCustomHostnameFallbackOrigin`
- `Auth0Connection` -- 8000–8999: Auth0 resources
- `Auth0Client`
- `Auth0EventStream`
- `Auth0ResourceServer`
- `Auth0Action`
- `Auth0Role`
- `OpenFgaStore` -- 9000–9999: OpenFGA resources Note: OpenFGA is Terraform-only - there is no Pulumi provider available. Pulumi modules for OpenFGA resources are pass-through placeholders.
- `OpenFgaAuthorizationModel`
- `OpenFgaRelationshipTuple`
- `OpenStackKeypair` -- 10000–10999: OpenStack resources
- `OpenStackNetwork`
- `OpenStackSubnet`
- `OpenStackRouter`
- `OpenStackRouterInterface`
- `OpenStackSecurityGroup`
- `OpenStackFloatingIp`
- `OpenStackNetworkPort`
- `OpenStackSecurityGroupRule`
- `OpenStackFloatingIpAssociate`
- `OpenStackInstance`
- `OpenStackServerGroup`
- `OpenStackVolume`
- `OpenStackVolumeAttach`
- `OpenStackProject`
- `OpenStackApplicationCredential`
- `OpenStackImage`
- `OpenStackRoleAssignment`
- `OpenStackLoadBalancer`
- `OpenStackLoadBalancerListener`
- `OpenStackLoadBalancerPool`
- `OpenStackLoadBalancerMember`
- `OpenStackLoadBalancerMonitor`
- `OpenStackDnsZone`
- `OpenStackDnsRecord`
- `ScalewayVpc`
- `ScalewayPrivateNetwork`
- `ScalewayPublicGateway`
- `ScalewayLoadBalancer`
- `ScalewayInstanceSecurityGroup`
- `ScalewayInstance`
- `ScalewayKapsuleCluster`
- `ScalewayKapsulePool`
- `ScalewayRdbInstance`
- `ScalewayRedisCluster`
- `ScalewayMongodbInstance`
- `ScalewayObjectBucket`
- `ScalewayBlockVolume`
- `ScalewayContainerRegistry`
- `ScalewayDnsZone`
- `ScalewayDnsRecord`
- `ScalewayServerlessFunction`
- `ScalewayServerlessContainer`
- `AliCloudLogProject`
- `AliCloudRamRole`
- `AliCloudRamPolicy`
- `AliCloudVpc`
- `AliCloudVswitch`
- `AliCloudSecurityGroup`
- `AliCloudEipAddress`
- `AliCloudNatGateway`
- `AliCloudApplicationLoadBalancer`
- `AliCloudNetworkLoadBalancer`
- `AliCloudVpnGateway`
- `AliCloudDnsZone`
- `AliCloudDnsRecord`
- `AliCloudPrivateDnsZone`
- `AliCloudStorageBucket`
- `AliCloudNasFileSystem`
- `AliCloudKmsKey`
- `AliCloudRdsInstance`
- `AliCloudPolardbCluster`
- `AliCloudRedisInstance`
- `AliCloudMongodbInstance`
- `AliCloudEcsInstance`
- `AliCloudContainerRegistry`
- `AliCloudKubernetesCluster`
- `AliCloudKubernetesNodePool`
- `AliCloudCdnDomain`
- `AliCloudFunction`
- `AliCloudSaeApplication`
- `AliCloudRocketmqInstance`
- `AliCloudCenInstance`
- `OciVcn`
- `OciSubnet`
- `OciSecurityGroup`
- `OciCompartment`
- `OciIdentityPolicy`
- `OciDynamicGroup`
- `OciComputeInstance`
- `OciContainerEngineCluster`
- `OciContainerEngineNodePool`
- `OciContainerInstance`
- `OciApplicationLoadBalancer`
- `OciNetworkLoadBalancer`
- `OciDynamicRoutingGateway`
- `OciPublicIp`
- `OciAutonomousDatabase`
- `OciDbSystem`
- `OciMysqlDbSystem`
- `OciPostgresqlDbSystem`
- `OciRedisCluster`
- `OciNosqlTable`
- `OciObjectStorageBucket`
- `OciFileSystem`
- `OciBlockVolume`
- `OciKmsVault`
- `OciKmsKey`
- `OciVaultSecret`
- `OciBastion`
- `OciFunctionsApplication`
- `OciApiGateway`
- `OciStreamPool`
- `OciQueue`
- `OciAlarm`
- `OciLogGroup`
- `OciDnsZone`
- `OciDnsRecord`
- `OciNetworkFirewall`
- `OciDevopsProject`
- `HetznerCloudSshKey`
- `HetznerCloudPlacementGroup`
- `HetznerCloudFirewall`
- `HetznerCloudNetwork`
- `HetznerCloudPrimaryIp`
- `HetznerCloudFloatingIp`
- `HetznerCloudServer`
- `HetznerCloudVolume`
- `HetznerCloudSnapshot`
- `HetznerCloudCertificate`
- `HetznerCloudLoadBalancer`
- `HetznerCloudDnsZone`

### spec.pod.initContainers[].env.secrets[].valueFrom.env

`string`

### spec.pod.initContainers[].env.secrets[].valueFrom.name

`string` · required

- rule: {"required":true}

### spec.pod.initContainers[].env.secrets[].valueFrom.fieldPath

`string`

### spec.pod.initContainers[].env.envFrom

`[]EnvFromSource`

Bulk import of environment variables from ConfigMaps or Secrets.

### spec.pod.initContainers[].env.envFrom[].prefix

`string`

Optional prefix prepended to each imported key name.
For example, prefix "APP_" with key "PORT" produces env var "APP_PORT".

### spec.pod.initContainers[].env.envFrom[].configMapRef

`ConfigMapRef`

Import all keys from a ConfigMap.

### spec.pod.initContainers[].env.envFrom[].configMapRef.name

`string` · required

Name of the ConfigMap.

- rule: {"required":true}

### spec.pod.initContainers[].env.envFrom[].configMapRef.optional

`bool`

If true, the ConfigMap is allowed to not exist without blocking pod startup.

### spec.pod.initContainers[].env.envFrom[].secretRef

`SecretRef`

Import all keys from a Secret.

### spec.pod.initContainers[].env.envFrom[].secretRef.name

`string` · required

Name of the Secret.

- rule: {"required":true}

### spec.pod.initContainers[].env.envFrom[].secretRef.optional

`bool`

If true, the Secret is allowed to not exist without blocking pod startup.

### spec.pod.initContainers[].resources

`ContainerResources`

CPU and memory requests and limits. Requests drive scheduling and are what the
pod is guaranteed; limits are the ceiling enforced at runtime (CPU is throttled,
memory overage is OOM-killed). Omitting limits entirely leaves the container
unbounded — acceptable for batch work on dedicated nodes, risky on shared ones.

### spec.pod.initContainers[].resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.pod.initContainers[].resources.limits.cpu

`string`

### spec.pod.initContainers[].resources.limits.memory

`string`

### spec.pod.initContainers[].resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.pod.initContainers[].resources.requests.cpu

`string`

### spec.pod.initContainers[].resources.requests.memory

`string`

### spec.pod.initContainers[].livenessProbe

`Probe`

Liveness probe: restarts the container when it fails. Detects deadlocked or
wedged processes. Keep it strictly about "is the process alive" — checking
downstream dependencies here turns a dependency blip into a restart storm.

### spec.pod.initContainers[].livenessProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.pod.initContainers[].livenessProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.pod.initContainers[].livenessProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.pod.initContainers[].livenessProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.pod.initContainers[].livenessProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.pod.initContainers[].livenessProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.pod.initContainers[].livenessProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.pod.initContainers[].livenessProbe.httpGet.portNumber

`int32`

### spec.pod.initContainers[].livenessProbe.httpGet.portName

`string`

### spec.pod.initContainers[].livenessProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.pod.initContainers[].livenessProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.pod.initContainers[].livenessProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.pod.initContainers[].livenessProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.pod.initContainers[].livenessProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.pod.initContainers[].livenessProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.pod.initContainers[].livenessProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.pod.initContainers[].livenessProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.pod.initContainers[].livenessProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.pod.initContainers[].livenessProbe.tcpSocket.portNumber

`int32`

### spec.pod.initContainers[].livenessProbe.tcpSocket.portName

`string`

### spec.pod.initContainers[].livenessProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.pod.initContainers[].livenessProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.pod.initContainers[].livenessProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.pod.initContainers[].readinessProbe

`Probe`

Readiness probe: removes the pod from Service endpoints while it fails. This is
the probe that makes rolling updates zero-downtime — traffic only reaches pods
that report ready.

### spec.pod.initContainers[].readinessProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.pod.initContainers[].readinessProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.pod.initContainers[].readinessProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.pod.initContainers[].readinessProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.pod.initContainers[].readinessProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.pod.initContainers[].readinessProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.pod.initContainers[].readinessProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.pod.initContainers[].readinessProbe.httpGet.portNumber

`int32`

### spec.pod.initContainers[].readinessProbe.httpGet.portName

`string`

### spec.pod.initContainers[].readinessProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.pod.initContainers[].readinessProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.pod.initContainers[].readinessProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.pod.initContainers[].readinessProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.pod.initContainers[].readinessProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.pod.initContainers[].readinessProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.pod.initContainers[].readinessProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.pod.initContainers[].readinessProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.pod.initContainers[].readinessProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.pod.initContainers[].readinessProbe.tcpSocket.portNumber

`int32`

### spec.pod.initContainers[].readinessProbe.tcpSocket.portName

`string`

### spec.pod.initContainers[].readinessProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.pod.initContainers[].readinessProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.pod.initContainers[].readinessProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.pod.initContainers[].startupProbe

`Probe`

Startup probe: holds off liveness and readiness checking until the app has
started, so slow-booting applications are not killed mid-initialization. Size
`failure_threshold × period_seconds` to the worst-case startup time.

### spec.pod.initContainers[].startupProbe.initialDelaySeconds

`int32`

Number of seconds after the container has started before liveness or readiness probes are initiated.
Defaults to 0 seconds. Minimum value is 0.

### spec.pod.initContainers[].startupProbe.periodSeconds

`int32`

How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.

### spec.pod.initContainers[].startupProbe.timeoutSeconds

`int32`

Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.

### spec.pod.initContainers[].startupProbe.successThreshold

`int32`

Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.

### spec.pod.initContainers[].startupProbe.failureThreshold

`int32`

Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.

### spec.pod.initContainers[].startupProbe.httpGet

`HTTPGetAction`

HTTPGet specifies the http request to perform.

### spec.pod.initContainers[].startupProbe.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.pod.initContainers[].startupProbe.httpGet.portNumber

`int32`

### spec.pod.initContainers[].startupProbe.httpGet.portName

`string`

### spec.pod.initContainers[].startupProbe.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.pod.initContainers[].startupProbe.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.pod.initContainers[].startupProbe.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.pod.initContainers[].startupProbe.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.pod.initContainers[].startupProbe.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.pod.initContainers[].startupProbe.grpc

`GRPCAction`

GRPC specifies an action involving a GRPC port.

### spec.pod.initContainers[].startupProbe.grpc.port

`int32`

Port number of the gRPC service.
Number must be in the range 1 to 65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.pod.initContainers[].startupProbe.grpc.service

`string`

Service is the name of the service to check.
If not specified, the default behavior defined by gRPC is used.
For standard gRPC health checks, leave empty to check overall server health.

### spec.pod.initContainers[].startupProbe.tcpSocket

`TCPSocketAction`

TCPSocket specifies an action involving a TCP port.

### spec.pod.initContainers[].startupProbe.tcpSocket.portNumber

`int32`

### spec.pod.initContainers[].startupProbe.tcpSocket.portName

`string`

### spec.pod.initContainers[].startupProbe.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.pod.initContainers[].startupProbe.exec

`ExecAction`

Exec specifies a command to execute inside the container.

### spec.pod.initContainers[].startupProbe.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.pod.initContainers[].volumeMounts

`[]VolumeMount`

Volume mounts for this container. Each entry both declares the mount path and
carries its volume source (ConfigMap, Secret, HostPath, EmptyDir, or PVC); the
module derives the pod-level volume list from the union of all containers'
mounts, de-duplicating by name — so two containers sharing an EmptyDir simply
declare the same mount name and source.

### spec.pod.initContainers[].volumeMounts[].name

`string` · required

Name of the volume mount. Must be unique within the container.
Used to correlate with the volume definition.

- rule: {"required":true}

### spec.pod.initContainers[].volumeMounts[].mountPath

`string` · required

Path within the container at which the volume should be mounted.
Must be an absolute path.

- rule: {"required":true}

### spec.pod.initContainers[].volumeMounts[].readOnly

`bool`

Whether the volume should be mounted read-only.
Default is false.

### spec.pod.initContainers[].volumeMounts[].subPath

`string`

Path within the volume from which the container's volume should be mounted.
Defaults to "" (volume's root).
Useful for mounting a subdirectory of a volume.

### spec.pod.initContainers[].volumeMounts[].configMap

`ConfigMapVolumeSource`

ConfigMap volume source.
Use this to mount a ConfigMap as a file or directory.

### spec.pod.initContainers[].volumeMounts[].configMap.name

`string` · required

Name of the ConfigMap to mount.
Can reference a ConfigMap defined in spec.config_maps or an existing one in the namespace.

- rule: {"required":true}

### spec.pod.initContainers[].volumeMounts[].configMap.key

`string`

Specific key from the ConfigMap to mount as a single file.
If not specified, all keys are mounted as files in the directory.

### spec.pod.initContainers[].volumeMounts[].configMap.path

`string`

If key is specified, this is the filename to use for the mounted file.
Defaults to the key name if not specified.
Example: key="config" path="app.yaml" mounts the "config" key as "app.yaml"

### spec.pod.initContainers[].volumeMounts[].configMap.defaultMode

`int32`

Mode bits to use on created files. Must be a value between 0 and 0777.
Defaults to 0644.
Use 0755 (493 in decimal) for executable scripts.

### spec.pod.initContainers[].volumeMounts[].secret

`SecretVolumeSource`

Secret volume source.
Use this to mount a Secret as a file or directory.

### spec.pod.initContainers[].volumeMounts[].secret.name

`string` · required

Name of the Secret to mount.

- rule: {"required":true}

### spec.pod.initContainers[].volumeMounts[].secret.key

`string`

Specific key from the Secret to mount as a single file.
If not specified, all keys are mounted as files in the directory.

### spec.pod.initContainers[].volumeMounts[].secret.path

`string`

If key is specified, this is the filename to use for the mounted file.
Defaults to the key name if not specified.

### spec.pod.initContainers[].volumeMounts[].secret.defaultMode

`int32`

Mode bits to use on created files. Must be a value between 0 and 0777.
Defaults to 0644.

### spec.pod.initContainers[].volumeMounts[].hostPath

`HostPathVolumeSource`

HostPath volume source.
Use this to mount a file or directory from the host node's filesystem.
Common for DaemonSets that need access to node-level resources.

### spec.pod.initContainers[].volumeMounts[].hostPath.path

`string` · required

Path on the host to mount.

- rule: {"required":true}

### spec.pod.initContainers[].volumeMounts[].hostPath.type

`string`

Type of the host path.
Valid values:
  "" - Empty string (default) means no check is performed before mounting
  "DirectoryOrCreate" - Create directory if it doesn't exist
  "Directory" - Directory must exist
  "FileOrCreate" - Create file if it doesn't exist
  "File" - File must exist
  "Socket" - UNIX socket must exist
  "CharDevice" - Character device must exist
  "BlockDevice" - Block device must exist

- rule: Type must be one of: "", "DirectoryOrCreate", "Directory", "FileOrCreate", "File", "Socket", "CharDevice", "BlockDevice"

### spec.pod.initContainers[].volumeMounts[].emptyDir

`EmptyDirVolumeSource`

EmptyDir volume source.
Use this for temporary storage that is erased when the pod is removed.
Useful for scratch space, caching, or sharing data between containers.

### spec.pod.initContainers[].volumeMounts[].emptyDir.medium

`string`

Medium for the empty directory.
"" (default) uses the node's default medium (typically disk).
"Memory" uses a tmpfs (RAM-backed filesystem).

Memory-backed volumes are faster but:
- Count against container memory limits
- Are lost on node restart
- Should have sizeLimit set to prevent OOM

- rule: Medium must be either "" or "Memory"

### spec.pod.initContainers[].volumeMounts[].emptyDir.sizeLimit

`string`

Size limit for the empty directory.
Format: Kubernetes quantity (e.g., "1Gi", "500Mi").
Only strictly enforced when medium is "Memory".
For disk-backed volumes, this is a best-effort limit.

### spec.pod.initContainers[].volumeMounts[].pvc

`PvcVolumeSource`

PersistentVolumeClaim volume source.
Use this to mount an existing PVC.
For StatefulSets, this can reference a volumeClaimTemplate.

### spec.pod.initContainers[].volumeMounts[].pvc.claimName

`string` · required

Name of the PersistentVolumeClaim to mount.
For StatefulSets, this can be the name of a volumeClaimTemplate.

- rule: {"required":true}

### spec.pod.initContainers[].volumeMounts[].pvc.readOnly

`bool`

Whether the PVC should be mounted read-only.
Default is false.

### spec.pod.initContainers[].lifecycle

`WorkloadContainerLifecycle`

Lifecycle hooks. `post_start` runs immediately after the container starts (the
container is not Running until it completes); `pre_stop` runs before termination
and is the standard lever for connection draining — e.g. a short sleep that keeps
the endpoint serving while load balancers converge on the terminating state.

### spec.pod.initContainers[].lifecycle.postStart

`WorkloadLifecycleHandler`

Runs immediately after the container is created. The container does not reach
Running until the hook completes; a failing post_start kills the container per
its restart policy.

### spec.pod.initContainers[].lifecycle.postStart.exec

`ExecAction`

Execute a command inside the container; non-zero exit fails the hook.

### spec.pod.initContainers[].lifecycle.postStart.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.pod.initContainers[].lifecycle.postStart.httpGet

`HTTPGetAction`

Perform an HTTP GET against the container.

### spec.pod.initContainers[].lifecycle.postStart.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.pod.initContainers[].lifecycle.postStart.httpGet.portNumber

`int32`

### spec.pod.initContainers[].lifecycle.postStart.httpGet.portName

`string`

### spec.pod.initContainers[].lifecycle.postStart.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.pod.initContainers[].lifecycle.postStart.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.pod.initContainers[].lifecycle.postStart.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.pod.initContainers[].lifecycle.postStart.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.pod.initContainers[].lifecycle.postStart.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.pod.initContainers[].lifecycle.postStart.tcpSocket

`TCPSocketAction`

Open a TCP connection against the container.

### spec.pod.initContainers[].lifecycle.postStart.tcpSocket.portNumber

`int32`

### spec.pod.initContainers[].lifecycle.postStart.tcpSocket.portName

`string`

### spec.pod.initContainers[].lifecycle.postStart.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.pod.initContainers[].lifecycle.postStart.sleep

`SleepAction`

Pause for a fixed number of seconds — the idiomatic pre_stop drain hook,
with no shell or sleep binary required in the image.

### spec.pod.initContainers[].lifecycle.postStart.sleep.seconds

`int64`

Seconds to sleep.

- rule: {"int64":{"gte":"0"}}

### spec.pod.initContainers[].lifecycle.preStop

`WorkloadLifecycleHandler`

Runs before the container is terminated by the kubelet (pod deletion, rolling
update, eviction). The termination grace period starts BEFORE the hook runs, so
keep `pod.termination_grace_period_seconds` larger than the hook's worst-case
duration. The classic zero-downtime pattern is a short sleep here so the pod
keeps serving while endpoint removal propagates.

### spec.pod.initContainers[].lifecycle.preStop.exec

`ExecAction`

Execute a command inside the container; non-zero exit fails the hook.

### spec.pod.initContainers[].lifecycle.preStop.exec.command

`[]string`

Command is the command line to execute inside the container.
The command is run in the container's root filesystem.
The command's exit status is used to determine the health:
- 0: Success
- Non-zero: Failure

### spec.pod.initContainers[].lifecycle.preStop.httpGet

`HTTPGetAction`

Perform an HTTP GET against the container.

### spec.pod.initContainers[].lifecycle.preStop.httpGet.path

`string`

Path to access on the HTTP server.
Defaults to '/'.

### spec.pod.initContainers[].lifecycle.preStop.httpGet.portNumber

`int32`

### spec.pod.initContainers[].lifecycle.preStop.httpGet.portName

`string`

### spec.pod.initContainers[].lifecycle.preStop.httpGet.host

`string`

Host name to connect to, defaults to the pod IP.
You probably want to set "Host" in http_headers instead.

### spec.pod.initContainers[].lifecycle.preStop.httpGet.scheme

`string`

Scheme to use for connecting to the host (HTTP or HTTPS).
Defaults to HTTP.

### spec.pod.initContainers[].lifecycle.preStop.httpGet.httpHeaders

`[]HTTPHeader`

Custom headers to set in the request.
HTTP allows repeated headers.

### spec.pod.initContainers[].lifecycle.preStop.httpGet.httpHeaders[].name

`string`

The header field name.

### spec.pod.initContainers[].lifecycle.preStop.httpGet.httpHeaders[].value

`string`

The header field value.

### spec.pod.initContainers[].lifecycle.preStop.tcpSocket

`TCPSocketAction`

Open a TCP connection against the container.

### spec.pod.initContainers[].lifecycle.preStop.tcpSocket.portNumber

`int32`

### spec.pod.initContainers[].lifecycle.preStop.tcpSocket.portName

`string`

### spec.pod.initContainers[].lifecycle.preStop.tcpSocket.host

`string`

Host name to connect to, defaults to the pod IP.

### spec.pod.initContainers[].lifecycle.preStop.sleep

`SleepAction`

Pause for a fixed number of seconds — the idiomatic pre_stop drain hook,
with no shell or sleep binary required in the image.

### spec.pod.initContainers[].lifecycle.preStop.sleep.seconds

`int64`

Seconds to sleep.

- rule: {"int64":{"gte":"0"}}

### spec.pod.initContainers[].securityContext

`WorkloadContainerSecurityContext`

Container-level security hardening. Settings here override the pod-level
security context for this container only.

### spec.pod.initContainers[].securityContext.privileged

`bool`

Runs the container with full host access — equivalent to root on the node.
Required by some node-level agents (device managers, network plugins). Never
combine with untrusted images.

### spec.pod.initContainers[].securityContext.runAsUser

`int64` · optional (explicit presence)

UID the container process runs as. Overrides the image's USER directive.

### spec.pod.initContainers[].securityContext.runAsGroup

`int64` · optional (explicit presence)

Primary GID the container process runs as.

### spec.pod.initContainers[].securityContext.runAsNonRoot

`bool` · optional (explicit presence)

Refuses to start the container if its effective user is root. The standard
baseline hardening — it catches images that silently default to UID 0.

### spec.pod.initContainers[].securityContext.readOnlyRootFilesystem

`bool` · optional (explicit presence)

Mounts the container's root filesystem read-only. Pair with EmptyDir mounts for
paths the app must write (e.g. /tmp).

### spec.pod.initContainers[].securityContext.allowPrivilegeEscalation

`bool` · optional (explicit presence)

Whether the process can gain more privileges than its parent (setuid binaries,
file capabilities). The restricted Pod Security Standard requires this to be
false. Always true when `privileged` is set, so leave it unset in that case.

### spec.pod.initContainers[].securityContext.capabilities

`WorkloadCapabilities`

Linux capabilities to add or drop. The restricted profile drops ALL and adds
back only NET_BIND_SERVICE when needed. Capability names are uppercase without
the CAP_ prefix (e.g. "NET_ADMIN", "SYS_TIME").

### spec.pod.initContainers[].securityContext.capabilities.add

`[]string`

Capabilities to add (e.g. "NET_BIND_SERVICE").

### spec.pod.initContainers[].securityContext.capabilities.drop

`[]string`

Capabilities to drop. Use ["ALL"] as the hardened baseline.

### spec.pod.initContainers[].securityContext.seccompProfile

`WorkloadSeccompProfile`

Seccomp syscall filter for the container. "RuntimeDefault" is the hardened
baseline; "Localhost" selects a node-local profile file via `localhost_profile`.

- rule: localhost_profile is required when type is "Localhost" and must be empty otherwise

### spec.pod.initContainers[].securityContext.seccompProfile.type

`string` · required

Profile type: "RuntimeDefault" (the container runtime's default filter — the
recommended baseline), "Unconfined" (no filtering), or "Localhost" (a profile
file installed on the node, named via localhost_profile).

- rule: Seccomp profile type must be one of "RuntimeDefault", "Unconfined", or "Localhost"
- rule: {"required":true}

### spec.pod.initContainers[].securityContext.seccompProfile.localhostProfile

`string`

Path of the profile file relative to the node's seccomp profile root. Required
when (and only meaningful when) type is "Localhost".

### spec.pod.labels

`map<string, string>`

Extra labels stamped on the POD TEMPLATE (and therefore every pod), merged with
the workload's own selector and governance labels. This is where cross-cutting
markers pods must carry go — e.g. `azure.workload.identity/use: "true"` for AKS
workload identity, or mesh sidecar-injection toggles. Keys here must not collide
with the selector labels the module derives.

### spec.pod.annotations

`map<string, string>`

Extra annotations stamped on the pod template — e.g. prometheus.io scrape hints
or mesh proxy tuning. Changing pod-template annotations rolls pods, which is
also the idiomatic config-reload lever.

### spec.pod.scheduling

`WorkloadScheduling`

Where pods may schedule: node selection, taint tolerations, affinity rules, and
zone/host spreading. Omit to let the scheduler place pods anywhere.

### spec.pod.scheduling.nodeSelector

`map<string, string>`

Simple hard node filter: every listed label must match the node. The right tool
for "run on the GPU pool" — reach for node_affinity only when you need operators
(In/NotIn/Exists) or soft preferences.

### spec.pod.scheduling.tolerations

`[]WorkloadToleration`

Taint tolerations. A toleration does not attract pods to tainted nodes — it only
permits scheduling there; pair with node_selector or affinity to target them.

### spec.pod.scheduling.tolerations[].key

`string`

Taint key to tolerate. Empty key with operator "Exists" tolerates every taint.

### spec.pod.scheduling.tolerations[].operator

`string`

How key/value match: "Equal" (default — value must match too) or "Exists"
(key presence alone matches).

- rule: Toleration operator must be either "Equal" or "Exists"

### spec.pod.scheduling.tolerations[].value

`string`

Taint value to match when operator is "Equal".

### spec.pod.scheduling.tolerations[].effect

`string`

Which taint effect is tolerated: "NoSchedule", "PreferNoSchedule", or
"NoExecute". Empty tolerates all effects for the key.

- rule: Toleration effect must be one of "NoSchedule", "PreferNoSchedule", or "NoExecute"

### spec.pod.scheduling.tolerations[].tolerationSeconds

`int64` · optional (explicit presence)

For "NoExecute" taints only: how many seconds already-running pods stay bound
after the taint appears. Unset means tolerate forever.

### spec.pod.scheduling.nodeAffinity

`WorkloadNodeAffinity`

Expressive node selection: hard requirements and weighted soft preferences over
node labels.

### spec.pod.scheduling.nodeAffinity.required

`[]WorkloadNodeSelectorTerm`

Hard requirement. The outer list ORs its terms; expressions within one term AND.

### spec.pod.scheduling.nodeAffinity.required[].matchExpressions

`[]WorkloadNodeSelectorRequirement` · required

- rule: {"repeated":{"minItems":"1"}}
- rule: In/NotIn require at least one value, Gt/Lt exactly one, Exists/DoesNotExist none

### spec.pod.scheduling.nodeAffinity.required[].matchExpressions[].key

`string` · required

Node label key, e.g. "topology.kubernetes.io/zone".

- rule: {"required":true}

### spec.pod.scheduling.nodeAffinity.required[].matchExpressions[].operator

`string` · required

Operator: "In"/"NotIn" (value set), "Exists"/"DoesNotExist" (key presence), or
"Gt"/"Lt" (single integer value, as strings — the Kubernetes API convention).

- rule: Operator must be one of "In", "NotIn", "Exists", "DoesNotExist", "Gt", or "Lt"
- rule: {"required":true}

### spec.pod.scheduling.nodeAffinity.required[].matchExpressions[].values

`[]string`

Values for the operator: required non-empty for In/NotIn, exactly one integer
string for Gt/Lt, and must be empty for Exists/DoesNotExist.

### spec.pod.scheduling.nodeAffinity.preferred

`[]WorkloadPreferredNodeSelectorTerm`

Weighted soft preferences.

### spec.pod.scheduling.nodeAffinity.preferred[].weight

`int32`

Preference weight, 1–100. Higher weights dominate placement scoring.

- rule: {"int32":{"lte":100,"gte":1}}

### spec.pod.scheduling.nodeAffinity.preferred[].term

`WorkloadNodeSelectorTerm` · required

- rule: {"required":true}

### spec.pod.scheduling.nodeAffinity.preferred[].term.matchExpressions

`[]WorkloadNodeSelectorRequirement` · required

- rule: {"repeated":{"minItems":"1"}}
- rule: In/NotIn require at least one value, Gt/Lt exactly one, Exists/DoesNotExist none

### spec.pod.scheduling.nodeAffinity.preferred[].term.matchExpressions[].key

`string` · required

Node label key, e.g. "topology.kubernetes.io/zone".

- rule: {"required":true}

### spec.pod.scheduling.nodeAffinity.preferred[].term.matchExpressions[].operator

`string` · required

Operator: "In"/"NotIn" (value set), "Exists"/"DoesNotExist" (key presence), or
"Gt"/"Lt" (single integer value, as strings — the Kubernetes API convention).

- rule: Operator must be one of "In", "NotIn", "Exists", "DoesNotExist", "Gt", or "Lt"
- rule: {"required":true}

### spec.pod.scheduling.nodeAffinity.preferred[].term.matchExpressions[].values

`[]string`

Values for the operator: required non-empty for In/NotIn, exactly one integer
string for Gt/Lt, and must be empty for Exists/DoesNotExist.

### spec.pod.scheduling.podAffinity

`WorkloadPodAffinity`

Attract pods toward nodes/zones already running matching pods (co-location with
a cache, for example).

### spec.pod.scheduling.podAffinity.required

`[]WorkloadPodAffinityTerm`

Hard rules — unschedulable until satisfied. Use sparingly; they can deadlock rollouts.

### spec.pod.scheduling.podAffinity.required[].matchLabels

`map<string, string>` · required

Labels of the pods to match against — for self-anti-affinity, the workload's own
selector labels (exported as the `selector_labels` stack output).

- rule: {"map":{"minPairs":"1"}}

### spec.pod.scheduling.podAffinity.required[].topologyKey

`string` · required

Node label defining the domain: "kubernetes.io/hostname" separates by node,
"topology.kubernetes.io/zone" by zone.

- rule: {"required":true}

### spec.pod.scheduling.podAffinity.required[].namespaces

`[]string`

Namespaces whose pods are considered. Empty means the workload's own namespace.

### spec.pod.scheduling.podAffinity.preferred

`[]WorkloadWeightedPodAffinityTerm`

Weighted soft rules — the scheduler's tiebreakers.

### spec.pod.scheduling.podAffinity.preferred[].weight

`int32`

Preference weight, 1–100.

- rule: {"int32":{"lte":100,"gte":1}}

### spec.pod.scheduling.podAffinity.preferred[].term

`WorkloadPodAffinityTerm` · required

- rule: {"required":true}

### spec.pod.scheduling.podAffinity.preferred[].term.matchLabels

`map<string, string>` · required

Labels of the pods to match against — for self-anti-affinity, the workload's own
selector labels (exported as the `selector_labels` stack output).

- rule: {"map":{"minPairs":"1"}}

### spec.pod.scheduling.podAffinity.preferred[].term.topologyKey

`string` · required

Node label defining the domain: "kubernetes.io/hostname" separates by node,
"topology.kubernetes.io/zone" by zone.

- rule: {"required":true}

### spec.pod.scheduling.podAffinity.preferred[].term.namespaces

`[]string`

Namespaces whose pods are considered. Empty means the workload's own namespace.

### spec.pod.scheduling.podAntiAffinity

`WorkloadPodAffinity`

Repel pods from nodes/zones already running matching pods — the classic
high-availability pattern is anti-affinity on the workload's own labels across
`kubernetes.io/hostname`.

### spec.pod.scheduling.podAntiAffinity.required

`[]WorkloadPodAffinityTerm`

Hard rules — unschedulable until satisfied. Use sparingly; they can deadlock rollouts.

### spec.pod.scheduling.podAntiAffinity.required[].matchLabels

`map<string, string>` · required

Labels of the pods to match against — for self-anti-affinity, the workload's own
selector labels (exported as the `selector_labels` stack output).

- rule: {"map":{"minPairs":"1"}}

### spec.pod.scheduling.podAntiAffinity.required[].topologyKey

`string` · required

Node label defining the domain: "kubernetes.io/hostname" separates by node,
"topology.kubernetes.io/zone" by zone.

- rule: {"required":true}

### spec.pod.scheduling.podAntiAffinity.required[].namespaces

`[]string`

Namespaces whose pods are considered. Empty means the workload's own namespace.

### spec.pod.scheduling.podAntiAffinity.preferred

`[]WorkloadWeightedPodAffinityTerm`

Weighted soft rules — the scheduler's tiebreakers.

### spec.pod.scheduling.podAntiAffinity.preferred[].weight

`int32`

Preference weight, 1–100.

- rule: {"int32":{"lte":100,"gte":1}}

### spec.pod.scheduling.podAntiAffinity.preferred[].term

`WorkloadPodAffinityTerm` · required

- rule: {"required":true}

### spec.pod.scheduling.podAntiAffinity.preferred[].term.matchLabels

`map<string, string>` · required

Labels of the pods to match against — for self-anti-affinity, the workload's own
selector labels (exported as the `selector_labels` stack output).

- rule: {"map":{"minPairs":"1"}}

### spec.pod.scheduling.podAntiAffinity.preferred[].term.topologyKey

`string` · required

Node label defining the domain: "kubernetes.io/hostname" separates by node,
"topology.kubernetes.io/zone" by zone.

- rule: {"required":true}

### spec.pod.scheduling.podAntiAffinity.preferred[].term.namespaces

`[]string`

Namespaces whose pods are considered. Empty means the workload's own namespace.

### spec.pod.scheduling.topologySpreadConstraints

`[]WorkloadTopologySpreadConstraint`

Even distribution of replicas across topology domains (zones, hosts). Preferred
over hostname anti-affinity for large replica counts because skew is bounded
rather than binary.

### spec.pod.scheduling.topologySpreadConstraints[].maxSkew

`int32`

Maximum allowed difference in matching-pod counts between any two domains.
1 is the strictest even spread.

- rule: {"int32":{"gte":1}}

### spec.pod.scheduling.topologySpreadConstraints[].topologyKey

`string` · required

Node label defining the domains to spread across (e.g.
"topology.kubernetes.io/zone").

- rule: {"required":true}

### spec.pod.scheduling.topologySpreadConstraints[].whenUnsatisfiable

`string` · required

What happens when the constraint cannot be met: "DoNotSchedule" (hard — pod
stays Pending) or "ScheduleAnyway" (soft — scheduler minimizes skew).

- rule: whenUnsatisfiable must be either "DoNotSchedule" or "ScheduleAnyway"
- rule: {"required":true}

### spec.pod.scheduling.topologySpreadConstraints[].matchLabels

`map<string, string>`

Labels selecting the pods counted per domain. Omit to have the module default to
the workload's own selector labels — self-spreading, the overwhelmingly common
intent.

### spec.pod.scheduling.schedulerName

`string`

Hand pods to a non-default scheduler installed in the cluster. Leave empty for
the standard scheduler.

### spec.pod.securityContext

`WorkloadPodSecurityContext`

Pod-level security context: the user/group identity and filesystem ownership
every container inherits unless it overrides them in its own security context.

### spec.pod.securityContext.runAsUser

`int64` · optional (explicit presence)

UID all container processes run as unless overridden per container.

### spec.pod.securityContext.runAsGroup

`int64` · optional (explicit presence)

Primary GID all container processes run as unless overridden per container.

### spec.pod.securityContext.runAsNonRoot

`bool` · optional (explicit presence)

Refuse to start any container whose effective user is root.

### spec.pod.securityContext.fsGroup

`int64` · optional (explicit presence)

GID that owns mounted volumes and is added to every container's supplemental
groups — the standard fix for "permission denied" on persistent volumes written
by non-root apps.

### spec.pod.securityContext.fsGroupChangePolicy

`string`

When volume ownership is re-chowned to fs_group: "Always" (default) or
"OnRootMismatch" (skip the recursive chown when the root already matches —
dramatically faster pod starts on large volumes).

- rule: fsGroupChangePolicy must be either "Always" or "OnRootMismatch"

### spec.pod.securityContext.supplementalGroups

`[]int64`

Additional group IDs applied to all container processes.

### spec.pod.securityContext.sysctls

`[]WorkloadSysctl`

Kernel parameters set for the pod. Only safe sysctls (or those the cluster
administrator has allow-listed on the kubelet) are admitted.

### spec.pod.securityContext.sysctls[].name

`string` · required

Sysctl name, e.g. "net.core.somaxconn".

- rule: {"required":true}

### spec.pod.securityContext.sysctls[].value

`string` · required

Sysctl value, e.g. "1024".

- rule: {"required":true}

### spec.pod.securityContext.seccompProfile

`WorkloadSeccompProfile`

Pod-wide seccomp profile; containers may override with their own.

- rule: localhost_profile is required when type is "Localhost" and must be empty otherwise

### spec.pod.securityContext.seccompProfile.type

`string` · required

Profile type: "RuntimeDefault" (the container runtime's default filter — the
recommended baseline), "Unconfined" (no filtering), or "Localhost" (a profile
file installed on the node, named via localhost_profile).

- rule: Seccomp profile type must be one of "RuntimeDefault", "Unconfined", or "Localhost"
- rule: {"required":true}

### spec.pod.securityContext.seccompProfile.localhostProfile

`string`

Path of the profile file relative to the node's seccomp profile root. Required
when (and only meaningful when) type is "Localhost".

### spec.pod.terminationGracePeriodSeconds

`int64` · optional (explicit presence)

Seconds the kubelet waits between SIGTERM and SIGKILL at pod termination.
Kubernetes defaults to 30 when unset. Size it to cover pre_stop hooks plus the
app's own drain time — the grace clock starts before the hook runs.

- rule: {"int64":{"gte":"0"}}

### spec.pod.dnsPolicy

`string`

Pod DNS resolution policy. "ClusterFirst" (the Kubernetes default) resolves
cluster services first; "Default" inherits the node's resolver;
"ClusterFirstWithHostNet" is what host-network pods need to keep resolving
cluster services; "None" hands control entirely to `dns_config`.

- rule: DNS policy must be one of "ClusterFirst", "ClusterFirstWithHostNet", "Default", or "None"

### spec.pod.dnsConfig

`WorkloadPodDnsConfig`

Custom DNS parameters merged into (or, with dns_policy "None", replacing) the
generated resolv.conf.

### spec.pod.dnsConfig.nameservers

`[]string`

Nameserver IPs (max 3 total after merging with the policy's servers).

- rule: {"repeated":{"maxItems":"3"}}

### spec.pod.dnsConfig.searches

`[]string`

Search domains for hostname lookup.

### spec.pod.dnsConfig.options

`[]WorkloadPodDnsConfigOption`

resolv.conf options, e.g. name "ndots" value "2". Value may be empty for flags.

### spec.pod.dnsConfig.options[].name

`string` · required

- rule: {"required":true}

### spec.pod.dnsConfig.options[].value

`string`

### spec.pod.hostAliases

`[]WorkloadHostAlias`

Static host-file entries injected into every container's /etc/hosts — the
escape hatch for names that no reachable DNS serves.

### spec.pod.hostAliases[].ip

`string` · required

The IP the hostnames resolve to.

- rule: {"required":true,"string":{"ip":true}}

### spec.pod.hostAliases[].hostnames

`[]string` · required

Hostnames mapped to the IP.

- rule: {"repeated":{"minItems":"1"}}

### spec.pod.hostNetwork

`bool`

Run pods in the node's network namespace. A DaemonSet agent pattern (node
monitoring, CNI components). Host-network pods share the node's port space —
combine with `dns_policy: ClusterFirstWithHostNet` if they resolve cluster
services.

### spec.pod.hostPid

`bool`

Share the node's PID namespace (process-visibility agents). Requires elevated
trust in the workload.

### spec.pod.priorityClassName

`string`

PriorityClass name controlling scheduling and eviction precedence. References a
cluster-scoped PriorityClass (e.g. "system-cluster-critical" or one created by a
KubernetesPriorityClass resource).

### spec.pod.runtimeClassName

`string`

RuntimeClass selecting an alternate container runtime configuration (e.g. gVisor
or Kata sandboxing) installed by the cluster administrator.

### spec.availability

`KubernetesStatefulSetAvailability`

Availability configuration: replica count, disruption budget, and rollout
safeguards. Omit for a single-replica set with Kubernetes defaults. StatefulSets
deliberately have no autoscaling here — stateful members join and leave through
application-aware procedures, not HPA.

### spec.availability.replicas

`int32` · optional (explicit presence)

Desired replica count. Defaults to 1. Scaling a stateful system is not free:
new members must sync data, and removed members' PVCs follow
`pvc_retention_policy.when_scaled`.

- default: `1`
- rule: {"int32":{"gte":0}}

### spec.availability.podDisruptionBudget

`KubernetesStatefulSetPodDisruptionBudget`

PodDisruptionBudget guarding availability during voluntary disruptions (node
drains, cluster upgrades). Especially important for quorum-based systems — e.g.
a 3-member cluster typically sets min_available "2" so a drain can never break
quorum.

- rule: Set exactly one of min_available or max_unavailable (defaults to min_available: "1" when both are empty)

### spec.availability.podDisruptionBudget.enabled

`bool`

Enables PDB creation.

### spec.availability.podDisruptionBudget.minAvailable

`string`

Minimum pods that must stay available — absolute ("2") or percentage ("50%").
For a quorum system of N members, set the quorum size (e.g. "2" for a 3-member
cluster) so voluntary disruptions can never take availability below it.

### spec.availability.podDisruptionBudget.maxUnavailable

`string`

Maximum pods that may be down — absolute or percentage. Alternative to
min_available; do not set both.

### spec.availability.minReadySeconds

`int32` · optional (explicit presence)

Seconds a new pod must stay ready (no container crashes) before it counts as
available during rollouts. A cheap flap detector — 10–30s catches crash-on-first-
request regressions before the rollout proceeds to the next ordinal.

- rule: {"int32":{"gte":0}}

### spec.availability.revisionHistoryLimit

`int32` · optional (explicit presence)

How many old ControllerRevisions are retained for rollback. Kubernetes defaults
to 10; 0 disables rollback entirely.

- rule: {"int32":{"gte":0}}

### spec.volumeClaimTemplates

`[]KubernetesStatefulSetVolumeClaimTemplate`

PersistentVolumeClaim templates. Each replica gets its OWN PVC stamped from each
template, named <template-name>-<statefulset-name>-<ordinal>, and keeps it across
pod restarts and rescheduling — this is what makes the storage stateful. Mount a
template in a container via `volume_mounts` with a PVC source whose claim name is
the template's name.

### spec.volumeClaimTemplates[].name

`string` · required

Template name — the claim name containers reference in their PVC volume mounts.
Must be a lowercase DNS label.

- rule: Volume claim template name must be a lowercase DNS label (alphanumeric and hyphens, starting and ending with an alphanumeric character, e.g. "data")
- rule: {"required":true}

### spec.volumeClaimTemplates[].storageClass

`string`

StorageClass provisioning the volumes. Empty uses the cluster's default
StorageClass. Reference a class created by a KubernetesStorageClass resource to
pin performance characteristics (e.g. SSD-backed, expandable) explicitly.

### spec.volumeClaimTemplates[].size

`string` · required

Requested storage per replica, as a Kubernetes quantity (e.g. "10Gi", "500Mi").
Growing later requires a StorageClass with volume expansion enabled.

- rule: Size must be a Kubernetes quantity (e.g. "10Gi", "500Mi")
- rule: {"required":true}

### spec.volumeClaimTemplates[].accessModes

`[]string`

How the volume may be mounted. Defaults to ["ReadWriteOnce"] — one node at a
time, which suits per-replica stateful storage; the other modes require a
storage driver that supports them ("ReadWriteMany" for shared filesystems,
"ReadWriteOncePod" to pin access to a single pod).

- rule: Each access mode must be one of "ReadWriteOnce", "ReadOnlyMany", "ReadWriteMany", or "ReadWriteOncePod"

### spec.volumeClaimTemplates[].volumeMode

`string`

"Filesystem" (default) mounts a formatted filesystem at the mount path; "Block"
hands the container the raw block device — only for applications that manage
their own on-disk layout.

- rule: Volume mode must be either "Filesystem" or "Block"

### spec.updateStrategy

`KubernetesStatefulSetUpdateStrategy`

How pods are replaced when the spec changes. Omit for the Kubernetes default:
RollingUpdate, replacing pods one at a time from the highest ordinal down.

### spec.updateStrategy.type

`string`

"RollingUpdate" (default) replaces pods in reverse ordinal order, waiting for
each replacement to become ready; "OnDelete" only recreates a pod (with the new
template) when someone deletes it — full manual control, used when each member
needs operator-driven update steps.

- rule: Update strategy type must be either "RollingUpdate" or "OnDelete"

### spec.updateStrategy.partition

`int32` · optional (explicit presence)

Canary-by-ordinal: during a rolling update only pods with ordinal >= partition
are updated; pods below it keep the old template. E.g. with 5 replicas and
partition 4, only pod -4 gets the new version — validate it, then lower the
partition step by step (finally to 0) to roll the rest. Defaults to 0 (update
everything).

- rule: {"int32":{"gte":0}}

### spec.updateStrategy.maxUnavailable

`string`

Maximum pods that may be unavailable during a rolling update — absolute ("2") or
percentage ("10%"). Values above 1 update multiple ordinals in parallel. Requires
the MaxUnavailableStatefulSet feature gate, which is not enabled on every
cluster — leave empty for the universally supported one-at-a-time default, and
note it has little effect under the OrderedReady pod management policy.

### spec.podManagementPolicy

`string`

How pods are created and deleted during scale operations. "OrderedReady" (the
Kubernetes default) proceeds one pod at a time, waiting for each to be ready —
what most clustered systems need for safe bootstrap; "Parallel" launches and
deletes all pods at once — faster for systems that coordinate membership
themselves (e.g. Cassandra). Does not affect update ordering.

- rule: Pod management policy must be either "OrderedReady" or "Parallel"

### spec.pvcRetentionPolicy

`KubernetesStatefulSetPvcRetentionPolicy`

What happens to the PVCs stamped from `volume_claim_templates` when the
StatefulSet is deleted or scaled down. Omit for the Kubernetes default: retain
everything, requiring manual PVC cleanup.

### spec.pvcRetentionPolicy.whenDeleted

`string`

What happens to the PVCs when the StatefulSet itself is deleted. "Retain"
(default) keeps them — re-creating the StatefulSet with the same name re-adopts
the data; "Delete" removes them with the workload — irreversible data loss unless
the data is replicated or backed up elsewhere.

- rule: when_deleted must be either "Retain" or "Delete"

### spec.pvcRetentionPolicy.whenScaled

`string`

What happens to the PVCs of excess replicas when the StatefulSet is scaled down.
"Retain" (default) keeps them, so scaling back up rejoins members with their data
intact — the safe choice for databases; "Delete" removes them on scale-down, so a
later scale-up starts those members empty (they must re-sync from peers).

- rule: when_scaled must be either "Retain" or "Delete"

### spec.ordinals

`KubernetesStatefulSetOrdinals`

Replica ordinal numbering. Omit for the default 0-based ordinals
(<name>-0, <name>-1, ...).

### spec.ordinals.start

`int32` · optional (explicit presence)

The first replica's index. Defaults to 0 (pods <name>-0 through <name>-(N-1)).
Set to number replicas from an alternate base (e.g. 1-indexed to match an
existing convention) or to migrate replicas between StatefulSets by carving up
disjoint ordinal ranges.

- rule: {"int32":{"gte":0}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesStatefulSet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | The namespace the workload was deployed into. |
| `status.outputs.stateful_set_name` | `string` | The name of the StatefulSet object as created in the cluster. |
| `status.outputs.service` | `string` | The headless governing Service of the StatefulSet — the Service that gives each replica its stable per-pod DNS name. Load-balanced client access also goes through this name. |
| `status.outputs.selector_labels` | `string` | The pod selector labels as a "k=v,k=v" string — the exact labels the Service selects on, ready for NetworkPolicy podSelectors, `kubectl get pods -l`, and pod-affinity terms in sibling workloads. |
| `status.outputs.port_forward_command` | `string` | Ready-to-run port-forward command for reaching the workload from a developer machine without any external exposure. ex: kubectl port-forward -n my-ns service/my-db 5432:5432 |
| `status.outputs.kube_endpoint` | `string` | In-cluster DNS endpoint of the Service — the handle exposure kinds (KubernetesIngress, KubernetesHttpRoute and the other Gateway API route kinds) and sibling workloads connect to. ex: my-db.my-ns.svc.cluster.local |
| `status.outputs.pod_dns_template` | `string` | Template for each replica's stable DNS name: "<name>-<ordinal>.<service>.<namespace>.svc.cluster.local". Substitute the ordinal to address a specific member — e.g. replica 0 of a StatefulSet "my-db" in namespace "my-ns" is "my-db-0.my-db.my-ns.svc.cluster.local". This is how clustered clients build their member lists. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.pod.serviceAccount` | KubernetesServiceAccount | `status.outputs.service_account_name` |
| `spec.pod.imagePullSecrets` | KubernetesSecret | `spec.name` |

## See Also

- [Overview](../README.md)
