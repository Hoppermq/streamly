resource "zitadel_project" "default" {
  name = var.project_name
  org_id = var.organization_id
  project_role_assertion = true
  project_role_check = true
}

resource "zitadel_project_role" "default" {
  for_each = var.roles

  org_id       = var.organization_id
  project_id   = zitadel_project.default.id
  display_name = each.value.display_name
  role_key     = each.value.key
}

resource "zitadel_machine_user" "default" {
  for_each = var.services

  org_id = var.organization_id

  user_name = "${each.value.name}-local"
  name      = each.value.name
  access_token_type = "ACCESS_TOKEN_TYPE_JWT"

  with_secret = true
}

resource "zitadel_machine_key" "default" {
  for_each = zitadel_machine_user.default

  org_id  = var.organization_id
  user_id = each.value.id
  key_type = "KEY_TYPE_JSON"
}

# Look up root user by username (created by Zitadel init config)
data "zitadel_human_users" "root" {
  org_id           = var.organization_id
  user_name        = "root@streamly.auth.localhost"
  user_name_method = "TEXT_QUERY_METHOD_EQUALS"
}
