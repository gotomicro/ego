package egorm

import (
	"context"
	"fmt"
	"github.com/gotomicro/ego/core/metric"
	"strconv"
	"time"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/trace"
	"github.com/gotomicro/ego/core/util/xcolor"
)

// Handler ...
type Handler func(*Scope)

// Interceptor ...
type Interceptor func(*DSN, string, *Config, *elog.Component) func(next Handler) Handler

func debugInterceptor(dsn *DSN, op string, options *Config, logger *elog.Component) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			fmt.Printf("%-50s[%s] => %s\n", xcolor.Green(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), xcolor.Green("Send: "+logSQL(scope.SQL, scope.SQLVars, true)))
			next(scope)
			if scope.HasError() {
				fmt.Printf("%-50s[%s] => %s\n", xcolor.Red(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), xcolor.Red("Erro: "+scope.DB().Error.Error()))
			} else {
				fmt.Printf("%-50s[%s] => %s\n", xcolor.Green(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), xcolor.Green("Affected: "+strconv.Itoa(int(scope.DB().RowsAffected))))
			}
		}
	}
}

func metricInterceptor(dsn *DSN, op string, options *Config, logger *elog.Component) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			beg := time.Now()
			next(scope)
			cost := time.Since(beg)

			// error metric
			if scope.HasError() {
				metric.LibHandleCounter.WithLabelValues(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr, "ERR").Inc()
				// todo sql语句，需要转换成脱密状态才能记录到日志
				if scope.DB().Error != ErrRecordNotFound {
					logger.Error("mysql err", elog.FieldErr(scope.DB().Error), elog.FieldName(dsn.DBName+"."+scope.TableName()), elog.FieldMethod(op))
				} else {
					logger.Warn("record not found", elog.FieldErr(scope.DB().Error), elog.FieldName(dsn.DBName+"."+scope.TableName()), elog.FieldMethod(op))
				}
			} else {
				metric.LibHandleCounter.Inc(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr, "OK")
			}

			metric.LibHandleHistogram.WithLabelValues(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr).Observe(cost.Seconds())

			if options.SlowThreshold > time.Duration(0) && options.SlowThreshold < cost {
				logger.Error(
					"slow",
					elog.FieldErr(errSlowCommand),
					elog.FieldMethod(op),
					elog.FieldExtMessage(logSQL(scope.SQL, scope.SQLVars, options.DetailSQL)),
					elog.FieldAddr(dsn.Addr),
					elog.FieldName(dsn.DBName+"."+scope.TableName()),
					elog.FieldCost(cost),
				)
			}
		}
	}
}

func logSQL(sql string, args []interface{}, containArgs bool) string {
	if containArgs {
		return bindSQL(sql, args)
	}
	return sql
}

func traceInterceptor(dsn *DSN, op string, options *Config, logger *elog.Component) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			if val, ok := scope.Get("_context"); ok {
				if ctx, ok := val.(context.Context); ok {
					span, _ := trace.StartSpanFromContext(
						ctx,
						"GORM", // TODO this op value is op or GORM
						trace.TagComponent("mysql"),
						trace.TagSpanKind("client"),
					)
					defer span.Finish()

					// 延迟执行 scope.CombinedConditionSql() 避免sqlVar被重复追加
					next(scope)

					span.SetTag("sql.inner", dsn.DBName)
					span.SetTag("sql.addr", dsn.Addr)
					span.SetTag("span.kind", "client")
					span.SetTag("peer.service", "mysql")
					span.SetTag("db.instance", dsn.DBName)
					span.SetTag("peer.address", dsn.Addr)
					span.SetTag("peer.statement", logSQL(scope.SQL, scope.SQLVars, options.DetailSQL))
					return
				}
			}

			next(scope)
		}
	}
}
