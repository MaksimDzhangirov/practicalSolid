package goodCode

import (
	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"
)

type PermissionProvider interface {
	Type() string
	GetPermissions(ctx *gin.Context) []string
}

type PermissionChecker struct {
	providers []PermissionProvider
	//
	// какие-то поля
	//
}

func (c *PermissionChecker) HasPermission(ctx *gin.Context, name string) bool {
	var permissions []string
	for _, provider := range c.providers {
		if ctx.GetString("authType") != provider.Type() {
			continue
		}

		permissions = provider.GetPermissions(ctx)
		break
	}

	var result []string
	linq.From(permissions).
			Where(func(permission interface{}) bool {
				return permission.(string) == name
			}).ToSlice(&result)

	return len(result) > 0
}