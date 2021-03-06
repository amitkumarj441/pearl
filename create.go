package pearl

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/mmcloughlin/pearl/buf"
	"github.com/mmcloughlin/pearl/check"
	"github.com/mmcloughlin/pearl/log"
	"github.com/mmcloughlin/pearl/ntor"
	"github.com/mmcloughlin/pearl/torcrypto"
	"github.com/pkg/errors"
)

// HandshakeType is an identifier for a circuit handshake type.
type HandshakeType uint16

// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L877-L880
//
//	   Recognized handshake types are:
//	       0x0000  TAP  -- the original Tor handshake; see 5.1.3
//	       0x0001  reserved
//	       0x0002  ntor -- the ntor+curve25519+sha256 handshake; see 5.1.4
//
const (
	HandshakeTypeTAP  HandshakeType = 0
	HandshakeTypeNTOR HandshakeType = 2
)

// Recognized HTAG values.
//
// Reference: https://github.com/torproject/torspec/blob/f9eeae509344dcfd1f185d0130a0055b00131cea/tor-spec.txt#L892-L894
//
//	   migration. See 5.1.2.1 below. Recognized HTAG values are:
//
//	       ntor -- 'ntorNTORntorNTOR'
//
const (
	HandshakeTagNTOR = "ntorNTORntorNTOR"
)

// TAP handshake data sizes.
//
// Reference: https://github.com/torproject/torspec/blob/f9eeae509344dcfd1f185d0130a0055b00131cea/tor-spec.txt#L1024-L1025
//
//	   Define TAP_C_HANDSHAKE_LEN as DH_LEN+KEY_LEN+PK_PAD_LEN.
//	   Define TAP_S_HANDSHAKE_LEN as DH_LEN+HASH_LEN.
//
const (
	HandshakeTAPClientLength = torcrypto.DiffieHellmanPublicSize + torcrypto.StreamCipherKeySize + torcrypto.PublicKeyPaddingSize
	HandshakeTAPServerLength = torcrypto.DiffieHellmanPublicSize + torcrypto.HashSize
)

// CreatedCell represents a CREATED cell.
//
// Reference: https://github.com/torproject/torspec/blob/0fd44031bfd6c6c822bfb194e54a05118c9625e2/tor-spec.txt#L896-L897
//
//	   The format of a CREATED cell is:
//	       HDATA     (Server Handshake Data)     [TAP_S_HANDSHAKE_LEN bytes]
//
type CreatedCell struct {
	CircID        CircID
	HandshakeData []byte
}

func (cell *CreatedCell) UnmarshalCell(c Cell) error {
	if c.Command() != CommandCreated {
		return ErrUnexpectedCommand
	}

	p := c.Payload()
	if len(p) < HandshakeTAPServerLength {
		return ErrShortCellPayload
	}

	cell.CircID = c.CircID()
	cell.HandshakeData = p[:HandshakeTAPServerLength]

	return nil
}

func (cell *CreatedCell) Payload() []byte {
	return cell.HandshakeData
}

// Create2Cell represents a CREATE2 cell.
type Create2Cell struct {
	CircID        CircID
	HandshakeType HandshakeType
	HandshakeData []byte
}

// ParseCreate2Cell parses a CREATE2 cell.
func ParseCreate2Cell(c Cell) (*Create2Cell, error) {
	// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L868-L871
	//
	//	   A CREATE2 cell contains:
	//	       HTYPE     (Client Handshake Type)     [2 bytes]
	//	       HLEN      (Client Handshake Data Len) [2 bytes]
	//	       HDATA     (Client Handshake Data)     [HLEN bytes]
	//
	if c.Command() != CommandCreate2 {
		return nil, ErrUnexpectedCommand
	}

	payload := c.Payload()
	n := len(payload)

	if n < 4 {
		return nil, errors.New("create2 cell too short")
	}

	htype := binary.BigEndian.Uint16(payload)
	hlen := binary.BigEndian.Uint16(payload[2:])

	if n < int(4+hlen) {
		return nil, errors.New("inconsistent create2 cell length")
	}

	return &Create2Cell{
		CircID:        c.CircID(),
		HandshakeType: HandshakeType(htype),
		HandshakeData: payload[4 : 4+hlen],
	}, nil
}

// Cell builds a cell from the CREATE2 payload.
func (c Create2Cell) Cell() (Cell, error) {
	hlen := len(c.HandshakeData)
	cell := NewFixedCell(c.CircID, CommandCreate2)
	payload := cell.Payload()

	binary.BigEndian.PutUint16(payload, uint16(c.HandshakeType))
	binary.BigEndian.PutUint16(payload[2:], uint16(hlen))
	copy(payload[4:], c.HandshakeData)

	return cell, nil
}

// Created2Cell represents a CREATED2 cell.
//
// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L873-L875
//
//	   A CREATED2 cell contains:
//	       HLEN      (Server Handshake Data Len) [2 bytes]
//	       HDATA     (Server Handshake Data)     [HLEN bytes]
//
type Created2Cell struct {
	CircID        CircID
	HandshakeData []byte
}

func (cell *Created2Cell) UnmarshalCell(c Cell) error {
	if c.Command() != CommandCreated2 {
		return ErrUnexpectedCommand
	}

	p := c.Payload()
	n := len(p)

	if n < 2 {
		return errors.New("created2 cell too short")
	}

	hlen := binary.BigEndian.Uint16(p)

	if n < int(2+hlen) {
		return errors.New("inconsistent created2 cell length")
	}

	cell.CircID = c.CircID()
	cell.HandshakeData = p[2 : 2+hlen]

	return nil
}

// Payload returns just the payload part of the CREATED2 cell.
func (c Created2Cell) Payload() []byte {
	n := len(c.HandshakeData)
	p := make([]byte, 2+n)
	binary.BigEndian.PutUint16(p, uint16(n))
	copy(p[2:], c.HandshakeData)
	return p
}

// Cell builds a cell from the CREATED2 payload.
func (c Created2Cell) Cell() (Cell, error) {
	cell := NewFixedCell(c.CircID, CommandCreated2)
	payload := cell.Payload()

	hlen := len(c.HandshakeData)
	binary.BigEndian.PutUint16(payload, uint16(hlen))
	copy(payload[2:], c.HandshakeData)

	return cell, nil
}

// REVIEW(mbm): CreateRequest as an interface?
type CreateRequest struct {
	CircID        CircID
	CreateType    Command
	HandshakeType HandshakeType
	HandshakeData []byte
}

// CreateFastHandler handles a received CREATE_FAST cell.
func CreateFastHandler(conn *Connection, c Cell) error {
	if c.Command() != CommandCreateFast {
		return ErrUnexpectedCommand
	}

	// Reference: https://github.com/torproject/torspec/blob/f66d1826c0b32d307898bba081dbf8ef598d4037/tor-spec.txt#L1139-L1141
	//
	//	   A CREATE_FAST cell contains:
	//
	//	       Key material (X)    [HASH_LEN bytes]
	//

	p := c.Payload()
	if len(p) < torcrypto.HashSize {
		return ErrShortCellPayload
	}
	X := p[:torcrypto.HashSize]

	// Reference: https://github.com/torproject/torspec/blob/f66d1826c0b32d307898bba081dbf8ef598d4037/tor-spec.txt#L1174-L1175
	//
	//	   If CREATE_FAST is used, both parties base their key material on
	//	   K0=X|Y.
	//

	Y := torcrypto.Rand(torcrypto.HashSize)
	s := append(X, Y...)

	k, err := BuildCircuitKeysKDFTOR(s)
	if err != nil {
		return errors.Wrap(err, "failed to build circuit keys")
	}

	err = LaunchCircuit(conn, c.CircID(), k)
	if err != nil {
		return errors.Wrap(err, "failed to launch circuit")
	}

	// Reference: https://github.com/torproject/torspec/blob/f66d1826c0b32d307898bba081dbf8ef598d4037/tor-spec.txt#L1143-L1148
	//
	//	   A CREATED_FAST cell contains:
	//
	//	       Key material (Y)    [HASH_LEN bytes]
	//	       Derivative key data [HASH_LEN bytes] (See 5.2.1 below)
	//
	//	   The values of X and Y must be generated randomly.
	//

	cell := NewFixedCell(c.CircID(), CommandCreatedFast)
	p = cell.Payload()
	copy(p, Y)
	copy(p[torcrypto.HashSize:], k.KH)

	err = conn.SendCell(cell)
	if err != nil {
		return errors.Wrap(err, "could not send created cell")
	}

	conn.logger.Info("circuit created")

	return nil
}

// CreateHandler handles a received CREATE cell.
func CreateHandler(conn *Connection, c Cell) error {
	if c.Command() != CommandCreate {
		return ErrUnexpectedCommand
	}

	// Reference: https://github.com/torproject/torspec/blob/f9eeae509344dcfd1f185d0130a0055b00131cea/tor-spec.txt#L883-L887
	//
	//	   The format of a CREATE cell is one of the following:
	//	       HDATA     (Client Handshake Data)     [TAP_C_HANDSHAKE_LEN bytes]
	//	   or
	//	       HTAG      (Client Handshake Type Tag) [16 bytes]
	//	       HDATA     (Client Handshake Data)     [TAP_C_HANDSHAKE_LEN-16 bytes]
	//

	req := CreateRequest{
		CircID:     c.CircID(),
		CreateType: CommandCreate,
	}

	p := c.Payload()
	if bytes.HasPrefix(p, []byte(HandshakeTagNTOR)) {
		req.HandshakeType = HandshakeTypeNTOR
		req.HandshakeData = p[16:HandshakeTAPClientLength]
	} else {
		req.HandshakeType = HandshakeTypeTAP
		req.HandshakeData = p[:HandshakeTAPClientLength]
	}

	return ProcessHandshake(conn, req)
}

// Create2Handler handles a received CREATE2 cell.
func Create2Handler(conn *Connection, c Cell) error {
	cr, err := ParseCreate2Cell(c)
	if err != nil {
		return errors.Wrap(err, "failed to parse create2 cell")
	}

	req := CreateRequest{
		CircID:        c.CircID(),
		CreateType:    CommandCreate,
		HandshakeType: cr.HandshakeType,
		HandshakeData: cr.HandshakeData,
	}

	return ProcessHandshake(conn, req)
}

// ProcessHandshake directs the creation request to the correct handshake
// mechanism.
func ProcessHandshake(conn *Connection, req CreateRequest) error {
	switch req.HandshakeType {
	case HandshakeTypeTAP:
		return ProcessHandshakeTAP(conn, req)
	case HandshakeTypeNTOR:
		return ProcessHandshakeNTOR(conn, req)
	}
	return errors.New("unknown handshake type")
}

// ProcessHandshakeTAP handles the "TAP" handshake.
func ProcessHandshakeTAP(conn *Connection, c CreateRequest) error {
	// Reference: https://github.com/torproject/torspec/blob/f9eeae509344dcfd1f185d0130a0055b00131cea/tor-spec.txt#L1027-L1037
	//
	//	   The payload for a CREATE cell is an 'onion skin', which consists of
	//	   the first step of the DH handshake data (also known as g^x).  This
	//	   value is encrypted using the "legacy hybrid encryption" algorithm
	//	   (see 0.4 above) to the server's onion key, giving a client handshake:
	//
	//	       PK-encrypted:
	//	         Padding                       [PK_PAD_LEN bytes]
	//	         Symmetric key                 [KEY_LEN bytes]
	//	         First part of g^x             [PK_ENC_LEN-PK_PAD_LEN-KEY_LEN bytes]
	//	       Symmetrically encrypted:
	//	         Second part of g^x            [DH_LEN-(PK_ENC_LEN-PK_PAD_LEN-KEY_LEN)
	//

	keys := conn.router.config.Keys
	pub, err := torcrypto.HybridDecrypt(keys.Onion, c.HandshakeData)
	if err != nil {
		return errors.Wrap(err, "failed to decrypt TAP handshake data")
	}

	dh, err := torcrypto.GenerateDiffieHellmanKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate DH key")
	}

	s, err := dh.ComputeSharedSecret(pub)
	if err != nil {
		return errors.Wrap(err, "could not compute shared secret")
	}

	// Reference: https://github.com/torproject/torspec/blob/f9eeae509344dcfd1f185d0130a0055b00131cea/tor-spec.txt#L1059-L1060
	//
	//	   Once both parties have g^xy, they derive their shared circuit keys
	//	   and 'derivative key data' value via the KDF-TOR function in 5.2.1.
	//
	k, err := BuildCircuitKeysKDFTOR(s)
	if err != nil {
		return errors.Wrap(err, "failed to build circuit keys")
	}

	err = LaunchCircuit(conn, c.CircID, k)
	if err != nil {
		return errors.Wrap(err, "failed to launch circuit")
	}

	// Reference: https://github.com/torproject/torspec/blob/f9eeae509344dcfd1f185d0130a0055b00131cea/tor-spec.txt#L1040-L1043
	//
	//	   The payload for a CREATED cell, or the relay payload for an
	//	   EXTENDED cell, contains:
	//	         DH data (g^y)                 [DH_LEN bytes]
	//	         Derivative key data (KH)      [HASH_LEN bytes]   <see 5.2 below>
	//

	cell := NewFixedCell(c.CircID, CommandCreated)
	p := cell.Payload()
	copy(p, dh.Public[:])
	copy(p[torcrypto.DiffieHellmanPublicSize:], k.KH)

	err = conn.SendCell(cell)
	if err != nil {
		return errors.Wrap(err, "could not send created cell")
	}

	conn.logger.Info("circuit created")

	return nil
}

func LaunchCircuit(conn *Connection, id CircID, k *CircuitKeys) error {
	fwd := k.ForwardCryptoState()
	back := k.BackwardCryptoState()
	circ := NewTransverseCircuit(conn, id, fwd, back, conn.logger)

	err := conn.circuits.AddWithID(id, circ.ForwardSender())
	if err != nil {
		check.Close(conn.logger, circ)
		return errors.Wrap(err, "failed to register circuit link")
	}

	return nil
}

// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L1075-L1077
//
//	      H_LENGTH  = 32.
//	      ID_LENGTH = 20.
//	      G_LENGTH  = 32
//
// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L1095-L1098
//
//	   and generates a client-side handshake with contents:
//	       NODEID      Server identity digest  [ID_LENGTH bytes]
//	       KEYID       KEYID(B)                [H_LENGTH bytes]
//	       CLIENT_PK   X                       [G_LENGTH bytes]
//
type ClientHandshakeDataNTOR []byte

func (h ClientHandshakeDataNTOR) ServerFingerprint() []byte { return h[:20] }
func (h ClientHandshakeDataNTOR) KeyID() []byte             { return h[20:52] }

func (h ClientHandshakeDataNTOR) ClientPK() [32]byte {
	var X [32]byte
	copy(X[:], h[52:84])
	return X
}

func ProcessHandshakeNTOR(conn *Connection, c CreateRequest) error {
	clientData := ClientHandshakeDataNTOR(c.HandshakeData)

	// Verify the fingerprint matches.
	got := clientData.ServerFingerprint()
	expect := conn.router.Fingerprint()
	ctx := log.WithBytes(conn.logger, "client_handshake_fingerprint", got)
	ctx = log.WithBytes(ctx, "server_fingerprint", expect)
	if !bytes.Equal(got, expect) {
		ctx.Notice("fingerprints do not match")
		return errors.New("incorrect server fingerprint")
	}
	ctx.Debug("verified server fingerprint")

	// Verify the NTOR key ID matches.
	got = clientData.KeyID()
	ntorKey := conn.router.config.Keys.Ntor
	expect = ntorKey.Public[:]
	ctx = conn.logger
	ctx = log.WithBytes(ctx, "client_handshake_keyid", got)
	ctx = log.WithBytes(ctx, "server_keyid", expect)
	if !bytes.Equal(got, expect) {
		ctx.Notice("ntor key ids do not match")
		return errors.New("incorrect ntor key id")
	}
	ctx.Debug("verified ntor key id")

	serverKeyPair, err := torcrypto.GenerateCurve25519KeyPair()
	if err != nil {
		return errors.Wrap(err, "failed to generate server key pair")
	}

	h := ntor.ServerHandshake{
		Public: ntor.Public{
			ID: conn.router.Fingerprint(),
			KX: clientData.ClientPK(),
			KY: serverKeyPair.Public,
			KB: ntorKey.Public,
		},
		Ky: serverKeyPair.Private,
		Kb: ntorKey.Private,
	}

	// Launch the circuit
	keys, err := BuildCircuitKeysNTOR(ntor.KDF(h))
	if err != nil {
		return errors.Wrap(err, "failed to build circuit keys")
	}

	err = LaunchCircuit(conn, c.CircID, keys)
	if err != nil {
		return errors.Wrap(err, "failed to launch circuit")
	}

	// Send reply
	//
	// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L1108-L1110
	//
	//	   The server's handshake reply is:
	//	       SERVER_PK   Y                       [G_LENGTH bytes]
	//	       AUTH        H(auth_input, t_mac)    [H_LENGTH bytes]
	//
	hd := NewServerHandshakeDataNTOR(serverKeyPair.Public, ntor.Auth(h))
	reply := &Created2Cell{
		CircID:        c.CircID,
		HandshakeData: hd,
	}

	err = BuildAndSend(conn, reply)
	if err != nil {
		return errors.Wrap(err, "could not send created2 cell")
	}

	conn.logger.Info("circuit created")

	return nil
}

type CircuitKeys struct {
	KH []byte // Derivative key data
	Df []byte // Forward digest
	Db []byte // Backward digest
	Kf []byte // Forward key
	Kb []byte // Backward key
}

func (k *CircuitKeys) ForwardCryptoState() *CircuitCryptoState {
	return NewCircuitCryptoState(k.Df, k.Kf)
}

func (k *CircuitKeys) BackwardCryptoState() *CircuitCryptoState {
	return NewCircuitCryptoState(k.Db, k.Kb)
}

// circuitKeySize is the number of bytes of key required for circuit keys.
const circuitKeySize = 2*torcrypto.StreamCipherKeySize + 3*torcrypto.HashSize

// BuildCircuitKeysKDFTOR builds circuit keys via the KDF-TOR method from a
// shared secret s.
func BuildCircuitKeysKDFTOR(s []byte) (*CircuitKeys, error) {
	d, err := torcrypto.KDFTOR(s, circuitKeySize)
	if err != nil {
		return nil, errors.Wrap(err, "key derivation error")
	}

	// Reference: https://github.com/torproject/torspec/blob/f9eeae509344dcfd1f185d0130a0055b00131cea/tor-spec.txt#L1181-L1184
	//
	//	   The first HASH_LEN bytes of K form KH; the next HASH_LEN form the forward
	//	   digest Df; the next HASH_LEN 41-60 form the backward digest Db; the next
	//	   KEY_LEN 61-76 form Kf, and the final KEY_LEN form Kb.  Excess bytes from K
	//	   are discarded.
	//
	k := &CircuitKeys{}
	k.KH, d = buf.Consume(d, torcrypto.HashSize)
	k.Df, d = buf.Consume(d, torcrypto.HashSize)
	k.Db, d = buf.Consume(d, torcrypto.HashSize)
	k.Kf, d = buf.Consume(d, torcrypto.StreamCipherKeySize)
	k.Kb, _ = buf.Consume(d, torcrypto.StreamCipherKeySize)

	return k, nil
}

// BuildCircuitKeysNTOR generates Circuit key material from r.
func BuildCircuitKeysNTOR(r io.Reader) (*CircuitKeys, error) {
	// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L1210-L1214
	//
	//	   When used in the ntor handshake, the first HASH_LEN bytes form the
	//	   forward digest Df; the next HASH_LEN form the backward digest Db; the
	//	   next KEY_LEN form Kf, the next KEY_LEN form Kb, and the final
	//	   DIGEST_LEN bytes are taken as a nonce to use in the place of KH in the
	//	   hidden service protocol.  Excess bytes from K are discarded.
	//
	d := make([]byte, circuitKeySize)
	_, err := io.ReadFull(r, d[:])
	if err != nil {
		return nil, errors.Wrap(err, "short read for circuit key material")
	}

	k := &CircuitKeys{}
	k.Df, d = buf.Consume(d, torcrypto.HashSize)
	k.Db, d = buf.Consume(d, torcrypto.HashSize)
	k.Kf, d = buf.Consume(d, torcrypto.StreamCipherKeySize)
	k.Kb, d = buf.Consume(d, torcrypto.StreamCipherKeySize)
	k.KH, _ = buf.Consume(d, torcrypto.HashSize)

	return k, nil
}

// ServerHandshakeDataNTOR represents server handshake data for the NTOR handshake.
//
// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L1108-L1110
//
//	   The server's handshake reply is:
//	       SERVER_PK   Y                       [G_LENGTH bytes]
//	       AUTH        H(auth_input, t_mac)    [H_LENGTH bytes]
//
type ServerHandshakeDataNTOR []byte

func NewServerHandshakeDataNTOR(Y [32]byte, auth []byte) ServerHandshakeDataNTOR {
	var b []byte
	b = append(b, Y[:]...)
	b = append(b, auth...)
	return b
}
