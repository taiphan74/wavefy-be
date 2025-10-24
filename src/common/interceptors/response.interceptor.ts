import {
  Injectable,
  NestInterceptor,
  ExecutionContext,
  CallHandler,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Response } from 'express';
import { HandlerResponse } from '../interfaces/handler-response.interface';

@Injectable()
export class ResponseInterceptor implements NestInterceptor {
  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    return next.handle().pipe(
      map((handlerResponse: HandlerResponse) => {
        const httpResponse = context.switchToHttp().getResponse<Response>();
        const httpStatusCode = httpResponse.statusCode;

        const responseMessage = handlerResponse?.message || 'Success';
        const responsePayload =
          handlerResponse?.data !== undefined
            ? handlerResponse.data
            : handlerResponse;
        const responseMeta = handlerResponse?.meta;

        const formattedResult: {
          status: number;
          message: string;
          data: unknown;
          meta?: unknown;
        } = {
          status: httpStatusCode,
          message: responseMessage,
          data: responsePayload,
        };

        if (responseMeta !== undefined) {
          formattedResult.meta = responseMeta;
        }

        return formattedResult;
      }),
    );
  }
}
