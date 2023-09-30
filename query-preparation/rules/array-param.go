package rules

import (
	"db-experiments/config"
	"strings"
)

type ArrayParamStrategy struct{}

func (s *ArrayParamStrategy) Apply(query string, param config.InputParameters) string {
	value := "ARRAY[" + strings.Join(param.Value, ",") + "]"
	return strings.ReplaceAll(query, "ARRAY["+param.Name+"]", value)
}
