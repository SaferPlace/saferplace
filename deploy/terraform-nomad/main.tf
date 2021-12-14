resource "nomad_namespace" "saferplace" {
  name        = var.namespace
  description = "Safer Place namespace"
}

resource "nomad_job" "saferplace" {
  jobspec = file("${path.module}/../nomad/saferplace.nomad")

  hcl2 {
    enabled = true

    # Use a reduced variable set until needed
    vars = {
      "region"       = var.region,
      "datacenters"  = var.datacenters,
      "namespace"    = var.namespace,
      "image"        = var.image,

    }

    # This is not working for some reason, so for now we will hardcode the
    # service tags inside the job.
    # vars {
    #   "service_tags" = var.tags,
    # }
  }
}
