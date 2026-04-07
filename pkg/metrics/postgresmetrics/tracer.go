package postgresmetrics

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type traceStartKey struct{}

type queryTraceData struct {
	start     time.Time
	queryType string
}

type MetricsTracer struct{}

func NewMetricsTracer() *MetricsTracer {
	return &MetricsTracer{}
}

func (t *MetricsTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	qtype := classifyQuery(data.SQL)

	return context.WithValue(ctx, traceStartKey{}, queryTraceData{
		start:     time.Now(),
		queryType: qtype,
	})
}

func (t *MetricsTracer) TraceQueryEnd(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	v := ctx.Value(traceStartKey{})
	if v == nil {
		return
	}

	traceData, ok := v.(queryTraceData)
	if !ok {
		return
	}

	DBQueryTotal.WithLabelValues(traceData.queryType).Inc()
	DBQueryDuration.WithLabelValues(traceData.queryType).
		Observe(time.Since(traceData.start).Seconds())

	if data.Err != nil {
		DBQueryErrorsTotal.WithLabelValues(traceData.queryType).Inc()
	}
}

func classifyQuery(sql string) string {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return "unknown"
	}

	fields := strings.Fields(sql)
	if len(fields) == 0 {
		return "unknown"
	}

	switch strings.ToUpper(fields[0]) {
	case "SELECT":
		return "select"
	case "INSERT":
		return "insert"
	case "UPDATE":
		return "update"
	case "DELETE":
		return "delete"
	case "WITH":
		return "with"
	default:
		return "other"
	}
}
