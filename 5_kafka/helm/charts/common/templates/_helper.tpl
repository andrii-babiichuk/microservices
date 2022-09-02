{{/*
Імʼя. Обрізається до 63-х символів, через обмеження Kubernetes. `nameOverride` дозволяє перевизначати імʼя.
*/}}
{{- define "common.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Повне імʼя додатку (сервісу). Обрізається до 63-х символів, через обмеження Kubernetes
*/}}
{{- define "common.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{ printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{/* Селектори */}}
{{- define "common.selectorLabels" -}}
app.kubernetes.io/name: {{ include "common.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/* Базові маркери. В даному випадку == селекторам, але як правило мають більше значень. */}}
{{- define "common.labels" -}}
{{ include "common.selectorLabels" . }}
{{- end }}
