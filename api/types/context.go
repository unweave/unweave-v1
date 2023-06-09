package types

// Context Keys are used to store and retrieve values from the context. These can be
// used for adding contextual information to logs, or for passing values between
// middleware. If using in middleware, you should be careful to only use keys in the until
// the last handler in the chain *at-most* and not rely on the context being passed
// further down the stack.
const (
	UserIDCtxKey     = "userID"
	AccountIDCtxKey  = "accountID"
	BuildIDCtxKey    = "buildID"
	ProjectIDCtxKey  = "projectID"
	ExecIDCtxKey     = "execID"
	ExecStatusCtxKey = "execStatus"
	ObserverCtxKey   = "observer"
)
