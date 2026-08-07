# Binary Payload

This preset creates a ConfigMap that carries a binary entry (`binaryData`, base64-encoded) alongside a regular text entry (`data`). Binary entries are for payloads that are not valid UTF-8 — icons, compiled files, keystores, serialized blobs.

## When to Use

- Distributing small binary artifacts to pods as mounted files (up to the 1MiB combined size cap)
- Payloads that would be corrupted by UTF-8 handling — anything that is bytes, not text
- Pairing a binary artifact with text metadata describing it, in one object

## Key Configuration Choices

- **`binaryData` values are base64** — supply the exact base64 encoding of the file's bytes (`base64 -w0 favicon.ico`); the schema validates base64 well-formedness before anything reaches the cluster
- **Binary keys mount as files only** — Kubernetes cannot expose `binaryData` keys as environment variables; consume them via a `configMap` volume
- **No key overlap** — a key may live in `data` OR `binaryData`, never both; the schema enforces the same rule the API server does, at validation time
- **Same key character rules as `data`** — alphanumeric, `-`, `_`, `.`; a file-style key (`favicon.ico`) becomes the mounted file's name

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace for the ConfigMap (must match the consuming workload's namespace) | Your namespace management |

The `favicon.ico` value is a small sample base64 string — replace it with your file's real encoding (`base64 -w0 <file>` on Linux, `base64 -i <file>` on macOS), and replace the `asset-manifest` text entry with metadata relevant to your payload (or remove it).

## Related Presets

- **01-app-config** — text configuration with scalar settings and a properties file
- **02-immutable-versioned** — locked, versioned config for production roll-forward
