import { CssBaseline } from '@mui/material'
import React from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, LoaderFunctionArgs, RouterProvider } from 'react-router-dom'
import IncidentView, { Props } from './routes/incident'
import Pending from './routes/pending'
import Root from './routes/root'
import { ReviewService } from '@saferplace/api/review/v1/review_connect'
import ErrorPage from './routes/error'

import { createPromiseClient, Interceptor } from '@bufbuild/connect'
import { createConnectTransport } from '@bufbuild/connect-web'
import { Incident } from '@saferplace/api/incident/v1/incident_pb'

const addReviewerEmailInterceptor: Interceptor = (next) => async (req) => {
  req.header.set('email', localStorage.getItem('email') ?? '')
  return await next(req)
}

const client = createPromiseClient(
  ReviewService,
  createConnectTransport({
    baseUrl: import.meta.env.VITE_BACKEND,
    interceptors: [addReviewerEmailInterceptor],
  }),
)

async function pendingLoader(): Promise<Incident[]> {
  const res = await client.incidentsWithoutReview({})
  return res.incidents
}

// Sending back the action to review incident is probably not the best choice
// but the react router seems to be focused on just the HTTP Form requests.
async function incidentLoader({params}: LoaderFunctionArgs): Promise<Props> {
  const res = await client.viewIncident({id: params.id})
  if (!res.incident) throw new Error("not found")
  return { incident: res.incident, onSubmit: client.reviewIncident }
}

const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    children: [
      {
        path: '/',
        element: <Pending />,
        loader: pendingLoader,
      }, {
        path: 'incident/:id',
        element: <IncidentView />,
        loader: incidentLoader,
      }
    ],
    errorElement: <ErrorPage />,
  },
], {
  basename: import.meta.env.BASE_URL,
})


ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <CssBaseline />
  </React.StrictMode>
)
