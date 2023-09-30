package rules

import (
	"db-experiments/config"
	"strings"
)

type INParamStrategy struct{}

func (s *INParamStrategy) Apply(query string, param config.InputParameters) string {
	value := "IN (" + strings.Join(param.Value, ",") + ")"
	return strings.ReplaceAll(query, "IN "+param.Name, value)
}
