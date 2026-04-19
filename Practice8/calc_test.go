package main

import (
	"testing"
)

//Add (from Steps 1 & 2) ───────────────────────────────────────────────────

func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

//Subtract 

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 10, 3, 7},
		{"positive minus zero", 5, 0, 5},
		{"negative minus positive", -3, 4, -7},
		{"both negative", -5, -2, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Divide 

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
		errMsg  string
	}{
		{"normal division", 10, 2, 5, false, ""},
		{"divide by one", 7, 1, 7, false, ""},
		{"negative dividend", -9, 3, -3, false, ""},
		{"both negative", -8, -2, 4, false, ""},
		{"zero dividend", 0, 5, 0, false, ""},
		{"divide by zero", 10, 0, 0, true, "division by zero"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Divide(%d, %d) expected error, got nil", tt.a, tt.b)
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Divide(%d, %d) error = %q; want %q", tt.a, tt.b, err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("Divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
				}
				if got != tt.want {
					t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
				}
			}
		})
	}
}
