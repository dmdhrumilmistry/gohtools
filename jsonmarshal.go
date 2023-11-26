package main

import (
	"encoding/json"
	"fmt"
)

type Foo struct {
	Bar string
	Baz string
}

func main() {
	f := Foo{"Bar", "Bax"}
	d, _ := json.Marshal(f)
	fmt.Println(string(d))
	json.Unmarshal(d, &f)
	fmt.Println(f)
}
