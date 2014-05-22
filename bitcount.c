/* ==========================================================================
   Bit Counting routines

   Author: Gurmeet Singh Manku    (manku@cs.stanford.edu)
   Date:   27 Aug 2002
   ========================================================================== */

/* Ivan:
// gcc -O4 -mpopcnt -mtune=core2 -march=core2 -fforce-addr -funroll-loops -frerun-cse-after-loop -frerun-loop-opt -malign-functions=4
// CPU0: Intel(R) Core(TM) i7-4770K CPU @ 3.50GHz (fam: 06, model: 3c, stepping: 03)
All BitCounts seem okay!  Starting speed trials
     iterated: 1000 million counts in 12540 ms for       79745 cnts/ms
       sparse: 1000 million counts in  6487 ms for      154154 cnts/ms
        dense: 1000 million counts in  7264 ms for      137665 cnts/ms
  precomputed: 1000 million counts in  1699 ms for      588582 cnts/ms
precomputed16: 1000 million counts in  1439 ms for      694927 cnts/ms
     parallel: 1000 million counts in  2017 ms for      495786 cnts/ms
        nifty: 1000 million counts in  2342 ms for      426985 cnts/ms
     nuonifty: 1000 million counts in  2108 ms for      474383 cnts/ms
      seander: 1000 million counts in  1700 ms for      588235 cnts/ms
          mit: 1000 million counts in  2508 ms for      398724 cnts/ms
       POPCNT: 1000 million counts in  1175 ms for      851064 cnts/ms
   */


#include <stdlib.h>
#include <stdio.h>
#include <limits.h>
#include <sys/time.h>

double kbmin = 1.0e10, kbmax = 0.0, kbsum = 0.0, etsum = 0.0;
static timevalfix(), timevalsub();

/* Iterated bitcount iterates over each bit. The while condition sometimes helps
   terminates the loop earlier */
int iterated_bitcount (unsigned int n)
{
    int count=0;
    while (n)
    {
        count += n & 0x1u ;
        n >>= 1 ;
    }
    return count ;
}


/* Sparse Ones runs proportional to the number of ones in n.
   The line   n &= (n-1)   simply sets the last 1 bit in n to zero. */
int sparse_ones_bitcount (unsigned int n)
{
    int count=0 ;
    while (n)
    {
        count++ ;
        n &= (n - 1) ;
    }
    return count ;
}


/* Dense Ones runs proportional to the number of zeros in n.
   It first toggles all bits in n, then diminishes count repeatedly */
int dense_ones_bitcount (unsigned int n)
{
    int count = 8 * sizeof(int) ;
    n ^= (unsigned int) -1 ;
    while (n)
    {
        count-- ;
        n &= (n - 1) ;
    }
    return count ;
}


/* Precomputed bitcount uses a precomputed array that stores the number of ones
   in each char. */
static int bits_in_char [256] ;

void compute_bits_in_char (void)
{
    unsigned int i ;
    for (i = 0; i < 256; i++)
        bits_in_char [i] = iterated_bitcount (i) ;
    return ;
}

int precomputed_bitcount (unsigned int n)
{
    // works only for 32-bit ints

    return bits_in_char [n         & 0xffu]
        +  bits_in_char [(n >>  8) & 0xffu]
        +  bits_in_char [(n >> 16) & 0xffu]
        +  bits_in_char [(n >> 24) & 0xffu] ;
}


/* Here is another version of precomputed bitcount that uses a precomputed array
   that stores the number of ones in each short. */

static char bits_in_16bits [0x1u << 16] ;

void compute_bits_in_16bits (void)
{
    unsigned int i ;
    for (i = 0; i < (0x1u<<16); i++)
        bits_in_16bits [i] = iterated_bitcount (i) ;
    return ;
}

int precomputed16_bitcount (unsigned int n)
{
    // works only for 32-bit int

    return bits_in_16bits [n         & 0xffffu]
        +  bits_in_16bits [(n >> 16) & 0xffffu] ;
}


/* Parallel   Count   carries   out    bit   counting   in   a   parallel
   fashion.   Consider   n   after    the   first   line   has   finished
   executing. Imagine splitting n into  pairs of bits. Each pair contains
   the <em>number of ones</em> in those two bit positions in the original
   n.  After the second line has finished executing, each nibble contains
   the  <em>number of  ones</em>  in  those four  bits  positions in  the
   original n. Continuing  this for five iterations, the  64 bits contain
   the  number  of ones  among  these  sixty-four  bit positions  in  the
   original n. That is what we wanted to compute. */

#define TWO(c) (0x1u << (c))
#define MASK(c) (((unsigned int)(-1)) / (TWO(TWO(c)) + 1u))
#define COUNT(x,c) ((x) & MASK(c)) + (((x) >> (TWO(c))) & MASK(c))

int parallel_bitcount (unsigned int n)
{
    n = COUNT(n, 0) ;
    n = COUNT(n, 1) ;
    n = COUNT(n, 2) ;
    n = COUNT(n, 3) ;
    n = COUNT(n, 4) ;
    /* n = COUNT(n, 5) ;    for 64-bit integers */
    return n ;
}


/* Nifty  Parallel Count works  the same  way as  Parallel Count  for the
   first three iterations. At the end  of the third line (just before the
   return), each byte of n contains the number of ones in those eight bit
   positions in  the original n. A  little thought then  explains why the
   remainder modulo 255 works. */

#define MASK_01010101 (((unsigned int)(-1))/3)
#define MASK_00110011 (((unsigned int)(-1))/5)
#define MASK_00001111 (((unsigned int)(-1))/17)

int nifty_bitcount (unsigned int n)
{
    n = (n & MASK_01010101) + ((n >> 1) & MASK_01010101) ;
    n = (n & MASK_00110011) + ((n >> 2) & MASK_00110011) ;
    n = (n & MASK_00001111) + ((n >> 4) & MASK_00001111) ;
    return n % 255 ;
}

/* NuoNifty was invented by Nuomnicron and is a minor variation on
   the nifty parallel count to avoid the mod operation */
#define MASK_0101010101010101 (((unsigned int)(-1))/3)
#define MASK_0011001100110011 (((unsigned int)(-1))/5)
#define MASK_0000111100001111 (((unsigned int)(-1))/17)
#define MASK_0000000011111111 (((unsigned int)(-1))/257)
#define MASK_1111111111111111 (((unsigned int)(-1))/65537)
int nuonifty_bitcount (unsigned int n)
{
  n = (n & MASK_0101010101010101) + ((n >> 1) & MASK_0101010101010101) ;
  n = (n & MASK_0011001100110011) + ((n >> 2) & MASK_0011001100110011) ;
  n = (n & MASK_0000111100001111) + ((n >> 4) & MASK_0000111100001111) ;
  n = (n & MASK_0000000011111111) + ((n >> 8) & MASK_0000000011111111) ;
  n = (n & MASK_1111111111111111) + ((n >> 16) & MASK_1111111111111111) ;

  return n;
}

/* Seander parallel count takes only 12 operations, which is the same
   as the lookup-table method, but avoids the memory and potential
   cache misses of a table. It is a hybrid between the purely parallel
   method above and the earlier methods using multiplies (in the
   section on counting bits with 64-bit instructions), though it
   doesn't use 64-bit instructions. The counts of bits set in the
   bytes is done in parallel, and the sum total of the bits set in the
   bytes is computed by multiplying by 0x1010101 and shifting right 24
   bits.  From http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel */
int seander_bitcount(unsigned int n)
{
  n = n - ((n >> 1) & 0x55555555);                        // reuse input as temporary
  n = (n & 0x33333333) + ((n >> 2) & 0x33333333);         // temp
  return(((n + (n >> 4) & 0xF0F0F0F) * 0x1010101) >> 24); // count
}

/* MIT Bitcount

   Consider a 3 bit number as being
        4a+2b+c
   if we shift it right 1 bit, we have
        2a+b
  subtracting this from the original gives
        2a+b+c
  if we shift the original 2 bits right we get
        a
  and so with another subtraction we have
        a+b+c
  which is the number of bits in the original number.

  Suitable masking allows the sums of the octal digits in a 32 bit
  number to appear in each octal digit.  This isn't much help unless
  we can get all of them summed together.  This can be done by modulo
  arithmetic (sum the digits in a number by molulo the base of the
  number minus one) the old "casting out nines" trick they taught in
  school before calculators were invented.  Now, using mod 7 wont help
  us, because our number will very likely have more than 7 bits set.
  So add the octal digits together to get base64 digits, and use
  modulo 63.  (Those of you with 64 bit machines need to add 3 octal
  digits together to get base512 digits, and use mod 511.)

  This is HAKMEM 169, as used in X11 sources.
  Source: MIT AI Lab memo, late 1970's.
*/
int mit_bitcount(unsigned int n)
{
    /* works for 32-bit numbers only */
    register unsigned int tmp;

    tmp = n - ((n >> 1) & 033333333333) - ((n >> 2) & 011111111111);
    return ((tmp + (tmp >> 3)) & 030707070707) % 63;
}

#ifdef __GNUC__
/*
 * This may use the SSE4 POPCNT instruction which will count the number
 * of ones in one CPU instruction.  Obviously you need a recent enough
 * version of gcc and a recent enough CPU to take advantage of this.
 */
int POPCNT_bitcount(unsigned int n)
{
  return(__builtin_popcount(n));
}
#endif

void verify_bitcounts (unsigned int x)
{
    int iterated_ones, sparse_ones, dense_ones ;
    int precomputed_ones, precomputed16_ones ;
    int parallel_ones, nifty_ones, seander_ones ;
    int nuonifty_ones, mit_ones ;
#ifdef __GNUC__
    int POPCNT_ones;
#endif

    iterated_ones      = iterated_bitcount      (x) ;
    sparse_ones        = sparse_ones_bitcount   (x) ;
    dense_ones         = dense_ones_bitcount    (x) ;
    precomputed_ones   = precomputed_bitcount   (x) ;
    precomputed16_ones = precomputed16_bitcount (x) ;
    parallel_ones      = parallel_bitcount      (x) ;
    nifty_ones         = nifty_bitcount         (x) ;
    nuonifty_ones      = nuonifty_bitcount      (x) ;
    seander_ones       = seander_bitcount       (x) ;
    mit_ones           = mit_bitcount           (x) ;
#ifdef __GNUC__
    POPCNT_ones        = POPCNT_bitcount        (x) ;
#endif

    if (iterated_ones != sparse_ones)
    {
        printf ("ERROR: sparse_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

    if (iterated_ones != dense_ones)
    {
        printf ("ERROR: dense_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

    if (iterated_ones != precomputed_ones)
    {
        printf ("ERROR: precomputed_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

    if (iterated_ones != precomputed16_ones)
    {
        printf ("ERROR: precomputed16_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

    if (iterated_ones != parallel_ones)
    {
        printf ("ERROR: parallel_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

    if (iterated_ones != nifty_ones)
    {
        printf ("ERROR: nifty_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

    if (iterated_ones != nuonifty_ones)
    {
        printf ("ERROR: nuonifty_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

    if (iterated_ones != seander_ones)
    {
        printf ("ERROR: nifty_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

    if (mit_ones != nifty_ones)
    {
        printf ("ERROR: mit_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }

#ifdef __GNUC__
    if (POPCNT_ones != nifty_ones)
    {
        printf ("ERROR: mit_bitcount (0x%x) not okay!\n", x) ;
        exit (0) ;
    }
#endif

    return ;
}

bitspeeds()
{
  timeone("iterated", iterated_bitcount);
  timeone("sparse", sparse_ones_bitcount);
  timeone("dense", dense_ones_bitcount);
  timeone("precomputed", precomputed_bitcount);
  timeone("precomputed16", precomputed16_bitcount);
  timeone("parallel", parallel_bitcount);
  timeone("nifty", nifty_bitcount);
  timeone("nuonifty", nuonifty_bitcount);
  timeone("seander", seander_bitcount);
  timeone("mit", mit_bitcount);
#ifdef __GNUC__
  timeone("POPCNT", POPCNT_bitcount);
#else
  printf("Cannot test GCC __builtin_popcount() function\n");
#endif
}

timeone(char *tname, int (*cntr)())
{
  struct timeval tstart, tend;
  const numtestsM = 1000;
  const numtests = numtestsM * 1000000;
  int i, mscount;

  gettimeofday(&tstart, 0);
  for(i=numtests; i; i--)
  {
    mscount = (*cntr)(i);
  }
  gettimeofday(&tend, 0);

  tend.tv_sec -= tstart.tv_sec;
  tend.tv_usec -= tstart.tv_usec;

  if (tend.tv_usec < 0)
  {
    tend.tv_sec--;
    tend.tv_usec += 1000000;
  }
  if (tend.tv_usec >= 1000000)
  {
    tend.tv_sec++;
    tend.tv_usec -= 1000000;
  }
  mscount = tend.tv_sec * 1000 + tend.tv_usec / 1000;

  printf("%13s: %d million counts in %5d ms for %11.0f cnts/ms\n", tname, numtestsM, mscount, ((double)numtests)/((double)mscount));
}



int main (void)
{
    int i ;

    compute_bits_in_char () ;
    compute_bits_in_16bits () ;

    verify_bitcounts (UINT_MAX) ;
    verify_bitcounts (0) ;

    for (i = 0 ; i < 100000 ; i++)
        verify_bitcounts (lrand48 ()) ;

    printf ("All BitCounts seem okay!  Starting speed trials\n") ;

    bitspeeds();

    return 0 ;
}

