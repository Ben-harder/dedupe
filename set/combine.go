package set

// Combine combines the passed sents into a single set and returns it.
func Combine(sets ...*Set) *Set {
	combined := New()
	for _, s := range sets {
		for k := range s.m {
			combined.Insert(k)
		}
	}
	return combined
}

// CombineFromChan expects n Sets to combine passed through setChan. Returns the resulting combined set.
func CombineFromChan(n int, setChan chan *Set) *Set {
	combined := New()
	for i := 0; i < n; i++ {
		s := <-setChan
		combined = Combine(combined, s)
	}
	return combined
}
