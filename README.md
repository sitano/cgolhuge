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

BenchmarkWorldStep2pages	      20	  77425383 ns/op
BenchmarkWorldStep1pages	      50	  37956369 ns/op
BenchmarkWorldReadPage2	     500	   5231006 ns/op
BenchmarkWorldRWPage	     500	   7288380 ns/op
BenchmarkWorldRWPageRaw8x8	   10000	    130273 ns/op
BenchmarkWorldRWPageRaw3x1	   50000	     60944 ns/op
BenchmarkWorldRPageRaw	  200000	     12766 ns/op
BenchmarkWorldRPageRaw_ConvUint64	   50000	     50514 ns/op
BenchmarkWorldRWPageRaw	  100000	     25576 ns/op

So, raw full scan of 16kb with 8x8 read and 1 write per op lasts 130000ns/page.
Raw read 1 byte per op 12766ns/page.

If page fill factor is <1% (16kb <-> 4 gliders (20)) => we can skip whole page with
it's def raw read speed 1 page read = 13k ns/page+lookup.
