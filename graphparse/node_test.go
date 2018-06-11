package graphparse

// import (
// 	"testing"
// 	"go/types"
// 	"go/token"
// 	"github.com/stretchr/testify/assert"
// )

// type mockObj struct {
// 	types.Object
// 	i int
// }

// func (obj mockObj) Pos() token.Pos {
// 	return token.Pos(obj.i)
// }

// func TestNodeIds(t *testing.T) {
// 	n1 := CreateNode(mockObj{nil, 0}, Func, "foo")
// 	n2 := LookupNode(mockObj{nil, 0})
// 	assert.Equal(t, n1, n2)
// }