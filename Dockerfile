ARG GO_VERSION=1.18-alpine3.15
ARG FROM_IMAGE=alpine:3.15

FROM golang:${GO_VERSION} AS builder

WORKDIR /app

COPY ./ /app

RUN apk update && \
  apk add ca-certificates gettext git make curl unzip && \
  rm -rf /tmp/* && \
  rm -rf /var/cache/apk/* && \
  rm -rf /var/tmp/*

RUN make build-for-container

FROM ${FROM_IMAGE}

COPY --from=builder /app/dist/argocd-actions-linux /bin/argocd-actions

ENTRYPOINT ["argocd-actions"]
