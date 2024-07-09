package file

import (
	"fmt"
	"strconv"
)

func ContainsToStr(target string, arr []string) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}

// func contains(target int, arr []int) bool {
// 	for _, val := range arr {
// 		if val == target {
// 			return true
// 		}
// 	}
// 	return false
// }

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
	}
	return i
}

func containsTo(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
