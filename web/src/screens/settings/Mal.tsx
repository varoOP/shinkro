// import {useQueryClient, useQuery, useMutation} from "@tanstack/react-query";
// import {useEffect, useState} from "react";
// import {MalQueryOptions} from "@api/queries";
// import {SettingsKeys} from "@api/query_keys";
// import {APIClient} from "@api/APIClient.ts";
// import {displayNotification} from "@components/notifications";

export const Mal = () => {
    // const queryClient = useQueryClient();
    // const [isReachable, setIsReachable] = useState<boolean | null>(null);
    //
    // const {data: mal, isLoading} = useQuery(MalQueryOptions());
    //
    // const isEmptyMal = !mal || Object.keys(mal).length === 0;
    //
    // useEffect(() => {
    //     if (!isEmptyMal) {
    //         APIClient.malauth
    //             .test()
    //             .then(() => setIsReachable(true))
    //             .catch(() => setIsReachable(false));
    //     } else {
    //         setIsReachable(null);
    //     }
    // }, [mal, isEmptyMal]);
    //
    // const updateMutation = useMutation({
    //     mutationFn: APIClient.malauth.storeOpts,
    //     onSuccess: () => {
    //         queryClient.invalidateQueries({queryKey: SettingsKeys.mal()});
    //         displayNotification({
    //             title: "Success",
    //             message: "MAL settings updated successfully.",
    //             type: "success",
    //         });
    //     },
    //     onError: (error) => {
    //         displayNotification({
    //             title: "Update failed",
    //             message: error.message || "Could not update MAL settings",
    //             type: "error",
    //         });
    //     },
    // })

    return (
        <></>
    )
}