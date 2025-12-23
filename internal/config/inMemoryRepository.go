package config

import (
	"fmt"
)

// Implémentation simple
type InMemoryConfig struct {
	data map[string]string
}

func NewInMemoryConfig() *InMemoryConfig {
	return &InMemoryConfig{
		data: map[string]string{
			// todo improve this
			"rootDir":  "sys",
			"buildDir": "bin",
			"cc":       "clang",
			"cflagsCommon": `-Wall -Wextra -Wpedantic \
		-Wshadow -Wcast-align -Wunused -Wold-style-definition \
		-Wmissing-prototypes -Wno-unused-parameter -Werror \
		-Wstrict-prototypes -Wpointer-arith -Wwrite-strings \
		-Wconversion -Wformat=2 -Wformat-security \
		-Wunreachable-code -Wundef -Wbad-function-cast \
		-Wdouble-promotion -Wmissing-include-dirs \
		-Winit-self -Wmissing-noreturn -fno-common \
		-fstack-protector-strong`,

			"cflagsRelease": `-O2 -DNDEBUG -DDEBUG_MEMORY=0 -fomit-frame-pointer -march=native -D_FORTIFY_SOURCE=2`,

			"cflagsDebug": `-g3 -O0 -DDEBUG -DDEBUG_MEMORY=1 -ftrapv`,
			"linkFlags":   "-lSDL2 -lm",
		},
	}
}

func (c *InMemoryConfig) Get(key string) (string, error) {
	if val, ok := c.data[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found: %s", key)
}

func (c *InMemoryConfig) GetAll() map[string]string {
	return c.data
}
