package logging

import "github.com/skiprco/go-utils/v2/metadata"

func fixtureMetadata() metadata.Metadata {
	return metadata.Metadata{
		"service_name":          "srv-test",
		"srv_test_test_key":     "before_value",
		"should_not_be_touched": "dont_touch_me",
	}
}
