package env

import (
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
)

type env struct {
	prefixs []string
}

func NewSource(prefixs ...string) config.Source {
	return &env{prefixs: prefixs}
}

func (e *env) Load() (kv []*config.KeyValue, err error) {
	return e.load(os.Environ()), nil
}

func (e *env) load(envStrings []string) []*config.KeyValue {
	var kv []*config.KeyValue
	for _, envstr := range envStrings {
		var k, v string
		subs := strings.SplitN(envstr, "=", 2)
		k = subs[0]
		if len(subs) > 1 {
			v = subs[1]
		}

		if len(e.prefixs) > 0 {
			p, ok := matchPrefix(e.prefixs, envstr)
			if !ok || len(p) == len(k) {
				continue
			}
			// trim prefix
			k = k[len(p):]
			if k[0] == '_' {
				k = k[1:]
			}
		}

		if len(k) != 0 {
			kv = append(kv, &config.KeyValue{
				Key:   k,
				Value: []byte(v),
			})
		}
	}
	return kv
}

func (e *env) Watch() (config.Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func matchPrefix(prefixs []string, s string) (string, bool) {
	for _, p := range prefixs {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}
	return "", false
}
