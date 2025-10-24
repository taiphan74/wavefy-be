import { IsString } from 'class-validator';
import { IsEmail } from 'class-validator';
import { IsBoolean } from 'class-validator';
import { IsEnum } from 'class-validator';
import { IsOptional } from 'class-validator';
import { SignupMethod } from './user.entity';

export class CreateUserDto {
  @IsString()
  username: string;

  @IsString()
  first_name: string;

  @IsString()
  last_name: string;

  @IsEmail()
  email: string;

  @IsOptional()
  @IsBoolean()
  email_verified?: boolean;

  @IsString()
  password: string;

  @IsEnum(SignupMethod)
  signup_method: SignupMethod;
}

export class UpdateUserDto {
  @IsOptional()
  @IsString()
  username?: string;

  @IsOptional()
  @IsString()
  first_name?: string;

  @IsOptional()
  @IsString()
  last_name?: string;

  @IsOptional()
  @IsEmail()
  email?: string;

  @IsOptional()
  @IsBoolean()
  email_verified?: boolean;

  @IsOptional()
  @IsString()
  password?: string;

  @IsOptional()
  @IsEnum(SignupMethod)
  signup_method?: SignupMethod;
}
