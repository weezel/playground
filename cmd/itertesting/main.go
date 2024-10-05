package main

import (
	"fmt"
	"iter"
)

func evenIter(s []int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for _, num := range s {
			if num%2 == 0 && !yield(num) {
				return
			}
		}
	}
}

func main() {
	l := []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	for i := range evenIter(l) {
		fmt.Println(i)
	}
}
