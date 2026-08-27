package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jb843051627/quasar-weave/internal/handler"
	"github.com/jb843051627/quasar-weave/internal/service"
	"github.com/jb843051627/quasar-weave/internal/store"
)

func main() {
	path := os.Getenv("QUASAR_WEAVE_DB")
	if path == "" {
		path = "data/quasar-weave.db"
	}
	repository, err := store.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	app := service.NewLab(repository)
	defer app.Close()
	if _, err := app.EnsureDefaultGate(context.Background()); err != nil {
		log.Fatal(err)
	}
	addr := os.Getenv("QUASAR_WEAVE_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("quasar-weave listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.New(app)))
}
