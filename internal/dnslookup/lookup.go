package dnslookup

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

func IterativeLookup(_question DNSQuestion, _header DNSHeader, cache *expirable.LRU[DNSQuestion, []DNSAnswer], debug bool) ([]byte, error) {
	server := RootServers[0]
	var itterations uint8 = 1

	for itterations < 16 {
		_, _, _, answers, authorities, _, err := sendDNSQuery(_header, _question.Name, 2, server)
		if debug {
			fmt.Printf("itterations: %+v Server: %+v\n", itterations, server)
		}

		if err != nil {
			log.Fatalln(err.Error())
			return nil, err
		}

		if len(answers) != 0 {

			_response, _header, _question, __answers, _authorities, _additionals, err := sendDNSQuery(_header, _question.Name, _question.Type, server)
			if err != nil {
				log.Fatalln(err.Error())
				return nil, err
			}
			if debug {
				fmt.Printf("\n\n\nQuestion: %+v\n", _question)
				fmt.Printf("itterations: %+v\n", itterations)
				fmt.Printf("Response Bytes: %+v\n", _response)
				fmt.Printf("Header: %+v\n", _header)
				fmt.Printf("Answers: %+v\n", __answers)
				fmt.Printf("Authorities: %+v\n", _authorities)
				fmt.Printf("Additionals: %+v\n", _additionals)
			}

			cache.Add(_question, __answers)
			return _response, nil
		} else if len(authorities) != 0 {
			server = authorities[0].RData
		} else {
			return nil, fmt.Errorf("Some Error Occured. Authoritative and answers are both empty Itterations: %+v Server: %+v\n", itterations, server)
		}

		itterations++
	}

	return nil, fmt.Errorf("itterations Excedded 16 Limit")
}

func sendDNSQuery(_header DNSHeader, domain string, qtype uint16, server string) ([]byte, DNSHeader, DNSQuestion, []DNSAnswer, []DNSAnswer, []DNSAnswer, error) {

	if !strings.HasSuffix(server, ":53") {
		server = server + ":53"
	}

	conn, err := net.Dial("udp", server)
	if err != nil {
		return nil, DNSHeader{}, DNSQuestion{}, nil, nil, nil, fmt.Errorf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Build DNS Query Packet
	query := new(bytes.Buffer)

	// Header Section (12 bytes)
	header := DNSHeader{
		ID:      _header.ID,
		Flags:   0x0100,
		QDCount: 1,
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	binary.Write(query, binary.BigEndian, header.ID)
	binary.Write(query, binary.BigEndian, header.Flags)
	binary.Write(query, binary.BigEndian, header.QDCount)
	binary.Write(query, binary.BigEndian, header.ANCount)
	binary.Write(query, binary.BigEndian, header.NSCount)
	binary.Write(query, binary.BigEndian, header.ARCount)

	// Question Section
	domainParts := encodeDomain(domain)
	query.Write(domainParts)
	binary.Write(query, binary.BigEndian, qtype)     // QTYPE
	binary.Write(query, binary.BigEndian, uint16(1)) // QCLASS (IN)

	question := DNSQuestion{
		Name:  domain,
		Type:  qtype,
		Class: 1,
	}

	// Send DNS Query
	_, err = conn.Write(query.Bytes())
	if err != nil {
		return nil, DNSHeader{}, DNSQuestion{}, nil, nil, nil, fmt.Errorf("Failed to send query: %v", err)
	}

	// Receive DNS Response
	buffer := make([]byte, 512)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, DNSHeader{}, DNSQuestion{}, nil, nil, nil, fmt.Errorf("Failed to read response: %v", err)
	}
	buffer = buffer[:n]

	rawHeader := buffer[:12]
	header.ID = binary.BigEndian.Uint16(rawHeader[0:2])
	header.Flags = binary.BigEndian.Uint16(rawHeader[2:4])
	header.QDCount = binary.BigEndian.Uint16(rawHeader[4:6])
	header.ANCount = binary.BigEndian.Uint16(rawHeader[6:8])
	header.NSCount = binary.BigEndian.Uint16(rawHeader[8:10])
	header.ARCount = binary.BigEndian.Uint16(rawHeader[10:12])

	questionEnd := 12 + len(domainParts) + 4
	offset := questionEnd

	answers := parseAnswers(buffer, &offset, int(header.ANCount))
	authorities := parseAnswers(buffer, &offset, int(header.NSCount))
	additionals := parseAnswers(buffer, &offset, int(header.ARCount))

	return buffer, header, question, answers, authorities, additionals, nil
}

func encodeDomain(domain string) []byte {
	parts := bytes.Split([]byte(domain), []byte("."))
	buffer := new(bytes.Buffer)

	for _, part := range parts {
		buffer.WriteByte(byte(len(part)))
		buffer.Write(part)
	}
	buffer.WriteByte(0) // Null terminator

	return buffer.Bytes()
}

func parseAnswers(response []byte, offset *int, count int) []DNSAnswer {
	var answers []DNSAnswer

	for i := 0; i < count; i++ {
		if *offset >= len(response) {
			break
		}

		name, nameEnd := decodeDomainName(response, *offset)
		*offset = nameEnd

		if *offset+10 > len(response) {
			break
		}

		answer := DNSAnswer{
			Name:     name,
			Type:     binary.BigEndian.Uint16(response[*offset : *offset+2]),
			Class:    binary.BigEndian.Uint16(response[*offset+2 : *offset+4]),
			TTL:      binary.BigEndian.Uint32(response[*offset+4 : *offset+8]),
			RDLength: binary.BigEndian.Uint16(response[*offset+8 : *offset+10]),
		}
		*offset += 10

		if *offset+int(answer.RDLength) > len(response) {
			break
		}

		if answer.Type == 2 || answer.Type == 5 || answer.Type == 12 {
			answer.RData, _ = decodeDomainName(response, *offset)
		} else {
			answer.RData = net.IP(response[*offset : *offset+int(answer.RDLength)]).String()
		}
		*offset += int(answer.RDLength)

		answers = append(answers, answer)
	}

	return answers
}

func decodeDomainName(response []byte, offset int) (string, int) {
	var name []string
	originalOffset := offset
	jumped := false

	for {
		if offset >= len(response) {
			break
		}

		length := int(response[offset])

		if length == 0 {
			offset++
			break
		}

		if length&0xC0 == 0xC0 {
			pointer := int(binary.BigEndian.Uint16(response[offset:offset+2]) & 0x3FFF)
			offset += 2
			if !jumped {
				originalOffset = offset
			}
			jumped = true
			offset = pointer
			continue
		}

		offset++
		if offset+length > len(response) {
			break
		}
		name = append(name, string(response[offset:offset+length]))
		offset += length
	}

	if jumped {
		return strings.Join(name, "."), originalOffset
	}

	return strings.Join(name, "."), offset
}
