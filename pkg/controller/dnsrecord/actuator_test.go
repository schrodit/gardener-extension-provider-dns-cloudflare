package dnsrecord_test

import (
	"bytes"
	"context"
	"fmt"
	"io"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	cfinstall "github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/apis/cloudflare/install"
	cfv1alpha1 "github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/apis/cloudflare/v1alpha1"
	"github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/controller/dnsrecord"
	"github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/dnsclient"
	"github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/utils/rand"
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
		Expect(extensionsv1alpha1.AddToScheme(scheme)).ToNot(HaveOccurred())
		Expect(gardencorev1beta1.AddToScheme(scheme)).ToNot(HaveOccurred())
		Expect(cfinstall.AddToScheme(scheme)).ToNot(HaveOccurred())

		codec := serializer.NewCodecFactory(scheme, serializer.EnableStrict)

		info, found := runtime.SerializerInfoForMediaType(codec.SupportedMediaTypes(), runtime.ContentTypeJSON)
		Expect(found).To(BeTrue(), "should be able to decode")

		encoder = codec.EncoderForVersion(info.Serializer, cfv1alpha1.SchemeGroupVersion)
	})

	It("no providerconfig -> default proxied to false", func(ctx context.Context) {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("test-%s", rand.RandStringRunes(5)),
			},
		}
		Expect(k8sManager.GetClient().Create(ctx, ns)).To(Succeed())
		record := &extensionsv1alpha1.DNSRecord{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: ns.Name,
			},
			Spec: extensionsv1alpha1.DNSRecordSpec{
				Name:       "test.example.com",
				RecordType: "A",
				Values:     []string{"127.0.0.1"},
				TTL:        Int(300),
			},
		}
		c := k8sManager.GetClient()
		Expect(c.Create(ctx, record)).To(Succeed())

		a := dnsrecord.NewActuator(k8sManager)

		dnsClient := dnsclient.NewFakeDNSClient(
			map[string]string{"example.com": "1"},
			make(map[string]dnsclient.Record),
		)
		Expect(a.ReconcileRecord(ctx, logr.Discard(), record.DeepCopy(), &extensionscontroller.Cluster{}, dnsClient)).NotTo(HaveOccurred())
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
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("test-%s", rand.RandStringRunes(5)),
			},
		}
		Expect(k8sManager.GetClient().Create(ctx, ns)).To(Succeed())
		record := &extensionsv1alpha1.DNSRecord{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: ns.Name,
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

		c := k8sManager.GetClient()
		Expect(c.Create(ctx, record)).To(Succeed())
		a := dnsrecord.NewActuator(k8sManager)

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
