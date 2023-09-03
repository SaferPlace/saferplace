import { ViewerService } from "@saferplace/api/viewer/v1/viewer_connect"
import { Region } from "@saferplace/api/viewer/v1/viewer_pb"
import { LatLngBounds } from "leaflet"
import useClient from "./client"
import React from "react"
import { Incident } from "@saferplace/api/incident/v1/incident_pb"

const MAX_REGION_LIMIT = 10

/**
 * regionsInBounds returns all the regions we should be getting in the specific map bounds.
 * @param bounds 
 * @returns regions
 */
export function regionsInBounds(bounds?: LatLngBounds): Region[] {
    if (!bounds) {
        return []
    }
    const maxNorth = Math.round(bounds.getNorth() * 100)
    const maxSouth = Math.round(bounds.getSouth() * 100)
    const maxWest = Math.round(bounds.getWest() * 100)
    const maxEast = Math.round(bounds.getEast() * 100)

    console.log(`getting all regions between n:${maxNorth}, s:${maxSouth}, w: ${maxWest}, e: ${maxEast}`)

    if (maxNorth - maxSouth > MAX_REGION_LIMIT
        || maxWest - maxEast > MAX_REGION_LIMIT) {
        console.info('trying to get too many regions')
        return []
    }

    const regions: Region[] = []
    for (let lat = maxNorth+1; lat > maxSouth-1; lat--) {
        for (let lon = maxWest-1; lon < maxEast+1; lon++) {
            regions.push(new Region({
                north: lat,
                south: lat - 1,
                west: lon,
                east: lon + 1,
            }))
        }
    }

    return regions
}

type IncidentsInRegion = {
    region: Region
    incidents: Incident[]
}

export function useIncidents(bounds?: LatLngBounds): Incident[] {
    const client = useClient(ViewerService)
    const [ incidentRegions, setIncidentRegions] = React.useState<IncidentsInRegion[]>([])
    const [ incidents, setIncidents ] = React.useState<Incident[]>([])

    React.useEffect(() => {
        if (!bounds) return
        setIncidents([])

        const regions = regionsInBounds(bounds)
        for (const region of regions) {
            const existing = incidentRegions.find(inc => inc.region.equals(region))
            if (existing) {
                console.debug(`reusing region ${region.toJsonString()}`)
                setIncidents(prev => [...prev, ...existing.incidents])
                continue
            }

            console.log(`getting region ${region.toJsonString()}`)
            client.viewInRegion({region})
                .then(resp => {
                    setIncidentRegions(prev => [...prev, {
                        region,
                        incidents: resp.incidents,
                    }])
                    setIncidents(prev => [...prev, ...resp.incidents])
                })
        }
    }, [client, incidentRegions, bounds])
    
    return incidents
} 
