import 'pinia';

declare module 'pinia' {
    export interface DefineStoreOptionsBase<S, Store> {
        persist?: boolean | PersistOptions;
    }
}

interface PersistOptions {
    enabled?: boolean;
    strategies?: Array<{
        key?: string;
        storage?: Storage;
        paths?: string[];
    }>;
}
