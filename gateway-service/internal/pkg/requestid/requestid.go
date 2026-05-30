package requestid

import "context"

const Header = "X-Request-Id"
const MetadataKey = "x-request-id"

type contextKey struct{}

func NewContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextKey{}, requestID)
}

func FromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(contextKey{}).(string)
	return requestID
}
