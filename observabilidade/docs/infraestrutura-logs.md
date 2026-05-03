# Infraestrutura de Logs

Este manual prepara o ambiente de infraestrutura para instrumentação e coleta de logs.

A ideia é subir, com Docker Compose, as ferramentas que recebem, armazenam e permitem visualizar os logs gerados pela aplicação.

## Stack

```text
app-observability-lab
  -> Docker logs
  -> Promtail
  -> Loki
  -> Grafana
```

## Ferramentas

**Promtail** coleta logs dos containers Docker e envia para o Loki.

Paralelo: lembra um agente coletor, como **Datadog Agent**, **Fluent Bit** ou **Filebeat**.

**Loki** armazena e consulta logs.

Paralelo: lembra a parte de logs do **Datadog Logs**, **Elasticsearch/ELK** ou **CloudWatch Logs**.

**Grafana** visualiza e explora os logs armazenados no Loki.

Paralelo: lembra dashboards e telas de busca do **Datadog**, **Kibana** ou **CloudWatch Logs Insights**.

**app-observability-lab** é a aplicação de teste que gera logs para o laboratório.

Paralelo: em uma aplicação **Spring Boot**, seria a aplicação escrevendo logs com Logback/SLF4J. O **Spring Actuator** é outro recurso: ele expõe endpoints como `/actuator/health` e `/actuator/metrics`, mas não é o agente coletor de logs.

## Arquivos

- `docker-compose.yml`: define os containers da stack.
- `infra/logs/loki-config.yaml`: configura o Loki.
- `infra/logs/promtail-config.yaml`: configura a coleta de logs dos containers Docker.
- `infra/metrics/grafana-datasources.yaml`: provisiona os Data Sources Loki e Prometheus no Grafana.
- `apps/app-observability-lab/Dockerfile`: gera a imagem da aplicação de teste.

## Arquivos de Configuração

Os arquivos em `infra/logs/` configuram as ferramentas da stack de logs.

### Grafana

Arquivo usado pelo Docker Compose atual:

```text
infra/metrics/grafana-datasources.yaml
```

Define o Loki como Data Source padrão do Grafana e também cadastra o Prometheus.

Pontos importantes:

- `apiVersion`: versão do formato de provisionamento do Grafana.
- `datasources`: lista de fontes de dados criadas automaticamente.
- `name: Loki`: nome exibido no Grafana.
- `type: loki`: informa que o Data Source é do tipo Loki.
- `url: http://loki:3100`: endereço interno do Loki na rede Docker Compose.
- `isDefault: true`: deixa o Loki selecionado como Data Source padrão.

### Loki

Arquivo:

```text
infra/logs/loki-config.yaml
```

Configura como o Loki recebe, indexa e armazena logs.

Pontos importantes:

- `auth_enabled: false`: desativa autenticação para o laboratório local.
- `server.http_listen_port: 3100`: porta HTTP do Loki.
- `common.storage.filesystem`: define armazenamento local em disco.
- `schema_config`: define o formato de indexação dos logs.
- `limits_config`: habilita recursos extras, como metadados estruturados.
- `ruler`: configura suporte a regras e alertas LogQL.

### Promtail

Arquivo:

```text
infra/logs/promtail-config.yaml
```

Configura como o Promtail descobre containers, lê logs e envia para o Loki.

Pontos importantes:

- `server`: porta interna de status/métricas do Promtail.
- `positions`: arquivo onde o Promtail salva até onde já leu os logs.
- `clients`: destino para envio dos logs, neste caso o Loki.
- `scrape_configs`: regras de coleta.
- `docker_sd_configs`: descobre containers Docker automaticamente.
- `pipeline_stages`: interpreta o formato dos logs do Docker.
- `relabel_configs`: transforma metadados Docker em labels como `container`, `stream` e `compose_service`.

## Pré-requisitos

Verifique o Docker:

```bash
docker version
```

Verifique o Docker Compose:

```bash
docker compose version
```

Em ambientes antigos, o comando pode ser `docker-compose`.

## Subir a Stack

Execute a partir da raiz do projeto:

```bash
docker compose up -d --build
```

Com Compose v1:

```bash
docker-compose up -d --build
```

Use `--build` quando quiser recriar a imagem da `app-observability-lab`.

## Verificar Containers

```bash
docker compose ps
```

Serviços esperados:

```text
grafana
loki
promtail
app-observability-lab
```

## Verificar Loki

```bash
curl http://localhost:3100/ready
```

Saída esperada:

```text
ready
```

A URL `http://localhost:3100` pode retornar `404` no navegador. Isso é normal: a raiz do Loki não é uma página web.

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

O Data Source `Loki` já é provisionado automaticamente pelo arquivo `infra/metrics/grafana-datasources.yaml`.

## Data Source Loki

Dentro da rede do Docker Compose, o Grafana acessa o Loki por:

```text
http://loki:3100
```

Não use `http://localhost:3100` dentro do Grafana. Dentro do container do Grafana, `localhost` aponta para o próprio container do Grafana, não para o Loki.

Resumo:

```text
Grafana container -> Loki container: http://loki:3100
Seu navegador/host -> Loki: http://localhost:3100
```

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

Depois que a infraestrutura estiver de pé, siga para [Aplicação app-observability-lab: Logs](aplicacao-go-logs.md).
