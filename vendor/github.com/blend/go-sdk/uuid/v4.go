package uuid

import "crypto/rand"

// V4 Create a new UUID version 4.
func V4() UUID {
	uuid := Empty()
	rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // set version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // set variant 2
	return uuid
}
