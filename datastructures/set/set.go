package set

type IntSet struct {
	data map[int]struct{}
}

func NewIntSet() *IntSet {
	data := make(map[int]struct{})
	return &IntSet{
		data: data,
	}
}
