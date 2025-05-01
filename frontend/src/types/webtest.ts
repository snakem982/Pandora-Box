export type WebTest = {
    id: string;       // 对应 `json:"id"`
    order: number;    // 对应 `json:"order"`
    title: string;    // 对应 `json:"title"`
    src: string;      // 对应 `json:"src"`
    testUrl: string;  // 对应 `json:"testUrl"`
    delay: number;    // 对应 `json:"delay"`
};
