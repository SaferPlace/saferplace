import { TranslationFile } from "."

export default {
    common: {
        email: "Email",
        backend: "Backend",
        cdn: "Content Delivery Network",
        description: "Description",
        submittedAtTime: "Submission Time",
        reportStatus: "Report Status",
    },
    action: {
        useEmail: "Use Email",
        useBackend: "Use Backend",
        useCDN: "Use CDN",
        viewIncidents: "View Incidents",
        viewIncident: "View Incident",
        submitReport: "Submit Report",
        takePhoto: "Take Photo",
        retakePhoto: "Retake Photo",
    },
    phrases: {
        addToHomeScreen: "Add SaferPlace to your homescreen for easier access!",
        reportSuccessfullySubmitted: "Successfully Submitted! Your report is now in review.",
        beforeYouReport: "Before you report",
        contactAuthoritiesFirst: "If your life or others is at risk, contact emergency services at 112 before creating a report!",
        usingReportLocation: "Using your location to create the report:",
        incidentDescriptionPlaceholder: "Describe what is happening",
        alertsNearby: "Alerts Nearby",
        week: "Week",
        day: "Day",
        hour: "Hour",
        showAdvancedOptions: 'Show Advanced Options',
        advancedDevelopmentOnly: 'These options are created only for development purposes and should typically not be changed unless you know what you are doing.',
    },
    resolution: {
        inReview: "In Review",
        accepted: "Accepted",
        alerted: "Alert Sent",
        rejected: "Rejected",
    },
} as TranslationFile
