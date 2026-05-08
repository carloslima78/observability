package controllers

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("app-observability-lab/controllers/tracing")

// ConfigurarTracing prepara o OpenTelemetry para enviar spans ao Tempo.
func ConfigurarTracing(ctx context.Context) (func(context.Context) error, error) {
	host := os.Getenv("OTEL_EXPORTER_JAEGER_AGENT_HOST")
	if host == "" {
		// Valor local util quando a aplicacao roda fora do Docker Compose.
		host = "localhost"
	}

	porta := os.Getenv("OTEL_EXPORTER_JAEGER_AGENT_PORT")
	if porta == "" {
		porta = "6831"
	}

	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(host),
			jaeger.WithAgentPort(porta),
		),
	)
	if err != nil {
		return nil, err
	}

	recurso, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", "app-observability-lab"),
			attribute.String("deployment.environment", "local"),
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		// BatchSpanProcessor envia em lotes; e o comportamento comum em aplicacoes reais.
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(recurso),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider.Shutdown, nil
}

// UsarTracing cria o span principal HTTP antes das rotas executarem.
func UsarTracing(servidor *gin.Engine) {
	servidor.Use(func(contexto *gin.Context) {
		if contexto.Request.URL.Path == "/metrics" {
			contexto.Next()
			return
		}

		nomeDoSpan := contexto.Request.Method + " " + contexto.Request.URL.Path
		ctx, span := tracer.Start(
			contexto.Request.Context(),
			nomeDoSpan,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		contexto.Request = contexto.Request.WithContext(ctx)
		contexto.Next()

		status := contexto.Writer.Status()
		rota := contexto.FullPath()
		if rota == "" {
			rota = "rota_nao_mapeada"
		}

		span.SetName(contexto.Request.Method + " " + rota)
		span.SetAttributes(
			attribute.String("http.method", contexto.Request.Method),
			attribute.String("http.route", rota),
			attribute.Int("http.status_code", status),
		)
		if status >= 500 {
			span.SetStatus(codes.Error, http.StatusText(status))
		}
	})
}

// RegistrarRotasDeTracing concentra rotas didaticas do pilar de tracing.
func RegistrarRotasDeTracing(servidor *gin.Engine) {
	// Cria um trace curto com spans filhos para visualizar hierarquia no Tempo.
	servidor.GET("/tracing/demo", func(contexto *gin.Context) {
		ctx, span := tracer.Start(contexto.Request.Context(), "tracing.demo")
		defer span.End()

		validarRequisicao(ctx)
		consultarDependencia(ctx)
		montarResposta(ctx)

		contexto.JSON(http.StatusOK, gin.H{
			"mensagem": "pilar de tracing ativo",
		})
	})

	// Simula um fluxo de negocio com spans filhos nomeados por etapa.
	servidor.GET("/tracing/order/:id", func(contexto *gin.Context) {
		ctx, span := tracer.Start(contexto.Request.Context(), "tracing.processar_pedido")
		defer span.End()

		orderID := contexto.Param("id")
		span.SetAttributes(attribute.String("order.id", orderID))

		validarPedido(ctx, orderID)
		consultarCliente(ctx)
		calcularFrete(ctx)
		processarPagamento(ctx)
		salvarPedido(ctx)

		contexto.JSON(http.StatusOK, gin.H{
			"mensagem": "pedido simulado processado",
			"order_id": orderID,
		})
	})

	// Marca um span com erro para estudar como falhas aparecem no Grafana Tempo.
	servidor.GET("/tracing/error", func(contexto *gin.Context) {
		ctx, span := tracer.Start(contexto.Request.Context(), "tracing.erro_simulado")
		defer span.End()

		err := errors.New("falha simulada no processamento")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		consultarDependencia(ctx)

		contexto.JSON(http.StatusInternalServerError, gin.H{
			"erro": err.Error(),
		})
	})

	// Cria uma etapa lenta para observar duracao do span no trace.
	servidor.GET("/tracing/latency", func(contexto *gin.Context) {
		ctx, span := tracer.Start(contexto.Request.Context(), "tracing.latencia_simulada")
		defer span.End()

		duracao := lerLatenciaDeTracing(contexto.Query("seconds"))
		span.SetAttributes(attribute.String("latency.simulated", duracao.String()))

		dormir(ctx, "tracing.sleep", duracao)

		contexto.JSON(http.StatusOK, gin.H{
			"latencia_simulada": duracao.String(),
		})
	})
}

func validarRequisicao(ctx context.Context) {
	dormir(ctx, "tracing.validar_requisicao", 20*time.Millisecond)
}

func consultarDependencia(ctx context.Context) {
	dormir(ctx, "tracing.consultar_dependencia", 80*time.Millisecond)
}

func montarResposta(ctx context.Context) {
	dormir(ctx, "tracing.montar_resposta", 15*time.Millisecond)
}

func validarPedido(ctx context.Context, orderID string) {
	_, span := tracer.Start(ctx, "pedido.validar")
	defer span.End()

	span.SetAttributes(attribute.String("order.id", orderID))
	time.Sleep(25 * time.Millisecond)
}

func consultarCliente(ctx context.Context) {
	dormir(ctx, "pedido.consultar_cliente", 70*time.Millisecond)
}

func calcularFrete(ctx context.Context) {
	dormir(ctx, "pedido.calcular_frete", 45*time.Millisecond)
}

func processarPagamento(ctx context.Context) {
	dormir(ctx, "pedido.processar_pagamento", 120*time.Millisecond)
}

func salvarPedido(ctx context.Context) {
	dormir(ctx, "pedido.salvar", 35*time.Millisecond)
}

// dormir cria um span filho simples: abre span, simula trabalho e fecha span.
func dormir(ctx context.Context, nomeDoSpan string, duracao time.Duration) {
	_, span := tracer.Start(ctx, nomeDoSpan)
	defer span.End()

	span.SetAttributes(attribute.String("step.duration", duracao.String()))
	time.Sleep(duracao)
}

func lerLatenciaDeTracing(valor string) time.Duration {
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
