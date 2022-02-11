package api

import (
	"context"
)

type Server interface {
	Run(port ...string) error
}

func StartServer(ctx context.Context, server Server) error {
	errChan := make(chan error)

	go func() {
		if err := server.Run(":8099"); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return nil
		}
	}
}
