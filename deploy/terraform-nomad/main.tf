resource "nomad_namespace" "saferplace" {
  name        = var.namespace
  description = "Safer Place namespace"
}

resource "nomad_job" "saferplace" {
  jobspec = file("${path.module}/../nomad/saferplace.nomad")

  hcl2 {
    enabled = true

    # Use a reduced variable set until needed
    # It looks like list(string) types cannot actually be list(string), but
    # have to be converted to strings. Let's see will it then convert back
    # to the correct type.
    vars = {
      "region"       = var.region,
      "datacenters"  = jsonencode(var.datacenters),
      "namespace"    = var.namespace,
      "image"        = var.image,
      "service_tags" = jsonencode(var.tags),
    }
  }
}
