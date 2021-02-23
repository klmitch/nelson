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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/klmitch/nelson/internal/parser"
)

func TestStateImplementsState(t *testing.T) {
	assert.Implements(t, (*parser.State)(nil), &state{})
}

func TestStateErrorBase(t *testing.T) {
	obj := &state{
		Text: "text",
	}

	result := obj.Error(nil)

	assert.ErrorIs(t, result, ErrInvalid)
}

func TestStateErrorWithError(t *testing.T) {
	obj := &state{
		Text: "text",
	}

	result := obj.Error(assert.AnError)

	assert.ErrorIs(t, result, ErrInvalid)
}

func TestStateGetBase(t *testing.T) {
	obj := &state{
		Text:      "foo12345bar",
		ExclStart: true,
		IPos:      3,
		State:     stateStart,
	}

	result, err := obj.Get(8)

	assert.NoError(t, err)
	assert.Equal(t, int64(12345), result)
	assert.Equal(t, &state{
		Text:      "foo12345bar",
		ExclStart: true,
		IPos:      3,
		State:     stateStart,
	}, obj)
}

func TestStateGetEmptyStart(t *testing.T) {
	obj := &state{
		Text:      "foo12345bar",
		ExclStart: true,
		IPos:      3,
		State:     stateStart,
	}

	result, err := obj.Get(3)

	assert.NoError(t, err)
	assert.Equal(t, int64(math.MinInt64), result)
	assert.Equal(t, &state{
		Text:       "foo12345bar",
		EmptyStart: true,
		ExclStart:  false,
		IPos:       3,
		State:      stateStart,
	}, obj)
}

func TestStateGetEmptyEnd(t *testing.T) {
	obj := &state{
		Text:      "foo12345bar",
		ExclStart: true,
		IPos:      3,
		State:     stateEnd,
	}

	result, err := obj.Get(3)

	assert.NoError(t, err)
	assert.Equal(t, int64(math.MaxInt64), result)
	assert.Equal(t, &state{
		Text:      "foo12345bar",
		ExclStart: true,
		ExclEnd:   true,
		IPos:      3,
		State:     stateEnd,
	}, obj)
}

func TestParseClosedClosed(t *testing.T) {
	result, err := Parse("[1,7]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 1,
		End:   8,
	}, result)
}

func TestParseClosedOpen(t *testing.T) {
	result, err := Parse("[1,7)")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 1,
		End:   7,
	}, result)
}

func TestParseOpenClosed(t *testing.T) {
	result, err := Parse("(1,7]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 2,
		End:   8,
	}, result)
}

func TestParseOpenOpen(t *testing.T) {
	result, err := Parse("(1,7)")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 2,
		End:   7,
	}, result)
}

func TestParseEmptyClosedClosed(t *testing.T) {
	result, err := Parse("[]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: math.MinInt64,
		End:   math.MaxInt64,
	}, result)
}

func TestParseEmptyOpenOpen(t *testing.T) {
	result, err := Parse("()")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: math.MinInt64,
		End:   math.MaxInt64,
	}, result)
}

func TestParseCommaClosedClosed(t *testing.T) {
	result, err := Parse("[,]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: math.MinInt64,
		End:   math.MaxInt64,
	}, result)
}

func TestParseOneClosedClosed(t *testing.T) {
	result, err := Parse("[5]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 5,
		End:   6,
	}, result)
}

func TestParseOneClosedOpet(t *testing.T) {
	result, err := Parse("[5)")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 5,
		End:   6,
	}, result)
}

func TestParseOneOpenClosed(t *testing.T) {
	result, err := Parse("(5]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 5,
		End:   6,
	}, result)
}

func TestParseOneOpenOpen(t *testing.T) {
	result, err := Parse("(5)")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 5,
		End:   6,
	}, result)
}

func TestParseMinClosedClosed(t *testing.T) {
	result, err := Parse("[,7]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: math.MinInt64,
		End:   8,
	}, result)
}

func TestParseMinClosedOpen(t *testing.T) {
	result, err := Parse("[,7)")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: math.MinInt64,
		End:   7,
	}, result)
}

func TestParseMinOpenClosed(t *testing.T) {
	result, err := Parse("(,7]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: math.MinInt64,
		End:   8,
	}, result)
}

func TestParseMinOpenOpen(t *testing.T) {
	result, err := Parse("(,7)")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: math.MinInt64,
		End:   7,
	}, result)
}

func TestParseMaxClosedClosed(t *testing.T) {
	result, err := Parse("[1,]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 1,
		End:   math.MaxInt64,
	}, result)
}

func TestParseMaxClosedOpen(t *testing.T) {
	result, err := Parse("[1,)")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 1,
		End:   math.MaxInt64,
	}, result)
}

func TestParseMaxOpenClosed(t *testing.T) {
	result, err := Parse("(1,]")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 2,
		End:   math.MaxInt64,
	}, result)
}

func TestParseMaxOpenOpen(t *testing.T) {
	result, err := Parse("(1,)")

	assert.NoError(t, err)
	assert.Equal(t, Interval{
		Start: 2,
		End:   math.MaxInt64,
	}, result)
}

func TestParseNoText(t *testing.T) {
	result, err := Parse("")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}

func TestParseBadInit(t *testing.T) {
	result, err := Parse("1, 7]")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}

func TestParseOverflow(t *testing.T) {
	result, err := Parse("[11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111]")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}

func TestParseBadSep(t *testing.T) {
	result, err := Parse("[1;7]")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}

func TestParseBadClose(t *testing.T) {
	result, err := Parse("[1,7>")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}

func TestParseExtraText(t *testing.T) {
	result, err := Parse("[1,7] ")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}

func TestParseExtraShort(t *testing.T) {
	result, err := Parse("[1,7")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}

func TestParseInverted(t *testing.T) {
	result, err := Parse("[7,1]")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}

func TestParseSameOpen(t *testing.T) {
	result, err := Parse("(1,1]")

	assert.ErrorIs(t, err, ErrInvalid)
	assert.Equal(t, Interval{}, result)
}
