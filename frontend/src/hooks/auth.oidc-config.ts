import {UserManager, type UserManagerSettings, WebStorageStateStore} from 'oidc-client-ts';
import {config} from "../config/env.ts";

const ZITADEL_ISSUER: string = config.zitadelURL
const CLIENT_ID: string = config.zitadelClientID

// Validate required config
if (!ZITADEL_ISSUER) {
  console.error('Zitadel issuer missing from env', 'value', ZITADEL_ISSUER)
}
if (!CLIENT_ID) {
  console.log('application client id is missing from env', 'value', CLIENT_ID)
}

console.log('🔐 OIDC Config:', {
  authority: ZITADEL_ISSUER,
  client_id: CLIENT_ID,
  redirect_uri: `${window.location.origin}/auth/callback`,
});


const projectid = '356531635715399939'

const userManagerSettings: UserManagerSettings = {
  authority: ZITADEL_ISSUER,
  client_id: CLIENT_ID,
  redirect_uri: `${window.location.origin}/auth/callback`,
  post_logout_redirect_uri: window.location.origin,
  response_type: 'code',
  scope: `openid profile email offline_access urn:zitadel:iam:org:project:id:${projectid}:aud`,
  userStore: new WebStorageStateStore({ store: window.localStorage }),
  automaticSilentRenew: true,
  silent_redirect_uri: `${window.location.origin}/auth/silent-callback`,
}

export const userManager = new UserManager(userManagerSettings)
