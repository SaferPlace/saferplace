import { createConnectTransport, createPromiseClient } from '@bufbuild/connect-web'
import { CssBaseline } from '@mui/material'
import React from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, LoaderFunctionArgs, RouterProvider } from 'react-router-dom'
import Incident, { Props } from './routes/incident'
import Pending from './routes/pending'
import Root from './routes/root'
import { ReviewService } from '@saferplace/api/review/v1/review_connectweb'
import { BasicIncidentDetails, ReviewIncidentRequest } from '@saferplace/api/review/v1/review_pb'
import { AuthProvider } from 'oidc-react'
import ErrorPage from './routes/error'

const client = createPromiseClient(
  ReviewService,
  createConnectTransport({
    baseUrl: import.meta.env.VITE_BACKEND,
  }),
)

async function pendingLoader(): Promise<BasicIncidentDetails[]> {
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
        element: <Incident />,
        loader: incidentLoader,
      }
    ],
    errorElement: <ErrorPage />,
  },
])


function App() {
  return (
    <AuthProvider
      authority={import.meta.env.VITE_OIDC_AUTHORITY}
      clientId={import.meta.env.VITE_OIDC_CLIENT_ID}
      redirectUri={import.meta.env.VITE_OIDC_REDIRECT_URL}
      scope='user:email'
      clientSecret={import.meta.env.VITE_OIDC_CLIENT_SECRET}
    >
      <RouterProvider router={router} />
      <CssBaseline />
    </AuthProvider>
  )
}


ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
