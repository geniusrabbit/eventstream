//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package kafka

import (
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
)

type subs struct {
	format converter.Converter
	stream eventstream.Streamer
}

func (s *subs) Handle(item interface{}) error {
	msg, err := eventstream.MessageDecode(item, s.format)
	if err == nil && s.stream.Check(msg) {
		err = s.stream.Put(msg)
	}
	return err
}
