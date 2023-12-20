package app

import (
	"fmt"

	"github.com/rasulov-emirlan/poc-auth/internal/storage/mongodb"
)

func (a *application) initDB() error {
	mdb, err := mongodb.NewRepoCombiner(a.ctx, a.cfg, a.logger)
	if err != nil {
		return fmt.Errorf("failed to init mongodb: %w", err)
	}

	a.mdb = mdb

	a.cleanupFuncs = append(a.cleanupFuncs, func() {
		if err := a.mdb.Close(a.ctx); err != nil {
			a.logger.ErrorContext(
				a.ctx,
				"failed to close mongodb connection",
				"err", err.Error(),
			)
		}
	})

	a.logger.InfoContext(a.ctx, "mongodb connection established")

	return nil
}
