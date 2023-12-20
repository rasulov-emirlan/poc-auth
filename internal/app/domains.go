package app

import (
	"fmt"

	"github.com/rasulov-emirlan/poc-auth/internal/domains/auth"
)

func (a *application) initDomains() error {
	authDomain, err := auth.NewService(a.ctx, auth.ServiceConfigs{
		UsersRepository: a.mdb.Users(),
		Logger:          a.logger,
		Cfg:             a.cfg,
	})

	if err != nil {
		return fmt.Errorf("failed to init auth domain: %w", err)
	}

	a.authDomain = authDomain

	a.logger.InfoContext(a.ctx, "domains initialized")

	return nil
}
