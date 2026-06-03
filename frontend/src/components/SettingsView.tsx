
import { Box, Typography } from '@mui/material';

export default function SettingsView() {
  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
        Settings
      </Typography>
      <Typography sx={{ color: '#9ca3af' }}>
        Configure your settings here.
      </Typography>
    </Box>
  );
}
