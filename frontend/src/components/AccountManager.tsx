
import { Box, Typography } from '@mui/material';

export default function AccountManager() {
  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
        Accounts
      </Typography>
      <Typography sx={{ color: '#9ca3af' }}>
        Manage your accounts here.
      </Typography>
    </Box>
  );
}
