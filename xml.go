package parseutils

import (
	"encoding/xml"
	"fmt"
	"io"
	S "strings"
)

type ConcreteParseResults_xml struct {
  // ParseTree ??
  NodeList   []xml.Token
  NodeDepths []int
	Raw string
}

func GetParseResults_xml(s string) (*ConcreteParseResults_xml, error) {
	var nl []xml.Token
	var e error
	nl, e = DoParse_xml(s)
	if e != nil {
		return nil, fmt.Errorf("pu.xml.parseResults: %w", e)
	}
	p := new(ConcreteParseResults_xml)
	p.NodeList = nl
	p.Raw = s
	return p, nil
}

// DoParse_xml takes a string, so we can assume that we can
// discard it after use cos the caller has another copy of it.
// To be safe, it copies every token using `xml.CopyToken(T)`.
func DoParse_xml(s string) (xtokens []xml.Token, err error) {
	var e error
	var T, TT xml.Token
	xtokens = make([]xml.Token, 0, 100)
	// println("(DD) XmlTokenizeBuffer:", s)

	r := S.NewReader(s)
	var parser = xml.NewDecoder(r)
	// Strict mode does not enforce XML namespace requirements. In parti-
	// cular it does not reject namespace tags that use undefined prefixes.
	// Such tags are recorded with the unknown prefix as the namespace URL.
	parser.Strict = false
	// When Strict == false, AutoClose is a set of elements to consider
	// closed immediately after they are opened, regardless of whether
	// an end element is present. For example, <br/>.
	// TODO Add anything for LwDITA ?
	parser.AutoClose = xml.HTMLAutoClose
	// Entity can map non-standard entity names to string replacements.
	// The parser is preloaded with the following standard XML mappings,
	// whether or not they are also provided in the actual map content:
	//	"lt": "<", "gt": ">", "amp": "&", "apos": "'", "quot": `"`
	// NOTE It doesn't do parameter entities, and we havnt necessarily
	// parsed any entities at all yet, so don't bother trying to use this.
	// NOTE If you dump all these, you find that there's a zillion of'em.
	parser.Entity = xml.HTMLEntity

	for {
		T, e = parser.Token()
		if e == io.EOF { break }
		if e != nil {
			return xtokens, fmt.Errorf("pu.xml.doParse: %w", e)
		}
		TT = xml.CopyToken(T)
		xtokens = append(xtokens, TT)
	}
	return xtokens, nil
}
