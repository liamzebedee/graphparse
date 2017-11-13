package graphparse

type rankPair struct {
	NodeId nodeid
	Rank float64
}

// A slice of pairs that implements sort.Interface to sort by values
type rankPairList []rankPair

func (p rankPairList) Len() int           { return len(p) }
func (p rankPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p rankPairList) Less(i, j int) bool { return p[i].Rank < p[j].Rank }