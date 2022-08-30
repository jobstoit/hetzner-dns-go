package dns

import "testing"

type assert struct {
	t *testing.T
}

func newAssert(t *testing.T) *assert {
	return &assert{
		t: t,
	}
}

func (a assert) NotNil(x any) bool {
	if x == nil {
		a.t.Error("expected value but got nil")
	}

	return x != nil
}

func (a assert) NoError(err error) bool {
	if err != nil {
		a.t.Errorf("unexpected error: %v", err)
	}

	return err == nil
}

func (a assert) Error(err error) bool {
	if err == nil {
		a.t.Error("missing expected error")
	}

	return err != nil
}

func (a assert) EqStr(expected, actual string) bool {
	if expected != actual {
		a.t.Errorf("expected '%s' but got '%s'", expected, actual)
	}

	return expected == actual
}

func (a assert) EqInt(expected, actual int) bool {
	if expected != actual {
		a.t.Errorf("expected %d but got %d", expected, actual)
	}

	return expected == actual
}
