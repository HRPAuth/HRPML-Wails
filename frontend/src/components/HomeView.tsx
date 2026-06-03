
import { Box, Typography } from '@mui/material';

interface HomeViewProps {
  selectedVersion: string | null;
  onVersionSelect: (version: string) => void;
}

export default function HomeView(_props: HomeViewProps) {
  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
        Play
      </Typography>
      <Typography sx={{ color: '#9ca3af' }}>
        Welcome to Samuel Client! Select a version to get started.
      </Typography>
    </Box>
  );
}
