//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package context

import "github.com/demdxx/gocast"

type options map[string]interface{}

func (o options) Item(key string, def interface{}) interface{} {
	if nil != o {
		if v, ok := o[key]; ok && nil != o {
			return v
		}
	}
	return def
}

func (o options) Int(key string, def int64) int64 {
	return gocast.ToInt64(o.Item(key, def))
}

func (o options) Float(key string, def float64) float64 {
	return gocast.ToFloat64(o.Item(key, def))
}

func (o options) Bool(key string, def bool) bool {
	return gocast.ToBool(o.Item(key, def))
}

func (o options) String(key, def string) string {
	return gocast.ToString(o.Item(key, def))
}
