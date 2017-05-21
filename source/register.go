//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package source

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/kafka"
	"github.com/geniusrabbit/notificationcenter/nats"
)

// Register stream subscriber
func Register(name, connection string) (err error) {
	var sub notificationcenter.Subscriber
	if sub, err = newSubscriber(connection); nil == err {
		err = notificationcenter.Register(name, sub)
	}
	return
}

// Subscribe handler
func Subscribe(name string, handle notificationcenter.Handler) error {
	return notificationcenter.Subscribe(name, handle)
}

func newSubscriber(connection string) (notificationcenter.Subscriber, error) {
	var url, err = url.Parse(connection)
	if nil != err {
		return nil, err
	}
	var params = url.Query()

	switch url.Scheme {
	case "kafka":
		return kafka.NewSubscriber(
			strings.Split(url.Host, ","),
			url.Path[1:],
			strings.Split(params.Get("topics"), ","),
		)
	case "nats":
		return nats.NewSubscriber(
			strings.Split(params.Get("topics"), ","),
			"nats://"+url.Host,
		)
	}
	return nil, fmt.Errorf("Undefined log scheme: %s", url.Scheme)
}
