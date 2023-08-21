import { Coordinates, Incident } from "@saferplace/api/incident/v1/incident_pb"
import { MapContainer, Marker, TileLayer, LayerGroup, Popup } from 'react-leaflet'
import { useMapEvents } from 'react-leaflet/hooks'
import { Box, Button, Skeleton } from "@mui/material"
import React from 'react'
import { usePosition } from "../../hooks/position"
import useClient from "../../hooks/client"
import { ViewerService } from "@saferplace/api/viewer/v1/viewer_connect"
import 'leaflet/dist/leaflet.css'
import { useTranslation } from "react-i18next"
import { useNavigate } from "react-router-dom"

export type Props = {
    center?: Coordinates
    setCenter: (coordinates: Coordinates) => void
    setZoom: (zoom: number) => void
}

export default function Map() {
    const [ initialPosition ] = usePosition()
    const [ incidents, setIncidents ] = React.useState<Incident[]>([])
    const [ center, setCenter ] = React.useState<Coordinates | undefined>()
    const { t } = useTranslation()
    const navigate = useNavigate()
    
    const [ zoom, setZoom ] = React.useState<number>(13)
    const client = useClient(ViewerService)

    React.useEffect(() => {
        if (!initialPosition) {
            return
        }
        setCenter(new Coordinates(initialPosition))
    }, [initialPosition])

    React.useEffect(() => {
        if (!center) {
            return
        }
        console.debug(`viewing incident at ${center.lat}, ${center.lon} at zoom ${zoom} with radius ${zoomToRadius(center.lat, zoom)}m`)
        client.viewInRadius({
            radius: zoomToRadius(center.lat, zoom), // Static until we know how to convert zoom to radius
            center,
        })
            .then(resp => setIncidents(resp.incidents))
            .catch(err => console.error(err))
    }, [center])
    
    return (
        <Box sx={{
            '.leaflet-container': {
                height: '100vh',
                width: '100vw',
            },
        }}>
            { center ? (
                <MapContainer center={latlon(center)} zoom={zoom}>
                <TileLayer
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                <LayerGroup>
                    { incidents.map(incident => (
                        <Marker key={incident.id} position={latlon(incident.coordinates)}>
                            <Popup>
                                <Button
                                    onClick={() => navigate(`/incident/${incident.id}`)}
                                >
                                    {t('action:viewIncident')}
                                </Button>
                            </Popup>
                        </Marker>
                    )) }
                    {/* Enable for debugging */}
                    {/* <Circle center={latlon(center)} radius={zoomToRadius(center?.lat ?? 0, zoom)} /> */}
                </LayerGroup>
                <MapDescendant
                    center={new Coordinates(initialPosition)}
                    setCenter={setCenter}
                    setZoom={setZoom}
                />
            </MapContainer>
            ) : (
                <Skeleton />
            )}
            
        </Box>
        
    )
}

function MapDescendant({setCenter, setZoom }: Props) {
    const map = useMapEvents({
        zoomend: () => {
            setZoom(map.getZoom())
        },
        moveend: () => {
            const {lat, lng} = map.getCenter()
            setCenter(new Coordinates({lat, lon: lng}))
        },
    })
    return null
}

function latlon(coords: Coordinates | undefined): [number, number] {
    return [coords?.lat || 0, coords?.lon || 0]
}

function zoomToRadius(lat: number, zoom: number): number {
    return ((80_000_000*Math.cos(lat * Math.PI / 180)) / Math.pow(2, zoom)) * 2
}
