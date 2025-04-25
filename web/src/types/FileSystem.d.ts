export interface FileSystem {
    name: string,
    path: string,
    is_dir: boolean
}

interface LogFileResponse {
    name: string;
    path: string;
    size: string;
    size_human: string;
    modified_at: string;
}
