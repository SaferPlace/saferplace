import { Alert, AlertTitle } from '@mui/material'
import { useRouteError } from 'react-router-dom'

export default function Error() {
  const error = useRouteError() as Error

  return (
    <Alert>
      <AlertTitle>Something went wrong!</AlertTitle>
      {error.message}
    </Alert>
  )
}
