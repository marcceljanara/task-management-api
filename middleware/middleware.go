package middleware

import (
	"context"
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/helper"
	"marcceljanara/task-management-api/service"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

func ValidateJWT(jwtService service.JWTService, next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		tokenString, err := getTokenFromRequest(request)
		if err != nil {
			_ = helper.WriteErrorResponse(writer, err)
			return
		}

		userID, err := jwtService.ValidateToken(tokenString)
		if err != nil || strings.TrimSpace(userID) == "" {
			_ = helper.WriteErrorResponse(writer, exception.New(exception.ErrUnauthorized, "invalid or expired token"))
			return
		}

		ctx := context.WithValue(request.Context(), userIDContextKey, userID)
		next(writer, request.WithContext(ctx), params)
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return "", false
	}
	return userID, true
}

func getTokenFromRequest(request *http.Request) (string, error) {
	cookie, err := request.Cookie("access_token")
	if err == nil && strings.TrimSpace(cookie.Value) != "" {
		return cookie.Value, nil
	}

	authorization := strings.TrimSpace(request.Header.Get("Authorization"))
	if authorization == "" {
		return "", exception.New(exception.ErrUnauthorized, "access token is required")
	}

	fields := strings.Fields(authorization)
	if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
		return fields[1], nil
	}

	return authorization, nil
}
