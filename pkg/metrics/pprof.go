package metrics

import (
	"context"
	"log"
	"net/http"
	"sync"
)

func RunPprof(ctx context.Context, wg *sync.WaitGroup) {
	srv := &http.Server{
		Addr: "localhost:6061",
	}

	go func() {
		log.Println("📊 pprof server started at http://localhost:6061/debug/pprof/")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ pprof server error: %v\n", err)
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
	log.Println("🛑 Shutting down pprof server...")
	wg.Done()
}
