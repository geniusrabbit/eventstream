//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package converter

import (
	"encoding/json"

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
)

// ByName decoder object
func ByName(name string) Converter {
	switch name {
	case "bson":
		return BSON
	}
	return JSON
}
