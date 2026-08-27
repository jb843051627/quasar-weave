package service

import "context"

func (l *Lab) Readiness(ctx context.Context) error {
	if err := l.ensureOpen(); err != nil {
		return err
	}
	_, err := l.store.Count(ctx, "observation")
	return err
}
