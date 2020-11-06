package metadata

func fixtureMetadata() Metadata {
	return Metadata{
		"PascalKey": "PascalValue",
		"snake_key": "snake_value",
		"kebab-key": "kebab-value",
	}
}
