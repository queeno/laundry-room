package loader

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Loader struct {
	app    App
	ctx    context.Context
	cancel context.CancelFunc
}

func NewLoader(app App) *Loader {
	ctx, cancel := context.WithCancel(context.Background())

	return &Loader{
		app:    app,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (l Loader) Run() error {
	defer l.cancel()

	exitCh := make(chan struct{})

	go func() {
		// This starts the app
		// and will wait until the child process is done
		l.app.Run(l.ctx)

		exitCh <- struct{}{}
	}()

	termCh := make(chan os.Signal, 2)
	signal.Notify(termCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-termCh
		l.cancel()
	}()

	<-exitCh
	return nil
}
