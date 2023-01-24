package parseutils

import (
	"fmt"
	S "strings"

	XU "github.com/fbaube/xmlutils"
	"golang.org/x/net/html"
)

type ParserResults_html struct {
	RootNode  *html.Node
	NodeSlice []*html.Node
	XU.CommonCPR
}

func GenerateParserResults_html(s string) (*ParserResults_html, error) {
	var root *html.Node
	var e error
	root, e = DoParseTree_html(s)
	if e != nil {
		return nil, fmt.Errorf("pu.html.parseResults.ParseTree: %w", e)
	}
	var nl []*html.Node
	var il []int
	var fp []*XU.FilePosition
	nl, il, fp, e = FlattenParseTree_html(root)
	if e != nil {
		return nil, fmt.Errorf("pu.html.parseResults.flattenTree: %w", e)
	}
	p := new(ParserResults_html)
	p.CommonCPR = *XU.NewCommonCPR()
	p.RootNode = root
	p.NodeSlice = nl
	p.NodeDepths = il
	p.FilePosns = fp
	if fp == nil {
		panic("OOPS fp")
	}
	p.CPR_raw = s
	return p, nil
}

// DoParseTree_html returns the parse tree for the HTML from the given string.
// The input is assumed to be UTF-8 encoded.
func DoParseTree_html(s string) (*html.Node, error) {
	var root *html.Node
	var e error
	root, e = html.Parse(S.NewReader(s))
	if e != nil {
		return nil, fmt.Errorf("parseutils.html.DoParseTree: %w", e)
	}
	return root, nil
}

// var HNdTypes = []string{"nil", "Blk", "Inl", "Doc"}

// hn = HTML Node
var hnList []*html.Node
var hnDepths []int
var hnFPosns []*XU.FilePosition
var hnError error
var hnWalkLevel int

func FlattenParseTree_html(pHN *html.Node) ([]*html.Node, []int, []*XU.FilePosition, error) {
	hnList = make([]*html.Node, 0)
	hnDepths = make([]int, 0)
	hnFPosns = make([]*XU.FilePosition, 0)
	HtmlWalk(pHN, wf_gatherTreeNodes_html)
	return hnList, hnDepths, hnFPosns, hnError
}

// wf_aWalker_html is a func type called when walking a tree of `html.Node`.
// Package `golang.org/x/net/html` does not define walking functions, so
// we define our own in this file.
// `entering` is set `true` before children are walked, `false` after.
// If wf_aWalker_html returns error, Walk function immediately stop walking.
type wf_aWalker_html func(n *html.Node, entering bool) // returns (WalkStatus, error)

// HtmlWalk walks a AST tree by the depth first search algorithm.
func HtmlWalk(n *html.Node, walker wf_aWalker_html) {
	walker(n, true)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		HtmlWalk(c, walker)
	}
	walker(n, false)
}

// wf_gatherTreeNodes_html implements `HtmlWalker`:
// type html.Walker func(n Node, entering bool) (WalkStatus, error)
func wf_gatherTreeNodes_html(n *html.Node, in bool) {
	if in {
		hnWalkLevel += 1
	} else {
		hnWalkLevel -= 1
		return
	}
	hnList = append(hnList, n)
	hnDepths = append(hnDepths, hnWalkLevel)
	hnFPosns = append(hnFPosns, &XU.FilePosition{0, 0, 0})
}

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

func KVpairsFromAttributes_html(atts []html.Attribute) []KVpair {
	//?? var stratts []strattribute
	for _, attr := range atts {
		println("HtmlAttr:", "NS", attr.Namespace, "Key", attr.Key, "Val", attr.Val)
		// litter.Dump(attr)
		// if ok,_ := []uint8{
		/* =================================
			strattr := new(strattribute)
			strattr.Name = string(attr.Name)
			switch attr.Value.(type) {
			case []uint8:
				strattr.Value = string(attr.Value.([]uint8))
			case [][]uint8:
				strattr.Value = ""
				var bbbb [][]byte
				var bb []byte
				bbbb = attr.Value.([][]byte)
				for _, bb = range bbbb {
					strattr.Value += string(bb) // attr.Value.([]uint8))
				}
			}
			stratts = append(stratts, *strattr)
		}
		return stratts
		*/
	}
	return nil
}

var NodeTypeString = []string{
	"Error", "CData", "Docmt", "Elmnt", "Comnt", "Doctype",
}

func NTstring(nt html.NodeType) string {
	return NodeTypeString[nt]
}
