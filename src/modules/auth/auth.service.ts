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
import { UserResponseDto } from '../users/user-response.dto';

@Injectable()
export class AuthService {
  constructor(
    private userService: UserService,
    @InjectRepository(User)
    private userRepository: Repository<User>,
  ) {}

  async register(registerDto: RegisterDto): Promise<UserResponseDto> {
    const { username, email, password } = registerDto;

    const exists = await this.userService.existsByEmailOrUsername(
      email,
      username,
    );
    if (exists) {
      throw new ConflictException('User already exists');
    }

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

  async login(loginDto: LoginDto): Promise<UserResponseDto> {
    const { email, password } = loginDto;

    const user = await this.userRepository
      .createQueryBuilder('user')
      .addSelect('user.password_hash')
      .where('user.email = :email', { email })
      .getOne();

    if (!user) {
      throw new UnauthorizedException('Invalid credentials');
    }

    const isPasswordValid = await bcrypt.compare(password, user.password_hash);
    if (!isPasswordValid) {
      throw new UnauthorizedException('Invalid credentials');
    }

    return new UserResponseDto(user);
  }
}
