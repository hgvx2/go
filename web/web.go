package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func main()  {
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login", login) // 设置访问的路由
	err := http.ListenAndServe(":9090", nil)
	if nil != err {
		fmt.Println("ListenAndServer: ", err)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm();
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val: ", strings.Join(v, " "))
	}
	fmt.Fprintf(w, "Hello astaxie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("E:\\GoWork\\im\\src\\web\\login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// 请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
