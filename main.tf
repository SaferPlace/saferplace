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

variable "tags" {
  description = "List of service tags"
  type        = list(string)
}

provider "nomad" {
  address = var.nomad_address
}

module "saferplace" {
  source = "./deploy/terraform-nomad"

  region      = var.region
  datacenters = var.datacenters
  namespace   = var.namespace
  image       = var.image
  tags        = var.tags
}
