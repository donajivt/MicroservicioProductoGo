package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleRaw, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"message": "Rol no encontrado"})
			c.Abort()
			return
		}

		role, ok := roleRaw.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "Formato de rol inválido"})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if allowed == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"message": "No tienes permiso para esta operación"})
		c.Abort()
	}
}
