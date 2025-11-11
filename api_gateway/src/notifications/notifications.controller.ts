import { Body, Controller, Get, HttpCode, HttpStatus, Param, Post } from '@nestjs/common';
import { NotificationsService } from './notifications.service';
import { CreateNotificationDto } from './dto/create-notification.dto';
import { UpdateStatusDto } from './dto/update-status.dto';

@Controller('notifications')
export class NotificationsController {
    constructor(private readonly notificationsService: NotificationsService){}

    @Post('send')
    @HttpCode(HttpStatus.OK)
    async sendNotification(@Body() dto: CreateNotificationDto) {
        return this.notificationsService.createNotification(dto);
    }


    @Post(':notification_type/status')
    @HttpCode(HttpStatus.OK)
    async updateStatus(
        @Param('notification_type') type: string,
        @Body() dto: UpdateStatusDto,
    ) {
        return this.notificationsService.updateNotificationStatus(type, dto);
    }

    @Get(':id/status')
    async getStatus(@Param('id') id: string) {
        return this.notificationsService.getNotificationStatus(id);
    }
}
