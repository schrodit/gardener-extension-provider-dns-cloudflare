<p>Packages:</p>
<ul>
<li>
<a href="#cloudflare.dns.provider.extensions.gardener.cloud%2fv1alpha1">cloudflare.dns.provider.extensions.gardener.cloud/v1alpha1</a>
</li>
</ul>
<h2 id="cloudflare.dns.provider.extensions.gardener.cloud/v1alpha1">cloudflare.dns.provider.extensions.gardener.cloud/v1alpha1</h2>
<p>
<p>Package v1alpha1 contains the Cloudflare provider configuration API resources.</p>
</p>
Resource Types:
<ul><li>
<a href="#cloudflare.dns.provider.extensions.gardener.cloud/v1alpha1.DnsRecordConfig">DnsRecordConfig</a>
</li></ul>
<h3 id="cloudflare.dns.provider.extensions.gardener.cloud/v1alpha1.DnsRecordConfig">DnsRecordConfig
</h3>
<p>
<p>DnsRecordConfig defines the configuration for the Cloudflare dns provider.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
cloudflare.dns.provider.extensions.gardener.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>DnsRecordConfig</code></td>
</tr>
<tr>
<td>
<code>proxied</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Proxied specifies whether the traffic for that domain should be proxied by Cloudflare.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <a href="https://github.com/ahmetb/gen-crd-api-reference-docs">gen-crd-api-reference-docs</a>
</em></p>
