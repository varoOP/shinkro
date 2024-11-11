import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
// import { Buffer } from "buffer";
import { MantineProvider } from "@mantine/core";
import { theme } from "@app/theme";
import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import "@app/pattern.css";
import { App } from "./App";
import { InitializeGlobalContext } from "./utils/Context";

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
    <MantineProvider theme={theme} defaultColorScheme="dark">
      <App />
    </MantineProvider>
  </StrictMode>
);
