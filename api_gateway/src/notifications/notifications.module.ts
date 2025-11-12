import { Module } from '@nestjs/common';
import { HttpModule } from '@nestjs/axios';
import { NotificationsService } from './notifications.service';
import { NotificationsController } from './notifications.controller';
import { QueueModule } from 'src/queue/queue.module';
import { RedisModule } from 'src/redis/redis.module';

@Module({
  imports: [QueueModule, RedisModule, HttpModule],
  providers: [NotificationsService],
  controllers: [NotificationsController]
})
export class NotificationsModule {}
