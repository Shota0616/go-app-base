import React from 'react';
import Footer from '/src/components/Footer';
import Header from '/src/components/Header';
import RoutesConfig from '/src/routes/Routes';
import '/src/i18n';
import { CustomThemeProvider } from '/src/context/ThemeContext';
import { UserProvider } from '/src/context/UserContext';

const App = () => {
  return (
    <UserProvider>
      <CustomThemeProvider>
        <Header />
        {/* RoutesConfigを呼び出し（route設定） */}
        <RoutesConfig />
        <Footer />
      </CustomThemeProvider>
    </UserProvider>
  );
};

export default App;