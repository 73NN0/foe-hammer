package domain

type Config interface {
	Get(key string) (string, bool)
	GetAll() map[string]string
}
