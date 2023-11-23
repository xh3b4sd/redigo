package prefix

func WithKeys(prefix string, keys ...string) string {
	newKey := prefix

	for _, k := range keys {
		newKey += ":" + k
	}

	if prefix == "" {
		newKey = newKey[1:]
	}

	return newKey
}
