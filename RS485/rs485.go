// Copyright 2021 Abhijit Bose. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
// Use of this source code is governed by a Apache 2.0 license that can be found
// in the LICENSE file.

// Package rs485 provides additional control functionality over the serial port driver
//
// Typically the RTS signal is used to provide as a mechanism to control
// transmit / receive enable. This package helps to achieve this.
// Note: That this package only support half duplex RS485 links only.

package rs485

import (
	"fmt"
	"time"

	serial "github.com/boseji/goSerialPort"
)

// Control is a type of function that can receive a boolean value and
// apply that to a paricular serial port handshake pin (DTR / RTS).
type Control func(bool) error

// Port provides a way to control the Halfduplex communication on RS485
type Port struct {
	port        serial.Port
	delayBefore time.Duration
	delayAfter  time.Duration
	sig         Control
}

// New creates a new Port that is configured for the timing and signalling
// requirements for halfduplex RS485
func New(port serial.Port, delayBefore, delayAfter time.Duration, sig Control) *Port {
	if port == nil {
		return nil
	}
	return &Port{
		port:        port,
		delayBefore: delayBefore,
		delayAfter:  delayAfter,
		sig:         sig,
	}
}

// Write implemantion of io.Writer interface
func (p *Port) Write(b []byte) (n int, err error) {
	if p == nil || len(b) == 0 {
		return 0, fmt.Errorf("failed to write empty / un-initialized port")
	}

	// Startup
	p.sig(true) // Activate the Signal
	if p.delayBefore != 0 {
		time.Sleep(p.delayBefore)
	}

	// Transmit
	n, err = p.port.Write(b)

	// End
	defer func() {
		if p.delayAfter != 0 {
			time.Sleep(p.delayAfter)
		}
		p.sig(false) // Deactivate
	}()

	return
}

// Close implementation of io.Closer interface
func (p *Port) Close() error {
	if p == nil {
		return serial.ErrNotOpen
	}
	return p.port.Close()
}

// Read implementation of io.Reader interface
func (p *Port) Read(b []byte) (n int, err error) {
	if p == nil || len(b) == 0 {
		return 0, fmt.Errorf("failed to read empty / un-initialized port")
	}

	// Startup
	p.sig(false) // Activate the Signal

	n, err = p.port.Read(b)
	return
}
