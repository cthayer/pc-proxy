# can specify as many rules as needed
rules {
  # can be: "block" or "allow" (default: "block")
  access = "block"

  # can be "host", "path", or "url" (default: "host")
  type = "host"

  # regex pattern to match against the value of `type`
  pattern = ""

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
