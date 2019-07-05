package namespace

const copsNamespaceTemplate string = `
apiVersion: coreops.conplement.cloud/v1
kind: CopsNamespace
metadata:
  name: {{ namespaceName }}
spec:
  namespaceAdminUsers:
{{ usernames }}`
