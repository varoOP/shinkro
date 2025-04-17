import {showNotification} from "@mantine/notifications";

type NotificationType = "success" | "error" | "info";

interface NotificationProps {
    type: NotificationType;
    title: string;
    message?: string;
}

export const displayNotification = ({
                                        type,
                                        title,
                                        message = "",
                                    }: NotificationProps) => {
    let color;

    switch (type) {
        case "success":
            color = "green";
            break;
        case "error":
            color = "red";
            break;
        case "info":
            color = "blue";
            break;
        default:
            color = "blue";
    }

    showNotification({
        title,
        message,
        color,
        withCloseButton: true,
        autoClose: 3000,
    });
};
