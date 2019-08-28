package loader

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"move-health-check/pkg/aggregator"
	"move-health-check/pkg/service"
)

// Helper functions
func returnEmptyContextWithCancel() context.Context {
	ctx, _ := context.WithCancel(context.Background())
	return ctx
}

func returnCancelFunc() context.CancelFunc {
	_, cancel := context.WithCancel(context.Background())
	return cancel
}

func newTestLoader() *Loader {
	// An instance of the application will only run for max 100ms for test purposes
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100)*time.Millisecond)

	return &Loader{
		app:    testApp{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func sendSignalToRun(t *testing.T, termSignal syscall.Signal, millisecBeforeKill int) bool {
	l := newTestLoader()
	defer l.cancel()

	go func() {
		_ = l.Run("") // using mock app
	}()

	exitChan := make(chan bool)

	go func() {
		for {
			select {
			case <-l.ctx.Done():
				if l.ctx.Err() == context.DeadlineExceeded {
					// The application was still running after 5s
					t.Log("Context deadline exceeded")
					exitChan <- false
					return
				}
				// termSignal caused Run to exit!
				t.Log("Application correctly terminated")
				exitChan <- true
				return
			case <-time.After(time.Duration(10) * time.Millisecond):
				t.Log("Waiting for the application to terminate...")
			}
		}
	}()

	time.Sleep(time.Duration(millisecBeforeKill) * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), termSignal)

	return <-exitChan
}

// New()
func TestNewLoader(t *testing.T) {
	loader := NewLoader()
	assert.IsType(t, &Loader{}, loader, "Loader not created")
	assert.IsType(t, returnCancelFunc(), loader.cancel, "CancelFunc not created when loader gets created")
}

func TestSelectAppAggregator(t *testing.T) {
	loader := NewLoader()
	loader.selectApp("aggregator")
	assert.IsType(t, aggregator.Aggregator{}, loader.app, "App not created when mode set to 'aggregator'")
}

func TestSelectAppService(t *testing.T) {
	loader := NewLoader()
	loader.selectApp("service")
	assert.IsType(t, service.Service{}, loader.app, "App not created when mode set to 'service'")
}

func TestSelectAppNoMode(t *testing.T) {
	loader := NewLoader()
	loader.selectApp("")
	assert.Nil(t, loader.app, "App created even though no mode was passed")
}

func TestSelectAppWrongMode(t *testing.T) {
	loader := NewLoader()
	loader.selectApp("simon")
	assert.Nil(t, loader.app, "App created even though no valid mode was passed")
}

// Run()
func TestRunWithWrongMode(t *testing.T) {
	err := NewLoader().Run("WrongMode")
	assert.Error(t, err, "Passing mode 'simon' should have caused an error", err)
}

func TestRunWithSIGTERM(t *testing.T) {
	assert.True(t, sendSignalToRun(t, syscall.SIGTERM, 30), "SIGTERM didn't stop the process after 30ms")
}

func TestRunWithSIGINT(t *testing.T) {
	assert.True(t, sendSignalToRun(t, syscall.SIGINT, 30), "SIGINT didn't stop the process after 30ms")
}

func TestRunWithSIGWINCH(t *testing.T) {
	assert.False(t, sendSignalToRun(t, syscall.SIGWINCH, 30), "SIGWINCH caused the process to stop after 30ms")
}

func TestRunWithNoSignal(t *testing.T) {
	assert.False(t, sendSignalToRun(t, syscall.Signal(0), 0), "Application doesn't seem to run")
}
