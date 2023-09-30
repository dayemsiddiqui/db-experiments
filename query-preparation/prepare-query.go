package query_preparation

import (
	"db-experiments/config"
	"db-experiments/query-preparation/rules"
)

type ParamStrategy interface {
	Apply(query string, param config.InputParameters) string
}

func PrepareQuery(cfg *config.Config, queryConfig *config.QueryConfig) string {
	query := queryConfig.Query

	strategies := []ParamStrategy{
		&rules.INParamStrategy{},
		&rules.ArrayParamStrategy{},
		&rules.DefaultParamStrategy{},
	}

	for _, param := range cfg.Parameters {
		for _, strategy := range strategies {
			query = strategy.Apply(query, param)
		}
	}
	return query
}
