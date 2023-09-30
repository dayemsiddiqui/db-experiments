package execution

import "db-experiments/config"

type RunQueryConfig struct {
	QueryConfig config.QueryConfig
	Config      *config.Config
}
