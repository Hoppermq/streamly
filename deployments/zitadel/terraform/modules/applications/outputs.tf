output "applications" {
  description = "Map of application names to their client IDs"
  value = {
    for key, app in zitadel_application_oidc.default : key => {
      client_id = app.client_id
      name      = app.name
    }
  }

  sensitive = true
}
