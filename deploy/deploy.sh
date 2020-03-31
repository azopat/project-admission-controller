#!/bin/bash

# Please update these variables
WEBHOOK_PROJECT=project-admission-validation
MAX_ALLOWED_PROJECT_PER_USER="2"
IMAGE="docker.io/azopat/admission"



CA_BUNDLE=`oc get cm/service-ca -o=jsonpath='{.data.ca-bundle\.crt}'  -n openshift-config-managed | base64`

echo "---
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.openshift.io/serving-cert-secret-name: project-admission-controller
  labels:
    app: project-admission-controller
  name: project-admission-controller-svc
spec:
  ports:
  - name: 443-tcp
    port: 443
    protocol: TCP
    targetPort: 8443
  selector:
    app: project-admission-controller
    deploymentconfig: project-admission-controller
  type: ClusterIP
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: project-admission-controller
  name: project-admission-controller-sa
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: namespace-reader
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: project-admission-webhook
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: namespace-reader
subjects:
- kind: ServiceAccount
  name: project-admission-controller-sa
  namespace: $WEBHOOK_PROJECT
---
apiVersion: apps.openshift.io/v1
kind: DeploymentConfig
metadata:
  labels:
    app: project-admission-controller
  name: project-admission-controller
spec:
  replicas: 1
  selector:
    app: project-admission-controller
    deploymentconfig: project-admission-controller
  template:
    metadata:
      labels:
        app: project-admission-controller
        deploymentconfig: project-admission-controller
    spec:
      containers:
      - env:
        - name: CERT_FOLDER
          value: /app/cert
        - name: MAX_ALLOWED_PROJECT
          value: \"$MAX_ALLOWED_PROJECT_PER_USER\"
        - name: PORT
          value: \"8443\"
        image: docker.io/azopat/admission
        imagePullPolicy: Always
        name: controller
        ports:
        - containerPort: 8443
          protocol: TCP
        resources:
          limits:
            cpu: 250m
            memory: 300Mi
        volumeMounts:
        - mountPath: /app/cert
          name: volume-cert
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      serviceAccount: project-admission-controller-sa
      volumes:
      - name: volume-cert
        secret:
          defaultMode: 420
          secretName: project-admission-controller
  triggers:
  - type: ConfigChange
---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: project-validation-webhook
webhooks:
- admissionReviewVersions:
  - v1beta1
  clientConfig:
    caBundle: $CA_BUNDLE
    service:
      name: project-admission-controller-svc
      namespace: project-admission-validation
      path: /validate-project
      port: 443
  failurePolicy: Ignore
  matchPolicy: Equivalent
  name: project-validation-webhook.cluster.local
  namespaceSelector: {}
  objectSelector: {}
  rules:
  - apiGroups:
    - project.openshift.io
    apiVersions:
    - '*'
    operations:
    - CREATE
    resources:
    - projectrequests
    scope: '*'
  sideEffects: None
  timeoutSeconds: 10" > admission_controller.yml
