//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	cloudflare "github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/apis/cloudflare"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*DnsRecordConfig)(nil), (*cloudflare.DnsRecordConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_DnsRecordConfig_To_cloudflare_DnsRecordConfig(a.(*DnsRecordConfig), b.(*cloudflare.DnsRecordConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cloudflare.DnsRecordConfig)(nil), (*DnsRecordConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cloudflare_DnsRecordConfig_To_v1alpha1_DnsRecordConfig(a.(*cloudflare.DnsRecordConfig), b.(*DnsRecordConfig), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_DnsRecordConfig_To_cloudflare_DnsRecordConfig(in *DnsRecordConfig, out *cloudflare.DnsRecordConfig, s conversion.Scope) error {
	out.Proxied = in.Proxied
	return nil
}

// Convert_v1alpha1_DnsRecordConfig_To_cloudflare_DnsRecordConfig is an autogenerated conversion function.
func Convert_v1alpha1_DnsRecordConfig_To_cloudflare_DnsRecordConfig(in *DnsRecordConfig, out *cloudflare.DnsRecordConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_DnsRecordConfig_To_cloudflare_DnsRecordConfig(in, out, s)
}

func autoConvert_cloudflare_DnsRecordConfig_To_v1alpha1_DnsRecordConfig(in *cloudflare.DnsRecordConfig, out *DnsRecordConfig, s conversion.Scope) error {
	out.Proxied = in.Proxied
	return nil
}

// Convert_cloudflare_DnsRecordConfig_To_v1alpha1_DnsRecordConfig is an autogenerated conversion function.
func Convert_cloudflare_DnsRecordConfig_To_v1alpha1_DnsRecordConfig(in *cloudflare.DnsRecordConfig, out *DnsRecordConfig, s conversion.Scope) error {
	return autoConvert_cloudflare_DnsRecordConfig_To_v1alpha1_DnsRecordConfig(in, out, s)
}
