package go_ds

const slots = 10

type Stack struct {
	top  int
	data []interface{}
}

func (s *Stack) Push(data interface{}) {
	if s.top++; s.top > len(s.data) {
		s.data = append(s.data, make([]interface{}, slots)...)
	}
	s.data[s.top-1] = data
}

func (s *Stack) Pop() interface{} {
	if s.top--; s.top < 0 {
		return nil
	}

	return s.data[s.top]
}

func NewStack() *Stack {
	return &Stack{
		data: make([]interface{}, slots),
	}
}
