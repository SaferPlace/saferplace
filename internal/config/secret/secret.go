package secret

import "log/slog"

// Secret hides secrets which do not want to log by accident.
type Secret string

// LogValue implements slog.LogValuer.
// It avoids revealing the token.
func (Secret) LogValue() slog.Value {
	return slog.StringValue("<redacted>")
}
