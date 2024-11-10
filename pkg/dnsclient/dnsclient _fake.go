package dnsclient

import (
	"context"
)

// Record represents a DNS record.
type Record struct {
	ZoneID     string
	Name       string
	RecordType string
	Rrdata     string
	Opts       DNSRecordOptions
}

// FakeDNSClient is a fake implementation of the DNSClient interface.
type FakeDNSClient struct {
	Zones   map[string]string
	Records map[string]Record
}

// NewFakeDNSClient creates a new FakeDNSClient.
func NewFakeDNSClient(
	zones map[string]string,
	records map[string]Record,
) *FakeDNSClient {
	return &FakeDNSClient{
		Zones:   zones,
		Records: records,
	}
}

// GetManagedZones returns a map of all managed zone DNS names mapped to their IDs, composed of the project ID and
// their user assigned resource names.
func (c *FakeDNSClient) GetManagedZones(_ context.Context) (map[string]string, error) {
	zones := make(map[string]string)

	for k, v := range c.Zones {
		zones[k] = v
	}

	return zones, nil
}

func (c *FakeDNSClient) getRecords() map[string]Record {
	r := map[string]Record{}
	for k, v := range c.Records {
		r[k] = v
	}
	return r
}

// CreateOrUpdateRecordSet creates or updates the resource recordset with the given name, record type, rrdatas, and ttl
// in the managed zone with the given name or ID.
func (c *FakeDNSClient) CreateOrUpdateRecordSet(
	_ context.Context,
	zoneID,
	name,
	recordType string,
	rrdatas []string,
	opts DNSRecordOptions,
) error {
	records := c.getRecords()
	for _, rrdata := range rrdatas {
		c.Records[rrdata] = Record{
			ZoneID:     zoneID,
			Name:       name,
			RecordType: recordType,
			Rrdata:     rrdata,
			Opts:       opts,
		}
		delete(records, rrdata)
	}

	// delete undefined data
	for rrdata := range records {
		delete(c.Records, rrdata)
	}
	return nil
}

// DeleteRecordSet deletes the resource recordset with the given name and record type
// in the managed zone with the given name or ID.
func (c *FakeDNSClient) DeleteRecordSet(_ context.Context, zoneID, name, recordType string) error {
	for rrdata, record := range c.Records {
		if record.ZoneID == zoneID && record.RecordType == recordType && record.Name == name {
			delete(c.Records, rrdata)
		}
	}
	return nil
}
