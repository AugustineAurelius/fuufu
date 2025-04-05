package skiplist

type Iterator struct {
	current *Node
}

func (it *Iterator) Next() {
	if it.current != nil {
		it.current = it.current.next[0].Load()
	}
}

func (it *Iterator) Key() []byte {
	if it.current != nil {
		return it.current.key
	}
	return nil
}

func (it *Iterator) Value() []byte {
	if it.current != nil {
		return it.current.value
	}
	return nil
}

func (it *Iterator) Valid() bool {
	return it.current != nil
}
