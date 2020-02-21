package gorm

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func (b *Backend) Begin() *Backend {
	res := &Backend{
		DB:      b.DB.Begin(),
		context: b.context,
	}
	res.DB.SetLogger(newLogger(res.context, res.LogLevel))
	return res
}

func (b *Backend) End(rErr error) error {
	if e := recover(); e != nil {
		rErr = multierror.Append(rErr, fmt.Errorf("panic: %+v", e))
		defer panic(rErr)
	}
	if rErr != nil {
		if err := b.DB.Rollback().Error; err != nil {
			return multierror.Append(rErr, fmt.Errorf("rolling back: %w", err))
		}
		return rErr
	}
	if err := b.DB.Commit().Error; err != nil {
		return fmt.Errorf("committing changes: %w", err)
	}
	return nil
}
