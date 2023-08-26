import React from "react";
import { usePosition } from "../hooks/position";
import useClient from "../hooks/client";
import { ViewerService } from "@saferplace/api/viewer/v1/viewer_connect";
import { Box, ToggleButton, ToggleButtonGroup, Typography } from "@mui/material";
import { useTranslation } from "react-i18next";
import { Timestamp } from "@bufbuild/protobuf";
import IncidentList from "./incidentlist";
import { Incident } from "@saferplace/api/incident/v1/incident_pb";
import { Region } from "@saferplace/api/viewer/v1/viewer_pb";

export default function AlertingIncidents() {
    const { t } = useTranslation()
    // Get the users position but abstract it away to just contain 3 blocks around the position
    const [position]  = usePosition()
    const client = useClient(ViewerService)
    const [incidents, setIncidents] = React.useState<Incident[]>([])

    const [ since, setSince ] = React.useState<string>('day')

    React.useEffect(() => {
        const sinceTimestamp = new Date(Date.now())
        let numberOfHours = 24
        switch (since) {
            case 'week': numberOfHours = 24*7; break
            case 'day': numberOfHours = 24; break
            case 'hour': numberOfHours = 1; break
        }
        sinceTimestamp.setHours(sinceTimestamp.getHours() - numberOfHours)

        const region = new Region({
            north: (Math.round((position?.lat ?? 0) * 100) / 100),
            south: ((Math.round((position?.lat ?? 0) * 100)-1) / 100),
            west: ((Math.round((position?.lon ?? 0) * 100)-1) / 100),
            east: ((Math.round((position?.lon ?? 0) * 100)) / 100),
        })
        console.info(`region: ${region.toJsonString()}, since ${sinceTimestamp.toLocaleDateString()}`)
        client.viewAlerting({
            since: Timestamp.fromDate(sinceTimestamp),
            region,
        })
            .then(resp => setIncidents(resp.incidents))
    }, [client, position, since])

    return (
        <Box>
            <Typography variant='h2'>{t('phrases:alertsNearby')}</Typography>
            <ToggleButtonGroup
                fullWidth
                color="primary"
                value={since}
                exclusive
                onChange={(_, value) => setSince(value)}
                aria-label="Platform"
            >
                <ToggleButton value="week">{t('phrases:week')}</ToggleButton>
                <ToggleButton value="day">{t('phrases:day')}</ToggleButton>
                <ToggleButton value="hour">{t('phrases:hour')}</ToggleButton>
            </ToggleButtonGroup>
            <IncidentList incidents={incidents} />
        </Box>
    )
}