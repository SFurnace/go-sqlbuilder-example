package log

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	basic "pers.drcz/test/regex/basic"
)

/* store log fields in context */

type zapFieldsKey struct{}

var zapFields zapFieldsKey // log fields store in context

func PutLogFields(ctx context.Context, fields ...zap.Field) {
	_ = basic.StoreVal(ctx, zapFields, append(LogFields(ctx), fields...))
}

func LogFields(ctx context.Context) []zap.Field {
	v, _ := basic.GetVal(ctx, zapFields).([]zap.Field)
	return v
}

func LogFields2Interfaces(ctx context.Context) []interface{} {
	fs := LogFields(ctx)
	result := make([]interface{}, 0, len(fs))
	for i := range fs {
		result = append(result, fs[i])
	}
	return result
}

/* Zap Helpers */

const (
	SessionId = "session-id"
)

func BackgroundCtxWithRandomId() context.Context {
	ctx := basic.Background()
	PutRandomId(ctx)
	return ctx
}

func PutId(ctx context.Context, id string) {
	PutLogFields(ctx, zap.String(SessionId, id))
}

func PutRandomId(ctx context.Context) {
	v, _ := uuid.NewRandom()
	PutLogFields(ctx, zap.String(SessionId, v.String()))
}

func RetrieveSessionId(ctx context.Context) string {
	fs := LogFields(ctx)
	for i := range fs {
		if fs[i].Key == SessionId {
			return fs[i].String
		}
	}
	return ""
}
