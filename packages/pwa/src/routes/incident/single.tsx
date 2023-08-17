import { Incident } from "@saferplace/api/incident/v1/incident_pb"
import { Card } from "@mui/material"
import { useLoaderData } from "react-router-dom"

export default function IncidentDetails() {
    const incident = useLoaderData() as Incident

    return (
        <Card>
            {incident.description}
            {incident.timestamp?.toDate.toString()}
        </Card>
    )
}