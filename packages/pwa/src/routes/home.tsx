import { Link, Stack } from "@mui/material"
import { useTranslation } from "react-i18next"
import { Link as RouterLink } from 'react-router-dom'

export default function Home() {
    const { t } = useTranslation()

    return (
        <Stack>
            <Link component={RouterLink} to='/incidents'>{t('action:viewIncidents')}</Link>
            <Link component={RouterLink} to='/report'>{t('action:submitReport')}</Link>
            <Link component={RouterLink} to='/map'>{t('action:viewMap')}</Link>
        </Stack>
    )
}
