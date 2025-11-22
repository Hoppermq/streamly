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
