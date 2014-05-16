package main

import (
	"testing"
	"runtime"
)

var a int

func BenchmarkStatic16kbAssignByte(b *testing.B) {
	runtime.GC()
	arr := make([]byte, 128 * 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range arr {
			arr[j] = 0
		}
	}
}

func BenchmarkStatic16kbAssignUint64(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 2048)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range arr {
			arr[j] = 0
		}
	}
}

func BenchmarkStatic16kbReadRangeByte(b *testing.B) {
	runtime.GC()
	arr := make([]byte, 128 * 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range arr {
			b := arr[j]
			if b > 0 {
				a ++
			}
		}
	}
}

func BenchmarkStatic16kbReadIndexByte2Uint64(b *testing.B) {
	runtime.GC()
	arr := make([]byte, 128 * 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 128 * 128; j ++ {
			b := uint64(arr[j])
			if b > 0 {
				a = int(b)
			}
		}
	}
}

func BenchmarkStatic16kbReadIndexByte(b *testing.B) {
	runtime.GC()
	arr := make([]byte, 128 * 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 128 * 128; j ++ {
			b := arr[j]
			if b > 0 {
				a ++
			}
		}
	}
}

func BenchmarkStatic16kbReadIndexUint32(b *testing.B) {
	runtime.GC()
	arr := make([]uint32, 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 4096; j ++ {
			b := arr[j]
			if b > 0 {
				a ++
			}
		}
	}
}

func BenchmarkStatic16kbReadIndexUint64(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 2048)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 2048; j ++ {
			b := arr[j]
			if b > 0 {
				a ++
			}
		}
	}
}

func BenchmarkStatic16kbReadIndexByte8x8(b *testing.B) {
	runtime.GC()
	arr := make([]byte, 128 * 128)
	b.ResetTimer()
	for i := 0 ; i < b.N; i ++ {
		for j := 128 + 1; j < 128 * 128 - 128 - 2; j ++ {
			b := arr[j - 128 - 1] + arr[j - 128] + arr[j - 128 + 1]
			b += arr[j - 1] + arr[j] + arr[j + 1]
			b += arr[j + 128 - 1] + arr[j + 128] + arr[j + 128 + 1]
			if b > 0 {
				a ++
			}
		}
	}
}

func BenchmarkStatic16kbReadIndexUInt643x1(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 2048)
	b.ResetTimer()
	for i := 0 ; i < b.N; i ++ {
		for j := 16; j < 16 * 16 - 16 - 1; j ++ {
			b := arr[j - 16] + arr[j] + arr[j + 16]
			b1:= b & 0xff + (b >> 8) & 0xff + (b >> 16) & 0xff
			if b1 > 0 {
				a ++
			}
		}
	}
}

func BenchmarkStatic16kbReadIndexUInt643x1_fcall(b *testing.B) {
	runtime.GC()
	arr := make([]uint64, 2048)
	b.ResetTimer()
	for i := 0 ; i < b.N; i ++ {
		for j := 16; j < 16 * 16 - 16 - 1; j ++ {
			b := arr[j - 16] + arr[j] + arr[j + 16]
			b1:= b & 0xff + (b >> 8) & 0xff + (b >> 16) & 0xff
			fcallA(b1)
		}
	}
}

func fcallA(t uint64) {
	if t > 0 {
		a ++
	}
}

func BenchmarkStatic16kbAlloc(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	for i := 0 ; i < b.N; i ++ {
		arr := make([]byte, 128 * 128)
		if (arr[0] > 0) {
			a++
		}
	}
}
