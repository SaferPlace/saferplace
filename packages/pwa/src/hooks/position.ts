import { PartialMessage } from "@bufbuild/protobuf"
import { Coordinates } from "@saferplace/api/incident/v1/incident_pb"
import React from "react"

export function usePosition(): [PartialMessage<Coordinates>?, Error?] {
    const [position, setPosition] = React.useState<PartialMessage<Coordinates> | undefined>()
    const [ error, setError ] = React.useState<Error | undefined>()

    React.useEffect(() => {
        getPosition()
            .then(setPosition)
            .catch(setError)
    }, [])

    return [position, error]
}

export async function getPosition(): Promise<PartialMessage<Coordinates>> {
    return new Promise((resolve, reject) => {
        if (!navigator.geolocation) {
            return reject('No geolocation')
        }

        navigator.geolocation.getCurrentPosition(
            (position) => {
                console.info(position.coords)
                resolve({lat: position.coords.latitude, lon: position.coords.longitude})
            },
            (error) => {
                console.error(error)
                // If we cannot get the location default to The Criminal Court of Justice
                resolve({lat: 53.34868617902951, lon: -6.29567143778413})
            },
        )
    })
}