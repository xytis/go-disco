terraform {
  backend "local" {
    path = ".dev/terraform.tfstate"
  }
}

provider "docker" {}

variable "consul_http_port" {
  type    = "string"
  default = "8500"
}

variable "consul_dns_port" {
  type    = "string"
  default = "8600"
}

resource "docker_container" "consul" {
  name  = "consul"
  image = "${docker_image.consul.latest}"

  ports {
    internal = 8500
    external = "${var.consul_http_port}"
  }

  ports {
    internal = 8600
    external = "${var.consul_dns_port}"
  }
}

resource "docker_image" "consul" {
  name = "consul"
}
