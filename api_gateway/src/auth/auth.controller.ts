import { Controller, Post, Get, Body, Headers, HttpException, UseGuards } from '@nestjs/common';
import axios from 'axios';
import { AuthGuard } from './auth.guard';

@Controller()
export class AuthController {
  private readonly userServiceUrl = process.env.USER_SERVICE_URL || 'http://localhost:3001';

  // Register - No auth required
  @Post('users/register')
  async register(@Body() body: any) {
    try {
      const response = await axios.post(
        `${this.userServiceUrl}/v1/users/register`,
        body,
        { timeout: 10000 }
      );
      return response.data;
    } catch (error) {
      throw new HttpException(
        error.response?.data || 'Registration failed',
        error.response?.status || 500
      );
    }
  }

  // Login - No auth required
  @Post('users/login')
  async login(@Body() body: any) {
    try {
      const response = await axios.post(
        `${this.userServiceUrl}/v1/users/login`,
        body,
        { timeout: 10000 }
      );
      return response.data;
    } catch (error) {
      throw new HttpException(
        error.response?.data || 'Login failed',
        error.response?.status || 500
      );
    }
  }

  // Get current user - Auth required
  @Get('users/me')
  @UseGuards(AuthGuard)
  async getCurrentUser(@Headers('authorization') auth: string) {
    try {
      const response = await axios.get(
        `${this.userServiceUrl}/v1/users/me`,
        {
          headers: { Authorization: auth },
          timeout: 5000,
        }
      );
      return response.data;
    } catch (error) {
      throw new HttpException(
        error.response?.data || 'Failed to get user',
        error.response?.status || 500
      );
    }
  }
}