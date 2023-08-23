/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_BACKEND: string
    readonly VITE_CDN: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}
