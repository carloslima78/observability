package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegistrarRotasDeTracing reserva as rotas do futuro pilar de tracing.
func RegistrarRotasDeTracing(servidor *gin.Engine) {
	servidor.GET("/tracing/demo", func(contexto *gin.Context) {
		contexto.JSON(http.StatusOK, gin.H{
			"mensagem": "pilar de tracing sera implementado em uma proxima etapa",
		})
	})
}
