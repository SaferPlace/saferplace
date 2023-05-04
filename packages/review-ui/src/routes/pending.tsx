import { Card, CardContent, CardHeader, Stack, Typography } from '@mui/material'
import React from 'react'
import { useLoaderData, Link } from 'react-router-dom'
import { BasicIncidentDetails } from '@saferplace/api/review/v1/review_pb'
import { ArrowForwardIos } from '@mui/icons-material'

export default function Pending() {
  const incidents = useLoaderData() as BasicIncidentDetails[]
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
            title={new Date(Number(incident.timestamp) * 1000).toString()}
            subheader={(new Date(Number(incident.timestamp) * 1000)).toString()}
          />
          <CardContent>
            
            <Typography>{incident.description}</Typography>
          </CardContent>
        </Card>
      ))}
    </Stack>
  )
}
