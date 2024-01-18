package config

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

const defaultLimit = 10

type ServiceEndpoints struct {
	Base      string              `json:"base"`
	Endpoints map[string]Endpoint `json:"endpoints"`
}

func (se *ServiceEndpoints) UnmarshalJSON(rawServiceData []byte) error {
	if len(rawServiceData) == 0 {
		logrus.Error("error service endpoints data is missing")
		return errServiceEndpointDataMissing
	}

	var serviceEndpoints struct {
		Base      string              `json:"base"`
		Endpoints map[string]Endpoint `json:"endpoints"`
	}

	err := json.Unmarshal(rawServiceData, &serviceEndpoints)
	if err != nil {
		logrus.Errorf("error unmarshaling service endpoints data: %v", err)
		return errUnmarshallingServiceData
	}

	se.Base = serviceEndpoints.Base
	se.Endpoints = serviceEndpoints.Endpoints

	for key, endpointData := range se.Endpoints {
		endpointData.SetBaseURL(se.Base)
		se.Endpoints[key] = endpointData
	}

	return nil
}

// GetEndpoint returns the service endpoint based on the given alias
func (se *ServiceEndpoints) GetEndpoint(endpointAlias string) (Endpoint, error) {
	endpointData, endpointExists := se.Endpoints[endpointAlias]
	if !endpointExists {
		return Endpoint{}, fmt.Errorf("%w: %s", errEndpointDoesNotExist, endpointAlias)
	}

	return endpointData, nil
}

type Endpoint struct {
	Path        string       `json:"path"`
	Method      string       `json:"method"`
	QueryParams *QueryParams `json:"query_params"`
	baseURL     string
}

func (e *Endpoint) SetBaseURL(base string) {
	e.baseURL = base
}

func (e *Endpoint) GetURL() string {
	return e.baseURL + e.Path
}

type QueryParams struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (qp *QueryParams) UnmarshalJSON(rawData []byte) error {
	if len(rawData) == 0 {
		return nil
	}

	var queryParams struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}
	err := json.Unmarshal(rawData, &queryParams)
	if err != nil {
		return fmt.Errorf("%w: %w", errUnmarshallingQueryParams, err)
	}

	qp.Limit = queryParams.Limit
	if qp.Limit == 0 {
		qp.Limit = defaultLimit
	}
	qp.Offset = queryParams.Offset

	return nil
}

// ToMap converts QueryParams struct into a map
func (qp *QueryParams) ToMap() map[string]string {
	paramsMap := make(map[string]string)

	paramsMap["limit"] = fmt.Sprintf("%v", qp.Limit)
	paramsMap["offset"] = fmt.Sprintf("%v", qp.Offset)

	return paramsMap
}
