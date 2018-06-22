package main

import "fmt"

func main(){
	//title := "129分钟 / 111分钟(中国大陆)"
	title := "/"  //48-57
	rTitle := ([]rune)(title)
	for _, m := range rTitle {
		fmt.Println(m)
	}
}
