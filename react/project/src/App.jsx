import React from 'react';
import Footer from '/src/components/Footer';
import Header from '/src/components/Header';
import RoutesConfig from '/src/routes/Routes';

const App = () => {
  return (
    <>
      <Header />
      {/* RoutesConfigを呼び出し（route設定） */}
      <RoutesConfig />
      <Footer />
    </>
  );
};

export default App;