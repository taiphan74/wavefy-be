import {
  Injectable,
  ConflictException,
  UnauthorizedException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { User, SignupMethod } from '../users/user.entity';
import { UserService } from '../users/user.service';
import { RegisterDto, LoginDto } from './auth.dto';

@Injectable()
export class AuthService {
  constructor(
    private userService: UserService,
    @InjectRepository(User)
    private userRepository: Repository<User>,
  ) {}

  async register(registerDto: RegisterDto): Promise<User> {
    const { username, email, password } = registerDto;

    // Check if user exists
    const exists = await this.userService.existsByEmailOrUsername(
      email,
      username,
    );
    if (exists) {
      throw new ConflictException('User already exists');
    }

    // Create user via userService
    const createUserDto = {
      username,
      email,
      password,
      first_name: '',
      last_name: '',
      signup_method: SignupMethod.LOCAL,
    };
    return this.userService.create(createUserDto);
  }

  async login(loginDto: LoginDto): Promise<User> {
    const { email, password } = loginDto;

    // Find user
    const user = await this.userRepository.findOne({ where: { email } });
    if (!user) {
      throw new UnauthorizedException('Invalid credentials');
    }

    // Check password
    const isPasswordValid = await bcrypt.compare(password, user.password_hash);
    if (!isPasswordValid) {
      throw new UnauthorizedException('Invalid credentials');
    }

    return user;
  }
}
