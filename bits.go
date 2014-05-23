package main

// http://stackoverflow.com/questions/109023/how-to-count-the-number-of-set-bits-in-a-32-bit-integer
// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
// http://gurmeet.net/puzzles/fast-bit-counting-routines/
// http://aggregate.ee.engr.uky.edu/MAGIC/#Population%20Count%20%28Ones%20Count%29
// http://en.wikipedia.org/wiki/Hamming_weight
func NumberOfSetBits(n uint64) uint64 {
	tmp := n - ((n >> 1) & 0x7777777777777777) - ((n >> 2) & 0x3333333333333333) - ((n >> 3) & 0x1111111111111111)
	return ((tmp + (tmp >> 4) ) & 0x0F0F0F0F0F0F0F0F) % 255
}

// __builtin_popcount
//go:noescape
func PopCnt(n uint64) uint64
