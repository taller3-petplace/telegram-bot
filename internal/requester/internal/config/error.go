package config

import "errors"

var (
	errEndpointDoesNotExist       = errors.New("error endpoint does not exist")
	errUnmarshallingQueryParams   = errors.New("error unmarshalling query params")
	errServiceEndpointDataMissing = errors.New("error service endpoints data is missing")
	errUnmarshallingServiceData   = errors.New("error unmarshalling service data")
)
