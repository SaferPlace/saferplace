import { PromiseClient, createPromiseClient } from "@bufbuild/connect"
import { createConnectTransport } from "@bufbuild/connect-web"
import { ServiceType } from "@bufbuild/protobuf"
import React from "react"

export default function useClient<T extends ServiceType>(service: T): PromiseClient<T> {
    return React.useMemo(() => {
        const backend = localStorage.getItem('backend') ?? 'https://api.safer.place'
        const transport = createConnectTransport({
            baseUrl: backend
        })
        console.debug(`connecting to ${backend}/${service.typeName}`)
        return createPromiseClient(service, transport)
    }, [service])
}