package parseutils

/* IFC
type mcfile.NodeStringer interface {
	NodeEcho(int) string
	NodeInfo(int) string
	NodeDebug(int) string
	NodeCount() int }
type ParserResults_html struct {
	RootNode  *html.Node
	NodeSlice []*html.Node
	XU.CommonCPR }
html.Node.Render(w io.Writer, n *Node) error
*/

import (
	"fmt"
	SU "github.com/fbaube/stringutils"
	"golang.org/x/net/html"
	S "strings"
)

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

/* REF
type Node struct {
Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node
Type      NodeType
DataAtom  atom.Atom
Data      string
Namespace string
Attr      []Attribute }
*/

func (p *ParserResults_html) NodeDebug(i int) string {
	if i >= len(p.NodeSlice) {
		return "(indexOverrun)"
	}
	h := *(p.NodeSlice[i])
	// return fmt.Sprintf("|%+v|", h)
	return fmt.Sprintf("|tp:%d:%s,(%s:%v),ns:%s,kids:%s,atts:%v|",
		h.Type, NodeTypeString[h.Type], h.Data, h.DataAtom,
		h.Namespace, SU.Yn(h.FirstChild != nil), h.Attr)
}
