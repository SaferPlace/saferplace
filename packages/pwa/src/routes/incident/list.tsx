
import { Incident, Resolution } from '@saferplace/api/incident/v1/incident_pb'
import { Card, CardActionArea, CardContent, CardHeader, CardMedia, Container, Stack, Typography } from "@mui/material"
import WarningAmberIcon from '@mui/icons-material/WarningAmber'
import { useLoaderData, useNavigate } from "react-router-dom"

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
            <Stack spacing={2} margin={2}>
                {incidents.map(incident => (
                    <IncidentCard key={incident.id} incident={incident} />
                ))}
            </Stack>
            <Typography>
                {center.lat}, { center.lon } { radius }m
            </Typography>
        </Container>
    )
}

type IncidentCardProps = {
    incident: Incident
}

/**
 * IncidentCard is the individual card showing the incident in the list.
 * @param incident Incident to display
 */
function IncidentCard({incident}: IncidentCardProps) {
    const navigate = useNavigate()
    const timestamp = incident.timestamp?.toDate()

    return (
        <Card>
            <CardActionArea
                onClick={() => navigate(`/incident/${incident.id}`)}
                sx={{
                    display: 'flex',
                    justifyContent: 'flex-start'
                }}
            >
                <CardMedia
                    component={WarningAmberIcon}
                    color={incident.resolution == Resolution.ALERTED ? 'error' : 'warning'}
                    sx={{
                        marginInlineStart: 4
                    }}
                />
                <CardContent>
                    <CardHeader
                        title={incident.description}
                        subheader={`${timestamp?.toLocaleDateString()} ${timestamp?.toLocaleTimeString()} `}
                    />
                </CardContent>
            </CardActionArea>
        </Card>
    )
}
