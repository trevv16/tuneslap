import { Configuration } from './runtime';
import { getStoredToken } from '@/utils/token';

/**
 * Creates a Configuration instance with token management integration
 * The token is retrieved from localStorage and will be used for bearer authentication
 */
export function getApiConfig(): Configuration {
  const basePath = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8082/api/v1';
  
  return new Configuration({
    basePath,
    accessToken: () => {
      const token = getStoredToken();
      return token || '';
    },
  });
}
