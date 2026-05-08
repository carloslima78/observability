# Infraestrutura de Traces

Este manual prepara o Tempo e o Grafana para visualizar traces.

## Stack

```text
app-observability-lab
  -> OpenTelemetry
  -> Tempo
  -> Grafana
```

## Ferramentas

**OpenTelemetry** instrumenta a aplicação e gera spans.

**Tempo** recebe e armazena traces.

**Grafana** consulta o Tempo e mostra a árvore de spans.

## Arquivos

- `docker-compose.yml`: sobe Tempo junto com a stack.
- `infra/traces/tempo-config.yaml`: configura o recebimento de traces.
- `infra/metrics/grafana-datasources.yaml`: provisiona o Data Source Tempo.
- `app/controllers/tracing.go`: gera traces na aplicação Go.

## Configuração

O Tempo recebe spans no protocolo Jaeger compacto:

```text
tempo:6831/udp
```

A aplicação recebe o destino por variáveis:

```text
OTEL_EXPORTER_JAEGER_AGENT_HOST=tempo
OTEL_EXPORTER_JAEGER_AGENT_PORT=6831
```

Lembrete de bolso:

```text
host dentro do Docker Compose usa nome do servico, nao localhost
```

## Subir a Stack

Execute na raiz do laboratório:

```bash
docker compose up -d --build
```

Serviços esperados:

```text
grafana
loki
prometheus
promtail
tempo
app-observability-lab
```

## Verificar Tempo

```bash
curl "http://localhost:3200/ready"
```

Saída esperada:

```text
ready
```

## Gerar Traces

```bash
curl "http://localhost:8080/tracing/demo"
curl "http://localhost:8080/tracing/order/123"
curl "http://localhost:8080/tracing/error"
curl "http://localhost:8080/tracing/latency?seconds=0.5"
```

## Ver No Grafana

Acesse:

```text
http://localhost:3000
```

Credenciais:

```text
Usuario: admin
Senha: admin
```

No Grafana:

```text
Explore -> Tempo -> Search
```

Procure pelo serviço:

```text
app-observability-lab
```

Consultas TraceQL úteis:

```text
{name="GET /tracing/demo"}
{name="GET /tracing/order/:id"}
{name="tracing.processar_pedido"}
```

## Próximo Manual

Depois que a infraestrutura estiver de pé, siga para [Aplicação app-observability-lab: Tracing](aplicacao-go-traces.md).
