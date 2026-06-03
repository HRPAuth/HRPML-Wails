
import { Box, Typography } from '@mui/material';

export default function VersionManager() {
  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
        Versions
      </Typography>
      <Typography sx={{ color: '#9ca3af' }}>
        Manage your game versions here.
      </Typography>
    </Box>
  );
}
