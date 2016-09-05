package main

/***************************************************
 * QuickSort implemention in Go
 ***************************************************/

import (
	"math/rand"
	"flag"
	"fmt"
)

// dump integer slice
func dump(nums []int) {
	for _, num := range nums {
		fmt.Printf("%d ", num)
	}
	fmt.Printf("\n")
}

// Hoare partitioning
// Choose pivot at random
func partition(nums []int) int {
	p := nums[rand.Int() % len(nums)]
	i := 0
	j := len(nums) - 1

	for {
		for i < len(nums) && nums[i] < p {
			i++
		}
		for j >= 0 && nums[j] > p {
			j--
		}
		if i >= j {
			break
		}
		nums[i], nums[j] = nums[j], nums[i]
		i++
		j--
	}
	return j
}

// QuickSort
func qsort(nums []int, comp chan bool) {
	if len(nums) > 0 {
	    comp1 := make(chan bool)
	    comp2 := make(chan bool)
	    
	    // spawn two concurrent qsorts
		p := partition(nums)
		go qsort(nums[:p], comp1)
		go qsort(nums[p+1:], comp2)
		
		// wait for both qsorts to finish
		count := 0
		for count != 2 {
	        select {
	        case <- comp1:
	            count++
	        case <- comp2:
	            count++
	        }
	    }
	}
	comp <- true
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

func main() {
	var n int
	var m int
	var d bool

	flag.IntVar(&n, "n", 10, "Number of elements to sort")
	flag.IntVar(&m, "m", 100, "Max value of element")
	flag.BoolVar(&d, "d", false, "Dump sorted array")
	flag.Parse()
	
	nums := generate(n, m)
	comp := make(chan bool)
	go qsort(nums, comp)
	<- comp
	if d {
    	dump(nums)
    }
}

