package history

// History is a struct that stores the history of user inputs
type History struct {
	inputs map[int]string // map of input history
	cursor int            // current cursor position
}

// NewHistory returns a new History struct
func NewHistory() *History {
	return &History{
		map[int]string{},
		0,
	}
}

// Reset resets the history
func (h *History) Reset() *History {
	h.inputs = map[int]string{}
	h.cursor = 0

	return h
}

// Add adds a new input to the history
func (h *History) Add(input string) *History {
	h.cursor = len(h.inputs)
	h.inputs[h.cursor] = input

	return h
}

// GetAll returns all the inputs in the history
func (h *History) GetAll() map[int]string {
	return h.inputs
}

// GetCursor returns the current cursor position
func (h *History) GetCursor() int {
	return h.cursor
}

// GetPrevious returns the previous input
func (h *History) GetPrevious() *string {
	if input, ok := h.inputs[h.cursor]; ok {
		h.cursor--
		return &input
	}

	return nil
}

// GetNext returns the next input
func (h *History) GetNext() *string {
	if input, ok := h.inputs[h.cursor+1]; ok {
		h.cursor++
		return &input
	}

	return nil
}
