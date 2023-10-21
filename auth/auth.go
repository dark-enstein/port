package auth

import (
	"github.com/dark-enstein/gauth"
	"time"
)

func NewJWT(id, subject, audience, issuer string, dur time.Duration, signedbyte []byte) (gauth.JWTStrand, error) {
	jReq := gauth.UserJWTRequest{
		ID:            id,
		Subject:       subject,
		Audience:      audience,
		ExpirationDur: &dur,
		Issuer:        issuer,
		SignedByte:    signedbyte,
	}
	gauthJwt, err := gauth.NewJWT(&jReq)
	if err != nil {
		return "", err
	}
	return gauthJwt.String(), nil
}

func DecodeJWT(tokenString string, secret []byte) (*gauth.JWT, error) {
	tkn := gauth.JWTStrand(tokenString)
	decode, err := tkn.Decode(secret)
	if err != nil {
		return nil, err
	}
	return decode, nil
}

//
//func authMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		tokenString := r.Header.Get("Authorization")
//		if tokenString == "" {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
//
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("unexpected signing method")
//			}
//			return []byte("secret"), nil
//		})
//
//		if err != nil {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		if !token.Valid {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		claims, ok := token.Claims.(jwt.MapClaims)
//		if !ok {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		userID, ok := claims["user_id"].(string)
//		if !ok {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), "user_id", userID)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
