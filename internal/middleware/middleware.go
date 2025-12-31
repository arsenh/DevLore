package middleware

import (
	"net/http"
	"time"

	"github.com/arsenh/DevLore/internal/logger"
	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logger.L().WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
			"took":   time.Since(start),
		}).Info("request handled")
	})
}
