# holavonat.is

Real-time train tracking visualization for Hungarian public transport. This application collects and displays train positions, delays, schedules, and other related data through an undocumented public transport API.

## Features
- Real-time train tracking
- Delay information and schedule data
- Interactive map visualization
- Historical data collection
- Flexible data distribution:
  - S3-compatible storage support (e.g., Cloudflare R2)
  - Local file system storage
- Highly configurable client application:
  - Adjustable update intervals
  - Configurable API endpoints and parameters

## Honorable Mention

This project pays tribute to [holavonat.hu](https://holavonat.hu), the first open-source project that pioneered third-party train tracking visualization in Hungary. Unfortunately, due to concerns about potential governmental actions, the original author had to discontinue the service. Their courage to innovate and their commitment to open-source principles continue to inspire this and similar projects.


## Configuration

The application is configured using a YAML file. Here's a detailed explanation of the configuration parameters:

### Basic Settings
```yaml
GraphqlEndpoint: "https://example.com/graphql"  # Upstream API endpoint
```

### Output Configuration
```yaml
Output:
  NamePrefix: "data"                # Base filename for outputs
  Format:
    JSON: true                      # Currently only JSON is supported (future: ProtoBuf)
  Archive: true                     # Enable ISO8601 suffixed archive files
```
When Archive is enabled, files are saved as: `{NamePrefix}_{ISO8601}.json`

### Distribution Modes

#### S3-Compatible Storage
```yaml
ObjectStorage:
  Compression: "br"                 # Supported: gzip, br, zstd (default: none)
  AccessKeyID: "<access-key>"
  SecretAccessKey: "<secret-key>"
  BucketName: "<bucket-name>"
  ObjectPath: "<optional-path>"
  EndpointURL: "<s3-endpoint>"
  PublicEndpointURL: "https://cdn.example.com"
```

#### File System
```yaml
File:
  Path: "/path/to/output"          # Output directory (created if not exists)
```

### API Communication
```yaml
Headers:
  User-Agent: "holavonatis/v0.0.1 (https://instance.example.com/)"
  Referrer: "https://instance.example.com/"
  Origin: "https://instance.example.com"
```

### Schedule Configuration
```yaml
Cron:
  Mode: "fix"                      # Options: "fix" or "window"
  Duration: "second"               # Options: "second", "minute", "hour"
  Fix:
    Interval: 30                   # Run every 30 seconds
  # OR for window mode:
  # Window:
  #   Min: 30                      # Minimum interval
  #   Max: 45                      # Maximum interval (random selection)
```

### Source Information
```yaml
Source:
  Origin: "https://instance.example.com/"    # Instance visualization URL
  Latest: "https://instance.example.com/"    # Base URL for latest data
  Schema:
    Version: "v0.0.1"                       # Schema version
    Link: "https://instance.example.com/schema.json"
    Format: "json"                          # Currently only JSON supported
```

### Network Settings
```yaml
Network:
  Proxy: "socks5://127.0.0.0:3124"         # Optional proxy for API communication
```

For a complete example, see [config_example.yaml](config_example.yaml).

## Live Data Source

**Currently Available**: A live data feed is accessible at https://cdn.holavonat.is/train_data_v3.json

### Data Access
- **CORS Policy**: Unrestricted (`*`) - accessible from any domain
- **Format**: JSON (following the v3 schema)
- **Update Frequency**: Real-time updates (typically every 60 seconds)
- **Public Use**: Available for third-party applications and research

### Important Disclaimers

⚠️ **Data Reliability Warning**: While this data source is currently available for public use, please be aware that:

- **No Service Guarantee**: The data feed may become unavailable at any time without notice
- **Data Accuracy**: The authors cannot guarantee the accuracy, completeness, or timeliness of the data
- **No Liability**: Users of this data source assume all responsibility for its use
- **Third-Party Risk**: The underlying data comes from an undocumented API and is subject to the same risks outlined in the main disclaimer

### Usage Recommendations
- Implement proper error handling for data unavailability
- Cache data locally when possible to reduce dependency
- Consider this a convenience service, not a guaranteed API
- Always have fallback mechanisms in your applications

For schema information, see [internal/api/schema.go](internal/api/schema.go).

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). This means:

- You can freely use, modify, and distribute this software
- If you modify and distribute the software, you must:
  - Make your modifications available under the AGPL-3.0
  - Provide the complete source code to users who interact with the software over a network
  - Include the original copyright notice and license

For the complete license terms, see the [LICENSE](LICENSE) file.

## Third-Party Licenses and Attributions

This project uses the following third-party data and services:

### Map Data
- **OpenStreetMap**: © OpenStreetMap contributors. The map data is available under the [Open Data Commons Open Database License (ODbL)](https://www.openstreetmap.org/copyright).
- **OpenRailwayMap**: © OpenRailwayMap and contributors. The railway data is available under the [Open Data Commons Open Database License (ODbL)](https://www.openrailwaymap.org/imprint-en.html).

### Libraries
- **Leaflet**: © Vladimir Agafonkin. Leaflet is used for map rendering and is available under the [BSD 2-Clause License](https://github.com/Leaflet/Leaflet/blob/main/LICENSE).

For complete licensing information of all dependencies, please refer to the vendor directory and respective package licenses.

## Disclaimer

⚠️ **IMPORTANT: PLEASE READ CAREFULLY** ⚠️

By using this software, you acknowledge and agree to the following:

1. This application interfaces with an undocumented API that was not officially intended for third-party use.

2. **Risks**:
   - The API provider may restrict or block access at any time without notice
   - The API provider may take legal action against end users
   - The API provider may seek damages for unauthorized API usage
   - The application may cease to function without prior notice

3. **No Warranty**:
   - The software is provided "AS IS" without warranty of any kind
   - The authors make no warranties regarding functionality, reliability, or legality
   - No guarantee of continued functionality or service availability

4. **Liability and Indemnification**:
   - Users assume all risks associated with using this software
   - The authors are not liable for any damages or legal consequences
   - Users agree to indemnify and hold harmless the authors from any claims

By using this software, you accept all these terms and conditions. If you do not agree, you must not use this software.

For complete terms, please see the [EULA](internal/config/eula.go) file.
