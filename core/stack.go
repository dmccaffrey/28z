package core

type (
	Stack[T any] struct {
		top    *node[T]
		length int
	}
	node[T any] struct {
		value T
		prev  *node[T]
	}
)

// Create a new stack
func NewCoreValueStack() Stack[CoreValue] {
	return Stack[CoreValue]{nil, 0}
}

// Return the number of items in the stack
func (this *Stack[T]) Len() int {
	return this.length
}

// View the top item on the stack
func (this *Stack[T]) Peek() *T {
	if this.length == 0 {
		return nil
	}
	return &this.top.value
}

// Pop the top item of the stack and return it
func (this *Stack[T]) Pop() *T {
	if this.length == 0 {
		return nil
	}

	n := this.top
	this.top = n.prev
	this.length--
	return &n.value
}

// Push a value onto the top of the stack
func (this *Stack[T]) Push(value T) {
	n := &node[T]{value, this.top}
	this.top = n
	this.length++
}

func (s *Stack[T]) ToArray() []T {
	result := make([]T, s.length)
	curr := s.top
	i := 0
	for curr != nil {
		result[i] = curr.value
		curr = curr.prev
		i++
	}
	return result
}
