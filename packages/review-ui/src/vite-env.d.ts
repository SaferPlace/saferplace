/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_BACKEND: string
    readonly VITE_CDN: string
    readonly VITE_OIDC_AUTHORITY: string
    readonly VITE_OIDC_CLIENT_ID: string
    readonly VITE_OIDC_REDIRECT_URL: string
    readonly VITE_OIDC_CLIENT_SECRET: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}
