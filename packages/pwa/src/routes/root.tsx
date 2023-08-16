import { Outlet, useNavigate } from "react-router-dom"
import {useUser} from "../hooks/user"
import React from "react"


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
        <Outlet />
    )
}