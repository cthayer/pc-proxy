# pc-proxy
Parental control proxy for controlling access to websites

Allows blocking access to websites by matching the host, url, or path of the website being requested.

Also allows a blocked site to be accessed by providing a password in the browser.

## Configuration

Configuration is passed to the service by passing the path to a configuration file when you start the the service.  The configuration file must be in `JSON` or `HCL` format.

```bash
pc-proxy --config-file /path/to/config.hcl
```

Here is an example configuration file:

```hcl
# can specify as many rules as needed
rules {
  # can be: "block" or "allow" (default: "block")
  access = "block"

  # can be "host", "path", or "url" (default: "host")
  type = "host"

  # regex pattern to match against the value of `type`
  pattern = "example\\.com"

  # when using `access = "block"` allow the site to be accessed if the correct password is provided (default: true)
  passwordBypass = true
}

tls {
  # enable support for TLS (default: false)
  enabled = false
  cert = ""
  key = ""

  # allowed ciphers (default: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384")
  ciphers = "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
}

logging {
  # what logging level to use.  can be: "error", "warn", "info", or "debug" (default: "info")
  level = "info"

  # can be: "console" or "json" (default: "console")
  encoding = "console"
}

listen {
  host = "0.0.0.0"
  port = 80
  tlsPort = 443
}
```

## Building

[Mage](https://magefile.org/) and [gox](https://github.com/mitchellh/gox) are used for building binaries.  They must be installed on the system running the build.

To build the binary execute the following command:

```bash
build/build.sh
```

The binaries, zip archives, and SHA256 checksum files will be written to the `build/bin` directory.

Set the `BUILD_DOCKER` env variable to `1` to build docker images for the binaries.

```bash
BUILD_DOCKER=1 build/build.sh
```

### Cleaning up

To clean up the results of a build and start fresh, run the following command:

```bash
mage clean
```

### Version

Set the version for a binary by editing the `build/versions.json` file.
