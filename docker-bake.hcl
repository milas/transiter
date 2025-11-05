target "default" {
  tags = ["docker.io/milas/transiter"]
  output = ["type=image,compression=zstd"]
  target = "minimal"
  inherits = RELEASE ? ["_release"] : []
}

target "codegen" {
  name = "codegen-${tgt}"
  matrix = {
    tgt = ["go", "docs"]
  }
  target = "codegen-${tgt}"
  output = ["type=local,dest=."]
}

target "_release" {
  platforms = ["linux/amd64", "linux/arm64"]
}

variable "RELEASE" {
  type    = boolean
  default = false
}
