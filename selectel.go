package selectel

import (
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/selectel"
	"go.uber.org/zap"
)

// Provider lets Caddy read and manipulate DNS records hosted by Selectel.
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
	repl := caddy.NewReplacer()

	p.Provider.User = repl.ReplaceAll(p.Provider.User, "")
	p.Provider.Password = repl.ReplaceAll(p.Provider.Password, "")
	p.Provider.AccountId = repl.ReplaceAll(p.Provider.AccountId, "")
	p.Provider.ProjectName = repl.ReplaceAll(p.Provider.ProjectName, "")

	// Resolve enable_debug_logging if it was set via an env placeholder.
	// The field may already be true (bare directive) or false (default).
	// The Caddyfile parser sets p.Provider.EnableDebugLogging = true for the
	// bare form; the env-placeholder form sets the string which is then
	// evaluated during Provision via the replacer below (see UnmarshalCaddyfile).

	// Inject the Caddy zap logger so provider logs flow through Caddy's
	// configured log routing (file sinks, JSON format, level filters, etc.).
	p.Provider.Logger = &zapAdapter{logger: ctx.Logger()}

	return nil
}

// zapAdapter bridges the selectel.Logger interface to Caddy's zap logger.
// All provider messages (INFO, ERROR, DEBUG) are routed at the Debug level
// in Caddy's logging system. The [INFO] / [ERROR] / [DEBUG] prefix embedded
// in each message by the provider preserves the original severity for human
// readers. Caddy's own log level filter controls whether these messages are
// emitted at all.
type zapAdapter struct {
	logger *zap.Logger
}

func (a *zapAdapter) Printf(format string, v ...any) {
	a.logger.Sugar().Debugf(format, v...)
}

// UnmarshalCaddyfile parses the provider configuration from a Caddyfile block.
//
// Supported directives:
//
//	user              <value>
//	password          <value>
//	account_id        <value>
//	project_name      <value>
//	enable_debug_logging [true|false|1|0]
//
// All values support the {env.VAR} placeholder syntax.
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "user":
				if p.Provider.User != "" {
					return d.Err("user already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.User = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}

			case "password":
				if p.Provider.Password != "" {
					return d.Err("password already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.Password = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}

			case "account_id":
				if p.Provider.AccountId != "" {
					return d.Err("account_id already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.AccountId = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}

			case "project_name":
				if p.Provider.ProjectName != "" {
					return d.Err("project_name already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.ProjectName = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}

			case "enable_debug_logging":
				if d.NextArg() {
					// Value provided: supports {env.VAR} resolved at Provision time,
					// or a literal true/false/1/0.
					val := d.Val()
					if d.NextArg() {
						return d.ArgErr()
					}
					// Resolve environment placeholders eagerly (before Provision).
					val = caddy.NewReplacer().ReplaceAll(val, "")
					enabled, err := parseBool(val)
					if err != nil {
						return d.Errf("enable_debug_logging: %v", err)
					}
					p.Provider.EnableDebugLogging = enabled
				} else {
					// Bare directive: enable debug logging.
					p.Provider.EnableDebugLogging = true
				}

			default:
				return d.Errf("unrecognized subdirective %q", d.Val())
			}
		}
	}

	if p.Provider.User == "" {
		return d.Err("missing required directive: user")
	}
	if p.Provider.Password == "" {
		return d.Err("missing required directive: password")
	}
	if p.Provider.AccountId == "" {
		return d.Err("missing required directive: account_id")
	}
	if p.Provider.ProjectName == "" {
		return d.Err("missing required directive: project_name")
	}

	return nil
}

// parseBool parses a string value as a boolean for the enable_debug_logging
// directive. Accepted: "true", "1" (enabled), "false", "0" (disabled).
func parseBool(s string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return b, nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
