package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Data struct {
	Name string
	Age  int
}
type Ret struct {
	Code  int
	Param string
	Msg   string
	Data  []Data
	Count int
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	// data := Data{Name: "why", Age: 18}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ret := new(Ret)
	id := req.FormValue("id")
	//id := req.PostFormValue('id')

	ret.Code = 0
	ret.Param = id
	ret.Msg = "success"
	result := []Data{}
	for i := 0; i < 115; i++ {
		data := Data{
			Name: fmt.Sprintf("name%d", i),
			Age:  i,
		}
		result = append(result, data)
	}

	page, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("id: %s, %d", id, page)

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * 10
	end := offset + 10
	fmt.Println(11111, offset, end)
	// fmt.Println("result", result)

	ret.Data = result[offset:end]
	ret.Count = len(result)
	ret_json, _ := json.Marshal(ret)

	io.WriteString(w, string(ret_json))
}

func HelloServer1(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world1!\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/hello1", HelloServer1)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
