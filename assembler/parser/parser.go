// mov $42, r1 -> (instruction literal hex value, register)

// move [$42 + (!loc - $1F)], r1 -> (instruction [literal hex value + (!variable - )])
package parser

import (
	"github.com/alecthomas/participle/v2"
)
