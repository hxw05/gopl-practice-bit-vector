package main

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint
}

const size = 32 << (^uint(0) >> 63)

// Has 检验集合中是否有数字x
func (s *IntSet) Has(x int) bool {
	// why uint(...)?
	word, bit := x/size, uint(x%size)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add 向集合中添加数字x
func (s *IntSet) Add(x int) {
	word, bit := x/size, uint(x%size)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith 将集合与t做并集运算，将改变原集合的值
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String 以类似于{1 2 3}的形式返回集合
func (s IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := range size {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", size*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len 返回集合的长度
func (s *IntSet) Len() (len int) {
	for _, word := range s.words {
		for i := range size {
			if word&(1<<uint(i)) != 0 {
				len++
			}
		}
	}
	return len
}

// Remove 从集合中移除元素x
func (s *IntSet) Remove(x int) {
	word, bit := x/size, uint(x%size)
	if word >= len(s.words) {
		return
	}
	s.words[word] &^= (1 << bit)
}

// Clear 清空集合
func (s *IntSet) Clear() {
	s.words = []uint{}
}

// Copy 返回集合的一个副本
func (s *IntSet) Copy() *IntSet {
    wordsCpy := make([]uint, len(s.words))
    copy(wordsCpy, s.words)
    return &IntSet{words: wordsCpy}
}

func (s *IntSet) AddAll(xs ...int) {
	for _, x := range xs {
		s.Add(x)
	}
}

// IntersectWith 将当前集合与t做交集运算
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
		} else {
			break
		}
	}
}

// DifferenceWith 将当前集合与t做差集运算
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &^= tword
		} else {
			break
		}
	}
}

// SymmetricDifference 返回s和t的对称差
func (s *IntSet) SymmetricDifference(t *IntSet) *IntSet {
	r := &IntSet{}
	minLen := min(len(s.words), len(t.words))

	for i := range minLen {
		r.words = append(r.words, t.words[i]^s.words[i])
	}

	for i := minLen; i < len(s.words); i++ {
		r.words = append(r.words, s.words[i])
	}

	for i := minLen; i < len(t.words); i++ {
		r.words = append(r.words, t.words[i])
	}

	return r
}

func (s *IntSet) Elems() (elems []int) {
	for i, word := range s.words {
		for j := range size {
			if word&(1<<uint(j)) != 0 {
				elems = append(elems, i*size+j)
			}
		}
	}
	return elems
}

func main() {
	fmt.Printf("当前平台字长: %d 位\n\n", size)

	// --- 1. 基础操作演示 ---
	fmt.Println("=== 1. 基础操作演示 ===")
	var s1 IntSet
	fmt.Printf("新建空集合 s1: %s (长度: %d)\n", s1, s1.Len())

	// 添加单个元素
	s1.Add(1)
	s1.Add(5)
	s1.Add(10)
	fmt.Printf("添加 1, 5, 10 后 s1: %s (长度: %d)\n", s1, s1.Len())

	// 批量添加元素
	s1.AddAll(15, 20, 25)
	fmt.Printf("批量添加 15, 20, 25 后 s1: %s (长度: %d)\n", s1, s1.Len())

	// 检查元素存在性
	fmt.Printf("Has(5): %t, Has(6): %t\n", s1.Has(5), s1.Has(6))

	// 移除元素
	s1.Remove(10)
	fmt.Printf("移除 10 后 s1: %s (长度: %d)\n", s1, s1.Len())

	// 移除不存在的元素（应无错误）
	s1.Remove(100)
	fmt.Printf("尝试移除不存在的 100 后 s1: %s\n", s1)

	// 清空集合
	s1.Clear()
	fmt.Printf("清空后 s1: %s (长度: %d)\n\n", s1, s1.Len())

	// --- 2. 集合运算演示 ---
	fmt.Println("=== 2. 集合运算演示 ===")
	var sA, sB IntSet
	sA.AddAll(1, 2, 3, 4, 5)
	sB.AddAll(4, 5, 6, 7, 8)
	fmt.Printf("集合 A: %s\n", sA)
	fmt.Printf("集合 B: %s\n", sB)

	// 并集 (Union) - 修改原集合
	var sUnion IntSet
	sUnion.AddAll(1, 2, 3) // 先加点数据
	fmt.Printf("原集合 (1,2,3) 与 B 做并集...\n")
	sUnion.UnionWith(&sB)
	fmt.Printf("并集结果 (UnionWith): %s\n", sUnion)

	// 交集 (Intersect) - 修改原集合
	sA.IntersectWith(&sB)
	fmt.Printf("A 与 B 做交集 (IntersectWith): %s\n", sA) // A 变为 {4 5}

	// 差集 (Difference) - 修改原集合
	var sDiff IntSet
	sDiff.AddAll(1, 2, 3, 4, 5)
	sDiff.DifferenceWith(&sB) // 从 {1,2,3,4,5} 中减去 {4,5,6,7,8}
	fmt.Printf("A 与 B 做差集 (DifferenceWith): %s\n", sDiff) // 结果应为 {1 2 3}

	// 对称差 (SymmetricDifference) - 返回新集合
	var sSymA, sSymB IntSet
	sSymA.AddAll(1, 2, 3)
	sSymB.AddAll(3, 4, 5)
	sSymDiff := sSymA.SymmetricDifference(&sSymB)
	fmt.Printf("{1,2,3} 与 {3,4,5} 的对称差: %s\n\n", sSymDiff) // 结果应为 {1 2 4 5}

	// --- 3. 副本与元素提取 ---
	fmt.Println("=== 3. 副本与元素提取 ===")
	var sOriginal IntSet
	sOriginal.AddAll(10, 20, 30)
	fmt.Printf("原始集合: %s\n", sOriginal)

	// 复制集合
	sCopy := sOriginal.Copy()
	fmt.Printf("复制集合: %s\n", sCopy)

	// 修改副本，验证深拷贝
	sCopy.Add(40)
	sCopy.Remove(10)
	fmt.Printf("修改副本后 (加40, 减10): %s\n", sCopy)
	fmt.Printf("原始集合保持不变: %s\n", sOriginal)

	// 提取所有元素
	elems := sOriginal.Elems()
	fmt.Printf("原始集合的元素切片: %v\n\n", elems)

	// --- 4. 边界与压力测试 ---
	fmt.Println("=== 4. 边界与压力测试 ===")
	var sEdge IntSet
	// 添加大数值，测试动态扩容
	sEdge.Add(99999999)
	sEdge.Add(1)
	sEdge.Add(64)   // 字边界
	sEdge.Add(63)   // 字边界前一位
	sEdge.Add(128)  // 跨两个字
	fmt.Printf("包含大数值的集合: %s\n", sEdge)
	fmt.Printf("集合长度校验: %d\n\n", sEdge.Len())
}