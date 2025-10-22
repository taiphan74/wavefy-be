// Đây là file chính của ứng dụng NestJS Wavefy
// Import các module cần thiết từ NestJS và các thư viện khác
import { NestFactory } from '@nestjs/core';
import { ValidationPipe } from '@nestjs/common';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import helmet from 'helmet';
import { AppModule } from './app.module';
import { ResponseInterceptor } from './common/interceptors/response.interceptor';

// Hàm khởi động ứng dụng bất đồng bộ
async function bootstrap() {
  // Tạo instance của ứng dụng NestJS từ AppModule
  const app = await NestFactory.create(AppModule);

  // Sử dụng helmet để bảo mật ứng dụng
  app.use(helmet());
  // Áp dụng ValidationPipe toàn cục để validate dữ liệu đầu vào
  app.useGlobalPipes(new ValidationPipe());
  // Đăng ký ResponseInterceptor toàn cục để chuẩn hóa response
  app.useGlobalInterceptors(new ResponseInterceptor());

  // Nếu đang ở môi trường development, cấu hình Swagger để tạo tài liệu API
  if (process.env.NODE_ENV === 'development') {
    const config = new DocumentBuilder()
      .setTitle('Wavefy API')
      .setDescription('The Wavefy API description')
      .setVersion('1.0')
      .addTag('wavefy')
      .build();
    const document = SwaggerModule.createDocument(app, config);
    SwaggerModule.setup('api', app, document);
  }

  // Lắng nghe trên port được chỉ định trong env hoặc mặc định 3000
  await app.listen(process.env.PORT ?? 3000);
  // Lấy URL của ứng dụng và in ra console
  const url = await app.getUrl();
  console.log(`Backend: ${url}`);
  // Nếu ở development, in thêm URL của tài liệu API
  if (process.env.NODE_ENV === 'development') {
    console.log(`API docs: ${url}/api`);
  }
}
// Gọi hàm bootstrap để khởi động ứng dụng
bootstrap();
