Re-currant
====

A poor man's gitops operator, which blindly applies kubernetes yaml synced from the git repo.
Supports Kustomize, custom commands to deploy, `/reload` endpoint

# Getting started

Create `kustomization.yaml` (see `manifests/kustomization.yaml` as an example):
```
# Use remote manifests as base
bases:
  - github.com/vrutkovs/re-currant//manifests?ref=master
# Add additional resources, e,g. cluster-admin permissions
resources:
  - cluster-admin.yaml
generatorOptions:
  disableNameSuffixHash: true
secretGenerator:
  - name: recurrant
    namespace: recurrant
    literals:
      # Repo with manifests to apply
      - GIT_SYNC_REPO=https://github.com/vrutkovs/k8s-podhunt
      # Repo branch
      - GIT_SYNC_REF=master
      # Sync period
      - GIT_SYNC_WAIT=10
      # Backoff time before re-attempting the apply
      - GIT_SYNC_WEBHOOK_TIMEOUT=30s
      # Subdir with manifests to apply
      - RECURRANT_SUBDIR=manifests
      # Optional: use kustomize to apply manifests
      - RECURRANT_USE_KUSTOMIZE=true
      # Optional: use oc instead of kubectl
      - RECURRANT_USE_OC=true
```

# Why not an operator?

Creating an operator to setup re-currant in particular namespaces and granular permissions would be very useful. Although at this stage creating CRDs requires `cluster-admin` permissions, which might not work for some installs. This app is deliberately kept simple and requires minimal permissions. Other permissions (like cross-namespace applies) could be added via additional rolebindings to re-currant serviceaccount.

# Features

## Push model

The pod exposes `/reload` endpoint, which restarts `git-sync` sidecar and makes it re-pull the tracked branch and apply changes.

Note, that re-currant is meant to be kept simple, it basically runs `kubectl apply`. Some deployments may require a pipeline setup, in this case [Tekton](https://tekton.dev/) would be a better choice.

## Custom command

Set `RECURRANT_COMMAND` to use a custom command to start the deploy (make sure it uses `RECURRANT_SUBDIR` env var)

## Git repository credentials

Copy ssh key to `id_rsa` file and use `ssh-keyscan github.com > known_hosts` to get setup a known hosts file.

Kustomize file which reconfigures git-sync to use SSH key would look similar to:

```
# Use remote base
bases:
  - github.com/vrutkovs/re-currant//manifests?ref=master
patchesStrategicMerge:
  - |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: recurrant
      namespace: recurrant
    spec:
      template:
        spec:
          containers:
            - name: git-sync
              env:
                - name: GIT_SYNC_SSH
                  value: "true"
              volumeMounts:
                - name: git-secret
                  mountPath: /etc/git-secret
          volumes:
            - name: git-secret
              secret:
                secretName: git-creds
                defaultMode: 444
generatorOptions:
  disableNameSuffixHash: true
secretGenerator:
  - name: git-creds
    namespace: recurrant
    files:
      - ssh=id_rsa
      - known-hosts=known_hosts
  - name: recurrant
    namespace: recurrant
    literals:
      - GIT_SYNC_REPO=git@github.com:vrutkovs/ocp-gitops
      - GIT_SYNC_REF=master
      - GIT_SYNC_WAIT=10
      - GIT_SYNC_WEBHOOK_TIMEOUT=30s
      - RECURRANT_SUBDIR=namespaces/gitops
      - RECURRANT_USE_KUSTOMIZE=false
      - RECURRANT_USE_OC=true
```
