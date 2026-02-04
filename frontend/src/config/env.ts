const getEnvVar = (key: string, defaultValue?: string): string => {
  const value = import.meta.env[key] || defaultValue;
  if (!value) {
    console.warn(`Environment variable ${key} is not set`);
  }
  return value || '';
};

export const config = {
  zitadelURL: getEnvVar('VITE_ZITADEL_ISSUER', 'http://auth.localhost:8080'),
  zitadelClientID: getEnvVar('VITE_ZITADEL_CLIENT_ID', ''),
  apiURL: getEnvVar('VITE_API_URL', 'http://localhost:8080'),
  environment: getEnvVar('VITE_APP_ENV', 'development'),
} as const;
