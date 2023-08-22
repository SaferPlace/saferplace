import { Container } from "@mui/material";
import { Outlet } from "react-router-dom";

export default function Page() {
    return (
        <Container sx={{marginBlock: 2}}>
            <Outlet />
        </Container>
    )
}
