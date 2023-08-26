export type TranslationFile = {
    common: Partial<{
        email: string
        backend: string
        cdn: string
        description: string
        submittedAtTime: string
        reportStatus: string
    }>
    action: Partial<{
        useEmail: string
        useBackend: string
        useCDN: string
        viewIncidents: string
        viewIncident: string
        submitReport: string
        retakePhoto: string
        takePhoto: string
    }>,
    phrases: Partial<{
        addToHomeScreen: string
        reportSuccessfullySubmitted: string
        beforeYouReport: string
        contactAuthoritiesFirst: string
        usingReportLocation: string
        incidentDescriptionPlaceholder: string
        alertsNearby: string
        week: string
        day: string
        hour: string
    }>,
    /** resolution are not partial as we need the description for every one */
    resolution: {
        inReview: string
        accepted: string
        alerted: string
        rejected: string
    }
}
