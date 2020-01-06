package somfy

import "time"

const (
	WakeUpPulseLength             = 9415 * time.Microsecond
	SilenceAfterWakeUpPulseLength = 89565 * time.Microsecond
	HardwareSyncPulseLength       = 2560 * time.Microsecond
	SoftwareSyncPulseLength       = 4550 * time.Microsecond
	// Data is Manchester encoded: 1=rising edge, 0=falling edge. One bit is encoded per pulse length (~1208 us)
	DataHalfPulseLength      = 640 * time.Microsecond
	InterFrameGapPulseLength = 30415 * time.Microsecond

	DefaultEncryptionKey = 0xA7
)

type Pulse struct {
	IsHigh bool
	Length time.Duration
}

type Control byte

// 4-bit Control codes, this indicates the button that is pressed
const (
	MyValue      Control = 0x1
	UpValue      Control = 0x2
	DownValue    Control = 0x4
	ProgramValue Control = 0x8
)

type PayloadData struct {
	// Most significant 4-bit are always 0xA, Least Significant bits is a linear counter
	EncryptionKey byte
	Control       Control
	// 16-bit rolling code (big-endian) increased with every button press
	RollingCode uint16
	// 24-bit identifier of sending device (little-endian)
	Address uint32
}

type FrameData [7]byte

func NewPayload(key byte) PayloadData {
	return PayloadData{
		EncryptionKey: key,
	}
}

// The checksum is calculated by doing a XOR of all nibbles of the frame.
// To verify a frame check if the value returned by the checksum algorithm is equal to 0, if not the frame is corrupt.
func (f *FrameData) AddCheckSum() {
	var checkSum byte
	for i := 0; i < 7; i++ {
		checkSum = checkSum ^ f[i] ^ (f[i] >> 4)
	}
	checkSum &= 0b1111

	f[1] |= checkSum
}

// The payload data is obfuscated by doing an XOR between the byte to obfuscate and the previous obfuscated byte
func (f *FrameData) Obfuscate() {
	for i := 1; i < 7; i++ {
		f[i] ^= f[i-1]
	}
}
