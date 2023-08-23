import { PartialMessage } from "@bufbuild/protobuf";
import { Alert, AlertTitle, Button, Card, Skeleton, Stack, TextField, Typography } from "@mui/material"
import { useTranslation } from "react-i18next";
import { SendReportRequest, SendReportResponse } from '@saferplace/api/report/v1/report_pb'
import { useLoaderData, useNavigate } from "react-router-dom"
import React from "react"
import { usePosition } from "../hooks/position"
import { getEndpoint } from "../hooks/client";
import PhotoCapture from "../components/photocapture";

export type Props = {
    submit: (request: PartialMessage<SendReportRequest>) => Promise<SendReportResponse>
}

export default function Report() {
    const props = useLoaderData() as Props
    const [ coordinates ] = usePosition()
    const { t } = useTranslation()
    const navigate = useNavigate()

    const [ submitted, setSubmitted ] = React.useState<boolean>(false)
    const [ description, setDescription ] = React.useState<string>('')
    const [ error, setError ] = React.useState<Error | null>(null)
    const [ image, setImage ] = React.useState<File|undefined>()

    const uploadImage = async (): Promise<string> => {
        if (!image) { return '' }
        const data = new FormData()
        data.append('image', image)
        return fetch(`${getEndpoint()}/v1/upload`, {
            method: 'POST',
            body: data,
        })
            .then(resp => resp.text())
    }

    const onSubmit = async() => {
        setSubmitted(true)
        setError(null)

        let imageID = ''
        // Only try to upload the image if it has been specified.
        if (image) {
            imageID = await uploadImage()
                .catch(err => setError(err))
                ?? ''
            if (imageID === '') return
        }

        // TODO: This is temporary only as the API doesn't yet support submitting
        //       images.
        console.info('uploaded image for incident', imageID)
        
        props.submit({
            incident: {
                description,
                coordinates,
                imageId: imageID,
            }
        })
            .then(resp => {
                navigate(`/incident/${resp.id}?isNewReport=true`)
            })
            .catch(err => {
                // Allow the user to retry sending
                setSubmitted(false)
                setError(err)
            })
    }

    return (
        <Card>
            <Stack spacing={2} padding={2}>
                <Typography variant='h5'>{t('action:submitReport')}</Typography>
                <Alert severity="warning">
                    <AlertTitle>{t('phrases:beforeYouReport')}</AlertTitle>
                    <Typography>{t('phrases:contactAuthoritiesFirst')}</Typography>
                </Alert>
                
                <Typography>{t('phrases:usingReportLocation')}</Typography>
                <Typography>
                    { coordinates ? (
                        `${coordinates.lat}, ${coordinates.lon}`
                    ) : (
                        <Skeleton />
                    )}
                </Typography>
                
                <TextField
                    rows={4}
                    disabled={submitted}
                    placeholder={t('phrases:incidentDescriptionPlaceholder')}
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                />
                <PhotoCapture image={image} setImage={setImage} />
                <Button
                    variant='contained'
                    color='success'
                    onClick={onSubmit}
                    disabled={submitted}
                >
                    {t('action:submitReport')}
                </Button>
                { error && (
                    <Alert severity='error'>
                        {/* TODO: Error i18n (#94) */}
                        {error?.message}
                    </Alert>
                )}
            </Stack>
        </Card>
    )
}
