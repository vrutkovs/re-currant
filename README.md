Re-currant
====

A poor man's gitops operator, which blindly applies kustomize files

Create `kustomization.yaml` (see `manifests/kustomization.yaml` as an example):
```
# Create namespace, imagestream, deployment config etc.
bases:
  - manifests/
# Apply patches to change recurrant namespace, config etc.
patchesStrategicMerge:
  - test-patches/route.yaml
generatorOptions:
  disableNameSuffixHash: true
secretGenerator:
  - name: recurrant
    literals:
      # Repo with manifests to apply
      - GIT_SYNC_REPO=https://github.com/vrutkovs/k8s-podhunt
      # Repo branch
      - GIT_SYNC_REF=master
      # Subdir in /tmp/git to checkout the repo
      - GIT_SYNC_CHECKOUT=repo
      # Sync period
      - GIT_SYNC_WAIT=10
      # Backoff time before re-attempting the apply
      - GIT_SYNC_WEBHOOK_TIMEOUT=30
      # Subdir with manifests to apply
      - RECURRANT_SUBDIR=manifests
      # Use kustomize to apply manifests
      - RECURRANT_USE_KUSTOMIZE=true
```

# Git repository credentials

TODO: use kustomize to create a secret, mount in recurrant pod and reconfigure git-sync to use it
