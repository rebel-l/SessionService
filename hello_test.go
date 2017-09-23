package SessionService

import "testing"

func TestMainEntry(t *testing.T) {
	main()
}

func TestMultiply(t *testing.T) {
	if (multiply(2,3) != 6) {
		t.Fatalf("2 x 3 is not 6")
	}
}