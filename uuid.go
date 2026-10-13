package main

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// generator produces UUIDv7 values (RFC 9562) that are strictly increasing
// for the lifetime of the process, even when many are requested within the
// same clock tick.
type generator struct {
	// last is the previous timestamp: Unix milliseconds in the upper bits,
	// a 12-bit sub-millisecond fraction in the lower 12.
	last int64
	now  func() time.Time
}

func newGenerator() *generator {
	return &generator{now: time.Now}
}

// next returns a UUIDv7 in canonical 8-4-4-4-12 lowercase form.
//
// Layout: 48 bits of Unix milliseconds, 4 version bits, 12 bits of
// sub-millisecond time (RFC 9562 section 6.2, method 3), 2 variant bits,
// then 62 random bits.
func (g *generator) next() (string, error) {
	var u [16]byte
	if _, err := rand.Read(u[8:]); err != nil {
		return "", err
	}

	ns := g.now().UnixNano()
	ts := (ns/1e6)<<12 | (ns%1e6)<<12/1e6
	if ts <= g.last {
		// Clock didn't advance (or went backwards); bump the fraction,
		// carrying into the millisecond field if it overflows.
		ts = g.last + 1
	}
	g.last = ts

	ms, frac := ts>>12, ts&0xfff
	u[0] = byte(ms >> 40)
	u[1] = byte(ms >> 32)
	u[2] = byte(ms >> 24)
	u[3] = byte(ms >> 16)
	u[4] = byte(ms >> 8)
	u[5] = byte(ms)
	u[6] = 0x70 | byte(frac>>8)
	u[7] = byte(frac)
	u[8] = 0x80 | u[8]&0x3f

	var s [36]byte
	hex.Encode(s[0:8], u[0:4])
	s[8] = '-'
	hex.Encode(s[9:13], u[4:6])
	s[13] = '-'
	hex.Encode(s[14:18], u[6:8])
	s[18] = '-'
	hex.Encode(s[19:23], u[8:10])
	s[23] = '-'
	hex.Encode(s[24:36], u[10:16])
	return string(s[:]), nil
}
