package main

import (
  "fmt"
  "time"
)

func main() {
  //时间戳
  t := time.Now()
  fmt.Println(t.Weekday().String())
}
