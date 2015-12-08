package the_platinum_searcher

type matchState interface {
	transition(matched bool) matchState
	reset() matchState
	isBefore() bool
	isMatching() bool
	isAfter() bool
}

func newMatchState() matchState {
	return stateBeforeMatch{}
}

type stateBeforeMatch struct{}

func (s stateBeforeMatch) transition(matched bool) matchState {
	if matched {
		return stateMatching{}
	} else {
		return s
	}
}

func (s stateBeforeMatch) reset() matchState {
	return s
}

func (s stateBeforeMatch) isBefore() bool {
	return true
}

func (s stateBeforeMatch) isMatching() bool {
	return false
}

func (s stateBeforeMatch) isAfter() bool {
	return false
}

type stateMatching struct{}

func (s stateMatching) transition(matched bool) matchState {
	if matched {
		return s
	} else {
		return stateAfterMatch{}
	}
}

func (s stateMatching) reset() matchState {
	return stateBeforeMatch{}
}

func (s stateMatching) isBefore() bool {
	return false
}

func (s stateMatching) isMatching() bool {
	return true
}

func (s stateMatching) isAfter() bool {
	return false
}

type stateAfterMatch struct{}

func (s stateAfterMatch) transition(matched bool) matchState {
	if matched {
		return stateMatching{}
	} else {
		return s
	}
}

func (s stateAfterMatch) reset() matchState {
	return stateBeforeMatch{}
}

func (s stateAfterMatch) isBefore() bool {
	return false
}

func (s stateAfterMatch) isMatching() bool {
	return false
}

func (s stateAfterMatch) isAfter() bool {
	return true
}
