import React from 'react'
import {
  AppBar,
  Box,
  Container,
  TextField,
  Toolbar,
  Typography,
} from '@mui/material'
import { Outlet } from 'react-router-dom'

export default function Root() {
  const [ email, setEmail ] = React.useState<string>(localStorage.getItem('email') ?? '')

  const updateEmail = (value: string) => {
    localStorage.setItem('email', value)
    setEmail(value)
  }

  return (
    <Box>
      <AppBar position='fixed'>
        <Toolbar sx={{ display: 'flex', justifyContent: 'space-between'}}>
          <Typography variant='h4'>SaferPlace Review</Typography>
          <Box>
          <TextField
            variant='outlined'
            value={email}
            onChange={e => updateEmail(e.target.value)}
            placeholder='Reviewer Email'
          />
          </Box>
        </Toolbar>
      </AppBar>
      <Container>
        <Toolbar />
        <Toolbar />
        <Outlet />
      </Container>
    </Box>
  )
}
