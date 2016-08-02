package subagent

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mathpl/agentx/pkg/pdu"
	"golang.org/x/net/context"
)

type SubAgent struct {
	c         net.Conn
	sessionID uint32

	openRequests     map[uint32]chan<- *pdu.Response
	openRequestMutex sync.Mutex

	timeout   time.Duration
	byteOrder binary.ByteOrder

	packetCnt uint32
}

func NewSubAgent(c net.Conn, sessionID uint32) *SubAgent {
	a := &SubAgent{c: c, openRequests: make(map[uint32]chan<- *pdu.Response)}
	a.sessionID = sessionID
	a.packetCnt += 100

	if uc, ok := c.(*net.UnixConn); ok {
		uc.SetReadBuffer(0)
		uc.SetWriteBuffer(0)
	}
	return a
}

func (a *SubAgent) Send(hpc *pdu.HeaderPacket) (<-chan *pdu.Response, error) {
	if hpc.Header.PacketID == 0 {
		hpc.Header.PacketID = atomic.AddUint32(&a.packetCnt, 1)
	}
	if hpc.Header.SessionID == 0 {
		hpc.Header.SessionID = a.sessionID
	}
	// hpc.Header.TransactionID should be handled by the client

	// prep open request before sending it
	var respChan chan *pdu.Response
	if hpc.Header.Type != pdu.TypeResponse {
		respChan = make(chan *pdu.Response, 1)
		a.openRequestMutex.Lock()
		a.openRequests[hpc.Header.PacketID] = respChan
		a.openRequestMutex.Unlock()
	}

	a.c.SetWriteDeadline(time.Now().Add(time.Second))
	_, err := hpc.Write(a.byteOrder, a.c)
	if err != nil {
		a.CloseRequest(hpc.Header.PacketID)
		return nil, err
	}

	return respChan, err
}

func (a *SubAgent) Run(ctx context.Context, wg *sync.WaitGroup, start time.Time, pduChan chan<- *pdu.HeaderPacket) {
	defer a.c.Close()
	defer wg.Done()

	wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			log.Infof("SessionID %d Canceled by context", a.sessionID)
			return
		default:
			a.c.SetReadDeadline(time.Now().Add(time.Minute))
			hp, err := pdu.ReadHeaderPacket(a.c)

			if err == io.EOF {
				log.Infof("SessionID %d recevied EOF, terminating.", a.sessionID)
				return
			}

			if err != nil {
				log.Errorf("SessionID %d got an error reading: %s", a.sessionID, err)
				return
			}

			responsePacket := &pdu.Response{UpTime: time.Now().Sub(start)}
			if hp.Header.Type != pdu.TypeOpen && hp.Header.SessionID != a.sessionID {
				log.Errorf("Invalid session for hander. Received:%d Has:%d", hp.Header.SessionID, a.sessionID)
				responsePacket.Error = pdu.ErrorNotOpen
			}

			if responsePacket.Error == pdu.ErrorNone {
				switch p := hp.Packet.(type) {
				case *pdu.Open:
					// Save session settings
					a.timeout = p.Timeout.Duration
					a.byteOrder = hp.Header.Flags.GetByteOrder()

					// Repo
					hp.Header.Flags = 0
					hp.Header.SessionID = a.sessionID
					hp.Context = nil

					v := pdu.Variable{}
					v.Type = pdu.VariableTypeOctetString
					v.Name = p.ID
					v.Value = p.Description.Text
					responsePacket.Variables = pdu.Variables{v}
				case *pdu.Register:
					pduChan <- hp
					continue
				case *pdu.AddAgentCaps:
					pduChan <- hp
					continue
				case *pdu.Response:
					a.handleResponse(hp.Header, p)
					continue
				case *pdu.Ping:
				default:
					log.Errorf("Unhandled: %s", p)
				}
			}

			resp := &pdu.HeaderPacket{Header: hp.Header, Packet: responsePacket, Context: hp.Context}
			_, err = a.Send(resp)
			if err != nil {
				log.Errorf("Error sending pdu: %s", err)
			}
		}
	}
}

func (a *SubAgent) CloseRequest(packetID uint32) {
	a.openRequestMutex.Lock()
	defer a.openRequestMutex.Unlock()
	delete(a.openRequests, packetID)
}

func (a *SubAgent) getResponseChanFromPacketID(packetID uint32) (chan<- *pdu.Response, error) {
	a.openRequestMutex.Lock()
	defer a.openRequestMutex.Unlock()
	respChan, found := a.openRequests[packetID]
	if !found {
		return nil, fmt.Errorf("Unexpected response for sessionsID:%d packetID:%d", a.sessionID, packetID)
	}
	return respChan, nil
}

func (a *SubAgent) handleResponse(h *pdu.Header, r *pdu.Response) {
	respChan, err := a.getResponseChanFromPacketID(h.PacketID)
	if err != nil {
		log.Error(err)
		return
	}
	respChan <- r
	close(respChan)
	a.CloseRequest(h.PacketID)
}
