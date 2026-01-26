resource "zitadel_application_oidc" "default" {
  for_each = var.applications
  project_id    = var.project_id
  redirect_uris = var.redirect_uris

  version = each.value.version
  grant_types    = each.value.grant_types
  name           = each.value.name
  response_types = each.value.response_types
  auth_method_type = each.value.auth_method_type
  access_token_type = each.value.access_token_type
  app_type = each.value.type
}
