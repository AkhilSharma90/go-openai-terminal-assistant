package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistory(t *testing.T) {
	// TestNewHistory tests the NewHistory function.
	t.Run("NewHistory", func(t *testing.T) {
		h := NewHistory()
		assert.NotNil(t, h)
		assert.Equal(t, 0, len(h.GetAll()))
		assert.Equal(t, 0, h.GetCursor())
	})

	// TestReset tests the Reset function.
	t.Run("Reset", func(t *testing.T) {
		h := NewHistory()
		h.Add("input1").Add("input2")
		h.Reset()
		assert.Equal(t, 0, len(h.GetAll()))
		assert.Equal(t, 0, h.GetCursor())
	})

	// TestAdd tests the Add function.
	t.Run("Add", func(t *testing.T) {
		h := NewHistory()
		h.Add("input1").Add("input2")
		assert.Equal(t, 2, len(h.GetAll()))
	})

	// TestGetAll tests the GetAll function.
	t.Run("GetAll", func(t *testing.T) {
		h := NewHistory()
		h.Add("input1").Add("input2")
		all := h.GetAll()
		assert.Equal(t, 2, len(all))
		assert.Equal(t, "input1", all[0])
		assert.Equal(t, "input2", all[1])
	})

	// TestGetPrevious tests the GetPrevious function.
	t.Run("GetPrevious", func(t *testing.T) {
		h := NewHistory()
		h.Add("input1").Add("input2")
		prev := h.GetPrevious()
		assert.NotNil(t, prev)
		assert.Equal(t, "input2", *prev)
		prev = h.GetPrevious()
		assert.NotNil(t, prev)
		assert.Equal(t, "input1", *prev)
	})

	// TestGetNext tests the GetNext function.
	t.Run("GetNext", func(t *testing.T) {
		h := NewHistory()
		h.Add("input1").Add("input2")
		h.GetPrevious()
		next := h.GetNext()
		assert.NotNil(t, next)
		assert.Equal(t, "input2", *next)
	})

	// TestGetOutOfBounds tests the GetOutOfBounds function.
	t.Run("GetOutOfBounds", func(t *testing.T) {
		h := NewHistory()
		h.Add("input1").Add("input2")
		h.GetPrevious()
		h.GetPrevious()
		prev := h.GetPrevious()
		assert.Nil(t, prev)

		h.GetNext()
		h.GetNext()
		next := h.GetNext()
		assert.Nil(t, next)
	})
}
