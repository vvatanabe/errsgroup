package errsgroup

import (
	"context"
	"errors"
	"testing"
)

func Test_Group(t *testing.T) {
	cases := []struct {
		errs []error
	}{
		{errs: []error{}},
		{errs: []error{errors.New("err-1")}},
		{errs: []error{errors.New("err-1"), errors.New("err-2"), errors.New("err-3")}},
	}
	for _, tc := range cases {
		g := NewGroup(LimitSize(10), ErrorChanelSize(10))
		for _, err := range tc.errs {
			err := err
			g.Go(func() error {
				return err
			})
		}
		if errs := g.Wait(); len(errs) != len(tc.errs) {
			t.Errorf("Could not match error size. want: %d, result: %d", len(tc.errs), len(errs))
		}
	}
}

func Test_WithContext(t *testing.T) {
	g, ctx := WithContext(context.Background(), LimitSize(10), ErrorChanelSize(10))
	for _, err := range []error{errors.New("err-1"), errors.New("err-2"), errors.New("err-3")} {
		err := err
		g.Go(func() error {
			return err
		})
	}
	if errs := g.Wait(); len(errs) != 3 {
		t.Errorf("Could not match error size. want: %d, result: %d", 3, len(errs))
	}
	canceled := false
	select {
	case <-ctx.Done():
		canceled = true
	default:
	}
	if !canceled {
		t.Errorf("Could not cancel")
	}
}
