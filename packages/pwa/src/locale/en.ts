import { TranslationFile } from "."

export default {
    common: {
        email: "Email",
        backend: "Backend",
        description: "Description",
        submittedAtTime: "Submission Time",
        reportStatus: "Report Status",
    },
    action: {
        useEmail: "Use Email",
        useBackend: "Use Backend",
        viewIncidents: "View Incidents",
        submitReport: "Submit Report",
    },
    phrases: {
        addToHomeScreen: "Add SaferPlace to your homescreen for easier access!",
        reportSuccessfullySubmitted: "Successfully Submitted! Your report is now in review.",
        beforeYouReport: "Before you report",
        contactAuthoritiesFirst: "If your life or others is at risk, contact emergency services at 112 before creating a report!",
        usingReportLocation: "Using your location to create the report:",
        incidentDescriptionPlaceholder: "Describe what is happening",
    },
    resolution: {
        inReview: "In Review",
        accepted: "Accepted",
        alerted: "Alert Sent",
        rejected: "Rejected",
    },
} as TranslationFile
