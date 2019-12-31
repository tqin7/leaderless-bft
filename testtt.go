package main

import (
	"encoding/json"
	"fmt"
)


type t struct {
	S string 	
	I int
	Q int64
}

func main() {


	var abc t
	abc.S = "123"
	abc.I = 1
	abc.Q = 7788
	fmt.Println(abc)
	buf, _ := json.Marshal(abc)
	//fmt.Println(reflect.TypeOf(buf))
	fmt.Println(string(buf))
	str_buf := string(buf)
	fmt.Println(str_buf)
	var bcd t
	_ = json.Unmarshal(buf, &bcd)
	//fmt.Println(bcd)



}
