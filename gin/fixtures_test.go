package gin

func fixtureMetadata() map[string]string {
	return map[string]string{
		"TestKey1":  "test_key_1_success",
		"Test-Key2": "test-key-2-success",
		"Test_key3": "TestKey3Success",
	}
}
