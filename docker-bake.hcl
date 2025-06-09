target "default" {
  tags = ["docker.io/milas/transiter"]
  output = ["type=image,compression=zstd"]
  target = "minimal"
  inherits = RELEASE ? ["_release"] : []
}

target "_release" {
  platforms = ["linux/amd64", "linux/arm64"]
}

variable "RELEASE" {
  type    = boolean
  default = false
}
