import React from "react";
import Drawer from "@components/Drawer";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import { useTheme, DarkThemeProvider } from "@components/ThemeContext";

const App = () => {
  return (
    <DarkThemeProvider>
      <ThemedApp />
    </DarkThemeProvider>
  );
};

const ThemedApp = () => {
  const { darkMode } = useTheme();

  const themeOptions = React.useMemo(
    () =>
      createTheme({
        palette: {
          mode: darkMode ? "dark" : "light",
          primary: {
            main: "#2E51A2",
          },
          secondary: {
            main: "#EBAF00",
          },
        },
      }),
    [darkMode]
  );

  return (
    <ThemeProvider theme={themeOptions}>
      <div>
        <Drawer />
      </div>
    </ThemeProvider>
  );
};

export default App;
