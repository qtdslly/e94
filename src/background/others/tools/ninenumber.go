package main

import (
	"fmt"
	"sort"
)

func main(){
	numbers := []int{1,1,3,4,5,6,7,8,9}

	if !isLegal(numbers){
		fmt.Println("无解")
	}
}

/*
功能    :检查给定的9个数字是否有解
算法描述:9个数字，9宫格里横竖之和都相等，那么一定符合下面的填法
294
753
618
将给定的数字按大小升序，按上面的排法检验各横竖之和是否等于最大值+最小值+中间数之后即可
*/
func isLegal(numbers []int)bool{
	sort.Ints(numbers)
	sum := numbers[0] + numbers[4] + numbers[8]
	if numbers[1] + numbers[8] + numbers[3] != sum{
		return false
	}
	if numbers[6] + numbers[4] + numbers[2] != sum{
		return false
	}
	if numbers[5] + numbers[0] + numbers[7] != sum{
		return false
	}
	if numbers[1] + numbers[6] + numbers[5] != sum{
		return false
	}
	if numbers[3] + numbers[2] + numbers[7] != sum{
		return false
	}
	if numbers[1] + numbers[4] + numbers[7] != sum{
		return false
	}
	if numbers[5] + numbers[4] + numbers[3] != sum{
		return false
	}
	return true
}