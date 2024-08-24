package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/soypat/goeda"
)

func main() {
	in := goeda.NewNet("IN")
	out := goeda.NewNet("OUT")
	gnd := goeda.NewNet("GND")
	circuit := divider(in, out, gnd, 0.33, 1e3, 0.1)
	fmt.Println(circuit.String())
}

func divider(in, out, gnd goeda.Net, factor, multiplier, error float64) *goeda.Circuit {
	// r1 := newresistor(factor-())
	var b big.Rat
	b.SetFloat64(factor)
	d := float64(b.Denom().Uint64())
	n := float64(b.Num().Uint64())
	if (factor-n/d)/factor > error {
		panic("bad")
	}

	c := goeda.Circuit{}
	r1 := newresistor((d - n) * multiplier)
	r2 := newresistor(n * multiplier)
	c.AddConnections(
		goeda.Join(in, r1.Pad(1)),
		goeda.Join(out, r1.Pad(2), r2.Pad(1)),
		goeda.Join(gnd, r2.Pad(2)),
	)
	return &c
}

func newresistor(R float64) Resistor {
	x := Resistor{Value: R, BaseComponent: goeda.BaseComponent{Name: getAssignmentVarname(1)}}
	x.SetPads(
		goeda.NewPad(1, "p1"),
		goeda.NewPad(2, "p2"),
	)
	return x
}

type Resistor struct {
	Value float64
	goeda.BaseComponent
}

func getAssignmentVarname(skipStackFrames int) string {
	if skipStackFrames < 0 {
		return ""
	}
	var buf [1024]byte
	n := runtime.Stack(buf[:], false)
	file := buf[:n]
	for i := 0; i < 4+skipStackFrames*2; i++ {
		_, file, _ = bytes.Cut(file, []byte("\n"))
	}
	file, _, ok := bytes.Cut(file, []byte{'\n'})
	if !ok {
		return ""
	}
	file = bytes.TrimSpace(file)
	space := bytes.LastIndex(file, []byte{' '})
	file = file[:space]
	lineIdx := bytes.LastIndex(file, []byte{':'})
	lineNum, _ := strconv.Atoi(string(file[lineIdx+1:]))
	filename := string(file[:lineIdx])
	fp, err := os.Open(filename)
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	currentLine := 0
	for currentLine < lineNum && scanner.Scan() {
		currentLine++
	}
	line := scanner.Text()
	line = strings.TrimSpace(line)
	v, _, _ := strings.Cut(line, " ")
	return v
}
