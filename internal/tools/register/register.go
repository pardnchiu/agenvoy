package toolRegister

import (
	"context"
	"encoding/json"
	"fmt"

	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

type Handler func(ctx context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error)

type GroupHandler func(ctx context.Context, e *toolTypes.Executor, name string, args json.RawMessage) (string, error)

var handlerMap = map[string]Handler{}

func Register(name string, h Handler) {
	handlerMap[name] = h
}

func Dispatch(ctx context.Context, e *toolTypes.Executor, name string, args json.RawMessage) (string, error) {
	handler, ok := handlerMap[name]
	if !ok {
		return "", fmt.Errorf("not exist: %s", name)
	}
	return handler(ctx, e, args)
}
