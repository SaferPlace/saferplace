import { PromiseClient, createPromiseClient, Interceptor } from '@bufbuild/connect'
import { createConnectTransport } from '@bufbuild/connect-web'
import { ServiceType } from '@bufbuild/protobuf'
import React from 'react'

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
 * emailInterceptor is a temporary interceptor designed to send the specified email
 * address until we have user authentication. In the future, the intercerceptor on the server
 * side will remove this header, authenticate, then re-add based on the verified user credentials.
 */
const emailInterceptor: Interceptor = (next) => async (req) => {
    emailHeaders(req.header)
    return await next(req)
}

function emailHeaders(headers: Headers) {
    headers.set('X-Email', localStorage.getItem('email') ?? '')
}

/**
 * getClient creates a new client
 * @param service that the client should connect to.
 */
export function getClient<T extends ServiceType>(service: T): PromiseClient<T> {
    const backend = getEndpoint()
    const transport = createConnectTransport({
        baseUrl: backend,
        interceptors: [emailInterceptor],
    })
    console.debug(`connecting to ${backend}/${service.typeName}`)
    return createPromiseClient(service, transport)
}

export async function uploadImage(image?: File): Promise<string> {
    if (!image) { return '' }
        const headers = new Headers()
        emailHeaders(headers)
        const body = new FormData()
        body.append('image', image)
        return fetch(`${getEndpoint()}/v1/upload`, {
            method: 'POST',
            body,
            headers,
        })
            .then(resp => resp.text())
}
