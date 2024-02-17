package main

import "fmt"

func main() {
	slice := []int{8, 91, 7, 5, -8, 31, 88, 99, 12, 34, 14, 2, -8}
	fmt.Println("Original Slice:")
	PrintSlice(slice)

	fmt.Println("\nSorted Slice:")
	SortSlice(slice)
	PrintSlice(slice)

	fmt.Println("\nIncremented Odd Positions:")
	IncrementOdd(slice)
	PrintSlice(slice)

	fmt.Println("\nReversed Slice:")
	ReverseSlice(slice)
	PrintSlice(slice)

	fmt.Println("\nSlice after appending functions:")
	functions := appendFunc(SortSlice, IncrementOdd, ReverseSlice)
	functions(slice)
	PrintSlice(slice)
}

func SortSlice(slice []int) {
	quicksort(slice, 0, len(slice)-1)
}

func quicksort(slice []int, low, high int) {
	if low < high {
		pivot := partition(slice, low, high)
		quicksort(slice, low, pivot-1)
		quicksort(slice, pivot+1, high)
	}
}

func partition(slice []int, low, high int) int {
	pivot := slice[high]
	i := low - 1
	for j := low; j < high; j++ {
		if slice[j] < pivot {
			i++
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
	slice[i+1], slice[high] = slice[high], slice[i+1]
	return i + 1
}

func IncrementOdd(slice []int) {
	for i := 1; i < len(slice); i += 2 {
		slice[i]++
	}
}

func ReverseSlice(slice []int) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func PrintSlice(slice []int) {
	for _, val := range slice {
		fmt.Print(val, " ")
	}
	fmt.Println()
}

func appendFunc(dst func([]int), src ...func([]int)) func([]int) {
	return func(slice []int) {
		dst(slice)
		for _, fn := range src {
			fn(slice)
		}
	}
}
