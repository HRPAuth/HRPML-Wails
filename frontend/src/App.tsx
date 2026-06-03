import { useState } from 'react';
import { Box, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText, AppBar, Toolbar, Typography, IconButton } from '@mui/material';
import { Gamepad, Home, Settings, Group, Layers, Menu, BugReport } from '@mui/icons-material';
import HomeView from './components/HomeView';
import VersionManager from './components/VersionManager';
import AccountManager from './components/AccountManager';
import ModManager from './components/ModManager';
import SettingsView from './components/SettingsView';
import DebugView from './components/DebugView';

const drawerWidth = 240;

interface NavItem {
  id: string;
  label: string;
  icon: typeof Home;
}

const navItems: NavItem[] = [
  { id: 'home', label: 'Play', icon: Home },
  { id: 'versions', label: 'Versions', icon: Gamepad },
  { id: 'accounts', label: 'Accounts', icon: Group },
  { id: 'mods', label: 'Mods', icon: Layers },
  { id: 'settings', label: 'Settings', icon: Settings },
  { id: 'debug', label: 'Debug', icon: BugReport },
];

export default function App() {
  const [selectedView, setSelectedView] = useState<string>('home');
  const [selectedVersion, setSelectedVersion] = useState<string | null>(null);
  const [mobileOpen, setMobileOpen] = useState(false);

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const renderView = () => {
    switch (selectedView) {
      case 'home':
        return <HomeView selectedVersion={selectedVersion} onVersionSelect={setSelectedVersion} />;
      case 'versions':
        return <VersionManager />;
      case 'accounts':
        return <AccountManager />;
      case 'mods':
        return <ModManager />;
      case 'settings':
        return <SettingsView />;
      case 'debug':
        return <DebugView />;
      default:
        return <HomeView selectedVersion={selectedVersion} onVersionSelect={setSelectedVersion} />;
    }
  };

  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar
        position="fixed"
        sx={{
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
          backgroundColor: '#282c34',
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: 'none' } }}
          >
            <Menu />
          </IconButton>
          <Typography variant="h6" noWrap component="div">
            Samuel Client
          </Typography>
        </Toolbar>
      </AppBar>

      <Box
        component="nav"
        sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
        aria-label="mailbox folders"
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          <Box onClick={handleDrawerToggle} sx={{ textAlign: 'center' }}>
            <Typography variant="h6" sx={{ my: 2 }}>Samuel Client</Typography>
            <List>
              {navItems.map((item) => (
                <ListItem key={item.id}>
                  <ListItemButton selected={selectedView === item.id}>
                    <ListItemIcon>
                      <item.icon />
                    </ListItemIcon>
                    <ListItemText primary={item.label} />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          </Box>
        </Drawer>

        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth, backgroundColor: '#212529' },
          }}
          open
        >
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="h6" sx={{ color: 'white' }}>Samuel Client</Typography>
          </Box>
          <List>
            {navItems.map((item) => (
              <ListItem key={item.id}>
                <ListItemButton
                  selected={selectedView === item.id}
                  onClick={() => setSelectedView(item.id)}
                  sx={{ '&.Mui-selected': { backgroundColor: '#4CAF50' } }}
                >
                  <ListItemIcon sx={{ color: selectedView === item.id ? 'white' : '#9ca3af' }}>
                    <item.icon />
                  </ListItemIcon>
                  <ListItemText primary={item.label} sx={{ color: selectedView === item.id ? 'white' : '#d1d5db' }} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
        </Drawer>
      </Box>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          backgroundColor: '#1a1d21',
          minHeight: '100vh',
        }}
      >
        <Toolbar />
        {renderView()}
      </Box>
    </Box>
  );
}
