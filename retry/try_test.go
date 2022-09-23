package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNew_Ctx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var i int
	err := New(
		func() (err error) {
			i++
			switch i {
			case 1:
				return errors.New("error")
			case 2:
				cancel()
			}
			return
		},
		Ctx(ctx),
		Timeout(1*time.Second),
	)

	fmt.Printf("%v\n", err)
}

func TestNew_Timeout(t *testing.T) {
	var i int
	err := New(
		func() (err error) {
			i++
			if i <= 2 {
				time.Sleep(3 * time.Second)
			}
			return
		},
		Timeout(1*time.Second),
	)

	fmt.Printf("%v\n", err)
}

func TestNew_Error(t *testing.T) {
	err := New(
		func() (err error) {
			return errors.New("error")
		},
		Timeout(1*time.Second),
	)

	fmt.Printf("%v\n", err)
}

func TestNew_Panic(t *testing.T) {
	err := New(
		func() (err error) {
			panic("panic")
			return
		},
		Timeout(1*time.Second),
	)

	fmt.Printf("%v\n", err)
}

func TestNew_TimeoutAndError(t *testing.T) {
	var i int
	err := New(
		func() (err error) {
			i++
			switch i {
			case 1:
				return errors.New("error")
			case 2:
				time.Sleep(3 * time.Second)
			}
			return
		},
		Timeout(1*time.Second),
	)

	fmt.Printf("%v\n", err)
}
