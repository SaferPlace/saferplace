import { LoaderFunctionArgs } from "react-router-dom"
import { getClient } from "../../hooks/client"
import { ViewerService } from "@saferplace/api/viewer/v1/viewer_connect"
import { IncidentsProps } from "./list"
import { Props as SingleProps } from './single'
import { getPosition } from "../../hooks/position"

function getQuery(request: Request, query: string): string | null {
    const url = new URL(request.url)
    return url.searchParams.get(query)
}

/**
 * incidentInRadiusLoader gets the incident in a given radius from the center. It uses the
 * search parameters to get the radius, lattitude and longitude.
 */
export async function incidentsInRadiusLoader({request}: LoaderFunctionArgs): Promise<IncidentsProps> {
    const client = getClient(ViewerService)

    const user = await getPosition()
    console.debug('user position', user)

    const radius = Number.parseFloat(getQuery(request, 'radius') ?? '10000000')
    const lat = ((): number => {
        const p = getQuery(request, 'lan')
        return p ? Number.parseFloat(p) : user.lat ?? 0
    })()
    const lon = ((): number => {
        const p = getQuery(request, 'lon')
        return p ? Number.parseFloat(p) : user.lon ?? 0
    })()

    return client.viewInRadius({
        radius: radius,
        center: {
            lat,
            lon,
        },
    })
        .then(resp => ({incidents: resp.incidents, radius, center: {lat,lon}} as IncidentsProps))
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
