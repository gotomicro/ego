package egorm

import (
	"sync"
)

var instances = sync.Map{}

// Range 遍历所有实例
func Range(fn func(name string, db *Component) bool) {
	instances.Range(func(key, val interface{}) bool {
		return fn(key.(string), val.(*Component))
	})
}

// Configs
func Configs() map[string]interface{} {
	var rets = make(map[string]interface{})
	instances.Range(func(key, val interface{}) bool {
		return true
	})

	return rets
}

// Stats
func Stats() (stats map[string]interface{}) {
	stats = make(map[string]interface{})
	instances.Range(func(key, val interface{}) bool {
		name := key.(string)
		db := val.(*Component)

		stats[name] = db.DB.DB().Stats()
		return true
	})

	return
}
