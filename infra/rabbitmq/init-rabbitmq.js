import amqp from 'amqplib';

const RABBITMQ_URL = process.env.RABBITMQ_URL || 'amqp://rabbitmq:5672';

async function setupQueues() {
  try {
    console.log('Connecting to RabbitMQ...');
    const connection = await amqp.connect(RABBITMQ_URL);
    const channel = await connection.createChannel();
    
    console.log('Setting up exchanges...');
    
    // Declare dead letter exchange
    await channel.assertExchange('dlx.notifications', 'direct', { durable: true });
    console.log('Created dead letter exchange: dlx.notifications');
    
    // Declare main exchange
    await channel.assertExchange('notifications.direct', 'direct', { durable: true });
    console.log('Created main exchange: notifications.direct');
    
    console.log('Setting up queues...');
    
    // Declare dead letter queue
    await channel.assertQueue('failed.queue', { durable: true });
    await channel.bindQueue('failed.queue', 'dlx.notifications', 'failed');
    console.log('Created queue: failed.queue');
    
    // Email queue with DLX
    await channel.assertQueue('email.queue', {
      durable: true,
      deadLetterExchange: 'dlx.notifications',
      deadLetterRoutingKey: 'failed',
      messageTtl: 300000, // 5 minutes
    });
    await channel.bindQueue('email.queue', 'notifications.direct', 'email');
    console.log('Created queue: email.queue');
    
    // Push queue with DLX
    await channel.assertQueue('push.queue', {
      durable: true,
      deadLetterExchange: 'dlx.notifications',
      deadLetterRoutingKey: 'failed',
      messageTtl: 300000,
    });
    await channel.bindQueue('push.queue', 'notifications.direct', 'push');
    console.log('Created queue: push.queue');
    
    console.log('\n RabbitMQ setup complete!');
    console.log('\n Queue Summary:');
    console.log('  Exchange: notifications.direct');
    console.log('  - email.queue → Email Service');
    console.log('  - push.queue → Push Service');
    console.log('  - failed.queue → Dead Letter Queue');
    console.log('\n Management UI: http://localhost:15672');
    console.log('   Credentials: guest / guest\n');
    
    await channel.close();
    await connection.close();
    process.exit(0);
  } catch (error) {
    console.error('Error setting up RabbitMQ:', error);
    process.exit(1);
  }
}

setupQueues();