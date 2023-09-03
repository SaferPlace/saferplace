import { Add } from "@mui/icons-material";
import { Fab, List, Stack } from "@mui/material";
import { useTranslation } from "react-i18next";
import { NavLink } from "react-router-dom";
import { ResponsiveStyleValue } from '@mui/system'
import { LatLngBounds } from "leaflet";

type Params = {
    direction?: ResponsiveStyleValue<"row" | "row-reverse" | "column" | "column-reverse">
    bounds?: LatLngBounds
    hideListButton?: boolean
}

export default function ActionStack({direction, bounds, hideListButton}: Params) {
    const { t } = useTranslation()

    const params = new URLSearchParams
    if (bounds) {
        params.set('north', bounds.getNorth().toPrecision(4))
        params.set('south', bounds.getSouth().toPrecision(4))
        params.set('east', bounds.getEast().toPrecision(4))
        params.set('west', bounds.getWest().toPrecision(4))
    }

    return (
        <Stack
            direction={direction}
            spacing={2}
        >
            { !hideListButton && (
                <Fab
                    variant='extended'
                    component={NavLink}
                    to={`/incidents?${params.toString()}`}
                >
                    <List sx={{ marginInlineEnd: 1 }} />
                    {t('action:viewIncidents')}
                </Fab>
            )}
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