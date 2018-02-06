package testify

// TODO: can be outsourced as external package

type AssertTime struct {
	Expected int64
	Actual int64
}

func (at *AssertTime) GreaterThanOrEqual() bool {
	if at.Actual >= at.Expected {
		return true
	}
	return false
}
