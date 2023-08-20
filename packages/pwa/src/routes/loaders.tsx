import { Props as ReportProps } from './report'
import { ReportService } from '@saferplace/api/report/v1/report_connect'
import { getClient } from "../hooks/client"


export async function reportLoader(): Promise<ReportProps> {
    const client = getClient(ReportService)

    return {
        submit: client.sendReport,
    }
}