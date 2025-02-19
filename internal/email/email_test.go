package email

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantEmail string
		wantErr   bool
	}{
		{
			name:      "Valid email",
			email:     "valid@email.com",
			wantEmail: "valid@email.com",
			wantErr:   false,
		},
		{
			name:      "Valid email with name",
			email:     "Valid Mail <valid@email.com>",
			wantEmail: "valid@email.com",
			wantErr:   false,
		},
		{
			name:      "Invalid email with name",
			email:     "Valid Mail valid@email.com>",
			wantEmail: "",
			wantErr:   true,
		},
		{
			name:      "Invalid email with missing @-sign",
			email:     "invalidemail.com",
			wantEmail: "",
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addr, err := ValidateEmail(test.email)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr = %v", err, test.wantErr)
			}
			if addr != test.wantEmail {
				t.Errorf("ValidateEmail() string = %v, wantEmail = %v", addr, test.email)
			}
		})
	}
}
