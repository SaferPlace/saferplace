import { Box, Button, Container, Paper, Stack, TextField, Toolbar, Typography } from "@mui/material"
import React  from "react"
import { useNavigate } from "react-router-dom"
import { useTranslation } from 'react-i18next'

export default function Login() {
    const [email, setEmail] =  React.useState<string>(localStorage.getItem('email') ?? '')
    const [backend, setBackend] = React.useState<string>(localStorage.getItem('backend') ?? import.meta.env.VITE_BACKEND)
    const [cdn, setCDN] = React.useState<string>(localStorage.getItem('cdn') ?? import.meta.env.VITE_CDN)
    const navigate = useNavigate()
    const { t } = useTranslation()

    React.useEffect(() => {
        const themeColor = document.querySelector('meta[name="theme-color"]')?.getAttribute('content') ?? ''
        const backgroundColor = document.body.style.backgroundColor

        document.querySelector('meta[name="theme-color"]')?.setAttribute('content', 'blue')
        document.body.style.backgroundColor = 'blue'

        return () => {
            document.querySelector('meta[name="theme-color"]')?.setAttribute('content', themeColor)
            document.body.style.backgroundColor = backgroundColor
        }
    }, [])

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
                    <Paper sx={{ padding: 4 }}>
                        <Stack spacing={2} direction='column'>
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
                </Stack>
            </Container>
        </Box>
    )
}
