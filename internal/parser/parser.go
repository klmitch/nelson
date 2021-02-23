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

package parser

// State describes a parser state.  The Parse function loops over a
// string, repeatedly calling the State.Parse method until all
// characters have been processed.
type State interface {
	// Parse processes a single character from the input.
	Parse(pos int, char rune) error
}

// Parse loops over characters in a string, applying the State.Parse
// method repeatedly until all characters have been processed.
func Parse(text string, s State) error {
	// Loop over the text
	for pos, char := range text {
		if err := s.Parse(pos, char); err != nil {
			return err
		}
	}

	return nil
}
