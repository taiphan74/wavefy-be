import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { User } from './user.entity';
import { CreateUserDto, UpdateUserDto } from './user.dto';
import { UserServiceAbstract } from './user.service.interface';

@Injectable()
export class UserService extends UserServiceAbstract {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
  ) {
    super();
  }

  async findAll(): Promise<User[]> {
    return this.userRepository.find();
  }

  async findOne(id: string): Promise<User | null> {
    return this.userRepository.findOneBy({ id });
  }

  async create(createUserDto: CreateUserDto): Promise<User> {
    const { password, ...rest } = createUserDto;
    const password_hash = await bcrypt.hash(password, 10);
    const user = this.userRepository.create({ ...rest, password_hash });
    return this.userRepository.save(user);
  }

  async update(id: string, updateUserDto: UpdateUserDto): Promise<User | null> {
    const { password, ...rest } = updateUserDto;
    if (password) {
      const password_hash = await bcrypt.hash(password, 10);
      await this.userRepository.update(id, { ...rest, password_hash });
    } else {
      await this.userRepository.update(id, rest);
    }
    return this.findOne(id);
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
