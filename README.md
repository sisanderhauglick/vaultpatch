# vaultpatch

> CLI tool to diff and apply HashiCorp Vault secret changes across environments

---

## Installation

```bash
go install github.com/yourusername/vaultpatch@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpatch/releases).

---

## Usage

**Diff secrets between two environments:**

```bash
vaultpatch diff --src secret/prod/app --dst secret/staging/app
```

**Apply changes from a diff file:**

```bash
vaultpatch apply --patch changes.patch --target secret/staging/app
```

**Export secrets to a patch file:**

```bash
vaultpatch export --path secret/prod/app --out changes.patch
```

> Requires `VAULT_ADDR` and `VAULT_TOKEN` environment variables to be set, or a configured `~/.vault-token`.

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.12+

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)