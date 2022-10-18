package clientsecretpolicyfields

type ClientSecretPolicyField int64

const (
	MaxLength ClientSecretPolicyField = iota
	MinDigits
	MinLength
	MinLowerCaseChars
	MinSpecialChars
	MinUpperCaseChars
)

func (s ClientSecretPolicyField) String() string {
	switch s {
	case MaxLength:
		return "max_length"
	case MinDigits:
		return "min_digits"
	case MinLength:
		return "min_length"
	case MinLowerCaseChars:
		return "min_lower_case_chars"
	case MinSpecialChars:
		return "min_special_chars"
	case MinUpperCaseChars:
		return "min_upper_case_chars"
	}
	return "unknown"
}
