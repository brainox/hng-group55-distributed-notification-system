import { Body, Controller, Get, HttpCode, HttpStatus, Param, Post, UseGuards, Request } from '@nestjs/common';
import { NotificationsService } from './notifications.service';
import { CreateNotificationDto } from './dto/create-notification.dto';
import { UpdateStatusDto } from './dto/update-status.dto';
import { AuthGuard } from '../auth/auth.guard';

@Controller('notifications')
export class NotificationsController {
    constructor(private readonly notificationsService: NotificationsService){}

    @Post()
    // @UseGuards(AuthGuard)
    @HttpCode(HttpStatus.OK)
    async sendNotification(@Body() dto: CreateNotificationDto, @Request() req: any) {
        console.log('Authenticated user:', req.user);
        return this.notificationsService.createNotification(dto, req.user);
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
    // @UseGuards(AuthGuard)
    async getStatus(@Param('id') id: string, @Request() req: any) {
        return this.notificationsService.getNotificationStatus(id, req.user);
    }
}
