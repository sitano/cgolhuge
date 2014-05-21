package main

import (
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

var rleXYformat, _ = regexp.Compile("x\\W*=\\W*([\\d+])\\W*,\\W*y\\W*=\\W*([\\d+])")
var rleLFformat, _ = regexp.Compile("[0-9!ob$]+")

func LoadRLE(v View, x uint64, y uint64, filename string) {
	var wr io.Writer
	var bb AABB

	w := uint64(0)
	h := uint64(0)
	wbuf := []byte{0}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("LoadRLE error for View(%v) at (%d, %d) of %s", v, x, y, filename))
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty
		if len(line) < 1 {
			continue
		}
		// Skip comments
		if line[0] == '#' {
			if len(line) == 1 {
				continue
			}
			if line[1] != 'C' && line[1] != 'N' && line[1] != 'O' {
				panic(fmt.Sprintf("LoadRLE error for View(%v) at (%d, %d) of %s: unknown format at %s", v, x, y, filename, line))
			}
			continue
		}
		// Find width / height
		if rleXYformat.MatchString(line) {
			wh := rleXYformat.FindStringSubmatch(line)
			// Parse
			n, _ := strconv.Atoi(wh[1])
			w = uint64(n)
			n, _  = strconv.Atoi(wh[2])
			h = uint64(n)
			// Prepare
			bb = NewXYWH(x, y, w, h)
			vu := v.(ViewUtil)
			wr = vu.Writer(bb)
			continue
		}
		if wr != nil && rleLFformat.MatchString(line) {
			rep := 1
			rep_str := ""
			vx := bb.MinX
			vy := bb.MinY
			pc := ' '
			for _, c := range line {
				if c == '!' {
					return
				}
				if c >= '0' && c <= '9' {
					rep_str += string(c)
				}
				if c == '$' {
					if vx > bb.MinX && vx <= bb.MaxX {
						c = 'b'
						rep = int(bb.MaxX - vx + 1)
						rep_str = ""
					}
					if vx == bb.MinX && pc == '$' {
						c = 'b'
						rep = int(bb.SizeX())
						rep_str = ""
					}
				}
				if c == 'b' || c == 'o' {
					rep = 1
					// Repeater
				    if len(rep_str) > 0 {
						rep, _ = strconv.Atoi(rep_str)
						rep_str = ""
					}
					// Write
					for ri := 0; ri < rep; ri ++ {
						if c == 'b' {
							wbuf[0] = DEAD
						}
						if c == 'o' {
							wbuf[0] = LIFE
						}
						n, _ := wr.Write(wbuf)
						if n != 1 {
							panic(fmt.Sprintf("LoadRLE error for View(%v) at (%d, %d) of %s: failed to write %s", v, x, y, filename, string(c)))
						}
						vx ++
						if vx > bb.MaxX {
							vx = bb.MinX
							vy ++
						}
					}
				}
				pc = c
			}
		}
	}
}
