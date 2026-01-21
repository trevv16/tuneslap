import { Configuration } from './runtime';
import { getStoredToken } from '@/utils/token';

const PRODUCTION_API_URL = 'https://api.tuneslap.com/api/v1';

/**
 * Creates a Configuration instance for client-side requests with authentication
 */
export function getApiConfig(): Configuration {
  const basePath = process.env.NEXT_PUBLIC_API_URL || PRODUCTION_API_URL;
  
  return new Configuration({
    basePath,
    accessToken: () => {
      const token = getStoredToken();
      return token || '';
    },
  });
}

/**
 * Creates a Configuration instance for server-side requests (SSR/ISR)
 * Uses INTERNAL_API_URL for Docker container-to-container communication,
 * falls back to NEXT_PUBLIC_API_URL for external access
 */
export function getServerApiConfig(): Configuration {
  const basePath = process.env.INTERNAL_API_URL || process.env.NEXT_PUBLIC_API_URL || PRODUCTION_API_URL;
  return new Configuration({ basePath });
}
