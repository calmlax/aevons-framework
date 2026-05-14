package utils

// SliceContains 判断切片中是否包含某元素
func SliceContains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// SliceIndex 返回元素第一次出现的下标，不存在返回 -1
func SliceIndex[T comparable](slice []T, item T) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// SliceUnique 去重，保持原顺序
func SliceUnique[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// SliceFilter 过滤切片，保留满足条件的元素
func SliceFilter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// SliceMap 对切片每个元素做转换
func SliceMap[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// SliceReduce 对切片做聚合
func SliceReduce[T any, R any](slice []T, init R, fn func(R, T) R) R {
	acc := init
	for _, v := range slice {
		acc = fn(acc, v)
	}
	return acc
}

// SliceChunk 将切片按指定大小分块
func SliceChunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	var chunks [][]T
	for size < len(slice) {
		slice, chunks = slice[size:], append(chunks, slice[:size])
	}
	return append(chunks, slice)
}

// SliceReverse 反转切片，返回新切片
func SliceReverse[T any](slice []T) []T {
	n := len(slice)
	result := make([]T, n)
	for i, v := range slice {
		result[n-1-i] = v
	}
	return result
}

// SliceToMap 将切片转换为 map，key 由 fn 提取
func SliceToMap[T any, K comparable](slice []T, fn func(T) K) map[K]T {
	result := make(map[K]T, len(slice))
	for _, v := range slice {
		result[fn(v)] = v
	}
	return result
}
