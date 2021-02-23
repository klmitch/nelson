// Copyright (c) 2021 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package interval

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode"

	"github.com/klmitch/nelson/internal/parser"
)

// ErrInvalid indicates an error parsing an interval expression.
var ErrInvalid = errors.New("invalid interval")

// Parser states.
const (
	stateInit  = iota // Initial state
	stateStart        // Reading start number
	stateSep          // Looking for comma separator
	stateEnd          // Reading end number
	stateClose        // Expecting the closer
	stateDone         // No more expected
)

// state describes the parser state.
type state struct {
	Text       string   // The text being parsed
	Ival       Interval // The interval being constructed
	EmptyStart bool     // Flag indicating a start was not provided
	ExclStart  bool     // Start interval is exclusive
	ExclEnd    bool     // End interval is exclusive
	IPos       int      // The starting position of an integer
	State      int      // State of the parse
}

// Error constructs a parser error.
func (s *state) Error(err error) error {
	if err == nil {
		return fmt.Errorf("%w %q", ErrInvalid, s.Text)
	}
	return fmt.Errorf("%w %q: %s", ErrInvalid, s.Text, err)
}

// Get extracts the integer from the interval expression.
func (s *state) Get(pos int) (int64, error) {
	// Is it empty?
	if s.IPos == pos {
		if s.State == stateStart {
			s.EmptyStart = true
			s.ExclStart = false
			return math.MinInt64, nil
		}

		s.ExclEnd = true
		return math.MaxInt64, nil
	}

	return strconv.ParseInt(s.Text[s.IPos:pos], 10, 64)
}

// Parse processes a single character from the input.
func (s *state) Parse(pos int, char rune) error {
	switch s.State {
	case stateInit:
		if char == '(' {
			s.ExclStart = true
		} else if char != '[' {
			return s.Error(nil)
		}
		s.State = stateStart
		s.IPos = pos + 1

	case stateStart, stateEnd:
		if !unicode.IsDigit(char) && !(s.IPos == pos && (char == '+' || char == '-')) {
			tmp, err := s.Get(pos)
			if err != nil {
				return s.Error(err)
			}
			if s.State == stateStart {
				s.Ival.Start = tmp
				s.State = stateSep
			} else {
				s.Ival.End = tmp
				s.State = stateClose
			}
			return s.Parse(pos, char)
		}

	case stateSep:
		if char == ',' {
			s.State = stateEnd
			s.IPos = pos + 1
			return nil
		}
		if char == ')' || char == ']' {
			if s.EmptyStart {
				s.Ival.End = math.MaxInt64
			} else {
				s.ExclStart = false
				s.Ival.End = s.Ival.Start + 1
			}
			s.ExclEnd = true
			s.State = stateClose
			return s.Parse(pos, char)
		}
		return s.Error(nil)

	case stateClose:
		if char == ')' {
			s.ExclEnd = true
		} else if char != ']' {
			return s.Error(nil)
		}
		s.State = stateDone

	case stateDone:
		return s.Error(nil)
	}

	return nil
}

// Parse parses a string into an Interval.
func Parse(text string) (Interval, error) {
	// Construct the state
	s := &state{
		Text: text,
	}

	// Text has to be at least 2 characters
	if len(text) < 2 {
		return Interval{}, s.Error(nil)
	}

	// Parse the text
	if err := parser.Parse(text, s); err != nil {
		return Interval{}, err
	}

	// Make sure we finished processing
	if s.State != stateDone {
		return Interval{}, s.Error(nil)
	}

	// Now, we need to canonicalize the interval
	if s.ExclStart {
		s.Ival.Start++
	}
	if !s.ExclEnd {
		s.Ival.End++
	}
	if s.Ival.End <= s.Ival.Start {
		return Interval{}, s.Error(nil)
	}

	return s.Ival, nil
}
