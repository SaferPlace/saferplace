import { Add } from "@mui/icons-material";
import { Fab, List, Stack } from "@mui/material";
import { Coordinates } from "@saferplace/api/incident/v1/incident_pb";
import { useTranslation } from "react-i18next";
import { NavLink } from "react-router-dom";
import { ResponsiveStyleValue } from '@mui/system'

type Params = {
    direction?: ResponsiveStyleValue<"row" | "row-reverse" | "column" | "column-reverse">
    radius?: number
    center?: Coordinates
    showList?: boolean
}

export default function ActionStack({direction, radius, center}: Params) {
    const { t } = useTranslation()
    return (
        <Stack
            direction={direction}
            spacing={2}
        >
            <Fab
                variant='extended'
                component={NavLink}
                to={`/incidents?radius=${radius}&lat=${center?.lat ?? 0}&lon=${center?.lon ?? 0}`}
            >
                <List sx={{ marginInlineEnd: 1 }} />
                {t('action:viewIncidents')}
            </Fab>
            <Fab
                component={NavLink}
                variant='extended'
                color='primary'
                to='/report'
            >
                <Add sx={{ marginInlineEnd: 1 }} />
                {t('action:submitReport')}
            </Fab>
        </Stack>
)
}