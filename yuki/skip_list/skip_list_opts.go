package skiplist

type SkipListOpt func(*SkipList)

func WithMaxHeight(maxHeight int) SkipListOpt {
	return func(sl *SkipList) {
		sl.maxHeight = maxHeight
	}
}
func WithBranchingFactor(branchingFactor int) SkipListOpt {
	return func(sl *SkipList) {
		sl.branchingFactor = branchingFactor
	}
}
func WithCMPFunc(cmp CompareFunc) SkipListOpt {
	return func(sl *SkipList) {
		sl.cmp = cmp
	}
}
