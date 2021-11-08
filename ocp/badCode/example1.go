package badCode

import (
	"net/http"

	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"
)

type PermissionChecker struct {
	//
	// какие-то поля
	//
}

func (c *PermissionChecker) HasPermission(ctx *gin.Context, name string) bool {
	var permissions []string
	switch ctx.GetString("authType") {
	case "jwt":
		permissions = c.extractPermissionsFromJwt(ctx.Request.Header)
	case "basic":
		permissions = c.getPermissionsForBasicAuth(ctx.Request.Header)
	case "applicationKey":
		permissions = c.getPermissionsForApplicationKey(ctx.Query("applicationKey"))
	}

	var result []string
	linq.From(permissions).
		Where(func(permission interface{}) bool {
			return permission.(string) == name
		}).ToSlice(&result)

	return len(result) > 0
}

func (c *PermissionChecker) getPermissionsForApplicationKey(key string) []string {
	var result []string
	//
	// получаем права доступа из key
	//
	return result
}

func (c *PermissionChecker) getPermissionsForBasicAuth(h http.Header) []string {
	var result []string
	//
	// получаем права доступа из заголовка
	//
	return result
}

func (c *PermissionChecker) extractPermissionsFromJwt(h http.Header) []string {
	var result []string
	//
	// извлекаем права доступа из JWT
	//
	return result
}
