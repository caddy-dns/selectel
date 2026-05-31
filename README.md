# Selectel DNS v2 module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with [Selectel DNS v2 API](https://developers.selectel.ru/docs/cloud-services/dns_api/dns_api_actual/).

## Caddy module name

```
dns.providers.selectel
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "selectel",
        "user": "SELECTEL_USER",
        "password": "SELECTEL_PASSWORD",
        "account_id": "SELECTEL_ACCOUNT_ID",
        "project_name": "SELECTEL_PROJECT_NAME"
      }
    }
  }
}
```

or with the Caddyfile:

```
# globally
{
	acme_dns selectel {
		user {env.SELECTEL_USER}
		password {env.SELECTEL_PASSWORD}
		account_id {env.SELECTEL_ACCOUNT_ID}
		project_name {env.SELECTEL_PROJECT_NAME}
	}
}
```

```
# one site
tls {
	dns selectel {
		user {env.SELECTEL_USER}
		password {env.SELECTEL_PASSWORD}
		account_id {env.SELECTEL_ACCOUNT_ID}
		project_name {env.SELECTEL_PROJECT_NAME}
	}
}
```

### Optional directives

| Directive | Description |
|---|---|
| `enable_debug_logging` | Enables verbose DEBUG-level log output for all DNS operations. INFO and ERROR messages are always emitted when a logger is set. |

Selectel [Service user](https://my.selectel.ru/iam/users_management/users?type=service) management
