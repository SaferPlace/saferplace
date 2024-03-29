import { Coordinates } from "@saferplace/api/incident/v1/incident_pb"
import { MapContainer, Marker, TileLayer, LayerGroup, Popup } from 'react-leaflet'
import { useMapEvents } from 'react-leaflet/hooks'
import { Box, Button, Skeleton } from "@mui/material"
import React from 'react'
import { usePosition } from "../../hooks/position"
import 'leaflet/dist/leaflet.css'
import { useTranslation } from "react-i18next"
import { useNavigate } from "react-router-dom"
import ActionStack from "../../components/actionstack"
import { LatLngBounds } from "leaflet"
import { useIncidents } from "../../hooks/incidents"

export default function Map() {
    const [ initialPosition ] = usePosition()
    const [ center, setCenter ] = React.useState<Coordinates | undefined>()
    const { t } = useTranslation()
    const navigate = useNavigate()
    const [bounds, setBounds] = React.useState<LatLngBounds | undefined>()

    const incidents = useIncidents(bounds)

    // Wait for the initialPosition to become available
    React.useEffect(() => {
        if (!initialPosition) return
        if (center) return

        setCenter(new Coordinates(initialPosition))
    }, [center, initialPosition])
    
    return (
        <Box sx={{
            display: 'flex',
            flexGrow: 1,
            '.leaflet-container': (theme) => ({
                display: 'flex',
                width: '100vw',
                minWidth: '100%',
                maxWidht: '100%',
                height: window.innerHeight - Number(theme.components?.MuiToolbar?.defaultProps?.style?.height || 54), // 
                minHeight: '100%',
                maxHeight: '100%',
                flexGrow: 1,
            }),
        }}>
            <Box
                sx={(theme) => ({
                    zIndex: 10000,
                    position: 'absolute',
                    bottom: theme.spacing(2),
                    right: theme.spacing(2),
                })}
            >
                <ActionStack
                    direction={{ xs: 'column', md: 'row'}}
                    bounds={bounds}
                />
            </Box>
            { center ? (
                <MapContainer center={latlon(center)} zoom={13}>
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
                    </LayerGroup>
                    <MapControl
                        setCenter={setCenter}
                        setBounds={setBounds}
                    />
                </MapContainer>
            ) : (
                <Skeleton width='100%' height='100%' />
            )}
            
        </Box>
        
    )
}

export type MapControlProps = {
    setCenter: (coordinates: Coordinates) => void
    setBounds: (bounds: LatLngBounds) => void
}

function MapControl({setCenter, setBounds }: MapControlProps) {
    const map = useMapEvents({
        moveend: () => {
            const {lat, lng} = map.getCenter()
            // TODO: Add a diff based algorithm to only get the new regions which changed.
            setBounds(map.getBounds())
            setCenter(new Coordinates({lat, lon: lng}))
        },
    })
    return null
}

function latlon(coords: Coordinates | undefined): [number, number] {
    return [coords?.lat || 0, coords?.lon || 0]
}