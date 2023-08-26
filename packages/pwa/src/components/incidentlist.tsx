import { Card, CardActionArea, CardContent, CardHeader, CardMedia, Stack } from "@mui/material"
import WarningAmberIcon from '@mui/icons-material/WarningAmber'
import { Incident, Resolution } from "@saferplace/api/incident/v1/incident_pb"
import { useNavigate } from "react-router-dom"

type Props = {
    incidents: Incident[]
}

export default function IncidentList({incidents}: Props) {
    return (
        <Stack spacing={2} margin={2}>
            {incidents.map(incident => (
                <IncidentCard key={incident.id} incident={incident} />
            ))}
        </Stack>
    )
}

type IncidentCardProps = {
    incident: Incident
}

/**
 * IncidentCard is the individual card showing the incident in the list.
 * @param incident Incident to display
 */
export function IncidentCard({incident}: IncidentCardProps) {
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
