package goeda

import "fmt"

type Circuit struct {
	conns []Connection
}

func (c *Circuit) String() (s string) {
	for i, conn := range c.conns {
		s += "("
		if conn.net.label != "" {
			s += conn.net.label + " "
		}
		for j, p := range conn.pads {
			s += p.String()
			if j != len(conn.pads)-1 {
				s += " "
			}
		}
		s += ")"
		if i != len(c.conns)-1 {
			s += " "
		}
	}
	return s
}

func (c *Circuit) AddConnections(conns ...Connection) {
	c.conns = append(c.conns, conns...)
}

type Connection struct {
	net  Net
	pads []Pad
}

func JoinPads(pads ...Pad) Connection {
	return Connection{pads: pads}
}

func Join(net Net, pads ...Pad) (conn Connection) {
	return Connection{net: net, pads: pads}
}

type Component interface {
	Pad(num int) Pad
	ForEachPad(fn func(Pad) error) error
}

type Node interface {
	Name() string
	isNode()
}

// Compile time check of interface implementation.
var (
	_ Component = &BaseComponent{}
	_ Node      = &node{}
)

type BaseComponent struct {
	Name string
	pads []Pad
}

type Net struct {
	*node // The embedded pointer ensures unique comparison validity and absolute data immutability.
}

type Pad struct {
	*node
	parent Component
}

func NewNet(label string) Net {
	return Net{node: &node{num: 0, label: label}}
}

func (p Pad) Num() int { return p.num }

func NewPad(num int, label string) Pad {
	return Pad{node: &node{num: num, label: label}}
}

func (c *BaseComponent) SetPads(pads ...Pad) {
	last := -1
	for i := range pads {
		if last >= pads[i].num { // Guarantees sorted order and pad uniqueness.
			panic("pads must be sorted by number in increasing order and >=0")
		}
		last = pads[i].num
		pads[i].parent = c
	}
	c.pads = pads
}

func (c *BaseComponent) Pad(num int) Pad {
	for i := range c.pads {
		if c.pads[i].num == num {
			return c.pads[i]
		}
	}
	return NonexistPad()
}

func (c *BaseComponent) ForEachPad(fn func(p Pad) error) error {
	for i := range c.pads {
		err := fn(c.pads[i])
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	ne = NewPad(-2, "ENOENT")
)

// NonexistPad is the not exist pad. Returned when seeking nonexistent pad.
func NonexistPad() Pad {
	return ne
}

type node struct {
	num   int
	label string
}

func (n *node) Name() string { return n.label }

func (n *node) isNode() {}

func (n *node) goString() string {
	if n.num < 0 {
		return n.label
	}
	return fmt.Sprintf("{%s:%d}", n.label, n.num)
}

func (n *node) String() string {
	if n.num < 0 {
		return n.label
	}
	return fmt.Sprintf("%s:%d", n.label, n.num)
}
