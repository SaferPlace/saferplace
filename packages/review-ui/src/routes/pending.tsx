import { Card, CardContent, CardHeader, Stack, Typography } from '@mui/material'
import React from 'react'
import { useLoaderData, Link } from 'react-router-dom'
import { ArrowForwardIos } from '@mui/icons-material'
import { Incident } from '@saferplace/api/incident/v1/incident_pb'

export default function Pending() {
  const incidents = useLoaderData() as Incident[]
  return (
    <Stack spacing={1}>
      <Typography variant='h4'>Incidents Pending Review</Typography>
      {incidents.map(incident => (
        <Card
          key={incident.id}
          component={Link}
          to={`incident/${incident.id}`}
          sx={{textDecoration: 'none'}}
        >
          <CardHeader
            action={<ArrowForwardIos />}
            title={incident.timestamp?.toDate().toString()}
          />
          <CardContent>
            
            <Typography>{incident.description}</Typography>
          </CardContent>
        </Card>
      ))}
    </Stack>
  )
}
