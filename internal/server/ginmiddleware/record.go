package ginmiddleware

import (
	"fmt"
	"github.com/MR5356/aurora/internal/domain/system"
	"github.com/MR5356/aurora/internal/domain/user"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Record() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		start := time.Now()

		defer func() {
			cost := time.Since(start).Microseconds()
			httpCode := ctx.Writer.Status()
			clientIP := ctx.ClientIP()
			clientUA := ctx.Request.UserAgent()
			method := ctx.Request.Method

			entry := logrus.WithFields(logrus.Fields{
				"cost":   cost,
				"method": method,
				"code":   httpCode,
				"ip":     clientIP,
				"path":   path,
			})

			if len(ctx.Errors) > 0 {
				entry.Error(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
			} else if path != "/api/v1/metrics" {
				msg := fmt.Sprintf("user-agent: %s", clientUA)
				switch {
				case httpCode >= http.StatusInternalServerError:
					entry.Error(msg)
				case httpCode >= http.StatusBadRequest:
					entry.Warn(msg)
				default:
					entry.Debug(msg)
				}
			}

			token := ginutil.GetToken(ctx)

			var uid string
			u, err := user.GetJWTService().ParseToken(token)
			if err == nil {
				uid = u.ID
			}

			err = system.GetService().InsertRecord(&system.Record{
				UserID:    uid,
				Path:      path,
				Method:    method,
				Code:      strconv.Itoa(httpCode),
				ClientIP:  clientIP,
				UserAgent: clientUA,
				Cost:      cost,
				IsApi:     strings.HasPrefix(path, "/api/"),
			})

			if err != nil {
				logrus.Errorf("insert record failed, error: %v", err)
			}
		}()

		ctx.Next()
	}
}
