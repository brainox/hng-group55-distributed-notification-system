import { Injectable, OnModuleInit, OnModuleDestroy, Logger } from '@nestjs/common';
import amqp, { Channel, Connection } from 'amqplib';
import { QueueMessage } from '../notifications/interfaces/queue-message.interface';

@Injectable()
export class QueueService implements OnModuleInit, OnModuleDestroy {
  private connection;
  private channel;
  private readonly logger = new Logger(QueueService.name);
  private readonly rabbitmqUrl = process.env.RABBITMQ_URL || 'amqp://localhost:5672';
  private readonly exchange = 'notifications.direct';

  async onModuleInit() {
    await this.connect();
  }

  private async connect() {
    try {
      this.logger.log('Connecting to RabbitMQ...');
      this.connection = await amqp.connect(this.rabbitmqUrl);
      this.channel = await this.connection.createChannel();
      
      // Ensure exchange exists
      await this.channel.assertExchange(this.exchange, 'direct', { durable: true });

      // Ensure queues exist and are bound
    //   await this.channel.assertQueue('email', { durable: true });
    //   await this.channel.assertQueue('push', { durable: true });
    //   await this.channel.bindQueue('email', this.exchange, 'email');
    //   await this.channel.bindQueue('push', this.exchange, 'push');
      this.logger.log('Connected to RabbitMQ');

    //   this.connection.on('close', () => {
    //     this.logger.warn('RabbitMQ connection closed. Reconnecting...');
    //     setTimeout(() => this.connect(), 5000);
    //   });
    } catch (error) {
      this.logger.error('Failed to connect to RabbitMQ:', error);
      throw error;
    }
  }

  async publishToEmailQueue(message: QueueMessage): Promise<void> {
    return this.publish('email', message);
  }

  async publishToPushQueue(message: QueueMessage): Promise<void> {
    return this.publish('push', message);
  }

  private async publish(routingKey: string, message: QueueMessage): Promise<void> {
    try {
    //   if (!this.channel) await this.connect();
      const messageBuffer = Buffer.from(JSON.stringify(message));
      
      const published = this.channel.publish(
        this.exchange,
        routingKey,
        messageBuffer,
        {
          persistent: true,
          contentType: 'application/json',
          timestamp: Date.now(),
        }
      );

      if (published) {
        this.logger.log(`Published message to ${routingKey}: ${message.notification_id}`);
      } else {
        throw new Error('Failed to publish message');
      }
    } catch (error) {
      this.logger.error(`Failed to publish to ${routingKey}:`, error);
      throw error;
    }
  }

  async onModuleDestroy() {
    try {
      await this.channel?.close();
      await this.connection?.close();
      this.logger.log('RabbitMQ connection closed');
    } catch (error) {
      this.logger.error('Error closing RabbitMQ connection:', error);
    }
  }
}
