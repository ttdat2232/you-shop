package background

import "github.com/rs/zerolog"

type JobFunc func()

func Go(logger zerolog.Logger, fn JobFunc) {
	logger = logger.With().Str("Background", "Job").Logger()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error().Interface("recover", r).Msg("panic recovered")
			}
		}()
		fn()
	}()
}
