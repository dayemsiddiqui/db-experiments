package rules

import (
	"db-experiments/config"
	"fmt"
	"strings"
)

type DefaultParamStrategy struct{}

func (s *DefaultParamStrategy) Apply(query string, param config.InputParameters) string {
	value := fmt.Sprintf("%q", param.Value)
	return strings.ReplaceAll(query, param.Name, value)
}
