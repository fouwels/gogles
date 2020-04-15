package ioman

import (
	"fmt"
	"log"
	"sync"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"

	"github.com/kaelanfouwels/iodrivers/i2c/sfm3000"
	"github.com/kaelanfouwels/iodrivers/spi/mcp3208"
	"github.com/kaelanfouwels/iodrivers/spi/mcp4921"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const _sampleRate = (1 * time.Second) / 1000 //1KHz

const _flow1Address = 0x40
const _flow1IsAir = false
const _i2cBus = "/dev/i2c-1"

const _adcSpiBus = "/dev/spidev0.1"
const _dacSpiBus = "/dev/spidev0.0"
const _spiSpeed = physic.Frequency(1 * physic.MegaHertz)

//IOMan ..
type IOMan struct {
	sensors  *phySensors
	o        DataPacket //Output data variables - external buffer. DataPacket is copied from working buffer to output buffer at end of io cycle.
	moutputs sync.Mutex
}

type phySensors struct {
	Flow *sfm3000.SFM3000
	ADC  *mcp3208.Mcp3208
	DAC  *mcp4921.Mcp4921
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
		return nil, fmt.Errorf("Failed to initialize periph.io host: %v", err)
	}

	//Initialize I2C
	logf("ioman:initialize", "Initializing I2C on bus %v", _i2cBus)
	i2cbus, err := i2creg.Open(_i2cBus)
	if err != nil {
		return nil, fmt.Errorf("Failed to create I2C device: %w", err)
	}
	//Initialize SPI1
	logf("ioman:initialize", "Initializing SPI on bus %v at %v ", _adcSpiBus, _spiSpeed)
	s1bus, err := spireg.Open(_adcSpiBus)
	if err != nil {
		return nil, fmt.Errorf("Failed to open SPI bus %v: %v", _adcSpiBus, err)
	}

	conn1, err := s1bus.Connect(_spiSpeed, spi.Mode0, 8)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect SPI: %v", err)
	}

	//Initialize SPI2
	logf("ioman:initialize", "Initializing SPI on bus %v at %v ", _dacSpiBus, _spiSpeed)
	s2bus, err := spireg.Open(_dacSpiBus)
	if err != nil {
		return nil, fmt.Errorf("Failed to open SPI bus %v: %v", _adcSpiBus, err)
	}

	conn2, err := s2bus.Connect(_spiSpeed, spi.Mode0, 8)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect SPI: %v", err)
	}

	//Initialize SFM3000
	logf("ioman:initialize", "Initializing SFM3000 on I2C bus %v 0x%x as %v", _i2cBus, _flow1Address, "FLOW1")

	dev := i2c.Dev{
		Bus:  i2cbus,
		Addr: _flow1Address,
	}

	flow1, err := sfm3000.NewSFM3000(&dev, _flow1Address, _flow1IsAir, "FLOW1")
	if err != nil {
		return nil, fmt.Errorf("Failed to create SFM3000: %w", err)
	}

	//Initialize MCP3208
	logf("ioman:initialize", "Initializing MCP3208 on SPI bus %v as %v", _adcSpiBus, "ADC1")
	adc1, err := mcp3208.NewMcp3208(conn1, "ADC1")
	if err != nil {
		return nil, fmt.Errorf("Failed to create MCP3208: %v", err)
	}

	//Initialize MCP4921
	logf("ioman:initialize", "Initializing MCP4921 on SPI bus %v as %v", _dacSpiBus, "DAC1")
	dac1, err := mcp4921.NewMcp4921(conn2, "DAC1", mcp4921.EnumBufferedTrue, mcp4921.EnumOutputGain1x, mcp4921.EnumShutdownModeActive)
	if err != nil {
		log.Fatalf("Failed to create MCP4921: %v", err)
	}

	return &phySensors{
		Flow: flow1,
		ADC:  adc1,
		DAC:  dac1,
	}, nil
}

func (io *IOMan) selftest(sensors *phySensors) error {

	//Test SFM3000
	logf("ioman:selftest", "Testing %v", sensors.Flow.Label())

	err := sensors.Flow.SoftReset()
	if err != nil {
		return fmt.Errorf("Failed to soft reset %v: %w", sensors.Flow.Label(), err)
	}
	logf("ioman:selftest", "%v soft-reset OK", sensors.Flow.Label())
	time.Sleep(100 * time.Millisecond) // Wait for sensor to reset

	serial, err := sensors.Flow.GetSerial()
	if err != nil {
		return fmt.Errorf("Failed to get %v serial number: %w", sensors.Flow.Label(), err)
	}
	logf("ioman:selftest", "%v serial number OK: %v", sensors.Flow.Label(), serial)

	_, _, _, _ = sensors.Flow.GetValue() //First value is expected to be garbage

	value, crc, _, err := sensors.Flow.GetValue()
	if err != nil {
		return fmt.Errorf("Failed to get %v value: %w", sensors.Flow.Label(), err)
	}
	logf("ioman:selftest", "%v flow value OK: %v crc %v", sensors.Flow.Label(), value, crc)

	//Test MCP3208
	logf("ioman:selftest", "Testing %v", sensors.ADC.Label())

	_, _, _ = sensors.ADC.GetValues(0, 4) //First value is expected to be garbage
	time.Sleep(50 * time.Millisecond)

	vals, _, err := sensors.ADC.GetValues(0, 4)
	if err != nil {
		return fmt.Errorf("Failed to get %v values: %w", sensors.ADC.Label(), err)
	}
	logf("ioman:selftest", "%v adc values OK: %v ", sensors.ADC.Label(), vals)

	//Test MCP4921
	logf("ioman:selftest", "Testing %v", sensors.DAC.Label())
	err = sensors.DAC.Write(0) //Write closed
	if err != nil {
		return fmt.Errorf("Failed to write %v to %v", 0, sensors.DAC.Label())
	}
	return nil
}

//Start ..
func (io *IOMan) Start(cherr chan<- error) {
	logf("ioman:start", "Starting at %v hz", 1/_sampleRate.Seconds())
	lt := time.NewTicker(_sampleRate)
	defer lt.Stop()
	cont := newController(_sampleRate)

	for range lt.C {

		// Read inputs
		fval, fcrc, tstamp, ferr := io.sensors.Flow.GetValue()
		flow := Flow{
			Val:       fval,
			CRC:       fcrc,
			Timestamp: tstamp,
			Err:       ferr,
		}

		ivals, tstamp, ferr := io.sensors.ADC.GetValues(0, 4)
		adc := ADC{
			Vals:      ivals,
			Err:       ferr,
			Timestamp: tstamp,
		}

		sensors := Sensors{
			Flow: flow,
			ADC:  adc,
		}

		d := DataPacket{
			Sensors:   sensors,
			Timestamp: time.Now(),
		}

		// If inputs read ok
		if d.Sensors.Flow.Err == nil && d.Sensors.ADC.Err == nil {
			d.Stats.OkReads++

			cont.buffers(sensors)
			state := cont.states(sensors)
			calculated := cont.calculate(sensors)

			d.State = state
			d.Calculated = calculated
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
