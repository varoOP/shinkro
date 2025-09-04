interface User {
    id: number;
    username: string;
    admin: boolean;
}

interface UserUpdate {
    username_current: string;
    password_current: string;
    password_new: string;
}

interface CreateUserRequest {
    username: string;
    password: string;
    admin: boolean;
}
