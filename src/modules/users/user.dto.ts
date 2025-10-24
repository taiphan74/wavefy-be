import { IsString, IsEmail, IsBoolean, IsEnum, IsOptional } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';
import { SignupMethod } from './user.entity';

export class CreateUserDto {
  @ApiProperty({ example: 'johndoe', description: 'Unique username of the user' })
  @IsString()
  username: string;

  @ApiProperty({ example: 'John', description: 'First name of the user' })
  @IsString()
  first_name: string;

  @ApiProperty({ example: 'Doe', description: 'Last name of the user' })
  @IsString()
  last_name: string;

  @ApiProperty({ example: 'john@example.com', description: 'Email address of the user' })
  @IsEmail()
  email: string;

  @ApiProperty({ example: false, description: 'Indicates whether the email has been verified', required: false })
  @IsOptional()
  @IsBoolean()
  email_verified?: boolean;

  @ApiProperty({ example: 'strongPassword123', description: 'User password' })
  @IsString()
  password: string;

  @ApiProperty({ enum: SignupMethod, example: SignupMethod.LOCAL, description: 'Signup method used by the user' })
  @IsEnum(SignupMethod)
  signup_method: SignupMethod;
}

export class UpdateUserDto {
  @ApiProperty({ example: 'newusername', description: 'Updated username', required: false })
  @IsOptional()
  @IsString()
  username?: string;

  @ApiProperty({ example: 'John', description: 'Updated first name', required: false })
  @IsOptional()
  @IsString()
  first_name?: string;

  @ApiProperty({ example: 'Doe', description: 'Updated last name', required: false })
  @IsOptional()
  @IsString()
  last_name?: string;

  @ApiProperty({ example: 'newemail@example.com', description: 'Updated email address', required: false })
  @IsOptional()
  @IsEmail()
  email?: string;

  @ApiProperty({ example: true, description: 'Updated email verification status', required: false })
  @IsOptional()
  @IsBoolean()
  email_verified?: boolean;

  @ApiProperty({ example: 'newPassword123', description: 'Updated user password', required: false })
  @IsOptional()
  @IsString()
  password?: string;

  @ApiProperty({ enum: SignupMethod, example: SignupMethod.GOOGLE, description: 'Updated signup method', required: false })
  @IsOptional()
  @IsEnum(SignupMethod)
  signup_method?: SignupMethod;
}