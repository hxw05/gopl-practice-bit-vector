package main

import (
	"testing"
)

func TestIntSet_BasicOperations(t *testing.T) {
	s := new(IntSet)

	// 测试 Add 和 Has
	s.Add(1)
	s.Add(64) // 测试跨 word 边界 (假设 size=64, 64 是下一个 word 的第 0 位)
	s.Add(100)

	if !s.Has(1) {
		t.Error("Expected Has(1) to be true")
	}
	if !s.Has(64) {
		t.Error("Expected Has(64) to be true")
	}
	if s.Has(2) {
		t.Error("Expected Has(2) to be false")
	}

	// 测试 Len
	if s.Len() != 3 {
		t.Errorf("Expected Len() to be 3, got %d", s.Len())
	}

	// 测试 Remove
	s.Remove(1)
	if s.Has(1) {
		t.Error("Expected Has(1) to be false after Remove")
	}
	if s.Len() != 2 {
		t.Errorf("Expected Len() to be 2 after Remove, got %d", s.Len())
	}

	// 测试 Remove 不存在的元素 (不应 panic)
	s.Remove(999)
}

func TestIntSet_ClearAndCopy(t *testing.T) {
	s := new(IntSet)
	s.AddAll(1, 2, 3)

	// 测试 Copy
	cpy := s.Copy()
	if !cpy.Has(1) {
		t.Error("Copy failed: missing element 1")
	}

	// 修改原集合，确保副本不受影响
	s.Clear()
	if !cpy.Has(1) || !cpy.Has(2) || !cpy.Has(3) {
		t.Error("Copy should be independent of original after Clear")
	}

	// 测试 Clear
	s.Add(5)
	s.Clear()
	if s.Len() != 0 {
		t.Error("Clear failed: length should be 0")
	}
}

func TestIntSet_UnionAndIntersect(t *testing.T) {
	s1 := new(IntSet)
	s1.AddAll(1, 2, 3)

	s2 := new(IntSet)
	s2.AddAll(3, 4, 5)

	// 测试 UnionWith
	s1.UnionWith(s2)
	expectedUnion := "{1 2 3 4 5}"
	if s1.String() != expectedUnion {
		t.Errorf("UnionWith failed. Expected: %s, Got: %s", expectedUnion, s1.String())
	}

	// 重置数据测试 IntersectWith
	s1 = new(IntSet)
	s1.AddAll(1, 2, 3)
	s2 = new(IntSet)
	s2.AddAll(3, 4, 5)

	s1.IntersectWith(s2)
	expectedIntersect := "{3}"
	if s1.String() != expectedIntersect {
		t.Errorf("IntersectWith failed. Expected: %s, Got: %s", expectedIntersect, s1.String())
	}
}

func TestIntSet_DifferenceAndSymmetric(t *testing.T) {
	// 差集: A - B (在 A 中但不在 B 中)
	s1 := new(IntSet)
	s1.AddAll(1, 2, 3, 4)
	s2 := new(IntSet)
	s2.AddAll(3, 4, 5, 6)

	s1.DifferenceWith(s2)
	if s1.String() != "{1 2}" {
		t.Errorf("DifferenceWith failed. Expected {1 2}, Got: %s", s1.String())
	}

	// 对称差: (A U B) - (A n B) 或者 (A-B) U (B-A)
	s3 := new(IntSet)
	s3.AddAll(1, 2, 3)
	s4 := new(IntSet)
	s4.AddAll(3, 4, 5)

	symDiff := s3.SymmetricDifference(s4)
	// 期望结果: 1, 2 (来自 s3), 4, 5 (来自 s4)
	// 顺序取决于实现，这里 String() 会按顺序输出
	expectedSym := "{1 2 4 5}"
	if symDiff.String() != expectedSym {
		t.Errorf("SymmetricDifference failed. Expected %s, Got: %s", expectedSym, symDiff.String())
	}
}

func TestIntSet_Elems(t *testing.T) {
	s := new(IntSet)
	s.AddAll(10, 20, 30)

	elems := s.Elems()
	if len(elems) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(elems))
	}

	// 简单检查元素是否存在
	expectedMap := map[int]bool{10: true, 20: true, 30: true}
	for _, e := range elems {
		if !expectedMap[e] {
			t.Errorf("Unexpected element %d in Elems()", e)
		}
	}
}

// 边界测试：空集合操作
func TestIntSet_Empty(t *testing.T) {
	s1 := new(IntSet)
	s2 := new(IntSet)

	// 空集合并集
	s1.UnionWith(s2)
	if s1.Len() != 0 {
		t.Error("Union of empty sets should be empty")
	}

	// 空集交集
	s1.Add(1)
	s2.IntersectWith(s1) // s2 is empty
	if s2.Len() != 0 {
		t.Error("Intersection with empty set should be empty")
	}
}
