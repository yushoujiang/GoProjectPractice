package Pratice

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("k", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	fmt.Fprintf(w, "Hello World")

}

func StartMyWebServer() {

	http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/login", loginTest)

	err := http.ListenAndServe(":9091", nil)

	if err != nil {
		log.Fatal("MyWebError:", err)
	}
}

func loginTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	method := r.Method
	if method == "GET" {

		curTime := time.Now().Unix()
		fmt.Println("curTime=", curTime)
		value := md5.New()
		fmt.Println("value1=", value)
		io.WriteString(value, strconv.FormatInt(curTime, 10))
		fmt.Println("value2=", value)
		token := fmt.Sprintf("%x", value.Sum(nil))

		fmt.Println("token=", token)
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, token)
	} else if method == "POST" {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])

		fmt.Println("12333P:", r.Form.Get("username"))

	}
}
