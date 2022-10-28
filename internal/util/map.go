package util

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// MergeStringMap merge two map
// Deprecated: this function will be moved to internal package, user should not use it any more.
func MergeStringMap(dest, src map[string]interface{}) {
	for sk, sv := range src {
		tv, ok := dest[sk]
		if !ok {
			// val不存在时，直接赋值
			dest[sk] = sv
			continue
		}

		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)
		if svType != tvType {
			fmt.Println("continue, type is different")
			continue
		}

		switch ttv := tv.(type) {
		case map[interface{}]interface{}:
			tsv := sv.(map[interface{}]interface{})
			ssv := ToMapStringInterface(tsv)
			stv := ToMapStringInterface(ttv)
			MergeStringMap(stv, ssv)
			dest[sk] = stv
		case map[string]interface{}:
			MergeStringMap(ttv, sv.(map[string]interface{}))
			dest[sk] = ttv
		default:
			dest[sk] = sv
		}
	}
}

// ToMapStringInterface cast map[interface{}]interface{} to map[string]interface{}
// Deprecated: this function will be moved to internal package, user should not use it any more.
func ToMapStringInterface(src map[interface{}]interface{}) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

// DeepSearchInMap deep search in map
// Deprecated: this function will be moved to internal package, user should not use it any more.
func DeepSearchInMap(m map[string]interface{}, paths ...string) map[string]interface{} {
	mtmp := make(map[string]interface{})
	for k, v := range m {
		mtmp[k] = v
	}
	for _, k := range paths {
		m2, ok := mtmp[k]
		if !ok {
			m3 := make(map[string]interface{})
			mtmp[k] = m3
			mtmp = m3
			continue
		}

		m3, err := cast.ToStringMapE(m2)
		if err != nil {
			m3 = make(map[string]interface{})
			mtmp[k] = m3
		}
		// continue search
		mtmp = m3
	}
	return mtmp
}
