package gorenogymodbus

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/goburrow/modbus"
)

type ModbusClient struct {
	Client modbus.Client
}

func NewModbusClient(logger *log.Logger, address string) (*ModbusClient, error) {
	// Modbus RTU/ASCII
	handler := modbus.NewRTUClientHandler(address)
	handler.BaudRate = 9600
	handler.SlaveId = 1
	handler.Timeout = 1 * time.Second
	handler.StopBits = 1
	handler.DataBits = 8
	handler.Parity = "N"
	handler.Logger = logger

	err := handler.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect modbus handler: %w", err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	return &ModbusClient{client}, nil
}

func (mc *ModbusClient) ReadData() ([]byte, error) {
	var (
		dataStartAddress uint16 = 0x100
		dataQuantity     uint16 = 34
	)

	res, err := mc.readHoldingRegisters(dataStartAddress, dataQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to read holding registers: %w", err)
	}

	return res, nil
}

func (mc *ModbusClient) readHoldingRegisters(address uint16, quantity uint16) (results []byte, err error) {
	res, err := mc.Client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type ChargingState int

const (
	ChargingDeactivated ChargingState = iota
	ChargingActivated
	MPPTChargingMode
	EqualizingChargingMode
	BoostChargingMode
	FloatingChargingMode
	CurrentLimitingOverPower
)

func (cs ChargingState) String() string {
	switch cs {
	case ChargingDeactivated:
		return "charging deactivated"
	case ChargingActivated:
		return "charging activated"
	case MPPTChargingMode:
		return "mppt charging mode"
	case EqualizingChargingMode:
		return "equalizing charging mode"
	case BoostChargingMode:
		return "boost charging mode"
	case FloatingChargingMode:
		return "floating charging mode"
	case CurrentLimitingOverPower:
		return "current limiting overpower"
	default:
		return "unknown"
	}
}

type ControllerFault int

const (
	NoFault ControllerFault = iota
	ChargeMOSShortCircuit
	AntiReverseMOSShort
	SolarPanelReverselyConnected
	SolarPanelWorkingPointOvervoltage
	SolarPanelCounterCurrent
	PhotovoltaicInputSideOverVoltage
	PhotovoltaicInputSideShortCircuit
	PhotovoltaicInputOverpower
	AmbientTemperatureTooHigh
	ControllerTemperatureTooHigh
	LoadOverPowerOrLoadOverCurrent
	LoadShortCircuit
	BatteryUnderVoltage
	BatteryOverVoltage
	BatteryOverDischarge
)

func (cf ControllerFault) String() string {
	switch cf {
	case NoFault:
		return "no fault"
	case ChargeMOSShortCircuit:
		return "charge mos short circuit"
	case AntiReverseMOSShort:
		return "anti reverse mos short"
	case SolarPanelReverselyConnected:
		return "solar panel reversely connected"
	case SolarPanelWorkingPointOvervoltage:
		return "solar panel working point overvoltage"
	case SolarPanelCounterCurrent:
		return "solar panel counter current"
	case PhotovoltaicInputSideOverVoltage:
		return "photovoltaic input side over voltage"
	case PhotovoltaicInputSideShortCircuit:
		return "photovoltaic input side short circuit"
	case PhotovoltaicInputOverpower:
		return "photovoltaic input overpower"
	case AmbientTemperatureTooHigh:
		return "ambient temperature too high"
	case ControllerTemperatureTooHigh:
		return "controller temperature too high"
	case LoadOverPowerOrLoadOverCurrent:
		return "load over power or load over current"
	case LoadShortCircuit:
		return "load short circuit"
	case BatteryUnderVoltage:
		return "battery under voltage"
	case BatteryOverVoltage:
		return "battery over voltage"
	case BatteryOverDischarge:
		return "battery over discharge"
	default:
		return "unknown"
	}
}

type DynamicControllerInformation struct {
	BatteryCapacitySOC                  int     `json:"battery_capacity_soc"`                    // 0x100
	BatteryVoltage                      float64 `json:"battery_voltage"`                         // 0x101
	ChargingCurrent                     float64 `json:"charging_current"`                        // 0x102
	ControllerTemperature               int     `json:"controller_temperature"`                  // 0x103 ?
	BatteryTemperature                  int     `json:"battery_temperature"`                     // 0x103 ?
	StreetLightLoadVoltage              float64 `json:"street_light_load_voltage"`               // 0x104
	StreetLightLoadCurrent              float64 `json:"street_light_load_current"`               // 0x105
	StreetLightLoadPower                float64 `json:"street_light_load_power"`                 // 0x106
	SolarPanelVoltage                   float64 `json:"solar_panel_voltage"`                     // 0x107
	SolarPanelCurrent                   float64 `json:"solar_panel_current"`                     // 0x108
	ChargingPower                       float64 `json:"charging_power"`                          // 0x109
	BatteryMinimumVoltageCurrentDay     float64 `json:"battery_minimum_voltage_current_day"`     // 0x10B
	BatteryMaximumVoltageCurrentDay     float64 `json:"battery_maximum_voltage_current_day"`     // 0x10C
	MaximumChargingCurrentCurrentDay    float64 `json:"maximum_charging_current_current_day"`    // 0x10D
	MaximumDischargingCurrentCurrentDay float64 `json:"maximum_discharging_current_current_day"` // 0x10E
	MaximumChargingPowerCurrentDay      float64 `json:"maximum_charging_power_current_day"`      // 0x10F
	MaximumDischargingPowerCurrentDay   float64 `json:"maximum_discharging_power_current_day"`   // 0x110
	ChargingAmpHoursCurrentDay          float64 `json:"charging_amp_hours_current_day"`          // 0x111
	DischargingAmpHoursCurrentDay       float64 `json:"discharging_amp_hours_current_day"`       // 0x112
	PowerGenerationCurrentDay           float64 `json:"power_generation_current_day"`            // 0x113
	PowerConsumptionCurrentDay          float64 `json:"power_consumption_current_day"`           // 0x114
	TotalOperatingDays                  int     `json:"total_operating_days"`                    // 0x115
	TotalBatteryOverDischarges          int     `json:"total_battery_over_discharges"`           // 0x116
	TotalBatteryFullCharges             int     `json:"total_battery_full_charges"`              // 0x117
	TotalChargingAmpHours               float64 `json:"total_charging_amp_hours"`                // 0x118-119
	TotalDischargingAmpHours            float64 `json:"total_discharging_amp_hours"`             // 0x11A-11B
	CumulativePowerGeneration           float64 `json:"cumulative_power_generation"`             // 0x11C-11D
	CumulativePowerConsumption          float64 `json:"cumulative_power_consumption"`            // 0x11E-11F
	StreetLightStatus                   bool    `json:"street_light_status"`                     // 0x120 (eight higher bits)
	StreetLightBrightness               int     `json:"street_light_brightness"`                 // 0x120 (eight higher bits)
	ChargingState                       string  `json:"charging_state"`                          // 0x120 (eight lower bits)
	ControllerFault                     string  `json:"controller_fault"`                        // 0x121-122
}

func Parse(dataBytes []byte) DynamicControllerInformation {
	return DynamicControllerInformation{
		BatteryCapacitySOC:                  int(binary.BigEndian.Uint16(dataBytes[0:2])),                 // 0x100
		BatteryVoltage:                      float64(binary.BigEndian.Uint16(dataBytes[2:4])) * 0.1,       // 0x101
		ChargingCurrent:                     float64(binary.BigEndian.Uint16(dataBytes[4:6])) * 0.01,      // 0x102
		ControllerTemperature:               int(int8(dataBytes[6])),                                      // 0x103 first byte
		BatteryTemperature:                  int(int8(dataBytes[7])),                                      // 0x103 second byte
		StreetLightLoadVoltage:              float64(binary.BigEndian.Uint16(dataBytes[8:10])) * 0.1,      // 0x104
		StreetLightLoadCurrent:              float64(binary.BigEndian.Uint16(dataBytes[10:12])) * 0.01,    // 0x105
		StreetLightLoadPower:                float64(binary.BigEndian.Uint16(dataBytes[12:14])),           // 0x106
		SolarPanelVoltage:                   float64(binary.BigEndian.Uint16(dataBytes[14:16])) * 0.1,     // 0x107
		SolarPanelCurrent:                   float64(binary.BigEndian.Uint16(dataBytes[16:18])) * 0.01,    // 0x108
		ChargingPower:                       float64(binary.BigEndian.Uint16(dataBytes[18:20])),           // 0x109
		BatteryMinimumVoltageCurrentDay:     float64(binary.BigEndian.Uint16(dataBytes[22:24])) * 0.1,     // 0x10B
		BatteryMaximumVoltageCurrentDay:     float64(binary.BigEndian.Uint16(dataBytes[24:26])) * 0.1,     // 0x10C
		MaximumChargingCurrentCurrentDay:    float64(binary.BigEndian.Uint16(dataBytes[26:28])) * 0.01,    // 0x10D
		MaximumDischargingCurrentCurrentDay: float64(binary.BigEndian.Uint16(dataBytes[28:30])) * 0.01,    // 0x10E
		MaximumChargingPowerCurrentDay:      float64(binary.BigEndian.Uint16(dataBytes[30:32])),           // 0x10F
		MaximumDischargingPowerCurrentDay:   float64(binary.BigEndian.Uint16(dataBytes[32:34])),           // 0x110
		ChargingAmpHoursCurrentDay:          float64(binary.BigEndian.Uint16(dataBytes[34:36])),           // 0x111
		DischargingAmpHoursCurrentDay:       float64(binary.BigEndian.Uint16(dataBytes[36:38])),           // 0x112
		PowerGenerationCurrentDay:           float64(binary.BigEndian.Uint16(dataBytes[38:40])) * 10000.0, // 0x113 (deciwatt/hour conversion)
		PowerConsumptionCurrentDay:          float64(binary.BigEndian.Uint16(dataBytes[40:42])) * 10000.0, // 0x114 (deciwatt/hour conversion)
		TotalOperatingDays:                  int(binary.BigEndian.Uint16(dataBytes[42:44])),               // 0x115
		TotalBatteryOverDischarges:          int(binary.BigEndian.Uint16(dataBytes[44:46])),               // 0x116
		TotalBatteryFullCharges:             int(binary.BigEndian.Uint16(dataBytes[46:48])),               // 0x117
		TotalChargingAmpHours:               float64(binary.BigEndian.Uint32(dataBytes[48:52])),           // 0x118-119
		TotalDischargingAmpHours:            float64(binary.BigEndian.Uint32(dataBytes[52:56])),           // 0x11A-11B
		CumulativePowerGeneration:           float64(binary.BigEndian.Uint32(dataBytes[56:60])) * 10000.0, // 0x11C-11D (deciwatt/hour conversion)
		CumulativePowerConsumption:          float64(binary.BigEndian.Uint32(dataBytes[60:64])) * 10000.0, // 0x11E-11F (deciwatt/hour conversion)
		StreetLightStatus:                   dataBytes[64]&0x80 != 0,                                      // 0x120 (eight higher bits)
		StreetLightBrightness:               int(dataBytes[64] & 0x7F),                                    // 0x120 (eight higher bits) may or may not be correct logic
		ChargingState:                       getChargingState(dataBytes[65]).String(),
		ControllerFault:                     getControllerFault(dataBytes[66:70]).String(),
	}
}

func getChargingState(b byte) ChargingState {
	switch b {
	case 0x00:
		return ChargingDeactivated
	case 0x01:
		return ChargingActivated
	case 0x02:
		return MPPTChargingMode
	case 0x03:
		return EqualizingChargingMode
	case 0x04:
		return BoostChargingMode
	case 0x05:
		return FloatingChargingMode
	case 0x06:
		return CurrentLimitingOverPower
	default:
		return -1 // unimplemented
	}
}

func getControllerFault(b []byte) ControllerFault {
	return NoFault // currently unimplemented
}
