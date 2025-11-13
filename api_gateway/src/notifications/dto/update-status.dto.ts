import { IsDateString, IsEnum, IsNotEmpty, IsOptional, IsString } from "class-validator";

export enum NotificationStatus {
    DELIVERED = 'delivered',
    PENDING = 'pending',
    FAILED = 'failed',
}

export class UpdateStatusDto {
    @IsString()
    @IsNotEmpty()
    notification_id: string;

    @IsEnum(NotificationStatus)
    @IsNotEmpty()
    status: NotificationStatus;

    @IsDateString()
    @IsOptional()
    timestamp?: string;

    @IsString()
    @IsOptional()
    error?: string;
}