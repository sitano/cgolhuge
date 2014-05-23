// #include "${GOROOT}/cmd/ld/textflag.h"
// Don't insert stack check preamble.
// #define NOSPLIT	4

#define POPCNT_AX_AX BYTE $0xf3; BYTE $0x48; BYTE $0x0f; BYTE $0xb8; BYTE $0xc0

// func PopCnt(n uint64) uint64
// void PopCnt(unsigned long long *ret, unsigned long long n);
// TEXT ·PopCnt(SB),NOSPLIT,$0
TEXT ·PopCnt(SB),4,$0-8
    MOVQ n+0(FP), AX
//  POPCNT AX, AX
    BYTE $0xf3; BYTE $0x48; BYTE $0x0f; BYTE $0xb8; BYTE $0xc0
    MOVQ AX, ret+8(FP)
    RET
