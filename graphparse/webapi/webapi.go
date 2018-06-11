package webapi

import (
	// "encoding/json"
	"net/http"
	// "fmt"
	// "time"
	"github.com/liamzebedee/graphparse/graphparse"
	"github.com/gin-contrib/cors"
	// "bytes"
	// "strings"
	// "go/ast"
	// "go/token"
	// "strconv"
	// "bufio"

	// "github.com/gorilla/mux"
	"github.com/gin-gonic/gin"
)


func webAPI() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	// gin.SetMode(gin.ReleaseMode)

	graph := router.Group("/graph/public")
	graphMw := graphMiddleware{}
	graphMw.Populate()
	graph.GET("/:name", graphMw.Middleware, getGraph)

	// src := graph.PathPrefix("/src/").Subrouter()
	// src.Path("/ast-range").
	// Queries(
	// 	"start", "{start:[0-9]+}",
	// 	"end",   "{end:[0-9]+}",
	// ).HandlerFunc(getASTRange)  
	// src.Path("/").HandlerFunc(getSrc)

	return router
}

func WebAPI(port string) {
	router := webAPI()
	router.Run(":"+port)
}

// Middleware
// ==========

func corsEnabledHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}




// Main
// ====

func getWelcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome."))
}



// Source
// ================

type astRangeRes struct {
	Output string `json:"output"`
}

// func getASTRange(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	posStart, _ := strconv.ParseInt(vars["start"], 0, 64)
// 	posEnd, _ := strconv.ParseInt(vars["end"], 0, 64)
// 	fmt.Println("Showing code from", vars["start"], "to", vars["end"])

// 	_, nodes, _ := prog.PathEnclosingInterval(token.Pos(posStart), token.Pos(posEnd))
// 	buf := new(bytes.Buffer)
// 	ast.Fprint(buf, fset, nodes[:1], ast.NotNilFilter)
	
// 	res := astRangeRes{
// 		Output: buf.String(),
// 	}
// 	json.NewEncoder(w).Encode(res)
// }

// func getSrc(w http.ResponseWriter, r *http.Request) {
// 	// TODO for testing remove later.
// 	for f, src := range fileLookup {
// 		if f == "parse.go" {
// 			json.NewEncoder(w).Encode(src)
// 			return
// 		}
// 	}

// 	http.Error(w, "example file not found", http.StatusBadRequest)
// }



// Graph
// ===============

func getGraph(ctx *gin.Context) {
	graph, _ := ctx.Get("graph")
	// fmt.Println(graph)
	graph = graph.(*graphparse.Graph).ToJson()
	ctx.JSON(http.StatusOK, graph)
}