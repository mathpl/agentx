package pdu

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/mathpl/typed"

	log "github.com/Sirupsen/logrus"
)

type HeaderPacket struct {
	Header  *Header
	Context *Context
	Packet  Packet
}

func ReadHeaderPacket(r io.Reader) (*HeaderPacket, error) {
	var err error

	t := time.Now()

	buf := make([]byte, 1024)
	rbh := typed.NewReadBuffer(buf)
	if _, err = rbh.FillFrom(r, HeaderSize); err != nil {
		return nil, err
	}

	hp := &HeaderPacket{}
	hp.Header, err = ReadHeader(rbh)
	if err != nil {
		return nil, err
	}

	rb := typed.NewReadBuffer(buf[HeaderSize:])

	if _, err = rb.FillFrom(r, int(hp.Header.PayloadLength)); err != nil {
		return nil, err
	}

	// Optional context
	if hp.Header.Flags.GetNonDefaultContext() {
		hp.Context = &Context{}
		hp.Context.Read(hp.Header.Flags, rb)
	}

	switch hp.Header.Type {
	case TypeOpen:
		hp.Packet = &Open{}
	case TypeRegister:
		hp.Packet = &Register{}
	case TypePing:
		hp.Packet = &Ping{}
	case TypeAddAgentCaps:
		hp.Packet = &AddAgentCaps{}
	case TypeResponse:
		hp.Packet = &Response{}
	default:
		return nil, fmt.Errorf("Unhandled type: %d", hp.Header.Type)
	}

	err = hp.Packet.Read(hp.Header.Flags, rb)
	if err != nil {
		return nil, err
	}

	// flush padding
	rb.ReadBytes(rb.BytesRemaining())

	log.Debugf("%s Received %d byte packets in %s:\n%s", time.Now(), HeaderSize+hp.Header.PayloadLength, time.Now().Sub(t), hex.Dump(buf[:HeaderSize+hp.Header.PayloadLength]))

	return hp, rb.Err()
}

func (hp *HeaderPacket) Write(bo binary.ByteOrder, w io.Writer) (int, error) {
	t := time.Now()

	buf := make([]byte, 1024)
	wb := typed.NewWriteBuffer(buf)

	hp.Header.Version = 1
	hp.Header.Type = hp.Packet.Type()

	payloadSizeBytes, err := hp.Header.Write(bo, wb)
	if err != nil {
		return -1, err
	}
	headerSize := wb.BytesWritten()

	if hp.Context != nil {
		err = hp.Context.Write(bo, wb)
		if err != nil {
			return -1, err
		}
	}

	err = hp.Packet.Write(bo, wb)
	if err != nil {
		return -1, err
	}

	payloadLen := wb.BytesWritten() - headerSize

	// needs padding
	extra := payloadLen % 4
	for i := extra; i < 4 && extra != 0; i += 1 {
		wb.WriteSingleByte(0)
	}

	payloadSizeBytes.UpdateEndian(bo, uint32(wb.BytesWritten()-headerSize))

	i, err := wb.FlushTo(w)

	log.Debugf("%s Sent %d byte packets in %s:\n%s", time.Now(), wb.BytesWritten(), time.Now().Sub(t).String(), hex.Dump(buf[:wb.BytesWritten()]))

	return i, err
}
