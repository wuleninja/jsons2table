//------------------------------------------------------------------------------
// the code here is about how we link a new property into an existing definition
//------------------------------------------------------------------------------

package main

import "fmt"

// inserts this property amongst its siblings in the common definition
// Within its original file map, the property was chained like this : prop_a -> prop -> prop_b
// But in the common definition, there might be stuff between prop_a and prop_b, brought by other file maps.
// So the goal here is to correctly find the position of our prop, i.e. obtain something
// like this : prop_a -> prop_a1 -> prop_a2 -> ... -> name -> ... -> prop_b2 -> prop_b1 -> prop_b.
// We're going to use the alphabetic order to find the more suitable place for our prop.
func (newProp *chainedProperty) link(originalProp *chainedProperty) {

	// detecting a serious bug
	if originalProp.previous == nil && originalProp.next == nil {
		panic(fmt.Errorf("original property '%s' was chained to nothing", originalProp))
	}

	debug("\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	debug("==> trying to link property: %s", newProp)

	// the properties already chained in the common definition
	commonProperties := newProp.owner.chainedProperties

	// if our original property had nothing preceding it, we just link it just before the root
	if originalProp.previous == nil {

		root := newProp.owner.oneChainedProperty().root()
		debug("--> linking before root %s", root)
		root.linkAfter(newProp, true)
		return
	}

	// our property had a property before it, that gives a hint about where to start
	// NB: it exists since we necessarily dealt with it previously
	firstPossiblePreviousProp := commonProperties[originalProp.previous.name]

	debug("--> first previous = %s", firstPossiblePreviousProp)
	// what's the max after we can reach ?
	var lastPossibleNextProp *chainedProperty

	// let's look at the next property that exists in the common definition
	for currentNext := originalProp.next; currentNext != nil && lastPossibleNextProp == nil; currentNext = currentNext.next {
		debug("--> last next scan = %s", currentNext)
		lastPossibleNextProp = commonProperties[currentNext.name]
	}

	// if still we have no clue about where to stop, then right after the first previous might be good
	if lastPossibleNextProp == nil {
		debug("--> last possible next prop not found, so getting the next of the first possible previous prop")
		lastPossibleNextProp = firstPossiblePreviousProp.next
	}

	debug("--> last next      = %s", lastPossibleNextProp)

	// in any case, if the first previous was last, or touches the last next property,
	// then we just have to squeeze our property in between !
	if firstPossiblePreviousProp.touches(lastPossibleNextProp) {
		debug("--> 'first previous' and 'last next' already touch each other!")
		if lastPossibleNextProp == nil {
			newProp.linkAfter(firstPossiblePreviousProp, true)
			return
		}
		newProp.insertAfter(firstPossiblePreviousProp)
		return
	}

	debug("--> no easy linking done here, so entering scanning loop")

	// so we basically look for the best previous property for our current prop,
	// starting from the firstPossiblePreviousProp, and finishing at most with the lastPossibleNextProp
	for previous := firstPossiblePreviousProp; !previous.touches(lastPossibleNextProp); previous = previous.next {

		debug("--> in the loop : can we sit betwen '%s' and '%s' ?", previous, previous.next)

		// has the property right after the current previous property a better fit in the alphabetical sense ?
		if next := previous.next; next.name > newProp.name {

			debug("--> in the loop : YEAH!")

			// yeah, so we stop right here
			newProp.insertAfter(previous)
			return
		}

		debug("--> in the loop : nope!")
	}

	debug("--> exited the scanning loop : inserting just right before the last next")

	// we're here because we reached lastPossibleNextProp; so let's just insert our property before this
	newProp.insertBefore(lastPossibleNextProp)
}
