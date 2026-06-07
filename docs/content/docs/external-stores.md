# External Stores

::: warning Coming in 0.4
External store support is planned for the 0.4 release. This page documents the intended design.
:::

External stores allow kredenv to fetch secrets from remote secrets managers instead of (or in addition to) the local vault.

## Planned Providers

- **AWS SSM Parameter Store** — fetch secrets from AWS Systems Manager
- **HashiCorp Vault** — fetch secrets from a self-hosted Vault instance
- **Azure Key Vault** — fetch secrets from Azure Key Vault
- **GCP Secret Manager** — fetch secrets from Google Cloud Secret Manager

## Intended Design

External stores will be declared in `kredsfile.yaml` under a `stores` block:

```yaml
stores:
  - name: production
    type: aws_ssm
    region: us-east-1
    path_prefix: /myapp/production

secrets:
  - key: DATABASE_PASSWORD
    store: production

  - key: API_KEY
    # no store = local vault (default behaviour)
```

The `store` field on a secret points to a named entry in `stores`. Secrets without a `store` field continue to use the local vault as today.

Authentication credentials for external stores will be configured globally via `kredenv configure` — separate from the kredsfile, not committed to version control.

## In the Meantime

For CI and production environments, use your platform's native secret injection:

- **GitHub Actions** — [Encrypted secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- **GitLab CI** — [CI/CD variables](https://docs.gitlab.com/ee/ci/variables/)
- **AWS** — IAM roles and instance profiles
- **Kubernetes** — Secrets or external-secrets operator

kredenv is designed for developer machines. For production secret management, use a dedicated remote secrets manager.
