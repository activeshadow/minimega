// supports completely decentralized message passing, both to a set of nodes 
// as well as broadcast. Meshage is design for resiliency, and automatically 
// updates routes and topologies when nodes in the mesh fail. Meshage 
// automatically maintains density health - as nodes leave the mesh, adjacent 
// nodes will connect to others in the mesh to maintain a minimum degree for 
// resiliency. 
// 
// Meshage is decentralized - Any node in the mesh is capable of initiating and
// receiving messages of any type. This also means that any node is capable of 
// issuing control messages that affect the topology of the mesh.
// 
// Meshage is secure and resilient - All messages are signed and encrypted by 
// the sender to guarantee authenticity and integrity. Nodes on the network 
// store public keys of trusted agents, who may send messages signed and 
// encrypted with a corresponding private key. This is generally done by the 
// end user. Compromised nodes on the mesh that attempt denial of service 
// through discarding messages routed through them are automatically removed 
// from the network by neighbor nodes.  
package meshage

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	log "minilog"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	RECEIVE_BUFFER = 1024
	REBROADCAST_WAIT = 10
)

const (
	SET = iota
	BROADCAST
)

const (
	UNION = iota
	INTERSECTION
	MESSAGE
	HANDSHAKE
	HANDSHAKE_SOLICITED
	ACK
)

// A Node object contains the network information for a given node. Creating a 
// Node object with a non-zero degree will cause it to begin broadcasting for 
// connections automatically.
type Node struct {
	name        string              // node name. Must be unique on a network.
	degree      uint                // degree for this node, set to 0 to force node to not broadcast
	mesh        map[string][]string // adjacency list for the known topology for this node
	routes      map[string]string   // one-hop routes for every node on the network, including this node
	receive     chan Message        // channel of incoming messages. A program will read this channel for incoming messages to this node

	clients     map[string]*client // list of connections to this node
	meshLock sync.Mutex
	clientLock sync.Mutex
	degreeLock  sync.Mutex
	messagePump chan Message
	port	int
	timeout time.Duration

	errors chan error
}

// A Message is the payload for all message passing, and contains the user 
// specified message in the Body field.
type Message struct {
	Type	int
	Recipients   []string    // list of recipients if MessageType = MESSAGE_SET, undefined if broadcast
	Source       string      // name of source node
	CurrentRoute []string    // list of hops for an in-flight message
	ID           uint64      // sequence id
	Command      int         // union, intersection, message
	Body         interface{} // message body
}

func init() {
	gob.Register(map[string][]string{})
}

// NewNode returns a new node and receiver channel with a given name and 
// degree. If degree is non-zero, the node will automatically begin 
// broadcasting for connections.
func NewNode(name string, degree uint, port int, timeout int) (*Node, chan Message, chan error) {
	n := &Node{
		name:        name,
		degree:      degree,
		mesh:        make(map[string][]string),
		routes:      make(map[string]string),
		receive:     make(chan Message, RECEIVE_BUFFER),
		clients:     make(map[string]*client),
		messagePump: make(chan Message, RECEIVE_BUFFER),
		port: port,
		errors:      make(chan error),
		timeout: time.Duration(timeout) * time.Second,
	}
	go n.connectionListener()
	go n.broadcastListener()
	go n.messageHandler()
	go n.checkDegree()

	return n, n.receive, n.errors
}

// check degree emits connection requests when our number of connected clients is below the degree threshold
func (n *Node) checkDegree() {
	// check degree only if we're not already running
	n.degreeLock.Lock()
	defer n.degreeLock.Unlock()

	var backoff uint = 1
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for uint(len(n.clients)) < n.degree {
		log.Debugln("soliciting connections")
		b := net.IPv4(255, 255, 255, 255)
		addr := net.UDPAddr{
			IP:   b,
			Port: n.port,
		}
		socket, err := net.DialUDP("udp4", nil, &addr)
		if err != nil {
			log.Errorln(err)
			n.errors <- err
			break
		}
		message := fmt.Sprintf("meshage:%s", n.name)
		_, err = socket.Write([]byte(message))
		if err != nil {
			log.Errorln(err)
			n.errors <- err
			break
		}
		wait := r.Intn(1 << backoff)
		time.Sleep(time.Duration(wait) * time.Second)
		if backoff < 7 { // maximum wait won't exceed 128 seconds
			backoff++
		}
	}
}

// broadcastListener listens for broadcast connection requests and attempts to connect to that node
func (n *Node) broadcastListener() {
	listenAddr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: n.port,
	}
	ln, err := net.ListenUDP("udp4", &listenAddr)
	if err != nil {
		log.Errorln(err)
		n.errors <- err
		return
	}
	for {
		d := make([]byte, 1024)
		read, _, err := ln.ReadFromUDP(d)
		data := strings.Split(string(d[:read]), ":")
		if len(data) != 2 {
			err = fmt.Errorf("gor malformed udp data: %v\n", data)
			log.Errorln(err)
			n.errors <- err
			continue
		}
		if data[0] != "meshage" {
			err = fmt.Errorf("got malformed udp data: %v\n", data)
			log.Errorln(err)
			n.errors <- err
			continue
		}
		host := data[1]
		if host == n.name {
			log.Debugln("got solicitation from myself, dropping")
			continue
		}
		log.Debug("got solicitation from %v\n", host)
		go n.dial(host, true)
	}
}

// connectionListener accepts incoming connections and hands new connections to a connection handler
func (n *Node) connectionListener() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", n.port))
	if err != nil {
		n.errors <- err
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Errorln(err)
			n.errors <- err
			continue
		}
		n.handleConnection(conn)
	}
}

// handleConnection creates a new client and issues a handshake. It adds the client to the list
// of clients only after a successful handshake
func (n *Node) handleConnection(conn net.Conn) {
	c := newClient(conn, n.timeout)

	var solicited bool
	if uint(len(n.clients)) < n.degree {
		solicited = true
	} else {
		solicited = false
	}

	if ok, err := c.sendHandshake(solicited, n.name, n.mesh); ok {
		// valid connection, add it to the client roster
		n.clientLock.Lock()
		defer n.clientLock.Unlock()

		n.clients[c.name] = c

		go n.receiveHandler(c)
	} else {
		if err != nil {
			n.errors <- err
		}
	}
}

func (n *Node) receiveHandler(c *client) {
	for {
		m, err := c.receive()
		if err != nil {
			log.Debugln("disconnecting from client")
			break
		}
		log.Debug("receiveHandler got: %#v\n", m)
		n.messagePump <- m
	}

	n.clientLock.Lock()

	// remove the client from our client list, and broadcast an intersection announcement about this connection
	delete(n.clients, c.name)

	mesh := make(map[string][]string)
	mesh[n.name] = []string{c.name}
	mesh[c.name] = []string{n.name}
	n.intersect(mesh)
	n.clientLock.Unlock()

	// let everyone know about the new topology
	u := Message{
		Source:       n.name,
		CurrentRoute: []string{n.name},
		ID:           n.sequenceID(),
		Command:      INTERSECTION,
		Body:         mesh,
	}
	log.Debug("receiveHandler broadcasting topology: %v\n", u)
	n.broadcast(u)

	// make sure we keep up the necessary degree
	n.checkDegree()
}

// Dial connects a node to another, regardless of degree. Returned error is nil 
// if successful.
func (n *Node) Dial(addr string) {
	n.dial(addr, false)
}

// Hangup disconnects from a connected client and announces the disconnect to the
// topology.
func (n *Node) Hangup(client string) error {
	c, ok := n.clients[client]
	if !ok {
		return fmt.Errorf("no such client")
	}
	c.hangup()
	return nil
}

func (n *Node) dial(host string, solicited bool) {
	addr := fmt.Sprintf("%s:%d", host, n.port)
	log.Debug("Dialing: %v\n", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		n.errors <- err
	}
	c := newClient(conn, n.timeout)

	if ok, mesh, err := c.recvHandshake(solicited, n.clients, n.name); ok {
		n.clientLock.Lock()
		n.clients[c.name] = c
		go n.receiveHandler(c)

		// add this new connection to the mesh and union with our mesh
		mesh[n.name] = append(mesh[n.name], c.name)
		mesh[c.name] = append(mesh[c.name], n.name)
		n.union(mesh)
		n.clientLock.Unlock()
		// let everyone know about the new topology
		u := Message{
			Source:       n.name,
			CurrentRoute: []string{n.name},
			ID:           n.sequenceID(),
			Command:      UNION,
			Body:         n.mesh,
		}
		log.Debug("dial broadcasting topology: %#v\n", u)
		n.broadcast(u)
	} else {
		if err != nil {
			n.errors <- err
		}
	}
}

// SetDegree sets the degree for a given node. Setting degree == 0 will cause the 
// node to stop broadcasting for connections.
func (n *Node) SetDegree(d uint) {
	n.degree = d
}

// Degree returns the current degree
func (n *Node) Degree() uint {
	return n.degree
}

// union merges a mesh with the local one and eliminates redundant connections
// union can also generate intersection messages - it checks the client list
// to ensure that union messages do not alter what it knows about its own 
// connections. If a discrepancy is found, it broadcasts an intersection to
// fix the discrepancy.
func (n *Node) union(m map[string][]string) {
	n.meshLock.Lock()
	defer n.meshLock.Unlock()

	log.Debug("union mesh: %v\n", m)

	n.dropRoutes()

	// merge everything, sort each bin, and eliminate duplicate entries
	for k, v := range m {
		n.mesh[k] = append(n.mesh[k], v...)
		sort.Strings(n.mesh[k])
		var nl []string
		for _, j := range n.mesh[k] {
			if len(nl) == 0 {
				nl = append(nl, j)
				continue
			}
			if nl[len(nl)-1] != j {
				nl = append(nl, j)
			}
		}
		n.mesh[k] = nl
	}
	log.Debug("new mesh is: %v\n", n.mesh)

	// check to make sure that our client list matches the connections
	// listed in the mesh
	intersection_mesh := make(map[string][]string)
	for _, v := range n.mesh[n.name] {
		if _, ok := n.clients[v]; !ok {
			intersection_mesh[n.name] = append(intersection_mesh[n.name], v)
			intersection_mesh[v] = append(intersection_mesh[v], n.name)
		}
	}

	if len(intersection_mesh) != 0 {
		n.intersect_locked(intersection_mesh)
		go func() {
			u := Message{
				Source:       n.name,
				CurrentRoute: []string{n.name},
				ID:           n.sequenceID(),
				Command:      INTERSECTION,
				Body:         intersection_mesh,
			}
			log.Debug("found union conflicts, broadcasting new intersection %v\n", intersection_mesh)
			s := rand.NewSource(time.Now().UnixNano())
			r := rand.New(s)
			wait := r.Intn(REBROADCAST_WAIT)
			time.Sleep(time.Duration(wait) * time.Second)
			n.broadcast(u)
		}()
	}

	// if the mesh we're unioning isn't exactly the same as our current mesh, then broadcast again
	if !n.meshCompare(m) {
		u := Message{
			Source: n.name,
			CurrentRoute: []string{n.name},
			ID: n.sequenceID(),
			Command: UNION,
			Body: n.mesh,
		}
		log.Debugln("union rebroadcast")
		go func() {
			s := rand.NewSource(time.Now().UnixNano())
			r := rand.New(s)
			wait := r.Intn(REBROADCAST_WAIT)
			time.Sleep(time.Duration(wait) * time.Second)
			n.broadcast(u)
		}()
	}
}

// compare a mesh to our internal one and return true if they match
func (n *Node) meshCompare(m map[string][]string) bool {
	if len(n.mesh) != len(m) {
		return false
	}
	for k,v := range m {
		if vn, ok := n.mesh[k]; ok {
			sort.Strings(vn)
			sort.Strings(v)
			if len(vn) != len(v) {
				return false
			}
			for i, j := range v {
				if vn[i] != j {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}
// intersect (this isn't actually an intersection function...) removes the 
// nodes given from the topology.
func (n *Node) intersect(m map[string][]string) {
	n.meshLock.Lock()
	n.intersect_locked(m)
	n.meshLock.Unlock()
}

func (n *Node) intersect_locked(m map[string][]string) {
	log.Debug("intersect mesh: %v\n", m)
	n.dropRoutes()
	for k, v := range m {
		// remove all of v from key k
		var nv []string
		for _, x := range n.mesh[k] {
			found := false
			for _, y := range v {
				if x == y {
					found = true
					break
				}
			}
			if !found {
				nv = append(nv, x)
			}
		}
		n.mesh[k] = nv

		// if key k is now empty, then remove key k
		if len(n.mesh[k]) == 0 {
			delete(n.mesh, k)
			if k == n.name {
				log.Debug("disconnected from all clients! %#v\n", n.clients)
				n.mesh = make(map[string][]string)
			}
		}
	}
	log.Debug("new mesh is: %v\n", n.mesh)
}

// Send a message according to the parameters set in the message. 
// Users will generally use the Set and Broadcast methods instead of Send.
func (n *Node) SendMessage(m Message) {
	n.clientLock.Lock()
	defer n.clientLock.Unlock()

	// we want to duplicate the message for each slice of recipients that follow a like route from this node
	route_slices := make(map[string][]string)

	log.Debug("sending message to %d clients\n", len(m.Recipients))
	for _, v := range m.Recipients {
		log.Debug("sending to %v\n", v)

		// don't send to ourselves
		if v == n.name {
			continue
		}

		// make sure we have a route to this client
		var route string
		var ok bool
		if route, ok = n.routes[v]; !ok {
			n.updateRoute(v)
			if route, ok = n.routes[v]; !ok {
				err := fmt.Errorf("no route to host: %v", v)
				log.Errorln(err)
				n.errors <- err
				continue
			}
		}
		route_slices[route] = append(route_slices[route], v)
	}

	log.Debug("route slices: %#v\n", route_slices)
	for k, v := range route_slices {
		m.Recipients = v
		// get the client for this route
		if c, ok := n.clients[k]; ok {
			go n.sendOne(c, m)
		} else {
			err := fmt.Errorf("mismatched client list and topology, something is very wrong: %v, %#v", v, n.clients)
			log.Errorln(err)
			n.errors <- err
		}
	}
}

// broadcastSend sends a broadcast message to all connected clients
func (n *Node) broadcast(m Message) {
	m.Recipients = []string{}
	for k, _ := range n.mesh {
		log.Debug("adding broadcast recipient: %v\n", k)
		m.Recipients = append(m.Recipients, k)
	}
	n.SendMessage(m)
}

func (n *Node) sendOne(c *client, m Message) {
	err := c.send(m)
	if err != nil {
		log.Errorln(err)
		n.errors <- err
	}
}

// Send a message to a list of recipients.
func (n *Node) Send(recipients []string, body interface{}) {
	u := Message{
		Source:       n.name,
		Recipients:   recipients,
		CurrentRoute: []string{n.name},
		ID:           n.sequenceID(),
		Command:      MESSAGE,
		Body:         body,
	}
	log.Debug("send message %#v\n", u)
	n.SendMessage(u)
}

// Broadcast sends a broadcast message to all connected nodes.
func (n *Node) Broadcast(body interface{}) {
	u := Message{
		Source:       n.name,
		CurrentRoute: []string{n.name},
		ID:           n.sequenceID(),
		Command:      MESSAGE,
		Body:         body,
	}
	log.Debug("broadcasting message %#v\n", u)
	n.broadcast(u)
}

// Return a sequence ID for this node and automatically increment the ID
func (n *Node) sequenceID() uint64 {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	id := uint64(r.Int63())
	log.Debug("set id: %v", id)
	return id
}

// messageHandler receives messages on a channel from any clients and processes them.
// Some messages are rebroadcast, or sent along other routes. Messages intended for this
// node are sent along the receive channel to the user.
func (n *Node) messageHandler() {
	for {
		m := <-n.messagePump
		log.Debug("messageHandler: %#v\n", m)
		m.CurrentRoute = append(m.CurrentRoute, n.name)

		// do we also handle it?
		// TODO: is there a better way to slice up this slice?
		var new_recipients []string
		for _, i := range m.Recipients {
			if i == n.name {
				go n.handleMessage(m)
			} else {
				new_recipients = append(new_recipients, i)
			}
		}
		m.Recipients = new_recipients

		go n.SendMessage(m)
	}
}

// handleMessage parses a message intended for this node.
// If the message is a control message, we process it here, if it's
// a regular message, we put it on the receive channel.
func (n *Node) handleMessage(m Message) {
	log.Debug("handleMessage: %#v\n", m)
	switch m.Command {
	case UNION:
		n.union(m.Body.(map[string][]string))
	case INTERSECTION:
		n.intersect(m.Body.(map[string][]string))
	case MESSAGE:
		n.receive <- m
	default:
		n.errors <- fmt.Errorf("handleMessage: invalid message type")
	}
}

// Mesh returns an adjacency list containing the known mesh. The adjacency list
// is a map[string][]string containing all connections to a node given as the
// key.
// The returned map is a copy of the internal mesh, and modifying is will not
// affect the mesh.
func (n *Node) Mesh() map[string][]string {
	n.meshLock.Lock()
	defer n.meshLock.Unlock()

	ret := make(map[string][]string)
	for k, v := range n.mesh {
		ret[k] = v
	}
	return ret
}
