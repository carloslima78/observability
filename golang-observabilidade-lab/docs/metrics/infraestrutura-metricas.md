# Infraestrutura de Métricas

Este manual prepara o ambiente de infraestrutura para coleta e visualização de métricas.

A instrumentação da aplicação Go está descrita em [Aplicação app-observability-lab: Métricas](aplicacao-go-metricas.md).

## Stack

```text
app-observability-lab
  -> GET /metrics
  -> Prometheus
  -> Grafana
```

## Ferramentas

**Prometheus** coleta métricas expostas pela aplicação.

Paralelo: lembra o papel de coleta e armazenamento de séries temporais em ferramentas como **Datadog Metrics**, **CloudWatch Metrics** ou **New Relic Metrics**.

**Grafana** visualiza as métricas armazenadas no Prometheus.

Paralelo: lembra dashboards operacionais usados para acompanhar taxa de erro, latência, tráfego e sinais de saúde.

**app-observability-lab** é a aplicação de teste que expõe `/metrics`.

## Arquivos

- `docker-compose.yml`: sobe Prometheus junto com a stack do laboratório.
- `infra/metrics/prometheus.yml`: configura o Prometheus para coletar métricas da aplicação.
- `infra/metrics/grafana-datasources.yaml`: provisiona os Data Sources Loki e Prometheus no Grafana.
- `app/controllers/metrics.go`: expõe o endpoint `/metrics` usado pelo Prometheus.

## Arquivos de Configuração

### Prometheus

Arquivo:

```text
infra/metrics/prometheus.yml
```

Define como o Prometheus coleta métricas.

Pontos importantes:

- `scrape_interval: 5s`: frequência de coleta.
- `job_name: app-observability-lab`: nome do alvo no Prometheus.
- `metrics_path: /metrics`: rota da aplicação usada para coleta.
- `targets: app-observability-lab:8080`: endereço da aplicação dentro da rede Docker Compose.

Resumo:

```text
Prometheus container
  -> http://app-observability-lab:8080/metrics
```

Dentro do Docker Compose, o Prometheus usa o nome do serviço `app-observability-lab`, não `localhost`.

### Grafana

Arquivo:

```text
infra/metrics/grafana-datasources.yaml
```

Provisiona os Data Sources do Grafana.

Pontos importantes:

- `name: Loki`: Data Source usado para logs.
- `name: Prometheus`: Data Source usado para métricas.
- `url: http://prometheus:9090`: endereço interno do Prometheus na rede Docker Compose.
- `url: http://loki:3100`: endereço interno do Loki na rede Docker Compose.

## Subir a Stack

Execute a partir da raiz do laboratório:

```bash
docker compose up -d --build
```

Serviços esperados:

```text
grafana
loki
prometheus
promtail
app-observability-lab
```

## Verificar Prometheus

Acesse:

```text
http://localhost:9090
```

No Prometheus:

1. Abra **Status**.
2. Abra **Target health**.
3. Verifique se `app-observability-lab` aparece como `UP`.

Essa tela é o primeiro ponto de validação da coleta antes de montar consultas ou dashboards no Grafana.

Consulta rápida pelo navegador ou terminal:

```bash
curl "http://localhost:9090/api/v1/targets"
```

## Verificar Endpoint da Aplicação

Antes de olhar o Prometheus, você pode conferir se a aplicação está expondo métricas:

```bash
curl "http://localhost:8080/metrics"
```

Procure por métricas como:

```text
app_http_requests_total
app_http_request_duration_seconds_bucket
go_goroutines
```

## Acessar Grafana

Abra:

```text
http://localhost:3000
```

Credenciais:

```text
Usuário: admin
Senha: admin
```

O Data Source `Prometheus` já é provisionado automaticamente pelo arquivo `infra/metrics/grafana-datasources.yaml`.

## Criar Dashboard de Métricas

No Grafana:

1. Clique em **Dashboards**.
2. Clique em **New**.
3. Clique em **New dashboard**.
4. Clique em **Add visualization**.
5. Selecione o Data Source **Prometheus**.
6. Informe uma consulta PromQL.
7. Escolha a visualização.
8. Clique em **Apply**.
9. Clique em **Save dashboard**.

Sugestão de painéis:

- `Requests por segundo`: `sum(rate(app_http_requests_total[1m]))`
- `Erros por status`: `sum by (status) (rate(app_http_errors_total[1m]))`
- `% HTTP 500`: `100 * sum(rate(app_http_errors_total{status="500"}[1m])) / sum(rate(app_http_requests_total[1m]))`
- `Latência p95`: `histogram_quantile(0.95, sum by (le) (rate(app_http_request_duration_seconds_bucket[5m])))`
- `Latência p95 por rota`: `histogram_quantile(0.95, sum by (route, le) (rate(app_http_request_duration_seconds_bucket[5m])))`
- `Requisições ativas`: `app_http_active_requests`
- `Pedidos por minuto`: `sum(rate(app_orders_created_total[1m])) * 60`
- `Goroutines Go`: `go_goroutines`

## Parar a Stack

```bash
docker compose stop
```

Para iniciar novamente:

```bash
docker compose start
```

## Remover a Stack

Remove containers e rede:

```bash
docker compose down
```

Remove também os volumes locais:

```bash
docker compose down -v
```

## Próximo Manual

Depois que a infraestrutura estiver de pé, siga para [Aplicação app-observability-lab: Métricas](aplicacao-go-metricas.md).
