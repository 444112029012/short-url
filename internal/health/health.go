package health

// Status returns a fixed OK string for CI smoke / health checks.
func Status() string {
	return "ok"
}
