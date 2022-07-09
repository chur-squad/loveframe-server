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

// WithCdnEndpointForDash returns dependency injection method for dash cdn endpoint
func WithCdnEndpointForDash(endpoint string) OptionFunc {
	return func(m *photoMaker) {
		m.cdnEndpoint = endpoint
	}
}