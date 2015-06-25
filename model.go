package flash2

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

type mErrors struct {
	Errors modelErrors `json:"errors"`
}

type modelErrors map[string][]string

// ModelBase structure for base model with error valiadtions
//	type User struct {
// 			ID   int64
// 	    Name string
// 	    flash2.ModelBase
// 	}
type ModelBase struct {
	Errors mErrors `sql:"-" json:"-"`
}

// ResetErrors clean all model errors
func (m *ModelBase) ResetErrors() {
	m.Errors = mErrors{Errors: make(modelErrors)}
}

// AddError adding error to model
func (m *ModelBase) AddError(f string, t string) {
	if m.IsValid() {
		m.Errors = mErrors{Errors: make(modelErrors)}
	}
	m.Errors.Errors[f] = append(m.Errors.Errors[f], t)
}

// IsValid returns true if no errors on model
func (m *ModelBase) IsValid() bool {
	return len(m.Errors.Errors) == 0
}

// Valid placeholder for validation function
//  func (u *User) Valid() bool {
//  	u.ValidatePresence("Name", u.Name)
//  	return u.IsValid()
//  }
func (m *ModelBase) Valid() bool {
	return m.IsValid()
}

// GetErrors returns model errors
func (m *ModelBase) GetErrors() modelErrors {
	return m.Errors.Errors
}

// SetErrors set model errors
func (m *ModelBase) SetErrors(e modelErrors) {
	m.Errors.Errors = e
}

// ValidatePresence validates string for presence
// 	m.ValidatePresence("Name", m.Name)
func (m *ModelBase) ValidatePresence(f, v string) {
	if utf8.RuneCountInString(v) == 0 {
		m.AddError(f, "can't be blank")
	}
}

// ValidateLength validates string min, max length. -1 for any
// 	m.ValidateLength("password", m.Password, 6, 18) // min 6, max 18
func (m *ModelBase) ValidateLength(f, v string, min, max int) {
	if min > 0 {
		if utf8.RuneCountInString(v) < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if utf8.RuneCountInString(v) > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateInt validates int min, max. -1 for any
// 	m.ValidateInt("number", 10, -1, 11)  // max 18
func (m *ModelBase) ValidateInt(f string, v, min, max int) {
	if min > 0 {
		if v < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if v > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateInt64 validates int64 min, max. -1 for any
// 	m.ValidateInt64("number", 10, 6, -1) // min 6
func (m *ModelBase) ValidateInt64(f string, v, min, max int64) {
	if min > 0 {
		if v < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if v > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateFloat32 validates float32 min, max. -1 for any
// 	m.ValidateFloat32("number", 10.2, -1, 11)
func (m *ModelBase) ValidateFloat32(f string, v, min, max float32) {
	if min > 0 {
		if v < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if v > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateFloat64 validates float64 min, max. -1 for any
// 	m.ValidateFloat64("number", 10.2, -1, 11)
func (m *ModelBase) ValidateFloat64(f string, v, min, max float64) {
	if min > 0 {
		if v < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if v > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateFormat validates string format with regex string
// 	m.ValidateFormat("ip address", u.IP, `\A(\d{1,3}\.){3}\d{1,3}\z`)
func (m *ModelBase) ValidateFormat(f, v, reg string) {
	if r, _ := regexp.MatchString(reg, v); !r {
		m.AddError(f, "invalid format")
	}
}

// BaseModel interface
type BaseModel interface {
	Valid() bool
	AddError(string, string)
	SetErrors(modelErrors)
	GetErrors() modelErrors
	ResetErrors()
}
