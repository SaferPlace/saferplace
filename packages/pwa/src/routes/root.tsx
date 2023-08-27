import { NavLink, Outlet, useNavigate } from "react-router-dom"
import {useUser} from "../hooks/user"
import React from "react"
import { AppBar, Box, IconButton, Toolbar, Typography } from "@mui/material"
import MapIcon from '@mui/icons-material/Map';

export default function Root() {
    const [user] = useUser()
    const navigate = useNavigate()

    React.useEffect(() => {
        if (user === '') {
            console.info('user is not authenticated, redirecting to login')
            navigate('/login')
        }
    }, [navigate, user])
   
    return (
        <Box>
            <AppBar>
                <Toolbar>
                    <Typography
                        variant='h4'
                        flexGrow={1}
                        component={NavLink}
                        to='/'
                        color='inherit'
                        sx={{ textDecoration: 'none' }}
                    >
                        SaferPlace
                    </Typography>
                    <IconButton
                        component={NavLink}
                        to='/map'
                        color='inherit'
                    >
                        <MapIcon />
                    </IconButton>
                </Toolbar>
            </AppBar>
            <Toolbar />
            <Outlet />
        </Box>
        
    )
}
