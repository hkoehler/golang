package main

/***************************************************
 * Concurrent QuickSort implemention in Go
 * Copyright (C) 2016, Heiko Koehler
 ***************************************************/

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// dump integer slice
func dump(nums []int) {
	for _, num := range nums {
		fmt.Printf("%d ", num)
	}
	fmt.Printf("\n")
}

// 3-way partitioning
// Return lt and gt element
// Choose pivot at random
func partition(nums []int) (int, int) {
	p := nums[rand.Int()%len(nums)]
	i := 0              // incremented for nums[i] <= p
	lt := 0             // incremented for nums[i] < p
	gt := len(nums) - 1 // decremented for nums[i] > p

	for i <= gt {
		if nums[i] < p {
			nums[i], nums[lt] = nums[lt], nums[i]
			i++
			lt++
		} else if nums[i] > p {
			nums[i], nums[gt] = nums[gt], nums[i]
			gt--
		} else {
			i++
		}
	}
	return lt, gt
}

// QuickSort
func qsort(nums []int, s int, comp chan bool) {
	if len(nums) > 1 {
		lt, gt := partition(nums)

		// spawn two concurrent qsorts if slice bigger than s
		if len(nums) >= s {
			comp1 := make(chan bool)
			comp2 := make(chan bool)

			go qsort(nums[:lt], s, comp1)
			go qsort(nums[gt:], s, comp2)

			count := 0
			for count != 2 {
				select {
				case <-comp1:
					count++
				case <-comp2:
					count++
				}
			}
		} else {
			qsort(nums[:lt], s, nil)
			qsort(nums[gt:], s, nil)
		}
	}

	if comp != nil {
		comp <- true
	}
}

// generate array populated with random numbers
// max: Maximum int value
func generate(n int, max int) []int {
	nums := make([]int, n)

	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Int() % (max + 1)
	}
	return nums
}

// verify data is sorted
func verify(nums []int) bool {
	for i := 1; i < len(nums); i++ {
		if nums[i-1] > nums[i] {
			return false
		}
	}
	return true
}

func main() {
	var n int
	var m int
	var s int
	var d bool
	var v bool

	flag.IntVar(&n, "n", 10, "Number of elements to sort")
	flag.IntVar(&m, "m", 100, "Max value of element")
	flag.IntVar(&s, "s", 10000, "Min slice size for spawning go routine")
	flag.BoolVar(&d, "d", false, "Dump sorted array")
	flag.BoolVar(&v, "v", false, "Verify array is sorted")
	flag.Parse()

	nums := generate(n, m)
	now := time.Now()
	qsort(nums, s, nil)
	delta := time.Now().UnixNano() - now.UnixNano()
	fmt.Printf("Sorted %d elements in %d us\n", n, delta/1000)
	if v && !verify(nums) {
		fmt.Println("Array not sorted!")
	}
	if d {
		dump(nums)
	}
}
