package ctx

import "context"

// KeyRequestID is used to uniquely reference each request
const KeyRequestID string = "requestID"

// RequestID gets the requestID from the context
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(KeyRequestID).(string)
	return requestID
}

// SetRequestID sets the requestID in the context
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, KeyRequestID, requestID)
}
