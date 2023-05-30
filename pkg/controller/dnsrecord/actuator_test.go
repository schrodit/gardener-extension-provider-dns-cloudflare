package dnsrecord_test

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/controller/dnsrecord"
	"github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/dnsclient"
)

func Int(i int64) *int64 {
	return &i
}

var _ = Describe("Actuator", func() {

	It("no providerconfig -> default dns to false", func(ctx context.Context) {
		scheme := runtime.NewScheme()
		extensionsv1alpha1.AddToScheme(scheme)
		gardencorev1beta1.AddToScheme(scheme)

		record := &extensionsv1alpha1.DNSRecord{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: extensionsv1alpha1.DNSRecordSpec{
				Name:       "test.example.com",
				RecordType: "A",
				Values:     []string{"127.0.0.1"},
				TTL:        Int(300),
			},
		}

		a := dnsrecord.NewActuator()
		inject.ClientInto(fake.NewClientBuilder().WithScheme(scheme).WithObjects(record).Build(), a)

		dnsClient := dnsclient.NewFakeDNSClient(
			map[string]string{"example.com": "1"},
			make(map[string]dnsclient.Record),
		)
		Expect(a.ReconcileRecord(ctx, logr.Discard(), record, &extensionscontroller.Cluster{}, dnsClient)).NotTo(HaveOccurred())
		Expect(dnsClient.Records).To(HaveKeyWithValue("127.0.0.1", dnsclient.Record{
			ZoneID:     "1",
			Name:       "test.example.com",
			RecordType: "A",
			Rrdata:     "127.0.0.1",
			Opts: dnsclient.DNSRecordOptions{
				TTL:     300,
				Proxied: false,
			},
		}))
	})

})
