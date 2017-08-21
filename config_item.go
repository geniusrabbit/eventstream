//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

import "github.com/demdxx/gocast"

// ConfigItem object
type ConfigItem map[string]interface{}

// Item object from config
func (c ConfigItem) Item(key string, def interface{}) (v interface{}) {
	if v, _ = c[key]; v == nil {
		return def
	}
	return
}

// String item by config
func (c ConfigItem) String(key, def string) string {
	return gocast.ToString(c.Item(key, def))
}

// Int item by config
func (c ConfigItem) Int(key string, def int64) int64 {
	return gocast.ToInt64(c.Item(key, def))
}

// Float item by config
func (c ConfigItem) Float(key string, def float64) float64 {
	return gocast.ToFloat64(c.Item(key, def))
}

// ConfigItem from item
func (c ConfigItem) ConfigItem(key string, def ConfigItem) ConfigItem {
	v := c.Item(key, def)
	if v != nil {
		v, _ = v.(ConfigItem)
	}
	return v.(ConfigItem)
}
