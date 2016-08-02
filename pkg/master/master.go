package master

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mathpl/agentx/pkg/pdu"
	"github.com/mathpl/agentx/pkg/subagent"

	"golang.org/x/net/context"
)

type Master struct {
	registry      map[string]uint32
	registryMutex sync.Mutex

	sessionCnt     uint32
	sessions       map[uint32]*subagent.SubAgent
	cancelSessions map[uint32]func()
	pduChan        chan *pdu.HeaderPacket

	start time.Time
	l     *net.UnixListener
}

func NewMasterAgent(proto, addr string) (*Master, error) {
	m := &Master{registry: make(map[string]uint32),
		sessions:       make(map[uint32]*subagent.SubAgent),
		cancelSessions: make(map[uint32]func()),
		pduChan:        make(chan *pdu.HeaderPacket, 100),
		start:          time.Now(),
	}

	unixAddr, err := net.ResolveUnixAddr(proto, addr)
	if err != nil {
		return nil, err
	}

	m.l, err = net.ListenUnix(proto, unixAddr)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Master) Run(ctx context.Context) error {
	defer m.l.Close()
	defer m.closeAllSessions()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				hp := <-m.pduChan
				if err := m.handlePdu(hp); err != nil {
					log.Errorf("Error handling pdu: %s", err)
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return errors.New("Master Agent canceled by context")
		default:
			// Don't wait forever
			m.l.SetDeadline(time.Now().Add(time.Second))

			conn, err := m.l.Accept()
			if e, ok := err.(net.Error); ok && (e.Timeout() || e.Temporary()) {
				continue
			}

			if err != nil {
				return err
			}

			m.startSession(ctx, wg, conn)
		}
	}

	return nil
}

func (m *Master) Get(ctx context.Context, objectIDs []string) (pdu.Variables, error) {
	reqs := make(map[*subagent.SubAgent]pdu.Ranges)
	for _, objectID := range objectIDs {
		oid, err := pdu.ParseOID(objectID)
		if err != nil {
			return nil, err
		}

		subagent, err := m.getSubAgentByOid(objectID)
		if err != nil {
			log.Warnf("OID without subagent: %s", objectID)
			continue
		}

		if subagent == nil {
			continue
		}

		f := pdu.ObjectIdentifier{Include: 1, Prefix: 0, Subidentifiers: oid}
		t := pdu.ObjectIdentifier{}
		reqs[subagent] = append(reqs[subagent], pdu.Range{From: f, To: t})
	}

	if len(reqs) == 0 {
		return nil, errors.New("No subagent with oid registered")
	}

	timeout := time.After(time.Second * 2)

	cases := make([]reflect.SelectCase, 2, len(reqs)+1)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timeout)}
	cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())}

	for a, pduRanges := range reqs {
		var r pdu.Ranges
		for _, pduRange := range pduRanges {
			r = append(r, pduRange)
		}
		h := &pdu.Header{Version: 1}
		p := &pdu.GetNext{SearchRanges: r}
		get := &pdu.HeaderPacket{Header: h, Packet: p, Context: nil}
		respChan, err := a.Send(get)
		if err != nil {
			log.Errorf("Error sending pdu packet: %s", err)
			continue
		}
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(respChan)})
	}

	var vars pdu.Variables
	remaining := len(reqs)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if ok {
			if chosen == 0 {
				return vars, errors.New("timeout waiting for response(s)")
			}

			if chosen == 1 {
				return vars, errors.New("canceled by context")
			}

			cases[chosen].Chan = reflect.ValueOf(nil)
			remaining -= 1

			p, ok := value.Interface().(*pdu.Response)
			if !ok {
				log.Error("Got empty response")
			}

			vars = append(vars, p.Variables...)
		}
	}

	return vars, nil
}

func (m *Master) getSubAgentByOid(oid string) (*subagent.SubAgent, error) {
	m.registryMutex.Lock()
	defer m.registryMutex.Unlock()

	sessionID, found := m.registry[oid]
	if !found {
		return nil, fmt.Errorf("%s not registered by a subagent", oid)
	}

	subagent, found := m.sessions[sessionID]
	if !found {
		return nil, fmt.Errorf("%s registered by a missing subagent with sessionsID:%d", oid, sessionID)
	}

	return subagent, nil
}

func (m *Master) handlePdu(hp *pdu.HeaderPacket) error {
	var err error
	switch p := hp.Packet.(type) {
	case *pdu.Register:
		pduError := m.register(hp.Header.SessionID, p)

		hp.Header.Flags = 0

		response := &pdu.Response{}
		response.UpTime = time.Now().Sub(m.start)
		response.Error = pduError
		response.Index = 0
		response.Variables = pdu.Variables{pdu.Variable{Type: pdu.VariableTypeNull, Name: p.Subtree}}
		err = m.sendResponse(hp, response)
	case *pdu.AddAgentCaps:
		log.Infof("SessionID:%d %s", hp.Header.SessionID, p.Description)

		response := &pdu.Response{}
		response.Error = pdu.ErrorNone
		response.Index = 0
		response.Variables = pdu.Variables{pdu.Variable{Type: pdu.VariableTypeNull, Name: p.ID}}

		err = m.sendResponse(hp, response)
	}

	return err
}

func (m *Master) register(sessionID uint32, reg *pdu.Register) pdu.Error {
	m.registryMutex.Lock()
	defer m.registryMutex.Unlock()

	foundSessionID, found := m.registry[reg.Subtree.String()]
	if found && foundSessionID != sessionID {
		log.Errorf("OID:%s already registered by session %d", reg.Subtree.String(), foundSessionID)
		return pdu.ErrorDuplicateRegistration
	}

	if !found {
		m.registry[reg.Subtree.String()] = sessionID
		log.Infof("SessionID:%d registered %s", sessionID, reg.Subtree.String())
	}

	return pdu.ErrorNone
}

func (m *Master) unregister(sessionID uint32) pdu.Error {
	m.registryMutex.Lock()
	defer m.registryMutex.Unlock()

	if cancel, found := m.cancelSessions[sessionID]; found {
		cancel()
	}
	delete(m.cancelSessions, sessionID)
	delete(m.sessions, sessionID)

	var oidToDelete []string
	for oid, oidSessionID := range m.registry {
		if sessionID == oidSessionID {
			oidToDelete = append(oidToDelete, oid)
		}
	}

	for _, oid := range oidToDelete {
		delete(m.registry, oid)
	}

	log.Debugf("SessionID:%d unregistered", sessionID)
	return pdu.ErrorNone
}

func (m *Master) sendResponse(hp *pdu.HeaderPacket, resp *pdu.Response) error {
	resp.UpTime = time.Now().Sub(m.start)
	r := &pdu.HeaderPacket{Header: hp.Header, Packet: resp, Context: nil}

	m.registryMutex.Lock()
	subAgent, found := m.sessions[hp.Header.SessionID]
	m.registryMutex.Unlock()
	if !found {
		return fmt.Errorf("No sub agent for session %d found.", hp.Header.SessionID)
	}

	_, err := subAgent.Send(r)
	if err != nil {
		return fmt.Errorf("Unable to send response on sessions %d: %s", hp.Header.SessionID, err)
	}

	return nil
}

func (m *Master) startSession(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
	sessionCtx, cancel := context.WithCancel(ctx)

	m.sessionCnt += 1
	sessionID := m.sessionCnt

	m.registryMutex.Lock()
	defer m.registryMutex.Unlock()
	if _, found := m.sessions[sessionID]; found {
		conn.Close()
		return
	}

	s := subagent.NewSubAgent(conn, sessionID)
	m.sessions[sessionID] = s
	m.cancelSessions[sessionID] = cancel

	go func() {
		log.Infof("Starting session %d", sessionID)
		defer m.unregister(sessionID)
		s.Run(sessionCtx, wg, m.start, m.pduChan)
	}()
}

func (m *Master) closeAllSessions() {
	m.registryMutex.Lock()
	defer m.registryMutex.Unlock()
	for _, cancel := range m.cancelSessions {
		cancel()
	}
	m.cancelSessions = make(map[uint32]func())
	m.sessions = make(map[uint32]*subagent.SubAgent)
	m.registry = make(map[string]uint32)

	close(m.pduChan)
}
