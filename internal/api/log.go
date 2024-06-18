package api

func WithLogging() ApiOption {
	return newFuncApiOption(func(a *Api) {
		a.loggingEnabled = true
	})
}
