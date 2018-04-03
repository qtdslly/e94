package main

import (
	"fmt"
	"strings"
	"strconv"

	//"background/stock/tools/util"
)

func GetTransCount(money,price float64)(result int64){
	result,_ = strconv.ParseInt(strings.Split(fmt.Sprint( money / price / 100 ),".")[0],10,64)
	return result * 100
}

func main(){

	//fmt.Println(GetTransCount(5235.78,3.2))
	//util.SendEmail("APAAA","<h3>BBBB</h3>")

	ss := "19950903"
	sDate := ss[0:4] + "-" + ss[4:6] + "-" + ss[6:8]

	fmt.Print(sDate)
}