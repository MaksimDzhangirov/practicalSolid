package moreExamples

import (
	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"
)

type PermissionProvider interface {
	Type() string
	GetPermissions(ctx *gin.Context) []string
}

type PermissionChecker struct {
	//
	// какие-то поля
	//
}

func (c *PermissionChecker) HasPermission(ctx *gin.Context, provider PermissionProvider, name string) bool {
	permissions := provider.GetPermissions(ctx)

	var result []string
	linq.From(permissions).
			Where(func(permission interface{}) bool {
				return permission.(string) == name
			}).ToSlice(&result)

	return len(result) > 0
}