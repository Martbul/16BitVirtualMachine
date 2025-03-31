package memorymapper

import "fmt"

// Interface for devices that can read/write memory.
type MemoryDevice interface {
	GetUint16(address int) uint16
	GetUint8(address int) uint8
	SetUint16(address int, value uint16)
	SetUint8(address int, value uint8)
}

// Region represents a mapped memory region.
type Region struct {
	Device MemoryDevice
	Start  int
	End    int
	Remap  bool
}

type MemoryMapper struct {
	Regions []Region
}

func NewMemoryMapper() *MemoryMapper {
	return &MemoryMapper{
		Regions: []Region{},
	}
}

// Map adds a new memory-mapped region.
func (m *MemoryMapper) Map(device MemoryDevice, start, end int, remap bool) func() {
	region := Region{
		Device: device,
		Start:  start,
		End:    end,
		Remap:  remap,
	}
	m.Regions = append([]Region{region}, m.Regions...)

	// Return an unmap function to remove the region later
	return func() {
		newRegions := []Region{}
		for _, r := range m.Regions {
			if r != region {
				newRegions = append(newRegions, r)
			}
		}
		m.Regions = newRegions
	}
}

// FindRegion locates the memory region for a given address.
func (m *MemoryMapper) FindRegion(address int) (*Region, error) {
	for _, region := range m.Regions {
		if address >= region.Start && address <= region.End {
			return &region, nil
		}
	}
	return nil, fmt.Errorf("no memory region found for address %d", address)
}

// A method for getting the uint16 value at a given address in memory
func (m *MemoryMapper) GetUint16(address int) (uint16, error) { //INFO: remaping the address for the current regeon(if regeion is 0x1000 - 0x1080, that means that for this particular regeorn the address 0x1000 is 0x0000)
	region, err := m.FindRegion(address)
	if err != nil {
		return 0, err
	}

	finalAddress := address
	if region.Remap {
		finalAddress = address - region.Start
	}

	return region.Device.GetUint16(finalAddress), nil
}

// A method for getting the uint8 value at a given address in memory
func (m *MemoryMapper) GetUint8(address int) (uint8, error) {
	region, err := m.FindRegion(address)
	if err != nil {
		return 0, err
	}

	finalAddress := address
	if region.Remap {
		finalAddress = address - region.Start
	}

	return region.Device.GetUint8(finalAddress), nil
}

// A method for setting a uint16 value at a given address in memory
func (m *MemoryMapper) SetUint16(address int, value uint16) error {
	region, err := m.FindRegion(address)
	if err != nil {
		return err
	}

	finalAddress := address
	if region.Remap {
		finalAddress = address - region.Start
	}

	region.Device.SetUint16(finalAddress, value)
	return nil
}

func (m *MemoryMapper) SetUint8(address int, value uint8) error {
	region, err := m.FindRegion(address)
	if err != nil {
		return err
	}

	finalAddress := address
	if region.Remap {
		finalAddress = address - region.Start
	}

	region.Device.SetUint8(finalAddress, value)
	return nil
}
