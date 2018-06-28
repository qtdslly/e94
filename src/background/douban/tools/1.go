package main

import "fmt"

func main(){
	//title := "129分钟 / 111分钟(中国大陆)"
	//title := "0123456789"  //48-57
	//title := "-" //45
	//title := "+"  //43
	//title := "ABCDEFGWYZ" //65-90
	title := "abcwyz"  //97-122
	//48-57 45 43 65-90 97-122
	rTitle := ([]rune)(title)
	for _, m := range rTitle {
		fmt.Println(m)
	}
}
