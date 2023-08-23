import { PromiseClient, createPromiseClient } from "@bufbuild/connect"
import { createConnectTransport } from "@bufbuild/connect-web"
import { ServiceType } from "@bufbuild/protobuf"
import React from "react"

export default function useClient<T extends ServiceType>(service: T): PromiseClient<T> {
    return React.useMemo(() => getClient(service), [service])
}

export function getEndpoint(): string {
    return localStorage.getItem('backend') ?? import.meta.env.VITE_BACKEND
}

export function getCDNEndpoint(): string {
    return localStorage.getItem('cdn') ?? import.meta.env.VITE_CDN
}

/**
 * getClient creates a new client
 * @param service that the client should connect to.
 */
export function getClient<T extends ServiceType>(service: T): PromiseClient<T> {
    const backend = getEndpoint()
    const transport = createConnectTransport({
        baseUrl: backend
    })
    console.debug(`connecting to ${backend}/${service.typeName}`)
    return createPromiseClient(service, transport)
}
