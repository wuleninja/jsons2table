//------------------------------------------------------------------------------
// keeping some stats for each property
//------------------------------------------------------------------------------

package main

type statKind string

const (
	statKindTEXT     statKind = "text"
	statKindCATEGORY statKind = "category"
	statKindBOOLEAN  statKind = "boolean"
	statKindDATE     statKind = "date"
	statKindNUMBER   statKind = "number"
)

type stat struct {
	owner       *chainedProperty
	valueCounts map[string]int
	kind        statKind
}
