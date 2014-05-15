Conway's Game of Life demo
========

This is demo solution of Conway's game of life problem simulated in a huge world.

[![Build Status](https://travis-ci.org/sitano/cgolhuge.png)](https://travis-ci.org/sitano/cgolhuge)
[![wercker status](https://app.wercker.com/status/91da64038f15c8fd4fdc8acca0101828/s/ "wercker status")](https://app.wercker.com/project/bykey/91da64038f15c8fd4fdc8acca0101828)

## Features

* Huge world (2^64^2) support
* World mapped on to closed sphere (edges are neighborhoods)
* Very sparsed life (empty mostly)
* Life locality
* Start state 100x100
* Persistent state (load / save)
* UI client + editor <=> server

# Getting Started

## Installing cgolhuge

### Manual build

You need to have the Go environment installed. To build and install cgolhuge, simply run:

```
$ go get github.com/sitano/cgolhuge
$ go install github.com/sitano/cgolhuge
```

You can build a release with

```
$ make release
```

`build/` then contains all you need.

## Running

## API

State:
* /world/setup
* /world/state[/x/y/w/h]
* /world/step[/{counts}]

Persistance:
* /world/save
* /world/load/{what}

## Possible solutions to big world problem

* virtual memory pages mapping
* compression
* (x, y) coord per point

## Possible optimizations

* hashlife
* cache on quad tree nodes

## Implementation (version 1)

* virtual memory pages 2^n (128x128 to match start window size = 16kb page)
* views: quad tree view (coords? origin)
* world ticker bitwised, 2 bit per life
* save / load rle???
* data compression

### Version 1 benchmarks (Mac Book Air)

```
BenchmarkWorldStep2pages          20      77425383 ns/op
BenchmarkWorldStep1pages          50      37956369 ns/op
BenchmarkWorldReadPage2          500       5231006 ns/op
BenchmarkWorldRWPage         500       7288380 ns/op
BenchmarkWorldRWPageRaw8x8     10000        130273 ns/op
BenchmarkWorldRWPageRaw3x1     50000         60944 ns/op
BenchmarkWorldRPageRaw    200000         12766 ns/op
BenchmarkWorldRPageRaw_ConvUint64      50000         50514 ns/op
BenchmarkWorldRWPageRaw       100000         25576 ns/op
```

So, raw full scan of 16kb with 8x8 read and 1 write per op lasts 130000ns/page.
Raw read 1 byte per op 12766ns/page.

If page fill factor is <1% (16kb <-> 4 gliders (20)) => we can skip whole page with
it's def raw read speed 1 page read = 13k ns/page+lookup.

### Version 1 bench + norealloc on page query

```
BenchmarkWorldStep2pages          50      30619028 ns/op
BenchmarkWorldStep1pages         100      14154527 ns/op
BenchmarkWorldReadPage2         1000       2108376 ns/op
BenchmarkWorldRWPage         500       4579450 ns/op
BenchmarkWorldRWPageRaw8x8     10000        137833 ns/op
BenchmarkWorldRWPageRaw3x1     50000         61063 ns/op
BenchmarkWorldRPageRaw    200000         12787 ns/op
BenchmarkWorldRPageRaw_ConvUint64      50000         50690 ns/op
BenchmarkWorldRWPageRaw       100000         25515 ns/op
```

Min read speed 12k ns/page -> 1 ms = 10^6 ns / 12k ns ~ 76 scans of empty pages per 1 ms -> 76k/sec max
See cpu profile for this at pprof.

#### TODO

Engine:

```
- VM to hold page tiles themselves
- Is it faster to alloc new 4-16kb page or to memset0 reclaimed one?
- Is it faster to read and sum [4kb]int or [16kb]byte?
- Reuse go compiler opt flags (http://dave.cheney.net/2012/10/07/notes-on-exploring-the-compiler-flags-in-the-go-compiler-suite)

- Const page size and bits count and mask
- Const world size
- Store pages in a tree in page coords
- Eliminate all math.Max/Min checks
- Than fix contains Point method =
- Eliminate any alloc during step
- Eliminate z layers in a View api (use new pages for processing)
- Efficient raw page stepping:
    - Hold rect of life on page (to process only small amount of bytes)
    - Accumulate changes count per page processing, if it is zero, skip it.
    - Use raw array read / write, no div / mul, process inner rect of page first
    - Write separate processing for page edges
    - Count life on edges
    - If edges have life, check only ness adjacent pages
    - Use only fixed count of query pages from tree per step to check adj pages (cache them)
    - Do not process DEAD cells, process only LIFE inside active RECT inside PAGE
```

API:

```
- View: Reader, Writer
- World is a View
- RLE, LIF
- MvXY by any value, use MASK and SHIFT to calc page
- PrintPattern to any pos
- MirrorXY any rect
- GetPoints + InRect
```
