package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/coffemanfp/chat/config"
	"github.com/golang-jwt/jwt"
	muxhandlers "github.com/gorilla/handlers"
)

func logginMiddleware(next http.Handler) http.Handler {
	return muxhandlers.LoggingHandler(os.Stdout, next)
}

func verifyJWTMiddleware(conf config.ConfigInfo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return authHandler{next, conf}
	}
}

type authHandler struct {
	h    http.Handler
	conf config.ConfigInfo
}

func (a authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header["Authorization"] != nil {
		tokenString := r.Header["Authorization"][0][7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(a.conf.Server.SecretKey), nil
		})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You're Unauthorized due to error parsing the JWT"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "id", int(claims["id"].(float64)))
			a.h.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("You're Unauthorized due to No token in the header"))
		if err != nil {
			return
		}
	}
}

// func verifyJWT(endpointHandler func(writer http.ResponseWriter, r *http.Request)) http.HandlerFunc {
// 	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
// 		if request.Header["Authorization"] != nil {
// 			token, err := jwt.Parse(request.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
// 				_, ok := token.Method.(*jwt.SigningMethodECDSA)
// 				if !ok {
// 					writer.WriteHeader(http.StatusUnauthorized)
// 					_, err := writer.Write([]byte("You're Unauthorized"))
// 					if err != nil {
// 						return nil, err
// 					}
// 				}
// 				return "", nil

// 			})
// 			// parsing errors result
// 			if err != nil {
// 				writer.WriteHeader(http.StatusUnauthorized)
// 				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
// 				if err2 != nil {
// 					return
// 				}

// 			}
// 			// if there's a token
// 			if token.Valid {
// 				endpointHandler(writer, request)
// 			} else {
// 				writer.WriteHeader(http.StatusUnauthorized)
// 				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
// 				if err != nil {
// 					return
// 				}
// 			}
// 		} else {
// 			writer.WriteHeader(http.StatusUnauthorized)
// 			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
// 			if err != nil {
// 				return
// 			}
// 		}
// 		// response for if there's no token header
// 	})
// }
