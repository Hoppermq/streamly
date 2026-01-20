output "project_id" {
  description = "The ID of the created project"
  value       = zitadel_project.default.id
}

output "service_user_ids" {
  description = "Map of service names to their user IDs"
  value       = {
    for k, v in zitadel_machine_user.default : k => v.id
  }
}

output "role_ids" {
  description = "Map of role keys to their IDs"
  value       = {
    for k, v in zitadel_project_role.default : k => v.id
  }
}

output "service_credentials" {
  description = "Service account credentials (client_id and client_secret)"
  sensitive   = true
  value = {
    for k, key in zitadel_machine_key.default : k => {
      client_id     = key.user_id
      client_secret = key.key_details
    }
  }
}

output "root_admin_credentials" {
  description = "Root Admin human user account credentials (user_id and PAT)"
  sensitive   = true
  value = {
    user_id = data.zitadel_human_users.root.id
    pat     = zitadel_personal_access_token.root.token
  }
}
