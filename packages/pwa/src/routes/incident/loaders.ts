import { LoaderFunctionArgs } from "react-router-dom"
import { getClient } from "../../hooks/client"
import { ViewerService } from "@saferplace/api/viewer/v1/viewer_connect"
import { IncidentsProps } from "./list"
import { Props as SingleProps } from './single'
import { getPosition } from "../../hooks/position"
import { regionsInBounds } from "../../hooks/incidents"
import { LatLngBounds } from "leaflet"
import { Incident } from "@saferplace/api/incident/v1/incident_pb"

function getQuery(request: Request, query: string): string | null {
    const url = new URL(request.url)
    return url.searchParams.get(query)
}

export async function incidentsInRegionLoader({request}: LoaderFunctionArgs): Promise<IncidentsProps> {
    const client = getClient(ViewerService)

    const user = await getPosition()
    console.debug('user position', user)

    const north = ((): number => {
        const p = getQuery(request, 'north')
        return p ? Number.parseFloat(p): (user.lat ?? 0)
    })()
    const south = ((): number => {
        const p = getQuery(request, 'south')
        return p ? Number.parseFloat(p): (user.lat ?? 0)
    })()
    const west = ((): number => {
        const p = getQuery(request, 'west')
        return p ? Number.parseFloat(p): (user.lon ?? 0)
    })()
    const east = ((): number => {
        const p = getQuery(request, 'east')
        return p ? Number.parseFloat(p): (user.lon ?? 0)
    })()

    const bounds = new LatLngBounds([south, west], [north, east])
    const incidents: Incident[] = []

    for (const region of regionsInBounds(bounds)) {
        console.info('getting region', region.toJsonString())
        client.viewInRegion({region})
            .then(resp => incidents.push(...resp.incidents))
            .catch(err => console.error(err))
    }

    console.log(incidents)
    
    return { incidents, bounds }
}

/**
 * incidentLoader loads the specified incident based on the ID.
 * @param Params contains the ID of the incident we want to load.
 * @returns Incident
 * @throws Error if we could not get an incident
 */
export async function incidentLoader({request, params}: LoaderFunctionArgs): Promise<SingleProps> {
    const isNewReport = getQuery(request, 'isNewReport') === 'true'
    const client = getClient(ViewerService)
    // return client.viewIncident({id: params.id ?? ''})
    //     .then(resp => ({resp.incident!} as SingleProps))
    return client.viewIncident({
        id: params.id ?? ''
    })
        .then(resp => ({incident: resp.incident, isNewReport} as SingleProps))
}
