ARG VERSION=1.23.2-bookworm

FROM golang:${VERSION}

ENV APP_ENV=development

COPY .cobra.yaml .cobra.yaml

RUN \
  go install github.com/spf13/cobra-cli@latest && \
  go install github.com/goreleaser/goreleaser/v2@latest && \
  mv .cobra.yaml ~/.cobra.yaml
