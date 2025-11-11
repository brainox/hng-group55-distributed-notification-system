import { Module } from '@nestjs/common';
import { NotificationsService } from './notifications.service';
import { NotificationsController } from './notifications.controller';
import { QueueModule } from 'src/queue/queue.module';
import { RedisModule } from 'src/redis/redis.module';

@Module({
  imports: [QueueModule, RedisModule],
  providers: [NotificationsService],
  controllers: [NotificationsController]
})
export class NotificationsModule {}
