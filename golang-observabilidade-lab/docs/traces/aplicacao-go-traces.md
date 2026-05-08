# Aplicação app-observability-lab: Tracing

Este guia explica como a aplicação gera traces com OpenTelemetry.

A infraestrutura com Tempo e Grafana está em [Infraestrutura de Traces](infraestrutura-traces.md).

## Conceito Rápido

**Trace** é a história completa de uma operação.

**Span** é uma etapa dessa história.

Exemplo:

```text
GET /tracing/order/123
  -> tracing.processar_pedido
  -> pedido.validar
  -> pedido.consultar_cliente
  -> pedido.calcular_frete
  -> pedido.processar_pagamento
  -> pedido.salvar
```

Lembrete de bolso:

```text
trace = fluxo inteiro
span = uma etapa do fluxo
span filho = uma etapa dentro de outra etapa
```

## Fluxo

```text
Aplicação Go
  -> OpenTelemetry SDK
  -> Jaeger exporter
  -> Grafana Tempo
  -> Grafana
```

O OpenTelemetry instrumenta a aplicação. O Tempo armazena os traces. O Grafana visualiza.

## Arquivo Principal

```text
app/controllers/tracing.go
```

Este arquivo concentra o pilar de tracing:

- `ConfigurarTracing`: configura OpenTelemetry.
- `UsarTracing`: cria o span principal HTTP.
- `RegistrarRotasDeTracing`: registra rotas de estudo.
- funções como `validarPedido`, `consultarCliente` e `processarPagamento`: criam spans filhos.

## Pontos Importantes do Código

### Configuração

```text
ConfigurarTracing
```

Cria o exporter e o `TracerProvider`.

Decisão do laboratório:

```text
OpenTelemetry -> protocolo Jaeger -> Tempo
```

Esse caminho é simples para estudo local. Em ambientes modernos de produção, OTLP costuma ser a escolha preferida.

### Middleware

```text
UsarTracing
```

Cria um span principal para cada requisição HTTP.

Ele ignora `/metrics`, porque scrape do Prometheus não representa fluxo de usuário.

Ele adiciona atributos como:

```text
http.method
http.route
http.status_code
```

Esse é o span de borda da aplicação.

### Spans Manuais

As funções do controller criam spans filhos.

Exemplo:

```text
pedido.processar_pagamento
```

Esses spans mostram onde o tempo foi gasto dentro da regra de negócio.

## Rotas

### Demo

```text
GET /tracing/demo
```

Exemplo:

```bash
curl "http://localhost:8080/tracing/demo"
```

Cria um trace pequeno:

```text
tracing.demo
  -> tracing.validar_requisicao
  -> tracing.consultar_dependencia
  -> tracing.montar_resposta
```

### Pedido

```text
GET /tracing/order/:id
```

Exemplo:

```bash
curl "http://localhost:8080/tracing/order/123"
```

Cria um fluxo de negócio simulado.

### Erro

```text
GET /tracing/error
```

Exemplo:

```bash
curl "http://localhost:8080/tracing/error"
```

Marca o span com erro usando:

```text
RecordError
SetStatus(codes.Error)
```

### Latência

```text
GET /tracing/latency?seconds=0.5
```

Exemplo:

```bash
curl "http://localhost:8080/tracing/latency?seconds=0.5"
```

Cria uma etapa lenta para observar duração no Grafana.

## Dependências Go

Se precisar atualizar manualmente as dependências do pilar:

```bash
go get go.opentelemetry.io/otel@v1.31.0
go get go.opentelemetry.io/otel/sdk@v1.31.0
go get go.opentelemetry.io/otel/exporters/jaeger@v1.17.0
go mod tidy
```

No laboratório, essas dependências já ficam registradas em:

```text
app/go.mod
app/go.sum
```

## Como Testar

Suba a stack:

```bash
docker compose up -d --build
```

Gere traces:

```bash
curl "http://localhost:8080/tracing/demo"
curl "http://localhost:8080/tracing/order/123"
curl "http://localhost:8080/tracing/error"
curl "http://localhost:8080/tracing/latency?seconds=0.5"
```

Acesse:

```text
http://localhost:3000
```

No Grafana:

```text
Explore -> Data Source Tempo -> Search
```

Procure pelo serviço:

```text
app-observability-lab
```

Consultas TraceQL rápidas:

```text
{name="GET /tracing/demo"}
{name="GET /tracing/order/:id"}
{name="tracing.processar_pedido"}
```

## Estrutura Básica de Um Span

Um bom span deve ter:

- nome claro
- início e fim
- duração
- status
- atributos úteis
- erro, quando existir

Lembrete:

```text
nome bom de span parece etapa de trabalho, nao frase generica
```

Exemplos bons:

```text
pedido.validar
pedido.consultar_cliente
pedido.processar_pagamento
```
