Conway's Game of Life demo
========

This is demo solution of Conway's game of life problem simulated in a huge world.

[![Build Status](https://travis-ci.org/sitano/cgolhuge.png)](https://travis-ci.org/sitano/cgolhuge)
[![wercker status](https://app.wercker.com/status/91da64038f15c8fd4fdc8acca0101828/s/ "wercker status")](https://app.wercker.com/project/bykey/91da64038f15c8fd4fdc8acca0101828)

## Features

* Huge world (2^64^2) support
* Sparsed life (empty mostly)
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

```
go build && ./cgolhuge -load ./pattern/glider_gun.lif -lx 20 -ly 20 -idle 100
```

## Available solutions to big world problem

* virtual memory pages mapping
* global compression
* (x, y) coord per point

## Optimizations todo

* hashlife
* cache on quad tree nodes

## Implementation (version 2)

* virtual memory pages 2^n (64x64, 1 bit/life, 1 page = 64 * uint64 = 64 * 8 byte = 512 byte)
* views: quad tree view
* load rle / lif

### Version 2 benchmarks (Mac Book Air) / single page with glider gun

```
go test -run no -bench World
PASS
BenchmarkWorldGliderGun	  100000	     18816 ns/op
ok  	github.com/sitano/cgolhuge	2.128s
```

#### TODO

Engine:

```

Version 2:

- Skip stride middle bits by 16 based on or of all 3 lines
- MemSet with C ext instead of realloc new buf on step
- Do not set old life - set dead instead (to reduce bits ops)
- No new buff for no changes during page step
- Page reclaimation

```
