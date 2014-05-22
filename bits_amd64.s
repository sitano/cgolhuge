// #include "${GOROOT}/cmd/ld/textflag.h"
// Don't insert stack check preamble.
// #define NOSPLIT	4
// func PopCnt(n uint64) uint64
// TEXT ·PopCnt(SB),NOSPLIT,$0
// 0x0000000000400a00 <+0>:	popcnt %edi,%eax
// 0x0000000000400a04 <+4>:	retq
TEXT ·PopCnt(SB),4,$8
//          POPCNT n+0(FP), X0
//          MOVSD X0, ret+8(FP)
          RET
