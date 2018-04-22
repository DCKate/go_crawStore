package main

import (
	"go_findstuff/server"
	"log"
	"net/http"
)

const ServeAddr = "localhost:55555"

func main() {
	// url path: /search
	// support GET, POST
	// Parameter	key(must) : keyword for search
	//				limit(option) : the number of returning products
	//				store(option) : assign which  store, now is support rt-mart and carrefour
	// ex.
	// curl http://localhost:55555/search\?key\="豬肉"\&limit\=20\&store\="rt-mart"
	// curl -XPOST -d'{"key":"豬肉","limit":20,"store":["rt-mart"]}' http://localhost:55555/search
	http.HandleFunc("/search", server.ProductsSearchHandler)
	// http.HandleFunc("/detail", server.ProductDetailHandler)
	log.Fatal(http.ListenAndServe(ServeAddr, nil))

}
