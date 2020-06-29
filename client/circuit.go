package client

import (
	"net/url"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/pkg/errors"
)

//CircuitBreakerConfig hold circuit breaker config for a host, service connects to
type CircuitBreakerConfig struct {
	// BaseURL of the host for which circuit breaker needs to be configured
	BaseURL string
	//Circuit breaker config for corresponding service
	CircuitConfig hystrix.CommandConfig
}

// Register - register circuit breaker for all hosts, service connects to.
func Register(cb []CircuitBreakerConfig) error {

	for _, cbCfg := range cb {
		if cbCfg.BaseURL != "" {

			u, err := url.Parse(cbCfg.BaseURL)
			if err != nil {
				return errors.Wrapf(err, "Failed to parse url: %s", cbCfg.BaseURL)
			}

			//Register cb using hostname
			hostname := u.Host
			hystrix.ConfigureCommand(hostname, cbCfg.CircuitConfig)
		} else {
			return errors.New("Missing BaseURL in circuit breaker config")
		}
	}
	return nil
}
