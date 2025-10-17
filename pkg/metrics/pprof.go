package metrics

import (
	"context"
	"log"
	"net/http"
	"sync"
)

func RunPprof(ctx context.Context, wg sync.WaitGroup) {
	srv := &http.Server{
		Addr: "localhost:6061",
	}

	// run server ใน background goroutine
	go func() {
		log.Println("📊 pprof server started at http://localhost:6061/debug/pprof/")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ pprof server error: %v\n", err)
		}
	}()

	// รอ context ถูก cancel เพื่อ shutdown server
	<-ctx.Done()
	log.Println("🛑 Shutting down pprof server...")
	_ = srv.Shutdown(ctx)
	wg.Done()
}
