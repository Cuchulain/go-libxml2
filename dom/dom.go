package dom

import (
	"fmt"
	"sync"

	"github.com/Cuchulain/go-libxml2/clib"
	"github.com/Cuchulain/go-libxml2/types"
	"github.com/Cuchulain/go-libxml2/xpath"
)

var docPool sync.Pool

func init() {
	SetupXPathCallback()
	docPool = sync.Pool{}
	docPool.New = func() interface{} {
		return Document{}
	}
}

func SetupXPathCallback() {
	xpath.WrapNodeFunc = WrapNode
}

func WrapDocument(n uintptr) *Document {
	doc := docPool.Get().(Document)
	doc.mortal = false
	doc.ptr = n
	return &doc
}

func wrapNamespace(n uintptr) *Namespace {
	ns := Namespace{}
	ns.ptr = n
	return &ns
}

func wrapAttribute(n uintptr) *Attribute {
	attr := Attribute{}
	attr.ptr = n
	return &attr
}

func wrapCDataSection(n uintptr) *CDataSection {
	cdata := CDataSection{}
	cdata.ptr = n
	return &cdata
}

func wrapComment(n uintptr) *Comment {
	comment := Comment{}
	comment.ptr = n
	return &comment
}

func wrapElement(n uintptr) *Element {
	el := Element{}
	el.ptr = n
	return &el
}

func wrapText(n uintptr) *Text {
	txt := Text{}
	txt.ptr = n
	return &txt
}

// WrapNode is a function created with the sole purpose of allowing
// go-libxml2 consumers that can generate a C.xmlNode pointer to
// create libxml2.Node types, e.g. go-xmlsec.
func WrapNode(n uintptr) (types.Node, error) {
	switch typ := clib.XMLGetNodeTypeRaw(n); typ {
	case clib.AttributeNode:
		return wrapAttribute(n), nil
	case clib.ElementNode:
		return wrapElement(n), nil
	case clib.TextNode:
		return wrapText(n), nil
	case clib.CDataSectionNode:
		return wrapCDataSection(n), nil
	default:
		return nil, fmt.Errorf("unknown node: %d", typ)
	}
}
