package main

import (
	"context"
	"log"
	"os"

	"app-observability-lab/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Em containers, logs devem ir para stdout/stderr para Docker, Promtail e Loki coletarem.
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC)

	finalizarTracing, err := controllers.ConfigurarTracing(context.Background())
	if err != nil {
		log.Printf("level=error msg=%q", "falha ao configurar tracing")
	} else {
		defer finalizarTracing(context.Background())
	}

	servidor := gin.Default()

	controllers.UsarTracing(servidor)

	controllers.RegistrarRotasDeMetricas(servidor)
	controllers.RegistrarRotasDeLogs(servidor)
	controllers.RegistrarRotasDeTracing(servidor)

	log.Fatal(servidor.Run(":8080"))
}
