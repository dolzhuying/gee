package main

import (
	"fmt"
	"net/http"

	"base3/gee"
)

func main(){
	r:=gee.New()
	addr:="127.0.0.1:9999"
	r.GET("/",func(w http.ResponseWriter,r *http.Request){
		fmt.Fprintf(w, "url path= %q\n", r.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	r.Run(addr)

}

