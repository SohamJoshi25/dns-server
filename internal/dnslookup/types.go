package dnslookup

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type DNSQuestion struct {
	Name  string
	Type  uint16
	Class uint16
}

type DNSAnswer struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    string
}

type Record struct {
	Name   string
	Type   uint16
	Answer string
}
