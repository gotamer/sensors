/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package ener314rt

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/sensors/protocol/ook"
	_ "github.com/djthorpe/sensors/protocol/openthings"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register pimote using GPIO
	gopi.RegisterModule(gopi.Module{
		Name:     "sensors/ener314rt",
		Requires: []string{"gpio", "sensors/rfm69/spi", "sensors/protocol/ook", "sensors/protocol/openthings"},
		Type:     gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			// GPIO pin configurations
			config.AppFlags.FlagUint("gpio.reset", 25, "Reset Pin (Logical)")
			config.AppFlags.FlagUint("gpio.led1", 27, "Green LED Pin (Logical)")
			config.AppFlags.FlagUint("gpio.led2", 22, "Red LED Pin (Logical)")

			// MiHome flags
			config.AppFlags.FlagString("mihome.cid", "", "20-bit Command Device ID (hexadecimal)")
			config.AppFlags.FlagUint("mihome.repeat", 0, "Command TX Repeat")
			config.AppFlags.FlagFloat64("mihome.tempoffset", 0, "Temperature Calibration Value")

			// Default spi.slave to 1
			if err := config.AppFlags.SetUint("spi.slave", 1); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			if gpio, ok := app.ModuleInstance("gpio").(gopi.GPIO); !ok {
				return nil, fmt.Errorf("Missing or invalid GPIO module")
			} else if radio, ok := app.ModuleInstance("sensors/rfm69/spi").(sensors.RFM69); !ok {
				return nil, fmt.Errorf("Missing or invalid Radio module")
			} else if ookproto, ok := app.ModuleInstance("sensors/protocol/ook").(sensors.ProtoOOK); !ok {
				return nil, fmt.Errorf("Missing or invalid OOK module")
			} else if otproto, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.ProtoOT); !ok {
				return nil, fmt.Errorf("Missing or invalid OT module")
			} else {
				config := MiHome{
					GPIO:     gpio,
					Radio:    radio,
					OOK:      ookproto,
					OT:       otproto,
					PinReset: gopi.GPIO_PIN_NONE,
					PinLED1:  gopi.GPIO_PIN_NONE,
					PinLED2:  gopi.GPIO_PIN_NONE,
				}
				if reset, _ := app.AppFlags.GetUint("gpio.reset"); reset > 0 && reset <= 0xFF {
					config.PinReset = gopi.GPIOPin(reset)
				}
				if led1, _ := app.AppFlags.GetUint("gpio.led1"); led1 > 0 && led1 <= 0xFF {
					config.PinLED1 = gopi.GPIOPin(led1)
				}
				if led2, _ := app.AppFlags.GetUint("gpio.led2"); led2 > 0 && led2 <= 0xFF {
					config.PinLED2 = gopi.GPIOPin(led2)
				}
				if cid, exists := app.AppFlags.GetString("mihome.cid"); exists {
					config.CID = cid
				}
				if repeat, exists := app.AppFlags.GetUint("mihome.repeat"); exists {
					config.Repeat = repeat
				}
				if tempoffset, exists := app.AppFlags.GetFloat64("mihome.tempoffset"); exists {
					config.TempOffset = float32(tempoffset)
				}
				return gopi.Open(config, app.Logger)
			}
		},
	})
}
