// #include "${GOROOT}/cmd/ld/textflag.h"
// Don't insert stack check preamble.

#define NOSPLIT 4

// 45nm Next Generation Intel® Core™2 Processor Family (Penryn) and Intel® Streaming SIMD Extensions 4 (Intel® SSE4)
// https://software.intel.com/en-us/articles/45nm-next-generation-intel-coret-2-processor-family-penryn-and-intel-streaming-simd-extensions-4-intel-sse4/
// POPCNT r64, r/m64 - F3 REX.W 0F B8
#define POPCNT_AX_AX BYTE $0xf3; BYTE $0x48; BYTE $0x0f; BYTE $0xb8; BYTE $0xc0
#define POPCNT_8RSP_AX BYTE $0xf3; BYTE $0x48; BYTE $0x0f; BYTE $0xb8; BYTE $0x44; BYTE $0x24; BYTE $0x08

// func PopCnt(n uint64) uint64
// void PopCnt(unsigned long long *ret, unsigned long long n);
TEXT ·PopCnt(SB),NOSPLIT,$0-8
    POPCNT_8RSP_AX
    MOVQ AX, ret+8(FP)
    RET
