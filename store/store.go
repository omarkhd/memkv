package store

type Store interface {
	Keys() []string
	Put(k, v string)
	Get(k string) string
	Delete(k string)
}
