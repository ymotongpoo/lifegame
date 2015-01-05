// See details of lifegame in Japanese.
//    http://ja.wikipedia.org/wiki/%E3%83%A9%E3%82%A4%E3%83%95%E3%82%B2%E3%83%BC%E3%83%A0
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Interval is display refresh interval.
const Interval = time.Second / 10

// Field holds cell data.
type Field struct {
	cs   [][]bool // field's memory
	w, h int      // field's width and height
}

// NewField returns a field which has w x h cells.
func NewField(h, w int) *Field {
	cs := make([][]bool, h)
	for i := range cs {
		cs[i] = make([]bool, w)
	}
	return &Field{cs: cs, w: w, h: h}
}

// Set sets cell's status.
func (f *Field) Set(r, c int, b bool) error {
	if r < 0 || r >= f.h || c < 0 || c >= f.w {
		return errors.New("out of field")
	}
	f.cs[r][c] = b
	return nil
}

// Alive confirm if specified cell is alive.
// This is utility function to check outbound field.
func (f *Field) Alive(r, c int) bool {
	r = (r + f.h) % f.h
	c = (c + f.w) % f.w
	return f.cs[r][c]
}

// NextGen returns if specified the cell of r & c will be alive
// in next generation.
func (f *Field) NextGen(r, c int) bool {
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (i != 0 || j != 0) && f.Alive(r+i, c+j) {
				alive++
			}
		}
	}
	return alive == 3 || alive == 2 && f.Alive(r, c)
}

// Print display one generation status to stdout.
func (f *Field) Print() {
	for _, r := range f.cs {
		bufr := make([]byte, f.w)
		for j, c := range r {
			if c {
				bufr[j] = 'o'
			} else {
				bufr[j] = ' '
			}
		}
		fmt.Println(string(bufr))
	}

}

// Life holds current and next generation field.
type Life struct {
	cur, next *Field
	gen       int
}

// NewLife create new lifegame buffer.
func NewLife(h, w int, init [][]bool) *Life {
	cur := NewField(h, w)
	next := NewField(h, w)
	cur.cs = init // TODO(ymotongpoo): confirm matrix size.
	return &Life{cur: cur, next: next, gen: 0}
}

// Next calculates each state of all cells in current field and set it in next.
// Swaps cur and next after calculation and proceed generation counter.
func (l *Life) Next() {
	for i, r := range l.cur.cs {
		for j := range r {
			l.next.Set(i, j, l.cur.NextGen(i, j))
		}
	}
	l.cur = l.next
	l.next = NewField(l.cur.w, l.cur.h)
	l.gen++
}

// Print display current generation status.
func (l *Life) Print() {
	cmd := exec.Command("clear") // TODO(ymotongpoo): Work out way to clear terminal on Windows.
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Printf("---------- %vth generation\n", l.gen)
	l.cur.Print()
}

func main() {
	fmt.Println("Lifegame")
	init := make([][]bool, 10)
	for i := range init {
		init[i] = make([]bool, 10)
	}

	// Beacon pattern
	// TODO(ymotongpoo): Add function to read initial state from text file.
	init[1][1] = true
	init[1][2] = true
	init[2][1] = true
	init[3][4] = true
	init[4][3] = true
	init[4][4] = true

	l := NewLife(10, 10, init)

	ticker := time.Tick(Interval)
	for range ticker {
		l.Print()
		l.Next()
	}
}
