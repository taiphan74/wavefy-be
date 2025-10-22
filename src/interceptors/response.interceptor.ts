import { Injectable, NestInterceptor, ExecutionContext, CallHandler } from '@nestjs/common';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable()
export class ResponseInterceptor implements NestInterceptor {
    intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
        return next.handle().pipe(
            map((handlerResponse) => {
                const httpResponse = context.switchToHttp().getResponse();
                const httpStatusCode = httpResponse.statusCode;

                const responseMessage = handlerResponse?.message || 'Success';
                const responsePayload =
                    handlerResponse?.data !== undefined ? handlerResponse.data : handlerResponse;
                const responseMeta = handlerResponse?.meta;

                const formattedResult: any = {
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