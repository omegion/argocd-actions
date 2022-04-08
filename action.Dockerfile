ARG ARGOCD_ACTIONS_VERSION=latest

FROM ghcr.io/omegion/argocd-actions:${ARGOCD_ACTIONS_VERSION}

ENTRYPOINT ["argocd-actions"]

