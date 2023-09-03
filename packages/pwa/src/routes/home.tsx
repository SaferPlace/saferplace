import { Box, Stack } from "@mui/material"
import AlertingIncidents from "../components/alertingincidents"
import ActionStack from "../components/actionstack"

export default function Home() {
    return (
        <Stack spacing={2}>
            <Box sx={(theme) => ({
                zIndex: 10000,
                position: 'absolute',
                bottom: theme.spacing(2),
                right: theme.spacing(2)
            })}>
                <ActionStack hideListButton direction='row' />
            </Box>
            <AlertingIncidents />
        </Stack>
    )
}
