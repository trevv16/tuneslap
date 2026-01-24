// MSW server setup for Node.js test environment
import { setupServer } from 'msw/node'
import { handlers } from './msw-handlers'

export const server = setupServer(...handlers)
