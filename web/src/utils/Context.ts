import type { StateWithValue } from "react-ridge-state";
import { newRidgeState } from "react-ridge-state";
import { useMantineColorScheme } from "@mantine/core";

interface SettingsType {
  debug: boolean;
  scrollOnNewLog: boolean;
  indentLogLines: boolean;
  hideWrappedText: boolean;
}

interface AuthInfo {
  username: string;
  isLoggedIn: boolean;
  admin: boolean;
}

// Default values
const AuthContextDefaults: AuthInfo = {
  username: "",
  isLoggedIn: false,
  admin: false,
};

const SettingsContextDefaults: SettingsType = {
  debug: false,
  scrollOnNewLog: false,
  indentLogLines: false,
  hideWrappedText: false,
};

// eslint-disable-next-line
function ContextMerger<T extends {}>(
  key: string,
  defaults: T,
  ctxState: StateWithValue<T>
) {
  let values = structuredClone(defaults);

  const storage = localStorage.getItem(key);
  if (storage) {
    try {
      const json = JSON.parse(storage);
      if (json === null) {
        console.warn(
          `JSON localStorage value for '${key}' context state is null`
        );
      } else {
        values = { ...values, ...json };
      }
    } catch (e) {
      console.error(`Failed to merge ${key} context state: ${e}`);
    }
  }

  ctxState.set(values);
}

const AuthKey = "shinkro_user_auth";
const SettingsKey = "shinkro_settings";

export const InitializeGlobalContext = () => {
  ContextMerger<AuthInfo>(AuthKey, AuthContextDefaults, AuthContext);
  ContextMerger<SettingsType>(
    SettingsKey,
    SettingsContextDefaults,
    SettingsContext
  );
};

function DefaultSetter<T>(name: string, newState: T, prevState: T) {
  try {
    localStorage.setItem(name, JSON.stringify(newState));
  } catch (e) {
    console.error(
      `An error occurred while trying to modify '${name}' context state: ${e}`
    );
    console.warn(`  --> prevState: ${prevState}`);
    console.warn(`  --> newState: ${newState}`);
  }
}

export const AuthContext = newRidgeState<AuthInfo>(AuthContextDefaults, {
  onSet: (newState, prevState) => DefaultSetter(AuthKey, newState, prevState),
});

export const SettingsContext = newRidgeState<SettingsType>(
  SettingsContextDefaults,
  {
    onSet: (newState, prevState) => {
      DefaultSetter(SettingsKey, newState, prevState);
    },
  }
);

export const useThemeToggle = () => {
  const { colorScheme, setColorScheme, clearColorScheme } =
    useMantineColorScheme();

  const toggleTheme = () => {
    setColorScheme(colorScheme === "dark" ? "light" : "dark");
  };

  return { colorScheme, setColorScheme, toggleTheme, clearColorScheme };
};
/**
 * Updates the meta theme color based on the current theme state.
 * Used by Safari to color the compact tab bar on both iOS and MacOS.
 */
// const updateMetaThemeColor = (darkTheme: boolean) => {
//   const color = darkTheme ? '#121315' : '#f4f4f5';
//   let metaThemeColor: HTMLMetaElement | null = document.querySelector('meta[name="theme-color"]');
//   if (!metaThemeColor) {
//     metaThemeColor = document.createElement('meta') as HTMLMetaElement;
//     metaThemeColor.name = "theme-color";
//     document.head.appendChild(metaThemeColor);
//   }

//   metaThemeColor.content = color;
// };
