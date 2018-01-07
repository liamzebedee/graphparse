package graphparse

import (
	"strconv"
	"fmt"
)

type rankPair struct {
	NodeId nodeid
	Rank float64
}

// A slice of pairs that implements sort.Interface to sort by values
type rankPairList []rankPair

func (p rankPairList) Len() int           { return len(p) }
func (p rankPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p rankPairList) Less(i, j int) bool { return p[i].Rank < p[j].Rank }




func pointerToStr(ptr interface{}) string {
	return fmt.Sprintf("%p", ptr)
}
func pointerToId(ptr interface{}) nodeid {
	if ptr == nil {
		panic("ptr cannot be nil")
	}
	if i, err := strconv.ParseInt(pointerToStr(ptr), 0, 64); err != nil {
		panic(err)
	} else {
		return nodeid(i)
	}
}