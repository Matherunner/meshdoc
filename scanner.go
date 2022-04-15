package meshdoc

import "bufio"

type LineScanner struct {
	scanner *bufio.Scanner
	i       int
	eof     bool
	peekBuf string
	peeked  bool
}

func NewLineScanner(scanner *bufio.Scanner) *LineScanner {
	eof := !scanner.Scan()
	return &LineScanner{scanner: scanner, eof: eof}
}

func (s *LineScanner) Scan() bool {
	s.i++
	s.eof = !s.scanner.Scan()
	s.peeked = false
	return !s.eof
}

func (s *LineScanner) Line() string {
	if s.peeked {
		return s.peekBuf
	}
	return s.scanner.Text()
}

func (s *LineScanner) LineNumber() int {
	return s.i
}

func (s *LineScanner) EOF() bool {
	return s.eof
}

func (s *LineScanner) Peek() bool {
	if s.eof {
		return false
	}
	s.peekBuf = s.scanner.Text()
	s.peeked = true
	return true
}
