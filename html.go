package parseutils

import (
	"fmt"
	S "strings"

	"golang.org/x/net/html"
)

type ConcreteParseResults_html struct {
	ParseTree  *html.Node
	NodeList   []*html.Node
	NodeDepths []int
	CPR_raw    string
}

func GetConcreteParseResults_html(s string) (*ConcreteParseResults_html, error) {
	var root *html.Node
	var e error
	root, e = DoParseTree_html(s)
	if e != nil {
		return nil, fmt.Errorf("pu.html.parseResults.ParseTree: %w", e)
	}
	var nl []*html.Node
	var il []int
	nl, il, e = FlattenParseTree_html(root)
	if e != nil {
		return nil, fmt.Errorf("pu.html.parseResults.flattenTree: %w", e)
	}
	p := new(ConcreteParseResults_html)
	p.ParseTree = root
	p.NodeList = nl
	p.NodeDepths = il
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

var HNdTypes = []string{"nil", "Blk", "Inl", "Doc"}

// hn = HTML Node
var hnList []*html.Node
var hnDepths []int
var hnError error
var hnWalkLevel int

func FlattenParseTree_html(pHN *html.Node) ([]*html.Node, []int, error) {
	hnList = make([]*html.Node, 0)
	hnDepths = make([]int, 0)
	HtmlWalk(pHN, wf_gatherTreeNodes_html)
	return hnList, hnDepths, hnError
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
}

func (pCPR *ConcreteParseResults_html) GetAllByAnyTag(ss []string) []*html.Node {
	if ss == nil || len(ss) == 0 {
		return nil
	}
	var ret = make([]*html.Node, 0)
	for _, p := range pCPR.NodeList {
		panic(fmt.Sprintf("OOPS: %p", p))
	}
	return ret
}

// GetAllByTag returns a slice of `*html.Node`. It checks the basic tag only,
// not any namespace. Note that these tag lookup func's default to searching
// the `ListNodesP`, not the tree of `Node`s.
func (pCPR *ConcreteParseResults_html) GetAllByTag(s string) []*html.Node {
	if s == "" {
		return nil
	}
	var ret = make([]*html.Node, 0)
	for _, p := range pCPR.NodeList {
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
	"Err", "Txt", "Doc", "Elm", "Cmt", "Doctype",
}

func NTstring(nt html.NodeType) string {
	return NodeTypeString[nt]
}

/*
	type html.Node struct {
    Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

    Type      NodeType
    DataAtom  atom.Atom
    Data      string
    Namespace string
    Attr      []Attribute
	}
	type my.HtmlToken struct {
		NodeDepth    int // from node walker
		NodeType     string
		NodeKind     string
		NodeKindEnum ast.NodeKind
		NodeKindInt  int
		// NodeText is the text of the MD node,
		//  and it is not present for all nodes.
		NodeText string
		// DitaTag and HtmlTag are the equivalent LwDITA and (X)HTML tags,
		// possibly with an attribute specified too. sDitaTag is authoritative;
		// sHtmlTag is provided mainly as an aid to understanding the code.
		DitaTag string
		NodeNumeric      int // Headings, Emphasis, ...?
	}

	var ndType html.NodeType = n.Type()
	p.NodeType = HNdTypes[ndType]
	p.NodeKindEnum = n.Kind()
	p.NodeKindInt = int(p.NodeKindEnum)
	var SB S.Builder
	var pfx = S.Repeat(" * ", p.NodeDepth-1)
	var sDump string
	// Stack of headers. [0] = the document.

	if DEBUG { println("\n---", pfx, "---") }

	switch p.NodeKindEnum {

	case html.KindAutoLink:
		// https://github.github.com/gfm/#autolinks
		// Autolinks are absolute URIs and email addresses btwn < and >.
		// They are parsed as links, with the link target reused as the link label.
		p.NodeKind = "KindAutoLink"
		p.DitaTag = "xref"
		p.HtmlTag = "a@href"
		n2 := n.(*html.AutoLink)
		sDump = litter.Sdump(*n2)
		// type AutoLink struct {
		//   BaseInline
		//   Type is a type of this autolink.
		//   AutoLinkType AutoLinkType
		//   Protocol specified a protocol of the link.
		//   Protocol []byte
		//   value *Text
		// }
		// w.WriteString(`<a href="`)
		// url := n.URL(source)
		// label := n.Label(source)
		// if n.AutoLinkType == html.AutoLinkEmail &&
		//    !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
		//   w.WriteString("mailto:")
		// }
		// w.Write(util.EscapeHTML(util.URLEscape(url, false)))
		// w.WriteString(`">`)
		// w.Write(util.EscapeHTML(label))
		// w.WriteString(`</a>`)
	case html.KindBlockquote:
		p.NodeKind = "KindBlockquote"
		p.DitaTag = "?blockquote"
		p.HtmlTag = "blockquote"
		n2 := n.(*html.Blockquote)
		sDump = litter.Sdump(*n2)
		// type Blockquote struct {
		//   BaseBlock
		// }
		// w.WriteString("<blockquote>\n")
	case html.KindCodeBlock:
		p.NodeKind = "KindCodeBlock"
		p.DitaTag = "?pre+?code"
		p.HtmlTag = "pre+code"
		n2 := n.(*html.CodeBlock)
		sDump = litter.Sdump(*n2)
		// type CodeBlock struct {
		//   BaseBlock
		// }
		// w.WriteString("<pre><code>")
		// r.writeLines(w, source, n)
	case html.KindCodeSpan:
		p.NodeKind = "KindCodeSpan"
		p.DitaTag = "?code"
		p.HtmlTag = "code"
		// // n2 := n.(*html.CodeSpan)
		// // sDump = litter.Sdump(*n2)
		// type CodeSpan struct {
		//   BaseInline
		// }
		// w.WriteString("<code>")
		// for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		//   segment := c.(*html.Text).Segment
		//   value := segment.Value(source)
		//   if bytes.HasSuffix(value, []byte("\n")) {
		//     r.Writer.RawWrite(w, value[:len(value)-1])
		//     if c != n.LastChild() {
		//       r.Writer.RawWrite(w, []byte(" "))
		//     }
		//   } else {
		//     r.Writer.RawWrite(w, value)
	case html.KindDocument:
		// Note that metadata comes btwn this
		// start-of-document tag and the content ("body").
		p.NodeKind = "KindDocument"
		p.DitaTag = "topic"
		p.HtmlTag = "html"
	case html.KindEmphasis:
		p.NodeKind = "KindEmphasis"
		// iLevel 2 | iLevel 1
		p.DitaTag = "b|i"
		p.HtmlTag = "strong|em"
		n2 := n.(*html.Emphasis)
		sDump = litter.Sdump(*n2)
		p.NodeNumeric = n2.Level
		// type Emphasis struct {
		//   BaseInline
		//   Level is a level of the emphasis.
		//   Level int
		// }
		// tag := "em"
		// if n.Level == 2 {
		//   tag = "strong"
		// }
		// if entering {
		//   w.WriteByte('<')
		//   w.WriteString(tag)
		//   w.WriteByte('>')
	case html.KindFencedCodeBlock:
		p.NodeKind = "KindFencedCodeBlock"
		p.DitaTag = "?code"
		p.HtmlTag = "code"
		n2 := n.(*html.FencedCodeBlock)
		sDump = litter.Sdump(*n2)
		// type FencedCodeBlock struct {
		//   BaseBlock
		//   Info returns a info text of this fenced code block.
		//   Info *Text
		//   language []byte
		// }
		// w.WriteString("<pre><code")
		// language := n.Language(source)
		// if language != nil {
		//   w.WriteString(" class=\"language-")
		//   r.Writer.Write(w, language)
	case html.KindHTMLBlock:
		p.NodeKind = "KindHTMLBlock"
		p.DitaTag = "?htmlblock"
		p.HtmlTag = "?htmlblock"
		n2 := n.(*html.HTMLBlock)
		sDump = litter.Sdump(*n2)
		// type HTMLBlock struct {
		//   BaseBlock
		//   Type is a type of this html block.
		//   HTMLBlockType HTMLBlockType
		//   ClosureLine is a line that closes this html block.
		//   ClosureLine textm.Segment
		// }
		// if r.Unsafe {
		//   l := n.Lines().Len()
		//   for i := 0; i < l; i++ {
		//     line := n.Lines().At(i)
		//     w.Write(line.Value(source))
		//   }
		// } else {
		//   w.WriteString("<!-- raw HTML omitted -->\n")
	case html.KindHeading:
		p.NodeKind = "KindHeading"
		p.DitaTag = "?"
		p.HtmlTag = "h%d"
		n2 := n.(*html.Heading)
		sDump = litter.Sdump(*n2)
		p.NodeNumeric = n2.Level
		// type Heading struct {
		//   BaseBlock
		//   Level returns a level of this heading.
		//   This value is between 1 and 6.
		//   Level int
		// }
	// w.WriteString("<h")
	// w.WriteByte("0123456"[n.Level])
	case html.KindImage:
		p.NodeKind = "KindImage"
		p.DitaTag = "image"
		p.HtmlTag = "img"
		n2 := n.(*html.Image)
		sDump = litter.Sdump(*n2) + " Dest. " + string(n2.Destination)
		// type Image struct {
		//   baseLink
		// }
		// w.WriteString("<img src=\"")
		// if r.Unsafe || !IsDangerousURL(n.Destination) {
		//   w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		// }
		// w.WriteString(`" alt="`)
		// w.Write(n.Text(source))
		// w.WriteByte('"')
		// if n.Title != nil {
		//   w.WriteString(` title="`)
		//   r.Writer.Write(w, n.Title)
		//   w.WriteByte('"')
		// }
		// if r.XHTML {
		//   w.WriteString(" />")
		// } else {
		//   w.WriteString(">")
	case html.KindLink:
		p.NodeKind = "KindLink"
		p.DitaTag = "xref"
		p.HtmlTag = "a@href"
		n2 := n.(*html.Link)
		sDump = litter.Sdump(*n2) + " Dest. " + string(n2.Destination)
		// type Link struct {
		//   baseLink
		// }
		// w.WriteString("<a href=\"")
		// if r.Unsafe || !IsDangerousURL(n.Destination) {
		//   w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		// }
		// w.WriteByte('"')
		// if n.Title != nil {
		//   w.WriteString(` title="`)
		//   r.Writer.Write(w, n.Title)
		//   w.WriteByte('"')
		// }
		// w.WriteByte('>')
	case html.KindList:
		p.NodeKind = "KindList"
		n2 := n.(*html.List)
		sDump = litter.Sdump(*n2)
		if n2.IsOrdered() {
			p.DitaTag = "ol"
			p.HtmlTag = "ol"
		} else {
			p.DitaTag = "ul"
			p.HtmlTag = "ul"
		}
		// type List struct {
		//   BaseBlock
		//   Marker is a markar character like '-', '+', ')' and '.'.
		//   Marker byte
		//   IsTight is a true if this list is a 'tight' list.
		//   See https://spec.commonmark.org/0.29/#loose for details.
		//   IsTight bool
		//   Start is an initial number of this ordered list.
		//   If this list is not an ordered list, Start is 0.
		//   Start int
		// }
		// tag := "ul"
		// if n.IsOrdered() {
		//   tag = "ol"
		// }
		// w.WriteByte('<')
		// w.WriteString(tag)
		// if n.IsOrdered() && n.Start != 1 {
		//   fmt.Fprintf(w, " start=\"%d\">\n", n.Start)
		// } else {
		//   w.WriteString(">\n")
	case html.KindListItem:
		p.NodeKind = "KindListItem"
		n2 := n.(*html.ListItem)
		sDump = litter.Sdump(*n2)
		p.DitaTag = "li"
		p.HtmlTag = "li"
		// type ListItem struct {
		//   BaseBlock
		//   Offset is an offset potision of this item.
		//   Offset int
		// }
		// w.WriteString("<li>")
		// fc := n.FirstChild()
		// if fc != nil {
		//   if _, ok := fc.(*html.TextBlock); !ok {
		//     w.WriteByte('\n')
	case html.KindParagraph:
		p.NodeKind = "KindParagraph"
		p.DitaTag = "p"
		p.HtmlTag = "p"
		// // n2 := n.(*html.Paragraph)
		// // sDump = litter.Sdump(*n2)
		// type Paragraph struct {
		//   BaseBlock
		// }
		// w.WriteString("<p>")
	case html.KindRawHTML:
		p.NodeKind = "KindRawHTML"
		p.DitaTag = "?rawhtml"
		p.HtmlTag = "?rawhtml"
		n2 := n.(*html.RawHTML)
		sDump = litter.Sdump(*n2)
		// type RawHTML struct {
		//   BaseInline
		//   Segments *textm.Segments
		// }
		// if r.Unsafe {
		// n := node.(*html.RawHTML)
		// l := n.Segments.Len()
		// for i := 0; i < l; i++ {
		//   segment := n.Segments.At(i)
		//   w.Write(segment.Value(source))
		// }
	case html.KindText:
		p.NodeKind = "KindText"
		n2 := n.(*html.Text)
		p.DitaTag = "?text"
		p.HtmlTag = "?text"
		// // sDump = litter.Sdump(*n2)
		// type Text struct {
		//   BaseInline
		//   Segment is a position in a source text.
		//   Segment textm.Segment
		//   flags uint8
		// }
		segment := n2.Segment
		// p.NodeText = fmt.Sprintf("KindText:\n | %s", string(TheReader.Value(segment)))
		p.NodeText = /* fmt.Sprintf("KindText:\n | %s", * / string(TheReader.Value(segment)) //)
		/*
			if n.IsRaw() {
				r.Writer.RawWrite(w, segment.Value(TheSource))
			} else {
				r.Writer.Write(w, segment.Value(TheSource))
				if n.HardLineBreak() || (n.SoftLineBreak() && r.HardWraps) {
					if r.XHTML {
						w.WriteString("<br />\n")
					} else {
						w.WriteString("<br>\n")
					}
				} else if n.SoftLineBreak() {
					w.WriteByte('\n')
				}
			}
		* /
	case html.KindTextBlock:
		p.NodeKind = "KindTextBlock"
		p.DitaTag = "?textblock"
		p.HtmlTag = "?textblock"
		// // n2 := n.(*html.TextBlock)
		// // sDump = litter.Sdump(*n2)
		// type TextBlock struct {
		//   BaseBlock
		// }
		// if _, ok := n.NextSibling().(html.Node); ok && n.FirstChild() != nil {
		//   w.WriteByte('\n')
	case html.KindThematicBreak:
		p.NodeKind = "KindThematicBreak"
		p.DitaTag = "hr"
		p.HtmlTag = "hr"
		// type ThemanticBreak struct {
		//   BaseBlock
		// }
		// if r.XHTML {
		//   w.WriteString("<hr />\n")
		// } else {
		//   w.WriteString("<hr>\n")
	default:
		p.NodeKind = "KindUNK"
		p.DitaTag = "UNK"
		p.HtmlTag = "UNK"
	}
	// SB.Write(pfx)
	SB.Write([]byte(fmt.Sprintf("%s:%d:%s // <%s:%s:lvl%d>\n",
		p.NodeType, p.NodeKindEnum, p.NodeKind,
		p.DitaTag, p.HtmlTag, p.NodeNumeric)))
	if len(n.Attributes()) > 0 {
		SAs := strattributesFromAttributes(n.Attributes())
		for _, sa := range SAs {
			SB.Write(SU.B(sa.String()))
		}
	}
	if DEBUG {
		print(SB.String())
		if p.NodeText != "" {
			println("NodeText<<", p.NodeText, ">>")
		}
		if len(sDump) > 1 {
			println(sDump)
			println(" ")
		}
		if ndType == html.TypeBlock {
			println("Block.RawText<< ")
			// TODO Suppress print of all-whitespace
			for i := 0; i < n.Lines().Len(); i++ {
				line := n.Lines().At(i)
				fmt.Printf("%s", line.Value(TheSourceAfr))
				// litter.Dump(theSRC)
			}
			println(" >>")
		}
	}
	// litter.Dump(txt.Segment)
	// }
	// }
	// Text(source []byte) []byte
	// Lines() *textm.Segments (block only)
	// IsRaw() bool
	// Attributes() []Attribute
	// type Attribute struct { Name, Value []byte }

	return false, nil
}
*/

/*
func (p HtmlAST) Echo() string {
	return "HTML ECHO" // p.Node.String()
}
func (p HtmlAST) EchoTo(w io.Writer) {
	w.Write([]byte(p.Echo()))
}
func (p HtmlAST) String() string {
	// return p.Node.String()
	return fmt.Sprintf("%+v", p)
}
func (p HtmlAST) DumpTo(w io.Writer) {
	w.Write([]byte(p.String()))
}
*/
