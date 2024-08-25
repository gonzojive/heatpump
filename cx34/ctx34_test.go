package cx34

import "testing"

func TestAirConditioningMode(t *testing.T) {
	testCases := []struct {
		mode          AirConditioningMode
		wantIsCooling bool
		wantIsHeating bool
		wantString    string
	}{
		{AirConditioningModeCooling, true, false, "Cooling"},
		{AirConditioningModeHeating, false, true, "Heating"},
		{AirConditioningModeDHW, false, false, "Domestic Hot Water"},
		{AirConditioningModeCoolingAndDHW, true, false, "Cooling + Domestic Hot Water"},
		{AirConditioningModeHeatingAndDHW, false, true, "Heating + Domestic Hot Water"},
		{AirConditioningMode(42), false, false, "Unknown AirConditioningMode (42)"},
	}

	for _, tc := range testCases {
		t.Run(tc.mode.String(), func(t *testing.T) {
			if got := tc.mode.IsCooling(); got != tc.wantIsCooling {
				t.Errorf("IsCooling(%v) = %v, want %v", tc.mode, got, tc.wantIsCooling)
			}
			if got := tc.mode.IsHeating(); got != tc.wantIsHeating {
				t.Errorf("IsHeating(%v) = %v, want %v", tc.mode, got, tc.wantIsHeating)
			}
			if got := tc.mode.String(); got != tc.wantString {
				t.Errorf("String(%v) = %q, want %q", tc.mode, got, tc.wantString)
			}
		})
	}
}

func TestAirConditioningMode_parse(t *testing.T) {
	testCases := []struct {
		want          AirConditioningMode
		registerValue uint16
		wantErr       bool
	}{
		{AirConditioningModeCooling, 0, false},
		{AirConditioningModeHeating, 1, false},
		{AirConditioningModeDHW, 2, false},
		{AirConditioningModeCoolingAndDHW, 3, false},
		{AirConditioningModeHeatingAndDHW, 4, false},
		{AirConditioningMode(42), 42, true},
	}

	for _, tc := range testCases {
		t.Run(tc.want.String(), func(t *testing.T) {
			got, err := parseAirConditioningMode(tc.registerValue)
			if err != nil {
				if !tc.wantErr {
					t.Errorf("parseAirConditioningMode(%v) error = %v, wantErr = %v", tc.registerValue, err, tc.wantErr)
				}
				return
			}

			if got != tc.want {
				t.Errorf("parseAirConditioningMode(%v) = %q, want %q", tc.registerValue, got, tc.want)
			}
		})
	}
}
