# Aplicação app-observability-lab: Logs

Este guia explica como a aplicação `app-observability-lab` gera logs e como consultar esses logs no Grafana.

A infraestrutura completa com Docker Compose, Promtail, Loki e Grafana está descrita em [Infraestrutura de Logs](infraestrutura-logs.md).

## Objetivo da Aplicação

A aplicação fica em:

```text
apps/app-observability-lab
```

Ela usa **Gin** para expor rotas HTTP simples e escreve logs em `stdout`, que é a saída padrão do processo.

Quando a aplicação roda em container, o fluxo fica assim:

```text
Aplicação Go
  -> stdout do container
  -> Docker
  -> Promtail
  -> Loki
  -> Grafana
```

Ou seja: a aplicação não envia logs diretamente para o Loki. Ela apenas escreve logs; Docker e Promtail cuidam da coleta.

## Arquivos Principais

- `apps/app-observability-lab/main.go`: app HTTP com Gin que gera logs.
- `apps/app-observability-lab/go.mod`: dependências da aplicação, incluindo Gin.
- `apps/app-observability-lab/Dockerfile`: empacota a aplicação em imagem Docker.
- `docker-compose.yml`: executa a aplicação junto com Grafana, Loki e Promtail.

Os arquivos em `infra/logs/` não fazem parte da instrumentação da aplicação. Eles configuram a infraestrutura que coleta, armazena e visualiza os logs:

- `infra/logs/promtail-config.yaml`: coleta logs dos containers Docker.
- `infra/logs/loki-config.yaml`: configura o armazenamento dos logs.
- `infra/metrics/grafana-datasources.yaml`: conecta o Grafana ao Loki e ao Prometheus.

## Estrutura do Código

O código foi organizado por pilar de observabilidade. A `main` registra controllers específicos, e cada controller concentra as rotas do seu tema:

- `main`: configura `stdout`, cria o servidor Gin, registra os controllers e inicia a aplicação.
- `controllers/logs.go`: concentra as rotas e funções do pilar de logs.
- `controllers/metrics.go`: concentra as rotas e a instrumentação do pilar de métricas.
- `controllers/tracing.go`: reserva a estrutura do futuro pilar de tracing.
- `registrarLogPorNivel`, `criaLogEstruturado` e `iniciarHeartbeat`: funções usadas pelo controller de logs.

Fluxo principal:

```text
main
  -> log.SetOutput(os.Stdout)
  -> gin.Default
  -> RegistrarRotasDeLogs
  -> RegistrarRotasDeMetricas
  -> RegistrarRotasDeTracing
  -> servidor.Run(":8080")

controllers/logs.go
  -> servidor.GET("/")
  -> servidor.GET("/health")
  -> servidor.GET("/logs")
  -> servidor.GET("/logs/estruturado")
  -> registrarLogPorNivel("info", "aplicacao go iniciada na porta 8080")
  -> heartbeat periodico em debug
```

## Rotas da Aplicação

### Rota raiz

```text
GET /
```

Retorna uma mensagem simples e gera um log `debug`:

```text
level=debug msg="rota raiz acessada"
```

Exemplo:

```bash
curl "http://localhost:8080/"
```

### Health check

```text
GET /health
```

Retorna `ok` e gera um log `info`:

```text
level=info msg="health check recebido"
```

Exemplo:

```bash
curl "http://localhost:8080/health"
```

Paralelo: em Spring Boot, essa ideia lembra o `/actuator/health`. A diferença é que, aqui, criamos uma rota simples com Gin apenas para o laboratório.

### Gerar logs por level

Logs por level classificam a severidade ou intenção da mensagem. Eles ajudam a separar eventos comuns (`info`), detalhes de diagnóstico (`debug`), alertas (`warn`) e falhas (`error`).

```text
GET /logs?level=<nivel>&message=<mensagem>
```

Essa é a rota principal do laboratório. Ela lê `level` e `message` da URL e chama `registrarLogPorNivel`.

A função `registrarLogPorNivel` usa um `switch case` para decidir qual log escrever.

Níveis usados no laboratório:

```text
debug
info
warn
error
```

Exemplos de logs gerados:

```text
level=debug msg="debugando-fluxo"
level=info msg="pedido-criado"
level=warn msg="estoque-baixo"
level=error msg="falha-ao-processar"
```

Se o level enviado não for reconhecido, a aplicação registra como `info` para manter o exemplo simples.

Exemplo:

```bash
curl "http://localhost:8080/logs?level=error&message=falha-ao-processar"
```

Fluxo da chamada:

```text
GET /logs?level=error&message=falha-ao-processar
  -> handler da rota /logs
  -> registrarLogPorNivel
  -> log.Printf("level=error ...")
  -> stdout do container
  -> Docker
  -> Promtail
  -> Loki
  -> Grafana
```

### Gerar log estruturado

Logs estruturados carregam a mensagem junto com campos de contexto. Em ambientes reais, isso ajuda a investigar por serviço, squad, ambiente, conta cloud, região, pedido ou transação.

```text
GET /logs/estruturado
```

Essa rota chama a função `criaLogEstruturado`.

Ela simula um log de pedido em um ambiente real, incluindo contexto como serviço, squad, ambiente, conta AWS e região.

Exemplo:

```bash
curl "http://localhost:8080/logs/estruturado"
```

Exemplo de log gerado:

```json
{
  "time": "2026-05-02T12:00:00Z",
  "level": "INFO",
  "msg": "pedido criado",
  "service": "orders-api",
  "squad": "payments",
  "environment": "production",
  "aws_account_id": "123456789012",
  "region": "us-east-1",
  "order_id": "987"
}
```

Neste laboratório usamos `log/slog`, que já faz parte da biblioteca padrão do Go. Em projetos reais, também é comum encontrar bibliotecas como `zap`, `zerolog` e `logrus`.

Tradeoff importante:

- Campos como `service`, `squad`, `environment` e `region` ajudam a filtrar e agrupar logs.
- Campos muito variáveis, como `order_id`, `request_id` e `user_id`, geralmente não devem virar labels no Loki para evitar alta cardinalidade.
- Mesmo sem virar label, esses campos continuam úteis para busca textual e investigação.

## Subir a Stack

Execute a partir da raiz do projeto:

```bash
docker compose up -d --build
```

Em ambientes com Compose v1:

```bash
docker-compose up -d --build
```

Use `--build` quando quiser garantir que a imagem da `app-observability-lab` seja criada ou atualizada.

Se você já tinha subido a stack antes do rename de `app-go` para `app-observability-lab`, prefira recriar os containers:

```bash
docker compose down
docker compose up -d --build
```

## Gerar Logs

Depois que a stack estiver de pé, gere logs com:

```bash
curl "http://localhost:8080/logs?level=info&message=pedido-criado"
curl "http://localhost:8080/logs?level=debug&message=debugando-fluxo"
curl "http://localhost:8080/logs?level=warn&message=estoque-baixo"
curl "http://localhost:8080/logs?level=error&message=falha-ao-processar"
curl "http://localhost:8080/logs/estruturado"
```

Também é possível testar as outras rotas:

```bash
curl "http://localhost:8080/"
curl "http://localhost:8080/health"
```

Mesmo sem chamar endpoints, a aplicação gera automaticamente um heartbeat a cada 30 segundos:

```text
level=debug msg="heartbeat da aplicacao"
```

## Ver Logs pelo Docker

Antes de consultar no Grafana, você pode conferir se a aplicação está escrevendo logs:

```bash
docker compose logs -f app-observability-lab
```

## Consultar no Grafana

Acesse:

```text
http://localhost:3000
```

Credenciais:

```text
Usuário: admin
Senha: admin
```

No Grafana:

1. Clique em **Explore**.
2. Selecione o Data Source **Loki**.
3. Clique em **Code**.
4. Execute uma consulta LogQL.

Consulta pelos logs da aplicação:

```logql
{compose_service="app-observability-lab"}
```

Consulta por tipo do log:

```logql
{compose_service="app-observability-lab"} |= "level=error"
```

Consulta por mensagem:

```logql
{compose_service="app-observability-lab"} |= "pedido-criado"
```

Consulta pelo log estruturado de pedido:

```logql
{compose_service="app-observability-lab"} |= "orders-api" |= "payments"
```

## Criar Dashboard de Logs

Depois de validar as consultas no **Explore**, você pode salvar uma visão simples em dashboard.

No Grafana:

1. Clique em **Dashboards**.
2. Clique em **New**.
3. Clique em **New dashboard**.
4. Clique em **Add visualization**.
5. Selecione o Data Source **Loki**.
6. Em **Code**, informe a query:

```logql
{compose_service="app-observability-lab"}
```

7. Em **Visualization**, selecione **Logs**.
8. No título do painel, use `Logs da app de teste`.
9. Clique em **Apply**.
10. Clique em **Save dashboard**.

Para criar um painel focado em erros, adicione outra visualização com:

```logql
{compose_service="app-observability-lab"} |= "level=error"
```

Sugestão de painéis para este laboratório:

- `Todos os logs`: `{compose_service="app-observability-lab"}`
- `Logs de erro`: `{compose_service="app-observability-lab"} |= "level=error"`
- `Logs de warning`: `{compose_service="app-observability-lab"} |= "level=warn"`
- `Busca por pedido`: `{compose_service="app-observability-lab"} |= "pedido-criado"`

## Labels Úteis

O Promtail adiciona labels aos logs coletados. Para esta aplicação, as mais úteis são:

```text
container="app-observability-lab"
compose_service="app-observability-lab"
job="docker"
```

Essas labels vêm dos metadados do Docker e são criadas pelo Promtail conforme configurado em `infra/logs/promtail-config.yaml`.

Na prática, prefira começar por:

```logql
{compose_service="app-observability-lab"}
```

## Próximos Passos

Depois dos logs funcionando, podemos criar guias específicos para os próximos pilares:

- `../metrics/aplicacao-go-metricas.md`
- `../traces/aplicacao-go-traces.md`
