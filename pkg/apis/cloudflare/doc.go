// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +k8s:deepcopy-gen=package
// +k8s:defaulter-gen=TypeMeta

// Package v1alpha1 contains the cloudflare provider configuration API resources.
// +groupName=cloudflare.dns.provider.extensions.gardener.cloud

//go:generate ../../../hack/update-codegen.sh

package cloudflare // import "github.com/schrodit/gardener-extension-provider-dns-cloudflare/pkg/apis/cloudflare"
