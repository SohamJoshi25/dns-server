package dnsproxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/golang-lru/v2/expirable"
	dnsdb "github.com/sohamjoshi25/dns-server/internal/dnsdb"
	dnslookup "github.com/sohamjoshi25/dns-server/internal/dnslookup"
)

func HandleDNSRequest(conn *net.UDPConn, cache *expirable.LRU[dnslookup.DNSQuestion, []dnslookup.DNSAnswer]) {
	buffer := make([]byte, 512)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Could Not Read From UDP Datagram")
		return
	}

	header, question, err := parseDNSQuery(buffer[:n])
	if err != nil {
		fmt.Println("Failed to parse DNS query")
		return
	}

	if question.Type == 12 {
		return
	}

	answers, err := dnsdb.QueryDatabase(question.Name, question.Type)

	if err == nil {
		fmt.Println("Domain found in Database", question.Name, question.Class)
		response := buildDNSResponse(header, question, answers, true)
		conn.WriteToUDP(response, addr)
	} else {

		DNSanswers, ok := cache.Get(*question)

		if ok { //Check Cache

			answers := make([]string, len(DNSanswers))
			for i, DNSanswer := range DNSanswers {
				answers[i] = DNSanswer.RData
			}
			fmt.Println("Response from Cache", question.Name, question.Class)
			response := buildDNSResponse(header, question, answers, false)
			conn.WriteToUDP(response, addr)
			return

		} else {

			fmt.Println("Domain NOT found in database : Performing Iterative Lookup", question.Name, question.Class)

			answer, err := dnslookup.IterativeLookup(*question, *header, cache, false)
			if err != nil {
				fmt.Println("Failed to LookUp", err.Error())
				return
			}

			conn.WriteToUDP(answer, addr)
			return
		}

	}

}

func parseDNSQuery(data []byte) (*dnslookup.DNSHeader, *dnslookup.DNSQuestion, error) {
	header := &dnslookup.DNSHeader{}
	bufReader := bytes.NewReader(data)

	err := binary.Read(bufReader, binary.BigEndian, header)
	if err != nil {
		return nil, nil, err
	}

	qName, err := readDomainName(bufReader)
	if err != nil {
		return nil, nil, err
	}

	var qType, qClass uint16
	binary.Read(bufReader, binary.BigEndian, &qType)
	binary.Read(bufReader, binary.BigEndian, &qClass)

	question := &dnslookup.DNSQuestion{Name: qName, Type: qType, Class: qClass}
	return header, question, nil
}

func buildDNSResponse(header *dnslookup.DNSHeader, question *dnslookup.DNSQuestion, answers []string, authoritative bool) []byte {
	response := new(bytes.Buffer)
	ANCount := uint16(len(answers))

	if authoritative {
		header.Flags |= 0x8400 // Authoritative + Standard Response
	} else {
		header.Flags |= 0x8000 // Standard Response
	}
	header.Flags |= 0x0100 // Recursion Available ✅

	header.ANCount = ANCount

	// Manually write the header fields in proper order
	binary.Write(response, binary.BigEndian, header.ID)
	binary.Write(response, binary.BigEndian, header.Flags)
	binary.Write(response, binary.BigEndian, header.QDCount)
	binary.Write(response, binary.BigEndian, header.ANCount)
	binary.Write(response, binary.BigEndian, header.NSCount)
	binary.Write(response, binary.BigEndian, header.ARCount)

	// Write the original question back
	writeDomainName(response, question.Name)
	binary.Write(response, binary.BigEndian, question.Type)
	binary.Write(response, binary.BigEndian, question.Class)

	// Answer Section (Multiple Answers)
	for _, answer := range answers {
		// Compression pointer (offset to start of domain name in query)
		response.WriteByte(0xC0)
		response.WriteByte(0x0C)

		binary.Write(response, binary.BigEndian, question.Type) // Type
		binary.Write(response, binary.BigEndian, uint16(1))     // Class (IN)
		binary.Write(response, binary.BigEndian, uint32(60))    // TTL

		var dataBuffer bytes.Buffer
		var dataLength uint16

		if question.Type == 1 && net.ParseIP(answer).To4() != nil {
			ipBytes := net.ParseIP(answer).To4()
			dataLength = 4
			dataBuffer.Write(ipBytes)
		} else if question.Type == 28 && net.ParseIP(answer).To16() != nil {
			ipBytes := net.ParseIP(answer).To16()
			dataLength = 16
			dataBuffer.Write(ipBytes)
		} else if question.Type == 16 || question.Type == 99 {
			textBytes := []byte(answer)
			dataBuffer.WriteByte(byte(len(textBytes))) // Write Length byte
			dataBuffer.Write(textBytes)                // Write Text segment
			dataLength = uint16(dataBuffer.Len())
		} else {
			dataBuffer.Write([]byte(answer))
			dataLength = uint16(dataBuffer.Len())
		}

		binary.Write(response, binary.BigEndian, dataLength)
		response.Write(dataBuffer.Bytes())
	}

	return response.Bytes()
}

func readDomainName(buf *bytes.Reader) (string, error) {
	var labels []string
	for {
		length, _ := buf.ReadByte()
		if length == 0 {
			break
		}
		label := make([]byte, length)
		buf.Read(label)
		labels = append(labels, string(label))
	}
	return strings.Join(labels, "."), nil
}

func writeDomainName(buf *bytes.Buffer, domain string) {
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		buf.WriteByte(byte(len(part))) // Write label length
		buf.WriteString(part)          // Write label characters
	}
	buf.WriteByte(0) // Null terminator
}
