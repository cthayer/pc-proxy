# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "zoom\\.us"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "google\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "seesaw\\.(com|me)"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "getepic\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "clever\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "i-ready\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "mcgraw-hill\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "pathtoexcel\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "mheducation\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "brainpop\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "mysteryscience\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "allow"
  type = "host"
  pattern = "gonoodle\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "block"
  type = "host"
  pattern = "mathplayground\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "block"
  type = "host"
  pattern = "youtube\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "block"
  type = "host"
  pattern = "coolmathgames\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "block"
  type = "host"
  pattern = "cartoonnetwork(me)?\\.com"
  passwordBypass = true
}
# can specify as many rules as needed
rules {
  access = "block"
  type = "host"
  pattern = ".*"
  passwordBypass = true
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
}