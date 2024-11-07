import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
// import { Buffer } from "buffer";
import { MantineProvider, createTheme, rem} from "@mantine/core";
import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import "@app/pattern.css";
import { App } from "./App";
import { InitializeGlobalContext } from "./utils/Context";

const theme = createTheme({
  // Set primary color and shades for light and dark modes
  primaryColor: "blue", // Orange as the primary color

  // Shadows customization
  shadows: {
    md: "1px 1px 3px rgba(0, 0, 0, .25)",
    xl: "5px 5px 3px rgba(0, 0, 0, .25)",
  },

  // Headings typography
  headings: {
    fontFamily: "Open Sans, sans-serif",
    sizes: {
      h1: { fontSize: rem(36) }, // Customize h1 font size
    },
  },
});

declare global {
  interface Window {
    APP: APP;
  }
}

window.APP = window.APP || {};
// window.Buffer = Buffer;
InitializeGlobalContext();

const root = createRoot(document.getElementById("root")!);

root.render(
  <StrictMode>
    <MantineProvider defaultColorScheme="dark" theme={theme}>
      <App />
    </MantineProvider>
  </StrictMode>
);
