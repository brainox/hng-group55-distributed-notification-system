import { BadRequestException, Injectable, Logger, NotFoundException } from '@nestjs/common';
import { v4 as uuidv4} from 'uuid';
import { QueueService } from 'src/queue/queue.service';
import { RedisService } from 'src/redis/redis.service';
import { CreateNotificationDto, NotificationType } from './dto/create-notification.dto';
import { UpdateStatusDto, NotificationStatus } from './dto/update-status.dto';
import { QueueMessage } from './interfaces/queue-message.interface';

@Injectable()
export class NotificationsService {
    private readonly logger = new Logger(NotificationsService.name);

    constructor(
        private readonly queueService: QueueService,
        private readonly redisService: RedisService,
    ){}

    async createNotification(dto: CreateNotificationDto) {
        const notificationId = uuidv4();
        // const correlationId = `req-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

        this.logger.log(`Processing notification request: ${notificationId}`);

        // TODO: Call User Service to get user email/push token
        // const user = await this.getUserFromService(dto.user_id);
        // recipient = dto.type === 'email' ? user.email : user.push_token;
        // For now, recipient is mocked
        const recipient = dto.notification_type === NotificationType.EMAIL 
        ? 'user@example.com' 
        : 'fcm-token-mock-123';

        this.logger.warn('Using mock recipient - integrate User Service');


        // TODO: Get template from Template Service
        // const template = await this.getTemplateFromService(dto.template_id);
        // Mock template for now
        const mockTemplate = {
            subject: dto.notification_type === 'email' ? 'Test Notification' : undefined,
            title: dto.notification_type === 'push' ? 'Test Notification' : undefined,
            body: 'Hello {{name}}, this is a test notification!',
        };


        // Render template with variables
        let renderedBody = mockTemplate.body;
        if (dto.variables) {
            Object.keys(dto.variables).forEach(key => {
                renderedBody = renderedBody.replace(`{{${key}}}`, dto.variables[key]);
            });
        }


        // Create queue message
        const queueMessage: QueueMessage = {
            notification_id: notificationId,
            notification_type: dto.notification_type,
            user_id: dto.user_id,
            recipient,
            subject: mockTemplate.subject,
            title: mockTemplate.title,
            body: renderedBody,
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
