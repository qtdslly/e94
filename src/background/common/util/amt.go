package util

import "strconv"

func GetAmt(amt string)(float64){
  result ,err :=strconv.ParseFloat(amt,10)
  if err != nil{
    return 0.00
  }
  return result
}
