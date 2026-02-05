import fastify from 'fastify';
import cors from '@fastify/cors';
import rateLimit from '@fastify/rate-limit';
import Redis from 'ioredis';
import amqp from 'amqplib';
import { v4 as uuid } from 'uuid';

// Types
interface Event {
  eventId: string;
  eventType: string;
  timestamp: string;
  traceId: string;
  payload: any;
}

// Config from env
const config = {
  port: parseInt(process.env.PORT || '8080'),
  host: process.env.HOST || '0.0.0.0',
  corsOrigins: process.env.CORS_ORIGINS?.split(',') || ['http://localhost:3000'],
  redisUrl: process.env.REDIS_URL || 'redis://localhost:6379',
  rabbitmqUrl: process.env.RABBITMQ_URL || 'amqp://localhost',
  cacheTtl: parseInt(process.env.CACHE_TTL || '300'),
  services: {
    catalog: process.env.CATALOG_BASE_URL || 'http://localhost:8081',
    users: process.env.USERS_BASE_URL || 'http://localhost:8082',
    library: process.env.LIBRARY_BASE_URL || 'http://localhost:8083',
    playback: process.env.PLAYBACK_BASE_URL || 'http://localhost:8085',
  }
};

// Initialize services
const app = fastify({ 
  logger: { level: process.env.LOG_LEVEL || 'info' }
});

let redis: Redis;
let rabbitMQ: amqp.Connection;
let channel: amqp.Channel;

// Initialize Redis
async function initRedis() {
  redis = new Redis(config.redisUrl);
  redis.on('error', (err) => {
    app.log.error('Redis connection error:', err);
  });
}

// Initialize RabbitMQ
async function initRabbitMQ() {
  try {
    rabbitMQ = await amqp.connect(config.rabbitmqUrl);
    channel = await rabbitMQ.createChannel();
    
    // Ensure exchanges exist
    await channel.assertExchange('nudex.events', 'topic', { durable: true });
    
    app.log.info('RabbitMQ connected successfully');
  } catch (error) {
    app.log.error('RabbitMQ connection error:', error);
  }
}

// Publish event
async function publishEvent(event: Event) {
  if (!channel) return;
  
  try {
    await channel.publish(
      'nudex.events',
      event.eventType,
      Buffer.from(JSON.stringify(event))
    );
  } catch (error) {
    app.log.error('Error publishing event:', error);
  }
}

// Register plugins
app.register(cors, {
  origin: config.corsOrigins,
  credentials: true
});

app.register(rateLimit, {
  max: parseInt(process.env.RATE_LIMIT_MAX || '100'),
  timeWindow: parseInt(process.env.RATE_LIMIT_WINDOW || '60000')
});

// Middleware for tracing
app.addHook('onRequest', async (request, reply) => {
  request.traceId = request.headers['x-trace-id'] as string || uuid();
  reply.header('x-trace-id', request.traceId);
});

// Health check
app.get('/health', async (request, reply) => {
  return { 
    status: 'healthy',
    timestamp: new Date().toISOString(),
    version: '1.0.0',
    services: {
      redis: redis.status,
      rabbitmq: channel ? 'connected' : 'disconnected'
    }
  };
});

// === API ROUTES ===

// Home feed
app.get('/api/feed/home', async (request, reply) => {
  const cacheKey = 'feed:home';
  
  try {
    // Check cache first
    const cached = await redis.get(cacheKey);
    if (cached) {
      return JSON.parse(cached);
    }
    
    // Fetch from catalog service
    const response = await fetch(`${config.services.catalog}/videos?limit=20&random=true`);
    const videos = await response.json();
    
    // Cache result
    await redis.setex(cacheKey, config.cacheTtl, JSON.stringify(videos));
    
    return videos;
  } catch (error) {
    app.log.error('Error fetching home feed:', error);
    reply.code(500).send({ error: 'Internal server error' });
  }
});

// Video detail
app.get('/api/videos/:id', async (request, reply) => {
  const { id } = request.params as { id: string };
  const cacheKey = `video:${id}`;
  
  try {
    // Check cache
    const cached = await redis.get(cacheKey);
    if (cached) {
      return JSON.parse(cached);
    }
    
    // Fetch from catalog service
    const response = await fetch(`${config.services.catalog}/videos/${id}`);
    if (!response.ok) {
      return reply.code(404).send({ error: 'Video not found' });
    }
    
    const video = await response.json();
    
    // Cache result
    await redis.setex(cacheKey, config.cacheTtl, JSON.stringify(video));
    
    return video;
  } catch (error) {
    app.log.error('Error fetching video:', error);
    reply.code(500).send({ error: 'Internal server error' });
  }
});

// Search
app.get('/api/search', async (request, reply) => {
  const { q } = request.query as { q: string };
  
  if (!q) {
    return reply.code(400).send({ error: 'Query parameter required' });
  }
  
  const cacheKey = `search:${encodeURIComponent(q)}`;
  
  try {
    // Check cache
    const cached = await redis.get(cacheKey);
    if (cached) {
      return JSON.parse(cached);
    }
    
    // Search in catalog service
    const response = await fetch(`${config.services.catalog}/videos/search?q=${encodeURIComponent(q)}`);
    const results = await response.json();
    
    // Cache for shorter time (search results change frequently)
    await redis.setex(cacheKey, 60, JSON.stringify(results));
    
    return results;
  } catch (error) {
    app.log.error('Error searching:', error);
    reply.code(500).send({ error: 'Internal server error' });
  }
});

// Mock auth login
app.post('/api/auth/login', async (request, reply) => {
  const { email, password } = request.body as { email: string; password: string };
  
  // Mock authentication - in real app, proxy to users service
  if (email && password) {
    const user = {
      id: uuid(),
      email,
      name: 'NUDEX User',
      avatar: null
    };
    
    const token = 'mock_jwt_token_' + uuid();
    
    // Publish user login event
    await publishEvent({
      eventId: uuid(),
      eventType: 'user.login',
      timestamp: new Date().toISOString(),
      traceId: request.traceId,
      payload: { userId: user.id, email }
    });
    
    return {
      user,
      token,
      expiresIn: '24h'
    };
  }
  
  reply.code(401).send({ error: 'Invalid credentials' });
});

// Favorites (proxy to library service)
app.get('/api/favorites', async (request, reply) => {
  const userId = request.headers['x-user-id'] as string || 'anonymous';
  
  try {
    const response = await fetch(`${config.services.library}/favorites`, {
      headers: { 'x-user-id': userId }
    });
    
    if (!response.ok) {
      return reply.code(response.status).send({ error: 'Failed to fetch favorites' });
    }
    
    return await response.json();
  } catch (error) {
    app.log.error('Error fetching favorites:', error);
    reply.code(500).send({ error: 'Internal server error' });
  }
});

// Add/Remove favorite
app.post('/api/favorites/:videoId', async (request, reply) => {
  const { videoId } = request.params as { videoId: string };
  const userId = request.headers['x-user-id'] as string || 'anonymous';
  
  try {
    const response = await fetch(`${config.services.library}/favorites`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'x-user-id': userId
      },
      body: JSON.stringify({ videoId })
    });
    
    if (!response.ok) {
      return reply.code(response.status).send({ error: 'Failed to update favorite' });
    }
    
    const result = await response.json();
    
    // Publish favorite event
    await publishEvent({
      eventId: uuid(),
      eventType: 'library.favorited',
      timestamp: new Date().toISOString(),
      traceId: request.traceId,
      payload: { userId, videoId, action: result.action }
    });
    
    return result;
  } catch (error) {
    app.log.error('Error updating favorite:', error);
    reply.code(500).send({ error: 'Internal server error' });
  }
});

// Start server
async function start() {
  try {
    await initRedis();
    await initRabbitMQ();
    
    await app.listen({ 
      port: config.port, 
      host: config.host 
    });
    
    app.log.info(`🚀 NUDEX API Gateway running on http://${config.host}:${config.port}`);
  } catch (error) {
    app.log.error('Error starting server:', error);
    process.exit(1);
  }
}

// Graceful shutdown
process.on('SIGTERM', async () => {
  app.log.info('SIGTERM received, shutting down gracefully');
  
  if (redis) redis.disconnect();
  if (rabbitMQ) await rabbitMQ.close();
  
  await app.close();
  process.exit(0);
});

start();