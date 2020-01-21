package parseutils

import "fmt"

type KVpair struct { Key, Value string }

func (kvp KVpair) String() string {
	return fmt.Sprintf("%s<%s> ", kvp.Key, kvp.Value)
}
