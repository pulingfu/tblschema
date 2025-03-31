package dvaplugin

import (
	"sort"
	"sync"
)

// Slicer  切片处理
type Slicer[T any] struct {
	data []T

	lock sync.Mutex
}

// NewSlicer 从切片创建切片对象
func NewSlicer[T any](input []T) *Slicer[T] {
	return &Slicer[T]{
		data: input,
		lock: sync.Mutex{},
	}
}

func (s *Slicer[T]) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.data)
}

func (s *Slicer[T]) Data() []T {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.data
}

// func (s *Slicer[T]) ResSet() *Slicer[T] {
// 	if s.resSet == nil {
// 		s.resSet = NewSlicer(make([]T, 0))
// 	}
// 	return s.resSet
// }

// InSilce  判断元素是否在数组中的方法
func (s *Slicer[T]) InSilce(item T, equal func(a, b T) bool) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, v := range s.data {
		if equal(v, item) {
			return true
		}
	}
	return false
}

// Contains 判断是否有元素符合条件
func (s *Slicer[T]) Contains(equal func(b T) bool) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, v := range s.data {
		if equal(v) {
			return true
		}
	}
	return false
}

// Divide 分割数组
func (s *Slicer[T]) Divide(_f_div func(T) bool) (hit *Slicer[T], miss *Slicer[T]) {
	s.lock.Lock()
	defer s.lock.Unlock()

	hits := make([]T, 0, len(s.data))
	misses := make([]T, 0, len(s.data))
	for _, v := range s.data {
		if _f_div(v) {
			hits = append(hits, v)
		} else {
			misses = append(misses, v)
		}
	}
	return NewSlicer(hits), NewSlicer(misses)
}

// Take 根据条件返回第一个符合条件的元素
func (s *Slicer[T]) Take(_f func(T) bool) T {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, v := range s.data {
		if _f(v) {
			return v
		}
	}
	return *new(T)
}

// Find 查找数组中符合条件的元素
func (s *Slicer[T]) Find(_f func(T) bool) *Slicer[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	var matches = make([]T, 0, len(s.data))
	for _, v := range s.data {
		if _f(v) {
			matches = append(matches, v)
		}
	}
	s.data = matches
	return s
}

func (s *Slicer[T]) _pop_idx(idx int) T {
	if idx < 0 || idx >= len(s.data) {
		return *new(T)
	}
	x := s.data[idx]
	if len(s.data) == 1 {
		s.data = make([]T, 0)
	} else {
		var new = make([]T, len(s.data)-1)
		copy(new[:idx], s.data[:idx])
		if idx < len(s.data)-1 {
			copy(new[idx:], s.data[idx+1:])
		}
		s.data = new
	}
	return x
}

// PopIdx 取出指定位置的元素并且在切片中删除
func (s *Slicer[T]) PopIdx(idx int) T {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s._pop_idx(idx)
}

// PopHead 取出第一个元素并且在切片中删除
func (s *Slicer[T]) PopHead() T {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s._pop_idx(0)
}

// PopHead 取出最后一个元素并且在切片中删除
func (s *Slicer[T]) PopTail() T {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s._pop_idx(len(s.data) - 1)
}

// Append 向切片尾部添加元素
func (s *Slicer[T]) Append(item ...T) *Slicer[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	var new = make([]T, len(s.data)+len(item))
	copy(new, s.data)
	copy(new[len(s.data):], item)
	s.data = new
	return s
}

// Prepend 向切片头部添加元素
func (s *Slicer[T]) Prepend(item ...T) *Slicer[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	var new = make([]T, len(s.data)+len(item))
	copy(new[len(item):], s.data)
	copy(new, item)
	s.data = new
	return s
}

// Delete 删除切片中的指定索引元素
func (s *Slicer[T]) RemoveByIdx(idx int) *Slicer[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	s._pop_idx(idx)
	return s
}

// InsertIdx 在指定索引位置插入元素
// |idx<=0 从头部插入
// |idx>=len(s.data) 从尾部插入
// |idx=1 表示从索引1之前插入 以此类推
func (s *Slicer[T]) InsertIdx(idx int, item ...T) *Slicer[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	var new = make([]T, 0, len(s.data)+len(item))
	if idx < 0 {
		new = append(new, item...)
		new = append(new, s.data...)
	} else if idx >= len(s.data) {
		new = append(new, s.data...)
		new = append(new, item...)
	} else {
		new = append(new, s.data[:idx]...)
		new = append(new, item...)
		if idx < len(s.data) {
			new = append(new, s.data[idx:]...)
		}
	}
	s.data = new
	return s
}

// Page 分页
// |offset 偏移量，跳过前offset个
// |limit 每页数量
func (s *Slicer[T]) Page(offset, limit int) *Slicer[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	if offset >= len(s.data) || limit <= 0 {
		s.data = make([]T, 0)
		return s
	}
	if offset > 0 {
		s.data = s.data[offset:]
	}
	if limit >= len(s.data) {
		return s
	}
	s.data = s.data[:limit]
	return s
}

// Sort 排序
// |a>b 降序
// |a<b 升序
func (s *Slicer[T]) Sort(_f func(a, b T) bool) *Slicer[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	sort.Slice(s.data, func(i, j int) bool {
		return _f(s.data[i], s.data[j])
	})
	return s
}

// func (s *Slicer[T]) Map(transform func(T) any) any {
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	result := make([]any, 0, len(s.data))
// 	for i, v := range s.data {
// 		result[i] = transform(v)
// 	}
// 	return result
// }

// Map 方法
func Map[T, R any](input []T, transform func(T) R) []R {
	result := make([]R, 0, len(input))
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

// Reduce 方法
func Reduce[T, R any](input []T, transform func(R, T) R) R {
	result := *new(R)
	for _, v := range input {
		result = transform(result, v)
	}
	return result
}
