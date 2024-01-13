package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

const defaultLimit = 10

type ServiceEndpoints struct {
	Base      string              `json:"base"`
	Endpoints map[string]Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Path        string       `json:"path"`
	Method      string       `json:"method"`
	QueryParams *QueryParams `json:"query_params"`
}

type QueryParams struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (qp *QueryParams) UnmarshalYAML(value *yaml.Node) error {
	//TODO implement me
	panic("implement me")
}

/*func (qp *QueryParams) UnmarshalJSON(rawData []byte) error {
	if len(rawData) == 0 {
		return nil
	}

	qp.Limit = 100
	qp.Offset = 0

	return nil
}*/

func (qp *QueryParams) ToMap() map[string]string {
	paramsMap := make(map[string]string)

	limit := qp.Limit
	// ToDo: add this logic in UnmarshalJSON
	if limit == 0 {
		limit = defaultLimit
	}

	paramsMap["limit"] = fmt.Sprintf("%v", limit)
	paramsMap["offset"] = fmt.Sprintf("%v", qp.Offset)

	return paramsMap
}
