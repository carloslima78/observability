# Aplicação app-observability-lab: Tracing

Este guia será usado no futuro módulo de tracing.

A aplicação já possui um controller reservado para este pilar:

```text
app/controllers/tracing.go
```

Rota provisória:

```text
GET /tracing/demo
```

Exemplo:

```bash
curl "http://localhost:8080/tracing/demo"
```

Por enquanto, essa rota apenas confirma que a estrutura do pilar de tracing existe.

Em uma próxima etapa, este manual deve cobrir:

- criação de traces
- criação de spans
- propagação de contexto
- integração com OpenTelemetry
- visualização em Tempo ou Jaeger
