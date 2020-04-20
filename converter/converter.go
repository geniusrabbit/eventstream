//
// @project geniusrabbit::eventstream 2017, 2019-2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2020
//

package converter

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"github.com/pkg/errors"
)

var (
	errUnsupportedDecodeType = errors.New(`unsupported converter decode type`)
	errUnsupportedEncodeType = errors.New(`unsupported converter encode type`)
)

// Converter interaface
type Converter interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// Fnk decoder wrapper
type Fnk struct {
	name    string
	encoder func(v interface{}) ([]byte, error)
	decoder func(data []byte, v interface{}) error
}

func (f Fnk) String() string {
	return f.name
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
	JSON Converter = Fnk{name: "json", encoder: json.Marshal, decoder: json.Unmarshal}
	BSON Converter = Fnk{name: "bson", encoder: bson.Marshal, decoder: bson.Unmarshal}
	RAW  Converter = Fnk{
		name: "raw",
		encoder: func(v interface{}) ([]byte, error) {
			switch b := v.(type) {
			case []byte:
				return b, nil
			case string:
				return []byte(b), nil
			}
			return nil, errors.Wrapf(errUnsupportedEncodeType, "[raw] %T", v)
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
			return errors.Wrapf(errUnsupportedDecodeType, "[raw] %T", v)
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
