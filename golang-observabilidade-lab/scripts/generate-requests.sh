#!/bin/bash

# Encerra o script quando ocorrer erro, variavel nao definida ou falha em pipeline.
# Isso evita que o teste continue gerando dados inconsistentes sem percebermos.
set -euo pipefail

# Endereco base da aplicacao.
# Pode ser sobrescrito para testar outro host, porta ou ambiente.
# Exemplo: BASE_URL="http://observabilityindeep.local"
BASE_URL="${BASE_URL:-http://localhost:8080}"

# Codigos HTTP usados para alimentar as metricas de erro.
# Use espacos para informar mais de um status.
# Exemplo: ERROR_CODES="500 503 404"
ERROR_CODES="${ERROR_CODES:-500 400 404}"

# Intervalo entre cada ciclo de requisicoes.
# Valores menores geram mais trafego; para estudo local, mantenha algo simples.
SLEEP_SECONDS="${SLEEP_SECONDS:-1}"

while true; do
  # Cada codigo gera uma resposta HTTP diferente na rota didatica de erro.
  # Isso permite comparar counters e percentuais por status no Prometheus/Grafana.
  for error_code in $ERROR_CODES; do
    # Golden Signal: Errors.
    # Alimenta app_http_errors_total e app_http_requests_total com o status escolhido.
    curl -s -o /dev/null -X GET "$BASE_URL/metrics/error/$error_code"
    echo "Requisicao enviada para: $BASE_URL/metrics/error/$error_code"

    # Golden Signal: Latency.
    # Alimenta o histogram app_http_request_duration_seconds com atraso artificial.
    curl -s -o /dev/null -X GET "$BASE_URL/metrics/latency?seconds=0.3"
    echo "Requisicao enviada para: $BASE_URL/metrics/latency?seconds=0.3"

    # Golden Signal: Traffic e metrica de negocio.
    # Alimenta app_http_requests_total e app_orders_created_total.
    curl -s -o /dev/null -X POST "$BASE_URL/metrics/orders"
    echo "Requisicao enviada para: $BASE_URL/metrics/orders"

    # Controla a frequencia do teste para nao transformar o laboratorio local em carga excessiva.
    sleep "$SLEEP_SECONDS"
  done
done
