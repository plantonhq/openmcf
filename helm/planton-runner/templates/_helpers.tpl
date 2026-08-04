{{/*
Expand the name of the chart.
*/}}
{{- define "planton-runner.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "planton-runner.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "planton-runner.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels.
*/}}
{{- define "planton-runner.labels" -}}
helm.sh/chart: {{ include "planton-runner.chart" . }}
{{ include "planton-runner.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- with .Values.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Selector labels.
*/}}
{{- define "planton-runner.selectorLabels" -}}
app.kubernetes.io/name: {{ include "planton-runner.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use.
*/}}
{{- define "planton-runner.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "planton-runner.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Namespace where build PipelineRuns land and the log streamer watches:
the configured tektonNamespace, or the runner's own namespace when empty
(mirroring the runner binary's own fallback).
*/}}
{{- define "planton-runner.buildNamespace" -}}
{{- default .Release.Namespace .Values.build.tektonNamespace }}
{{- end }}

{{/*
Fail fast at install time when the build block cannot work as configured --
a clear error beats a build worker that silently never polls.
*/}}
{{- define "planton-runner.validateBuild" -}}
{{- if .Values.build.enabled }}
{{- if not (or (eq .Values.runner.executionMode "temporal") (eq .Values.runner.executionMode "dual")) }}
{{- fail "build.enabled requires runner.executionMode \"temporal\" or \"dual\" -- the build worker is a Temporal worker" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Container image with tag.
*/}}
{{- define "planton-runner.image" -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag }}
{{- printf "%s:%s" .Values.image.repository $tag }}
{{- end }}

{{/*
Name of the Kubernetes Secret containing the runner credentials.
Returns the existingSecret name if set, otherwise generates a name from the release.
*/}}
{{- define "planton-runner.credentialsSecretName" -}}
{{- if .Values.credentials.existingSecret }}
{{- .Values.credentials.existingSecret }}
{{- else }}
{{- printf "%s-credentials" (include "planton-runner.fullname" .) }}
{{- end }}
{{- end }}

{{/*
Key within the credentials Secret that holds the JSON content.
Returns the existingSecretKey if an external secret is used, otherwise "credentials.json".
*/}}
{{- define "planton-runner.credentialsSecretKey" -}}
{{- if .Values.credentials.existingSecret }}
{{- default "credentials.json" .Values.credentials.existingSecretKey }}
{{- else }}
{{- "credentials.json" }}
{{- end }}
{{- end }}
