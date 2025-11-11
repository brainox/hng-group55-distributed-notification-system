export interface QueueMessage {
    notification_id: string;
    notification_type: 'email' | 'push';
    user_id: string;
    recipient: string;
    subject?: string;
    title?: string;
    body: string;
    variables: {
        name: string;
        link: string;
        meta?: Record<string, any>;
    };    
    template_code: string;
    priority: number,
    metadata: {
        timestamp: string;
        retry_count?: number;
    };
}