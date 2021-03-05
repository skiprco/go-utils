package validation

const errorDomain = "go_utils"
const errorSubDomain = "validation"

// ErrorInvalidCountryCode provided country code is invalid.
const ErrorInvalidCountryCode = "invalid_country_code"

// ErrorNotAPhoneNumber indicates the provided phone number is not recognised as one.
const ErrorNotAPhoneNumber = "not_a_phone_number"

// ErrorInvalidPhoneNumber indicates the provided phone number has the correct format,
// but is symantically incorrect.
const ErrorInvalidPhoneNumber = "invalid_phone_number"

// ErrorNotAMobilePhoneNumber indicates the provided phone number is not a mobile number.
const ErrorNotAMobilePhoneNumber = "not_a_mobile_phone_number"

// ErrorEndTimeBeforeStartTime indicates the provided end time is before
// the provided start time. End time should be equal to or after start time.
const ErrorEndTimeBeforeStartTime = "end_time_before_start_time"
