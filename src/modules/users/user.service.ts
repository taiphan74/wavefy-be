// modules/users/user.service.ts
import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { User } from './user.entity';
import { CreateUserDto, UpdateUserDto } from './user.dto';
import { UserServiceAbstract } from './user.service.interface';
import { UserResponseDto } from './user-response.dto';

@Injectable()
export class UserService extends UserServiceAbstract {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
  ) {
    super();
  }

  async findAll(): Promise<UserResponseDto[]> {
    const users = await this.userRepository.find();
    return users.map((user) => new UserResponseDto(user));
  }

  async findOne(id: string): Promise<UserResponseDto | null> {
    const user = await this.userRepository.findOneBy({ id });
    return user ? new UserResponseDto(user) : null;
  }

  async create(createUserDto: CreateUserDto): Promise<UserResponseDto> {
    const { password, ...rest } = createUserDto;
    const password_hash = await bcrypt.hash(password, 10);
    const user = this.userRepository.create({ ...rest, password_hash });
    const savedUser = await this.userRepository.save(user);
    return new UserResponseDto(savedUser);
  }

  async update(id: string, updateUserDto: UpdateUserDto): Promise<UserResponseDto | null> {
    const { password, ...rest } = updateUserDto;
    if (password) {
      const password_hash = await bcrypt.hash(password, 10);
      await this.userRepository.update(id, { ...rest, password_hash });
    } else {
      await this.userRepository.update(id, rest);
    }
    const updatedUser = await this.findOne(id);
    return updatedUser ? new UserResponseDto(updatedUser as any) : null;
  }

  async remove(id: string): Promise<void> {
    await this.userRepository.delete(id);
  }

  async existsByEmailOrUsername(
    email: string,
    username: string,
  ): Promise<boolean> {
    const user = await this.userRepository.findOne({
      where: [{ email }, { username }],
    });
    return !!user;
  }
}
