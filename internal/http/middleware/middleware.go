package middleware

import (
	tokens "cors/pkg/token"
	"github.com/casbin/casbin/v2"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Middleware(c *gin.Context) {

	allow, err := CheckPermission(c.Request)

	if err != nil {
		if valid, ok := err.(*jwt.ValidationError); ok && valid.Errors == jwt.ValidationErrorExpired {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "token was expired",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
		return
	} else if !allow {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
		return
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	id, _ := tokens.GetIdFromToken(c.Request)
	c.Set("user_id", id)
	email, _ := tokens.GetEmailFromToken(c.Request)
	c.Set("email", email)
	c.Next()
}

func TimingMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start)
	c.Writer.Header().Set("X-Response-Time", duration.String())
}

func CheckPermission(r *http.Request) (bool, error) {
	role, err := GetRole(r)
	if err != nil {
		log.Println("Error while getting role from token: ", err)
		return false, err
	}

	method := r.Method
	path := r.URL.Path

	enforcer, err := casbin.NewEnforcer("auth.conf", "auth.csv")
	if err != nil {
		log.Println(err)
		return false, err
	}
	allowed, err := enforcer.Enforce(role, path, method)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return allowed, nil
}

func GetRole(r *http.Request) (string, error) {
	tokenStr := r.Header.Get("Authorization")

	if tokenStr == "" {
		return "unauthorized", nil
	} else if strings.Contains(tokenStr, "Basic") {
		return "unauthorized", nil
	}

	claims, err := tokens.ExtractClaim(tokenStr)
	if err != nil {
		log.Println("Error while extracting claims: ", err)
		return "unauthorized", err
	}

	return claims["role"].(string), nil
}
