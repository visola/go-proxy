package listener

// Initialize loads all listeners from the pesisted state ensuring at least one exists
func Initialize() error {
	if loadErr := LoadFromPersistedState(); loadErr != nil {
		return loadErr
	}

	return ensureListener()
}
