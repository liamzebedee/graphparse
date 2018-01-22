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


type arraylist struct {
	counter int
	backing map[int]interface{}
}

func newArrayList() *arraylist {
	return &arraylist{
		counter: 0,
		backing: make(map[int]interface{}),
	}
}

func (l *arraylist) append(v interface{}) {
	l.backing[l.counter] = v
	l.counter++
	// fmt.Println("app", l.counter)
}

func (l *arraylist) delete(i int) {
	// fmt.Println("del", i)
	delete(l.backing, i)
}

// func (l arraylist) toArray() (arr []interface{}) {
// 	arr = []interface{}{}
// 	for _, v := range l.backing {
// 		arr = append(arr, v)
// 	}
// 	return arr
// }

func (l *arraylist) Map() map[int]interface{} {
	return l.backing
}

func (l *arraylist) get(i int) interface{} {
	return l.backing[i]
}
