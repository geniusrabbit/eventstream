//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package kafka

import (
	"log"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
)

type subs struct {
	debug  bool
	format converter.Converter
	stream eventstream.Streamer
}

func (s *subs) Handle(item interface{}) error {
	msg, err := eventstream.MessageDecode(item, s.format)
	if err == nil && s.stream.Check(msg) {
		err = s.stream.Put(msg)
	} else if s.debug {
		log.Println("[nats] decode message", err, item)
	}
	return err
}
