# 📦 cfg — Flexible Configuration Loader for Go

`cfg` is a lightweight yet powerful Go package for loading, managing, and saving application configuration files.
It supports **local and remote JSON configurations**, **environment variable interpolation**, and provides a rich set of accessor functions for structured configuration data.

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
  - Flags and more
- 🔄 Environment variable interpolation using `${VAR_NAME}` syntax
- 💾 Save configuration back to disk (local files only)
- 🔁 Hot **reload** configuration without restarting the application
- 🧠 Access helpers for common lookups (e.g., `GetDatabaseInfo`, `GetEndpointInfo`)
- 🧰 Type-safe flag retrieval with `GetFlag[T]`

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

> ⚠️ Saving is only supported for **local configuration files**.

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
| `GetNotificationInfo(id string)` | Get notification by ID (uses default if empty) |
| `GetSourceInfo(id string)` | Get source by ID |
| `GetOAuthInfo(id string)` | Get OAuth provider by ID |

---

## 🏷 Flags

The configuration can include a list of `Flags` for miscellaneous values.
You can retrieve them safely with type inference:

```go
value := cfg.GetFlag[string](conf.Flags, "SOME_KEY")
timeout := cfg.GetFlag[int](conf.Flags, "TIMEOUT")
```

If the flag is missing or can't be converted, a **zero value** is returned.

---

## 🧠 Environment Variable Interpolation

All the following fields support `${ENV}` placeholders:

- `DatabaseInfo.ConnectionString`
- `EndpointInfo.Address`, `APIKey`, `Token`
- `OAuthProviderInfo.IconUrl`, `ProviderHost`, `ProviderWebUri`, `ProviderApiUri`
- `NotificationInfo.APIHost`, `Login`, `Password`, `SenderAddress`, `ReplyTo`
- `CacheInfo.Address`, `Password`

Values are interpolated **once on load**, and then restored before saving.

---

## 🧰 Struct Overview

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
	Cache                 *CacheInfo
	// ... and more
}
```

Each section (e.g. `DatabaseInfo`, `EndpointInfo`) has its own well-defined fields for easy JSON marshalling and unmarshalling.

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

## 📝 License

MIT — feel free to use, modify, and contribute.

---

## 🤝 Contributing

Pull requests and issues are welcome.
Please make sure to run `go fmt ./...` before submitting any changes.

---

## 🧠 Author

Developed and maintained by **[Elizalde Baguinon](https://github.com/eaglebush)**.
