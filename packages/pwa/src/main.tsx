import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider
} from 'react-router-dom'
import Root from './routes/root'
import { CssBaseline } from '@mui/material'
import Login from './routes/login'
import './i18n'
import Home from './routes/home'
import IncidentList from './routes/incident/list'
import Incident from './routes/incident/single'
import { incidentLoader, incidentsInRadiusLoader } from './routes/incident/loaders'
import Report from './routes/report'
import { reportLoader } from './routes/loaders'

const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    children: [
      {
        path: '/',
        Component: Home,
      }, {
        path: '/incidents',
        loader: incidentsInRadiusLoader,
        Component: IncidentList,
      }, {
        path: 'incident/:id',
        loader: incidentLoader,
        Component: Incident,
      }, {
        path: 'report',
        loader: reportLoader,
        Component: Report,
      }
    ]
  }, {
    path: '/login',
    Component: Login,
  }
])

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <CssBaseline />
  </React.StrictMode>,
)
