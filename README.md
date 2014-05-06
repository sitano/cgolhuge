Conway's Game of Life demo
========

This is demo solution of Conway's game of life problem simulated in a huge world.

Features
========

* Huge world (2^64^2) support
* World mapped on to closed sphere (edges are neighborhoods)
* Very sparsed life (empty mostly)
* Life locality
* Start state 100x100
* Persistent state (load / save)
* UI client + editor <=> server

Quick start
========

In progress...

API
========

State:
* /world/setup
* /world/state[/x/y/w/h]
* /world/step[/{counts}]
 
Persistance:
* /world/save
* /world/load/{what}

Possible solutions
========

* virtual memory pages mapping
* compression
* (x, y) coord per point

Possible optimizations
========

* hashlife
* cache on quad tree nodes

Implementation
========

* virtual memory pages 2^n
* views: pages hashed / quad tree view (coords? origin)
* world ticker bitwised, 3 bit per life
* save / load rle???
* data compression

Future scalability options
========

* quad trees
* supervised router
