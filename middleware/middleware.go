package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	httpHelper "github.com/zhughes3/website/http"
	"strconv"
	"strings"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			httpHelper.NewErrorResponse(w, http.StatusUnauthorized, "Authentication failure")
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := VerifyToken(tokenString)
		if err != nil {
			httpHelper.NewErrorResponse(w, http.StatusUnauthorized, "Error verifying JWT token: " + err.Error())
			return
		}

		userId := strconv.FormatFloat(claims.(jwt.MapClaims)["user_id"].(float64), 'g', 1, 64)
		r.Header.Set("userId", userId)
		next.ServeHTTP(w, r)
	})

}