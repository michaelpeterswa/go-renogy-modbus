package gorenogymodbus_test

import (
	"testing"

	gorenogymodbus "github.com/michaelpeterswa/go-renogy-modbus"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSynthesize(t *testing.T) {
	tests := []struct {
		Name  string
		DCI   gorenogymodbus.DynamicControllerInformation
		Bytes []byte
	}{
		{
			Name: "simple",
			DCI: gorenogymodbus.DynamicControllerInformation{
				BatteryCapacitySOC:                  100,                                         // Percentage
				BatteryVoltage:                      decimal.NewFromFloat(13.6),                  // Volts
				ChargingCurrent:                     decimal.NewFromFloat(1.5),                   // Amperes
				ControllerTemperature:               25,                                          // Celsius
				BatteryTemperature:                  0,                                           // Celsius
				StreetLightLoadVoltage:              decimal.NewFromFloat(13.6),                  // Volts
				StreetLightLoadCurrent:              decimal.NewFromFloat(4),                     // Amperes
				StreetLightLoadPower:                decimal.NewFromFloat(54.4),                  // Watts
				SolarPanelVoltage:                   decimal.NewFromFloat(16.6),                  // Volts
				SolarPanelCurrent:                   decimal.NewFromFloat(1.1),                   // Amperes
				ChargingPower:                       decimal.NewFromFloat(18.26),                 // Watts
				BatteryMinimumVoltageCurrentDay:     decimal.NewFromFloat(0),                     // Volts
				BatteryMaximumVoltageCurrentDay:     decimal.NewFromFloat(13.2),                  // Volts
				MaximumChargingCurrentCurrentDay:    decimal.NewFromFloat(1.5),                   // Amperes
				MaximumDischargingCurrentCurrentDay: decimal.NewFromFloat(4),                     // Amperes
				MaximumChargingPowerCurrentDay:      decimal.NewFromFloat(19.8),                  // Watts
				MaximumDischargingPowerCurrentDay:   decimal.NewFromFloat(12),                    // Amperes
				ChargingAmpHoursCurrentDay:          decimal.NewFromFloat(4),                     // Amperes
				DischargingAmpHoursCurrentDay:       decimal.NewFromFloat(4),                     // Amperes
				PowerGenerationCurrentDay:           decimal.NewFromFloat(3),                     // Kilowatt/hours
				PowerConsumptionCurrentDay:          decimal.NewFromFloat(1),                     // Kilowatt/hours
				TotalOperatingDays:                  12,                                          // int
				TotalBatteryOverDischarges:          0,                                           // int
				TotalBatteryFullCharges:             10,                                          // int
				TotalChargingAmpHours:               decimal.NewFromFloat(10),                    // Amperes
				TotalDischargingAmpHours:            decimal.NewFromFloat(10),                    // Amperes
				CumulativePowerGeneration:           decimal.NewFromFloat(10),                    // Kilowatt/hours
				CumulativePowerConsumption:          decimal.NewFromFloat(10),                    // Kilowatt/hours
				StreetLightStatus:                   false,                                       // bool
				StreetLightBrightness:               0,                                           // Percentage
				ChargingState:                       gorenogymodbus.ChargingDeactivated.String(), // string
				ControllerFaults:                    nil,                                         // []string
			},
			Bytes: []byte{
				0x00, 0x64, 0x00, 0x88, 0x00,
				0x96, 0x19, 0x00, 0x00, 0x88,
				0x01, 0x90, 0x00, 0x36, 0x00,
				0xa6, 0x00, 0x6e, 0x00, 0x12,
				0x00, 0x00, 0x00, 0x00, 0x00,
				0x84, 0x00, 0x96, 0x01, 0x90,
				0x00, 0x13, 0x00, 0x0c, 0x00,
				0x04, 0x00, 0x04, 0x75, 0x30,
				0x75, 0x30, 0x00, 0x0c, 0x00,
				0x00, 0x00, 0x0a, 0x00, 0x00,
				0x00, 0x0a, 0x00, 0x00, 0x00,
				0x0a, 0x00, 0x01, 0x86, 0xa0,
				0x00, 0x01, 0x86, 0xa0, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.DCI.Synthesize()
			assert.NoError(t, err)
			assert.Equal(t, tc.Bytes, result)
		})
	}
}
