import { LoaderFunctionArgs } from "react-router-dom"
import { getClient } from "../../hooks/client"
import { ViewerService } from "@saferplace/api/viewer/v1/viewer_connect"
import { Incident } from "@saferplace/api/incident/v1/incident_pb"
import { IncidentsProps } from "./list"

async function getPosition(): Promise<{lat: number, lon: number}> {
    return new Promise((resolve, reject) => {
        if (!navigator.geolocation) {
            return reject('No geolocation')
        }

        navigator.geolocation.getCurrentPosition(
            (position) => console.log(position.coords),
            (error) => {
                console.error(error)
                // If we cannot get the location default to The Criminal Court of Justice
                resolve({lat: 53.34868617902951, lon: -6.29567143778413})
            },
        )
    })
}

/**
 * incidentInRadiusLoader gets the incident in a given radius from the center. It uses the
 * search parameters to get the radius, lattitude and longitude.
 */
export async function incidentsInRadiusLoader({request}: LoaderFunctionArgs): Promise<IncidentsProps> {
    const client = getClient(ViewerService)

    const user = await getPosition()
    console.debug('user position', user)
    
    const url = new URL(request.url)
    const radius = Number.parseFloat(url.searchParams.get('radius') ?? '10000000')
    const lat = ((): number => {
        const p = url.searchParams.get('lat')
        return p ? Number.parseFloat(p) : user.lat
    })()
    const lon = ((): number => {
        const p = url.searchParams.get('long')
        return p ? Number.parseFloat(p) : user.lon
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
export async function incidentLoader({params}: LoaderFunctionArgs): Promise<Incident> {
    const client = getClient(ViewerService)
    return client.viewIncident({id: params.id ?? ''})
        .then(resp => resp.incident!)
}