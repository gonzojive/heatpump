// Package modbusutil provides a wrapper around a modbus.Client that allows
// locking the client.
package modbusutil

import (
	"sync"

	"github.com/goburrow/modbus"
)

// ClientWithLock returns a wrapped modbus client that requires a lock for any
// method calls.
func ClientWithLock(c modbus.Client, l sync.Locker) modbus.Client {
	return &lockableClient{c, l}
}

type lockableClient struct {
	underlying modbus.Client
	lock       sync.Locker
}

// ReadCoils reads from 1 to 2000 contiguous status of coils in a
// remote device and returns coil status.
func (c *lockableClient) ReadCoils(address uint16, quantity uint16) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.ReadCoils(address, quantity)
}

// ReadDiscreteInputs reads from 1 to 2000 contiguous status of
// discrete inputs in a remote device and returns input status.
func (c *lockableClient) ReadDiscreteInputs(address uint16, quantity uint16) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.ReadDiscreteInputs(address, quantity)
}

// WriteSingleCoil write a single output to either ON or OFF in a
// remote device and returns output value.
func (c *lockableClient) WriteSingleCoil(address uint16, value uint16) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.WriteSingleCoil(address, value)
}

// WriteMultipleCoils forces each coil in a sequence of coils to either
// ON or OFF in a remote device and returns quantity of outputs.
func (c *lockableClient) WriteMultipleCoils(address uint16, quantity uint16, value []byte) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.WriteMultipleCoils(address, quantity, value)
}

// 16-bit access
// ReadInputRegisters reads from 1 to 125 contiguous input registers in
// a remote device and returns input registers.
func (c *lockableClient) ReadInputRegisters(address uint16, quantity uint16) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.ReadInputRegisters(address, quantity)
}

// ReadHoldingRegisters reads the contents of a contiguous block of
// holding registers in a remote device and returns register value.
func (c *lockableClient) ReadHoldingRegisters(address uint16, quantity uint16) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.ReadHoldingRegisters(address, quantity)
}

// WriteSingleRegister writes a single holding register in a remote
// device and returns register value.
func (c *lockableClient) WriteSingleRegister(address uint16, value uint16) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.WriteSingleRegister(address, value)
}

// WriteMultipleRegisters writes a block of contiguous registers
// (1 to 123 registers) in a remote device and returns quantity of
// registers.
func (c *lockableClient) WriteMultipleRegisters(address uint16, quantity uint16, value []byte) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.WriteMultipleRegisters(address, quantity, value)
}

// ReadWriteMultipleRegisters performs a combination of one read
// operation and one write operation. It returns read registers value.
func (c *lockableClient) ReadWriteMultipleRegisters(readAddress uint16, readQuantity uint16, writeAddress uint16, writeQuantity uint16, value []byte) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity, value)
}

// MaskWriteRegister modify the contents of a specified holding
// register using a combination of an AND mask, an OR mask, and the
// register's current contents. The function returns
// AND-mask and OR-mask.
func (c *lockableClient) MaskWriteRegister(address uint16, andMask uint16, orMask uint16) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.MaskWriteRegister(address, andMask, orMask)
}

//ReadFIFOQueue reads the contents of a First-In-First-Out (FIFO) queue
// of register in a remote device and returns FIFO value register.
func (c *lockableClient) ReadFIFOQueue(address uint16) (results []byte, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.underlying.ReadFIFOQueue(address)
}
