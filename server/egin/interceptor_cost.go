package egin

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/transport"
	"github.com/spf13/cast"
)

// XResCostTimer wrap gin response writer add start time
type XResCostTimer struct {
	gin.ResponseWriter
	start           time.Time
	ginCtx          *gin.Context
	enableHeaderApp bool
	prefixKey       string
}

// WriteHeader 如果写入header，需要这么处理
// ctx.Request = ctx.Request.WithContext(context)
func (w *XResCostTimer) WriteHeader(statusCode int) {
	// header必须在c.json响应。
	cost := float64(time.Since(w.start).Microseconds()) / 1000
	if w.enableHeaderApp {
		errStr := w.ginCtx.Errors.ByType(gin.ErrorTypePrivate).String()
		if errStr != "" {
			if econf.GetBool(w.prefixKey + ".enableResHeaderError") {
				w.Header().Set(eapp.EgoHeaderExpose()+"time", strconv.FormatFloat(cost, 'f', -1, 64)+"|"+errStr+"|"+eapp.Name())
			} else {
				w.Header().Set(eapp.EgoHeaderExpose()+"time", strconv.FormatFloat(cost, 'f', -1, 64)+"|err|"+eapp.Name())
			}
		} else {
			w.Header().Set(eapp.EgoHeaderExpose()+"time", strconv.FormatFloat(cost, 'f', -1, 64)+"|ok|"+eapp.Name())
		}
	} else {
		errStr := w.ginCtx.Errors.ByType(gin.ErrorTypePrivate).String()
		if errStr != "" {
			if econf.GetBool(w.prefixKey + ".enableResHeaderError") {
				w.Header().Set(eapp.EgoHeaderExpose()+"time", strconv.FormatFloat(cost, 'f', -1, 64)+"|"+errStr)
			} else {
				w.Header().Set(eapp.EgoHeaderExpose()+"time", strconv.FormatFloat(cost, 'f', -1, 64)+"|err")
			}
		} else {
			w.Header().Set(eapp.EgoHeaderExpose()+"time", strconv.FormatFloat(cost, 'f', -1, 64)+"|ok")
		}
	}
	for _, key := range transport.CustomHeaderKeys() {
		if value := cast.ToString(w.ginCtx.Request.Context().Value(key)); value != "" {
			// x-expose 需要在这里获取
			if strings.HasPrefix(key, eapp.EgoHeaderExpose()) {
				// 设置到ctx response header
				w.Header().Set(key, value)
			}
		}
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *XResCostTimer) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

// NewXResCostTimer middleware to add X-Res-Cost-Time
func NewXResCostTimer(prefixKey string, enableHeaderApp bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &XResCostTimer{
			ResponseWriter:  c.Writer,
			start:           time.Now(),
			ginCtx:          c,
			enableHeaderApp: enableHeaderApp,
			prefixKey:       prefixKey,
		}
		c.Writer = blw
		c.Next()
	}
}
