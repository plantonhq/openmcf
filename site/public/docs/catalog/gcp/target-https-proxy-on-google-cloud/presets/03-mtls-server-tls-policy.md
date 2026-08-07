---
title: "Mutual-TLS Frontend"
description: "An HTTPS frontend that authenticates CLIENTS, not just the server: a network security `ServerTlsPolicy` demands and validates client certificates during the handshake, on top of the normal server..."
type: "preset"
rank: "03"
presetSlug: "03-mtls-server-tls-policy"
componentSlug: "target-https-proxy-on-google-cloud"
componentTitle: "Target HTTPS Proxy on Google Cloud"
provider: "gcp"
icon: "package"
order: 3
---

# Mutual-TLS Frontend

An HTTPS frontend that authenticates CLIENTS, not just the server: a network security `ServerTlsPolicy` demands and validates client certificates during the handshake, on top of the normal server certificate. TLS early data stays disabled — replayable 0-RTT requests and client-certificate auth do not mix.

## When to Use

- B2B APIs where callers present organization-issued client certificates
- Zero-trust perimeters that require mutual authentication before any request reaches a backend

## Remix Notes

- The `ServerTlsPolicy` (with its trust config of CA anchors) is created in Network Security / Certificate Manager; this proxy attaches it by resource name.
- On Traffic Director (`proxyBind: true` + `INTERNAL_SELF_MANAGED`), `serverTlsPolicy` is the ONLY TLS mechanism — drop `sslCertificates` entirely there.
- Removing `serverTlsPolicy` later cleanly PATCHes the proxy back to plain server-side TLS.
