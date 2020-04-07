package ioman

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kaelanfouwels/iodrivers/i2c"
	"github.com/kaelanfouwels/iodrivers/i2c/sfm3000"
)

const _sfm3000Address = 0x40
const _sfm3000IsAir = false
const _sfm3000Name = "FLOW1"
const _i2cBus = 1
const _ioLoopTime = 100 * time.Millisecond

//IOMan ..
type IOMan struct {
	sfm  *sfm3000.SFM3000
	data DataPacket
	sync.Mutex
}

//NewIOMan ..
func NewIOMan() (*IOMan, error) {
	iom := IOMan{}

	iom.initialize()

	return &iom, nil
}

func (i *IOMan) initialize() error {

	//Initialize I2C
	logf("ioman", "Initializing I2C on bus %v", _i2cBus)
	i21, err := i2c.NewI2C(_i2cBus)
	if err != nil {
		return fmt.Errorf("Failed to create I2C device: %w", err)
	}

	//Initialize SFM3000
	logf("ioman", "Initializing SFM3000 %v on I2C bus %v, as 0x%x", _sfm3000Name, _i2cBus, _sfm3000Address)
	sfm, err := sfm3000.NewSFM3000(i21, _sfm3000Address, _sfm3000IsAir, _sfm3000Name)
	if err != nil {
		return fmt.Errorf("Failed to create SFM3000: %w", err)
	}

	logf("ioman", "Testing %v", sfm.Label())
	err = sfm.SoftReset()
	if err != nil {
		return fmt.Errorf("Failed to soft reset %v: %w", sfm.Label(), err)
	}
	logf("ioman", "%v soft-reset OK", sfm.Label())
	time.Sleep(50 * time.Millisecond)

	serial, err := sfm.GetSerial()
	if err != nil {
		return fmt.Errorf("Failed to get %v serial number: %w", sfm.Label(), err)
	}
	logf("ioman", "%v serial number OK: %v", sfm.Label(), serial)
	time.Sleep(50 * time.Millisecond)

	_, _, _ = sfm.GetValue() //First value is expected to be garbage
	time.Sleep(50 * time.Millisecond)

	value, crc, err := sfm.GetValue()
	if err != nil {
		return fmt.Errorf("Failed to get %v value: %w", sfm.Label(), err)
	}
	logf("ioman", "%v flow value OK: %v crc %v", sfm.Label(), value, crc)
	i.sfm = sfm
	logf("ioman", "All sensors initialized")

	return nil
}

//Start ..
func (i *IOMan) Start(cherr chan<- error) {
	lt := time.NewTicker(_ioLoopTime)
	defer lt.Stop()

	logf("ioman", "Starting io loop at %v Hz", 1/_ioLoopTime.Seconds())

	for range lt.C {
		sfmval, _, err := i.sfm.GetValue()
		if err != nil {
			cherr <- err
			return
		}

		i.Lock()
		i.data.Flow = sfmval
		i.data.Valid = true
		i.Unlock()
	}

	cherr <- fmt.Errorf("io loop ended unexpectedly")
}

//GetDataPacket ..
func (i *IOMan) GetDataPacket() DataPacket {
	i.Lock()
	defer i.Unlock()
	return i.data
}

//Destroy ..
func (*IOMan) Destroy() {

}

func logf(owner string, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[%v] %v", owner, message)
}
