import { Box, Button, Container, Paper, Stack, TextField, Toolbar, Typography } from "@mui/material"
import React  from "react"
import { useNavigate } from "react-router-dom"
import { useTranslation } from 'react-i18next'

export default function Login() {
    const [email, setEmail] =  React.useState<string>(localStorage.getItem('email') ?? '')
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

    const useEmail = () => {
        localStorage.setItem('email', email)
        navigate('/')
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
                                fullWidth type='email'
                                value={email}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                    setEmail(e.target.value)
                                }}
                            />
                            <Button
                                variant='contained'
                                fullWidth
                                onClick={useEmail}
                            >
                                {t('action:useEmail')}
                            </Button>
                            <Typography>{t('phrases:addToHomeScreen')}</Typography>
                        </Stack>
                    </Paper>
                </Stack>
            </Container>
        </Box>
    )
}