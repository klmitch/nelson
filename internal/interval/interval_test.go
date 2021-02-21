package interval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntervalStringBase(t *testing.T) {
	obj := Interval{
		Start: 1,
		End:   7,
	}

	result := obj.String()

	assert.Equal(t, "[1,7)", result)
}

func TestIntervalStringOne(t *testing.T) {
	obj := Interval{
		Start: 1,
		End:   2,
	}

	result := obj.String()

	assert.Equal(t, "[1]", result)
}

func TestIntervalIncludesLow(t *testing.T) {
	obj := Interval{
		Start: 1,
		End:   7,
	}

	result := obj.Includes(0)

	assert.False(t, result)
}

func TestIntervalIncludesStart(t *testing.T) {
	obj := Interval{
		Start: 1,
		End:   7,
	}

	result := obj.Includes(1)

	assert.True(t, result)
}

func TestIntervalIncludesMidpoint(t *testing.T) {
	obj := Interval{
		Start: 1,
		End:   7,
	}

	result := obj.Includes(4)

	assert.True(t, result)
}

func TestIntervalIncludesEndpoint(t *testing.T) {
	obj := Interval{
		Start: 1,
		End:   7,
	}

	result := obj.Includes(6)

	assert.True(t, result)
}

func TestIntervalIncludesEnd(t *testing.T) {
	obj := Interval{
		Start: 1,
		End:   7,
	}

	result := obj.Includes(7)

	assert.False(t, result)
}

func TestIntervalIncludesHigh(t *testing.T) {
	obj := Interval{
		Start: 1,
		End:   7,
	}

	result := obj.Includes(8)

	assert.False(t, result)
}
