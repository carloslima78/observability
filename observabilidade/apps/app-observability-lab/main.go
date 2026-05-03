package main

import (
	"log"
	"os"

	"app-observability-lab/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Em containers, logs devem ir para stdout/stderr para Docker, Promtail e Loki coletarem.
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC)

	servidor := gin.Default()

	controllers.RegistrarRotasDeMetricas(servidor)
	controllers.RegistrarRotasDeLogs(servidor)
	controllers.RegistrarRotasDeTracing(servidor)

	log.Fatal(servidor.Run(":8080"))
}
