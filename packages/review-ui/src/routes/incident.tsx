import { Box, Button, Card, CardContent, CardHeader, CardMedia, Stack, TextField, ToggleButton, ToggleButtonGroup, Typography } from '@mui/material'
import React from 'react'
import 'leaflet/dist/leaflet.css'
import { MapContainer, Marker, TileLayer } from 'react-leaflet'
import { PartialMessage } from '@bufbuild/protobuf'

import { ReviewIncidentRequest } from '@saferplace/api/review/v1/review_pb'
import * as ipb from '@saferplace/api/incident/v1/incident_pb'
import { Close, Done, Notifications } from '@mui/icons-material'
import { useLoaderData, useNavigate, useRevalidator } from 'react-router-dom'


export type Props = {
  incident: ipb.Incident
  onSubmit: (review: PartialMessage<ReviewIncidentRequest>) => void
}

function Review({incident, onSubmit}: Props) {
  console.debug('incident', incident)
  const navigate = useNavigate();
  const revalidator = useRevalidator();

  const [description, setDescription] = React.useState<string>(incident.description)
  const [comment, setComment] = React.useState<string>('')
  const [resolution, setResolution] = React.useState<ipb.Resolution>(incident.resolution)

  const submit = () => {
    onSubmit({
      id: incident.id,
      comment,
      resolution,
    })
    setComment('')
    if (resolution != ipb.Resolution.UNSPECIFIED) {
      navigate('/')
    } else {
      // Redirect to the same page
      revalidator.revalidate()
    }
  }
  return (
    <Stack spacing={2}>
      <Card>
      <CardMedia>
        <Box sx={{
          '.leaflet-container': {
            height: '50vh',
          },
        }}>
          <MapContainer center={latlon(incident.coordinates)} zoom={14}>
            <TileLayer
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            <Marker position={latlon(incident.coordinates)} />
          </MapContainer>
        </Box>
        <img src={`${import.meta.env.VITE_CDN}/images/${incident.imageId}`} width='100%' />
      </CardMedia>
      <CardHeader>
        <Typography variant='h4'>Review Incident {incident.id}</Typography>
      </CardHeader>
      <CardContent>
        <TextField
          label='Incident ID'
          value={incident.id}
          disabled
          fullWidth
          margin='normal'
        />
        <TextField
          label='Time Submitted'
          value={new Date(Number(incident.timestamp) * 1000)}
          disabled
          fullWidth
          margin='normal'
        />
        <TextField
          label='Description'
          value={description}
          multiline
          minRows={4}
          onChange={e => setDescription(e.target.value)}
          disabled // Disable until we can send back updated description
          fullWidth
          margin='normal'
        />
        <TextField
          label='Leave a Comment'
          multiline
          value={comment}
          minRows={4}
          onChange={e => setComment(e.target.value)}
          fullWidth
          margin='normal'
        />
        <ToggleButtonGroup
          exclusive
          value={resolution}
          onChange={(_, value) => setResolution(value)}
        >
          <ToggleButton value={ipb.Resolution.UNSPECIFIED}>
            Ignore
          </ToggleButton>
          <ToggleButton value={ipb.Resolution.REJECTED} color='error'>
            <Close />
            Reject
          </ToggleButton>
          <ToggleButton value={ipb.Resolution.ACCEPTED} color='success'>
            <Done />
            Accept
          </ToggleButton>
          <ToggleButton value={ipb.Resolution.ALERTED} color='warning'>
            <Notifications />
            Alert
          </ToggleButton>
        </ToggleButtonGroup>
        <Button
          onClick={() => submit()}
        >
          Submit Review
        </Button>
      </CardContent>
      </Card>

      <Typography variant='h3'>Comments</Typography>

      { incident.reviewerComments.length ? (
        <Stack spacing={1}>
        {incident.reviewerComments.map(comment => (
          <Card key={Number(comment.timestamp)}>
            <CardHeader
              title={comment.authorId}
              subheader={(new Date(Number(comment.timestamp) * 1000)).toString()}
            />
            <CardContent>
            <TextField
              label='Comment'
              multiline
              value={comment.message}
              minRows={3}
              onChange={e => setComment(e.target.value)}
              fullWidth
              margin='normal'
              disabled
            />
            </CardContent>
          </Card>
        ))}
        </Stack>
      ) : (
        <Typography>No Comments</Typography>
      )}
      
    </Stack>
  )
}

export default function Incident() {
  const {incident, onSubmit} = useLoaderData() as Props
 
  return <Review incident={incident} onSubmit={onSubmit} />
}

function latlon(coords: ipb.Coordinates | undefined): [number, number] {
  return [coords?.lat || 0, coords?.lon || 0]
}
