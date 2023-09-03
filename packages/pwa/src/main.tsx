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
import { incidentLoader, incidentsInRegionLoader } from './routes/incident/loaders'
import Report from './routes/report'
import { reportLoader } from './routes/loaders'
import Map from './routes/incident/map'
import Page from './routes/page'

const router = createBrowserRouter([
  // Routes not requiring Authentication
  {
    path: '/login',
    Component: Login,
  },
  // Routes Requiring Authentication
  {
    path: '/',
    Component: Root,
    children: [
      // Pages with nice borders etc
      {
        path: '/',
        Component: Page,
        children: [
          {
            path: '/',
            Component: Home,
          }, {
            path: '/incidents',
            loader: incidentsInRegionLoader,
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
      },
      // Fullpage applications
      {
        path: '/map',
        Component: Map,
      }
    ]
  }
])

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <CssBaseline />
  </React.StrictMode>,
)
