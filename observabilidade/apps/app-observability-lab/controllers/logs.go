package controllers

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RegistrarRotasDeLogs concentra as rotas usadas no pilar de logs.
func RegistrarRotasDeLogs(servidor *gin.Engine) {
	// Rota raiz: serve como teste rapido no navegador e gera um log DEBUG.
	servidor.GET("/", func(contexto *gin.Context) {
		registrarLogPorNivel("debug", "rota raiz acessada")
		contexto.String(http.StatusOK, "app-observability-lab com Gin\n")
	})

	// Health check: confirma que a aplicacao esta respondendo e gera um log INFO.
	servidor.GET("/health", func(contexto *gin.Context) {
		registrarLogPorNivel("info", "health check recebido")
		contexto.String(http.StatusOK, "ok")
	})

	// Gerador de logs: recebe level e message pela URL.
	// Exemplo: /logs?level=error&message=falha-ao-processar
	servidor.GET("/logs", func(contexto *gin.Context) {
		nivel := contexto.Query("level")
		mensagem := contexto.Query("message")

		registrarLogPorNivel(nivel, mensagem)

		contexto.JSON(http.StatusOK, gin.H{
			"mensagem": mensagem,
			"nivel":    nivel,
		})
	})

	// Log estruturado: exemplo de log com contexto de negocio e operacao.
	servidor.GET("/logs/estruturado", func(contexto *gin.Context) {
		criaLogEstruturado()

		contexto.JSON(http.StatusOK, gin.H{
			"mensagem": "log estruturado de pedido criado",
		})
	})

	registrarLogPorNivel("info", "aplicacao go iniciada na porta 8080")
	iniciarHeartbeat()
}

// registrarLogPorNivel recebe um nivel e uma mensagem, depois escreve o log
// no formato usado nas consultas do Loki: level=<nivel> msg=<mensagem>.
func registrarLogPorNivel(nivel string, mensagem string) {
	nivel = strings.ToLower(nivel)

	switch nivel {
	case "debug":
		// Aqui nasce um log DEBUG, util para detalhes de desenvolvimento.
		log.Printf("level=debug msg=%q", mensagem)
	case "info":
		// Aqui nasce um log INFO, usado para eventos esperados da aplicacao.
		log.Printf("level=info msg=%q", mensagem)
	case "warn":
		// Aqui nasce um log WARN, usado para alertas sem quebrar a aplicacao.
		log.Printf("level=warn msg=%q", mensagem)
	case "error":
		// Aqui nasce um log ERROR, usado quando algo falhou e precisa de atencao.
		log.Printf("level=error msg=%q", mensagem)
	default:
		// Se o nivel vier desconhecido, registramos como INFO para manter a demo simples.
		log.Printf("level=info msg=%q", mensagem)
	}
}

// criaLogEstruturado mostra um exemplo mais proximo de ambiente real.
// Alem da mensagem, o log carrega contexto como servico, squad, ambiente e conta AWS.
func criaLogEstruturado() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info(
		"pedido criado",
		"service", "orders-api",
		"squad", "payments",
		"environment", "production",
		"aws_account_id", "123456789012",
		"region", "us-east-1",
		"order_id", "987",
	)
}

// iniciarHeartbeat gera um log DEBUG periodico para mostrar que a aplicacao segue viva.
func iniciarHeartbeat() {
	// Em um goroutine separado, para nao bloquear o servidor HTTP.
	go func() {
		for {
			registrarLogPorNivel("debug", "heartbeat da aplicacao")
			time.Sleep(30 * time.Second)
		}
	}()
}
