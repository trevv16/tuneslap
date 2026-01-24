// This is a stub file for the OpenAPI generated runtime module
// The actual runtime.ts is generated during build time via `yarn generate-types`
// This stub allows tests to run without requiring the generated code

export class Configuration {
  basePath: string
  accessToken?: () => string

  constructor(config: { basePath?: string; accessToken?: () => string } = {}) {
    this.basePath = config.basePath ?? ''
    this.accessToken = config.accessToken
  }
}
