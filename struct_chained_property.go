//------------------------------------------------------------------------------
// chained properties allow us to keep track of the properties
// present in a JSON file and how they are ordered
//------------------------------------------------------------------------------

package main

import "reflect"

type chainedProperty struct {
	owner    *fileMap
	name     string
	kind     reflect.Kind
	previous *chainedProperty
	next     *chainedProperty
	addOn    bool // when a new property is inserted in an already existing common definition
}

func (thisProperty *chainedProperty) String() string {
	if thisProperty == nil {
		return "<nil>"
	}
	return thisProperty.name
}

func (thisProperty *chainedProperty) FullString() string {
	result := thisProperty.name
	for owner := thisProperty.owner; owner != nil; owner = owner.parent {
		result = owner.name + "." + result
	}
	return result
}

// chaining this property right after the given targeted property
func (thisProperty *chainedProperty) linkAfter(target *chainedProperty, verbose bool) {
	if verbose {
		debug("--> new linking : %s -> %s", target, thisProperty)
	}
	target.next = thisProperty
	thisProperty.previous = target
}

// inserting after the given property, and thus, before its next property
func (thisProperty *chainedProperty) insertAfter(target *chainedProperty) {
	targetNext := target.next
	thisProperty.linkAfter(target, false)
	targetNext.linkAfter(thisProperty, false)
	debug("--> insertion   : %s -> %s -> %s", target, thisProperty, targetNext)
}

// inserting before the given property, and thus, after its previous one
func (thisProperty *chainedProperty) insertBefore(target *chainedProperty) {
	targetPrevious := target.previous
	thisProperty.linkAfter(targetPrevious, false)
	target.linkAfter(thisProperty, false)
	debug("--> insertion : %s -> %s -> %s", targetPrevious, thisProperty, target)
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
