interface APIKey {
    name: string;
    key: string;
    scopes: string[];
    created_at: Date;
}

interface UserUpdate {
    username_current: string;
    username_new?: string;
    password_current?: string;
    password_new?: string;
}
