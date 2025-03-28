package dvaplugin

// Slicer  切片处理
type Slicer[T any] struct {
	Input []T
}

// NewSlicer 从切片创建切片对象
func NewSlicer[T any](input []T) *Slicer[T] {
	return &Slicer[T]{
		Input: input,
	}
}

// InSilce  判断元素是否在数组中的方法
func (s *Slicer[T]) InSilce(item T, equal func(a, b T) bool) bool {
	for _, v := range s.Input {
		if equal(v, item) {
			return true
		}
	}
	return false
}

// Contains 判断是否有元素符合条件
func (s *Slicer[T]) Contains(equal func(b T) bool) bool {
	for _, v := range s.Input {
		if equal(v) {
			return true
		}
	}
	return false
}

// Divide 分割数组
func (s *Slicer[T]) Divide(_f_div func(T) bool) (hits []T, misses []T) {
	hits = make([]T, 0, len(s.Input))
	misses = make([]T, 0, len(s.Input))
	for _, v := range s.Input {
		if _f_div(v) {
			hits = append(hits, v)
		} else {
			misses = append(misses, v)
		}
	}
	return hits, misses
}

// Take 根据条件返回第一个符合条件的元素
func (s *Slicer[T]) Take(_f func(T) bool) T {
	for _, v := range s.Input {
		if _f(v) {
			return v
		}
	}
	return *new(T)
}

// Find 查找数组中符合条件的元素
func (s *Slicer[T]) Find(_f func(T) bool) []T {
	var matches = make([]T, 0, len(s.Input))
	for _, v := range s.Input {
		if _f(v) {
			matches = append(matches, v)
		}
	}
	return matches
}

// Map 方法
func Map[T, R any](input []T, transform func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = transform(v)
	}
	return result
}

// DuplicateMerge 合并重复元素
// 使用dupf判断元素重复
// 再使用transform合并重复的元素，生成新的元素
// 新旧元素可以是不同类型
func DuplicateMerge[T, R any](input []T,
	dupf func(T) any,
	transform func(T, R) R,
) []R {
	seen := make(map[any]R)
	order := make([]any, 0, len(input))
	for _, v := range input {
		key := dupf(v)
		if _, ok := seen[key]; !ok {
			seen[key] = *new(R)
			order = append(order, key)
		}
		seen[key] = transform(v, seen[key])
	}

	result := make([]R, 0, len(seen))
	for _, v := range order {
		result = append(result, seen[v])
	}
	return result
}
