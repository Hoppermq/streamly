locals {
  roles = {
    super_admin = {
      key = "super_admin"
      display_name = "Super Admin"
    },
    admin = {
      key = "admin"
      display_name = "Admin"
    },
    technical = {
      key = "technical"
      display_name = "Technical"
    }
  }

  services = {
    auth = {
      scopes =["openid", "streamly:events:read", "streamly:events:write"],
      roles = ["technical"],
    },
    ingestor = {
      scopes =["openid", "streamly:events:read", "streamly:events:write"],
      roles = ["technical"]
    },
    query ={
      scopes = ["openid", "streamly:events:read"],
      roles = ["technical"],
    } ,
    processor = {
      scopes = ["openid", "streamly:events:read", "streamly:events:write"],
      roles = ["technical"]
    }
    platform = {
      scopes = ["openid", "streamly:events:read", "streamly:events:write"],
      roles = ["technical"]
    }
  }

  applications = {
    web-oidc = {
      type = "OIDC_APP_TYPE_USER_AGENT"  # SPA - no server-side code
      name = "web-oidc"
      response_types = [
        "OIDC_RESPONSE_TYPE_CODE"
      ]
      grant_types = [
        "OIDC_GRANT_TYPE_AUTHORIZATION_CODE",
        "OIDC_GRANT_TYPE_REFRESH_TOKEN"
      ]
      auth_method_type = "OIDC_AUTH_METHOD_TYPE_NONE"  # Public client - PKCE only
      access_token_type = "OIDC_TOKEN_TYPE_JWT"
      version = "OIDC_VERSION_1_0"
    }
  }

  # Derive service definitions for IAM module (machine user creation)
  service_definitions = {
    for key, config in local.services : key => {
      name = key
    }
  }

  # Derive service role mappings for service-accounts module
  service_role_mappings = {
    for key, config in local.services : key => {
      service_key = key
      roles = {
        for role in config.roles : role => role
      }
    }
  }
}
