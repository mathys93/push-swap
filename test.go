package pushswap

import fmt "fmt"	*

func (s *Stack) Push(value int) {
	s.values = append(s.values, value)
}