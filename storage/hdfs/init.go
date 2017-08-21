//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package hdfs

import (
	"net/url"
	"time"

	"github.com/vladimirvivien/gowfs"
)

// func init() {
// 	storage.RegisterConnector(hdfsConnect, "hdfs")
// }

func hdfsConnect(u *url.URL, debug bool) (interface{}, error) {
	// return hdfs.New(u.Host)

	conf := *gowfs.NewConfiguration()
	conf.Addr = u.Host
	conf.User = "hdfs"
	conf.ConnectionTimeout = time.Second * 15
	conf.DisableKeepAlives = false

	return gowfs.NewFileSystem(conf)
}
