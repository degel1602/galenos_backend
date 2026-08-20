package httpadapter

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed estampado_default.png
var estampadoDefault embed.FS

// Estampado maneja GET /api/v1/firmaperu/estampado.png: sirve la imagen PNG
// por defecto del estampado que el Firmador local descarga (imageToStamp)
// cuando el frontend no envía una imagen propia.
func (h *FirmaPeruHandler) Estampado(c *gin.Context) {
	img, err := estampadoDefault.ReadFile("estampado_default.png")
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/plain", []byte("estampado no disponible"))
		return
	}
	c.Data(http.StatusOK, "image/png", img)
}

// EstampadoUuid maneja GET /api/v1/firmaperu/estampado/:uuid: sirve la imagen
// PNG del estampado (logo del hospital) enviada por el frontend para ese
// proceso de firma. Sin imagen propia responde 404.
func (h *FirmaPeruHandler) EstampadoUuid(c *gin.Context) {
	img, err := h.service.ObtenerImagenEstampado(c.Request.Context(), c.Param("uuid"))
	if err != nil || len(img) == 0 {
		c.Data(http.StatusNotFound, "text/plain", []byte("estampado no encontrado"))
		return
	}
	c.Data(http.StatusOK, "image/png", img)
}
