package contract

import "sync"

var (
	ccenter = make(map[string]Contract)
	mutex   sync.Mutex
)

type Contract interface {
	Name() string
	Def() string
	Data(string, []interface{}) (string, error)
	Result(string, string) (interface{}, error)
}

func Register(c Contract) {
	mutex.Lock()
	defer mutex.Unlock()
	ccenter[c.Name()] = c
}

func Get(addr string) Contract {
	if c, ok := ccenter[addr]; ok {
		return c
	}
	return nil
}
