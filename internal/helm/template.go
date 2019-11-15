package helm

const helmTemplate string = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller-account
  namespace: {{ namespace }}
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: tiller-manager
  namespace: {{ namespace }}
rules:
- apiGroups: ["", "batch", "extensions", "apps"]
  resources: ["*"]
  verbs: ["*"]
- apiGroups: ["monitoring.coreos.com"]
  resources: ["servicemonitors", "prometheusrules"]
  verbs: ["*"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: tiller-binding
  namespace: {{ namespace }}
subjects:
- kind: ServiceAccount
  name: tiller-account
  namespace: {{ namespace }}
roleRef:
  kind: Role
  name: tiller-manager
  apiGroup: rbac.authorization.k8s.io
`
