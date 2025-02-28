package dnslookup

var RootServers []string = []string{
	"198.41.0.4",
	"199.9.14.201",
	"192.33.4.12",
	"199.7.91.13",
}

var RRTypeMap = map[uint16]string{
	1:   "A",
	2:   "NS",
	5:   "CNAME",
	6:   "SOA",
	12:  "PTR",
	15:  "MX",
	16:  "TXT",
	28:  "AAAA",
	33:  "SRV",
	39:  "DNAME",
	41:  "OPT",
	43:  "DS",
	46:  "RRSIG",
	47:  "NSEC",
	48:  "DNSKEY",
	50:  "NSEC3",
	51:  "NSEC3PARAM",
	52:  "TLSA",
	99:  "SPF",
	257: "CAA",
}

var RRClassMap = map[uint16]string{
	1: "IN",
	2: "CS",
	3: "CH",
	4: "HS",
}
