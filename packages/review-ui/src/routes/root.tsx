import React from 'react'
import {
  AppBar,
  Box,
  Container,
  Toolbar,
  Typography,
} from '@mui/material'
import { Outlet } from 'react-router-dom'

export default function Root() {
  return (
    <Box>
      <AppBar position='fixed'>
        <Toolbar>
          <Typography variant='h4'>SaferPlace Review</Typography>
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
