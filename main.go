package main

import (
	"github.com/liamzebedee/graphparse/graphparse"
)
// idea: actually model it like webpages
//       looking for a piece of code? use current scope as starting page
//       use the type system to autofill the vars? 

// other thing with VR:
// need a visual shape-based constraint/design language


// two approaches:
// - big global canonical lookup table of scope (bad, error-prone)
// - small recursive progressive graph AST approach (better)



// next goal: get methods on structs to be linked back to the struct type
// really what this entails is changing the canonicalisation of the id



func main() {
	graphparse.Stuff()
}