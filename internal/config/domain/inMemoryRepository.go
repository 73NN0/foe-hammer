package domain

import "fmt"

type MemoryConfig struct {
	data map[string]string
}

func NewMemoryConfig(datas ...[]string) *MemoryConfig {

	d := make(map[string]string)
	for _, data := range datas {
		if len(data) != 2 {
			fmt.Printf("data need to be in pair (k,v) got : %q ... skipping \n", data)
			continue
		}
		key := data[0]
		value := data[1]

		d[key] = value
	}

	return &MemoryConfig{
		data: d,
	}
}
