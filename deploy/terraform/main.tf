resource "nomad_namespace" "saferplace" {
  name        = var.namespace
  description = "Safer Place namespace"
}

module "saferplace" {
  source  = "Voytechnology/generic/nomad"
  version = "0.0.4"

  job_name    = "saferplace"
  namespace   = var.namespace
  region      = var.region
  datacenters = var.datacenters

  image = var.image
  ports = {
    "http" = {
      to     = 8080
      static = 0
    }
  }
  docker_username = "Voytechnology"
  docker_password = var.docker_password
  env = {
    "SAFERPLACE_ADDRESS_RESOLVERS_EIRCODE_TOKEN" = var.eircode_token
  }

  service_port = "http"
  service_tags = [
    "traefik.http.routers.saferplace.rule=Host(`safer.place`)",
    "prometheus.io/scrape=false"
  ]
}
