package main

import (
	"testing"
	"runtime"
)

func TestNumberOfSetBits(t *testing.T) {
	if NumberOfSetBits(0) != 0 {
		t.Errorf("Bits(%d) != %d\n", 0, 0)
	}
	if NumberOfSetBits(3) != 2 {
		t.Errorf("Bits(%d) != %d\n", 3, 2)
	}
	if NumberOfSetBits(^uint64(0)) != 64 {
		t.Errorf("Bits(%d) != %d, but %d\n", ^uint64(0), 64, NumberOfSetBits(^uint64(0)))
	}
}

func TestPopCnt(t *testing.T) {
	if PopCnt(0) != 0 {
		t.Errorf("Bits(%d) != %d\n", 0, 0)
	}
	if PopCnt(3) != 2 {
		t.Errorf("Bits(%d) != %d\n", 3, 2)
	}
	if PopCnt(^uint64(0)) != 64 {
		t.Errorf("Bits(%d) != %d, but %d\n", ^uint64(0), 64, PopCnt(^uint64(0)))
	}
}

func BenchmarkBits512bReadIndexUint64(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 64 /* 64 x 64 */, 64)
	arr_len := len(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < arr_len; j ++ {
			b := arr[j]
			if b > 0 {
				a ++
			}
		}
	}
}

func BenchmarkBits512bReadBits1_Empty(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 64 /* 64 x 64 */, 64)
	arr_len := len(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < arr_len; j ++ {
			b := arr[j]
			if NumberOfSetBits(b) > 0 {
				a ++
			}
		}
	}
}

func BenchmarkBits512bReadBits1_Full(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 64 /* 64 x 64 */, 64)
	arr_len := len(arr)
	for j := 0; j < arr_len; j ++ {
		arr[j] = ^uint64(0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < arr_len; j ++ {
			b := arr[j]
			if NumberOfSetBits(b) == 0 {
				a ++
			}
		}
	}
}

func BenchmarkBits512bReadBits2_Empty(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 64 /* 64 x 64 */, 64)
	arr_len := len(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < arr_len; j ++ {
			b := arr[j]
			if PopCnt(b) > 0 {
				a ++
			}
		}
	}
}

func BenchmarkBits512bReadBits2_Full(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 64 /* 64 x 64 */, 64)
	arr_len := len(arr)
	for j := 0; j < arr_len; j ++ {
		arr[j] = ^uint64(0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < arr_len; j ++ {
			b := arr[j]
			if PopCnt(b) == 0 {
				a ++
			}
		}
	}
}
