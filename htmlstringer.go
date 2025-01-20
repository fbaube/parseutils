package parseutils

/* IFC
type ParserResults_html struct {
	RootNode  *html.Node
	NodeSlice []*html.Node
	XU.CommonCPR }
html.Node.Render(w io.Writer, n *Node) error
*/

import (
	"fmt"
	// SU "github.com/fbaube/stringutils"
	"golang.org/x/net/html"
	S "strings"
)

// DataOfHtmlNode returns a string that should be
// the value of both [Node.Data] and [Node.DataAtom] .
// If they differ, a warning is issued. Note that if
// the tag is not recognized, DataAtom is left empty.
//
// TODO: Use [strings.Clone] ?
// .
func DataOfHtmlNode(n *html.Node) string {
	datom := n.DataAtom
	datomS := S.TrimSpace(datom.String())
	dataS := S.TrimSpace(n.Data)
	if dataS == datomS {
		return dataS
	}
	if dataS == "" {
		return datomS
	}
	if datomS == "" {
		return dataS
	}
	s := fmt.Sprintf("<<%s>> v <<%s>>", dataS, datomS)
	if datomS == "" {
		println("Unknown HTML tag:", dataS)
	} else {
		println("HtmlNode data mismatch!:", s)
	}
	return s
}

/*
func (p *ParserResults_html) NodeCount() int {
	return len(p.NodeSlice)
}

func (p *ParserResults_html) NodeEcho(i int) string {
	if i >= len(p.NodeSlice) {
		return "(indexOverrun)"
	}
	// var pBB *bytes.Buffer
	var pSB = new(S.Builder)
	// FIXME this call is recursive!
	// html.Render(pSB, p.NodeSlice[i])
	n := p.NodeSlice[i]
	FC, LC := n.FirstChild, n.LastChild
	n.FirstChild, n.LastChild = nil, nil
	html.Render(pSB, n)
	n.FirstChild, n.LastChild = FC, LC

	return pSB.String()
}

func (p *ParserResults_html) NodeInfo(i int) string {
	if i >= len(p.NodeSlice) {
		return "(indexOverrun)"
	}
	return fmt.Sprintf("<h[%d] lv%d,ch%d,%s>",
		i, p.NodeDepths[i], p.FilePosns[i].Pos, p.NodeDebug(i))
}
*/

/* REF
type Node struct {
Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node
Type      NodeType
DataAtom  atom.Atom
Data      string
Namespace string
Attr      []Attribute }
*/
/*
func (p *ParserResults_html) NodeDebug(i int) string {
	if i >= len(p.NodeSlice) {
		return "(indexOverrun)"
	}
	h := *(p.NodeSlice[i])
	// return fmt.Sprintf("|%+v|", h)
	return fmt.Sprintf("|tp:%d:%s,data:%s,ns:%s,kids:%s,atts:%v|",
		h.Type, NodeTypeString[h.Type], DataOfHtmlNode(&h),
		h.Namespace, SU.Yn(h.FirstChild != nil), h.Attr)
}
*/