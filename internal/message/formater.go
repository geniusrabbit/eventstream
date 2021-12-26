//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package message

// Formater processor
type Formater interface {
	Format(msg Message) (interface{}, error)
}
