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
	//				price_range(option) : 1 for (500-), 2 for (501-1000), 3 for (1001-3000)
	//									, 4 for (3001-10000), 5 for (10001-30000), 6 for (30001+)
	// ex.
	// curl http://localhost:55555/search\?key\="冰箱"\&limit\=20\&store\="rt-mart"&price_range\=5
	// curl -XPOST -d'{"key":"冰箱","limit":5,"store":["rt-mart"],"price_range":5}' http://localhost:55555/search
	http.HandleFunc("/search", server.ProductsSearchHandler)
	// http.HandleFunc("/detail", server.ProductDetailHandler)
	log.Fatal(http.ListenAndServe(ServeAddr, nil))

}
