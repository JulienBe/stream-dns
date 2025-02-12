package main

import (
	"fmt"

	dns "github.com/miekg/dns"
)

// Structure used to convert record
// from kafka in a JSON format into RR
type Record struct {
	Name     string
	Type     string
	Content  string
	Ttl      int
	Priority int
}

func A(rr string) *dns.A { r, _ := dns.NewRR(rr); return r.(*dns.A) }

func AAAA(rr string) *dns.AAAA { r, _ := dns.NewRR(rr); return r.(*dns.AAAA) }

func CNAME(rr string) *dns.CNAME { r, _ := dns.NewRR(rr); return r.(*dns.CNAME) }

func SOA(rr string) *dns.SOA { r, _ := dns.NewRR(rr); return r.(*dns.SOA) }

func MX(rr string) *dns.MX { r, _ := dns.NewRR(rr); return r.(*dns.MX) }

func NS(rr string) *dns.NS { r, _ := dns.NewRR(rr); return r.(*dns.NS) }

func TXT(rr string) *dns.TXT { r, _ := dns.NewRR(rr); return r.(*dns.TXT) }

func PTR(rr string) *dns.PTR { r, _ := dns.NewRR(rr); return r.(*dns.PTR) }

func recordToString(record Record) string {
	if record.Priority > 0 {
		return fmt.Sprintf("%s %d IN %s %d %s", record.Name, record.Ttl, record.Type, record.Priority, record.Content)
	} else {
		return fmt.Sprintf("%s %d IN %s %s", record.Name, record.Ttl, record.Type, record.Content)
	}
}

// - The primary name server for the domain
// - The responsible party for the domain
// - A timestamp that changes whenever you update your domain.
// - The number of seconds before the zone should be refreshed.
// - The number of seconds before a failed refresh should be retried.
// - The upper limit in seconds before a zone is considered no longer authoritative.
// - The negative result TTL
func recordSOAToString(record Record) string {
	return fmt.Sprintf("%s %d IN SOA %s", record.Name, record.Ttl, record.Content)
}

// RecordToRR converts a Record to a dns.RR.
// If the RR is not a type that this package uses,
// it returns a nil RR.
func RecordToRR(record Record) dns.RR {
	var rr dns.RR
	rtype := dns.StringToType[record.Type]
	recordstr := recordToString(record)

	switch rtype {
	case dns.TypeA:
		rr = A(recordstr)
	case dns.TypeAAAA:
		rr = AAAA(recordstr)
	case dns.TypeCNAME:
		rr = CNAME(recordstr)
	case dns.TypeSOA:
		rr = SOA(recordSOAToString(record))
	case dns.TypeMX:
		rr = MX(recordstr)
	case dns.TypeNS:
		rr = NS(recordstr)
	case dns.TypeTXT:
		rr = TXT(recordstr)
	case dns.TypePTR:
		rr = PTR(recordstr)
	default:
		return nil
	}

	return rr
}

// Convert slice of Record into a  slice of dns.RR
func MapRecordsIntoRRs(records []Record) (rrs []dns.RR, err error) {
	for _, r := range records {
		rrs = append(rrs, RecordToRR(r))
	}

	return
}
