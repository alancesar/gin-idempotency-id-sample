package api

import (
	"context"
)

type Server interface {
	Run(addr ...string) error
}

func StartServer(ctx context.Context, server Server, addr ...string) error {
	errChan := make(chan error)

	go func() {
		if err := server.Run(addr...); err != nil {
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
