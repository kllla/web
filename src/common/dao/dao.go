package dao

import "github.com/kllla/web/src/config"

type Dao interface {
	SetClientFromConfig(config config.Config)
	Close() error
}
