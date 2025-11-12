import { BadRequestException, Injectable, Logger, NotFoundException } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { v4 as uuidv4 } from 'uuid';
import { firstValueFrom } from 'rxjs';
import { QueueService } from 'src/queue/queue.service';
import { RedisService } from 'src/redis/redis.service';
import { CreateNotificationDto, NotificationType } from './dto/create-notification.dto';
import { UpdateStatusDto, NotificationStatus } from './dto/update-status.dto';
import { QueueMessage } from './interfaces/queue-message.interface';

@Injectable()
export class NotificationsService {
    private readonly logger = new Logger(NotificationsService.name);
    private readonly userServiceUrl: string;

    constructor(
        private readonly queueService: QueueService,
        private readonly redisService: RedisService,
        private readonly httpService: HttpService,
    ) {
        this.userServiceUrl = process.env.USER_SERVICE_URL || 'http://localhost:3001';
    }

    async createNotification(dto: CreateNotificationDto) {
        const notificationId = uuidv4();

        this.logger.log(`Processing notification request: ${notificationId}`);

        // Fetch user from User Service
        let recipient: string;
        let userName: string;
        
        try {
            const userUrl = `${this.userServiceUrl}/v1/users/${dto.user_id}`;
            this.logger.log(`Fetching user from: ${userUrl}`);
            
            const userResponse = await firstValueFrom(
                this.httpService.get(userUrl)
            );
            
            const user = userResponse.data.data;
            
            // Check user preferences
            if (dto.notification_type === NotificationType.EMAIL && !user.preferences?.email) {
                throw new BadRequestException('User has disabled email notifications');
            }
            
            if (dto.notification_type === NotificationType.PUSH && !user.preferences?.push) {
                throw new BadRequestException('User has disabled push notifications');
            }
            
            // Get recipient based on notification type
            if (dto.notification_type === NotificationType.EMAIL) {
                recipient = user.email;
                if (!recipient) {
                    throw new BadRequestException('User email not found');
                }
            } else {
                recipient = user.push_token;
                if (!recipient) {
                    throw new BadRequestException('User push token not configured');
                }
            }
            
            userName = user.full_name || user.name || 'User';
            
            // Add user name to variables if not provided
            if (!dto.variables.name) {
                dto.variables.name = userName;
            }
            
            this.logger.log(`User fetched successfully: ${userName} (${recipient})`);
            
        } catch (error) {
            if (error instanceof BadRequestException) {
                throw error;
            }
            this.logger.error(`Failed to fetch user: ${error.message}`);
            throw new NotFoundException(`User with ID ${dto.user_id} not found`);
        }

        // Create queue message without pre-rendered content
        // Email Service will fetch template from Template Service
        const queueMessage: QueueMessage = {
            notification_id: notificationId,
            notification_type: dto.notification_type,
            user_id: dto.user_id,
            recipient,
            subject: undefined,  // Let Email Service fetch from Template Service
            title: undefined,
            body: undefined,     // Let Email Service render with variables
            template_code: dto.template_code,
            variables: dto.variables,
            priority: dto.priority || 1,
            metadata: {
                timestamp: new Date().toISOString(),
                retry_count: 0,
            },
        };


        // Publish to appropriate queue
        if (dto.notification_type === 'email') {
            await this.queueService.publishToEmailQueue(queueMessage);
        } else {
            await this.queueService.publishToPushQueue(queueMessage);
        }

        // Store status in Redis
        await this.redisService.setNotificationStatus(notificationId, {
            status: NotificationStatus.PENDING,
            notification_type: dto.notification_type,
            user_id: dto.user_id,
            recipient,
            request_id: dto.request_id,
            created_at: new Date().toISOString(),
        });

        return {
            success: true,
            message: 'Notification queued successfully',
            data: {
                notification_id: notificationId,
                status: 'pending',
                request_id: dto.request_id,
            },
        };
    }

    async updateNotificationStatus(type: string, dto: UpdateStatusDto) {
        this.logger.log(`Updating status for ${dto.notification_id}: ${dto.status}`);

        // Get existing notification data
        const existingData = await this.redisService.getNotificationStatus(dto.notification_id);

        if (!existingData) {
            throw new NotFoundException(`Notification ${dto.notification_id} not found`);
        }

        // Update status
        const updatedData = {
            ...existingData,
            status: dto.status,
            updated_at: dto.timestamp || new Date().toISOString(),
            error: dto.error,
        };

        await this.redisService.setNotificationStatus(dto.notification_id, updatedData);

        return {
            success: true,
            message: 'Status updated successfully',
            data: {
                notification_id: dto.notification_id,
                status: dto.status,
            },
        };
    }

    async getNotificationStatus(notificationId: string) {
        const status = await this.redisService.getNotificationStatus(notificationId);

        if (!status) {
            throw new NotFoundException(`Notification ${notificationId} not found`);
        }

        return {
            success: true,
            data: {
                notification_id: notificationId,
                ...status,
            },
        };
    }

}
