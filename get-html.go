package parseutils

import (
	"fmt"
	"golang.org/x/net/html"
)

func (pCPR *ParserResults_html) GetAllByAnyTag(ss []string) []*html.Node {
	if ss == nil || len(ss) == 0 {
		return nil
	}
	var ret = make([]*html.Node, 0)
	for _, p := range pCPR.NodeSlice {
		panic(fmt.Sprintf("OOPS: %p", p))
	}
	return ret
}

// GetAllByTag returns a slice of `*html.Node`. It checks the basic tag only,
// not any namespace. Note that these tag lookup func's default to searching
// the `ListNodesP`, not the tree of `Node`s.
func (pCPR *ParserResults_html) GetAllByTag(s string) []*html.Node {
	if s == "" {
		return nil
	}
	var ret = make([]*html.Node, 0)
	for _, p := range pCPR.NodeSlice {
		panic(fmt.Sprintf("OOPS: %p", p))
	}
	return ret
}
