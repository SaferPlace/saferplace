import { Outlet, useNavigate } from "react-router-dom"
import {useUser} from "../hooks/user"
import React from "react"
import { Container } from "@mui/material"

export default function Root() {
    const [user] = useUser()
    const navigate = useNavigate()

    React.useEffect(() => {
        if (user === '') {
            console.info('user is not authenticated, redirecting to login')
            navigate('/login')
        }
    }, [user, navigate])
   
    return (
        <Container sx={{marginBlock: 2}}>
            <Outlet />
        </Container>
    )
}
