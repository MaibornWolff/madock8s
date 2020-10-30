{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "prometheus_exporter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "prometheus_exporter.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "prometheus_exporter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "prometheus_exporter.labels" -}}
app.kubernetes.io/name: {{ include "prometheus_exporter.name" . }}
helm.sh/chart: {{ include "prometheus_exporter.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
madock8s: {{ .Values.madock8sID | quote }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Common annotations
*/}}
{{- define "exporter.annotations" -}}
madock8s: {{ .Values.madock8sID | quote }}

{{- if .Values.annotations.gitlabExporter }}
madock8s.exporter/gitlabExporter.baseurl: {{ .Values.annotations.gitlabExporter.baseurl | quote }}
madock8s.exporter/gitlabExporter.path: {{ .Values.annotations.gitlabExporter.path | quote }}
madock8s.exporter/gitlabExporter.recursive: {{ .Values.annotations.gitlabExporter.recursive | quote }}
madock8s.exporter/gitlabExporter.ref: {{ .Values.annotations.gitlabExporter.ref | quote }}
madock8s.exporter/gitlabExporter.pattern: {{ .Values.annotations.gitlabExporter.pattern | quote }}
{{- end }}

{{- if .Values.annotations.envExporter }}
madock8s.exporter/envExporter: {{ .Values.annotations.envExporter.deployment | quote }}
{{- end }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "prometheus_exporter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "prometheus_exporter.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}
