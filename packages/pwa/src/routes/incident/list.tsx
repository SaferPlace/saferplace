
import { Incident } from '@saferplace/api/incident/v1/incident_pb'

import { useLoaderData } from "react-router-dom"
import IncidentList from '../../components/incidentlist'
import { Container, Typography } from '@mui/material'
import { LatLngBounds } from 'leaflet'

export type IncidentsProps = {
    incidents: Incident[]
    bounds: LatLngBounds
}

/**
 * Incidents lists the incidents
 */
export default function Incidents() {
    const {incidents, bounds} = useLoaderData() as IncidentsProps

    return (
        <Container>
            <IncidentList incidents={incidents} />
            <Typography>
                n: {bounds.getNorth()}, s: {bounds.getSouth()}, w: {bounds.getWest()}, e: {bounds.getEast()}
            </Typography>
        </Container>
    )
}

