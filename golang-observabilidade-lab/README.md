# Laboratório Go de Observabilidade

Este repositório reúne anotações e exemplos práticos para estudar observabilidade.

## Pilares da Observabilidade

### Logs

Logs são registros de eventos que acontecem dentro de uma aplicação.

Podemos pensar em um log como um "carimbo de tempo" que descreve algo relevante que ocorreu no sistema: uma requisição recebida, uma transação concluída, uma falha ao acessar um serviço externo ou uma decisão tomada pelo fluxo da aplicação.

Sem logs, diagnosticar problemas vira tentativa e erro. Com bons logs, conseguimos entender o comportamento do sistema e investigar a causa raiz de falhas que impactam usuários.

#### Níveis de Log

**Info** registra operações esperadas e bem-sucedidas.

Exemplo:

```text
level=info msg="pedido criado com sucesso"
```

**Error** registra falhas pontuais que não necessariamente param o sistema.

Exemplo:

```text
level=error msg="falha ao processar pagamento"
```

**Debug** registra detalhes úteis durante desenvolvimento, testes ou investigação.

Exemplo:

```text
level=debug msg="payload recebido para validacao"
```

**Critical** registra falhas graves em que o sistema perde capacidade de processar ou reagir corretamente, exigindo intervenção ou estratégias como fallback e circuit breaker.

Exemplo:

```text
level=critical msg="servico principal indisponivel"
```

#### Estrutura de um Bom Log

Além da mensagem principal, um log útil deve conter contexto para diagnóstico.

Campos comuns:

- timestamp
- level
- message
- URI ou rota da transação
- node, pod ou container
- ID do cliente
- ID do registro ou da transação

Exemplo:

```text
timestamp=2026-04-26T13:00:00Z level=info uri=/orders client_id=123 order_id=987 msg="pedido criado"
```

#### Por que Logar

Logs são a base para enxergar o que a aplicação está fazendo.

Eles ajudam a responder perguntas como:

- o que aconteceu?
- quando aconteceu?
- em qual fluxo aconteceu?
- qual cliente ou registro foi afetado?
- o erro foi pontual ou recorrente?

No laboratório atual, a aplicação Go escreve logs em `stdout`/`stderr`, o Docker armazena esses logs, o Promtail coleta, o Loki armazena e o Grafana permite consultar.

### Métricas

Métricas registram valores numéricos ao longo do tempo, como volume de requisições, taxa de erros, latência e estado atual da aplicação.

No laboratório atual, a aplicação Go expõe métricas em `/metrics`, o Prometheus coleta esses dados e o Grafana permite consultar e montar dashboards.

Fluxo:

```text
Aplicação
  -> /metrics
  -> Prometheus
  -> Grafana
```

### Traces

Traces serão estudados em uma etapa futura.

Fluxo esperado:

```text
Aplicação
  -> OpenTelemetry
  -> Tempo ou Jaeger
  -> Grafana
```

## Guias Práticos

### Logs

- [Infraestrutura de Logs](docs/logs/infraestrutura-logs.md)
- [Aplicação app-observability-lab: Logs](docs/logs/aplicacao-go-logs.md)

### Métricas

- [Infraestrutura de Métricas](docs/metrics/infraestrutura-metricas.md)
- [Aplicação app-observability-lab: Métricas](docs/metrics/aplicacao-go-metricas.md)

### Traces

- Infraestrutura de Traces: será criada em uma etapa futura.
- [Aplicação app-observability-lab: Tracing](docs/traces/aplicacao-go-traces.md)
