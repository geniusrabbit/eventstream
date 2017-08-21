//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

// Metric record type
const (
	TypeIncrement = iota
	TypeCounter
	TypeTiming
	TypeUnique
)

// Item metric message
type Item struct {
	Name string   `json:"name"`
	Type int      `json:"type"`
	Tags []string `json:"tags"`
}
