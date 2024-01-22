import React, { createContext, useState, useContext, ReactNode } from "react";

interface ThemeContextProps {
  darkMode: boolean;
  toggleDarkMode: () => void;
}

const ThemeContext = createContext<ThemeContextProps>({
  darkMode: false,
  toggleDarkMode: () => {},
});

interface ThemeProviderProps {
    children: ReactNode; 
  }
  
export const DarkThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  const [darkMode, setDarkMode] = useState(false);

  const toggleDarkMode = () => {
    setDarkMode(!darkMode);
  };

  return (
    <ThemeContext.Provider value={{ darkMode, toggleDarkMode }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => useContext(ThemeContext);
