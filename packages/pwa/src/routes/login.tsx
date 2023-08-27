import { Alert, Box, Button, Container, Fade, FormControlLabel, Paper, Stack, Switch, TextField, Toolbar, Typography, useTheme } from "@mui/material"
import React, { ChangeEvent }  from "react"
import { useNavigate } from "react-router-dom"
import { useTranslation } from 'react-i18next'
import { Warning } from "@mui/icons-material"

export default function Login() {
    const [email, setEmail] =  React.useState<string>(localStorage.getItem('email') ?? '')
    const [backend, setBackend] = React.useState<string>(localStorage.getItem('backend') ?? import.meta.env.VITE_BACKEND)
    const [cdn, setCDN] = React.useState<string>(localStorage.getItem('cdn') ?? import.meta.env.VITE_CDN)
    const navigate = useNavigate()
    const { t } = useTranslation()
    const theme = useTheme()
    const [ showAdvanced, setShowAdvanced ] = React.useState<boolean>(false)

    React.useEffect(() => {
        const backgroundColor = document.body.style.backgroundColor

        document.body.style.backgroundColor = theme.palette.primary.main

        return () => {
            document.body.style.backgroundColor = backgroundColor
        }
    }, [theme])

    const saveEmail = () => {
        localStorage.setItem('email', email)
        navigate('/')
    }

    const saveBackend = () => {
        localStorage.setItem('backend', backend)
    }

    const saveCDN = () => {
        localStorage.setItem('cdn', cdn)
    }

    const handleChangedAdvanced = (e: ChangeEvent<HTMLInputElement>) => {
        setShowAdvanced(e.target.checked)
    }

    return (
        <Box>
            <Toolbar />
            <Container>
                <Stack justifyContent='center' spacing={2}>
                    <Typography variant='h2' color='white'>SaferPlace</Typography>
                    <Paper sx={{padding: 4}}>
                        <Stack spacing={1}>
                            <TextField
                                label={t('common:email')}
                                variant='outlined'
                                fullWidth
                                type='email'
                                value={email}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                    setEmail(e.target.value)
                                }}
                            />
                            <Button
                                variant='contained'
                                fullWidth
                                onClick={saveEmail}
                            >
                                {t('action:useEmail')}
                            </Button>
                            <Typography>{t('phrases:addToHomeScreen')}</Typography>
                        </Stack>
                    </Paper>
                    <FormControlLabel
                        sx={{ justifyContent: 'center'}}
                        control={
                            <Switch
                                checked={showAdvanced}
                                onChange={handleChangedAdvanced}
                                color='success'
                            />
                        }
                        label={
                            <Typography color='white'>
                                {t('phrases:showAdvancedOptions')}
                            </Typography>
                        }
                    />
                    <Fade in={showAdvanced}>
                        <Paper sx={{ padding: 4 }}>
                            <Stack spacing={2} direction='column'>
                                <Alert variant='standard' color='warning' icon={<Warning />}>
                                    {t('phrases:advancedDevelopmentOnly')}
                                </Alert>
                                <TextField
                                    label={t('common:backend')}
                                    variant='outlined'
                                    fullWidth
                                    type='url'
                                    value={backend}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                        setBackend(e.target.value)
                                    }}
                                />
                                <Button
                                    variant='contained'
                                    fullWidth
                                    onClick={saveBackend}
                                >
                                    {t('action:useBackend')}
                                </Button>
                                <TextField
                                    label={t('common:cdn')}
                                    variant='outlined'
                                    fullWidth
                                    type='url'
                                    value={cdn}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                        setCDN(e.target.value)
                                    }}
                                />
                                <Button
                                    variant='contained'
                                    fullWidth
                                    onClick={saveCDN}
                                >
                                    {t('action:useCDN')}
                                </Button>
                            </Stack>
                        </Paper>
                    </Fade>
                </Stack>
            </Container>
        </Box>
    )
}
