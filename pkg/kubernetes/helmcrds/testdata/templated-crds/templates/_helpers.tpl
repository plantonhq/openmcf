{{- define "fixture.fullname" -}}
{{- default .Release.Name .Values.fullnameOverride -}}
{{- end -}}
