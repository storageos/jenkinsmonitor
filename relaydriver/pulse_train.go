package relaydriver

import (
	"github.com/davecheney/gpio"
	"sync"
	"time"
)

var pulseTrainMutex = sync.Mutex{}
var connected = false

func driverConnect() {
	pulseTrainMutex.Lock()
	defer pulseTrainMutex.Unlock()

	if !connected {
		sendPulseTrain()
		connected = true
	}
}

func driverDisconnect() {
	pulseTrainMutex.Lock()
	defer pulseTrainMutex.Unlock()

	if connected {
		sendPulseTrain()
		connected = false
	}
}

func sendPulseTrain() {
	pulsePin, err := gpio.OpenPin(HandshakePin, gpio.ModeOutput)
	if err != nil {
		panic("failed to open GPIO02 pin for output")
	}
	pulsePin.Clear()

	// Give the relay board some time to recover from the noise
	time.Sleep(time.Second)

	t := time.NewTicker(50 * time.Millisecond)
	defer t.Stop()

	for i := 0; i < 4; i++ {
		pulsePin.Set()
		<-t.C
		pulsePin.Clear()
		<-t.C
	}
}
