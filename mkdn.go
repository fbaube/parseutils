package parseutils

import (
	"fmt"

	XM "github.com/fbaube/xmlmodels"
	GM "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	_ "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	RRR "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// ParserResults_mkdn is a bit dodgy cos
// `ast.Node` is an interface, not a struct.
type ParserResults_mkdn struct {
	RootNode  ast.Node
	NodeSlice []ast.Node
	Reader    text.Reader
	XM.CommonCPR
}

// mn = MarkdownNode
var mnList []ast.Node
var mnDepths []int
var mnFilPosns []*XM.FilePosition
var mnError error
var mnWalkLevel int

// var theSRC []byte // string
// var NdKdNms []string
var TheSourceBfr []byte
var TheSourceAfr []byte
var TheReader text.Reader
var r RRR.Renderer

var MNdTypes = []string{"nil", "Blk", "Inl", "Doc"}

func GenerateParserResults_mkdn(s string) (*ParserResults_mkdn, error) {
	var root ast.Node
	var rdr text.Reader
	var e error
	root, rdr, e = DoParseTree_mkdn(s)
	if e != nil {
		return nil, fmt.Errorf("pu.mkdn.parseResults.ParseTree: %w", e)
	}
	var nl []ast.Node
	var il []int
	var fp []*XM.FilePosition
	nl, il, fp, e = FlattenParseTree_mkdn(root)
	if e != nil {
		return nil, fmt.Errorf("pu.mkdn.parseResults.flattenTree: %w", e)
	}
	p := new(ParserResults_mkdn)
	p.CommonCPR = *XM.NewCommonCPR()
	p.RootNode = root
	p.NodeSlice = nl
	p.NodeDepths = il
	p.FilePosns = fp
	p.Reader = rdr
	p.CPR_raw = s
	return p, nil
}

// DoParseTree_mkdn takes a string and returns the tree produced by the parser.
func DoParseTree_mkdn(s string) (ast.Node, text.Reader, error) {
	var GoldMarkDown GM.Markdown
	// var REND renderer.Renderer
	GoldMarkDown = GM.New(
		GM.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Strikethrough,
			extension.Table,
			extension.Linkify,
			extension.TaskList,
		),
		GM.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		GM.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithWriter(html.DefaultWriter), // os.Stdout),
		),
	)
	var TheParser parser.Parser
	var TheParseTree ast.Node
	TheSourceBfr = []byte(s) // p.CheckedPath.Raw)
	// println("TheSource:", litter.Sdump(s))
	// r = GM.DefaultRenderer() // GoldMarkDown.Renderer().(html.Renderer)
	r = GoldMarkDown.Renderer()
	TheParser = GoldMarkDown.Parser()
	TheReader = text.NewReader(TheSourceBfr)
	TheParseTree = TheParser.Parse(TheReader)
	TheSourceAfr = TheReader.Source()
	// Put a tree dump here ?
	TheReader.ResetPosition()
	return TheParseTree, TheReader, nil // pMTokzn, nil
}

func FlattenParseTree_mkdn(pMN ast.Node) ([]ast.Node, []int, []*XM.FilePosition, error) {
	mnList = make([]ast.Node, 0)
	mnDepths = make([]int, 0)
	mnFilPosns = make([]*XM.FilePosition, 0)
	e := ast.Walk(pMN, wf_gatherTreeNodes_mkdn)
	if e != nil {
		panic(e)
	}
	return mnList, mnDepths, mnFilPosns, nil
}

// wf_aWalker_mkdn is a func type called when walking a tree of `ast.Node`.
// `entering` is set `true` before children are walked, and `false` after.
// If wf_aWalker_mkdn returns error, the walk function immediately stops walking.
type wf_aWalker_mkdn func(n ast.Node, entering bool) // returns (ast.WalkStatus, error)

// wf_gatherTreeNodes_mkdn is ::
// type ast.Walker func(n Node, entering bool) (WalkStatus, error)
// NOTE `ast.Node` is an interface!
func wf_gatherTreeNodes_mkdn(n ast.Node, in bool) (ast.WalkStatus, error) {
	if in {
		mnWalkLevel += 1
	} else {
		mnWalkLevel -= 1
		return ast.WalkContinue, nil
	}
	mnList = append(mnList, n)
	mnDepths = append(mnDepths, mnWalkLevel)
	mnFilPosns = append(mnFilPosns, &XM.FilePosition{0, 0, 0})
	return ast.WalkContinue, nil
}

// GetAllByTag returns a slice of `ast.Node`. It checks the basic tag only,
// not any namespace. Note that these tag lookup func's default to searching
// the `ListNodesP`, not the tree of `Node`s.
func (pCPR *ParserResults_mkdn) GetAllByTag(s string) []ast.Node {
	if s == "" {
		return nil
	}
	var ret = make([]ast.Node, 0)
	for _, p := range pCPR.NodeSlice {
		panic(fmt.Sprintf("OOPS: %p", p))
	}
	return ret
}

func KVpairsFromAttributes_mkdn(atts []ast.Attribute) []KVpair {
	var KVpairs []KVpair
	for _, attr := range atts {
		// litter.Dump(attr)
		// if ok,_ := []uint8{
		kvp := new(KVpair)
		kvp.Key = string(attr.Name)
		switch attr.Value.(type) {
		case []uint8:
			kvp.Key = string(attr.Value.([]uint8))
		case [][]uint8:
			kvp.Key = ""
			var bbbb [][]byte
			var bb []byte
			bbbb = attr.Value.([][]byte)
			for _, bb = range bbbb {
				kvp.Value += string(bb) // attr.Value.([]uint8))
			}
		}
		KVpairs = append(KVpairs, *kvp)
	}
	return KVpairs
}
