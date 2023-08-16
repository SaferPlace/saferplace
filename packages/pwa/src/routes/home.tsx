import { Link } from "@mui/material"
import { useTranslation } from "react-i18next"
import { Link as RouterLink } from 'react-router-dom'

export default function Home() {
    const { t } = useTranslation()

    return (
        <Link component={RouterLink} to='/incidents'>{t('action:viewIncidents')}</Link>
    )
}