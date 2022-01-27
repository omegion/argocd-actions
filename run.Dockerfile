ARG GO_VERSION=1.17-alpine3.14
ARG FROM_IMAGE=alpine:3.14
ARG IMAGE_TAG

FROM ghcr.io/omegion/argocd-actions:${IMAGE_TAG}
