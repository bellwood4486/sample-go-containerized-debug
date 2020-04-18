package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "hello world")
	})

	srv := http.Server{Addr: ":8080"}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGHUP,
			syscall.SIGQUIT)

		switch <-sigChan {
		case syscall.SIGINT:
			fmt.Println("Recieved: SIGINT")
		case syscall.SIGTERM:
			fmt.Println("Recieved: SIGTERM")
		case syscall.SIGHUP:
			fmt.Println("Recieved: SIGHUP")
		case syscall.SIGQUIT:
			fmt.Println("Recieved: SIGQUIT")
		default:
			fmt.Print("Unknown signal")
		}

		// received a signal, shut down...
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	fmt.Printf("Launching server at %q ...\n", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
