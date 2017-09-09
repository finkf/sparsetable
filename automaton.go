package fsa

// Automaton is the basic interface for all automaton implementations.
type Automaton interface {
	// Delta does one transition from the current state(s) to the next.
	// Delta returns false if no transition could be done.
	Delta(byte) bool
	// Final returns true if a final state is active in the automaton and its
	// asscoiated data. If not it return nil, false.
	Final() (interface{}, bool)
	// Initialize initializes the automaton.
	// Intitialize should be the first function called on the automaton
	// before any matching.
	Initialize()
}

// Accepts tests if the given automaton accets the given string.
func Accepts(a Automaton, str string) (interface{}, bool) {
	a.Initialize()
	// use explicit loop to iterate over the bytes of the string
	for i := 0; i < len(str); i++ {
		// fmt.Printf("[%v] str[%v] = %v\n", str, i, str[i])
		if !a.Delta(str[i]) {
			return nil, false
		}
	}
	return a.Final()
}

// DeltaStar returns the longest possible accepted string in the automaton
// starting at str.
// It returns the data of the accepted string and the length of the match.
// If no string is accepted, nil and the first position in the string that
// equals sync is returned.
func DeltaStar(a Automaton, str string, sync byte) (interface{}, int) {
	a.Initialize()
	var data interface{}
	var pos int
	for i := range str {
		c := str[i]
		if pos == 0 && c == sync {
			pos = i
		}
		if !a.Delta(c) {
			return data, pos
		}
		if d, final := a.Final(); final {
			pos = i + 1
			data = d
		}
	}
	return data, pos
}
