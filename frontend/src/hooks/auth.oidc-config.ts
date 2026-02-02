import {UserManager, type UserManagerSettings, WebStorageStateStore} from 'oidc-client-ts';

const ZITADEL_ISSUER = "http://auth.localhost:8080"
const CLIENT_ID = "357291801155035402"

const userManagerSettings: UserManagerSettings = {
  authority: ZITADEL_ISSUER,
  client_id: CLIENT_ID,
  redirect_uri: `${window.location.origin}/auth/callback`,
  post_logout_redirect_uri: window.location.origin,
  response_type: 'code',
  scope: 'openid profile email offline_access',
  userStore: new WebStorageStateStore({ store: window.localStorage }),
  automaticSilentRenew: true,
  silent_redirect_uri: `${window.location.origin}/auth/silent-callback`,
}
export const userManager = new UserManager(userManagerSettings)
