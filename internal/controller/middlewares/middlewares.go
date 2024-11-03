package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zde37/instashop-task/internal/models"
	"github.com/zde37/instashop-task/pkg"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
)

// Auth middleware authenticates user requests using JWT
func Auth(token *pkg.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    "MISSING_AUTH_HEADER",
				Message: "Authorization header is required",
			})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    "INVALID_AUTH_HEADER",
				Message: "Invalid authorization header format. Use 'Bearer <token>'",
			})
			return
		}

		accessToken := fields[1]
		claims, err := token.VerifyToken(accessToken)
		if err != nil {
			var errResp models.ErrorResponse
			switch {
			case strings.Contains(err.Error(), "expired"):
				errResp = models.ErrorResponse{
					Code:    "TOKEN_EXPIRED",
					Message: "Access token has expired",
				}
			case strings.Contains(err.Error(), "invalid"):
				errResp = models.ErrorResponse{
					Code:    "INVALID_TOKEN",
					Message: "Invalid or malformed access token",
				}
			default:
				errResp = models.ErrorResponse{
					Code:    "AUTH_FAILED",
					Message: "Authentication failed",
					Details: err.Error(),
				}
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResp)
			return
		}

		// Store user information in context
		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.UserEmail)
		ctx.Set("user_role", claims.UserRole)

		ctx.Next()
	}
}

// AdminRequired middleware ensures the user has admin role
func AdminRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("user_role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    "NO_USER_ROLE",
				Message: "User role not found in context",
				Details: "Authentication middleware must be applied before this middleware",
			})
			return
		}

		if userRole != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
				Code:    "ADMIN_REQUIRED",
				Message: "This operation requires admin privileges",
			})
			return
		}

		ctx.Next()
	}
}
