//
// @project geniusrabbit::eventstream 2017, 2023
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2023
//

package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/demdxx/gocast/v2"
)

var (
	timeFormats = []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"2006-01-02",
		"01-02-2006",
		time.RFC1123Z,
		time.RFC3339Nano,
		time.UnixDate,
		time.RubyDate,
		time.RFC1123,
		time.RFC3339,
		time.RFC822,
		time.RFC850,
		time.RFC822Z,
	}
)

// ParseTime from string
func parseTime(tm string) (t time.Time, err error) {
	for _, f := range timeFormats {
		if t, err = time.Parse(f, tm); err == nil {
			break
		}
	}
	return t, err
}

func ip2EscapeString(ip net.IP) string {
	var (
		data    = make([]byte, 16)
		ipBytes = ip2Int(ip).Bytes()
	)
	for i := range ipBytes {
		data[15-i] = ipBytes[len(ipBytes)-i-1]
	}
	return escapeBytes(data, 0)
}

func escapeBytes(data []byte, size int) string {
	var buff bytes.Buffer
	for i, b := range data {
		if size > 0 && i > size {
			break
		}
		buff.WriteString(fmt.Sprintf("\\%03o", b))
	}
	for i := len(data); i < size; i++ {
		buff.WriteString(fmt.Sprintf("\\%03o", byte(0)))
	}
	return buff.String()
}

func bytesSize(data []byte, size int) []byte {
	if size < 1 {
		return data
	}
	if len(data) > size {
		return data[:size]
	}
	for i := len(data); i < size; i++ {
		data = append(data, 0)
	}
	return data
}

func ip2Int(ip net.IP) *big.Int {
	ipInt := big.NewInt(0)
	ipInt.SetBytes(ip)
	return ipInt
}

func valueToTime(v any) (tm time.Time) {
	switch vl := v.(type) {
	case nil:
	case int64:
		tm = time.Unix(vl, 0)
	case uint64:
		tm = time.Unix(int64(vl), 0)
	case float64:
		tm = time.Unix(int64(vl), 0)
	case string:
		tm, _ = parseTime(gocast.Str(v))
	default:
		tm, _ = parseTime(gocast.Str(v))
	}
	return tm
}

func valueUnixNanoToTime(v any) (tm time.Time) {
	switch vl := v.(type) {
	case nil:
	case int64:
		tm = time.Unix(0, vl)
	case uint64:
		tm = time.Unix(0, int64(vl))
	case float64:
		tm = time.Unix(0, int64(vl))
	case string:
		tm, _ = parseTime(gocast.Str(v))
	default:
		tm, _ = parseTime(gocast.Str(v))
	}
	return tm
}

func valueToIP(v any) (ip net.IP) {
	switch vl := v.(type) {
	case net.IP:
		ip = vl
	case uint:
		ip = make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, uint32(vl))
	case uint32:
		ip = make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, vl)
	case int:
		ip = make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, uint32(vl))
	case int32:
		ip = make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, uint32(vl))
	default:
		ip = net.ParseIP(gocast.Str(v))
	}
	return ip
}
