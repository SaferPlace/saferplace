import { PartialMessage } from "@bufbuild/protobuf";
import { Alert, AlertTitle, Button, Card, Skeleton, Stack, TextField, Typography } from "@mui/material"
import { useTranslation } from "react-i18next";
import { SendReportRequest, SendReportResponse } from '@saferplace/api/report/v1/report_pb'
import { useLoaderData, useNavigate } from "react-router-dom"
import React from "react"
import { usePosition } from "../hooks/position"

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

    const onSubmit = async() => {
        setSubmitted(true)
        setError(null)
        props.submit({
            incident: {
                description,
                coordinates,
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
                <Button
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
