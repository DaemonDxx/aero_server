package lks

type CookieCache interface {
	Get(key string) (map[string]string, bool)
	Put(key string, cookie map[string]string) error
	Delete(key string)
}
