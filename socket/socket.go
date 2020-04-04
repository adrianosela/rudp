package socket

import (
	"fmt"
	"net"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/atc"
	"github.com/adrianosela/rdtp/netwk"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/pckfactory"
	"github.com/pkg/errors"
)

// Socket represents a socket abstraction and carries all
// necessary info and statistics about the socket
type Socket struct {
	lAddr *rdtp.Addr // local rdtp address
	rAddr *rdtp.Addr // remote rdtp address

	txBytes uint32 // current sequence number
	rxBytes uint32 // current ack number

	atc *atc.AirTrafficCtrl
	pf  *pckfactory.PacketFactory

	In          chan *packet.Packet
	application net.Conn
}

// NewSocket returns a newly allocated socket
func NewSocket(lAddr, rAddr *rdtp.Addr, nw *netwk.Network, c net.Conn) (*Socket, error) {

	atctrl := atc.NewAirTrafficCtrl(func(p *packet.Packet) {
		nw.Send(rAddr.Host, p)
	})

	pf, err := pckfactory.New(
		uint16(lAddr.Port),
		uint16(rAddr.Port),
		func(p *packet.Packet) error {
			atctrl.Send(p)
			return nil
		},
		packet.MaxPayloadBytes)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize new packetfactory")
	}

	s := &Socket{
		lAddr:       lAddr,
		rAddr:       rAddr,
		application: c,
		atc:         atctrl,
		pf:          pf,
		In:          make(chan *packet.Packet),
	}

	go s.receiver()
	go s.sender()

	return s, nil
}

// ID returns the of unique identifier of the socket
func (s *Socket) ID() string {
	return fmt.Sprintf("%s %s", s.lAddr.String(), s.rAddr.String())
}

// LocalAddr returns the local network address.
func (s *Socket) LocalAddr() net.Addr {
	return s.lAddr
}

// RemoteAddr returns the remote network address.
func (s *Socket) RemoteAddr() net.Addr {
	return s.rAddr
}

// Close closes a socket
func (s *Socket) Close() error {
	return s.application.Close()
}

func (s *Socket) receiver() {
	for {
		p := <-s.In

		s.atc.Ack(p.AckNo)             // acknowledge received packet
		s.rxBytes += uint32(p.Length)  // keep track of stats
		s.application.Write(p.Payload) // pass packet to application layer
	}
}

func (s *Socket) sender() {
	buf := make([]byte, 1500)
	for {
		n, err := s.application.Read(buf)
		if err != nil {
			return
		}
		n, err = s.pf.Send(buf[:n])
		if err != nil {
			return
		}
		s.txBytes += uint32(n)
	}
}
