package photos

import (
	"net/http"
)

type Option interface {
	apply(m *photoMaker)
}

type OptionFunc func(maker *photoMaker)

func (o OptionFunc) apply(maker *photoMaker) {
	o(maker)
}

// WithCdnClient returns dependency injection method for cdn client
func WithCdnClient(client *http.Client) OptionFunc {
	return func(m *photoMaker) {
		m.cdnClient = client
	}
}

// WithCdnEndpoint returns dependency injection method for cdn (photo store)
func WithCdnEndpoint(endpoint string) OptionFunc {
	return func(m *photoMaker) {
		m.cdnEndpoint = endpoint
	}
}