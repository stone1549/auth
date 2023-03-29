package service

import "encoding/json"

// User holds information on a user.
type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	UserProfile
}

type Gender string

const (
	MALE      Gender = "male"
	FEMALE    Gender = "female"
	NONBINARY Gender = "non-binary"
	OTHER     Gender = "other"
)

func ParseGender(genderStr string) (Gender, error) {
	switch genderStr {
	case MALE.String():
		return MALE, nil
	case FEMALE.String():
		return FEMALE, nil
	case NONBINARY.String():
		return NONBINARY, nil
	case OTHER.String():
		return OTHER, nil
	default:
		return "", ErrorResponse{Status: 400, Message: "invalid gender argument"}
	}
}

//goland:noinspection GoMixedReceiverTypes
func (g Gender) String() string {
	return string(g)
}

// MarshalJSON must be a *value receiver* to ensure that a Suit on a parent object
// does not have to be a pointer in order to have it correctly marshaled.
//
//goland:noinspection GoMixedReceiverTypes
func (g Gender) MarshalJSON() ([]byte, error) {
	// It is assumed Suit implements fmt.Stringer.
	return json.Marshal(g.String())
}

// UnmarshalJSON must be a *pointer receiver* to ensure that the indirect from the
// parsed value can be set on the unmarshaling object. This means that the
// ParseSuit function must return a *value* and not a pointer.
//
//goland:noinspection GoMixedReceiverTypes
func (g *Gender) UnmarshalJSON(data []byte) (err error) {
	var gender string
	if err := json.Unmarshal(data, &gender); err != nil {
		return err
	}
	if *g, err = ParseGender(gender); err != nil {
		return err
	}
	return nil
}

type UserProfile struct {
	Gender Gender   `json:"gender"`
	Age    int      `json:"age"`
	Topics []string `json:"topics"`
}
