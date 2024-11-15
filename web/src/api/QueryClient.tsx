import { QueryCache, QueryClient } from "@tanstack/react-query";
import { AuthContext } from "@utils/Context";
import { redirect } from "@tanstack/react-router";
import { LoginRoute } from "@app/routes";
import { displayNotification } from "@components/notifications";

const MAX_RETRIES = 6;

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, query) => {
      console.error(`Caught error for query '${query.queryKey}': `, error);

      if (error.message === "Cookie expired or invalid.") {
        AuthContext.reset();
        redirect({
          to: LoginRoute.to,
          search: {
            redirect: location.href, // redirect after login
          },
        });
        return;
      } else {
        displayNotification({
          title: "Error",
          message: error?.message || "An error occurred",
          type: "error",
        });
      }
    },
  }),
  defaultOptions: {
    queries: {
      throwOnError: (error) => {
        return error.message !== "Cookie expired or invalid.";
      },
      retry: (failureCount, error) => {
        if (error.message === "Cookie expired or invalid.") {
          return false;
        }

        console.error(`Retrying query (N=${failureCount}): `, error);
        return failureCount <= MAX_RETRIES;
      },
    },
    mutations: {
      onError: (error) => {
        console.log("mutation error: ", error);

        if (error instanceof Response) {
          return;
        }

        const message =
          typeof error === "object" && "message" in error
            ? (error as Error).message
            : `${error}`;

        displayNotification({
          title: "Mutation Error",
          message: message,
          type: "error",
        });
      },
    },
  },
});
