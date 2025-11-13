import { IsNotEmpty, IsObject, IsOptional, IsString, IsUrl } from "class-validator";

export class UserDataDto {
    @IsString()
    @IsNotEmpty()
    name: string;

    @IsUrl()
    @IsNotEmpty()
    link: string;

    @IsObject()
    @IsOptional()
    meta?: Record<string, any>;    
}