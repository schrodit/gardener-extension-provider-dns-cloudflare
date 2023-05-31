package dnsrecord_test

import (
	"bytes"
	"context"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cfinstall "github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/apis/cloudflare/install"
	cfv1alpha1 "github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/apis/cloudflare/v1alpha1"
	"github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/controller/dnsrecord"
	"github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/dnsclient"
)

func Int(i int64) *int64 {
	return &i
}

var _ = Describe("Actuator", func() {

	var (
		scheme  = runtime.NewScheme()
		encoder runtime.Encoder
		encode  = func(obj runtime.Object) []byte {
			b := &bytes.Buffer{}
			Expect(encoder.Encode(obj, b)).To(Succeed())

			data, err := io.ReadAll(b)
			Expect(err).ToNot(HaveOccurred())

			return data
		}
	)

	BeforeEach(func() {
		extensionsv1alpha1.AddToScheme(scheme)
		gardencorev1beta1.AddToScheme(scheme)
		cfinstall.AddToScheme(scheme)

		codec := serializer.NewCodecFactory(scheme, serializer.EnableStrict)

		info, found := runtime.SerializerInfoForMediaType(codec.SupportedMediaTypes(), runtime.ContentTypeJSON)
		Expect(found).To(BeTrue(), "should be able to decode")

		kubernetes.GardenCodec.UniversalDecoder()

		encoder = codec.EncoderForVersion(info.Serializer, cfv1alpha1.SchemeGroupVersion)
	})

	It("no providerconfig -> default proxied to false", func(ctx context.Context) {
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

	It("providerconfig proxied true -> proxied is enabled", func(ctx context.Context) {
		record := &extensionsv1alpha1.DNSRecord{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: extensionsv1alpha1.DNSRecordSpec{
				Name:       "test.example.com",
				RecordType: "A",
				Values:     []string{"127.0.0.1"},
				TTL:        Int(300),
				DefaultSpec: extensionsv1alpha1.DefaultSpec{
					ProviderConfig: &runtime.RawExtension{Raw: encode(&cfv1alpha1.DnsRecordConfig{
						Proxied: true,
					})},
				},
			},
		}

		a := dnsrecord.NewActuator()
		inject.ClientInto(fake.NewClientBuilder().WithScheme(scheme).WithObjects(record).Build(), a)
		inject.SchemeInto(scheme, a)

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
				Proxied: true,
			},
		}))
	})

})
