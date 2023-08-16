import useClient from "../hooks/client"

import { ViewerService } from '@saferplace/api/viewer/v1/viewer_connect'
import { Incident } from '@saferplace/api/incident/v1/incident_pb'
import { Stack, Typography } from "@mui/material"
import React from "react"

export default function Incidents() {
    const client = useClient(ViewerService)

    const [incidents, setIncidents] = React.useState<Incident[]>([])

    React.useEffect(() => {
        client.viewInRadius({radius: 10000000000000, center: {lat: 0, lon: 0}})
            .then(res => setIncidents(res.incidents))
            .catch(err => console.error(err))
    }, [client])

    return (
        <Stack spacing={1}>
            {incidents.map(incident => (
                <Typography>{incident.description}</Typography>
            ))}
        </Stack>
    )
}