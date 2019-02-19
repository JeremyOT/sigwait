package sigwait

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// StopWaitable represents a service that can be stopped and waited upon.
type StopWaitable interface {
	// Stop and return a channel that is closed after cleanup, usually Wait()
	Stop() <-chan struct{}
	// Wait returns a channel that is closed after the task is fully stopped.
	Wait() <-chan struct{}
}

// ExitOnSignal attempts a clean exit on signal. If invoked twice, it will
// call os.Exit with -1.
func ExitOnSignal(s StopWaitable, sigChan <-chan os.Signal) {
	sig := <-sigChan
	log.Println("Received signal", sig)
	log.Print("Exiting...")
	select {
	case <-s.Stop():
		return
	case <-sigChan:
		log.Printf("Force quitting...")
		os.Exit(-1)
	}
}

// RunUntilSignal blocks until the StopWaitable stops. It triggers stop on any of
// syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT.
func RunUntilSignal(s StopWaitable) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	go ExitOnSignal(s, sigChan)
	<-s.Wait()
}
