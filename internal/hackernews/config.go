package hackernews

import "net/http"

type ConfigFunc func(c *config)

type config struct {
	roundTripper http.RoundTripper
}

func WithRoundTripper(roundTripper http.RoundTripper) ConfigFunc {
	return func(c *config) {
		c.roundTripper = roundTripper
	}
}
