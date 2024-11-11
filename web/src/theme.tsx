import { createTheme, Paper, virtualColor } from "@mantine/core";
export const theme = createTheme({
  primaryColor: "primary",

  colors: {
    plex: [
      "#FFE5B2",
      "#FFD280",
      "#FFBF4D",
      "#FFAC1A",
      "#EBAF00",
      "#DA9F00",
      "#C58F00",
      "#B17F00",
      "#9C6F00",
      "#885F00",
    ],
    mal: [
      "#D6E4FF",
      "#ADC8FF",
      "#84ABFF",
      "#5B8FFF",
      "#2E51A2",
      "#284792",
      "#233D82",
      "#1D3371",
      "#172961",
      "#121F51",
    ],

    primary: virtualColor({
      name: "primary",
      dark: "plex",
      light: "mal",
    }),
  },

  fontFamily: "Verdana, sans-serif",

  components: {
    Paper: Paper.extend({
      defaultProps: {
        mt: "sm",
        shadow: "md",
        radius: "xl",
        withBorder: true,
        p: "xl",
      },
      styles: {
        root: {
          minWidth: 450,
          minHeight: 250,
          margin: "0 auto",
          width: "100%",
        },
      },
    }),
  },
});
