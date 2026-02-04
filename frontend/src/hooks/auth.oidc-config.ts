import {UserManager, type UserManagerSettings, WebStorageStateStore} from 'oidc-client-ts';
import {config} from "../config/env.ts";

const ZITADEL_ISSUER: string = config.zitadelURL
const CLIENT_ID: string = config.zitadelClientID
const PROJECT_ID: string = config.zitadelProjectID

if (!ZITADEL_ISSUER) {
  console.error('Zitadel issuer missing from env', 'value', ZITADEL_ISSUER)
}
if (!CLIENT_ID) {
  console.error('application client id is missing from env', 'value', CLIENT_ID)
}
if (!PROJECT_ID) {
  console.error('Zitadel project ID missing from env', 'value', PROJECT_ID)
}

console.debug('🔐 OIDC Config:', {
  authority: ZITADEL_ISSUER,
  client_id: CLIENT_ID,
  project_id: PROJECT_ID,
  redirect_uri: `${window.location.origin}/auth/callback`,
});

const userManagerSettings: UserManagerSettings = {
  authority: ZITADEL_ISSUER,
  client_id: CLIENT_ID,
  redirect_uri: `${window.location.origin}/auth/callback`,
  post_logout_redirect_uri: window.location.origin,
  response_type: 'code',
  scope: `openid profile email offline_access urn:zitadel:iam:org:project:id:${PROJECT_ID}:aud`,
  userStore: new WebStorageStateStore({ store: window.localStorage }),
  automaticSilentRenew: true,
  silent_redirect_uri: `${window.location.origin}/auth/silent-callback`,
}

export const userManager = new UserManager(userManagerSettings)
