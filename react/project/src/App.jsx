import React from 'react';
import Footer from '/src/components/Footer';
import Header from '/src/components/Header';
import RoutesConfig from '/src/routes/Routes';
import '/src/i18n';
import { CustomThemeProvider } from '/src/context/ThemeContext';
import { UserProvider } from '/src/context/UserContext';
import { Box } from '@mui/material';

const App = () => {
  return (
    <UserProvider>
      <CustomThemeProvider>
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
          <Header />
          <Box component="main" sx={{ flexGrow: 1 }}>
            <RoutesConfig />
          </Box>
          <Footer />
        </Box>
      </CustomThemeProvider>
    </UserProvider>
  );
};

export default App;