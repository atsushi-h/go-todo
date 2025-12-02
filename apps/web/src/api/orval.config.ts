import { defineConfig } from 'orval'

export default defineConfig({
  api: {
    input: {
      target: '../../../api/openapi/openapi.yaml',
    },
    output: {
      mode: 'tags',
      target: 'generated',
      client: 'react-query',
      mock: false,
      clean: true,
      override: {
        mutator: {
          path: './axios-instance.ts',
          name: 'customAxiosInstance',
        },
        query: {
          version: 5,
          signal: true,
        },
      },
      headers: true,
    },
  },
})
