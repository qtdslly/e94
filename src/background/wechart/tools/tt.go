package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := "真相漩涡纺纱人真相旋涡"  // SpinningMan
	var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]$")
	fmt.Println(hzRegexp.MatchString(str))
}