package hustle

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Config is the bag of poo that everyone knows about.  Wheeeee!
type Config struct {
	HTTPAddr    string
	HTTPPubAddr string
	// HTTPSAddr string
	HubAddr   string
	WSAddr    string
	WSPubAddr string
	// WSSAddr   string
	StatsAddr    string
	StatsPubAddr string
}

// ProcessConfig wraps envconfig.Process, dangit.
func ProcessConfig(config *Config) error {
	return envconfig.Process("hustle", config)
}

// HTTPHost returns only the host from HTTPAddr
func (cfg *Config) HTTPHost() string {
	return hostPart(cfg.HTTPAddr)
}

// HTTPPort returns only the port from HTTPAddr
func (cfg *Config) HTTPPort() string {
	return portPart(cfg.HTTPAddr)
}

// HTTPPubHost returns only the host from HTTPPubAddr
func (cfg *Config) HTTPPubHost() string {
	return hostPart(cfg.HTTPPubAddr)
}

// HTTPPubPort returns only the port from HTTPPubAddr
func (cfg *Config) HTTPPubPort() string {
	return portPart(cfg.HTTPPubAddr)
}

// WSHost returns only the host from WSAddr
func (cfg *Config) WSHost() string {
	return hostPart(cfg.WSAddr)
}

// WSPort returns only the port from WSAddr
func (cfg *Config) WSPort() string {
	return portPart(cfg.WSAddr)
}

// WSPubHost returns only the host from WSPubAddr
func (cfg *Config) WSPubHost() string {
	return hostPart(cfg.WSPubAddr)
}

// WSPubPort returns only the port from WSPubAddr
func (cfg *Config) WSPubPort() string {
	return portPart(cfg.WSPubAddr)
}

func hostPart(addr string) string {
	host := strings.SplitN(addr, ":", 2)[0]
	if host == "" {
		host = "localhost"
	}
	return host
}

func portPart(addr string) string {
	return strings.SplitN(addr, ":", 2)[1]
}
