//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package converter

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// Converter interaface
type Converter interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// Fnk decoder wrapper
type Fnk struct {
	encoder func(v interface{}) ([]byte, error)
	decoder func(data []byte, v interface{}) error
}

// Marshal value to data
func (f Fnk) Marshal(v interface{}) ([]byte, error) {
	return f.encoder(v)
}

// Unmarshal data to value
func (f Fnk) Unmarshal(data []byte, v interface{}) error {
	return f.decoder(data, v)
}

// Converters
var (
	JSON Converter = Fnk{encoder: json.Marshal, decoder: json.Unmarshal}
	BSON Converter = Fnk{encoder: bson.Marshal, decoder: bson.Unmarshal}
	RAW  Converter = Fnk{
		encoder: func(v interface{}) ([]byte, error) {
			switch b := v.(type) {
			case []byte:
				return b, nil
			case string:
				return []byte(b), nil
			}
			return nil, fmt.Errorf("[raw] unsupported converter encode type %T", v)
		},
		decoder: func(data []byte, v interface{}) error {
			switch pt := v.(type) {
			case *interface{}:
				*pt = data
			case *[]byte:
				*pt = data
			case []byte:
				copy(pt, data)
			}
			return fmt.Errorf("[raw] unsupported converter decode type %T", v)
		},
	}
)

// ByName decoder object
func ByName(name string) Converter {
	switch name {
	case "bson":
		return BSON
	case "json":
		return JSON
	}
	return RAW
}
