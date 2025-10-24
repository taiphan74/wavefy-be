import { User } from './user.entity';

export class UserResponseDto {
  id: string;
  username: string;
  first_name: string;
  last_name: string;
  email: string;
  email_verified: boolean;
  signup_method: string;
  created_at: Date;
  updated_at: Date;

  constructor(user: User) {
    this.id = user.id;
    this.username = user.username;
    this.first_name = user.first_name;
    this.last_name = user.last_name;
    this.email = user.email;
    this.email_verified = user.email_verified;
    this.signup_method = user.signup_method;
    this.created_at = user.created_at;
    this.updated_at = user.updated_at;
  }
}
