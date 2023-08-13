package gorenogymodbus

import (
	"testing"

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
			result, err := getControllerFaults(tc.InputBytes)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			assert.ElementsMatch(t, tc.ExpectedErrors, result)
		})
	}
}

func TestBuildStreetLightStateAndBrightnessByte(t *testing.T) {
	tests := []struct {
		Name                  string
		StreetLightStatus     bool
		StreetLightBrightness int
		ResultByte            byte
		ShouldError           bool
	}{
		{
			Name:                  "load off and minimum brightness",
			StreetLightStatus:     false,
			StreetLightBrightness: 0,
			ResultByte:            0b00000000,
			ShouldError:           false,
		},
		{
			Name:                  "load on and max brightness",
			StreetLightStatus:     true,
			StreetLightBrightness: 100,
			ResultByte:            0b11100100,
			ShouldError:           false,
		},
		{
			Name:                  "load on and 50% brightness",
			StreetLightStatus:     true,
			StreetLightBrightness: 50,
			ResultByte:            0b10110010,
			ShouldError:           false,
		},
		{
			Name:                  "load on and 150% brightness, should error",
			StreetLightStatus:     true,
			StreetLightBrightness: 150,
			ResultByte:            0b10010110,
			ShouldError:           true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := buildStreetLightStatusAndBrightnessByte(tc.StreetLightStatus, tc.StreetLightBrightness)
			if !tc.ShouldError {
				assert.NoError(t, err)
				assert.Equal(t, tc.ResultByte, result)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestGetChargingState(t *testing.T) {
	tests := []struct {
		Name                  string
		Byte                  byte
		ExpectedChargingState ChargingState
	}{
		{
			Name:                  "minimum charging state",
			Byte:                  0x0,
			ExpectedChargingState: ChargingDeactivated,
		},
		{
			Name:                  "maximum charging state",
			Byte:                  0x6,
			ExpectedChargingState: CurrentLimitingOverPower,
		},
		{
			Name:                  "invalid charging state",
			Byte:                  0x7,
			ExpectedChargingState: 0x7,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			result := getChargingState(tc.Byte)

			assert.Equal(t, tc.ExpectedChargingState, result)
		})
	}
}
