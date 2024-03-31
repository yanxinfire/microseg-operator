{{/*
Return the proper operator image name
*/}}
{{- define "security.operator.image" -}}
{{- $registryName := (include "common.tplvalues.render" ( dict "value" .Values.image.registry "context" $)) -}}
{{- $repositoryName := (include "common.tplvalues.render" ( dict "value" .Values.image.repository "context" $)) -}}
{{- $imageName := .Values.image.name -}}
{{- $tag := .Values.image.tag | toString -}}

{{- printf "%s/%s/%s:%s" $registryName $repositoryName $imageName $tag -}}
{{- end -}}

Expand the name of the chart.
*/}}
{{- define "security.operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "security.operator.fullname" -}}
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
{{- define "security.operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/*
Common labels
*/}}
{{- define "security.operator.labels" -}}
app: {{ include "operator.name" . }}
chart: {{ include "operator.chart" . }}
release: {{ .Release.Name }}
heritage: {{ .Release.Service }}
{{- end -}}
