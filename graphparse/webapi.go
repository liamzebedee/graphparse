package graphparse

import (
	"encoding/json"
	"log"
	"net/http"
	"go/ast"
	"go/token"
	"strconv"
	"fmt"
	"bytes"

	"github.com/gorilla/mux"
)

const fileForExperiment = "subnet/server.go"

func WebAPI(port string) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/src", corsEnabledHeaders(showSrc))
	router.HandleFunc("/src/from/{start}/to/{end}", corsEnabledHeaders(getPos))

	log.Fatal(http.ListenAndServe(":" + port, router))
}

func corsEnabledHeaders(fn http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        fn(w, r)
    }
}

func showSrc(w http.ResponseWriter, r *http.Request) {
	src := fileLookup[fileForExperiment]
	json.NewEncoder(w).Encode(src)
}

type posInfo struct {
	// Nodes []ast.Node `json:"output"`
	Output string `json:"output"`
}

func getPos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	posStart, _ := strconv.ParseInt(vars["start"], 0, 32)
	posEnd, _ := strconv.ParseInt(vars["end"], 0, 32)
	fmt.Println("Showing code from", vars["start"], "to", vars["end"])

	_, nodes, _ := prog.PathEnclosingInterval(token.Pos(posStart), token.Pos(posEnd))

	buf := new(bytes.Buffer)
	ast.Fprint(buf, fset, nodes[:1], ast.NotNilFilter)
	
	res := posInfo{
		Output: buf.String(),
	}
	json.NewEncoder(w).Encode(res)
}
