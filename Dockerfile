FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

ENV APP_ENV=development

COPY .cobra.yaml ~/.cobra.yaml

RUN \
    go install github.com/spf13/cobra-cli@latest