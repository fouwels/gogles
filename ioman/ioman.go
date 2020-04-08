package ioman

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kaelanfouwels/iodrivers/i2c"
	"github.com/kaelanfouwels/iodrivers/i2c/sfm3000"
	"github.com/kaelanfouwels/iodrivers/spi/mcp3208"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const _flow1Address = 0x40
const _flow1IsAir = false
const _i2cBus = "/dev/i2c-1"
const _ioLoopTime = (1 * time.Second) / 1000 //1KHz

const _spiBus = "/dev/spidev0.1"
const _spiSpeed = physic.Frequency(1 * physic.MegaHertz)
const _adc1ChipSelect = 1

const _breathInFlowThreshold = 20
const _breathOutFlowThreshold = -20

//IOMan ..
type IOMan struct {
	sensors  *phySensors
	o        DataPacket //Output data variables - external buffer. DataPacket is copied from working buffer to output buffer at end of io cycle.
	moutputs sync.Mutex
}

type phySensors struct {
	Flow *sfm3000.SFM3000
	ADC  *mcp3208.Mcp3208
}

//NewIOMan ..
func NewIOMan() (*IOMan, error) {
	iom := IOMan{}

	logf("ioman", "Initializing sensors")
	sens, err := iom.initialize()
	if err != nil {
		logf("ioman", "Initialization failed")
		return nil, fmt.Errorf("Failed to initialize: %w", err)
	}

	logf("ioman", "Performing Self Test")
	err = iom.selftest(sens)
	if err != nil {
		logf("ioman", "Self test failed")
		return nil, fmt.Errorf("Failed to self test: %w", err)
	}

	iom.sensors = sens

	return &iom, nil
}

func (io *IOMan) initialize() (*phySensors, error) {

	//Initialize host - required for SPI driver
	logf("ioman:initialize", "Initializing host")
	_, err := host.Init()
	if err != nil {
		log.Fatalf("Failed to initialize periph.io host: %v", err)
	}

	//Initialize I2C
	logf("ioman:initialize", "Initializing I2C on bus %v", _i2cBus)
	i21, err := i2c.NewI2C(_i2cBus)
	if err != nil {
		return nil, fmt.Errorf("Failed to create I2C device: %w", err)
	}

	//Initialize SFM3000
	logf("ioman:initialize", "Initializing SFM3000 on I2C bus %v 0x%x as %v", _i2cBus, _flow1Address, "FLOW1")
	flow1, err := sfm3000.NewSFM3000(i21, _flow1Address, _flow1IsAir, "FLOW1")
	if err != nil {
		return nil, fmt.Errorf("Failed to create SFM3000: %w", err)
	}

	//Initialize SPI
	logf("ioman:initialize", "Initializing SPI on bus %v at %v ", _spiBus, _spiSpeed)
	s, err := spireg.Open(_spiBus)
	if err != nil {
		log.Fatalf("Failed to open SPI bus %v: %v", _spiBus, err)
	}
	conn, err := s.Connect(_spiSpeed, spi.Mode0, 8)
	if err != nil {
		log.Fatalf("Failed to connect SPI: %v", err)
	}

	//Initialize MCP3208
	logf("ioman:initialize", "Initializing MCP3208 on SPI bus %v CS %v as %v", _spiBus, _adc1ChipSelect, "ADC1")
	adc1, err := mcp3208.NewMcp3208(conn, "ADC1")
	if err != nil {
		log.Fatalf("Failed to create MCP3208: %v", err)
	}

	return &phySensors{
		Flow: flow1,
		ADC:  adc1,
	}, nil
}

func (io *IOMan) selftest(sensors *phySensors) error {

	// Test SFM3000

	// Test Reset
	logf("ioman:selftest", "Testing %v", sensors.Flow.Label())
	err := sensors.Flow.SoftReset()
	if err != nil {
		return fmt.Errorf("Failed to soft reset %v: %w", sensors.Flow.Label(), err)
	}
	logf("ioman:selftest", "%v soft-reset OK", sensors.Flow.Label())
	time.Sleep(50 * time.Millisecond)

	// Test GetSerial
	serial, err := sensors.Flow.GetSerial()
	if err != nil {
		return fmt.Errorf("Failed to get %v serial number: %w", sensors.Flow.Label(), err)
	}
	logf("ioman:selftest", "%v serial number OK: %v", sensors.Flow.Label(), serial)
	time.Sleep(50 * time.Millisecond)

	// Test GetValue
	_, _, _ = sensors.Flow.GetValue() //First value is expected to be garbage
	time.Sleep(50 * time.Millisecond)

	value, crc, err := sensors.Flow.GetValue()
	if err != nil {
		return fmt.Errorf("Failed to get %v value: %w", sensors.Flow.Label(), err)
	}
	logf("ioman:selftest", "%v flow value OK: %v crc %v", sensors.Flow.Label(), value, crc)

	//Test MCP3208
	logf("ioman:selftest", "Testing %v", sensors.ADC.Label())

	// Test GetValue
	_, _ = sensors.ADC.GetValues(0, 4) //First value is expected to be garbage
	time.Sleep(50 * time.Millisecond)

	vals, err := sensors.ADC.GetValues(0, 4)
	if err != nil {
		return fmt.Errorf("Failed to get %v values: %w", sensors.ADC.Label(), err)
	}
	logf("ioman:selftest", "%v adc values OK: %v ", sensors.ADC.Label(), vals)

	return nil
}

//Start ..
func (io *IOMan) Start(cherr chan<- error) {
	logf("ioman:start", "Starting at %v hz", 1/_ioLoopTime.Seconds())
	lt := time.NewTicker(_ioLoopTime)
	defer lt.Stop()

	i := internalPacket{} //Internal data variables

	for range lt.C {

		d := DataPacket{}

		// Read inputs
		fval, fcrc, ferr := io.sensors.Flow.GetValue()
		flow := Flow{
			Val: fval,
			CRC: fcrc,
			Err: ferr,
		}

		ivals, ferr := io.sensors.ADC.GetValues(0, 4)
		adc := ADC{
			Vals: ivals,
			Err:  ferr,
		}

		d.Sensors = Sensors{
			Flow: flow,
			ADC:  adc,
		}
		d.Timestamp = time.Now()

		// If inputs read ok
		if d.Sensors.Flow.Err == nil && d.Sensors.ADC.Err == nil {
			d.Stats.OkReads++

			io.states(&d, &i)
			io.calculate(&d, &i)
			d.Valid = true

		} else {
			d.Stats.FailedReads++

			d.Valid = false
		}

		// Copy to output registers
		io.moutputs.Lock()
		io.o = d
		io.moutputs.Unlock()
	}

	cherr <- fmt.Errorf("io loop ended unexpectedly")
}

func (io *IOMan) states(d *DataPacket, i *internalPacket) {
	breathin := d.Sensors.Flow.Val > _breathInFlowThreshold
	breathout := d.Sensors.Flow.Val < _breathOutFlowThreshold

	lastState := d.State
	state := StateRest

	if breathin {
		state = StateBreathingIn
		//logf("SPECIAL", "Breathing In %v", d.Sensors.Flow.Val)
	}
	if breathout {
		state = StateBreathingOut
		//logf("SPECIAL", "Breathing Out %v", d.Sensors.Flow.Val)
	}
	if breathin && breathout {
		state = StateError
		logf("ioman:states", "Breathing Error %v", d.Sensors.Flow.Val)
	}

	if state != lastState {
		i.lastStateChangeNmin1 = i.lastStateChangeN
		i.lastStateChangeN = d.Timestamp
	}
	i.lastState = lastState
	d.State = state
}

func (io *IOMan) calculate(d *DataPacket, i *internalPacket) {

	calc := Calculated{}

	// If moving into breathing state, reset counters
	if d.State == StateBreathingIn && i.lastState != StateBreathingIn {
		i.flowAverageTotal = 0
		i.flowAverageN = 0
	}

	// If in breathing state, increment counters
	if d.State == StateBreathingIn {
		i.flowAverageTotal += float64(d.Sensors.Flow.Val)
		i.flowAverageN++
	}

	// If leaving breathing state, calculate integrated flow.
	if d.State != StateBreathingIn && i.lastState == StateBreathingIn && d.State != StateError {
		duration := d.Timestamp.Sub(i.lastStateChangeNmin1)
		flow := (i.flowAverageTotal / float64(i.flowAverageN)) * duration.Minutes()

		calc.FlowIntegrated = flow
		calc.FlowIntegratedTimestamp = d.Timestamp
	}

	d.Calculated = calc
}

//GetDataPacket ..
func (io *IOMan) GetDataPacket() DataPacket {
	io.moutputs.Lock()
	defer io.moutputs.Unlock()
	return io.o
}

//Destroy ..
func (*IOMan) Destroy() {

}

func logf(owner string, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[%v] %v", owner, message)
}
