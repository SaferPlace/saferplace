import { Box, Stack } from "@mui/material"
import AlertingIncidents from "../components/alertingincidents"
import ActionStack from "../components/actionstack"

export default function Home() {
    return (
        <Stack spacing={2}>
            <Box sx={{display: 'flex', justifyContent: 'center'}}><ActionStack direction='row' /></Box>
            <AlertingIncidents />
        </Stack>
    )
}
