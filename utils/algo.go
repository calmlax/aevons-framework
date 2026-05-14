// Package utils 提供常见排序与查找算法的泛型实现。
// 注意：生产环境排序推荐使用标准库 slices.SortFunc，此处实现供学习与参考。
package utils

import "cmp"

// ------------------- 排序算法 -------------------

// BubbleSort 冒泡排序，时间复杂度 O(n²)，原地排序
func BubbleSort[T cmp.Ordered](slice []T) {
	n := len(slice)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-1-i; j++ {
			if slice[j] > slice[j+1] {
				slice[j], slice[j+1] = slice[j+1], slice[j]
				swapped = true
			}
		}
		if !swapped {
			break // 已有序，提前退出
		}
	}
}

// InsertionSort 插入排序，时间复杂度 O(n²)，小数据集表现好，原地排序
func InsertionSort[T cmp.Ordered](slice []T) {
	for i := 1; i < len(slice); i++ {
		key := slice[i]
		j := i - 1
		for j >= 0 && slice[j] > key {
			slice[j+1] = slice[j]
			j--
		}
		slice[j+1] = key
	}
}

// SelectionSort 选择排序，时间复杂度 O(n²)，原地排序
func SelectionSort[T cmp.Ordered](slice []T) {
	n := len(slice)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if slice[j] < slice[minIdx] {
				minIdx = j
			}
		}
		slice[i], slice[minIdx] = slice[minIdx], slice[i]
	}
}

// MergeSort 归并排序，时间复杂度 O(n log n)，返回新切片
func MergeSort[T cmp.Ordered](slice []T) []T {
	if len(slice) <= 1 {
		return slice
	}
	mid := len(slice) / 2
	left := MergeSort(slice[:mid])
	right := MergeSort(slice[mid:])
	return merge(left, right)
}

func merge[T cmp.Ordered](left, right []T) []T {
	result := make([]T, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// QuickSort 快速排序，平均时间复杂度 O(n log n)，原地排序
func QuickSort[T cmp.Ordered](slice []T) {
	if len(slice) <= 1 {
		return
	}
	pivot := slice[len(slice)/2]
	left, right := 0, len(slice)-1
	for left <= right {
		for slice[left] < pivot {
			left++
		}
		for slice[right] > pivot {
			right--
		}
		if left <= right {
			slice[left], slice[right] = slice[right], slice[left]
			left++
			right--
		}
	}
	QuickSort(slice[:right+1])
	QuickSort(slice[left:])
}

// HeapSort 堆排序，时间复杂度 O(n log n)，原地排序
func HeapSort[T cmp.Ordered](slice []T) {
	n := len(slice)
	for i := n/2 - 1; i >= 0; i-- {
		heapify(slice, n, i)
	}
	for i := n - 1; i > 0; i-- {
		slice[0], slice[i] = slice[i], slice[0]
		heapify(slice, i, 0)
	}
}

func heapify[T cmp.Ordered](slice []T, n, i int) {
	largest, left, right := i, 2*i+1, 2*i+2
	if left < n && slice[left] > slice[largest] {
		largest = left
	}
	if right < n && slice[right] > slice[largest] {
		largest = right
	}
	if largest != i {
		slice[i], slice[largest] = slice[largest], slice[i]
		heapify(slice, n, largest)
	}
}

// ------------------- 查找算法 -------------------

// BinarySearch 二分查找（要求切片已升序排列），返回下标，未找到返回 -1
func BinarySearch[T cmp.Ordered](slice []T, target T) int {
	lo, hi := 0, len(slice)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		switch {
		case slice[mid] == target:
			return mid
		case slice[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}

// LinearSearch 线性查找，返回第一个匹配的下标，未找到返回 -1
func LinearSearch[T comparable](slice []T, target T) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}
