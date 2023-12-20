package history

//COMPLETE

//since we're building a proper terminal tool, we need functionality like where we can press
//up button to access previous commands and toggle with the help of down arrow to next command
//this is exactly what's being handled in this file


// History is a struct that stores the history of user inputs
type History struct {
	inputs map[int]string // map of input history
	cursor int            // current cursor position
}

// NewHistory is just a helper function that returns a new History struct
//which is empty and you can add to it using the add function
func NewHistory() *History {
	return &History{
		map[int]string{},
		0,
	}
}

// Reset is another helper function that resets the history
func (h *History) Reset() *History {
	h.inputs = map[int]string{}
	h.cursor = 0

	return h
}

// Helper function Add - adds a new input to the history
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

// GetPrevious returns the previous input with respect to where the cursor 
//currently is
func (h *History) GetPrevious() *string {
	if input, ok := h.inputs[h.cursor]; ok {
		h.cursor--
		return &input
	}

	return nil
}

// GetNext returns the next input with respect to where the cursor
//currently is
func (h *History) GetNext() *string {
	if input, ok := h.inputs[h.cursor+1]; ok {
		h.cursor++
		return &input
	}

	return nil
}
