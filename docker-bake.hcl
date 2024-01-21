// Alpine version
variable "ALPINE_VERSION" {
  default = "latest"
}

variable "GOLANG_VERSION" {
  default = "latest"
}

target "args" {
  args = {
    ALPINE_VERSION = ALPINE_VERSION
    GOLANG_VERSION = GOLANG_VERSION
  }
}

target "platforms" {
  platforms = [
    "linux/386",
    "linux/amd64",
    "linux/arm64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/ppc64le",
    "linux/s390x"
  ]
}

// Special target: https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {
  tags = ["failed-pod-cleaner:local"]
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["args", "docker-metadata-action"]
}

target "image-local" {
  inherits = ["image"]
  output   = ["type=docker"]
}

target "image-all" {
  inherits = ["platforms", "image"]
}
