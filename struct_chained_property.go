//------------------------------------------------------------------------------
// chained properties allow us to keep track of the properties
// present in a JSON file and how they are ordered
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"reflect"
)

type chainedProperty struct {
	owner          *fileMap
	name           string
	path           path // i.e. with the path to this property when it's nested
	kind           reflect.Kind
	previous       *chainedProperty
	next           *chainedProperty
	addOn          bool        // when a new property is inserted in an already existing common definition
	index          int         // the global index for this property within the common definition
	conf           *configItem // the config associated with this property
	maxLength      int         // the maximum size of the column
	statistic      *stat       // this property seen as a statistical variable
	computed       bool        // if true, then is property is computed from other columns
	computationDef interface{} // if this is a computed property, then how to compute stuff is set here
}

func (thisProperty *chainedProperty) String() string {
	if thisProperty == nil {
		return "<nil>"
	}
	return thisProperty.name
}

func (thisProperty *chainedProperty) getPath() path {
	if thisProperty.path == "" {
		thisProperty.path = path(fmt.Sprintf("%s%s", thisProperty.owner.getPath(), thisProperty.name))
	}
	return thisProperty.path
}

// chaining this property right after the given targeted property
func (thisProperty *chainedProperty) linkAfter(target *chainedProperty, verbose bool) {
	if verbose {
		log("--> new linking : %s -> %s", target, thisProperty)
	}
	target.next = thisProperty
	thisProperty.previous = target
}

// inserting after the given property, and thus, before its next property
func (thisProperty *chainedProperty) insertAfter(target *chainedProperty) {
	targetNext := target.next
	thisProperty.linkAfter(target, false)
	if targetNext != nil {
		targetNext.linkAfter(thisProperty, false)
		log("--> insertion   : %s -> %s -> %s", target, thisProperty, targetNext)
	} else {
		log("--> insertion   : %s -> %s", target, thisProperty)
	}
}

// inserting before the given property, and thus, after its previous one
func (thisProperty *chainedProperty) insertBefore(target *chainedProperty) {
	targetPrevious := target.previous
	thisProperty.linkAfter(targetPrevious, false)
	target.linkAfter(thisProperty, false)
	log("--> insertion : %s -> %s -> %s", targetPrevious, thisProperty, target)
}

// finding the root
func (thisProperty *chainedProperty) root() *chainedProperty {
	if thisProperty.previous == nil {
		return thisProperty
	}
	return thisProperty.previous.root()
}

// is this propertu equal to the other ?
func (thisProperty *chainedProperty) equals(other *chainedProperty) bool {
	if (thisProperty == nil) != (other == nil) {
		return false
	}
	return thisProperty == nil || thisProperty.name == other.name
}

// is this property right before the other ?
func (thisProperty *chainedProperty) touches(other *chainedProperty) bool {
	return thisProperty.next.equals(other)
}
