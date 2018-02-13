package relaydriver

import "fmt"
import "github.com/davecheney/gpio"

const (
	HandshakePin = 2
	Relay1       = 3
	Relay2       = 4
	Relay3       = 14
	Relay4       = 15
)

type Driver struct {
	r1 gpio.Pin
	r2 gpio.Pin
	r3 gpio.Pin
	r4 gpio.Pin
}

func (d *Driver) SetHigh(relay int) {
	fmt.Printf("Set high %v\n", relay)

	switch relay {
	case Relay1:
		d.r1.Set()
	case Relay2:
		d.r2.Set()
	case Relay3:
		d.r3.Set()
	case Relay4:
		d.r4.Set()
	}
}

func (d *Driver) SetLow(relay int) {
	fmt.Printf("Set low %v\n", relay)
	switch relay {
	case Relay1:
		d.r1.Clear()
	case Relay2:
		d.r2.Clear()
	case Relay3:
		d.r3.Clear()
	case Relay4:
		d.r4.Clear()
	}
}

func (d *Driver) Shutdown() {
	d.r1.Clear()
	defer d.r1.Close()

	d.r2.Clear()
	defer d.r2.Close()

	d.r3.Clear()
	defer d.r3.Close()

	d.r4.Clear()
	defer d.r4.Close()

	// Disconnecting the relay board prior to running defered statements avoids
	// floating pin state
	driverDisconnect()
}

func NewDriver() (*Driver, error) {
	// Setup pins for the relays, prior to the handshake
	r1, err := gpio.OpenPin(Relay1, gpio.ModeOutput)
	if err != nil {
		return nil, fmt.Errorf("error opening relay 1, %v", err)
	}
	r2, err := gpio.OpenPin(Relay2, gpio.ModeOutput)
	if err != nil {
		return nil, fmt.Errorf("error opening relay 2, %v", err)
	}
	r3, err := gpio.OpenPin(Relay3, gpio.ModeOutput)
	if err != nil {
		return nil, fmt.Errorf("error opening relay 3, %v", err)
	}
	r4, err := gpio.OpenPin(Relay4, gpio.ModeOutput)
	if err != nil {
		return nil, fmt.Errorf("error opening relay 4, %v", err)
	}

	// Doing the handshake now avoids floating pin state getting observed by
	// the relay board.
	driverConnect()

	return &Driver{
		r1: r1,
		r2: r2,
		r3: r3,
		r4: r4,
	}, nil
}
