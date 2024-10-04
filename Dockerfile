FROM mcr.microsoft.com/devcontainers/go:1.3.0-1.23-bookworm

ENV APP_ENV=development

COPY .cobra.yaml ~/.cobra.yaml

RUN \
    go install github.com/spf13/cobra-cli@latest && \
    go install github.com/goreleaser/goreleaser/v2@latest && \
    chown -R vscode /go