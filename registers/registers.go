package registers

var Registers = []string{
	"ip", "acc",
	"r1", "r2", "r3", "r4",
	"r5", "r6", "r7", "r8",
	"sp", "fp", "mb", "im",
}

var Map map[string]int

func init() {
	Map = make(map[string]int)
	for i, name := range Registers {
		Map[name] = i
	}
}
