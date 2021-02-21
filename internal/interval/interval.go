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

import "fmt"

// Interval describes a interval of values.  A Interval is normalized
// to be a half-open interval, but the input text uses "[]" and "()"
// to indicate closed or open intervals, and anything in between.
type Interval struct {
	Start int64 // Start value of the interval (inclusive)
	End   int64 // End value of the interval (exclusive)
}

// String outputs a string version of the Interval object.
func (r Interval) String() string {
	// Handle the basic case
	if r.End <= r.Start+1 {
		return fmt.Sprintf("[%d]", r.Start)
	}

	// OK, construct the interval notation
	return fmt.Sprintf("[%d,%d)", r.Start, r.End)
}

// Includes tests to see if a specified number falls within the
// Interval.
func (r Interval) Includes(v int64) bool {
	return v >= r.Start && v < r.End
}
