package parseutils

import "fmt"

// KVpair is used for storign property values. 
type KVpair struct { Key, Value string }

func (kvp KVpair) String() string {
	return fmt.Sprintf("%s<%s> ", kvp.Key, kvp.Value)
}
