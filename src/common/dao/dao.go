package dao

import "github.com/kllla/web/src/config"

// Dao is the default actions required for all Daos
// independent of their data type
type Dao interface {
	SetClientFromConfig(config config.Config)
	Close() error
}
