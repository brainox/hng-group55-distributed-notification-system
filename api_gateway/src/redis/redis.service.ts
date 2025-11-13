import { Injectable, Logger, OnModuleDestroy, OnModuleInit } from '@nestjs/common';
import { createClient, RedisClientType } from 'redis';

@Injectable()
export class RedisService implements OnModuleInit, OnModuleDestroy {
    private client: RedisClientType;
    private readonly logger = new Logger(RedisService.name);
    private readonly redisUrl = process.env.REDIS_URL || 'redis://localhost:6379';

    async onModuleInit() {
        await this.connect();
    }

    private async connect() {
        try {
            this.logger.log('Connecting to Redis...');
            this.client = createClient({url: this.redisUrl });

            this.client.on('error', (err) => {
                this.logger.error('Redis error:', err);
            });

            await this.client.connect();
            this.logger.log('Connected to Redis');
        } catch (error) {
            this.logger.error('Failed to connect to Redis:', error);
            throw error;
        }
    }

    async setNotificationStatus(notificationId: string, status: any, ttl: number = 3600) {
        try {
            const key = `notification:${notificationId}`;
            await this.client.setEx(key, ttl, JSON.stringify(status));
            this.logger.log(`Stored status for ${notificationId}`);
        } catch (error) {
            this.logger.error(`Failed to store status for ${notificationId}:`, error);
            throw error;
        }
    }

    async getNotificationStatus(notificationId: string): Promise<any> {
        try {
            const key = `notification:${notificationId}`;
            const data = await this.client.get(key);
            return data ? JSON.parse(data): null;
        } catch (error) {
            this.logger.error(`Failed to get status for ${notificationId}:`, error);
            throw error;
        }
    }

    async cacheData(key: string, data: any, ttl: number = 300): Promise<void> {
        try {
            await this.client.setEx(key, ttl, JSON.stringify(data));
        } catch (error) {
            this.logger.error(`Failed to cache data for ${key}:`, error);
            throw error;
        }
    }

    async getCachedData(key: string): Promise<any> {
        try {
            const data = await this.client.get(key);
            return data ? JSON.parse(data) : null;
        } catch (error) {
            this.logger.error(`Failed to get cached data for ${key}:`, error);
            return null;
        }
    }

    async onModuleDestroy() {
        try {
            await this.client?.disconnect();
            this.logger.log('Redis connection closed');
        } catch (error) {
            this.logger.error('Error closing Redis connection:', error);
        }
    }
}
