package cache

import (
	"time"
)

// Options like sweeping functionality, cache size, etc
type Options struct {
	SweepInterval      time.Duration
	AutoReloadInterval time.Duration
}

func (os *Options) Apply(o ...Option) error {
	for _, o := range o {
		if err := o(os); err != nil {
			return err
		}
	}

	return nil
}

type Option func(options *Options) error

// WithSweeping make sure that sweeping with happen at specified interval
func WithSweeping(interval time.Duration) Option {
	return func(options *Options) error {
		options.SweepInterval = interval
		return nil
	}
}

// WithAutoReload make sure that cache got update will happen at specified interval
func WithAutoReload(interval time.Duration) Option {
	return func(options *Options) error {
		options.SweepInterval = interval
		return nil
	}
}
