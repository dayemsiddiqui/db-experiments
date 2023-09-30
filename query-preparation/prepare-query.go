package query_preparation

import (
	"db-experiments/config"
	"fmt"
	"strings"
)

type ParamStrategy interface {
	Apply(query string, param config.InputParameters) string
}

type INParamStrategy struct{}
type ArrayParamStrategy struct{}
type DefaultParamStrategy struct{}

func (s *INParamStrategy) Apply(query string, param config.InputParameters) string {
	return strings.ReplaceAll(query, "IN "+param.Name, renderINParam(param.Value))
}

func (s *ArrayParamStrategy) Apply(query string, param config.InputParameters) string {
	return strings.ReplaceAll(query, "ARRAY["+param.Name+"]", renderArrayParam(param.Value))
}

func (s *DefaultParamStrategy) Apply(query string, param config.InputParameters) string {
	return strings.ReplaceAll(query, param.Name, paramToString(param.Value))
}

func PrepareQuery(cfg *config.Config, queryConfig *config.QueryConfig) string {
	query := queryConfig.Query

	strategies := []ParamStrategy{
		&INParamStrategy{},
		&ArrayParamStrategy{},
		&DefaultParamStrategy{},
	}

	for _, param := range cfg.Parameters {
		for _, strategy := range strategies {
			query = strategy.Apply(query, param)
		}
	}
	return query
}

func paramToString(arr []string) string {
	return fmt.Sprintf("%q", arr)
}

func renderINParam(arr []string) string {
	return "IN (" + strings.Join(arr, ",") + ")"
}

func renderArrayParam(arr []string) string {
	return "ARRAY[" + strings.Join(arr, ",") + "]"
}
