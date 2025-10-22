import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  getHello(): object {
    return {
      name: 'Wavefy',
      version: '1.0.0',
      description: 'Wavefy Backend API',
      environment: process.env.NODE_ENV || 'development',
    };
  }
}
