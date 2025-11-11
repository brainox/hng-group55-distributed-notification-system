import { 
    IsEnum, 
    IsNotEmpty, 
    IsObject, 
    IsOptional, 
    IsString, 
    IsUUID,
    IsInt, 
    Min, 
    Max,
    ValidateNested 
} from 'class-validator';
import { Type } from 'class-transformer';
import { UserDataDto } from './user-data.dto';


export enum NotificationType {
  EMAIL = 'email',
  PUSH = 'push',
}

export class CreateNotificationDto {
    @IsEnum(['email', 'push'])
    @IsNotEmpty()
    notification_type: 'email' | 'push';

    @IsUUID()
    @IsNotEmpty()
    user_id: string;

    @IsString()
    @IsNotEmpty()
    template_code: string;


    @ValidateNested()
    @Type(() => UserDataDto)
    variables: UserDataDto;

    // @IsObject()
    // @IsNotEmpty()
    // // variables: UserDataDto;
    // variables: Record<string, string>;

    @IsString()
    @IsNotEmpty()
    request_id: string;

    @IsInt()
    @IsOptional()
    @Min(1)
    @Max(3)
    priority?: number; // 1 = high, 2 = normal, 3 = low

    @IsOptional()
    @IsObject()
    metadata?: Record<string, any>;
}
