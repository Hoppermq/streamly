resource "zitadel_user_grant" "service_grants" {
  for_each = var.services
  
  user_id    = var.service_user_ids[each.value.service_key]
  project_id = var.project_id
  role_keys  = values(each.value.roles)
}
