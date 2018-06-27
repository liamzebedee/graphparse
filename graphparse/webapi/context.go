package webapi

import (
	"github.com/gin-gonic/gin"
	// "context"
	// "fmt"
	"net/http"
	"github.com/liamzebedee/graphparse/graphparse"
	"encoding/base64"
)

type graphMiddleware struct {
	graphs map[string]*graphparse.Graph
}

// type contextKey int
// const (
// 	graph contextKey = iota
// )

// func NewContext(ctx context.Context, g *graphparse.Graph) context.Context {
// 	return context.WithValue(ctx, graph, g)
// }

// func FromContext(ctx context.Context) (*graphparse.Graph) {
// 	return ctx.Value(graph).(*graphparse.Graph)
// }

type graphReqKey struct {
	path string
}


func NewGraphReqKey(path string) graphReqKey {
	return graphReqKey{path}
}

func (k graphReqKey) Path() string {
	return k.path
}


func (middleware *graphMiddleware) Populate() {
	middleware.graphs = make(map[string]*graphparse.Graph)
}


func (middleware *graphMiddleware) getGraph(req graphReqKey) (*graphparse.Graph, error) {
	if g, ok := middleware.graphs[req.Path()]; !ok {
		// generate graph
		if g, err := graphparse.GenerateCodeGraph(req.Path()); err != nil {
			return nil, err
		} else {
			middleware.graphs[req.Path()] = g
			return g, nil
		}
	} else {
		return g, nil
	}
}

func (middleware *graphMiddleware) Middleware(ctx *gin.Context) {
	name, err := base64.StdEncoding.DecodeString(ctx.Param("name"))

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	if g, err := middleware.getGraph(NewGraphReqKey(string(name))); err == nil {
		ctx.Set("graph", g)
	} else {
		ctx.AbortWithError(http.StatusNotFound, err)
	}
}