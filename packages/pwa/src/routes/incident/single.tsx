import { Incident, Resolution } from "@saferplace/api/incident/v1/incident_pb"
import { Alert, AlertColor, Box, Card, CardContent, Stack, TextField, Typography } from "@mui/material"
import { useLoaderData } from "react-router-dom"
import { useTranslation } from "react-i18next"

import { AccessTime, Done, Warning, HighlightOff } from '@mui/icons-material'

export type Props = {
    incident: Incident
    isNewReport?: boolean
}

type Status = {
    description: string
    severity: AlertColor
    icon: React.ReactNode
}

export default function IncidentDetails() {
    const { incident, isNewReport } = useLoaderData() as Props
    const { t } = useTranslation() 

    const created = incident.timestamp?.toDate()

    let status: Status
    switch(incident.resolution) {
        case Resolution.UNSPECIFIED:
            status = {
                icon: <AccessTime />,
                severity: 'info',
                description: t('resolution:inReview'),
            }
            break
        case Resolution.ACCEPTED:
            status = {
                icon: <Done />,
                severity: 'success',
                description: t('resolution:accepted')
            }
            break
        case Resolution.ALERTED:
            status = {
                icon: <Warning />,
                severity: 'warning',
                description: t('resolution:alerted')
            }
            break
        case Resolution.REJECTED:
            status = {
                icon: <HighlightOff />,
                severity: 'error',
                description: t('resolution:rejected')
            }
            break
    }

    return (
        <Stack spacing={2}>
            { isNewReport && (
                <Alert
                    severity='success'
                    variant='filled'
                >
                    {t('phrases:reportSuccessfullySubmitted')}
                </Alert>
            )}
            <Card>
                <CardContent>
                    <Stack spacing={2}>
                        <Typography variant='h4'>{t('common:description')}</Typography>
                        <Box padding={2}>
                            <TextField value={incident.description} disabled={true} fullWidth />
                        </Box>
                        <Typography variant='h4'>{t('common:submittedAtTime')}</Typography>
                        <Box padding={2}>
                            {created?.toLocaleTimeString()} {created?.toDateString()}
                        </Box>
                        <Typography variant='h4'>{t('common:reportStatus')}</Typography>
                        <Box padding={2}>
                            <Alert icon={status.icon} severity={status.severity}>{status.description}</Alert>
                        </Box>
                    </Stack>
                </CardContent>
            </Card>
        </Stack>
    )
}
