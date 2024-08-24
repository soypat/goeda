package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/soypat/geda/io/dsn"
)

const dsnFile = `(pcb example.dsn
  (resolution um 10)
  (unit um)
  (structure
	(layer F.Cu
	  (type signal)
	  (property
		(index 0)
	  )
	)
	(layer B.Cu
	  (type signal)
	  (property
		(index 1)
	  )
	)
	(boundary
	  (path pcb 0  205105 -132080  106553 -132080  106553 -78253.2  205105 -78253.2
			205105 -132080)
	)

	(plane GND (polygon F.Cu 0  155067 -129413  151765 -129413  154051 -126492  154051 -121666
	143002 -121666  143002 -125984  145415 -129159  141732 -129159
	141859 -120777  155067 -120777  155067 -129413))

	(placement
	 (component "Package_SO:SOIC-8_3.9x4.9mm_P1.27mm"
	  (place U7 129094.000000 -110617.000000 front 180.000000 (PN MAX3072E))
	 )
	 (component "Package_TO_SOT_SMD:SOT-223"
	  (place Q3 177379.149000 -109851.275000 front 0.000000 (PN FQT13N06L))
	 )
	)
)`

func main() {
	const source = "example.dsn"
	// One can also use a file as a source:
	// fp, _ := os.Open(filename)
	fp := strings.NewReader(dsnFile) // Use string as source.
	var l dsn.Lexer
	err := l.Reset(source, fp)
	if err != nil {
		log.Fatal(err)
	}
	parser, err := dsn.NewParser(&l)
	if err != nil {
		log.Fatal(err)
	}
	decls, err := parser.ParseFilter(func(b []byte) bool {
		return true // parse all declarations.
	})
	if err != nil {
		log.Fatal(err)
	}
	for i := range decls {
		PrintDecl(decls[i])
	}
}

func PrintDecl(d dsn.Decl) {
	indent := strings.Repeat(" ", d.Depth())
	fmt.Print(indent, "(", d.Name())
	for _, arg := range d.Args() {
		fmt.Print(" ", arg.Literal)
	}
	children := d.Children()
	if len(children) > 0 {
		fmt.Print("\n")
		for _, child := range d.Children() {
			PrintDecl(child)
		}
		fmt.Print(indent, ")\n")
	} else {
		fmt.Print(")\n")
	}
}
