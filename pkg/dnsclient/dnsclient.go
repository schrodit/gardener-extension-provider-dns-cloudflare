package dnsclient

import (
	"context"
	"fmt"

	cloudflare "github.com/cloudflare/cloudflare-go"
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// API_TOKEN_FIELD is the field name in the secret data containing the api token.
const API_TOKEN_FIELD = "apiToken"

// DNSRecordOptions contains optional parameters for creating or updating a DNS record.
type DNSRecordOptions struct {
	TTL     int64
	Proxied bool
}

// DNSClient is an interface for interacting with a DNS provider.
type DNSClient interface {
	GetManagedZones(ctx context.Context) (map[string]string, error)
	CreateOrUpdateRecordSet(
		ctx context.Context, managedZone,
		name,
		recordType string,
		rrdatas []string,
		opts DNSRecordOptions) error
	DeleteRecordSet(ctx context.Context, managedZone, name, recordType string) error
}

type dnsClient struct {
	api *cloudflare.API
}

// NewDNSClient creates a new dns client with a cloudflare api token.
func NewDNSClient(_ context.Context, apiToken string) (DNSClient, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, err
	}

	return &dnsClient{
		api: api,
	}, nil
}

// NewDNSClientFromSecretRef creates a bew dns client from a secret containing an apiToken.
func NewDNSClientFromSecretRef(ctx context.Context, c client.Client, secretRef corev1.SecretReference) (DNSClient, error) {
	secret, err := extensionscontroller.GetSecretByReference(ctx, c, &secretRef)
	if err != nil {
		return nil, err
	}

	apiToken, ok := secret.Data[API_TOKEN_FIELD]
	if !ok {
		return nil, fmt.Errorf("no api token defined")
	}
	return NewDNSClient(ctx, string(apiToken))
}

// GetManagedZones returns a map of all managed zone DNS names mapped to their IDs, composed of the project ID and
// their user assigned resource names.
func (c *dnsClient) GetManagedZones(ctx context.Context) (map[string]string, error) {
	zones := make(map[string]string)

	result, err := c.api.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	for _, z := range result {
		zones[z.Name] = z.ID
	}

	return zones, nil
}

// CreateOrUpdateRecordSet creates or updates the resource recordset with the given name, record type, rrdatas, and ttl
// in the managed zone with the given name or ID.
func (c *dnsClient) CreateOrUpdateRecordSet(
	ctx context.Context,
	zoneID,
	name,
	recordType string,
	rrdatas []string,
	opts DNSRecordOptions,
) error {
	records, err := c.getRecordSet(ctx, name, zoneID)
	if err != nil {
		return err
	}
	for _, rrdata := range rrdatas {
		if record, ok := records[rrdata]; ok {
			if recordIsOutdated(record, opts) {
				if err := c.updateRecord(
					ctx, zoneID, record.ID, name, recordType, rrdata, opts,
				); err != nil {
					return err
				}
			}

			// entry already exists
			delete(records, rrdata)
			continue
		}
		if err := c.createRecord(ctx, zoneID, name, recordType, rrdata, opts); err != nil {
			return err
		}
		delete(records, rrdata)
	}

	// delete undefined data
	for _, record := range records {
		if err := c.deleteRecord(ctx, zoneID, record.ID, name, record.Content); err != nil {
			return err
		}
	}
	return nil
}

// DeleteRecordSet deletes the resource recordset with the given name and record type
// in the managed zone with the given name or ID.
func (c *dnsClient) DeleteRecordSet(ctx context.Context, zoneID, name, recordType string) error {
	records, err := c.getRecordSet(ctx, name, zoneID)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Type != recordType {
			continue
		}
		if err := c.deleteRecord(ctx, zoneID, record.ID, name, record.Content); err != nil {
			return err
		}
	}

	return nil
}

func (c *dnsClient) createRecord(
	ctx context.Context,
	zoneID,
	name,
	recordType,
	rrdata string,
	opts DNSRecordOptions,
) error {
	res, err := c.api.CreateDNSRecord(ctx, zoneID, cloudflare.DNSRecord{
		Name:    name,
		Type:    recordType,
		TTL:     int(opts.TTL),
		Content: rrdata,
		Proxied: &opts.Proxied,
	})
	if err != nil {
		return fmt.Errorf("Unable to set dns record for %s to %s: %w", name, rrdata, err)
	}
	if !res.Success {
		return fmt.Errorf("Unable to set dns record for %s to %s: %#v", name, rrdata, res.Errors)
	}
	return nil
}

func (c *dnsClient) updateRecord(
	ctx context.Context,
	zoneID,
	recordID,
	name,
	recordType,
	rrdata string,
	opts DNSRecordOptions,
) error {
	return c.api.UpdateDNSRecord(ctx, zoneID, recordID, cloudflare.DNSRecord{
		Name:    name,
		Type:    recordType,
		TTL:     int(opts.TTL),
		Content: rrdata,
		Proxied: &opts.Proxied,
	})
}

func (c *dnsClient) deleteRecord(ctx context.Context, zoneID, recordID, name, rrdata string) error {
	err := c.api.DeleteDNSRecord(ctx, zoneID, recordID)
	if err != nil {
		return fmt.Errorf("Unable to set dns record for %s to %s: %w", name, rrdata, err)
	}
	return nil
}

func (c *dnsClient) getZoneID(ctx context.Context, name string) (string, error) {
	zones, err := c.GetManagedZones(ctx)
	if err != nil {
		return "", err
	}
	zoneID, ok := zones[name]
	if !ok {
		return "", fmt.Errorf("No zone found for %s", name)
	}
	return zoneID, nil
}

// getRecordSets returns a map of rrdata to dns record for the given name.
func (c *dnsClient) getRecordSet(ctx context.Context, name, zoneID string) (map[string]cloudflare.DNSRecord, error) {
	results, err := c.api.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	records := make(map[string]cloudflare.DNSRecord, len(results))
	for _, record := range results {
		records[record.Content] = record
	}
	return records, nil
}

func recordIsOutdated(current cloudflare.DNSRecord, expected DNSRecordOptions) bool {
	return *current.Proxied != expected.Proxied || current.TTL != int(expected.TTL)
}
