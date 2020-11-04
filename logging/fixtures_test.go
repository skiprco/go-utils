package logging

import "github.com/skiprco/go-utils/v2/metadata"

func fixtureMetadata() metadata.Metadata {
	return metadata.Metadata{
		"test_service_test_key": "after_value",
		"should_not_be_touched": "dont_touch_me",
	}
}
