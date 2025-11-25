# 📦 cfg — Flexible Configuration Loader for Go

`cfg` is a lightweight yet powerful Go package for loading, managing, and saving application configuration files. It supports **local and remote JSON configurations**, **environment variable interpolation**, and provides a rich set of accessor functions for structured configuration data.

---

## ✨ Features

- ✅ Load configuration from **local files** or **HTTP(S) URLs**
- ✅ Supports structured configuration for:
  - Databases
  - API Endpoints
  - OAuth Providers
  - Directories
  - Notifications
  - Caching
  - Secrets
  - Queues
  - Flags
  - FlagGroups and more
- 🔄 Environment variable interpolation using `${VAR_NAME}` syntax
- 💾 Save configuration back to disk (local files only)
- 🔁 Hot **reload** configuration without restarting the application
- 🧠 Access helpers for common lookups (e.g., `GetDatabaseInfo`, `GetEndpointInfo`)
- 🧰 Type-safe flag retrieval with `GetFlag[T]`
- 🧩 Built-in defaults for common configuration values (e.g., IDs, cookie domain, secrets)

---

## 📥 Installation

```bash
go get github.com/eaglebush/config
```

> Replace `github.com/eaglebush/config` with your actual module path.

---

## 🧪 Basic Usage

### Load a Configuration

```go
package main

import (
	"fmt"
	"log"

	"github.com/eaglebush/config"
)

func main() {
	conf, err := cfg.Load("./config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Println("App Name:", *conf.ApplicationName)

	// Access default database
	db := conf.GetDatabaseInfo(*conf.DefaultDatabaseID)
	if db != nil {
		fmt.Println("DB ID:", db.ID)
		fmt.Println("Conn String:", db.ConnectionString)
	}
}
```

---

### Using Environment Variables

You can use `${VAR_NAME}` placeholders in your configuration file:

```json
{
	"Databases": [
		{
			"ID": "main",
			"ConnectionString": "Server=${DB_HOST};User Id=${DB_USER};Password=${DB_PASS};"
		}
	]
}
```

At runtime, `cfg` automatically interpolates these values from the environment.

```bash
export DB_HOST=localhost
export DB_USER=admin
export DB_PASS=secret
```

---

### Saving Configuration

```go
err := conf.Save()
if err != nil {
	log.Println("Save failed:", err)
}
```

> ⚠️ Saving is only supported for **local configuration files**. Remote URLs are read-only.

---

### Hot Reloading

Reload the configuration from the original source without restarting:

```go
err := conf.Reload()
if err != nil {
	log.Println("Reload failed:", err)
}
```

---

## 🧰 Lookup Functions

The `Configuration` struct exposes various helper methods:

| Method | Description |
|--------|-------------|
| `GetDatabaseInfo(id string)` | Get database by ID |
| `GetDatabaseInfoGroup(groupId string)` | Get databases by group |
| `GetEndpointInfo(id string)` | Get API endpoint by ID (uses default if empty) |
| `GetEndpointInfoGroup(groupId string)` | Get endpoints by group |
| `GetDirectory(groupId string)` | Get directory by group ID |
| `GetDirectoryItem(groupId, key string)` | Get specific directory item |
| `GetDomainInfo(name string)` | Get domain info |
| `GetFlagGroupFlags(groupId string)` | Gets flags from defined group |
| `GetNotificationInfo(id string)` | Get notification by ID (uses default if empty) |
| `GetSourceInfo(id string)` | Get source by ID |
| `GetOAuthInfo(id string)` | Get OAuth provider by ID |
| `GetSecretInfo(id string)` | Get secret by ID |
| `GetSecretInfoGroup(groupId string)` | Get secrets by group ID |

---

## 🏷 Flags

The configuration can include a list of `Flags` for miscellaneous values. You can retrieve them safely with type inference:

```go
value := cfg.GetFlag[string](conf.Flags, "SOME_KEY")
timeout := cfg.GetFlag[int](conf.Flags, "TIMEOUT")
```

If the flag is missing or can't be converted, a **zero value** is returned.

Additionally, `Configuration.Flag(key)` provides a convenient way to find a flag by key, ignoring underscores and dashes.

---

## 🧠 Environment Variable Interpolation

The following fields support `${ENV}` placeholders:

- `DatabaseInfo.ConnectionString`
- `EndpointInfo.Address`, `APIKey`, `Token`
- `OAuthProviderInfo.IconUrl`, `ProviderHost`, `ProviderWebUri`, `ProviderApiUri`
- `NotificationInfo.APIHost`, `Login`, `Password`, `SenderAddress`, `ReplyTo`
- `CacheInfo.Address`, `Password`

Values are interpolated **once on load**, and then **restored to original values** before saving.

### Example
```json
{
  "APIEndpoints": [
    {
      "ID": "DEFAULT",
      "Address": "https://${API_HOST}/v1",
      "APIKey": "${API_KEY}"
    }
  ]
}
```

---

## 🧩 Struct Overview

`Configuration` is the central struct:

```go
type Configuration struct {
	ApplicationID         *string
	ApplicationName       *string
	Databases             *[]DatabaseInfo
	APIEndpoints          *[]EndpointInfo
	Notifications         *[]NotificationInfo
	OAuths                *[]OAuthProviderInfo
	Directories           *[]DirectoryInfo
	Flags                 *[]Flag
	FlagGroups            *[]FlagGroup
	Cache                 *CacheInfo
	Secrets               *[]SecretInfo
	Sources               *[]SourceInfo
	Queue                 *QueueInfo
	CookieDomain          *string
	JWTSecret             *string
	Secure                *bool
	ReadTimeout           *int
	WriteTimeout          *int
}
```

### Supporting Types

- `DatabaseInfo` — defines connection details, pooling, and SQL formatting options.
- `EndpointInfo` — represents API endpoints with optional `APIKey`, `Token`, and attached secrets.
- `OAuthProviderInfo` — stores OAuth2 client and provider details.
- `NotificationInfo` — holds notification service connection settings.
- `CacheInfo` — caching backend information.
- `QueueInfo` — queue/streaming configuration.
- `SecretInfo` — grouped secrets for secure data.
- `DirectoryInfo` — configuration for grouped flags.

---

## ⚡ Example JSON Configuration

```json
{
	"ApplicationName": "MyApp",
	"Databases": [
		{
			"ID": "main",
			"DriverName": "sqlserver",
			"ConnectionString": "Server=${DB_HOST};User Id=${DB_USER};Password=${DB_PASS};"
		}
	],
	"APIEndpoints": [
		{
			"ID": "DEFAULT",
			"Name": "Main API",
			"Address": "https://api.example.com"
		}
	],
	"Flags": [
		{ "Key": "TIMEOUT", "Value": "30" }
	]
}
```

---

## ⚙️ Defaults

When a configuration is loaded, default values are applied where missing:

| Field | Default Value |
|--------|----------------|
| `CookieDomain` | `localhost` |
| `JWTSecret` | `defaultsecretkey` |
| `DatabaseInfo.StorageType` | `SERVER` |
| `DatabaseInfo.InterpolateTables` | `true` |
| `DatabaseInfo.StringEnclosingChar` | `'` |
| `DatabaseInfo.StringEscapeChar` | `\\` |
| `DatabaseInfo.ReservedWordEscapeChar` | `"` |

---

## ⚠️ Error Handling

| Error | Meaning |
|-------|----------|
| `ErrNoDataFromSource` | No data was read from the configuration source |
| `ErrSaveNotLocalFile` | Tried to save a remote (non-local) configuration |

---

## 🧩 Internal Helpers

### `interpolateEnvVars`
Replaces `${ENV_VAR}` tokens in text with values from `os.Getenv()`.

### `findByID`
Generic helper that finds structs in a slice by a string ID (case-insensitive).

---

## 🧠 Thread Safety
The `Configuration` type is **not thread-safe**. Use synchronization (like `sync.RWMutex`) when accessing or modifying configuration concurrently.

---

## 🧰 Testing Tips

- Use temporary files for `Load()`/`Save()` tests.
- Use `t.Setenv()` to mock environment variables.
- Confirm that saved JSON preserves **uninterpolated** values.
- Validate that in-memory state remains **interpolated** post-save.

---

## 🧾 License

MIT — feel free to use, modify, and contribute.

---

## 🤝 Contributing

Pull requests and issues are welcome. Please run `go fmt ./...` before submitting any changes.

---

## 👤 Author

Developed and maintained by **[Elizalde Baguinon](https://github.com/eaglebush)**.

