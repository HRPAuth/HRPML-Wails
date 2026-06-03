
import { Box, Typography } from '@mui/material';

export default function ModManager() {
  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
        Mods
      </Typography>
      <Typography sx={{ color: '#9ca3af' }}>
        Manage your mods here.
      </Typography>
    </Box>
  );
}
