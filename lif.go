package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"regexp"
	"strconv"
)

var lifXYformat, _ = regexp.Compile("#P[ ]+([-+]?[0-9]+)[ ]+([-+]?[0-9]+)")
var lifLFformat, _ = regexp.Compile("[.*]+")

func LoadLIF(v View, x uint64, y uint64, filename string) {
	vx := x
	ox := x
	vy := y

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
			if strings.HasPrefix(line, "#Life") {
				continue
			}
			if line[1] != 'C' && line[1] != 'N' && line[1] != 'O' && line[1] != 'P' && line[1] != 'D' {
				panic(fmt.Sprintf("LoadRLE error for View(%v) at (%d, %d) of %s: unknown format at %s", v, x, y, filename, line))
			}
			// Find width / height
			if lifXYformat.MatchString(line) {
				wh := lifXYformat.FindStringSubmatch(line)
				// Parse
				n, _ := strconv.Atoi(wh[1])
				if n >= 0 {
					vx = x + uint64(n)
				} else {
					vx = x - uint64(-1 * n)
				}
				ox = vx
				n, _ = strconv.Atoi(wh[2])
				if n >= 0 {
					vy = y + uint64(n)
				} else {
					vy = y - uint64(-1 * n)
				}
			}
			continue
		}
		if lifLFformat.MatchString(line) {
			for _, c := range line {
				if c == '.' {
					v.Set(vx, vy, DEAD)
				}
				if c == '*' {
					v.Set(vx, vy, LIFE)
				}
				vx ++
			}
			vx = ox
			vy ++
		}
	}
}
