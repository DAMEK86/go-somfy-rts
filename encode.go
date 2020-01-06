package somfy

import "fmt"

type Encoder interface {
	Marshal(data PayloadData, frameRepetition int) []Pulse
}

type encoder struct {
}

func NewEncoder() Encoder {
	return &encoder{}
}

func (e *encoder) Marshal(data PayloadData, frameRepetition int) []Pulse {
	var frame = FrameData{}
	frame[0] = data.EncryptionKey
	frame[1] = byte(data.Control << 4)
	frame[2] = byte(data.RollingCode >> 8)
	frame[3] = byte(data.RollingCode & 0xFF)
	frame[4] = byte(data.Address >> 16)
	frame[5] = byte((data.Address >> 8) & 0xFF)
	frame[6] = byte(data.Address & 0xFF)
	fmt.Printf("Frame:\t\t [%# x]\n", frame)

	frame.AddCheckSum()
	fmt.Printf("With cks:\t [%# x]\n", frame)

	frame.Obfuscate()
	fmt.Printf("Obfuscated:\t [%# x]\n", frame)

	return e.generateWaveForm(frame, frameRepetition)
}

func (e *encoder) generateWaveForm(frame FrameData, frameRepetition int) []Pulse {
	pulseWave := make([]Pulse, 0)

	pulseWave = append(pulseWave, Pulse{IsHigh: true, Length: WakeUpPulseLength})
	pulseWave = append(pulseWave, Pulse{IsHigh: false, Length: SilenceAfterWakeUpPulseLength})

	hwSyncPulsesFirst := e.hardwareSync(2)
	hwSyncPulsesSecond := e.hardwareSync(7)
	swSyncPulses := e.softwareSync()
	encodedFrame := e.encodeFrameToManchester(frame)
	interFrameGap := e.interFrameGap()

	for repeat := 0; repeat < frameRepetition+1; repeat++ {
		if repeat == 0 {
			pulseWave = append(pulseWave, hwSyncPulsesFirst...)
		} else {
			pulseWave = append(pulseWave, hwSyncPulsesSecond...)
		}

		pulseWave = append(pulseWave, swSyncPulses...)
		pulseWave = append(pulseWave, encodedFrame...)
		pulseWave = append(pulseWave, interFrameGap...)
	}

	return pulseWave
}

func (e *encoder) hardwareSync(repeats int) (pulses []Pulse) {
	for count := 0; count < repeats; count++ {
		pulses = append(pulses, Pulse{IsHigh: true, Length: HardwareSyncPulseLength})
		pulses = append(pulses, Pulse{IsHigh: false, Length: HardwareSyncPulseLength})
	}
	return pulses
}

func (e *encoder) softwareSync() (pulses []Pulse) {
	pulses = append(pulses, Pulse{IsHigh: true, Length: SoftwareSyncPulseLength})
	pulses = append(pulses, Pulse{IsHigh: false, Length: DataHalfPulseLength})
	return pulses
}

func (e *encoder) encodeFrameToManchester(frame FrameData) (pulses []Pulse) {
	for index := 0; index < 56; index++ {
		if (frame[index/8]>>(7-(index%8)))&1 > 0 {
			pulses = append(pulses, Pulse{IsHigh: false, Length: DataHalfPulseLength})
			pulses = append(pulses, Pulse{IsHigh: true, Length: DataHalfPulseLength})
		} else {
			pulses = append(pulses, Pulse{IsHigh: true, Length: DataHalfPulseLength})
			pulses = append(pulses, Pulse{IsHigh: false, Length: DataHalfPulseLength})
		}
	}
	/*for index := 0; index < len(frame); index++ {
		pulses = append(pulses, e.encodeByte(frame[index])...)
	}*/
	return pulses
}

func (e *encoder) encodeByte(data byte) (pulses []Pulse) {
	for bitIndex := 0; bitIndex < 8; bitIndex++ {
		bitValue := data & (1 << bitIndex)
		isHigh := bitValue == 0
		pulses = append(pulses, Pulse{IsHigh: isHigh, Length: DataHalfPulseLength})
		pulses = append(pulses, Pulse{IsHigh: !isHigh, Length: DataHalfPulseLength})
	}
	return pulses
}

func (e *encoder) interFrameGap() (pulses []Pulse) {
	pulses = append(pulses, Pulse{IsHigh: false, Length: InterFrameGapPulseLength})
	return pulses
}
