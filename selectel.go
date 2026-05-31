package selectel

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/selectel"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *selectel.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.selectel",
		New: func() caddy.Module { return &Provider{new(selectel.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.Provider.User = caddy.NewReplacer().ReplaceAll(p.Provider.User, "")
	p.Provider.Password = caddy.NewReplacer().ReplaceAll(p.Provider.Password, "")
	p.Provider.AccountId = caddy.NewReplacer().ReplaceAll(p.Provider.AccountId, "")
	p.Provider.ProjectName = caddy.NewReplacer().ReplaceAll(p.Provider.ProjectName, "")

	if p.Provider.ZonesCache == nil {
		p.Provider.ZonesCache = make(map[string]string)
	}

	return nil
}

func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "user":
				if p.Provider.User != "" {
					return d.Err("User already set")
				}
				if d.NextArg() {
					p.Provider.User = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "password":
				if p.Provider.Password != "" {
					return d.Err("Password already set")
				}
				if d.NextArg() {
					p.Provider.Password = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "account_id":
				if p.Provider.AccountId != "" {
					return d.Err("AccountId already set")
				}
				if d.NextArg() {
					p.Provider.AccountId = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "project_name":
				if p.Provider.ProjectName != "" {
					return d.Err("ProjectName already set")
				}
				if d.NextArg() {
					p.Provider.ProjectName = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "enable_debug_logging":
				p.Provider.EnableDebugLogging = true
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}

	if p.Provider.User == "" {
		return d.Err("missing User")
	}
	if p.Provider.Password == "" {
		return d.Err("missing Password")
	}
	if p.Provider.AccountId == "" {
		return d.Err("missing AccountId")
	}
	if p.Provider.ProjectName == "" {
		return d.Err("missing ProjectName")
	}

	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
