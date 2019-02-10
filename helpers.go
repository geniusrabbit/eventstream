//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
	"time"
	"unicode"
)

var (
	timeFormats = []string{
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
		if t, err = time.Parse(f, tm); nil == err {
			break
		}
	}
	return
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
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
