package main

import(
	"net/http"
	"fmt"
	"log"
)

func main(){
	http.HandleFunc("/",indexHandler)
	http.HandleFunc("/hello",helloHandler)
	log.Fatal(http.ListenAndServe(":9999",nil))


}

func indexHandler(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w,"url path= %q\n",r.URL.Path)
}

func helloHandler(w http.ResponseWriter,r *http.Request){
	for k,v:=range r.Header{
		fmt.Fprintf(w,"header[%q]= %q\n",k,v)
	}
}