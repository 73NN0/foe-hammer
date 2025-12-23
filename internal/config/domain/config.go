package config

import (
	"fmt"
)

type ConfigRepository interface {
	Get(key string) (string, error)
	GetAll() map[string]string
}
