package devices

import "fmt"

// ScreenDevice represents a simple terminal screen
type ScreenDevice struct {
}

// CreateScreenDevice initializes and returns a new ScreenDevice
func CreateScreenDevice() *ScreenDevice {
	return &ScreenDevice{}
}

// GetUint16 returns a dummy value (not used in this device)
func (s *ScreenDevice) GetUint16(address int) uint16 {
	return 0
}

// GetUint8 returns a dummy value (not used in this device)
func (s *ScreenDevice) GetUint8(address int) uint8 {
	return 0
}
func (s *ScreenDevice) SetUint8(address int, value uint8) {
	// Just printing the value for now, you can modify as needed
	fmt.Printf("SetUint8 called at address: 0x%X with value: %X\n", address, value)
}

// SetUint16 processes commands and writes characters to the screen
func (s *ScreenDevice) SetUint16(address int, data uint16) {
	command := (data & 0xFF00) >> 8
	characterValue := data & 0x00FF

	// Process commands
	switch command {
	case 0xFF:
		eraseScreen()
	case 0x01:
		setBold()
	case 0x02:
		setRegular()
	}

	// Calculate screen position
	x := (address % 16) + 1
	y := (address / 16) + 1
	moveTo(x*2, y)

	// Print character
	character := rune(characterValue)
	fmt.Print(string(character))
}

// eraseScreen clears the terminal screen
func eraseScreen() {
	fmt.Print("\x1b[2J") //INFO: Escape code
}

// moveTo moves the cursor to a specific (x, y) position
func moveTo(x, y int) {
	fmt.Printf("\x1b[%d;%dH", y, x)
}

// setBold sets the text to bold
func setBold() {
	fmt.Print("\x1b[1m")
}

// setRegular resets text formatting to regular
func setRegular() {
	fmt.Print("\x1b[0m")
}
