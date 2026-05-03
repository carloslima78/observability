package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CounterVec: contador acumulado separado por metodo, rota e status HTTP.
// Ajuda a responder: "quantas requisicoes chegaram em cada parte da aplicacao?"
var requisicoesHTTP = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_http_requests_total",
		Help: "Total de requisicoes HTTP recebidas pela aplicacao.",
	},
	[]string{"method", "route", "status"},
)

// CounterVec: contador acumulado apenas para respostas de erro.
// Ajuda a separar taxa de erro por status, como 400, 404, 500 e 503.
var errosHTTP = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_http_errors_total",
		Help: "Total de respostas HTTP com status de erro.",
	},
	[]string{"method", "route", "status"},
)

// HistogramVec: mede a distribuicao da duracao das requisicoes.
// Os buckets representam faixas de tempo usadas depois em consultas de latencia.
var duracaoHTTP = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "app_http_request_duration_seconds",
		Help:    "Duracao das requisicoes HTTP em segundos.",
		Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	},
	[]string{"method", "route", "status"},
)

// Gauge: valor atual que sobe e desce conforme requisicoes entram e terminam.
// Ajuda a observar se ha muitas requisicoes simultaneas em andamento.
var requisicoesAtivas = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "app_http_active_requests",
		Help: "Quantidade de requisicoes HTTP em andamento.",
	},
)

// Counter: exemplo simples de metrica de negocio.
// Em um sistema real, poderia representar pedidos, pagamentos ou cadastros.
var pedidosCriados = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "app_orders_created_total",
		Help: "Total de pedidos simulados criados pela aplicacao.",
	},
)

func init() {
	// MustRegister registra as metricas assim que a aplicacao sobe.
	// Se houver erro de configuracao, a aplicacao falha cedo em vez de esconder o problema.
	prometheus.MustRegister(requisicoesHTTP)
	prometheus.MustRegister(errosHTTP)
	prometheus.MustRegister(duracaoHTTP)
	prometheus.MustRegister(requisicoesAtivas)
	prometheus.MustRegister(pedidosCriados)
}

// RegistrarRotasDeMetricas concentra a instrumentacao do pilar de metricas.
func RegistrarRotasDeMetricas(servidor *gin.Engine) {
	// Middleware que mede todas as rotas da aplicacao, exceto o proprio /metrics.
	servidor.Use(medirRequisicaoHTTP())

	// /metrics e o endpoint que o Prometheus acessa para coletar as metricas.
	servidor.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Rota didatica para confirmar que o pilar de metricas esta ativo.
	servidor.GET("/metrics/demo", func(contexto *gin.Context) {
		contexto.JSON(http.StatusOK, gin.H{
			"mensagem": "pilar de metricas ativo",
		})
	})

	// Simula latencia para alimentar o Histogram app_http_request_duration_seconds.
	servidor.GET("/metrics/latency", func(contexto *gin.Context) {
		segundos := lerDuracaoSimulada(contexto.Query("seconds"))
		time.Sleep(segundos)

		contexto.JSON(http.StatusOK, gin.H{
			"latencia_simulada": segundos.String(),
		})
	})

	// Simula erros HTTP para alimentar counters e consultas por status.
	servidor.GET("/metrics/error/:status", func(contexto *gin.Context) {
		status, err := strconv.Atoi(contexto.Param("status"))
		if err != nil || status < 400 || status > 599 {
			contexto.JSON(http.StatusBadRequest, gin.H{
				"mensagem": "informe um status HTTP entre 400 e 599",
			})
			return
		}

		contexto.JSON(status, gin.H{
			"mensagem": "erro simulado para estudo de metricas",
			"status":   status,
		})
	})

	// Simula uma metrica de negocio, separada das metricas tecnicas HTTP.
	servidor.POST("/metrics/orders", func(contexto *gin.Context) {
		pedidosCriados.Inc()

		contexto.JSON(http.StatusCreated, gin.H{
			"mensagem": "pedido simulado criado",
		})
	})
}

// medirRequisicaoHTTP calcula metricas comuns de HTTP apos cada requisicao.
func medirRequisicaoHTTP() gin.HandlerFunc {
	return func(contexto *gin.Context) {
		// Evita medir o scrape do Prometheus como se fosse trafego de usuario.
		if contexto.Request.URL.Path == "/metrics" {
			contexto.Next()
			return
		}

		inicio := time.Now()
		// Gauge sobe no inicio da requisicao e desce ao terminar.
		requisicoesAtivas.Inc()
		defer requisicoesAtivas.Dec()

		// Executa a rota real antes de registrar status e duracao.
		contexto.Next()

		status := strconv.Itoa(contexto.Writer.Status())
		rota := contexto.FullPath()
		if rota == "" {
			// Evita usar caminhos dinamicos desconhecidos como label de alta cardinalidade.
			rota = "rota_nao_mapeada"
		}

		// Labels escolhidas por baixa cardinalidade: metodo, rota padronizada e status.
		labels := prometheus.Labels{
			"method": contexto.Request.Method,
			"route":  rota,
			"status": status,
		}

		requisicoesHTTP.With(labels).Inc()
		duracaoHTTP.With(labels).Observe(time.Since(inicio).Seconds())

		// Apenas status 4xx e 5xx entram no contador especifico de erros.
		if contexto.Writer.Status() >= 400 {
			errosHTTP.With(labels).Inc()
		}
	}
}

// lerDuracaoSimulada interpreta o parametro seconds da rota /metrics/latency.
// O limite de 5 segundos evita sleeps longos demais no laboratorio local.
func lerDuracaoSimulada(valor string) time.Duration {
	if valor == "" {
		return 300 * time.Millisecond
	}

	segundos, err := strconv.ParseFloat(valor, 64)
	if err != nil || segundos < 0 {
		return 300 * time.Millisecond
	}

	if segundos > 5 {
		segundos = 5
	}

	return time.Duration(segundos * float64(time.Second))
}
