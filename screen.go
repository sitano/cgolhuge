package main

import "os"
import "os/exec"
import "fmt"
import "regexp"
import "strconv"

var screenRow1, _ = regexp.Compile("([\\d]+) rows")
var screenRow2, _ = regexp.Compile("rows ([\\d]+)")
var screenCol1, _ = regexp.Compile("([\\d]+) columns")
var screenCol2, _ = regexp.Compile("columns ([\\d]+)")

// http://www.termsys.demon.co.uk/vtansi.htm
type Screen struct {
	rows, cols int
}

func NewScreen() *Screen {
	cmd := exec.Command("stty", "-a")
	cmd.Stdin = os.Stdin
	stty, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("Can't run stty -a, error = %v", err))
	}
	rows := 0
	cols := 0
	if screenRow1.Match(stty) {
		sm := screenRow1.FindSubmatch(stty)
		rows, _ = strconv.Atoi(string(sm[1]))
	}
	if screenRow2.Match(stty) {
		sm := screenRow2.FindSubmatch(stty)
		rows, _ = strconv.Atoi(string(sm[1]))
	}
	if screenCol1.Match(stty) {
		sm := screenCol1.FindSubmatch(stty)
		cols, _ = strconv.Atoi(string(sm[1]))
	}
	if screenCol2.Match(stty) {
		sm := screenCol2.FindSubmatch(stty)
		cols, _ = strconv.Atoi(string(sm[1]))
	}
	return &Screen{rows, cols}
}

func (s *Screen) Reset() {
	fmt.Printf("\033[2J\033c")
}

func (s *Screen) EraseUp() {
	fmt.Printf("\033[1J")
}

func (s *Screen) EraseDown() {
	fmt.Printf("\033[J")
}

func (s *Screen) EraseLine() {
	fmt.Printf("\033[2K")
}

func (s *Screen) Goto(row int, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func (s *Screen) GotoUp(count int) {
	fmt.Printf("\033[%dA", count)
}

func (s *Screen) GotoDown(count int) {
	fmt.Printf("\033[%dB", count)
}

func (s *Screen) GotoForward(count int) {
	fmt.Printf("\033[%dC", count)
}

func (s *Screen) GotoBack(count int) {
	fmt.Printf("\033[%dD", count)
}

func (s *Screen) PushCursor() {
	fmt.Printf("\033[7")
}

func (s *Screen) PopCursor() {
	fmt.Printf("\033[8")
}

func (s *Screen) Print(block ...interface {}) {
	fmt.Print(block...)
}

func (s *Screen) Printf(format string, block ...interface {}) {
	fmt.Printf(format, block...)
}

func (s *Screen) Println(block ...interface {}) {
	fmt.Println(block...)
}

func (s *Screen) PrintAt(row int, col int, block interface{}) {
	s.Goto(row, col)
	for _, c := range fmt.Sprintf("%v", block) {
		if c == '\n' {
			row ++
			s.Goto(row, col)
			continue
		}
		if c == '\r' {
			s.Goto(row, col)
			continue
		}
		fmt.Print(string(c))
	}
}

func (s *Screen) DisableInputBuffering() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
}

func (s *Screen) HideInputChars() {
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func (s *Screen) ShowInputChars() {
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	fmt.Println()
}
