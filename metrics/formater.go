//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

// Formater of metric message
type Formater interface {
	Format(it *Item) (interface{}, error)
}
