package retry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type TryFunc func() error

const (
	// DefaultTimeout 默认超时时间
	DefaultTimeout = 60 * time.Second
	// DefaultNum 默认执行次数
	DefaultNum = 3
	// DefaultInterval 默认重试间隔时间
	DefaultInterval = 1 * time.Second
)

func New(tryFunc TryFunc, opts ...Option) (err error) {
	opt := Options{
		ctx:      context.TODO(),
		timeout:  DefaultTimeout,
		num:      DefaultNum,
		interval: DefaultInterval,
	}

	for _, item := range opts {
		item(&opt)
	}

	var i = 0
	for i < opt.num {
		err = nil
		ch := make(chan error, 1)

		go func() {
			defer func() {
				if e := recover(); e != nil {
					ch <- fmt.Errorf("%v", e)
				}
			}()

			ch <- tryFunc()
		}()

		select {
		case err = <-ch:
			break
		case <-opt.ctx.Done():
			return errors.New("context cancelled")
		case <-time.After(opt.timeout + (100 * time.Millisecond)):
			err = errors.New("timeout")
		}

		if err == nil {
			return
		}

		log.Printf("index: %d, err: %v\n", i, err)
		time.Sleep(opt.interval)
		i++
	}

	funcPath := runtime.FuncForPC(reflect.ValueOf(tryFunc).Pointer()).Name()
	lastSlash := strings.LastIndex(funcPath, "/")
	funcName := funcPath[lastSlash+1:]
	return fmt.Errorf("func: %s, num: %d, err: %v", funcName, i, err)
}
