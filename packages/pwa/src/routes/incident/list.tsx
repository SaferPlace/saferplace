
import { Incident } from '@saferplace/api/incident/v1/incident_pb'

import { useLoaderData } from "react-router-dom"
import IncidentList from '../../components/incidentlist'
import { Container, Typography } from '@mui/material'

export type IncidentsProps = {
    incidents: Incident[]
    center: {lat: number, lon: number},
    radius: number
}

/**
 * Incidents lists the incidents
 */
export default function Incidents() {
    const {incidents, center, radius} = useLoaderData() as IncidentsProps

    return (
        <Container>
            <IncidentList incidents={incidents} />
            <Typography>
                {center.lat}, { center.lon } { radius }m
            </Typography>
        </Container>
    )
}

