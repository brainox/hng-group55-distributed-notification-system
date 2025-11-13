import { Injectable, CanActivate, ExecutionContext, UnauthorizedException } from '@nestjs/common';
import { Observable } from 'rxjs';
import axios from 'axios';

@Injectable()
export class AuthGuard implements CanActivate {
  private readonly userServiceUrl = process.env.USER_SERVICE_URL || 'http://localhost:8000';

  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request = context.switchToHttp().getRequest();
    const token = this.extractTokenFromHeader(request);

    if (!token) {
      throw new UnauthorizedException('No authentication token provided');
    }

    return this.validateToken(token, request);
  }

  private extractTokenFromHeader(request: any): string | undefined {
    const authHeader = request.headers.authorization;
    if (!authHeader) {
      return undefined;
    }

    const [type, token] = authHeader.split(' ');
    return type === 'Bearer' ? token : undefined;
  }

  private async validateToken(token: string, request: any): Promise<boolean> {
    try {
      // Call User Service to validate token
      const response = await axios.get(
        `${this.userServiceUrl}/v1/users/me`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
          timeout: 5000,
        }
      );

      // Attach user info to request for later use
      request.user = response.data.data;
      return true;

    } catch (error) {
      if (error.response?.status === 401) {
        throw new UnauthorizedException('Invalid or expired token');
      }
      throw new UnauthorizedException('Token validation failed');
    }
  }
}