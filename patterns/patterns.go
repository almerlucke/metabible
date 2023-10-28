package patterns

import (
	"io"

	xbnf "github.com/almerlucke/exbana/v2"
)

type RuneStream struct {
	values     []rune
	pos        int
	mismatches []*xbnf.Mismatch[rune, int]
}

func NewRuneStream(str string) *RuneStream {
	return &RuneStream{
		values: []rune(str),
		pos:    0,
	}
}

func (ts *RuneStream) Peek1() (rune, error) {
	if ts.pos < len(ts.values) {
		return ts.values[ts.pos], nil
	}

	return 0, io.EOF
}

func (ts *RuneStream) Peek(n int, buf []rune) (int, error) {
	i := 0
	p := ts.pos
	l := len(ts.values)

	for i < n && p < l {
		buf[i] = ts.values[p]
		p++
		i++
	}

	if i != n {
		return i, io.EOF
	}

	return i, nil
}

func (ts *RuneStream) Read1() (rune, error) {
	if ts.pos < len(ts.values) {
		v := ts.values[ts.pos]
		ts.pos++
		return v, nil
	}

	return 0, io.EOF
}

func (ts *RuneStream) Read(n int, buf []rune) (int, error) {
	i := 0
	l := len(ts.values)

	for i < n && ts.pos < l {
		buf[i] = ts.values[ts.pos]
		ts.pos++
		i++
	}

	if i != n {
		return i, io.EOF
	}

	return i, nil
}

func (ts *RuneStream) Skip(n int) (int, error) {
	m := len(ts.values) - ts.pos
	if n > m {
		n = m
	}

	ts.pos += n

	if n == m {
		return n, io.EOF
	}

	return n, nil
}

func (ts *RuneStream) Finished() bool {
	return ts.pos >= len(ts.values)
}

func (ts *RuneStream) Position() (int, error) {
	if ts.pos >= len(ts.values) {
		return ts.pos, io.EOF
	}

	return ts.pos, nil
}

func (ts *RuneStream) SetPosition(pos int) error {
	ts.pos = pos
	if ts.pos >= len(ts.values) {
		return io.EOF
	}

	return nil
}

func (ts *RuneStream) Log(mismatch *xbnf.Mismatch[rune, int]) {
	ts.mismatches = append(ts.mismatches, mismatch)
}

func (ts *RuneStream) Range(begin int, end int) ([]rune, error) {
	return ts.values[begin:end], nil
}

func (ts *RuneStream) Write(objs ...rune) error {
	ts.values = append(ts.values, objs...)

	return nil
}

func (ts *RuneStream) Runes() []rune {
	return ts.values
}

func (ts *RuneStream) Finish() error {
	return nil
}

func runeEntityEqual(o1 rune, o2 rune) bool {
	return o1 == o2
}

func runeMatch(r rune) *xbnf.UnitPattern[rune, int] {
	return xbnf.Unit[rune, int](func(obj rune) bool {
		return obj == r
	})
}

func runeFuncMatch(rf func(rune) bool) *xbnf.UnitPattern[rune, int] {
	return xbnf.Unit[rune, int](func(obj rune) bool {
		return rf(obj)
	})
}

func runeFuncMatchx(id string, rf func(rune) bool) *xbnf.UnitPattern[rune, int] {
	return xbnf.Unitx[rune, int](id, false, func(obj rune) bool {
		return rf(obj)
	})
}

func runeSeries(str string) *xbnf.SeriesPattern[rune, int] {
	return xbnf.Series[rune, int](runeEntityEqual, []rune(str)...)
}

func IsObjectReference(str string) bool {
	// example: recNI7VFXIGfqeeKj
	match, _, _ := xbnf.Concat[rune, int](runeSeries("rec"), xbnf.N[rune, int](runeFuncMatch(func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	}), 14)).Match(NewRuneStream(str), nil)

	return match
}
