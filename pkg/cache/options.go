package cache

import (
	"time"
)

// Options like sweeping functionality, cache size, etc
type Options struct {
	sweepInterval time.Duration
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
		options.sweepInterval = interval
		return nil
	}
}
