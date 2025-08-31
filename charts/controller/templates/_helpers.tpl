{{- define "name" -}}
gardener-extension-provider-dns-cloudflare
{{- end -}}

{{- define "labels.app.key" -}}
app.kubernetes.io/name
{{- end -}}
{{- define "labels.app.value" -}}
{{ include "name" . }}
{{- end -}}

{{- define "labels" -}}
{{ include "labels.app.key" . }}: {{ include "labels.app.value" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "tag" -}}
  {{- if .Values.image.tag }}
  {{- .Values.image.tag }}
  {{- else }}
  {{- .Chart.Version }}
  {{- end }}
{{- end }}

{{- define "image" -}}
  {{- $tag := include "tag" . }}
  {{- if hasPrefix "sha256:" $tag }}
  {{- printf "%s@%s" .Values.image.repository $tag }}
  {{- else }}
  {{- printf "%s:%s" .Values.image.repository $tag }}
  {{- end }}
{{- end }}

{{- define "deploymentversion" -}}
apps/v1
{{- end -}}

{{- define "priorityclassversion" -}}
scheduling.k8s.io/v1
{{- end -}}
