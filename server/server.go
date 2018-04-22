package server

import (
	"encoding/json"
	"go_findstuff/controller"
	"log"
	"net/http"
	"strconv"
)

type errorRespData struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func toJSON(m interface{}) []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return js
}

func parseRequesyData(r *http.Request) map[string]interface{} {
	mret := make(map[string]interface{})
	mtmp := make(map[string]interface{})
	switch r.Method {
	case "GET":
		vals := toJSON(r.URL.Query())
		err := json.Unmarshal([]byte(vals), &mtmp)
		if err != nil {
			panic(err)
		}
	case "POST":
		err := json.NewDecoder(r.Body).Decode(&mtmp)
		if err != nil {
			panic(err)
		}
	}

	for kk, vv := range mtmp {
		switch vv.(type) {
		case string:
			switch kk {
			case "key":
				mret[kk] = vv.(string)
			case "store":
				mret[kk] = []interface{}{vv.(string)}
			case "limit":
				lNum, err := strconv.Atoi(vv.(string))
				if err == nil {
					mret[kk] = lNum
				}
			}
		case []interface{}:
			switch kk {
			case "key":
				mret[kk] = vv.([]interface{})[0].(string)
			case "store":
				mret[kk] = vv
			case "limit":
				lNum, err := strconv.Atoi(vv.([]interface{})[0].(string))
				if err == nil {
					mret[kk] = lNum
				}
			}
		}
	}
	return mret
}

func ProductsSearchHandler(w http.ResponseWriter, r *http.Request) {
	mpara := parseRequesyData(r)
	log.Println(mpara)

	if _, ok := mpara["key"]; ok {
		ri, rdata := controller.ProductsSearchController(mpara)
		if ri == 0 {
			json.NewEncoder(w).Encode(&rdata)
			return
		}
		msg := string(toJSON(errorRespData{
			Code: -2,
			Msg:  "query not found",
		}))
		http.Error(w, msg, http.StatusNotFound)
	} else {
		msg := string(toJSON(errorRespData{
			Code: -1,
			Msg:  "missing key for search",
		}))
		http.Error(w, msg, http.StatusBadRequest)
	}
}

// func ProductDetailHandler(w http.ResponseWriter, r *http.Request) {

// }
