package test

import "testing"

func AssertError(t *testing.T, expected string, err error) {
	if err.Error() != expected {
		t.Errorf("failed assert error expected = %s, got %s", expected, err.Error())
	}
}

func AssertTrue(t *testing.T, val bool) {
	if val != true {
		t.Errorf("failed assert true")
	}
}

func AssertFalse(t *testing.T, val bool) {
	if val != false {
		t.Errorf("failed assert false")
	}
}

func AssertEqual(t *testing.T, expected interface{}, val interface{}) {
	if expected != val {
		t.Error("failed assert equal")
	}
}

func AssertNil(t *testing.T, val interface{}) {
	if val != nil {
		t.Errorf("failed assert nil")
	}
}

func AssertNotNil(t *testing.T, val interface{}) {
	if val == nil {
		t.Errorf("failed assert not nil")
	}

}
