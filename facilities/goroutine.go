package facilities

import (
	"context"
	"github.com/heron-sense/gadk/logger"
)

// GoWithRecover wraps a `go func()` with recover()
func GoWithRecover(ctx context.Context, handler func(), recoverHandler func(r interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if recoverHandler != nil {
					go func() {
						defer func() {
							if p := recover(); p != nil {
								logger.Error("GoWithRecover call recoverHandler recover panic")
							}
						}()
						recoverHandler(r)
					}()
				}
			}
		}()
		handler()
	}()
}
