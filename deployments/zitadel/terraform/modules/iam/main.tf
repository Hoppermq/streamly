resource "zitadel_project" "default" {
  name = var.project_name
  org_id = var.organization_id
  project_role_assertion = true
  project_role_check = true
}

resource "zitadel_project_role" "default" {
  project_id   = zitadel_project.default.id
  for_each = var.roles

  display_name = each.value.display_name
  role_key = each.value.key
}

resource "zitadel_machine_user" "default" {
  for_each = var.services

  org_id = var.organization_id

  user_name = "${each.value.name}-local"
  name      = each.value.name
  access_token_type = "ACCESS_TOKEN_TYPE_JWT"

  with_secret = true
}
