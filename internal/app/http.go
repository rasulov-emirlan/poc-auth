package app

import "github.com/rasulov-emirlan/poc-auth/internal/transport/rest"

func (a *application) initHttp() error {
	srvr := rest.NewServer(a.cfg)

	a.cleanupFuncs = append(a.cleanupFuncs, func() {
		if err := srvr.Shutdown(a.ctx); err != nil {
			a.logger.ErrorContext(
				a.ctx,
				"failed to shutdown http server",
				"err", err.Error(),
			)
		}

		a.logger.InfoContext(a.ctx, "http server stopped")
	})

	go func() {
		if err := srvr.Start(); err != nil {
			a.fatal("failed to start http server", err)
		}
	}()

	return nil
}
