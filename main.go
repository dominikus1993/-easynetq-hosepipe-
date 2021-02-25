package main

import "fmt"

func main() {
	mapp := map[string]interface{}{
		"21": 1331,
	}

	modify(mapp)

	elem, ok := mapp["21"].(int)

	fmt.Println("Witam", ok, elem)
}

func modify(headers map[string]interface{}) {
	headers["21"] = 2112
}
