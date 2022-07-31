# Use the safer place terraform module, so we have to replicate everything here.

terraform {
  # TODO: Eventually migrate to cloud block
  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "saferplace"

    workspaces {
      name = "services"
    }
  }
}

variable "nomad_address" {
  description = "Nomad address"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "Deployment region"
  type        = string
}

variable "datacenters" {
  description = "List of datacenters to deploy"
  type        = list(string)
}

variable "namespace" {
  description = "Deployment namespace"
  type        = string
}

variable "image" {
  description = "Deployment image"
  type        = string
}

variable "docker_password" {
  description = "Read only token for the container registry"
  type        = string
  sensitive   = true
}

variable "eircode_token" {
  description = "Token used to authenticate to the eircode service"
  type        = string
  sensitive   = true
}


provider "nomad" {
  address = var.nomad_address
}

module "saferplace" {
  source = "./deploy/terraform"

  region          = var.region
  datacenters     = var.datacenters
  namespace       = var.namespace
  image           = var.image
  docker_password = var.docker_password
  eircode_token   = var.eircode_token
}
