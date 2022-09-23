package retry

import (
	"context"
	"time"
)

type Options struct {
	// 上下文
	ctx context.Context
	// 超时控制
	timeout time.Duration
	// 执行次数
	num int
	// 间隔时间
	interval time.Duration
}

// Option db操作选项
type Option func(exec *Options)

// Ctx 设置上下文
func Ctx(ctx context.Context) Option {
	return func(exec *Options) {
		exec.ctx = ctx
	}
}

// Timeout 设置超时时间
func Timeout(timeout time.Duration) Option {
	return func(exec *Options) {
		exec.timeout = timeout
	}
}

// Num 尝试次数
func Num(num int) Option {
	return func(exec *Options) {
		exec.num = num
	}
}

// Interval 间隔时间
func Interval(interval time.Duration) Option {
	return func(exec *Options) {
		exec.interval = interval
	}
}
