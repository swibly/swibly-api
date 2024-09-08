package middleware

import "github.com/gin-gonic/gin"

// Why? It solves a very specific problem where responses should not be cached,
// ensuring that each request receives fresh data based on its parameters.
// Caused due to cloud caching headers or smt like that idk
//
// HACK: Please, future me, fix this when you got time, it's kinda hacky
func DisableCache(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "1")
	ctx.Next()
}
