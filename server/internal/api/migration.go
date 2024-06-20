package api

func WithMigrator(migrator Migrator) ApiOption {
	return newFuncApiOption(func(a *Api) {
		a.migrator = migrator
	})
}
