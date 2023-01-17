//
// @project geniusrabbit::eventstream 2017, 2019-2023
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2023
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
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// Fnk decoder wrapper
type Fnk struct {
	name    string
	encoder func(v any) ([]byte, error)
	decoder func(data []byte, v any) error
}

func (f Fnk) String() string {
	return f.name
}

// Marshal value to data
func (f Fnk) Marshal(v any) ([]byte, error) {
	return f.encoder(v)
}

// Unmarshal data to value
func (f Fnk) Unmarshal(data []byte, v any) error {
	return f.decoder(data, v)
}

// Converters
var (
	JSON Converter = Fnk{name: "json", encoder: json.Marshal, decoder: json.Unmarshal}
	BSON Converter = Fnk{name: "bson", encoder: bson.Marshal, decoder: bson.Unmarshal}
	RAW  Converter = Fnk{
		name: "raw",
		encoder: func(v any) ([]byte, error) {
			switch b := v.(type) {
			case []byte:
				return b, nil
			case string:
				return []byte(b), nil
			}
			return nil, errors.Wrapf(errUnsupportedEncodeType, "[raw] %T", v)
		},
		decoder: func(data []byte, v any) error {
			switch pt := v.(type) {
			case *any:
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
