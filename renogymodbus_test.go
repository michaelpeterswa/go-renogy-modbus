package gorenogymodbus_test

import (
	"testing"

	gorenogymodbus "github.com/michaelpeterswa/go-renogy-modbus"
	"github.com/stretchr/testify/assert"
)

func TestGetControllerFaults(t *testing.T) {
	tests := []struct {
		Name           string
		InputBytes     []byte
		ExpectedErrors []string
	}{
		{
			Name:           "no error flags set",
			InputBytes:     []byte{0b00000000, 0b00000000, 0b00000000, 0b00000000},
			ExpectedErrors: nil,
		},
		{
			Name:       "single error flag set",
			InputBytes: []byte{0b00001000, 0b00000000, 0b0000000, 0b00000000},
			ExpectedErrors: []string{
				"solar panel working point overvoltage",
			},
		},
		{
			Name:           "dual error flag set",
			InputBytes:     []byte{0b0011000, 0b00000000, 0b0000000, 0b00000000},
			ExpectedErrors: []string{"solar panel working point overvoltage", "solar panel reversely connected"},
		},
		{
			Name:       "all error flag set",
			InputBytes: []byte{0b01111111, 0b11111111, 0b10000000, 0b00000000},
			ExpectedErrors: []string{
				"battery over discharge", "battery over voltage", "battery under voltage", "load short circuit",
				"load over power or load over current", "controller temperature too high", "ambient temperature too high",
				"photovoltaic input overpower", "photovoltaic input side short circuit", "photovoltaic input side over voltage",
				"solar panel counter current", "solar panel working point overvoltage", "solar panel reversely connected",
				"anti reverse mos short", "charge mos short circuit",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := gorenogymodbus.GetControllerFaults(tc.InputBytes)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			assert.ElementsMatch(t, tc.ExpectedErrors, result)
		})
	}
}
