# Selectel DNS v2 module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It manages DNS records via the [Selectel DNS v2 API](https://docs.selectel.ru/en/api/dns-actual/).

## Caddy module name

```
dns.providers.selectel
```

## Config examples

### Caddyfile

```
# Globally (all sites)
{
	acme_dns selectel {
		user        {env.SELECTEL_USER}
		password    {env.SELECTEL_PASSWORD}
		account_id  {env.SELECTEL_ACCOUNT_ID}
		project_name {env.SELECTEL_PROJECT_NAME}
	}
}
```

```
# Per site
tls {
	dns selectel {
		user        {env.SELECTEL_USER}
		password    {env.SELECTEL_PASSWORD}
		account_id  {env.SELECTEL_ACCOUNT_ID}
		project_name {env.SELECTEL_PROJECT_NAME}
	}
}
```

### JSON

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

## Required directives

| Directive | Description |
|---|---|
| `user` | Selectel service-user login. |
| `password` | Selectel service-user password. |
| `account_id` | Selectel account ID (domain name). |
| `project_name` | Selectel project name that owns the DNS zones. |

All values support the `{env.VAR}` placeholder syntax.

## Optional directives

| Directive | Description |
|---|---|
| `enable_debug_logging` | Enables verbose DEBUG-level log output for all DNS operations (bare flag = `true`). Accepts an optional value: `true`, `false`, `1`, `0`, or an env placeholder that resolves to one of these (e.g. `{env.SELECTEL_DEBUG}`). |

### Debug logging

When `enable_debug_logging` is set, all provider messages (INFO, ERROR, and DEBUG) flow through Caddy's configured logging system at the `debug` level. Caddy's own log level setting controls whether they are emitted. The severity of each message is preserved as a `[INFO]` / `[ERROR]` / `[DEBUG]` prefix inside the message text.

```
{
    log {
        level DEBUG
    }
    acme_dns selectel {
        user        {env.SELECTEL_USER}
        password    {env.SELECTEL_PASSWORD}
        account_id  {env.SELECTEL_ACCOUNT_ID}
        project_name {env.SELECTEL_PROJECT_NAME}
        enable_debug_logging
    }
}
```

## Credentials

This module requires a Selectel [service user](https://my.selectel.ru/iam/users_management/users?type=service) with access to the project that owns the DNS zones.

The Selectel DNS v2 API uses a project-scoped IAM token (Keystone v3). The module obtains and automatically refreshes this token — you only need to supply the service-user credentials above.
